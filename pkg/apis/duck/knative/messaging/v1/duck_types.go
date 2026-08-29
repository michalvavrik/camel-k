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

// The types in this file are a partial copy of the Knative Eventing Messaging
// API (knative.dev/eventing/pkg/apis/messaging/v1), Apache License 2.0 by The
// Knative Authors. Only the subset Camel K produces/resolves is retained, with
// the field layout kept identical to preserve serialization.

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/apache/camel-k/v2/pkg/apis/duck/knative/pkg/apis"
	duckv1 "github.com/apache/camel-k/v2/pkg/apis/duck/knative/pkg/apis/duck/v1"
)

// +kubebuilder:object:root=true

// Subscription routes events received on a Channel to a DNS name and
// corresponds to the subscriptions.channels.knative.dev CRD.
type Subscription struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	Spec   SubscriptionSpec   `json:"spec"`
	Status SubscriptionStatus `json:"status,omitempty"`
}

// SubscriptionSpec specifies the Channel for incoming events, a Subscriber
// target for processing those events and where to put the result of the
// processing.
type SubscriptionSpec struct {
	// Reference to a channel that will be used to create the subscription.
	Channel duckv1.KReference `json:"channel"`

	// Subscriber is reference to (optional) function for processing events.
	// +optional
	Subscriber *duckv1.Destination `json:"subscriber,omitempty"`

	// Reply specifies (optionally) how to handle events returned from the Subscriber target.
	// +optional
	Reply *duckv1.Destination `json:"reply,omitempty"`
}

// SubscriptionStatus (computed) for a subscription.
type SubscriptionStatus struct {
	// inherits duck/v1 Status.
	duckv1.Status `json:",inline"`

	// PhysicalSubscription is the fully resolved values that this Subscription represents.
	PhysicalSubscription SubscriptionStatusPhysicalSubscription `json:"physicalSubscription,omitempty"`
}

// SubscriptionStatusPhysicalSubscription represents the fully resolved values
// for this Subscription.
type SubscriptionStatusPhysicalSubscription struct {
	// SubscriberURI is the fully resolved URI for spec.subscriber.
	// +optional
	SubscriberURI *apis.URL `json:"subscriberUri,omitempty"`

	// ReplyURI is the fully resolved URI for the spec.reply.
	// +optional
	ReplyURI *apis.URL `json:"replyUri,omitempty"`

	// DeadLetterSinkURI is the fully resolved URI for the dead letter sink.
	// +optional
	DeadLetterSinkURI *apis.URL `json:"deadLetterSinkUri,omitempty"`
}

// +kubebuilder:object:root=true

// SubscriptionList is a collection of Subscriptions.
type SubscriptionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Subscription `json:"items"`
}

// +kubebuilder:object:root=true

// Channel represents a generic Channel. It is normally used when we want a
// Channel but do not need a specific Channel implementation.
type Channel struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ChannelSpec   `json:"spec,omitempty"`
	Status ChannelStatus `json:"status,omitempty"`
}

// ChannelSpec holds the desired state of the Channel.
type ChannelSpec struct {
	// ChannelTemplate specifies which Channel CRD to use to create the backing
	// Channel for this generic Channel.
	// +optional
	ChannelTemplate *ChannelTemplateSpec `json:"channelTemplate,omitempty"`
}

// ChannelTemplateSpec identifies the underlying Channel implementation.
type ChannelTemplateSpec struct {
	metav1.TypeMeta `json:",inline"`
}

// ChannelStatus represents the current state of a Channel.
type ChannelStatus struct {
	// inherits duck/v1 Status.
	duckv1.Status `json:",inline"`

	// AddressStatus is the part where the Channel fulfills the Addressable contract.
	duckv1.AddressStatus `json:",inline"`
}

// +kubebuilder:object:root=true

// ChannelList is a collection of Channels.
type ChannelList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Channel `json:"items"`
}
