---
title: Attaching your local machine to a Kubernetes cluster
date: 2020-09-27
draft: true
---

provision a VM (digital ocean)
ssh to it: `curl -sfL https://get.k3s.io | sh -`
`k3s kubectl get node` to confirm

copy `/etc/rancher/k3s/k3s.yaml` to local machine
update server ip to use VM's external ip
use it as kubeconfig: either save to `~/.kube/config` or set env var `KUBECONFIG`

locally: kubectl get node

I'm using Ubuntu 20.04 which has a firewall (ufw) by default: ufw allow 6443

curl -sfL https://get.k3s.io | K3S_URL=https://myserver:6443 K3S_TOKEN=mynodetoken sh -
cat /var/lib/rancher/k3s/server/node-token

cd
curl -L -o k3s https://github.com/rancher/k3s/releases/download/v1.18.9%2Bk3s1/k3s
chmod +x ./k3s
./k3s agent --token=K1095fb48d627b9fcbf23f2e59d99f91c8298d7a43701e98415f4bdbcd74aa9fc27::server:27e2ebf5af5a8be04e7fbdd00e9554da --server=https://207.154.197.91:6443 --docker

kubectl create deployment hello-node --image=k8s.gcr.io/echoserver:1.4
kubectl get pod -o wide


kubectl port-forward  pod/hello-node-7bf657c596-dgtmk 8080
http://localhost:8080/

docker ps


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

kubectl label nodes amaterasu name=amaterasu
kubectl create -f ./local.yaml
kubectl get pod -o wide

docker exec -it k8s_test-container_local-pod_default_1e9ecfbb-6b03-493e-94c8-82b7dbb8bc13_0 bash
cd /host
date > now

cat /tmp/shared/now


kubectl taint nodes amaterasu personal=amaterasu:NoExecute
kubectl get pod -o wide

`hello-node*` is restarted on remote server
local-pod is gone!


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


kubectl create -f ./local.yaml


kubectl delete node amaterasu
kill k3s
i also teared down the remote cluster at this point to stop incurring charges



cleanup
rm -rf /etc/rancher /run/k3s /run/flannel /var/lib/rancher /var/lib/kubelet




dockerized variant

docker run --rm -it --name=k3s-server --privileged rancher/k3s:v1.18.9-k3s1 server
TOKEN=$(docker exec -it k3s-server /bin/sh -c "cat /var/lib/rancher/k3s/server/node-token")
SERVER_IP=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' k3s-server)
docker run --rm -it --name=k3s-agent --privileged rancher/k3s:v1.18.9-k3s1 agent --token "$TOKEN" --server "https://${SERVER_IP}:6443"


docker run --rm -it --name=k3s-agent --privileged -v /var/run/docker.sock:/var/run/docker.sock rancher/k3s:v1.18.9-k3s1 agent --token "$TOKEN" --server "https://${SERVER_IP}:6443" --node-name=amaterasu --docker


