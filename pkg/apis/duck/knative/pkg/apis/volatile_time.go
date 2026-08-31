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

// This file is a partial copy of knative.dev/pkg/apis/volatile_time.go
// (Apache License 2.0, The Knative Authors).

package apis

import (
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func init() {
	// Register a semantic-equality func that treats any two VolatileTime values
	// as equal. This mirrors the upstream knative.dev/pkg/apis init() and is the
	// whole point of the type: a change limited to a condition's
	// LastTransitionTime must not be seen as a status change by the operator's
	// StatusChangedPredicate, otherwise it triggers redundant reconciles.
	if err := equality.Semantic.AddFunc(
		func(VolatileTime, VolatileTime) bool {
			return true
		},
	); err != nil {
		panic(err)
	}
}

// VolatileTime wraps metav1.Time. It is used to represent times that are not
// semantically significant and thus should be ignored in equality checks.
type VolatileTime struct {
	Inner metav1.Time `json:",inline"`
}

// MarshalJSON implements the json.Marshaler interface.
func (t VolatileTime) MarshalJSON() ([]byte, error) {
	return t.Inner.MarshalJSON()
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *VolatileTime) UnmarshalJSON(b []byte) error {
	return t.Inner.UnmarshalJSON(b)
}
