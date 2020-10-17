---
title: Attaching your local machine to a Kubernetes cluster
date: 2020-10-17
---

Why on earth would you want to do that?! You may have a legitimate reason like
pooling spare compute capacity. Oooooor you're just goofing around in order to
understand things. Like me.

Specifically I was playing with [telepresence](https://www.telepresence.io/) (a
"VPN into the cluster") and a had an idea: why not just run a part of the
cluster locally. Then you can interact with the containers and containers can
interact with your host (host port, host path etc), no magic required.

# Setting up the cluster

I'm a big fan of [K3s](https://k3s.io/) (a lightweight kubernetes distribution)
for small clusters as it's very easy to set up and doesn't require much
resources to run. Perfect if I want to run part of it on my machine.

So for purposes of this testing I setup a single node remote cluster in the
cloud. Specifically I provisioned a VM on [Digital
Ocean](https://www.digitalocean.com/) and ran this on it

```bash
curl -sfL https://get.k3s.io | sh -
k3s kubectl get node # to confirm we're up
```

Then copied `/etc/rancher/k3s/k3s.yaml` to my machine and updated the IP to
point to the VM's external IP, point `KUBECONFIG` env variable to this file and
we're connected! Locally I can now also run o

```bash
kubectl get node
```

# Joining the cluster

K3s makes joining a worker node (an agent as they call it) super easy, you just
need to point it to the K8s API server and provide an authentication token.

We can get the token by reading `/var/lib/rancher/k3s/server/node-token` on the
VM.

And since I'm on Linux I can actually install K3s locally natively (server and
token values are placeholders here.

```bash
curl -sfL https://get.k3s.io | K3S_URL=https://myserver:6443 K3S_TOKEN=mynodetoken sh -
```

Et voila, we're now part of the cluster. This does make a bit of a mess on the
local machine so instead we can run it in a docker container.

```
... set TOKEN & SERVER_IP
docker run --rm -it --name=k3s-agent --privileged  rancher/k3s:v1.18.9-k3s1 agent --token "$TOKEN" --server "https://${SERVER_IP}:6443" --node-name=amaterasu
```

but we can do better, we can actually integrate with host docker so pods and
containers will show up as containers on host!

```
docker run --rm -it --name=k3s-agent --privileged -v /var/run/docker.sock:/var/run/docker.sock rancher/k3s:v1.18.9-k3s1 agent --token "$TOKEN" --server "https://${SERVER_IP}:6443" --node-name=amaterasu --docker
```

Note I'm using hostname `amaterasu` here because futher configs will refer to
it.

# Using the cluster

Let's run some workloads!

```bash
kubectl create deployment hello-node --image=k8s.gcr.io/echoserver:1.4
kubectl get pod -o wide
```

And can now connect to it

```bash
kubectl port-forward  pod/hello-node-7bf657c596-dgtmk 8080
http://localhost:8080/
```

AND see it running as containers on my local docker

```bah
docker ps
```

How about running some workloads that are _always_ on my local machine? I can
use a node selector for that!

First I label my node

```bash
kubectl label nodes amaterasu name=amaterasu
```

And then I can use this label as a selector for a pod that specifically
leverages something available only on my node. E.g. mounting a host path volume.

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: local-pod
spec:
  containers:
  - image: phusion/baseimage:0.11
    name: test-container
    volumeMounts:
    - mountPath: /host
      name: test-volume
  volumes:
  - name: test-volume
    hostPath:
      path: /tmp/shared
      type: Directory
  nodeSelector:
    name: amaterasu
```

And deploy with:

```
kubectl create -f ./local.yaml
kubectl get pod -o wide
```

Now I can get a shell in there with `kubectl exec -it local-pod -- bash` or I
can use docker directly:

```
docker exec -it k8s_test-container_local-pod_default_1e9ecfbb-6b03-493e-94c8-82b7dbb8bc13_0 bash
```

Once in I can confirm I've actually mounted a host path:

```bash
cd /host
date > now
```

And check the contents on my host machine:

```bash
cat /tmp/shared/now
```

Works like magic :)

## Avoiding other workloads

If you're running multiple workloads then the scheduler may decide to use your
local machine for them. And you might not like that as you only want to run your
specific workloads.

You can actually tell that to the scheduler by tainting your node with a
`NoExecute` flag.

```bash
kubectl taint nodes amaterasu personal=amaterasu:NoExecute
kubectl get pod -o wide
```

In my case pod `hello-node*` was actually running on my local machine and was
now killed and subsequently restarted on the remote cluster as it's managed by
it's parent deployment.

But `local-pod` is gone too!

We also need to specify a toleration so it will actually schedule despite the
taint.


```yaml
apiVersion: v1
kind: Pod
metadata:
  name: local-pod
spec:
  containers:
  - image: phusion/baseimage:0.11
    name: test-container
    volumeMounts:
    - mountPath: /host
      name: test-volume
  volumes:
  - name: test-volume
    hostPath:
      path: /tmp/shared
      type: Directory
  nodeSelector:
    name: amaterasu
  tolerations:
  - key: "personal"
    operator: "Equal"
    value: "amaterasu"
    effect: "NoExecute"
```

And repeat the deployment.

There you have it: natively running k8s workloads locally with local access and
full integration with the remote cluster. Not sure if it's a good idea though ;)

