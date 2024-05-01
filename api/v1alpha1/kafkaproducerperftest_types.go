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

type ProducerConfigSecretRef struct {
	Name string `json:"name"`
}

type ProducerPerfParams struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=5000
	RecordsCount int32 `json:"recordsCount"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=512000
	RecordSizeBytes int32 `json:"recordSizeBytes"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=-1
	Throughput int32 `json:"throughput"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=-1
	Acks int32 `json:"acks"`
}

type Topic struct {
	Name string `json:"name"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=true
	AutoCreate bool `json:"autoCreate"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=-1
	ReplicationFactor int32 `json:"replicas"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=1
	Partitions int32 `json:"partitions"`
}

type RemoteWriteMetadataConfig struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=true
	Send bool `json:"send"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:="1m"
	Send_interval string `json:"send_interval"`
	// +kubebuilder:validation:Optional
	Max_samples_per_send int `json:"max_samples_per_send"`
}

type RemoteWrite struct {
	URL string `json:"url"`
	// +kubebuilder:validation:Optional
	Bearer_Token string `json:"bearer_token,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:={send: true, send_interval: "15s", max_samples_per_send: 500}
	Metadata_config RemoteWriteMetadataConfig `json:"metadata_config"`
}

type PrometheusConfig struct {
	RemoteWrite []RemoteWrite `json:"remoteWrite"`
}

type MetricsCollector struct {
	// +kubebuilder:validation:Enum=prometheus
	Provider string `json:"provider"`
	// +kubebuilder:validation:Optional
	JMXPrometheusURLs []string         `json:"jmxPrometheusURLs,omitempty"`
	Config            PrometheusConfig `json:"config"`
}

type LabelSelector struct {
	MatchLabels map[string]string `json:"matchLabels"`
}

type TopologySpreadConstraint struct {
	TopologyKey string `json:"topologyKey"`
	// +kubebuilder:validation:Optional
	MaxSkew int32 `json:"maxSkew"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=DoNotSchedule;ScheduleAnyway
	WhenUnsatisfiable string        `json:"whenUnsatisfiable"`
	LabelSelector     LabelSelector `json:"labelSelector"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:="Ignore"
	// +kubebuilder:validation:Enum=Honor;Ignore
	NodeAffinityPolicy string `json:"nodeAffinityPolicy"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:="Ignore"
	// +kubebuilder:validation:Enum=Honor;Ignore
	NodeTaintsPolicy string `json:"nodeTaintsPolicy"`
}

// KafkaProducerPerfTestSpec defines the desired state of KafkaProducerPerfTest
type KafkaProducerPerfTestSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Count int32 `json:"count"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:="platformatorylabs/kafka-performance-suite:1.0.1"
	Image            string `json:"image,omitempty"`
	BootstrapServers string `json:"bootstrapServers"`
	// +kubebuilder:validation:Optional
	ProducerConfig map[string]string `json:"producerConfig,omitempty"`
	// +kubebuilder:validation:Optional
	ProducerConfigSecretRef ProducerConfigSecretRef `json:"producerConfigSecretRef,omitempty"`
	ProducerPerfParams      ProducerPerfParams      `json:"producerPerfParams"`
	Topic                   Topic                   `json:"topic"`
	MetricsCollector        MetricsCollector        `json:"metricsCollector"`
	// +kubebuilder:validation:Optional
	TopologySpreadConstraints []TopologySpreadConstraint `json:"topologySpreadConstraints,omitempty"`
}

// KafkaProducerPerfTestStatus defines the observed state of KafkaProducerPerfTest
type KafkaProducerPerfTestStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// KafkaProducerPerfTest is the Schema for the kafkaproducerperftests API
type KafkaProducerPerfTest struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KafkaProducerPerfTestSpec   `json:"spec,omitempty"`
	Status KafkaProducerPerfTestStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// KafkaProducerPerfTestList contains a list of KafkaProducerPerfTest
type KafkaProducerPerfTestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KafkaProducerPerfTest `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KafkaProducerPerfTest{}, &KafkaProducerPerfTestList{})
}
