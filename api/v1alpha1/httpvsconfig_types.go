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

// HttpvsConfigSpec defines the desired state of HttpvsConfig
type HttpvsConfigSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of HttpvsConfig. Edit httpvsconfig_types.go to remove/update
	Namespace string   `json:"namespace,omitempty"`
	Hosts     []string `json:"hosts,omitempty"`
	Route     []Route  `json:"route,omitempty"`
}
type Route struct {
	Service string   `json:"service,omitempty"`
	Port    uint32   `json:"port,omitempty"`
	Uri     []string `json:"uri,omitempty"`
}

// HttpvsConfigStatus defines the observed state of HttpvsConfig
type HttpvsConfigStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// HttpvsConfig is the Schema for the httpvsconfigs API
type HttpvsConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HttpvsConfigSpec   `json:"spec,omitempty"`
	Status HttpvsConfigStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// HttpvsConfigList contains a list of HttpvsConfig
type HttpvsConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HttpvsConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&HttpvsConfig{}, &HttpvsConfigList{})
}
