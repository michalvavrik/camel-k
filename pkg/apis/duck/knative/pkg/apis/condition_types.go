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

// This file is a partial copy of knative.dev/pkg/apis/condition_types.go
// (Apache License 2.0, The Knative Authors). The accessor methods are nil-safe,
// exactly as upstream, since Camel K relies on that behavior when a Knative
// Service has no Ready condition yet.

package apis

import (
	corev1 "k8s.io/api/core/v1"
)

// Conditions is the schema for the conditions portion of the payload.
type Conditions []Condition

// ConditionType is used to communicate the status of the reconciliation process.
type ConditionType string

const (
	// ConditionReady specifies that the resource is ready.
	ConditionReady ConditionType = "Ready"
	// ConditionSucceeded specifies that the resource has finished. For long-running resources.
	ConditionSucceeded ConditionType = "Succeeded"
)

// ConditionSeverity expresses the severity of a Condition Type failing.
type ConditionSeverity string

const (
	// ConditionSeverityError specifies that a failure of a condition type should be viewed as an error.
	ConditionSeverityError ConditionSeverity = ""
	// ConditionSeverityWarning specifies that a failure of a condition type should be viewed as a warning.
	ConditionSeverityWarning ConditionSeverity = "Warning"
	// ConditionSeverityInfo specifies that a failure of a condition type should be viewed as purely informational.
	ConditionSeverityInfo ConditionSeverity = "Info"
)

// Condition defines a readiness condition for a Knative resource.
type Condition struct {
	// Type of condition.
	// +required
	Type ConditionType `json:"type"`

	// Status of the condition, one of True, False, Unknown.
	// +required
	Status corev1.ConditionStatus `json:"status"`

	// Severity with which to treat failures of this type of condition.
	// When this is not specified, it defaults to Error.
	// +optional
	Severity ConditionSeverity `json:"severity,omitempty"`

	// LastTransitionTime is the last time the condition transitioned from one status to another.
	// +optional
	LastTransitionTime VolatileTime `json:"lastTransitionTime,omitempty"`

	// The reason for the condition's last transition.
	// +optional
	Reason string `json:"reason,omitempty"`

	// A human readable message indicating details about the transition.
	// +optional
	Message string `json:"message,omitempty"`
}

// IsTrue is true if the condition is True.
func (c *Condition) IsTrue() bool {
	if c == nil {
		return false
	}

	return c.Status == corev1.ConditionTrue
}

// IsFalse is true if the condition is False.
func (c *Condition) IsFalse() bool {
	if c == nil {
		return false
	}

	return c.Status == corev1.ConditionFalse
}

// IsUnknown is true if the condition is Unknown.
func (c *Condition) IsUnknown() bool {
	if c == nil {
		return true
	}

	return c.Status == corev1.ConditionUnknown
}

// GetReason returns the reason of the condition.
func (c *Condition) GetReason() string {
	if c == nil {
		return ""
	}

	return c.Reason
}

// GetMessage returns the message of the condition.
func (c *Condition) GetMessage() string {
	if c == nil {
		return ""
	}

	return c.Message
}
