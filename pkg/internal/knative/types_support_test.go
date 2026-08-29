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
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetSinkBinding(t *testing.T) {
	camelEnv := NewCamelEnvironment()
	svc1, err := BuildCamelServiceDefinition(
		"test",
		CamelEndpointKindSink,
		CamelServiceTypeChannel,
		url.URL{},
		"apiVersion",
		"InMemoryChannel",
	)
	require.NoError(t, err)
	camelEnv.Services = append(camelEnv.Services, svc1)
	svc := camelEnv.FindService("test",
		CamelEndpointKindSink,
		CamelServiceTypeChannel,
		"apiVersion",
		"InMemoryChannel",
	)
	assert.NotNil(t, svc)
	assert.False(t, svc.SinkBinding)
	camelEnv.SetSinkBinding("test",
		CamelEndpointKindSink,
		CamelServiceTypeChannel,
		"apiVersion",
		"InMemoryChannel",
	)
	svc = camelEnv.FindService("test",
		CamelEndpointKindSink,
		CamelServiceTypeChannel,
		"apiVersion",
		"InMemoryChannel",
	)
	assert.NotNil(t, svc)
	assert.True(t, svc.SinkBinding)
}

// TestToCamelProperties pins the application-properties mapping the knative trait
// serializes into the Camel Knative component config. This is the wire format the
// running Integration consumes, so it must stay byte-stable across the duck-type
// migration (#6354).
func TestToCamelProperties(t *testing.T) {
	env := NewCamelEnvironment()
	env.Services = append(env.Services,
		// index 0: non-sinkbinding endpoint sink with url + path + reply
		CamelServiceDefinition{
			Name:        "ep",
			ServiceType: CamelServiceTypeEndpoint,
			URL:         "http://ep.host/",
			Path:        "/mypath",
			Metadata: map[string]string{
				CamelMetaKnativeKind:       "Service",
				CamelMetaKnativeAPIVersion: "serving.knative.dev/v1",
				CamelMetaEndpointKind:      string(CamelEndpointKindSink),
				CamelMetaKnativeReply:      "false",
				CamelMetaKnativeName:       "ep",
			},
		},
		// index 1: event (broker) source -> objectName is emitted
		CamelServiceDefinition{
			Name:        "default",
			ServiceType: CamelServiceTypeEvent,
			URL:         "http://broker.host/",
			Metadata: map[string]string{
				CamelMetaKnativeKind:       "Broker",
				CamelMetaKnativeAPIVersion: "eventing.knative.dev/v1",
				CamelMetaEndpointKind:      string(CamelEndpointKindSource),
				CamelMetaKnativeName:       "default",
			},
		},
		// index 2: sinkbinding channel sink -> K_SINK / K_CE_OVERRIDES, no path
		CamelServiceDefinition{
			Name:        "ch",
			ServiceType: CamelServiceTypeChannel,
			SinkBinding: true,
			Metadata: map[string]string{
				CamelMetaKnativeKind:       "Channel",
				CamelMetaKnativeAPIVersion: "messaging.knative.dev/v1",
				CamelMetaEndpointKind:      string(CamelEndpointKindSink),
			},
		},
	)

	props := env.ToCamelProperties()

	p0 := "camel.component.knative.environment.resources[0]"
	assert.Equal(t, "ep", props[p0+".name"])
	assert.Equal(t, "endpoint", props[p0+".type"])
	assert.Equal(t, "Service", props[p0+".objectKind"])
	assert.Equal(t, "serving.knative.dev/v1", props[p0+".objectApiVersion"])
	assert.Equal(t, "sink", props[p0+".endpointKind"])
	assert.Equal(t, "false", props[p0+".reply"])
	assert.Equal(t, "http://ep.host/", props[p0+".url"])
	assert.Equal(t, "/mypath", props[p0+".path"])
	assert.NotContains(t, props, p0+".objectName", "objectName is only emitted for event services")
	assert.NotContains(t, props, p0+".ceOverrides", "ceOverrides is only emitted for sinkBinding services")

	p1 := "camel.component.knative.environment.resources[1]"
	assert.Equal(t, "event", props[p1+".type"])
	assert.Equal(t, "default", props[p1+".objectName"])
	assert.Equal(t, "http://broker.host/", props[p1+".url"])

	p2 := "camel.component.knative.environment.resources[2]"
	assert.Equal(t, "${K_SINK}", props[p2+".url"])
	assert.Equal(t, "${K_CE_OVERRIDES}", props[p2+".ceOverrides"])
	assert.NotContains(t, props, p2+".path", "sinkBinding services must not emit a path")
}

// TestCamelEnvironmentSerializeRoundTrip locks the JSON encoding of the Camel
// environment (KAMEL_KNATIVE_CONFIGURATION) and its idempotency.
func TestCamelEnvironmentSerializeRoundTrip(t *testing.T) {
	env := NewCamelEnvironment()
	env.Services = append(env.Services, CamelServiceDefinition{
		Name:        "ep",
		ServiceType: CamelServiceTypeEndpoint,
		URL:         "http://ep/",
		Path:        "/p",
		Metadata:    map[string]string{CamelMetaKnativeKind: "Service"},
		SinkBinding: true,
	})

	s, err := env.Serialize()
	require.NoError(t, err)
	assert.JSONEq(t,
		`{"services":[{"type":"endpoint","name":"ep","url":"http://ep/","path":"/p","metadata":{"knative.kind":"Service"},"sinkBinding":true}]}`,
		s,
	)

	var env2 CamelEnvironment
	require.NoError(t, env2.Deserialize(s))
	assert.Equal(t, env, env2)

	// Serialization must be idempotent.
	s2, err := env2.Serialize()
	require.NoError(t, err)
	assert.Equal(t, s, s2)
}
