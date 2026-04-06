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

package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	duckv1 "knative.dev/pkg/apis/duck/v1"
	servingv1 "knative.dev/serving/pkg/apis/serving/v1"

	v1 "github.com/apache/camel-k/v2/pkg/apis/camel/v1"
)

func TestKnativeServiceControllerRevisionFailed(t *testing.T) {
	integration := &v1.Integration{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-it",
			Namespace: "ns",
		},
		Status: v1.IntegrationStatus{
			Phase: v1.IntegrationPhaseRunning,
		},
	}

	ksvc := &servingv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-ksvc",
			Namespace: "ns",
		},
		Status: servingv1.ServiceStatus{
			Status: duckv1.Status{
				Conditions: duckv1.Conditions{
					{
						Type:    servingv1.ServiceConditionReady,
						Status:  corev1.ConditionFalse,
						Reason:  "RevisionFailed",
						Message: "container exited with code 1",
					},
				},
			},
		},
	}

	ctrl := &knativeServiceController{
		obj:         ksvc,
		integration: integration,
	}

	handled, err := ctrl.checkReadyCondition(context.Background())
	assert.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, v1.IntegrationPhaseError, integration.Status.Phase)
	readyCond := integration.Status.GetCondition(v1.IntegrationConditionReady)
	assert.Equal(t, corev1.ConditionFalse, readyCond.Status)
	assert.Equal(t, "container exited with code 1", readyCond.Message)
}

func TestKnativeServiceControllerNotReady(t *testing.T) {
	integration := &v1.Integration{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-it",
			Namespace: "ns",
		},
		Status: v1.IntegrationStatus{
			Phase: v1.IntegrationPhaseRunning,
		},
	}

	ksvc := &servingv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-ksvc",
			Namespace: "ns",
		},
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

	ctrl := &knativeServiceController{
		obj:         ksvc,
		integration: integration,
	}

	handled, err := ctrl.checkReadyCondition(context.Background())
	assert.NoError(t, err)
	assert.False(t, handled)
	assert.Equal(t, v1.IntegrationPhaseRunning, integration.Status.Phase)
}

func TestKnativeServiceControllerUpdateReadyConditionTrue(t *testing.T) {
	integration := &v1.Integration{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-it",
			Namespace: "ns",
		},
		Status: v1.IntegrationStatus{
			Phase: v1.IntegrationPhaseRunning,
		},
	}

	ksvc := &servingv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-ksvc",
			Namespace: "ns",
		},
		Status: servingv1.ServiceStatus{
			Status: duckv1.Status{
				Conditions: duckv1.Conditions{
					{
						Type:   servingv1.ServiceConditionReady,
						Status: corev1.ConditionTrue,
					},
				},
			},
		},
	}

	ctrl := &knativeServiceController{
		obj:         ksvc,
		integration: integration,
	}

	ready := ctrl.updateReadyCondition(1)
	assert.True(t, ready)
	readyCond := integration.Status.GetCondition(v1.IntegrationConditionReady)
	assert.Equal(t, corev1.ConditionTrue, readyCond.Status)
	assert.Equal(t, v1.IntegrationConditionKnativeServiceReadyReason, readyCond.Reason)
}

func TestKnativeServiceControllerUpdateReadyConditionFalse(t *testing.T) {
	integration := &v1.Integration{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-it",
			Namespace: "ns",
		},
		Status: v1.IntegrationStatus{
			Phase: v1.IntegrationPhaseRunning,
		},
	}

	ksvc := &servingv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-ksvc",
			Namespace: "ns",
		},
		Status: servingv1.ServiceStatus{
			Status: duckv1.Status{
				Conditions: duckv1.Conditions{
					{
						Type:    servingv1.ServiceConditionReady,
						Status:  corev1.ConditionFalse,
						Reason:  "NotYet",
						Message: "still deploying",
					},
				},
			},
		},
	}

	ctrl := &knativeServiceController{
		obj:         ksvc,
		integration: integration,
	}

	ready := ctrl.updateReadyCondition(0)
	assert.False(t, ready)
	readyCond := integration.Status.GetCondition(v1.IntegrationConditionReady)
	assert.Equal(t, corev1.ConditionFalse, readyCond.Status)
	assert.Equal(t, "NotYet", readyCond.Reason)
	assert.Equal(t, "still deploying", readyCond.Message)
}

func TestKnativeServiceControllerHasTemplateIntegrationLabel(t *testing.T) {
	ksvc := &servingv1.Service{
		Spec: servingv1.ServiceSpec{
			ConfigurationSpec: servingv1.ConfigurationSpec{
				Template: servingv1.RevisionTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							v1.IntegrationLabel: "my-it",
						},
					},
				},
			},
		},
	}

	ctrl := &knativeServiceController{obj: ksvc}
	assert.True(t, ctrl.hasTemplateIntegrationLabel())

	ksvcNoLabel := &servingv1.Service{}
	ctrlNoLabel := &knativeServiceController{obj: ksvcNoLabel}
	assert.False(t, ctrlNoLabel.hasTemplateIntegrationLabel())
}

func TestKnativeServiceControllerGetControllerName(t *testing.T) {
	ksvc := &servingv1.Service{
		ObjectMeta: metav1.ObjectMeta{Name: "my-ksvc"},
	}

	ctrl := &knativeServiceController{obj: ksvc}
	assert.Equal(t, "KnativeService/my-ksvc", ctrl.getControllerName())
}
