// Copyright Istio Authors
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

package controller

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/discovery/v1"

	"istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pkg/fuzz"
	"istio.io/istio/pkg/network"
)

func FuzzKubeController(f *testing.F) {
	fuzz.Fuzz(f, func(fg fuzz.Helper) {
		networkID := network.ID("fakeNetwork")
		fco := fuzz.Struct[FakeControllerOptions](fg)
		fco.SkipRun = true
		controller, _ := NewFakeControllerWithOptions(fg.T(), fco)
		controller.network = networkID

		p := fuzz.Struct[*corev1.Pod](fg)
		controller.pods.onEvent(p, model.EventAdd)
		s := fuzz.Struct[*corev1.Service](fg)
		controller.onServiceEvent(s, model.EventAdd)
		if fco.Mode == EndpointsOnly {
			e := fuzz.Struct[*corev1.Endpoints](fg)
			controller.endpoints.onEvent(e, model.EventAdd)
		} else {
			e := fuzz.Struct[*v1.EndpointSlice](fg)
			controller.endpoints.onEvent(e, model.EventAdd)
		}
	})
}
