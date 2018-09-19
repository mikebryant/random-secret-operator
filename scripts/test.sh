#!/bin/bash

cleanup() {
  set +e
  echo "Clean up..."
  kubectl delete -f examples/resource.yaml
  kubectl delete -f deploy/operator.yaml
  until ! kubectl get namespace randomsecret; do sleep 1; done
}
trap cleanup EXIT

set -ex

echo "Install metacontroller..."
kubectl apply -f https://raw.githubusercontent.com/GoogleCloudPlatform/metacontroller/master/manifests/metacontroller-rbac.yaml
kubectl apply -f https://raw.githubusercontent.com/GoogleCloudPlatform/metacontroller/master/manifests/metacontroller.yaml

echo "Install controller..."
kubectl apply -f deploy/operator.yaml
kubectl create configmap code -n randomsecret --from-file=sync.py

echo "Wait until CRD is available..."
until kubectl get randomsecrets; do sleep 1; done

echo "Create an object..."
kubectl apply -f examples/resource.yaml

echo "Wait for secret..."
until kubectl -n default get secret minimal -o yaml; do sleep 1; done
