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

// Package v1 contains a partial (duck-typed) schema of the Knative Serving API
// (serving.knative.dev/v1). It mirrors only the subset of the upstream
// knative.dev/serving types that Camel K produces and consumes, so that the
// operator no longer needs to depend on the full knative.dev/serving module.
//
// The struct shapes (field names, JSON tags and embeddings) are intentionally
// kept byte-for-byte identical to the upstream types to preserve the exact
// wire format of the Knative Service resources Camel K creates.
//
// +kubebuilder:object:generate=true
// +groupName=serving.knative.dev
package v1
