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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type ClientConfigSecretRef struct {
	Name string `json:"name"`
}

type E2EPerfParams struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=500000
	RecordsCount int32 `json:"recordsCount"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=512000
	RecordSizeBytes int32 `json:"recordSizeBytes"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:="all"
	Acks string `json:"acks"`
}

// KafkaE2EPerfTestSpec defines the desired state of KafkaE2EPerfTest
type KafkaE2EPerfTestSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Count int32 `json:"count"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:="platformatorylabs/kafka-performance-suite:1.0.1"
	Image            string `json:"image,omitempty"`
	BootstrapServers string `json:"bootstrapServers"`
	// +kubebuilder:validation:Optional
	ClientConfig map[string]string `json:"clientConfig,omitempty"`
	// +kubebuilder:validation:Optional
	ClientConfigSecretRef ClientConfigSecretRef `json:"clientConfigSecretRef,omitempty"`
	E2EPerfParams         E2EPerfParams         `json:"e2ePerfParams"`
	Topic                 Topic                 `json:"topic"`
	MetricsCollector      MetricsCollector      `json:"metricsCollector"`
	// +kubebuilder:validation:Optional
	TopologySpreadConstraints []TopologySpreadConstraint `json:"topologySpreadConstraints,omitempty"`
}

// KafkaE2EPerfTestStatus defines the observed state of KafkaE2EPerfTest
type KafkaE2EPerfTestStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// KafkaE2EPerfTest is the Schema for the kafkae2eperftests API
type KafkaE2EPerfTest struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KafkaE2EPerfTestSpec   `json:"spec,omitempty"`
	Status KafkaE2EPerfTestStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// KafkaE2EPerfTestList contains a list of KafkaE2EPerfTest
type KafkaE2EPerfTestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KafkaE2EPerfTest `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KafkaE2EPerfTest{}, &KafkaE2EPerfTestList{})
}
