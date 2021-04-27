// Copyright 2020 The Operator-SDK Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package metrics

import (
	"github.com/prometheus/client_golang/prometheus"

	sdkVersion "github.com/stackrox/rox/operator/pkg/operator-sdk/version"
)

const (
	subsystem = "helm_operator"
)

var (
	buildInfo = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Subsystem: subsystem,
			Name:      "build_info",
			Help:      "Build information for the helm-operator binary",
			ConstLabels: map[string]string{
				"commit":  sdkVersion.GitCommit,
				"version": sdkVersion.Version,
			},
		},
	)
)

func RegisterBuildInfo(r prometheus.Registerer) {
	buildInfo.Set(1)
	r.MustRegister(buildInfo)
}
