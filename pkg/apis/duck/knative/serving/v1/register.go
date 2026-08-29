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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"knative.dev/pkg/apis"
)

const (
	// KnativeServingGroup is the API group for Knative Serving resources.
	KnativeServingGroup = "serving.knative.dev"
	// KnativeServingVersion is the API version supported by these duck types.
	KnativeServingVersion = "v1"
	// GroupName is the API group for Knative Serving resources (alias kept for
	// parity with the upstream knative.dev/serving/pkg/apis/serving package).
	GroupName = KnativeServingGroup
)

// ServiceConditionReady is set when the service is configured and has available
// backends ready to receive traffic. It matches the upstream Knative constant.
const ServiceConditionReady = apis.ConditionReady

var (
	// SchemeGroupVersion is the group version used to register these objects.
	SchemeGroupVersion = schema.GroupVersion{Group: KnativeServingGroup, Version: KnativeServingVersion}
	// SchemeBuilder builds a scheme with the Knative Serving duck types.
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	// AddToScheme registers the Knative Serving duck types with a scheme.
	AddToScheme = SchemeBuilder.AddToScheme
)

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&Service{},
		&ServiceList{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)

	return nil
}
