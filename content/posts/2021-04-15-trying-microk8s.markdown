---
title: Trying out Microk8s
date: 2021-04-15
tags: ["trying-out"]
---

In another "running kuberets" flavoured post I'll be trying out
[microk8s](https://microk8s.io/)

> High availability K8s.
> Low-ops, minimal production Kubernetes,
> for devs, cloud, clusters, workstations, Edge and IoT.

In the past I've been using [k3s](https://k3s.io/) for small (mostly single
node) "clusters". The pitch is very similar. The few differences I can point
out now are:

- microk8s is a bit heavier - I don't think it will run on 1G vps, let's try it
  out
- microk8s has a better HA story - definately want to try this out
- k3s comes as a single "run everywhere" binary whereas microk8s requires
  [snap](https://snapcraft.io/docs/installing-snapd)

# Setting up

I'll be using [Hetzner Cloud](https://www.hetzner.com/cloud) to provision a few
small VMs for my laboratory. But I'll start with a single one.

Specs: Ubuntu 20.04 (so we get snap out of the box), 40GB local NVMe storage, 2x
vCPU, 4GB ram. And a private network so I can connect future instances to it as
well. A minute later I'm in

```
$ ssh root@162.55.52.75
Last login: Wed Apr 14 10:14:55 2021 from 46.123.240.50
root@microk8s-1:~#
```

Now let's follow the instructions on microk8s homepage

```
root@microk8s-1:~# snap install microk8s --classic

Command 'snap' not found, but can be installed with:

apt install snapd
```

Ok, apparently I need to still install snap first.

```
apt update
...
apt install -y snapd
...
```

# Installation and setup

Here we go again

```
root@microk8s-1:~# snap install microk8s --classic
...
2021-04-14T10:18:14+02:00 INFO Waiting for automatic snapd restart...
microk8s (1.20/stable) v1.20.5 from Canonical✓ installed
```

Step 2: wait for it :)

```
microk8s status --wait-ready

microk8s is running
high-availability: no
  datastore master nodes: 127.0.0.1:19001
  datastore standby nodes: none
addons:
  enabled:
    ha-cluster           # Configure high availability on the current node
  disabled:
    ambassador           # Ambassador API Gateway and Ingress
    cilium               # SDN, fast with full network policy
    dashboard            # The Kubernetes dashboard
    dns                  # CoreDNS
    fluentd              # Elasticsearch-Fluentd-Kibana logging and monitoring
    gpu                  # Automatic enablement of Nvidia CUDA
    helm                 # Helm 2 - the package manager for Kubernetes
    helm3                # Helm 3 - Kubernetes package manager
    host-access          # Allow Pods connecting to Host services smoothly
    ingress              # Ingress controller for external access
    istio                # Core Istio service mesh services
    jaeger               # Kubernetes Jaeger operator with its simple config
    keda                 # Kubernetes-based Event Driven Autoscaling
    knative              # The Knative framework on Kubernetes.
    kubeflow             # Kubeflow for easy ML deployments
    linkerd              # Linkerd is a service mesh for Kubernetes and other frameworks
    metallb              # Loadbalancer for your Kubernetes cluster
    metrics-server       # K8s Metrics Server for API access to service metrics
    multus               # Multus CNI enables attaching multiple network interfaces to pods
    portainer            # Portainer UI for your Kubernetes cluster
    prometheus           # Prometheus operator for monitoring and logging
    rbac                 # Role-Based Access Control for authorisation
    registry             # Private image registry exposed on localhost:32000
    storage              # Storage class; allocates storage from host directory
    traefik              # traefik Ingress controller for external access
```

Well apparently this worked :D

Step 3: enable services

```
root@microk8s-1:~# microk8s enable dashboard dns helm3 metrics-server registry
storage traefik
Enabling Kubernetes Dashboard
Enabling Metrics-Server
...
clusterrolebinding.rbac.authorization.k8s.io/traefik-ingress-controller created
service/traefik-web-ui created
ingress.networking.k8s.io/traefik-web-ui created
traefik ingress controller has been installed on port 8080
```

I must say this already feels much smoother than k3s setup.

