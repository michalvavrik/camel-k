/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// The types in this file are a partial copy of the Knative Eventing SinkBinding
// API (knative.dev/eventing/pkg/apis/sources/v1), Apache License 2.0 by The
// Knative Authors. Only the subset Camel K produces is retained, with the field
// layout kept identical to preserve serialization.

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	duckv1 "github.com/apache/camel-k/v2/pkg/apis/duck/knative/pkg/apis/duck/v1"
)

// +kubebuilder:object:root=true

// SinkBinding describes a Binding that is also a Source. The `sink` (from the
// Source duck) is resolved to a URL and then projected into the `subject` by
// augmenting the runtime contract of the referenced containers.
type SinkBinding struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SinkBindingSpec   `json:"spec"`
	Status SinkBindingStatus `json:"status"`
}

// SinkBindingSpec holds the desired state of the SinkBinding (from the client).
type SinkBindingSpec struct {
	// inherits duck/v1 SourceSpec, which currently provides:
	// * Sink - a reference to an object that will resolve to a uri to use as the sink.
	// * CloudEventOverrides - defines overrides to control the output format.
	duckv1.SourceSpec `json:",inline"`

	// inherits duck/v1 BindingSpec, which currently provides:
	// * Subject - references the resource(s) whose "runtime contract" should be
	//   augmented by Binding implementations.
	duckv1.BindingSpec `json:",inline"`
}

// SinkBindingStatus communicates the observed state of the SinkBinding (from the controller).
type SinkBindingStatus struct {
	// inherits duck/v1 Status, which currently provides:
	// * ObservedGeneration - the 'Generation' of the SinkBinding that was last processed by the controller.
	// * Conditions - the latest available observations of a resource's current state.
	duckv1.Status `json:",inline"`
}

// +kubebuilder:object:root=true

// SinkBindingList contains a list of SinkBinding.
type SinkBindingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []SinkBinding `json:"items"`
}
