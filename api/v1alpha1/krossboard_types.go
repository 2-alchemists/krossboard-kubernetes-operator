/*
Copyright 2022 2ALCHEMISTS SAS.

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

// KoaInstance defines a state of a kube-opex-analytics instance
type KoaInstance struct {
	Name               string `json:"name,omitempty"`
	ContainerPort      int64  `json:"containerPort,omitempty"`
	ClusterName        string `json:"clusterName,omitempty"`
	ClusterEndpointURL string `json:"clusterEndpoint,omitempty"`
}

// KbComponentInstance defines a the state of a Krossboard component instance
type KbComponentInstance struct {
	Name          string `json:"name,omitempty"`
	ContainerPort int64  `json:"containerPort,omitempty"`
}

// KrossboardSpec defines the desired state of Krossboard
type KrossboardSpec struct {
	// KrossboardUIImage sets Krossboard UI image
	//+kubebuilder:default="krossboard/krossboard-ui:latest"
	KrossboardUIImage string `json:"krossboardUIImage,omitempty"`

	// KrossboardDataProcessorImage sets the image of Krossboard Data Processor
	//+kubebuilder:default="krossboard/krossboard-data-processor:latest"
	KrossboardDataProcessorImage string `json:"krossboardDataProcessorImage,omitempty"`

	// KrossboardPersistentVolumeClaim sets the name of the persitent volume used for data
	//+kubebuilder:default="krossboard/krossboard-data-pvc"
	KrossboardPersistentVolumeClaim string `json:"krossboardPersistentVolumeClaim,omitempty"`

	// KoaImage sets the image of kube-opex-analytics
	//+kubebuilder:default="rchakode/kube-opex-analytics:latest"
	KoaImage string `json:"koaImage,omitempty"`

	// KrossboardSecretName is the name of the secret map for Krossbaord
	//+kubebuilder:default="krossboard-secrets"
	KrossboardSecretName string `json:"krossboardSecretName,omitempty"`

	// useGKEIdentity uses relying on GCP Workload Identity to get GKE credentials
	//+kubebuilder:default=false
	UseGKEIdentity bool `json:"useGKEIdentity,omitempty"`
}

// KrossboardStatus defines the observed state of Krossboard
type KrossboardStatus struct {
	// KoaInstances contains a list of kube-opex-analytics instances
	KoaInstances         []KoaInstance         `json:"koaInstances"`
	KbComponentInstances []KbComponentInstance `json:"kbComponentInstances"`
}

//+kubebuilder:object:root=true

// Krossboard is the Schema for the krossboards API
// +kubebuilder:subresource:status
type Krossboard struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KrossboardSpec   `json:"spec,omitempty"`
	Status KrossboardStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// KrossboardList contains a list of Krossboard
type KrossboardList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Krossboard `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Krossboard{}, &KrossboardList{})
}