```
microk8s is running
high-availability: no
  datastore master nodes: 127.0.0.1:19001
  datastore standby nodes: none
addons:
  enabled:
    dashboard            # The Kubernetes dashboard
    dns                  # CoreDNS
    ha-cluster           # Configure high availability on the current node
    helm3                # Helm 3 - Kubernetes package manager
    metrics-server       # K8s Metrics Server for API access to service metrics
    registry             # Private image registry exposed on localhost:32000
    storage              # Storage class; allocates storage from host directory
    traefik              # traefik Ingress controller for external access
  disabled:
    ambassador           # Ambassador API Gateway and Ingress
    cilium               # SDN, fast with full network policy
    fluentd              # Elasticsearch-Fluentd-Kibana logging and monitoring
    gpu                  # Automatic enablement of Nvidia CUDA
    helm                 # Helm 2 - the package manager for Kubernetes
    host-access          # Allow Pods connecting to Host services smoothly
    ingress              # Ingress controller for external access
    istio                # Core Istio service mesh services
    jaeger               # Kubernetes Jaeger operator with its simple config
    keda                 # Kubernetes-based Event Driven Autoscaling
    knative              # The Knative framework on Kubernetes.
    kubeflow             # Kubeflow for easy ML deployments
    linkerd              # Linkerd is a service mesh for Kubernetes and other frameworks
    metallb              # Loadbalancer for your Kubernetes cluster
    multus               # Multus CNI enables attaching multiple network interfaces to pods
    portainer            # Portainer UI for your Kubernetes cluster
    prometheus           # Prometheus operator for monitoring and logging
    rbac                 # Role-Based Access Control for authorisation
```

# Using

Step 4: start using (like k3s there is a `kubectl` packaged in)

```
root@microk8s-1:~# microk8s kubectl get namespaces
NAME                 STATUS   AGE
kube-system          Active   7m5s
kube-public          Active   7m4s
kube-node-lease      Active   7m4s
default              Active   7m4s
container-registry   Active   3m43s
traefik              Active   3m41s
root@microk8s-1:~# microk8s kubectl get nodes
NAME         STATUS   ROLES    AGE    VERSION
microk8s-1   Ready    <none>   7m8s   v1.20.5-34+40f5951bd9888a
```

Now how do I connect to this from my laptop? Quick search says microk8s config`
will outtput a kubeconfig file. And indeed it works

```
andraz@amaterasu /tmp/temp
 $ export KUBECONFIG=$(pwd)/kubeconfig
andraz@amaterasu /tmp/temp
 $ kubectl get nodes
NAME         STATUS   ROLES    AGE   VERSION
microk8s-1   Ready    <none>   12m   v1.20.5-34+40f5951bd9888a
```

Can we get a dashboard?

```

andraz@amaterasu /tmp/temp
 $ k get service -A | grep dashboard
kube-system          kubernetes-dashboard        ClusterIP   10.152.183.220   <none>        443/TCP                  12m
kube-system          dashboard-metrics-scraper   ClusterIP   10.152.183.160   <none>        8000/TCP                 12m
andraz@amaterasu /tmp/temp
 $ k proxy
Starting to serve on 127.0.0.1:8001
```

And indded I get the Web UI available [on
localhost](http://localhost:8001/api/v1/namespaces/kube-system/services/https:kubernetes-dashboard:/proxy/#/login)

How about resource usage?

```
root@microk8s-1:~# free -m
              total        used        free      shared  buff/cache   available
Mem:           3840        1154         331           1        2354        2524
Swap:             0           0           0
```

microk8s comes in at about 1.1G. Slightly heavier than [k3](http://k3s.io/)
(around 0.8G) but lighter than [RKE](https://rancher.com/products/rke/) (around
1.7G).

# Deploying workloads

Let's try to get a demo server up and running. I'll be using `nginx-hello` from
[nginx-demos](https://github.com/nginxinc/NGINX-Demos)


```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: hello
  labels:
    app: hello
spec:
  replicas: 3
  selector:
    matchLabels:
      app: hello
  template:
    metadata:
      labels:
        app: hello
    spec:
      containers:
      - name: nginx
        image: nginxdemos/hello:plain-text
        ports:
        - containerPort: 80
---
apiVersion: v1
kind: Service
metadata:
  name: hello
spec:
  selector:
    app: hello
  ports:
    - protocol: TCP
      port: 80
```

After applying I can curl this though `kubectl proxy`

```
andraz@amaterasu /tmp/temp
 $ curl http://localhost:8001/api/v1/namespaces/default/services/http:hello:/proxy/
Server address: 10.1.35.142:80
Server name: hello-767b8bc964-f4xfs
Date: 14/Apr/2021:15:55:14 +0000
URI: /
Request ID: 1218ef33c76f9bdc5616b8934f64898b
andraz@amaterasu /tmp/temp
 $ curl http://localhost:8001/api/v1/namespaces/default/services/http:hello:/proxy/
Server address: 10.1.35.141:80
Server name: hello-767b8bc964-scq6z
Date: 14/Apr/2021:15:55:16 +0000
URI: /
Request ID: a4cd09501f7fb10f457ef40ecc4f4d27
```

Note different server ips for subsequent requests. Multiple pods are in fact
serving traffic.

# Setting up ingress

I did enable `traefik` ingress controller so this should be pretty
straightforward.

First of all I set up a domain record `*.microk8s.edofic.com` pointing directly
to the IP of my single node.

Let's try it out

```
 $ curl http://demo.microk8s.edofic.com
