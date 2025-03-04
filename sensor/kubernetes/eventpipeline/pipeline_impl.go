package eventpipeline

import (
	"github.com/stackrox/rox/generated/internalapi/central"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/centralsensor"
	"github.com/stackrox/rox/pkg/concurrency"
	"github.com/stackrox/rox/pkg/env"
	"github.com/stackrox/rox/pkg/logging"
	"github.com/stackrox/rox/sensor/common"
	"github.com/stackrox/rox/sensor/common/detector"
	"github.com/stackrox/rox/sensor/common/reprocessor"
	"github.com/stackrox/rox/sensor/common/store/resolver"
	"github.com/stackrox/rox/sensor/kubernetes/eventpipeline/component"
)

var (
	log = logging.LoggerForModule()
)

type eventPipeline struct {
	output      component.OutputQueue
	resolver    component.Resolver
	listener    component.PipelineComponent
	detector    detector.Detector
	reprocessor reprocessor.Handler

	eventsC chan *central.MsgFromSensor
	stopSig concurrency.Signal
}

// Capabilities implements common.SensorComponent
func (*eventPipeline) Capabilities() []centralsensor.SensorCapability {
	return nil
}

// ProcessMessage implements common.SensorComponent
func (p *eventPipeline) ProcessMessage(msg *central.MsgToSensor) error {
	switch {
	case msg.GetPolicySync() != nil:
		return p.processPolicySync(msg.GetPolicySync())
	case msg.GetReassessPolicies() != nil:
		return p.processReassessPolicies()
	case msg.GetUpdatedImage() != nil:
		return p.processUpdatedImage(msg.GetUpdatedImage())
	case msg.GetReprocessDeployments() != nil:
		return p.processReprocessDeployments()
	case msg.GetReprocessDeployment() != nil:
		return p.processReprocessDeployment(msg.GetReprocessDeployment())
	case msg.GetInvalidateImageCache() != nil:
		return p.processInvalidateImageCache(msg.GetInvalidateImageCache())
	}
	return nil
}

// ResponsesC implements common.SensorComponent
func (p *eventPipeline) ResponsesC() <-chan *central.MsgFromSensor {
	return p.eventsC
}

// Start implements common.SensorComponent
func (p *eventPipeline) Start() error {
	// The order is important here, we need to start the components
	// that receive messages from other components first
	if err := p.output.Start(); err != nil {
		return err
	}

	if env.ResyncDisabled.BooleanSetting() {
		if err := p.resolver.Start(); err != nil {
			return err
		}
	}

	if err := p.listener.Start(); err != nil {
		return err
	}

	go p.forwardMessages()
	return nil
}

// Stop implements common.SensorComponent
func (p *eventPipeline) Stop(_ error) {
	defer close(p.eventsC)
	// The order is important here, we need to stop the components
	// that send messages to other components first
	p.listener.Stop(nil)
	if env.ResyncDisabled.BooleanSetting() {
		p.resolver.Stop(nil)
	}
	p.output.Stop(nil)
	p.stopSig.Signal()
}

func (p *eventPipeline) Notify(common.SensorComponentEvent) {}

// forwardMessages from listener component to responses channel
func (p *eventPipeline) forwardMessages() {
	for {
		select {
		case <-p.stopSig.Done():
			return
		case msg, more := <-p.output.ResponsesC():
			if !more {
				log.Error("Output component channel closed")
				return
			}
			p.eventsC <- msg
		}
	}
}

func (p *eventPipeline) processPolicySync(sync *central.PolicySync) error {
	log.Debug("PolicySync message received from central")
	return p.detector.ProcessPolicySync(sync)
}

func (p *eventPipeline) processReassessPolicies() error {
	log.Debug("ReassessPolicies message received from central")
	if err := p.detector.ProcessReassessPolicies(); err != nil {
		return err
	}
	if env.ResyncDisabled.BooleanSetting() {
		message := component.NewEvent()
		// TODO(ROX-14310): Add WithSkipResolving to the DeploymentReference (Revert: https://github.com/stackrox/stackrox/pull/5551)
		message.AddDeploymentReference(resolver.ResolveAllDeployments(),
			component.WithForceDetection())
		p.resolver.Send(message)
	}
	return nil
}

func (p *eventPipeline) processReprocessDeployments() error {
	log.Debug("ReprocessDeployments message received from central")
	if err := p.detector.ProcessReprocessDeployments(); err != nil {
		return err
	}
	if env.ResyncDisabled.BooleanSetting() {
		message := component.NewEvent()
		// TODO(ROX-14310): Add WithSkipResolving to the DeploymentReference (Revert: https://github.com/stackrox/stackrox/pull/5551)
		message.AddDeploymentReference(resolver.ResolveAllDeployments(),
			component.WithForceDetection())
		p.resolver.Send(message)
	}
	return nil
}

func (p *eventPipeline) processUpdatedImage(image *storage.Image) error {
	log.Debugf("UpdatedImage message received from central: image name: %s, number of components: %d", image.GetName().GetFullName(), image.GetComponents())
	if err := p.detector.ProcessUpdatedImage(image); err != nil {
		return err
	}
	if env.ResyncDisabled.BooleanSetting() {
		message := component.NewEvent()
		message.AddDeploymentReference(resolver.ResolveDeploymentsByImages(image),
			component.WithForceDetection(),
			component.WithSkipResolving())
		p.resolver.Send(message)
	}
	return nil
}

func (p *eventPipeline) processReprocessDeployment(req *central.ReprocessDeployment) error {
	log.Debug("ReprocessDeployment message received from central")
	if err := p.reprocessor.ProcessReprocessDeployments(req); err != nil {
		return err
	}
	if env.ResyncDisabled.BooleanSetting() {
		message := component.NewEvent()
		message.AddDeploymentReference(resolver.ResolveDeploymentIds(req.GetDeploymentIds()...),
			component.WithForceDetection(),
			component.WithSkipResolving())
		p.resolver.Send(message)
	}
	return nil
}

func (p *eventPipeline) processInvalidateImageCache(req *central.InvalidateImageCache) error {
	log.Debug("InvalidateImageCache message received from central")
	if err := p.reprocessor.ProcessInvalidateImageCache(req); err != nil {
		return err
	}
	if env.ResyncDisabled.BooleanSetting() {
		keys := make([]*storage.Image, len(req.GetImageKeys()))
		for i, image := range req.GetImageKeys() {
			keys[i] = &storage.Image{
				Id: image.GetImageId(),
				Name: &storage.ImageName{
					FullName: image.GetImageFullName(),
				},
			}
		}
		message := component.NewEvent()
		message.AddDeploymentReference(resolver.ResolveDeploymentsByImages(keys...),
			component.WithForceDetection(),
			component.WithSkipResolving())
		p.resolver.Send(message)
	}
	return nil
}
