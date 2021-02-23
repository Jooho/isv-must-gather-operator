/*
Copyright 2021 Jooho Lee.

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

// MustGatherSpec defines the desired state of MustGather
type MustGatherSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// MustGatherImgURL is the ISV operator must gather image url
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="ISV Operator Must-Gather Image URL",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:string", "urn:alm:descriptor:io.kubernetes:custom"}
	MustGatherImgURL string `json:"mustGatherImgURL,omitempty"`

	//TODO
	//NAMESPACE
	//DEBUG
}

// MustGatherStatus defines the observed state of MustGather
type MustGatherStatus struct {
	//DownloadURL is the endpoint to access downloag web page.
	DownloadURL string `json:"downloadURL,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// MustGather is the Schema for the mustgathers API
type MustGather struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MustGatherSpec   `json:"spec,omitempty"`
	Status MustGatherStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MustGatherList contains a list of MustGather
type MustGatherList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MustGather `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MustGather{}, &MustGatherList{})
}
