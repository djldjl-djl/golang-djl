/*
Copyright 2025.

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

package v1

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// DjlD1Spec defines the desired state of DjlD1.
type DjlD1Spec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Size            *int32           `json:"size"`
	Image           string           `json:"image"`
	ImagePullPolicy v1.PullPolicy    `json:"imagePullPolicy"`
	Ports           []v1.ServicePort `json:"ports"`
	ServerName      string           `json:"serverName"`
}

// DjlD1Status defines the observed state of DjlD1.
type DjlD1Status struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Ready    int32 `json:"ready"`
	Notready int32 `json:"notready"`
}

//通过注解自定义显示字段
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Ready",type=integer,JSONPath=".status.ready",description="Ready Pod count"
// +kubebuilder:printcolumn:name="NotReady",type=integer,JSONPath=".status.notready",description="Not Ready Pod count"
// DjlD1 is the Schema for the djld1s API.

type DjlD1 struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DjlD1Spec   `json:"spec,omitempty"`
	Status DjlD1Status `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// DjlD1List contains a list of DjlD1.
type DjlD1List struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DjlD1 `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DjlD1{}, &DjlD1List{})
}
