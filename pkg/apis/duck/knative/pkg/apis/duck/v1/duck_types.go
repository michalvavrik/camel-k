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

// The types in this file are a partial copy of knative.dev/pkg/apis/duck/v1
// (Apache License 2.0, The Knative Authors). Only the subset used by Camel K is
// retained, with the field layout kept identical to preserve serialization.

package v1

import (
	"github.com/apache/camel-k/v2/pkg/apis/duck/knative/pkg/apis"
	"github.com/apache/camel-k/v2/pkg/apis/duck/knative/pkg/tracker"
)

// Conditions is an alias of apis.Conditions, re-exported here to mirror the
// upstream knative.dev/pkg/apis/duck/v1 package.
type Conditions = apis.Conditions

// Status shows how we expect folks to embed Conditions in their Status field.
type Status struct {
	// ObservedGeneration is the 'Generation' of the Service that was last processed by the controller.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// Conditions the latest available observations of a resource's current state.
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions apis.Conditions `json:"conditions,omitempty" patchMergeKey:"type" patchStrategy:"merge"`

	// Annotations is additional Status fields for the Resource to save some additional State
	// as well as convey more information to the user. This is roughly akin to Annotations on any
	// k8s resource, just the reconciler conveying richer information outwards.
	Annotations map[string]string `json:"annotations,omitempty"`
}

// AddressStatus shows how we expect folks to embed Addressable in their Status field.
type AddressStatus struct {
	// Address is a single Addressable address. If Addresses is present, Address will be ignored by clients.
	// +optional
	Address *Addressable `json:"address,omitempty"`

	// Addresses is a list of addresses for different protocols (HTTP and HTTPS).
	// If Addresses is present, Address must be ignored by clients.
	// +optional
	Addresses []Addressable `json:"addresses,omitempty"`
}

// Addressable provides a generic mechanism for a custom resource definition to
// indicate a destination for message delivery.
type Addressable struct {
	// Name is the name of the address.
	// +optional
	Name *string `json:"name,omitempty"`

	// URL is the address URL.
	URL *apis.URL `json:"url,omitempty"`

	// CACerts is the Certification Authority (CA) certificates in PEM format
	// according to https://www.rfc-editor.org/rfc/rfc7468.
	// +optional
	CACerts *string `json:"CACerts,omitempty"`

	// Audience is the OIDC audience for this address.
	// +optional
	Audience *string `json:"audience,omitempty"`
}

// Destination represents a target of an invocation over HTTP.
type Destination struct {
	// Ref points to an Addressable.
	// +optional
	Ref *KReference `json:"ref,omitempty"`

	// URI can be an absolute URL(non-empty scheme and non-empty host) pointing
	// to the target or a relative URI. Relative URIs will be resolved using the
	// base URI retrieved from Ref.
	// +optional
	URI *apis.URL `json:"uri,omitempty"`

	// CACerts are Certification Authority (CA) certificates in PEM format
	// according to https://www.rfc-editor.org/rfc/rfc7468.
	// +optional
	CACerts *string `json:"CACerts,omitempty"`

	// Audience is the OIDC audience.
	// +optional
	Audience *string `json:"audience,omitempty"`
}

// KReference contains enough information to refer to another object.
type KReference struct {
	// Kind of the referent.
	Kind string `json:"kind"`

	// Namespace of the referent.
	// +optional
	Namespace string `json:"namespace,omitempty"`

	// Name of the referent.
	Name string `json:"name"`

	// API version of the referent.
	// +optional
	APIVersion string `json:"apiVersion,omitempty"`

	// Group of the API, without the version of the group. This can be used as an
	// alternative to the APIVersion.
	// +optional
	Group string `json:"group,omitempty"`

	// Address points to a specific Address Name.
	// +optional
	Address *string `json:"address,omitempty"`
}

// SourceSpec is the minimum resource shape to adhere to the Source Specification.
type SourceSpec struct {
	// Sink is a reference to an object that will resolve to a uri to use as the sink.
	Sink Destination `json:"sink,omitempty"`

	// CloudEventOverrides defines overrides to control the output format and
	// modifications of the event sent to the sink.
	// +optional
	CloudEventOverrides *CloudEventOverrides `json:"ceOverrides,omitempty"`
}

// CloudEventOverrides defines arguments for a Source that control the output
// format of the CloudEvents produced by the Source.
type CloudEventOverrides struct {
	// Extensions specify what attribute are added or overridden on the outbound event.
	// +optional
	Extensions map[string]string `json:"extensions,omitempty"`
}

// BindingSpec is the minimum resource shape to adhere to the Binding Specification.
type BindingSpec struct {
	// Subject references the resource(s) whose "runtime contract" should be
	// augmented by Binding implementations.
	Subject tracker.Reference `json:"subject"`
}