curl: (7) Failed to connect to demo.microk8s.edofic.com port 80: Connection refused
```

huh, this should return 404. Let's see what traefik is actually listening to

```
andraz@amaterasu /tmp/temp
 $ k get pod -A | grep traefik
traefik              traefik-ingress-controller-l57tb             1/1     Running   0          8h
andraz@amaterasu /tmp/temp
 $ k get pod -n traefik traefik-ingress-controller-l57tb -o json | jq '.spec.containers[] | .ports'
[
  {
    "containerPort": 8080,
    "hostPort": 8080,
    "name": "http",
    "protocol": "TCP"
  }
]
```

Apparently it's 8080... I can live with this for now.

```
 $ curl http://demo.microk8s.edofic.com:8080
404 page not found
```

Lift off! Now let's wire up some ingress

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: hello
spec:
  rules:
  - host: "demo.microk8s.edofic.com"
    http:
      paths:
      - pathType: Prefix
        path: "/"
        backend:
          service:
            name: hello
            port:
              number: 80
```

And apply. Should work now.

```

 $ curl http://demo.microk8s.edofic.com:8080
Server address: 10.1.35.141:80
Server name: hello-767b8bc964-scq6z
Date: 14/Apr/2021:17:05:10 +0000
URI: /
Request ID: bebd6693485e5581305da8c3e68060de
```

And indeed it does.

# Storage

I'll skip over storage as I can see from the description in status that it'sjj
just a host path provisioner - pretty much the same as K3s. If you want fancier
storage you need to bring your own (plenty of options).

Storage classes confirm there is nothing fancy going on

```
 $ k get storageclass
NAME                          PROVISIONER            RECLAIMPOLICY   VOLUMEBINDINGMODE   ALLOWVOLUMEEXPANSION   AGE
microk8s-hostpath (default)   microk8s.io/hostpath   Delete          Immediate           false                  8h
an
```

# Scaling / high availability

