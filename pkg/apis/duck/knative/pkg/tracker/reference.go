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

// Package tracker contains a partial, vendored copy of knative.dev/pkg/tracker
// (Apache License 2.0, The Knative Authors): only the Reference type used by the
// Camel K SinkBinding duck type, with an identical serialized form.
//
// +kubebuilder:object:generate=true
package tracker

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Reference is modeled after corev1.ObjectReference, but omits fields
// unsupported by the tracker, and permits us to extend things in
// divergent ways.
type Reference struct {
	// API version of the referent.
	// +optional
	APIVersion string `json:"apiVersion,omitempty"`

	// Kind of the referent.
	// +optional
	Kind string `json:"kind,omitempty"`

	// Namespace of the referent.
	// +optional
	Namespace string `json:"namespace,omitempty"`

	// Name of the referent.
	// Mutually exclusive with Selector.
	// +optional
	Name string `json:"name,omitempty"`

	// Selector of the referents.
	// Mutually exclusive with Name.
	// +optional
	Selector *metav1.LabelSelector `json:"selector,omitempty"`
}
