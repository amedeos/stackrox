package resources

import (
	"github.com/pkg/errors"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/set"
	"github.com/stackrox/rox/pkg/sync"
	"github.com/stackrox/rox/sensor/common/imagecacheutils"
	"github.com/stackrox/rox/sensor/common/selector"
	"github.com/stackrox/rox/sensor/common/store"
)

// DeploymentStore stores deployments.
type DeploymentStore struct {
	lock sync.RWMutex

	// Stores deployment IDs by namespaces.
	deploymentIDs map[string]map[string]struct{}
	// Stores deployments by IDs.
	deployments map[string]*deploymentWrap
}

// newDeploymentStore creates and returns a new deployment store.
func newDeploymentStore() *DeploymentStore {
	return &DeploymentStore{
		deploymentIDs: make(map[string]map[string]struct{}),
		deployments:   make(map[string]*deploymentWrap),
	}
}

func (ds *DeploymentStore) addOrUpdateDeployment(wrap *deploymentWrap) {
	ds.lock.Lock()
	defer ds.lock.Unlock()

	ds.addOrUpdateDeploymentNoLock(wrap)
}

func (ds *DeploymentStore) addOrUpdateDeploymentNoLock(wrap *deploymentWrap) {
	ids, ok := ds.deploymentIDs[wrap.GetNamespace()]
	if !ok {
		ids = make(map[string]struct{})
		ds.deploymentIDs[wrap.GetNamespace()] = ids
	}
	ids[wrap.GetId()] = struct{}{}

	ds.deployments[wrap.GetId()] = wrap
}

func (ds *DeploymentStore) removeDeployment(wrap *deploymentWrap) {
	ds.lock.Lock()
	defer ds.lock.Unlock()

	ids := ds.deploymentIDs[wrap.GetNamespace()]
	if ids == nil {
		return
	}
	delete(ids, wrap.GetId())
	delete(ds.deployments, wrap.GetId())
}

func (ds *DeploymentStore) getDeploymentsByIDs(_ string, idSet set.StringSet) []*deploymentWrap {
	ds.lock.RLock()
	defer ds.lock.RUnlock()

	deployments := make([]*deploymentWrap, 0, len(idSet))
	for id := range idSet {
		wrap := ds.deployments[id]
		if wrap != nil {
			deployments = append(deployments, wrap)
		}
	}
	return deployments
}

func (ds *DeploymentStore) getMatchingDeployments(namespace string, sel selector.Selector) (matching []*deploymentWrap) {
	ds.lock.RLock()
	defer ds.lock.RUnlock()

	ids := ds.deploymentIDs[namespace]
	if ids == nil {
		return
	}

	for id := range ids {
		wrap := ds.deployments[id]
		if wrap == nil {
			continue
		}

		if sel.Matches(selector.CreateLabelsWithLen(wrap.PodLabels)) {
			matching = append(matching, wrap)
		}
	}
	return
}

// CountDeploymentsForNamespace returns the number of deployments in a namespace
func (ds *DeploymentStore) CountDeploymentsForNamespace(namespace string) int {
	ds.lock.RLock()
	defer ds.lock.RUnlock()

	return len(ds.deploymentIDs[namespace])
}

// OnNamespaceDeleted reacts to a namespace deletion, deleting all deployments in this namespace from the store.
func (ds *DeploymentStore) OnNamespaceDeleted(namespace string) {
	ds.lock.Lock()
	defer ds.lock.Unlock()

	ids := ds.deploymentIDs[namespace]
	if ids == nil {
		return
	}

	for id := range ids {
		delete(ds.deployments, id)
	}
	delete(ds.deploymentIDs, namespace)
}

// GetAll returns all deployments.
func (ds *DeploymentStore) GetAll() []*storage.Deployment {
	ds.lock.RLock()
	defer ds.lock.RUnlock()

	var ret []*storage.Deployment
	for _, wrap := range ds.deployments {
		if wrap != nil {
			ret = append(ret, wrap.GetDeployment())
		}
	}
	return ret
}

