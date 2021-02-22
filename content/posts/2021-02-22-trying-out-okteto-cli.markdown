---
title: Trying out Okteto (CLI)
date: 2021-02-22
tags: ["trying-out"]
---

Today I'll be giving [Okteto]( https://github.com/okteto/okteto) a spin. From the docs

> Okteto: Tool to Develop Applications on Kubernetes

Or my personal pitch:

> Okteto delivers on the "my other computer is the cloud" trope

What does this mean? First of all: there is a distinction between

* Okteto Cloud & Okteto Enterprise: this is a server side solution that provides
  managed k8s namespaces, CD and other bells and whistles - will look into this
  in the future
* okteto - the tool: open source client to work with Okteto Cloud & Okteto
  Enterprise - but the core functionality works with any k8s cluster. This is the
  tool I'm trying out today.
* Okteto - the company (behind all of this)

# How it works

From the docs:

> Okteto allows you to develop inside a container. When you run okteto up your
> Kubernetes deployment is replaced by a development container that contains
> your development tools (e.g. maven and jdk, or npm, python, go compiler,
> debuggers, etc). This development container can be any docker image. The
> development container inherits the same secrets, configmaps, volumes or any
> other configuration value of the original Kubernetes deployment.

> In addition to that, okteto up will:

> Create a bidirectional file synchronization service to keep your changes up to
> date between your local filesystem and your development container.

> Automatic local and remote port forwarding using SSH, so you can access your
> cluster services via localhost or connect a remote debugger.

> Give you an interactive terminal to your development container, so you can
> build, test, and run your application as you would from a local terminal.o

Or my mental model: just how "docker for mac" allows you to develop locally on
you mac against a virtual machine it manages, okteto allows you to do the same
but against any kubernetes cluster.

Yes underlying mechanisms are very different, but it still allows you to run
containers (or deployments to be exact) and does the necessary wiring to make it
seem they are running on your machine. The same happens when you spin up a
docker container on a mac - it's effecively running on a foreign machine (that
just happens to be virtualised on you laptop), you can actually even use remote
docker hosts.

In my particular use case I'm using spare on-prem servers to offload heavy
services needed for development - that's running on kubernetes and is wired into
my local environment with port forwarding. K8s is used as a scheduler here so my
colleagues can share and more efficiently utilise underlying infra. But over
time the setup grew, offloading more and more components, to the point I'm only
running the component I'm actually working on locally. Which brings me to the
logical conclusion - can I run this on the server too? Might as well  make use
of the more performant hardware. I will confess: I'm not using okteto for my
day-to-day work yet, but I'm looking into it.

So this post is not describing my exact first impression but I genuinely haven't
tried doing this from scratch on linux for [an application that is not yet
packaged for k8s](/posts/2021-02-07-trying-out-buffalo/).
# Installation

Let's see the docs: https://okteto.com/docs/getting-started/installation/index.html

Installation method for linux is a "curl to shell", ugh

```sh
 $ curl https://get.okteto.com -sSfL | sh
```

Not a fan and highly sceptical it will actually work on NixOS. (I've been using
the brew version on mac previously). Let's read through the file and do a manual
approximation.


```bash
$ wget -O ~/bin/okteto https://github.com/okteto/okteto/releases/download/1.10.6/okteto-Linux-x86_64
$ chmod +x ~/bin/okteto
$ okteto --version
okteto version 1.10.6
```

Hooray for static builds, no need to compile from source :D

# A cluster

