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

package knative

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	corev1 "k8s.io/api/core/v1"
)

// TestCreateSinkBinding pins the full SinkBinding spec produced by camel-k. The
// field layout (Subject as a tracker.Reference, Sink as a duckv1.Destination ->
// KReference) is an external contract consumed by Knative's SinkBinding admission
// webhook, so a duck-type migration (#6354) must reproduce it 1:1.
func TestCreateSinkBinding(t *testing.T) {
	source := corev1.ObjectReference{
		APIVersion: "apps/v1",
		Kind:       "Deployment",
		Name:       "my-source",
		Namespace:  "my-ns",
	}
	target := corev1.ObjectReference{
		APIVersion: "messaging.knative.dev/v1",
		Kind:       "Channel",
		Name:       "my-channel",
	}

	sb := CreateSinkBinding(source, target)
	require.NotNil(t, sb)

	assert.Equal(t, "sources.knative.dev/v1", sb.APIVersion)
	assert.Equal(t, "SinkBinding", sb.Kind)
	assert.Equal(t, "my-ns", sb.Namespace)
	assert.Equal(t, "my-source", sb.Name)

	// Subject == source workload
	assert.Equal(t, "apps/v1", sb.Spec.Subject.APIVersion)
	assert.Equal(t, "Deployment", sb.Spec.Subject.Kind)
	assert.Equal(t, "my-source", sb.Spec.Subject.Name)

	// Sink == target addressable
	require.NotNil(t, sb.Spec.Sink.Ref)
	assert.Equal(t, "messaging.knative.dev/v1", sb.Spec.Sink.Ref.APIVersion)
	assert.Equal(t, "Channel", sb.Spec.Sink.Ref.Kind)
	assert.Equal(t, "my-channel", sb.Spec.Sink.Ref.Name)
}

// TestCreateSubscription pins the full Subscription spec (Channel KReference and
// the Subscriber Destination with a serving Service ref and an apis.URL path).
func TestCreateSubscription(t *testing.T) {
	channelRef := corev1.ObjectReference{
		APIVersion: "messaging.knative.dev/v1",
		Kind:       "Channel",
		Name:       "my-channel",
		Namespace:  "my-ns",
	}

	sub := CreateSubscription(channelRef, "my-service", "/channels/my-channel")
	require.NotNil(t, sub)

	assert.Equal(t, "messaging.knative.dev/v1", sub.APIVersion)
	assert.Equal(t, "Subscription", sub.Kind)
	assert.Equal(t, "my-ns", sub.Namespace)
	assert.Equal(t, "my-channel-my-service", sub.Name)

	assert.Equal(t, "messaging.knative.dev/v1", sub.Spec.Channel.APIVersion)
	assert.Equal(t, "Channel", sub.Spec.Channel.Kind)
	assert.Equal(t, "my-channel", sub.Spec.Channel.Name)

	require.NotNil(t, sub.Spec.Subscriber)
	require.NotNil(t, sub.Spec.Subscriber.Ref)
	assert.Equal(t, "serving.knative.dev/v1", sub.Spec.Subscriber.Ref.APIVersion)
	assert.Equal(t, "Service", sub.Spec.Subscriber.Ref.Kind)
	assert.Equal(t, "my-service", sub.Spec.Subscriber.Ref.Name)
	require.NotNil(t, sub.Spec.Subscriber.URI)
	assert.Equal(t, "/channels/my-channel", sub.Spec.Subscriber.URI.Path)

	// The apis.URL embedded in the Destination MUST marshal to a plain string
	// (custom MarshalJSON). A duck-type replacement must preserve this contract,
	// otherwise the emitted Subscription/Trigger CRs would be malformed.
	raw, err := json.Marshal(sub.Spec.Subscriber)
	require.NoError(t, err)
	assert.Contains(t, string(raw), `"uri":"/channels/my-channel"`)
}

func TestCreateKnativeServiceTrigger(t *testing.T) {
	broker := corev1.ObjectReference{
		APIVersion: "eventing.knative.dev/v1",
		Kind:       "Broker",
		Name:       "default",
		Namespace:  "my-ns",
	}

	trigger, err := CreateKnativeServiceTrigger(broker, "my-service", "evt.type", "/events/evt.type", map[string]string{
		"type": "evt.type",
	})
	require.NoError(t, err)
	require.NotNil(t, trigger)

	assert.Equal(t, "eventing.knative.dev/v1", trigger.APIVersion)
	assert.Equal(t, "Trigger", trigger.Kind)
	assert.Equal(t, "default-my-service-evttype", trigger.Name)
	assert.Equal(t, "default", trigger.Spec.Broker)

	require.NotNil(t, trigger.Spec.Subscriber.Ref)
	assert.Equal(t, "serving.knative.dev/v1", trigger.Spec.Subscriber.Ref.APIVersion)
	assert.Equal(t, "Service", trigger.Spec.Subscriber.Ref.Kind)
	assert.Equal(t, "my-service", trigger.Spec.Subscriber.Ref.Name)
	require.NotNil(t, trigger.Spec.Subscriber.URI)
	assert.Equal(t, "/events/evt.type", trigger.Spec.Subscriber.URI.Path)

	require.NotNil(t, trigger.Spec.Filter)
	assert.Equal(t, "evt.type", trigger.Spec.Filter.Attributes["type"])
}
