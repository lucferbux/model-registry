[Model registry server set up]: ../../bff/docs/dev-guide.md

## Deploying the Model Registry UI in a local cluster

For this guide, we will be using kind for locally deploying our cluster. See
the [Model registry server set up] guide for prerequisites on setting up kind 
and deploying the model registry server.

### Setup
#### 1. Create a kind cluster
Create a local cluster for running the MR UI using the following command:
```shell
kind create cluster
```

#### 2. Create kubeflow namespace
Create a namespace for model registry to run in, by default this is kubeflow, run:
```shell
kubectl create namespace kubeflow
```

#### 3. Deploy Model Registry UI to cluster
You can now deploy the UI and BFF to your newly created cluster using the kustomize configs in this directory:
```shell
cd clients/ui

kubectl apply -k manifests/base/ -n kubeflow
```

After a few seconds you should see 2 pods running (1 for BFF and 1 for UI):
```shell
kubectl get pods -n kubeflow
```
```
NAME                                  READY   STATUS    RESTARTS   AGE
model-registry-bff-746f674b99-bfvgs   1/1     Running   0          11s
model-registry-ui-58755c4754-zdrnr    1/1     Running   0          11s
```

#### 4. Access the Model Registry UI running in the cluster
Now that the pods are up and running you can access the UI.

First you will need to port-forward the UI service by running the following in it's own terminal:
```shell
kubectl port-forward service/model-registry-ui-service 8080:8080 -n kubeflow
```

You can then access the UI running in your cluster locally at http://localhost:8080/

To test the BFF separately you can also port-forward that service by running:
```shell
kubectl port-forward service/model-registry-bff-service 4000:4000 -n kubeflow
```

You can now make API requests to the BFF endpoints like:
```shell
curl http://localhost:4000/api/v1/model-registry
```
```
{
    "model_registry": null
}
```

### Troubleshooting

#### Running on macOS
When running locally on macOS you may find the pods fail to deploy, with one or more stuck in the `pending` state. This is usually due to insufficient memory allocated to your docker / podman virtual machine. You can verify this by running:
```shell
kubectl describe pods -n kubeflow
```
If you're experiencing this issue you'll see an output containing something similar to the following:
```
Events:
  Type     Reason            Age   From               Message
  ----     ------            ----  ----               -------
  Warning  FailedScheduling  29s   default-scheduler  0/1 nodes are available: 1 Insufficient memory. preemption: 0/1 nodes are available: 1 No preemption victims found for incoming pod.
```

To fix this, you'll need to increase the amount of memory available to the VM. This can be done through either the Podman Desktop or Docker Desktop GUI. 6-8GB of memory is generally a sufficient amount to use.

## Running with Kubeflow and Istio
Alternatively, if you'd like to run the UI and BFF pods with an Istio configuration for the KF Central Dashboard, you can apply the manifests by running:
```shell
kubectl apply -k overlays/kubeflow -n kubeflow
```