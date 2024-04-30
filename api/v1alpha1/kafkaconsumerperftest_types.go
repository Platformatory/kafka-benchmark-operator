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

type ConsumerConfigSecretRef struct {
	Name string `json:"name"`
}

type ConsumerPerfParams struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=10000
	MessagesCount int32 `json:"messagesCount"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=100000
	Timeout int32 `json:"timeout"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:="$POD_NAME"
	GroupID string `json:"groupID"`
}

// KafkaConsumerPerfTestSpec defines the desired state of KafkaConsumerPerfTest
type KafkaConsumerPerfTestSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Count int32 `json:"count"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:="platformatorylabs/kafka-performance-suite:1.0.1"
	Image            string `json:"image,omitempty"`
	BootstrapServers string `json:"bootstrapServers"`
	// +kubebuilder:validation:Optional
	ConsumerConfig map[string]string `json:"consumerConfig,omitempty"`
	// +kubebuilder:validation:Optional
	ConsumerConfigSecretRef ConsumerConfigSecretRef `json:"consumerConfigSecretRef,omitempty"`
	ConsumerPerfParams      ConsumerPerfParams      `json:"consumerPerfParams"`
	Topic                   Topic                   `json:"topic"`
	MetricsCollector        MetricsCollector        `json:"metricsCollector"`
	// +kubebuilder:validation:Optional
	TopologySpreadConstraints []TopologySpreadConstraint `json:"topologySpreadConstraints,omitempty"`
}

// KafkaConsumerPerfTestStatus defines the observed state of KafkaConsumerPerfTest
type KafkaConsumerPerfTestStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// KafkaConsumerPerfTest is the Schema for the kafkaconsumerperftests API
type KafkaConsumerPerfTest struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KafkaConsumerPerfTestSpec   `json:"spec,omitempty"`
	Status KafkaConsumerPerfTestStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// KafkaConsumerPerfTestList contains a list of KafkaConsumerPerfTest
type KafkaConsumerPerfTestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KafkaConsumerPerfTest `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KafkaConsumerPerfTest{}, &KafkaConsumerPerfTestList{})
}
