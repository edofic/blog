---
title: Trying out Okteto Cloud
date: 2021-02-27
tags: ["trying-out"]
---

This is a part two of [Trying out Okteto](/posts/2021-02-22-trying-out-okteto-cli/).
I'll be trying out [okteto.com](https://okteto.com/) - "The Kubernetes development platform".

At a first glance it looks like a hosted kubernetes solution + the CLI tool I
already tested + some bells and whistles. Let's find out what those are.

# First impressions

Login with github. Boom I'm in, no extra sign up required. Big thumbs up.

Looks like I got a managed namespace with 8G ram, 4cpus and 5G storage. Quite
generous for a free plan.

There is a convenient link to [Getting started](https://okteto.com/docs/getting-started/index.html).

# Getting started

Apparently those bells and whistles include build and deploy
pipeline....reminds me of Heroku :)

I'll deviate from the guide and will try to deploy my buffalo test app because
I actually want to dive deep and understand what's happening.

Let's have a look at [okteto-pipeline.yml from the example](https://github.com/okteto/movies/blob/master/okteto-pipeline.yml)

```yaml
icon: https://apps.okteto.com/movies/icon.png
deploy:
  - okteto build -t okteto.dev/api:${OKTETO_GIT_COMMIT} api
  - okteto build -t okteto.dev/frontend:${OKTETO_GIT_COMMIT} frontend
  - helm upgrade --install movies chart --set tag=${OKTETO_GIT_COMMIT}
devs:
  - api/okteto.yml
  - frontend/okteto.yml
```

Quite a bit going on here. Looking at `okteto build -h` localy it seems this
will build two Dockerfiles, push the images to okteto servers and then  do a
helm upgrade. This implies we need a helm chart for out app. Oh look, there is a
[chart for the example app](https://github.com/okteto/movies/tree/master/chart)

Let's try to mimic this and see how for I can get.

# The helm digression

I'll admit I was a bit intimidated at the prospect of packaging my app up with
help just to get up and running as I've only created a chart once previously and
even that was not from scratch but with some tweaking. But...how hard can it be?

First I copied over my test app since I'll probably be making some changes. I
already have a Dockerfile, okteto.yaml and database deployment configs. I think I
should create a chart next so I can try to deploy the app.

`Chart.yaml`

```yaml
apiVersion: v2
name: tryingout
description: Trying out Buffalo on Kubernetes
type: application
version: 0.1.0
appVersion: 1.0.0
```

Any my `values.yaml` (based of movies example)

```yaml
tag: dev

web:
  replicaCount: 1
  image: okteto.dev/web
```

Copied over notes, helpers (with name replacement). Not it's time to get to
business with the actual templates.  I'll just reuse my deployment for postgres
and just add in the labels (though I should be a stateful set if I was even half
serious).

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: db
  labels:
    component: db
    {{- include "tryingout.labels" . | nindent 4 }}
spec:
  replicas: 1
  selector:
    matchLabels:
      component: db
      {{- include "tryingout.labels" . | nindent 6 }}
  template:
    metadata:
      labels:
        component: db
        {{- include "tryingout.labels" . | nindent 8 }}
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
  labels:
    component: db
    {{- include "tryingout.labels" . | nindent 4 }}
spec:
  selector:
    component: db
    {{- include "tryingout.labels" . | nindent 4 }}
  ports:
    - protocol: TCP
      port: 5432
```

Any now for the meat and the potatoes - my web deployment. Should be pretty much the same as postgres, I just need to
grab the image from the values.

_A word of warning:_ there are errors here - I left them in intentionally so I can
document my full process. See the [git
repo](https://github.com/edofic/trying-out/tree/main/okteto/v2/chart) for latest
versions (which you can actually apply).

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: web
  labels:
    component: web
    {{- include "tryingout.labels" . | nindent 4 }}
spec:
  replicas: 1
  selector:
    matchLabels:
      component: web
      {{- include "tryingout.labels" . | nindent 6 }}
  template:
    metadata:
      labels:
        component: web
        {{- include "tryingout.labels" . | nindent 8 }}
    spec:
      containers:
      - name: web
        image: {{ .Values.api.image }}:{{ .Values.tag }}
        ports:
        - containerPort: 3000
---
apiVersion: v1
kind: Service
metadata:
  name: web
  labels:
    component: web
    {{- include "tryingout.labels" . | nindent 4 }}
spec:
  selector:
    component: web
    {{- include "tryingout.labels" . | nindent 4 }}
  ports:
    - protocol: TCP
      port: 3000
```

An ingress may come in handy as well.

```yaml
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: {{ include "tryingout.fullname" . }}
  labels:
    {{- include "tryingout.labels" . | nindent 4 }}
  annotations:
    dev.okteto.com/generate-host: "true"
spec:
  rules:
    - http:
        paths:
          - path: /
            backend:
              serviceName: web
              servicePort: 3000
```

And now tie it all together with a pipeline

```yaml
deploy:
  - okteto build -t okteto.dev/web:${OKTETO_GIT_COMMIT} web
  - helm upgrade --install tryingout chart --set tag=${OKTETO_GIT_COMMIT}
devs:
  - web/okteto.yml
```


Now how do I run this? :)