Now we need a cluster to run against. I have an experimentation setup I'll use
but you can also get a free managed namespace from
[okteto.com](https://okteto.com/) to keep in the spirit of the tool :)

Otherwise you can use [k3d](https://k3d.io/) (k3s in docker) to follow along if
you with (or simply enable kubernetes if you're using docker for desktop)

## K3d

```bash
$ wget -O ~/bin/k3d https://github.com/rancher/k3d/releases/download/v4.1.1/k3d-linux-amd64
$ chmod +x ~/bin/k3d
$ k3d version
k3d version v4.1.1
k3s version v1.20.2-k3s1 (default)
```
We now have the wrapper, let's create a local cluster (requires docker to run
some containers)

```bash
k3d cluster create
```

This stood up a cluster and created a context, just need to use it

follow the output of the command


```bash
nix-env -iA nixos.kubectl # I actually need to install kubectl first
kubectl config use-context k3d-k3s-default
kubectl cluster-info
```

And we're in business :)

# Packaging our app

I'll be using the dummy app I've created in [Trying out
Buffalo](/posts/2021-02-07-trying-out-buffalo/). All the code I produced is
available [on github](https://github.com/edofic/trying-out/tree/main/okteto/v1).

First off I'll set-up a dummy database. Since I don't really need persistence of
any sort of production-level quality I'll create a deployment. This is the
easiest way to get started for me (this is *FAR* from production grade setup).



```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: db
  labels:
    app: db
spec:
  replicas: 1
  selector:
    matchLabels:
      app: db
  template:
    metadata:
      labels:
        app: db
    spec:
      containers:
      - name: postgres
        image: postgres:13.1-alpine
        env:
        - name: POSTGRES_PASSWORD
          value: "postgres"
        ports:
        - containerPort: 5432
---
apiVersion: v1
kind: Service
metadata:
  name: db
spec:
  selector:
    app: db
  ports:
    - protocol: TCP
      port: 5432
```

And apply this config to my cluster

```bash
kubectl apply -f ./postgres.yaml
```

Are we up?

```bash
 $ kubectl get pod
 NAME                  READY   STATUS    RESTARTS   AGE
 db-5d5cd58cdf-6fvsd   1/1     Running   0          18s
```


Yes we are. Now let's use this pod for my local development via a port forward -
I'm setting myself up in a situation where I have my dependencies running in the
cluster so I can get some benefit from doing my development in the cluster a
well.


```bash
kubectl port-forward service/db 5432:5432
```
Now I need to create my database and run the dev server.

```
cd ../buffalo
buffalo pop create -a
buffalo dev
```

BOOM. We're using a remote db. (this is now serving on localhost)

# Enter okteto

Now that I have a dependency in the cluster, let's try to simplify my life using
`okteto`. First thing I like to do with new tools is check the build-in help

```bash
 $ okteto -h
Manage development containers

Usage:
  okteto [command]

Available Commands:
  analytics   Enable / Disable analytics
  build       Build (and optionally push) a Docker image
  create      Creates resources
  delete      Deletes resources
  doctor      Generates a zip file with the okteto logs
  down        Deactivates your development container
  exec        Execute a command in your development container
  help        Help about any command
  init        Automatically generates your okteto manifest file
  login       Log into Okteto
  namespace   Downloads k8s credentials for a namespace
  pipeline    Pipeline management commands
  push        Builds, pushes and redeploys source code to the target deployment
  restart     Restarts the deployments listed in the services field of the okteto manifest
  stack       Stack management commands
  status      Status of the synchronization process
  up          Activates your development container
  version     View the version of the okteto binary

Flags:
  -h, --help              help for okteto
  -l, --loglevel string   amount of information outputted (debug, info, warn, error) (default "warn")

Use "okteto [command] --help" for more information about a command.
```

Hum, we can try running `init` to get started - or you're sane and you read
through the readme to find [Super Quick
Start](https://github.com/okteto/okteto#super-quick-start) that tells you the
same: run `init` and then `up`.

```bash
 $ okteto init -h
Automatically generates your okteto manifest file

Usage:
  okteto init [flags]

Flags:
  -c, --context string     context target for generating the okteto manifest
  -f, --file string        path to the manifest file (default "okteto.yml")
  -h, --help               help for init
  -n, --namespace string   namespace target for generating the okteto manifest
  -o, --overwrite          overwrite existing manifest file

Global Flags:
  -l, --loglevel string   amount of information outputted (debug, info, warn, error) (default "warn")
```

Let's try the defaults - will run this in the buffalo application folder.

```
okteto init
```

Somewhat surprisingly, init was interactive, so no snippets here. Some
takeaways:

- auto detected we're using Go, very nice
- no deployment detected (my app is not deployed to cluster yet), let's go with
  defaults

And we got a new file: `okteto.yml`

```yaml
name: buffalo
image: okteto/golang:1
command: bash
securityContext:
  capabilities:
    add:
    - SYS_PTRACE
volumes:
- /go/pkg/
- /root/.cache/go-build/
sync:
- .:/usr/src/app
forward:
- 2345:2345
- 8080:8080
```

Let's try to dissect it:
- name: default
- image...will come back to this as there is more to unpack
- command, I'm guessing this is what you get as your dev shell - bash is a sane
  default
- PTRACE capabilites - this is usually needed to attach debuggers
- volumes: makes sense to create a volume for go packages and cache - maybe this
  maps to k8s pvcs and we get persistence this way (this is what I'd expect)
- sync: so we'll be syncing current dir to `/usr/src/app`, ok
- forwards: 8000 is a sane default for our app while 2345 is the default port
  for [delve](https://github.com/go-delve/delve) - a go debugger, confirming my
  PTRACE suspicion


Now let's dive deep into the image. You can look at the [layers on docker
hub](https://hub.docker.com/layers/okteto/golang/1.14/images/sha256-060f1ef05023e836f4711e1b41fe8141c7b5d048a8fcd9cecc3bed5af6604ce8?context=explore)

```
1 ADD file ... in /
2 CMD ["bash"]
3 /bin/sh -c set -eux; apt-get
4 /bin/sh -c set -ex; if
5 /bin/sh -c apt-get update &&
6 /bin/sh -c apt-get update &&
7 ENV PATH=/usr/local/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
8 ENV GOLANG_VERSION=1.14.15
9 /bin/sh -c set -eux; dpkgArch="$(dpkg
10 ENV GOPATH=/go
11 ENV PATH=/go/bin:/usr/local/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
12 /bin/sh -c mkdir -p "$GOPATH/src"
13 WORKDIR /go
14 WORKDIR /usr/src/app
15 COPY file:276084f05422a205579d0bb9852e23a9b7bb569ebacb1a59dd92a96a192bccc3 in /root/.bashrc
16 /bin/sh -c go get github.com/codegangsta/gin
17 CMD ["bash"]
```
Looks like a mostly straightforward go env + gin (tool for auto reloading).

But I already have an environment I want to package up - the one I set up
locally in the [buffalo post](/posts/2021-02-07-trying-out-buffalo/): buffalo,
go, node, yarn.

A quick search on docker hub reveals: https://hub.docker.com/r/gobuffalo/buffalo

Looking at [the layers](https://hub.docker.com/layers/gobuffalo/buffalo/v0.16.21/images/sha256-608ee13bdf5a71f37165d43e48ff89757faf4257b46f6979cf4fb3ff6d960681?context=explore)
this looks spot on. Let's try running it locally to verify

```bash
docker run --rm -it gobuffalo/buffalo:v0.16.21 sh
# buffalo version
INFO Buffalo version is: v0.16.20
# go version
go version go1.15.6 linux/amd64
# node --version
v12.20.1
# yarn --version
1.22.10
# pwd
/src
```

We're in business. Since the default okteto images doesn't seem to contain
anything okteto-specific I'll naively try to use the buffalo image and cross my
fingers (I'll also change default ports and sync path to better match buffalo)

```yaml
name: buffalo
image: gobuffalo/buffalo:v0.16.21
command: bash
securityContext:
  capabilities:
    add:
    - SYS_PTRACE
volumes:
- /go/pkg/
- /root/.cache/go-build/
sync:
- .:/src
forward:
- 2345:2345
- 3000:3000
```

That's the whole setup - let's see if it works.

```
 $ okteto up
Installing dependencies...
syncthing-linux-amd64-v1.13.1.tar.gz 10.36 MiB / 10.36 MiB [-----------------------------------------------------] 100.00% 928.87 KiB p/s
 ✓  Dependencies successfully installed
 ✓  Client certificates generated
 ✓  Development container activated
 ✓  Connected to your development container
 ✓  Files synchronized
    Namespace: andraz
    Name:      buffalo
    Forward:   2345 -> 2345
               3000 -> 3000

root@buffalo-bbb6d7f56-vlvng:/src#
```

Got an environment, now let's run our app and hope it works :)

```
buffalo dev
...

couldn't start a new transaction: could not create new transaction: failed to connect to `host=127.0.0.1 user=postgres
database=trying_out_development`: dial error (dial tcp 127.0.0.1:5432: connect: connection refused)
```

Welp, that didn't work. Need to update the database config to point to `db` now
instead of `localhost` (host field)

```
development:
  dialect: postgres
  database: trying_out_development
  user: postgres
  password: postgres
  host: db
  pool: 5
```

And we're in business! Okteto synced the changes in the config automatically.
But new problems arise, can you spot them?

```
INFO[2021-02-15T17:29:21Z] Starting application at http://127.0.0.1:3000
INFO[2021-02-15T17:29:21Z] Starting Simple Background Worker
```

Buffalo is listening on `127.0.0.1` so it will be only reachable from inside the
pod - no use for me trying to access it from the browser on my local machine.

we need to bind to external IP. Quick google search points me at [buffalo
docs](https://gobuffalo.io/en/docs/getting-started/config-vars) on how to setup
the host:

```
 $ ADDR=0.0.0.0 buffalo dev
```

et voila - development environment in the cluster. Now if I do a code change
okteto will pick it up, sync to the pod where buffalo server will
recompile/reload and I get seamleass development experience.

But more on that in the next installment where I go deeper and try out Okteto
Cloud as well.

# Volumes

Let's confirm my suspicion about PVCs:

```
 $ kubectl get pvc
NAME         STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
okteto-web   Bound    pvc-ead3ab73-9d8f-4e52-8990-bbe6a6d5b897   2Gi        RWO            okteto         4m8s
```

yup, we get persistent storage for our environment. Love it.

# Shutting down

Interestingly, just exiting the shell does nothing. You can `okteto up` again
and re-attach to a running pod. Great for e.g. switching WI-FIs.  Of course you
can shut down:

```
okteto down
```

this will still leave the PVC so we get to keep our cache (and other state).


# Conclusion

I like this technology very much as it seems it just generalizes existing
concepts with very little surprises. I'll do another installment of "trying out"
with Okteto Cloud and more hands on experience with development.

