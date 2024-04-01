/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package v1 contains API Schema definitions for the slv v1 API group
// +kubebuilder:object:generate=true
// +groupName=slv.oss.amagi.com
package v1

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"oss.amagi.com/slv/core/config"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

const (
	Group   = config.K8SLVGroup
	Version = config.K8SLVVersion
	Kind    = config.K8SLVKind
)

var (
	// GroupVersion is group version used to register these objects
	GroupVersion = schema.GroupVersion{Group: Group, Version: Version}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilder = &scheme.Builder{GroupVersion: GroupVersion}

	// AddToScheme adds the types in this group-version to the given scheme.
	AddToScheme = SchemeBuilder.AddToScheme
)