Let's read the docs again....

```
 $ okteto namespace
Authentication required. Do you want to log into Okteto? [y/n]: y
What is the URL of your Okteto instance? [https://cloud.okteto.com]:
Authentication will continue in your default browser
You can also open a browser and navigate to the following address:
https://cloud.okteto.com/auth/authorization-code?...
 ✓  Logged in as edofic
 ✓  Updated context 'cloud_okteto_com' in '/home/andraz/.kube/config'
```

Got a browser window to authenticate and I was back in the shell, logged in. Smooth.

Now how do I deploy from the terminal? There are [some docs](https://okteto.com/docs/cloud/deploy-from-terminal).

```
 $ okteto pipeline deploy
  x  failed to analyze git repo: repository does not exist
```

I think this means it too want to be in the repo root. Time to switch gears,
I'll try running steps manually. So we have a build and a helm chart install to
do.

I'll event set up `OKTETO_GIT_COMMIT` env variable so I can directly copy paste
the commands.

```
 $ OKTETO_GIT_COMMIT=$(git rev-parse HEAD)
 $ okteto build -t okteto.dev/web:${OKTETO_GIT_COMMIT} web
 .... a lot of output ....
 ✓  Image 'okteto.dev/web:d9f4f4850cb3e2cbb8fa623e5de68ab8419df304' successfully pushed
 ```

Looks like this actually pushed my "docker context" to the cloud and then ran
something like `docker build` & `docker push` there. Nice.

Let's try deploying my chart.

```
 $ helm upgrade --install tryingout chart --set tag=${OKTETO_GIT_COMMIT}
 Release "tryingout" does not exist. Installing it now.
 Error: template: tryingout/templates/web.yaml:22:25: executing
 "tryingout/templates/web.yaml" at <.Values.api.image>: nil pointer evaluating
 interface {}.image
 ```

 Time to debug my yamls :D

 A quick look at the error message and my config: I did a stupid :) There is
 `api.image` in the template but `web.image` in `values.yaml` - a copy paste
 gone wrong. Let's fix the template.

 relevant part

```yaml
    spec:
      containers:
      - name: web
        image: {{ .Values.web.image }}:{{ .Values.tag }}
        ports:
        - containerPort: 3000
```

And try again

```

 $ helm upgrade --install tryingout chart --set tag=${OKTETO_GIT_COMMIT}
Release "tryingout" does not exist. Installing it now.
NAME: tryingout
LAST DEPLOYED: Mon Feb 15 16:37:14 2021
NAMESPACE: edofic
STATUS: deployed
REVISION: 1
TEST SUITE: None
NOTES:
```
# Some success

And somthing popped up in the UI (no refresh needed!)

![my chart installed](/images/okteto/1_my_chart.webp)

I can see the status, I can destroy/redeploy.... How can I now try to open up
the app? I created an ingress so let's see if I can use it.

```
 $ kubectl get ingress
NAME        HOSTS                               ADDRESS        PORTS     AGE
tryingout   tryingout-edofic.cloud.okteto.net   35.238.195.0   80, 443   13m
```

Oh wow, that actually worked we got a domain provisioned.  Opening up this
hostname gives my my error page - migrations haven't run yet.

![migrations error](/images/okteto/2_error_migrations.webp)

I should probably do a helm post upgrade hook to run migrations...but I'll half
ass it and do it manually now since I'm trying to focus on okteto instead.

