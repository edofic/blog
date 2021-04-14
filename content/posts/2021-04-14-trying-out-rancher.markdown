---
title: Trying out Rancher
date: 2021-04-14
tags: ["trying-out"]
---

This time it's a bit different: this post has very little code and a lot of
screenshots as I was primarily working with a GUI.

Reading about (and working with) [K3s](https://k3s.io/) is almost impossible to
not hear about [Rancher](https://rancher.com/) - the company behind K3s. But
what had me slightly confused is their primary product: [Rancher (the
software)](https://rancher.com/products/rancher/).

The web page says

> Managed Kubernetes Cluster Operations

but that doesn't really tell me much. Neither do their presentations. So what
better way to get a feel for it than taking it for a spin :) (there is some
enterprise offerings but the base product is free).

# Getting started

After clicking the fetching [Getting Started](https://rancher.com/quick-start/)
button I'm presented with these rather short instructions:

![getting started](/images/rancher/01.png)

Ok, prerequisites first. I spun up a new VM on [Digital
Ocean](https://www.digitalocean.com/) and selected Ubuntu 20.04 (which is on the
list of supported distros and is not ancient), 8gb of ram and 4 vCPUs -
completely guessing here - no idea what I actually need I picked a "smallish VM
that should still handle any reasonable application".

Second step is getting [docker](https://www.docker.com/) on my machine, this is
where picking Ubuntu comes in handy:

```
snap install docker
```

Just a few minutes later I can actually proceed onto actually running rancher
via docker:

```
docker run --privileged -d --restart=unless-stopped -p 80:80 -p 443:443 rancher/rancher
```

And just like that we're up and listening on port 80!

```
root@ubuntu-s-4vcpu-8gb-amd-fra1-01:~# docker ps
CONTAINER ID        IMAGE               COMMAND             CREATED             STATUS              PORTS                                      NAMES
89ffd60c0892        rancher/rancher     "entrypoint.sh"     16 seconds ago      Up 15 seconds       0.0.0.0:80->80/tcp, 0.0.0.0:443->443/tcp   sweet_cannon
r
```

# First impressions

I navigate to port 80 of my new VM from my laptop and I'm greeted with a first
time setup. Nice :)

![initial setup](/images/rancher/02.png)

Let's go with multiple clusters, to see what this is all about.

![global dashboard](/images/rancher/03.png)

Huh. Apparently we got a local K3s cluster as well. I'm guessing "local" means
it's running on the same VM (there is no real other option anyway).

I'm wondering if you can also install rancher on an existing cluster. Quick
search suggests the answer [is
yes](https://rancher.com/docs/rancher/v2.x/en/installation/install-rancher-on-k8s/).

What happens if I click on "local"?

![04.png](/images/rancher/04.png)

Vaguely reminiscent of [K8s Web
UI](https://kubernetes.io/docs/tasks/access-application-cluster/web-ui-dashboard/). Looks useful.

# Creating clusters

But I'm not interested in a dashboard for a K3s instance - I want to try setting
up some clusters :D Let's try clicking add cluster.

![05.png](/images/rancher/05.png)

That's quite a list of providers already. Quite intriguing one is "Existing
nodes". Also interesting to me is the ability to connect managed clusters. You
really do get one central management tool for all your clusters. Huh you can
even automatically provision VMs on vSphere. But I'm going with DigitalOcean
today as this is where I'm doing my experiment. A happy accident really -
haven't checked what is supported upfront.

But I did do some digging at this point and apparently you can install extra
(even 3rd party) providers to extend this. E.g. I've found [this tutorial](https://jmrobles.medium.com/how-to-create-a-kubernetes-cluster-with-rancher-on-hetzner-3b2f7f0c037a) for setting up a cluster on [Hetzner](https://www.hetzner.com/) - which is another provider I sometimes use for my projects. Liking this extensibility a lot.

![06.png](/images/rancher/06.png)

First I need to configure a template for my nodes (since this is the first
time and I don't have any?). I'm going with something small as I don't have
plans for any substantial workloads on this cluster.

Note: I was also prompted to configure DO integration at this point by providing
an API token - no magic integration but super straightforward.

![07.png](/images/rancher/07.png)

Time to setup my node pools now. If I'm reading this right I'm configuring how
many VMs of what template to run and what services to run on them. There's even
validation (green checkmarks) that I have a setup that yields a functional
cluster. I won't bother with other options, let's just take defaults out for a
spin.

![08.png](/images/rancher/08.png)

And now we wait....

Meanwhile in DO console I can see droplets coming up

![09.png](/images/rancher/09.png)

After about 10min the cluster reports as "active". Great success!

![11.png](/images/rancher/11.png)

Clicking that I get an overview of my new nodes

![12.png](/images/rancher/12.png)

And some management too

![13.png](/images/rancher/13.png)

Doesn't seem much...but stop to think about it for a moment. This abstracts over
the underlying compute provider which is great if you're running multiple
clusters on a mix (e.g. AWS + on prem).

I can also pick "namespaces" from the menu and get an overview of k8s namespaces
on my cluster.

![14.png](/images/rancher/14.png)

And yes - there is still the dashboard.

![15.png](/images/rancher/15.png)

One interesting difference pops out: we're not running K3s anymore but
[RKE](https://rancher.com/products/rke/) - Rancher Kubernetes Engine.

> RKE is a CNCF-certified Kubernetes distribution that runs entirely within
> Docker containers. It solves the common frustration of installation complexity
> with Kubernetes by removing most host dependencies and presenting a stable
> path for deployment, upgrades, and rollbacks.

The key point here is that we got what is basically a managed cluster (k8s as
a service) using only a compute provider (yes I know DO also offers managed k8s
but that a different beast altogether).

# Features

With a few more clicks I got a metrics stack installed in the cluster and now
have more detailed stats.

![16.png](/images/rancher/16.png)

And a configured grafana instance too!

![17.png](/images/rancher/17.png)

I can get a kubeconfig file directly from the UI and use it to connect to my new
cluster like any other cluster

```
andraz@amaterasu /tmp
 $ export KUBECONFIG=$(pwd)/kubeconfig
andraz@amaterasu ~
 $ kubectl get nodes
NAME   STATUS   ROLES                      AGE   VERSION
sb2    Ready    controlplane,etcd,worker   10m   v1.20.5
sb3    Ready    controlplane,etcd,worker   10m   v1.20.5
sb4    Ready    controlplane,etcd,worker   10m   v1.20.5
```

But if I cannot be bothered and just want to poke around a bit there is even a
convenient web-based shell

![18.png](/images/rancher/18.png)

# Apps

Moving on the last element I see in the UI - Apps

![19.png](/images/rancher/19.png)

Looks like some form of prepackaged software.

![20.png](/images/rancher/20.png)

Digging deeper confirms this is in fact a UI for [helm](https://helm.sh/)
charts.

![21.png](/images/rancher/21.png)

Let's configure a Wordpress instance then. This should be a nice test if our
cluster is in fact fully functional and provides all we could possibly want

![22.png](/images/rancher/22.png)

Takes a moment to deploy the chart. I can now zoom in and see the actual
deployment in action

![23.png](/images/rancher/23.png)

Few more minutes and we're green.

![24.png](/images/rancher/24.png)

A quick `/etc/hosts` edit later I'm accessing my WP instance. I'm pointing my
browser at one of the nodes where I'm hitting an ingress controller that routes
to an internal service based on the virtual host (hence `/etc/hosts` edit) that
routes to the actual pod running somewhere in the cluster.

![25.png](/images/rancher/25.png)

# Removing a cluster

That was enough fun. Let's tear the cluster down to see if it cleans up after
itself.

![26.png](/images/rancher/26.png)

Entered a "Removing" state

![27.png](/images/rancher/27.png)

and few moments later is gone. So are all the traces from my DO dashboard.

# Closing thoughts

I think I get "managed kuberenetes cluster operations" now. Provisioning,
monitoring, day 2 ops, running some common workloads. Seems like the primary use
case is managing multiple clusters but it's also an interesting alternative for
providing "managed" clusters in places where you don't have this option.