So let's try out the magical HA setup! According to [the
docs](https://microk8s.io/docs/high-availability) it should be as simple as
creating a token on the existing node and then issuing a single command on new
nodes to join. But first I need to create some new nodes.

I created two more VMs just like the first one and ran through the same steps on
each one:

```
apt update
apt install -y snapd
snap install microk8s --classic
```

Then on my inital node

```
root@microk8s-1:~# microk8s add-node
From the node you wish to join to this cluster, run the following:
microk8s join 162.55.52.75:25000/0b9399975ac32273a17f210ae7e25802

If the node you are adding is not reachable through the default interface you can use one of the following:
 microk8s join 162.55.52.75:25000/0b9399975ac32273a17f210ae7e25802
 microk8s join 10.0.0.2:25000/0b9399975ac32273a17f210ae7e25802
```

And follow the instructions (using the private IP). First on the second node.

```
root@microk8s-2:~# microk8s join 10.0.0.2:25000/0b9399975ac32273a17f210ae7e25802
Contacting cluster at 10.0.0.2
Waiting for this node to finish joining the cluster. ..
root@microk8s-2:~# microk8s status
microk8s is running
high-availability: no
  datastore master nodes: 10.0.0.2:19001
  datastore standby nodes: none
...
root@microk8s-2:~# microk8s kubectl get nodes
NAME         STATUS   ROLES    AGE   VERSION
microk8s-1   Ready    <none>   8h    v1.20.5-34+40f5951bd9888a
microk8s-2   Ready    <none>   69s   v1.20.5-34+40f5951bd9888a
```

Apparently I have a two node cluster now. Let's join in the third node -
microk8s should automatically switch to HA mode.

```
root@microk8s-3:~# microk8s join 10.0.0.2:25000/0b9399975ac32273a17f210ae7e25802
Contacting cluster at 10.0.0.2
Failed to join cluster. Error code 500. Invalid token
```

Huh, looks like I need a new token. Fine.

```
root@microk8s-3:~# microk8s join 10.0.0.2:25000/cab6b9f7f538bbb5c7e1a3a3f6239c50
Contacting cluster at 10.0.0.2
Waiting for this node to finish joining the cluster. ..
```

After this is done I can check (on any node, outputs agree now)

```
root@microk8s-1:~# microk8s status
microk8s is running
high-availability: yes
  datastore master nodes: 10.0.0.2:19001 10.0.0.3:19001 10.0.0.4:19001
  datastore standby nodes: none
...
```

Yay: `high-availability: yes`. It worked. Looking from my laptop I now also see
three nodes.

```
 $ k get nodes
NAME         STATUS   ROLES    AGE     VERSION
microk8s-2   Ready    <none>   5m10s   v1.20.5-34+40f5951bd9888a
microk8s-1   Ready    <none>   9h      v1.20.5-34+40f5951bd9888a
microk8s-3   Ready    <none>   33s     v1.20.5-34+40f5951bd9888a
```

If I check pod placement I still see that everything is still runing on node 1

```
andraz@amaterasu /tmp/temp
 $ k get pod -o wide
NAME                     READY   STATUS    RESTARTS   AGE   IP            NODE         NOMINATED NODE   READINESS GATES
hello-767b8bc964-scq6z   1/1     Running   0          92m   10.1.35.141   microk8s-1   <none>           <none>
hello-767b8bc964-f4xfs   1/1     Running   0          92m   10.1.35.142   microk8s-1   <none>           <none>
hello-767b8bc964-sf4ds   1/1     Running   0          92m   10.1.35.143   microk8s-1   <none>           <none>
```

But if I kill any pod it should quickly be rescheduled around the cluster

```
andraz@amaterasu /tmp/temp
 $ k delete pod hello-767b8bc964-sf4ds
pod "hello-767b8bc964-sf4ds" deleted
andraz@amaterasu /tmp/temp
 $ k get pod -o wide
NAME                     READY   STATUS    RESTARTS   AGE   IP             NODE         NOMINATED NODE   READINESS GATES
hello-767b8bc964-scq6z   1/1     Running   0          92m   10.1.35.141    microk8s-1   <none>           <none>
hello-767b8bc964-f4xfs   1/1     Running   0          92m   10.1.35.142    microk8s-1   <none>           <none>
hello-767b8bc964-p9xw8   1/1     Running   0          9s    10.1.100.129   microk8s-2   <none>           <none>
andraz@amaterasu /tmp/temp
```

Great, everything behaving as expected.

# Load balancing

I want to try to kill my initial node to see HA in action but this will bork my
ingress as the IP is hard coded. Luckily Hetzner also provides managed load
balancers so I'll set one up. I can use labels to automatically pick up servers
to use as targets, use health checks to see which ones take traffic and even
expose port 80 and target 8080 internally via internal IPs. Yes I can target any
server as traefik ingress controller is provisioned on all of them and will
route traffic internally to wherever target pods are scheduled.

So all I really need to do know is to update my domain record and my cluster is
none the wiser.

Mind the lack of port 8080

```
 $ curl http://demo.microk8s.edofic.com
Server address: 10.1.35.142:80
Server name: hello-767b8bc964-f4xfs
Date: 14/Apr/2021:17:32:52 +0000
URI: /
Request ID: 237198f55f60c5cc4629750d1f4e3908
```

# Deleting the initial node

With my setup now truly HA let's try to kill the initial node. I simulated this
by doing a hard power off of `microk8s-1`.

Load balancer needed a few moment to detecd the node is down and then it stoped
routing traffic to it.

Microk8s needed a bit more time. Initially `kutectl get node` still showed no.1
as ready. But after a few minutes it turned `NotReady`

```
root@microk8s-2:~# microk8s kubectl get nodes
NAME         STATUS     ROLES    AGE   VERSION
microk8s-1   NotReady   <none>   9h    v1.20.5-34+40f5951bd9888a
microk8s-2   Ready      <none>   36m   v1.20.5-34+40f5951bd9888a
microk8s-3   Ready      <none>   32m   v1.20.5-34+40f5951bd9888a
```

However this means that my pods are still scheduled to node 1 so I get reduced
availability of my service. However the service objects correctly routes so
requests still go through. Not great though.

Reading the docs I now need manual intervention to tell kubernetes that this
node is not coming back.

```
root@microk8s-2:~# microk8s remove-node microk8s-1 --force
root@microk8s-2:~# microk8s kubectl get nodes
NAME         STATUS   ROLES    AGE   VERSION
microk8s-2   Ready    <none>   39m   v1.20.5-34+40f5951bd9888a
microk8s-3   Ready    <none>   35m   v1.20.5-34+40f5951bd9888a
```

Better. How about my pods?

```
root@microk8s-2:~# microk8s kubectl get pod -o wide
NAME                     READY   STATUS    RESTARTS   AGE     IP             NODE         NOMINATED NODE   READINESS GATES
hello-767b8bc964-p9xw8   1/1     Running   0          34m     10.1.100.129   microk8s-2   <none>           <none>
hello-767b8bc964-5787w   1/1     Running   0          3m38s   10.1.151.68    microk8s-3   <none>           <none>
hello-767b8bc964-bk5gd   1/1     Running   0          3m38s   10.1.151.65    microk8s-3   <none>           <none>
```

All rescheduled. So microk8s HA is not a silver bullet - it does not know about
underlying comput infrastructure. But it's really low-ops - just as the
description claims :D

# Closing thoughts

It's an interesting piece of technology. I definately see a use case where you
want to easily setup a cluster on a small number of manually managed machines.
May come use in the future.
