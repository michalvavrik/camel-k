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

// The types in this file are a partial copy of the Knative Serving API
// (knative.dev/serving/pkg/apis/serving/v1), which is licensed under the
// Apache License, Version 2.0 by The Knative Authors. Only the subset of
// fields used by Camel K is retained, and the field layout is kept identical
// to preserve the serialized form of the resources.

package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"knative.dev/pkg/apis"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

// +kubebuilder:object:root=true

// Service acts as a top-level container that manages a Route and Configuration
// which implement a network service.
type Service struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ServiceSpec   `json:"spec,omitempty"`
	Status ServiceStatus `json:"status,omitempty"`
}

// ServiceSpec represents the configuration for the Service object.
type ServiceSpec struct {
	// ConfigurationSpec holds the desired state of the Configuration (from the client).
	ConfigurationSpec `json:",inline"`

	// RouteSpec holds the desired state of the Route (from the client).
	RouteSpec `json:",inline"`
}

// ServiceStatus represents the Status stanza of the Service resource.
type ServiceStatus struct {
	duckv1.Status `json:",inline"`

	// ConfigurationStatusFields represents the current Configuration.
	ConfigurationStatusFields `json:",inline"`

	// RouteStatusFields represents the current Route.
	RouteStatusFields `json:",inline"`
}

// +kubebuilder:object:root=true

// ServiceList is a list of Service resources.
type ServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Service `json:"items"`
}

// ConfigurationSpec holds the desired state of the Configuration (from the client).
type ConfigurationSpec struct {
	// Template holds the latest specification for the Revision to be stamped out.
	Template RevisionTemplateSpec `json:"template"`
}

// ConfigurationStatusFields holds the fields of Configuration's status that
// are not generally shared.
type ConfigurationStatusFields struct {
	// LatestReadyRevisionName holds the name of the latest Revision stamped out
	// from this Configuration that has had its "Ready" condition become "True".
	LatestReadyRevisionName string `json:"latestReadyRevisionName,omitempty"`

	// LatestCreatedRevisionName is the last revision that was created from this
	// Configuration. It might not be ready yet, for that use LatestReadyRevisionName.
	LatestCreatedRevisionName string `json:"latestCreatedRevisionName,omitempty"`
}

// RevisionTemplateSpec describes the data a revision should have when created
// from a template.
type RevisionTemplateSpec struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec RevisionSpec `json:"spec,omitempty"`
}

// RevisionSpec holds the desired state of the Revision (from the client).
type RevisionSpec struct {
	corev1.PodSpec `json:",inline"`

	// ContainerConcurrency specifies the maximum allowed in-flight (concurrent)
	// requests per container of the Revision.
	ContainerConcurrency *int64 `json:"containerConcurrency,omitempty"`

	// TimeoutSeconds is the maximum duration in seconds that the request instance
	// is allowed to respond to a request.
	TimeoutSeconds *int64 `json:"timeoutSeconds,omitempty"`

	// ResponseStartTimeoutSeconds is the maximum duration in seconds that the
	// request routing layer will wait for a request delivered to a container to
	// begin sending any network traffic.
	ResponseStartTimeoutSeconds *int64 `json:"responseStartTimeoutSeconds,omitempty"`

	// IdleTimeoutSeconds is the maximum duration in seconds a request will be
	// allowed to stay open while not receiving any bytes from the user's application.
	IdleTimeoutSeconds *int64 `json:"idleTimeoutSeconds,omitempty"`
}

// RouteSpec holds the desired state of the Route (from the client).
type RouteSpec struct {
	// Traffic specifies how to distribute traffic over a collection of
	// revisions and configurations.
	Traffic []TrafficTarget `json:"traffic,omitempty"`
}

// RouteStatusFields holds the fields of Route's status that are not generally
// shared.
type RouteStatusFields struct {
	// URL holds the url that will distribute traffic over the provided traffic
	// targets. It generally has the form:
	// http[s]://{route-name}.{route-namespace}.{cluster-level-suffix}
	URL *apis.URL `json:"url,omitempty"`

	// Address holds the information needed for a Route to be the target of an event.
	Address *duckv1.Addressable `json:"address,omitempty"`

	// Traffic holds the configured traffic distribution.
	Traffic []TrafficTarget `json:"traffic,omitempty"`
}

// TrafficTarget holds a single entry of the routing table for a Route.
type TrafficTarget struct {
	// Tag is optionally used to expose a dedicated url for referencing this
	// target exclusively.
	Tag string `json:"tag,omitempty"`

	// RevisionName of a specific revision to which to send this portion of traffic.
	RevisionName string `json:"revisionName,omitempty"`

	// ConfigurationName of a configuration to whose latest revision we will send
	// this portion of traffic.
	ConfigurationName string `json:"configurationName,omitempty"`

	// LatestRevision may be optionally provided to indicate that the latest ready
	// Revision of the Configuration should be used for this traffic target.
	LatestRevision *bool `json:"latestRevision,omitempty"`

	// Percent indicates that percentage based routing should be used and the value
	// indicates the percent of traffic that is to be routed to this Revision or Configuration.
	Percent *int64 `json:"percent,omitempty"`

	// URL displays the URL for accessing named traffic targets. URL is displayed in
	// status, and is disallowed on spec. URL must contain a scheme (e.g. http://) and a hostname.
	URL *apis.URL `json:"url,omitempty"`
}
