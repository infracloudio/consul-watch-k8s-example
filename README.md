# consul-watch-k8s-example

## Using Consul Watch to track config changes in Consul KV (Key/Value) store

Tracking KV changes using watches can have two approaches:
a. Configuring Watches on the Consul Agent that is running as Daemonset.
b. Configuring Watches on the Consul Agent that is running as a sidecar.

## a. Configuring watches on the Consul Agent running as a Daemonset.

In this scenario, Consul Agent runs as a Daemonset on each node. Consul Agent is configured with watches that tracks KV changes for the Apps/Pods located on that particular node. Generally, watches are configured using CLI command `consul watch` or using json file placed in config directory of Consul Agent. Consider these watches use http handler to trigger notification of KV changes. Each watch will require an IP address of respective App/Pods to notify the KV changes. This scenario works fine until all the Pods are running perfectly. But, as Pods are ephemeral, there is a possibility of Pod going down. In such a case, this scenario fails to work as intended.

Suppose node 1 has Pod A for which a watch is configured to track KV-A and it uses http handler pointing to the App/Pod's IP (say 1.1.1.1) to notify KV changes. But if Pod A terminates and is scheduled on node 2 with a new IP in such a case, watch configured on node 1 will not be able to trigger or connect to Pod A due to changed IP. Moreover, a new watch must be configured in the Consul Agent of node 2 with the changed IP (say 2.2.2.2). This cannot be done dynamically as it requires knowledge of Pod’s existence, the KV that Pod needs to track and IP of Pod. In addition to this, Consul Agent must also be reloaded to incorporate new config. Reloading the Consul Agent may hamper the existing Consul processes. Hence, it is impractical to configure watches on the fly as it is not a straightforward procedure that can be followed to auto-configure watches in such scenarios.

The problem intensifies if another Pod B which tracks changes of KV-B for which the watch is configured by Consul Agent spawns up with previous IP of Pod A (say 1.1.1.1). This will create false triggers for Pod B. Problem worsens if Pod B is acquiring changed KV from the http handler itself (by using http POST methd). This can be.
Thus, to overcome theses issues, we consider the second approach i.e. sidecar approach.

## b. Configuring Watches on the Consul Agent that is running as a sidecar.

In this scenario, Consul Agent runs as a sidecar in each Pod. Watch is configured on this Pod to track changes of the KV and notify subsequently using http handler. In this case, http handler is pointed to the localhost in contrast to the Daemonset approach where handler is pointed to IP address. Thus, whenever, the Pod terminates or goes down a new Pod is created with same configurations (assuming it is running as deployment) that joins the Consul cluster and start tracking the same KV changes again. Steps to create this scanario are given below.

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

4. Deploy Application Pod
```
kubectl apply -f consul-app-client.yaml
```
- Note: This will deploy two containers in a single Pod, one with consul agent and other with the app itself


### Verfication
- Application at initialization will read config values from Consul KV store and print on its CLI.
- Consul watch is configured to check for changes in the specified KV. Whenever there's a change in KV, Consul watch will notify the application via http handler specified in Consul config. Application will delete the old values of KV and start printing the changed KV values.