Looking at buffalo generated Dockerfile to orient myself I came across this
helpful nugget

```
# Uncomment to run the binary in "production" mode:
# ENV GO_ENV=production
...
# Uncomment to run the migrations before running the binary:
# CMD /bin/app migrate; /bin/app
```

Let's try this to see if out app stands up, then actually uncomment this and
we're off to the races!

```
 $ kubectl exec -it  web-644d8d487f-gp9hq -- sh
/bin # export GO_ENV=production
/bin # /bin/app migrate
WARN[2021-02-15T15:59:59Z] Unless you set SESSION_SECRET env variable, your session storage is not protected!
[POP] 2021/02/15 15:59:59 warn - ignoring file .pop-tmp.md because it does not match the migration file pattern
[POP] 2021/02/15 16:00:00 info - 0.0901 seconds
2021/02/15 16:00:00 Failed to run migrations: Migrator: problem creating schema migrations: couldn't start a new transaction: could not create new transaction: failed to connect to `host=db user=postgres database=trying_out_production`: server error (FATAL: database "trying_out_production" does not exist (SQLSTATE 3D000))
```

Apparently we also need to create the database. Let's do it manually at first to
confirm.

```
 $ kubectl exec -it db-f9f647bd4-l9gf6 -- psql -U postgres
create databasepsql (13.1)
Type "help" for help.

postgres=# create database trying_out_production;
CREATE DATABASE
postgres=#
```

Then try to run migrations again.

```
/bin # /bin/app migrate
WARN[2021-02-15T16:01:43Z] Unless you set SESSION_SECRET env variable, your session storage is not protected!
[POP] 2021/02/15 16:01:43 warn - ignoring file .pop-tmp.md because it does not match the migration file pattern
[POP] 2021/02/15 16:01:44 info - > create_users
[POP] 2021/02/15 16:01:44 info - Successfully applied 1 migrations.
[POP] 2021/02/15 16:01:44 info - 0.2692 seconds
```

Let's enable this for our build by uncommennting lines in the dockerfile.
And do the build & helm upgrade dance again.


```
 $ OKTETO_GIT_COMMIT=$(git rev-parse HEAD)
 $ okteto build -t okteto.dev/web:${OKTETO_GIT_COMMIT} web
 .... a lot of output ....
 ✓  Image 'okteto.dev/web:88607af3ea6dd4e5ee631d77dee462c5cf82f344' successfully pushed
 ```

