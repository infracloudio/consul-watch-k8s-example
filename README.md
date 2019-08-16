# consul-watch-k8s-example

## Using Consul Watch to track config changes in Consul KV store

This example provides insight into Consul watches that uses http handlers to track and notify config changes stored in Consul KV store.

### Prerequisite

- A Kubernetes cluster deployed on minikube or any of the Cloud Providers
- Installed Helm on the Kubernetes Cluster (2.10+)

## Steps for creating the environment

1. Clone the repository
```
git clone https://github.com/rutu-k/consul-watch-k8s-example.git
```

2. Deploy Consul on Kubernetes Cluster
```
cd consul-watch-k8s-example
```
```
helm install --release myconsul ./consul_helm
```

- Once up and running, use port-forward to access the Consul dashboard at http://localhost:8500
```
kubectl port-forward myconsul-consul-server-0 8500:8500
```

- Check the Consul server details using its CLI
```
kubectl exec -it myconsul-consul-server-0 -- /bin/sh
```
```
consul members
```

3. Change/Add configs to Consul KV store

- Configs can be added to the KV store using dashboard or CLI
- Access the Dashboard and go to `Key/Value`  and create KV named `mykey1` in json format

4. Deploy Application pod
```
kubectl apply -f myconsul-app-client.yaml
```
- Note: This will deploy two containers in a single pod, one with consul agent and other with the app itself


### Verfication
- Application at initialization will read config values from Consul KV store and print on its CLI.
- Consul watch is configured to check for changes in the specified KV. Whenever there's a change in KV, Consul watch will notify the application via http handler specified in Consul config. Application will delete the old values of KV and start printing the changed Kv values.
