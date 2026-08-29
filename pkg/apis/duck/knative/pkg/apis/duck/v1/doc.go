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

// Package v1 contains a partial, vendored copy of the knative.dev/pkg/apis/duck/v1
// types used by the Camel K Knative duck types (Status, Addressable, Destination,
// KReference, SourceSpec, BindingSpec). Copied from the Apache License 2.0 sources
// authored by The Knative Authors, keeping the serialized form identical.
//
// +kubebuilder:object:generate=true
package v1
