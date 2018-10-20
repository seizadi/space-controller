/*
Copyright 2017 The Kubernetes Authors.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/api/core/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Space is a specification for a SpaceSpec resource
type Space struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	
	Spec   SpaceSpec   `json:"spec"`
	Status SpaceStatus `json:"status"`
}

// SpaceSpec is the spec for a Space resource
type SpaceSpec struct {
	Path       string            `json:"path"`
	SecretName string            `json:"secretName"`
	Type       corev1.SecretType `json:"type"`
	Secrets    map[string]string `json:"secrets"`
}

// SpaceStatus is the status for a Space resource
type SpaceStatus struct {
	AvailableSecrets int32 `json:"availablesecrets"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SpaceList is a list of Space resources
type SpaceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Space `json:"items"`
}
