#!/usr/bin/env bash

# Check for required tools
command -v docker >/dev/null 2>&1 || { echo >&2 "Docker is required but it's not installed. Aborting."; exit 1; }
command -v kubectl >/dev/null 2>&1 || { echo >&2 "kubectl is required but it's not installed. Aborting."; exit 1; }
command -v kind >/dev/null 2>&1 || { echo >&2 "kind is required but it's not installed. Aborting."; exit 1; }

echo "WARNING: You must have access to a cluster with kubeflow installed."

# Step 1: Deploy Model Registry UI to cluster
pushd  ../../manifests/kustomize/options/ui/overlays/kubeflow
echo "Deploying Model Registry UI..."
kubectl apply -n kubeflow -k .

# Wait for deployment to be available
echo "Waiting Model Registry UI to be available..."
kubectl wait --for=condition=available -n kubeflow deployment/model-registry-ui --timeout=1m

# Step 5: Port-forward the service
echo "Port-forwarding Kubeflow Central Dashboard..."
echo -e "\033[32mDashboard available in http://localhost:8080\033[0m"
kubectl port-forward svc/istio-ingressgateway -n istio-system 8080:80
