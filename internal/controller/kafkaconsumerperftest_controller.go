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

package controller

import (
	"context"
	"fmt"
	"gopkg.in/yaml.v2"
	"strconv"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	kafkav1alpha1 "github.com/Platformatory/kafka-benchmark-operator/api/v1alpha1"
	"github.com/Platformatory/kafka-benchmark-operator/internal/utils"
)

// KafkaConsumerPerfTestReconciler reconciles a KafkaConsumerPerfTest object
type KafkaConsumerPerfTestReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=kafka.platformatory.io,resources=kafkaconsumerperftests,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=kafka.platformatory.io,resources=kafkaconsumerperftests/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=kafka.platformatory.io,resources=kafkaconsumerperftests/finalizers,verbs=update
//+kubebuilder:rbac:groups=*,resources=secrets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=*,resources=jobs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=*,resources=configmaps,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the KafkaConsumerPerfTest object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.3/pkg/reconcile
func (r *KafkaConsumerPerfTestReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	kafkaConsumerPerfTestSync := &kafkav1alpha1.KafkaConsumerPerfTest{}
	if err := r.Get(ctx, req.NamespacedName, kafkaConsumerPerfTestSync); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	secretName := ""

	if kafkaConsumerPerfTestSync.Spec.ConsumerConfig != nil || kafkaConsumerPerfTestSync.Spec.ConsumerConfigSecretRef.Name == "" {
		if kafkaConsumerPerfTestSync.Spec.ConsumerConfig == nil {
			kafkaConsumerPerfTestSync.Spec.ConsumerConfig = make(map[string]string, 1)
		}
		kafkaConsumerPerfTestSync.Spec.ConsumerConfig["bootstrap.servers"] = kafkaConsumerPerfTestSync.Spec.BootstrapServers
		secretData := map[string][]byte{
			"kafka.properties": []byte(utils.MapToJavaProperties(kafkaConsumerPerfTestSync.Spec.ConsumerConfig))}
		secretName = fmt.Sprintf("%s-%s", kafkaConsumerPerfTestSync.Name, "kafka-consumer-config")

		// Create or update the Secret
		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      secretName,
				Namespace: req.Namespace,
			},
			Data: secretData,
		}
		// Check if Secret exists
		existingSecret := &corev1.Secret{}
		err := r.Get(ctx, types.NamespacedName{Name: secretName, Namespace: req.Namespace}, existingSecret)
		if err != nil && errors.IsNotFound(err) {
			// Secret does not exist, create it
			if err := r.Create(ctx, secret); err != nil {
				logger.Error(err, "Failed to create Secret")
				return ctrl.Result{}, err
			}
			logger.Info("Secret created successfully")
		} else {
			// Secret exists, update it
			existingSecret.Data = secret.Data
			if err := r.Update(ctx, existingSecret); err != nil {
				logger.Error(err, "Failed to update Secret")
				return ctrl.Result{}, err
			}
			logger.Info("Secret updated successfully")
		}
	} else {
		secretName = kafkaConsumerPerfTestSync.Spec.ConsumerConfigSecretRef.Name
	}

	prometheusConfig := map[string]interface{}{
		"global": map[string]string{
			"scrape_interval":     "15s",
			"evaluation_interval": "15s",
		},
		"scrape_configs": []map[string]interface{}{
			{
				"job_name": "jmx",
				"static_configs": []map[string]interface{}{
					{
						"targets": []string{"localhost:7071"},
						"labels": map[string]string{
							"env": kafkaConsumerPerfTestSync.Name,
						},
					},
					{
						"targets": kafkaConsumerPerfTestSync.Spec.MetricsCollector.JMXPrometheusURLs,
						"labels": map[string]string{
							"env": kafkaConsumerPerfTestSync.Name,
						},
					},
				},
				"relabel_configs": []map[string]interface{}{
					{
						"source_labels": []string{"__address__"},
						"target_label":  "hostname",
						"regex":         "([^:]+)(:[0-9]+)?",
						"replacement":   "${1}",
					},
				},
			},
		},
		"remote_write": kafkaConsumerPerfTestSync.Spec.MetricsCollector.Config.RemoteWrite,
	}

	prometheusYAMLBytes, err := yaml.Marshal(prometheusConfig)
	if err != nil {
		logger.Error(err, "Failed to marshall prometheus config to YAML")
		return ctrl.Result{}, err
	}

	configMapData := string(prometheusYAMLBytes)

	configMapName := fmt.Sprintf("%s-%s", kafkaConsumerPerfTestSync.Name, "prometheus-config")

	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: req.Namespace,
			Name:      configMapName,
		},
		Data: map[string]string{
			"prometheus.yml": configMapData,
		},
	}

	// Check if ConfigMap exists
	existingConfigMap := &corev1.ConfigMap{}
	err = r.Get(ctx, types.NamespacedName{Name: configMapName, Namespace: req.Namespace}, existingConfigMap)
	if err != nil && errors.IsNotFound(err) {
		// ConfigMap does not exist, create it
		if err := r.Create(ctx, configMap); err != nil {
			logger.Error(err, "Failed to create ConfigMap")
			return ctrl.Result{}, err
		}
		logger.Info("ConfigMap created successfully")
	} else {
		// ConfigMap exists, update it
		existingConfigMap.Data = configMap.Data
		if err := r.Update(ctx, existingConfigMap); err != nil {
			logger.Error(err, "Failed to update ConfigMap")
			return ctrl.Result{}, err
		}
		logger.Info("ConfigMap updated successfully")
	}

	sleep_time := 60

	if kafkaConsumerPerfTestSync.Spec.MetricsCollector.Provider == "prometheus" {
		for _, remote_url := range kafkaConsumerPerfTestSync.Spec.MetricsCollector.Config.RemoteWrite {
			if remote_url.Metadata_config != (kafkav1alpha1.RemoteWriteMetadataConfig{}) {
				send_interval, _ := time.ParseDuration(remote_url.Metadata_config.Send_interval)
				if int(send_interval.Seconds()) > sleep_time {
					sleep_time = int(send_interval.Seconds())
				}
			}
		}
	}

	sleep_time = sleep_time * 3

	consumer_job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      kafkaConsumerPerfTestSync.Name + "-consumer-perf",
			Namespace: req.Namespace,
		},
		Spec: batchv1.JobSpec{
			Completions: &kafkaConsumerPerfTestSync.Spec.Count,
			Parallelism: &kafkaConsumerPerfTestSync.Spec.Count,
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "consumer",
							Image: kafkaConsumerPerfTestSync.Spec.Image,
							Command: []string{
								"/bin/sh",
								"-c",
								"(./prometheus --enable-feature=agent --config.file=\"/prom/prometheus.yml\" --log.level=error &) && " +
									"kafka-consumer-perf-test --topic " + kafkaConsumerPerfTestSync.Spec.Topic.Name +
									" --bootstrap-server " + kafkaConsumerPerfTestSync.Spec.BootstrapServers +
									" --messages " + strconv.Itoa(int(kafkaConsumerPerfTestSync.Spec.ConsumerPerfParams.MessagesCount)) +
									" --timeout " + strconv.Itoa(int(kafkaConsumerPerfTestSync.Spec.ConsumerPerfParams.Timeout)) +
									" group.id=" + kafkaConsumerPerfTestSync.Spec.ConsumerPerfParams.GroupID +
									" --consumer.config /mnt/kafka.properties && " +
									fmt.Sprintf(`sleep %d`, sleep_time),
							},
							Env: []corev1.EnvVar{
								{
									Name: "KAFKA_OPTS",
									Value: "-javaagent:/usr/app/jmx_prometheus_javaagent-0.15.0.jar=7071:" +
										"/usr/app/kafka_client.yml",
								},
								{
									Name: "POD_IP",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "status.podIP",
										},
									},
								},
								{
									Name: "POD_NAME",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "metadata.name",
										},
									},
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "kafka-properties",
									MountPath: "/mnt",
								},
								{
									Name:      "prometheus-config",
									MountPath: "/prom",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "kafka-properties",
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName: secretName,
								},
							},
						},
						{
							Name: "prometheus-config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: configMapName,
									},
								},
							},
						},
					},
					RestartPolicy: "Never",
				},
			},
			BackoffLimit: utils.GetIntPointer(4),
		},
	}

	if kafkaConsumerPerfTestSync.Spec.Topic.AutoCreate == true {
		consumer_job.Spec.Template.Spec.InitContainers = []corev1.Container{
			{
				Name:  "topics",
				Image: kafkaConsumerPerfTestSync.Spec.Image,
				Command: []string{
					"/bin/sh",
					"-c",
					"kafka-topics --if-not-exists --topic " + kafkaConsumerPerfTestSync.Spec.Topic.Name + " --create --bootstrap-server " +
						kafkaConsumerPerfTestSync.Spec.BootstrapServers + " --replication-factor " +
						strconv.Itoa(int(kafkaConsumerPerfTestSync.Spec.Topic.ReplicationFactor)) + " --partitions " +
						strconv.Itoa(int(kafkaConsumerPerfTestSync.Spec.Topic.Partitions)) + " --command-config /mnt/kafka.properties",
				},
				VolumeMounts: []corev1.VolumeMount{
					{
						Name:      "kafka-properties",
						MountPath: "/mnt",
					},
				},
			},
		}
	}

	// TODO: Check if job exists and update job accordingly
	if err := r.Create(ctx, consumer_job); err != nil {
		logger.Error(err, "Failed to reconcile Job for KafkaConsumerPerf")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *KafkaConsumerPerfTestReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kafkav1alpha1.KafkaConsumerPerfTest{}).
		Complete(r)
}
