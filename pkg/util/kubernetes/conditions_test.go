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

package kubernetes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"knative.dev/pkg/apis"
	duckv1 "knative.dev/pkg/apis/duck/v1"
	servingv1 "knative.dev/serving/pkg/apis/serving/v1"
)

func TestGetKnativeServiceConditionReady(t *testing.T) {
	svc := servingv1.Service{
		Status: servingv1.ServiceStatus{
			Status: duckv1.Status{
				Conditions: duckv1.Conditions{
					{
						Type:    servingv1.ServiceConditionReady,
						Status:  corev1.ConditionTrue,
						Reason:  "AllGood",
						Message: "service is ready",
					},
				},
			},
		},
	}

	cond := GetKnativeServiceCondition(svc, servingv1.ServiceConditionReady)
	require.NotNil(t, cond)
	assert.True(t, cond.IsTrue())
	assert.False(t, cond.IsFalse())
	assert.Equal(t, "AllGood", cond.GetReason())
	assert.Equal(t, "service is ready", cond.GetMessage())
}

func TestGetKnativeServiceConditionFalseWithReason(t *testing.T) {
	svc := servingv1.Service{
		Status: servingv1.ServiceStatus{
			Status: duckv1.Status{
				Conditions: duckv1.Conditions{
					{
						Type:    servingv1.ServiceConditionReady,
						Status:  corev1.ConditionFalse,
						Reason:  "RevisionFailed",
						Message: "revision failed with exit code 1",
					},
				},
			},
		},
	}

	cond := GetKnativeServiceCondition(svc, servingv1.ServiceConditionReady)
	require.NotNil(t, cond)
	assert.True(t, cond.IsFalse())
	assert.False(t, cond.IsTrue())
	assert.Equal(t, "RevisionFailed", cond.GetReason())
	assert.Equal(t, "revision failed with exit code 1", cond.GetMessage())
}

func TestGetKnativeServiceConditionMissing(t *testing.T) {
	svc := servingv1.Service{}

	cond := GetKnativeServiceCondition(svc, servingv1.ServiceConditionReady)
	assert.Nil(t, cond)
}

func TestGetKnativeServiceConditionUnknownStatus(t *testing.T) {
	svc := servingv1.Service{
		Status: servingv1.ServiceStatus{
			Status: duckv1.Status{
				Conditions: duckv1.Conditions{
					{
						Type:   servingv1.ServiceConditionReady,
						Status: corev1.ConditionUnknown,
						Reason: "Deploying",
					},
				},
			},
		},
	}

	cond := GetKnativeServiceCondition(svc, servingv1.ServiceConditionReady)
	require.NotNil(t, cond)
	assert.False(t, cond.IsTrue())
	assert.False(t, cond.IsFalse())
	assert.Equal(t, "Deploying", cond.GetReason())
}

func TestGetKnativeServiceConditionWrongType(t *testing.T) {
	svc := servingv1.Service{
		Status: servingv1.ServiceStatus{
			Status: duckv1.Status{
				Conditions: duckv1.Conditions{
					{
						Type:   apis.ConditionType("ConfigurationsReady"),
						Status: corev1.ConditionTrue,
					},
				},
			},
		},
	}

	cond := GetKnativeServiceCondition(svc, servingv1.ServiceConditionReady)
	assert.Nil(t, cond)
}
