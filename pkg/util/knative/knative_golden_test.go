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
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	corev1 "k8s.io/api/core/v1"
)

// assertGolden marshals obj to indented JSON and compares it byte-for-byte with a
// checked-in golden file. Regenerate with CAMEL_K_TEST_UPDATE_GOLDEN=true.
//
// These goldens are the 1:1 anchor for the Knative duck-type migration (#6354):
// after the knative.dev types are replaced by hand-rolled duck types, the objects
// camel-k builds MUST serialize to exactly the same bytes, or these tests fail.
func assertGolden(t *testing.T, obj any, name string) {
	t.Helper()

	actual, err := json.MarshalIndent(obj, "", "  ")
	require.NoError(t, err)
	actual = append(actual, '\n')

	path := filepath.Join("testdata", name)
	if os.Getenv("CAMEL_K_TEST_UPDATE_GOLDEN") == "true" {
		require.NoError(t, os.MkdirAll("testdata", 0o755))
		require.NoError(t, os.WriteFile(path, actual, 0o600))
		return
	}

	expected, err := os.ReadFile(path)
	require.NoError(t, err, "missing golden %s; regenerate with CAMEL_K_TEST_UPDATE_GOLDEN=true", path)
	assert.Equal(t, string(expected), string(actual),
		"serialization changed for %s; if this change is intended, regenerate with "+
			"CAMEL_K_TEST_UPDATE_GOLDEN=true and review the diff carefully", name)
}

func TestGoldenSinkBinding(t *testing.T) {
	sb := CreateSinkBinding(
		corev1.ObjectReference{APIVersion: "apps/v1", Kind: "Deployment", Name: "my-source", Namespace: "my-ns"},
		corev1.ObjectReference{APIVersion: "messaging.knative.dev/v1", Kind: "Channel", Name: "my-channel"},
	)
	assertGolden(t, sb, "sinkbinding.golden.json")
}

func TestGoldenSubscription(t *testing.T) {
	sub := CreateSubscription(
		corev1.ObjectReference{APIVersion: "messaging.knative.dev/v1", Kind: "Channel", Name: "my-channel", Namespace: "my-ns"},
		"my-service", "/channels/my-channel",
	)
	assertGolden(t, sub, "subscription.golden.json")
}

func TestGoldenTrigger(t *testing.T) {
	trigger, err := CreateKnativeServiceTrigger(
		corev1.ObjectReference{APIVersion: "eventing.knative.dev/v1", Kind: "Broker", Name: "default", Namespace: "my-ns"},
		"my-service", "evt.type", "/events/evt.type", map[string]string{"type": "evt.type"},
	)
	require.NoError(t, err)
	assertGolden(t, trigger, "trigger.golden.json")
}