while this was building I looked up [docs for my postgres
image](https://hub.docker.com/_/postgres) and found out I can automate the
database creation part by setting `POSTGRES_DB` env variable. So I'll update the
chart as well.

From chart/templates/postgres.yaml

```
      containers:
      - name: postgres
        image: postgres:13.1-alpine
        env:
        - name: POSTGRES_PASSWORD
          value: "postgres"
        - name: POSTGRES_DB
          value: "trying_out_production"
        ports:
        - containerPort: 5432
```

I'll actually click destroy helm chart in the UI now to try it out and clean up
to see if everyting stands up from scratch.

```
helm upgrade --install movies chart --set tag=${OKTETO_GIT_COMMIT}
```

And we have liftoff!

![app is up](/images/okteto/3_up_in_cloud.webp)


Now I think my app is properly....err...well enough packaged and I can actually
try out pipelines. I'll kill the helm chart once more so I can try the
automation now.

# Trying pipelines

I'll start off by creating a new orphan branch and then copy over my files
(won't bother rewriting history).

```
 $ git checkout --orphan okteto-pipeline
Switched to a new branch 'okteto-pipeline'
 $ git reset .
 $ rm -rf *
 $ git checkout main -- ./okteto/v2
 $ mv okteto/v2/* ./
 $ rm -rf okteto
 $ git commit -am "Initial pipeline commit"
 ...
 $ git push -u origin HEAD
Enumerating objects: 78, done.
Counting objects: 100% (78/78), done.
Delta compression using up to 8 threads
Compressing objects: 100% (70/70), done.
Writing objects: 100% (78/78), 167.44 KiB | 1.86 MiB/s, done.
Total 78 (delta 2), reused 0 (delta 0), pack-reused 0
remote: Resolving deltas: 100% (2/2), done.
^[[Aremote:
remote: Create a pull request for 'okteto-pipeline' on GitHub by visiting:
remote:      https://github.com/edofic/trying-out/pull/new/okteto-pipeline
remote:
To github.com:edofic/trying-out.git
 * [new branch]      HEAD -> okteto-pipeline
```

Time to set up the pipeline

![deploy ui](/images/okteto/4_empty_namespace.webp)

UI is obvious enough...

![deploy setup](/images/okteto/5_deploy.webp)

Not much to configure...

![deploying](/images/okteto/6_deploying.webp)

Looks like it's running the same process as I was locally just before,
promising.

![deployed](/images/okteto/7_deployed.webp)

Success!
And out app is up again :D

# Development

Now let's try developing it! I did need to update my `web/okteto.yml` with the
new deployment name `web`.

```
 $ okteto up
 ✓  Development container activated
 ✓  Connected to your development container
 ✓  Files synchronized
    Namespace: edofic
    Name:      web
    Forward:   2345 -> 2345
               3000 -> 3000

root@web-67456fb97b-85xcj:/src#
```

Amazingly the UI picked this up too

![in development](/images/okteto/8_in_development.webp)

Let's run a development version

```
root@web-67456fb97b-85xcj:/src# buffalo dev
...downloading the world for the first time
...building the world for the first time
```

Meanwhile I checked the UI and I see I have some storage going on, these would
be my persistent work environments and go cache so I don't always need to
rebuild the world.

```
 $ kubectl get pvc
NAME         STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
okteto-web   Bound    pvc-ead3ab73-9d8f-4e52-8990-bbe6a6d5b897   2Gi        RWO            okteto         4m8s
```

meanwhile the world has been built...

```
INFO[2021-02-15T17:29:21Z] Starting application at http://127.0.0.1:3000
INFO[2021-02-15T17:29:21Z] Starting Simple Background Worker
```

sigh...this won't work as we need to bind to external IP in order to be
available from the cluster. Quick google search points me at [buffalo docs](https://gobuffalo.io/en/docs/getting-started/config-vars) on how to setup the host:

```
 $ ADDR=0.0.0.0 buffalo dev
```

And we're up! Of course I forgot to create my databases again.
Wait I minute, I do have a db I want to use. Let's edit `database.yaml` instead

```
development:
  dialect: postgres
  database: trying_out_production
  user: postgres
  password: postgres
  host: db
  pool: 5
```

This is a bit nonsensical but I'm just playing around. Note: I did the edit
locally and the file sync picked it up. Sweet.

```
 $ ADDR=0.0.0.0 buffalo dev
```


And I finally have my dev server, that's one time setup for you :)

Now let's do a silly edit. I'll change some text in `templates/users/index.plush.html`            `


```html
<h3 class="d-inline-block">My Users</h3>
```

File sync picked this up, synced, then dev server picked up and reloaded templates. By the time I switched from my editor to the browser and refreshed it was there. :o

Now let's save this.

```
 $ okteto down
 ✓  Development container deactivated
 i  Run 'okteto push' to deploy your code changes to the cluster
```

And we reverted back to the deployed image now. Let's try this suggested `okteto push`

```
 $ okteto push
 ... build commences
```

This apparently pushed my context again, build the image remotely and updated my deployment live. Interesting. Let's try the pipeline way as well. I'll do another silly change and push just my commits to see what happens.


```
  ... my edit
 $ git commit -am "New text change"
[okteto-pipeline 99b7f64] New text change
 1 file changed, 1 insertion(+), 1 deletion(-)
 $ git push
 ...
```

Nothing happened for 2min (maybe I'm just too inpatient) so I've hit redeploy in the UI. Did the trick - my changes are live. Much more quickly that with `okteto push`

# Charts

Once more nice thing: you can deploy some helm charts from a catalog straight form the UI.

![charts](/images/okteto/9_deploy_charts.webp)

# Conclusion

In case this wasn't clear from my other remarks and commands I've ran: okteto
gave my a fully fledged managed kubernetes namespace - this way I was able to
run my standard kubectl and helm commands.

All in all the CLI toool impresses me most. But Okteto Cloud is very nice as
well if you want a single opionionated full stack. You can cobble together
something similar for an existing cluster as well but it will require quite some
work and may have rough edges. (e.g. jenkins + argo cd + harbor registry)
