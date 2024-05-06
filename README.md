# kafka-benchmark-operator
The Kafka Benchmark Operator is a Kubernetes operator to run distributed benchmarks for Kafka clusters.

## Description

## Install

```sh
helm repo add platformatory https://platformatory.github.io/kafka-benchmark-operator
helm repo update

helm install kafka-benchmark-operator platformatory/kafka-benchmark-operator \
    -n kafka-benchmark-operator \
    --create-namespace
```

## Try it out

```sh
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo add confluentinc https://packages.confluent.io/helm
helm repo update

helm install kube-prometheus-stack prometheus-community/kube-prometheus-stack -n monitoring --create-namespace --set prometheus.prometheusSpec.enableRemoteWriteReceiver=true

helm upgrade --install \
  confluent-operator confluentinc/confluent-for-kubernetes -n confluent --create-namespace

kubectl apply -f https://raw.githubusercontent.com/confluentinc/confluent-kubernetes-examples/master/quickstart-deploy/confluent-platform-singlenode.yaml -n confluent

# Modify samples in config/samples
kubectl apply -k config/samples
```

## License

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

