---
title: Dockerless Services (with Nix)
---

So you want an "isolated" MySQL "service" but don't want to (or are not able
to) run Docker/rkt/LXC/whatever?

Here are the tl;dr steps:

    cd /tmp
    mkdir mysql-1
    cd mysql-1
    nix-shell -p mysql
    mysql_install_db --datadir=$(pwd)
    mysqld --datadir=$(pwd) --socket=$(pwd)/socket &
    mysqladmin -uroot -h127.0.0.1 password <REDACTED>
    mysql -uroot -p<REDACTED> -h127.0.0.1
    ...


## Why?

My motivation for this was running a database on a machine with a specially
tuned kernel that did not support containers.


## How it works?

Create a new folder to host all data for your service

    cd /tmp
    mkdir mysql-1
    cd mysql-1

Get the binary somewehere. I use [Nix](https://nixos.org/nix/) due to its pure
stateless nature - it allows me to have different conlicting versions on my
system without any container magic.

    nix-shell -p mysql

Initialize the data structure. Mind the `--datadir`

    mysql_install_db --datadir=$(pwd)

Launch the server. If you want to run multiple instances you also need to
specify a port number with `--port` here. Socket is also placed in the datadir
so the init does not try to pullute in `/run`

    mysqld --datadir=$(pwd) --socket=$(pwd)/socket &

Finish configuration and use. Need to provide the IP here since UNIX socket is
on a non-standard place.

    mysqladmin -uroot -h127.0.0.1 password <REDACTED>
    mysql -uroot -p<REDACTED> -h127.0.0.1

You can now stop with

    mysqladmin -uroot -p<REDACTED> -h127.0.0.1 shutdown


## Going faster

You can mount `tmpfs` for your datadir to make MySQL go really fast. **Note**
your data is not persisted to your harddrive anymore - use with caution.

```bash
mkdir mysql-1
mount -t tmpfs none mysql-1
cd mysql-1
...
```

Creating an in-memory filesystem is [slightly more
involved](https://www.tekrevue.com/tip/how-to-create-a-4gbs-ram-disk-in-mac-os-x/)
for MacOS users but still not too bad.


## Automate all the things

With the assumption that we want to store data under `/tmp` and that `/tmp` is
already `tmpfs` we can fully automate the script above:


    #! /usr/bin/env bash
    set -e
    PORT=${1:-3306}
    cd $(mktemp -d)
    mysql_install_db --datadir=$(pwd)
    mysqld --datadir=$(pwd) --socket=$(pwd)/socket --port=$PORT &
    sleep 3
    mysqladmin -uroot -h127.0.0.1 password root

Or even make it self-contained and pull in the right MySQL as well

    #! /usr/bin/env nix-shell
    #! nix-shell -i bash -p mysql55
    set -e
    PORT=${1:-3306}
    cd $(mktemp -d)
    mysql_install_db --datadir=$(pwd)
    mysqld --datadir=$(pwd) --socket=$(pwd)/socket --port=$PORT &
    sleep 3
    mysqladmin -uroot -h127.0.0.1 password root
