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

package trait

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1 "github.com/apache/camel-k/v2/pkg/apis/camel/v1"
	traitv1 "github.com/apache/camel-k/v2/pkg/apis/camel/v1/trait"
)

// assertKnativeServiceGolden marshals obj to indented JSON and compares it
// byte-for-byte with a checked-in golden file. Regenerate with
// CAMEL_K_TEST_UPDATE_GOLDEN=true.
//
// This is the 1:1 anchor for the serving side of the Knative duck-type migration
// (#6354): the Knative Service camel-k builds MUST serialize identically after the
// knative.dev/serving type is replaced by a hand-rolled duck type.
func assertKnativeServiceGolden(t *testing.T, obj any, name string) {
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

// TestGoldenKnativeService pins the full serialization of the Knative Service
// skeleton produced by the knative-service trait (autoscaling annotations, rollout,
// visibility, timeout, service account, revision template), isolating the serving
// type from the container/mount traits by calling getServiceFor directly.
func TestGoldenKnativeService(t *testing.T) {
	tr, ok := newKnativeServiceTrait().(*knativeServiceTrait)
	require.True(t, ok)

	target := 80
	minScale := 1
	maxScale := 5
	timeout := int64(44)
	tr.Class = "hpa.autoscaling.knative.dev"
	tr.Metric = "cpu"
	tr.Target = &target
	tr.MinScale = &minScale
	tr.MaxScale = &maxScale
	tr.RolloutDuration = "60s"
	tr.Visibility = "cluster-local"
	tr.TimeoutSeconds = &timeout
	tr.Annotations = map[string]string{"my-annotation": "my-value"}

	e := &Environment{
		Integration: &v1.Integration{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-it",
				Namespace: "my-ns",
			},
			Spec: v1.IntegrationSpec{
				ServiceAccountName: "my-sa",
				Traits:             v1.Traits{KnativeService: &traitv1.KnativeServiceTrait{}},
			},
		},
	}

	svc, err := tr.getServiceFor(e)
	require.NoError(t, err)
	assertKnativeServiceGolden(t, svc, "knative_service.golden.json")
}