// FindDeploymentIDsWithServiceAccount returns all deployment IDs in `namespace` that have ServiceAccountName matching `sa`.
func (ds *DeploymentStore) FindDeploymentIDsWithServiceAccount(namespace, sa string) []string {
	ds.lock.RLock()
	defer ds.lock.RUnlock()

	var match []string
	for id, wrap := range ds.deployments {
		if wrap.GetServiceAccount() == sa && wrap.GetNamespace() == namespace {
			match = append(match, id)
		}
	}
	return match
}

// FindDeploymentIDsByLabels returns a slice of deployments based on matching namespace and labels
func (ds *DeploymentStore) FindDeploymentIDsByLabels(namespace string, sel selector.Selector) (resIDs []string) {
	ds.lock.RLock()
	defer ds.lock.RUnlock()
	ids, found := ds.deploymentIDs[namespace]
	if !found || ids == nil {
		return
	}

	for id := range ids {
		wrap, found := ds.deployments[id]
		if !found || wrap == nil {
			continue
		}

		if sel.Matches(selector.CreateLabelsWithLen(wrap.GetPodLabels())) {
			resIDs = append(resIDs, id)
		}
	}
	return
}

func (ds *DeploymentStore) findDeploymentIDsByImageNoLock(image *storage.Image) set.Set[string] {
	ids := set.NewStringSet()
	for _, d := range ds.deployments {
		for _, c := range d.GetContainers() {
			if imagecacheutils.CompareImageCacheKey(c.GetImage(), image) {
				ids.Add(d.GetId())
				// The deployment id is already the set, we can break here
				break
			}
		}
	}
	return ids
}

// FindDeploymentIDsByImages returns a slice of deployment ids based on matching images
func (ds *DeploymentStore) FindDeploymentIDsByImages(images []*storage.Image) []string {
	ds.lock.RLock()
	defer ds.lock.RUnlock()
	ids := set.NewStringSet()
	for _, image := range images {
		ids = ids.Union(ds.findDeploymentIDsByImageNoLock(image))
	}
	return ids.AsSlice()
}

func (ds *DeploymentStore) getWrap(id string) *deploymentWrap {
	ds.lock.RLock()
	defer ds.lock.RUnlock()
	return ds.getWrapNoLock(id)
}

func (ds *DeploymentStore) getWrapNoLock(id string) *deploymentWrap {
	wrap := ds.deployments[id]
	return wrap
}

// Get returns deployment for supplied id.
func (ds *DeploymentStore) Get(id string) *storage.Deployment {
	wrap := ds.getWrap(id)
	return wrap.GetDeployment()
}

// GetBuiltDeployment returns a cloned deployment for supplied id and a flag if it is fully built.
func (ds *DeploymentStore) GetBuiltDeployment(id string) (*storage.Deployment, bool) {
	ds.lock.Lock()
	defer ds.lock.Unlock()
	wrap := ds.getWrapNoLock(id)
	if wrap == nil {
		return nil, false
	}
	return wrap.GetDeployment().Clone(), wrap.isBuilt
}

// BuildDeploymentWithDependencies creates storage.Deployment object using external object dependencies.
func (ds *DeploymentStore) BuildDeploymentWithDependencies(id string, dependencies store.Dependencies) (*storage.Deployment, error) {
	ds.lock.Lock()
	defer ds.lock.Unlock()

	// Get wrap with no lock since ds.lock.Lock() was already requested above
	wrap, found := ds.deployments[id]
	if !found || wrap == nil {
		return nil, errors.Errorf("deployment with ID %s doesn't exist in the internal deployment store", id)
	}

	wrap.updateServiceAccountPermissionLevel(dependencies.PermissionLevel)
	wrap.updatePortExposureSlice(dependencies.Exposures)
	if err := wrap.updateHash(); err != nil {
		return nil, err
	}

	// These properties are set when initially parsing a deployment/pod event as a deploymentWrap. Since secrets could
	// influence its values, we need to call this again with the same pods from the wrap. Inside this function we call
	// the registry store and update `IsClusterLocal` and `NotPullable` based on it. Meaning that if a pull secret was
	// updated, the value from this properties might need to be updated.
	wrap.populateDataFromPods(wrap.pods...)

	wrap.isBuilt = true
	ds.addOrUpdateDeploymentNoLock(wrap)
	return wrap.GetDeployment().Clone(), nil
}
