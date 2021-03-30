/*
Copyright 2021.

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
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// HostAndPath defines the hostname and filepath info
type HostAndPath struct {
	Hostname string `json:"hostname"`
	FilePath string `json:"filePath"`
}

// NatsCoSpec defines the desired state of NatsCo
type NatsCoSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	CoType       string        `json:"coType,omitempty"`
	NatsServers  []string      `json:"natsServer,omitempty"`
	Source       HostAndPath   `json:"source"`
	Destinations []HostAndPath `json:"destinations"`
}

// NatsCoStatus defines the observed state of NatsCo
type NatsCoStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Source      string            `json:"source"`
	Destination map[string]string `json:"destination"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// NatsCo is the Schema for the natscoes API
type NatsCo struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NatsCoSpec   `json:"spec,omitempty"`
	Status NatsCoStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// NatsCoList contains a list of NatsCo
type NatsCoList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NatsCo `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NatsCo{}, &NatsCoList{})
}
