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

// The types in this file are a partial copy of the Knative Eventing API
// (knative.dev/eventing/pkg/apis/eventing/v1), Apache License 2.0 by The Knative
// Authors. Only the subset Camel K produces/resolves is retained, with the field
// layout kept identical to preserve serialization.

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	duckv1 "github.com/apache/camel-k/v2/pkg/apis/duck/knative/pkg/apis/duck/v1"
)

// +kubebuilder:object:root=true

// Trigger represents a request to have events delivered to a Subscriber from a
// Broker's event pool.
type Trigger struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec defines the desired state of the Trigger.
	Spec TriggerSpec `json:"spec,omitempty"`

	// Status represents the current state of the Trigger. This data may be out of date.
	// +optional
	Status TriggerStatus `json:"status,omitempty"`
}

// TriggerSpec defines the desired state of the Trigger.
type TriggerSpec struct {
	// Broker is the broker that this trigger receives events from.
	Broker string `json:"broker,omitempty"`

	// Filter is the filter to apply against all events from the Broker. Only
	// events that pass this filter will be sent to the Subscriber.
	// +optional
	Filter *TriggerFilter `json:"filter,omitempty"`

	// Subscriber is the addressable that receives events from the Broker that
	// pass the Filter. It is required.
	Subscriber duckv1.Destination `json:"subscriber"`
}

// TriggerFilter allows filtering events destined for a Trigger's Subscriber.
type TriggerFilter struct {
	// Attributes filters events by exact match on event context attributes.
	// +optional
	Attributes TriggerFilterAttributes `json:"attributes,omitempty"`
}

// TriggerFilterAttributes is a map of context attribute names to values for
// filtering by equality.
type TriggerFilterAttributes map[string]string

// TriggerStatus represents the current state of a Trigger.
type TriggerStatus struct {
	// inherits duck/v1 Status.
	duckv1.Status `json:",inline"`
}

// +kubebuilder:object:root=true

// TriggerList is a collection of Triggers.
type TriggerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Trigger `json:"items"`
}

// +kubebuilder:object:root=true

// Broker collects a pool of events that are consumable using Triggers. Brokers
// provide a well-known endpoint for event delivery that senders can use with
// minimal knowledge of the event routing strategy.
type Broker struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec defines the desired state of the Broker.
	Spec BrokerSpec `json:"spec,omitempty"`

	// Status represents the current state of the Broker. This data may be out of date.
	// +optional
	Status BrokerStatus `json:"status,omitempty"`
}

// BrokerSpec defines the desired state of the Broker.
type BrokerSpec struct {
	// Config is a KReference to the configuration that specifies configuration
	// options for this Broker.
	// +optional
	Config *duckv1.KReference `json:"config,omitempty"`
}

// BrokerStatus represents the current state of a Broker.
type BrokerStatus struct {
	// inherits duck/v1 Status.
	duckv1.Status `json:",inline"`

	// AddressStatus is the part where the Broker fulfills the Addressable contract.
	duckv1.AddressStatus `json:",inline"`
}

// +kubebuilder:object:root=true

// BrokerList is a collection of Brokers.
type BrokerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Broker `json:"items"`
}
