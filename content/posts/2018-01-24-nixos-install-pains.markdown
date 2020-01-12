---
title: NixOS install & labor pains
date: 2018-01-24
---

# The Back Story

My first foray into world of Linux happened with Red Hat Linux 6 (original, not
the enterprise one). It was magical but it didn't stick. See, I was a kid who
liked video games. So it was back to Windows until I finally got serious about
programming and discovered I the development environment on Linux. So I
formatted my hard drive and installed Ubuntu. I quickly learned absorbed
information and slowly got bored. When Unity came out it was time to switch. I
dabbled a bit with Mint but ultimately landed on ArchLinux. I really like the
*RTFM or GTFO* philosophy of Arch because it means you learn stuff and you
learn it quickly. Because you just have to. Oh, and there is great
documentation. And everything was great. Until I borked my system doing an
upgrade or some other system package operation. A few times. Btrfs to the
rescue!  But there must be a better way. (I was doing manual snapshots of the
root filesystem before any potentially destructive operation.

Then I heard my colleagues talking about this amazing package manager called
Nix and the derived distribution - [NixOS](https://nixos.org).  It's tag lines
are (respectively:

> The Purely Functional Package Manager

and

> The Purely Functional Linux Distribution

# The Install

Being the Haskell geek I am, I needed to give it a spin. I tried the Virtual
Box image but it was weird. It has KDE and I couldn't really figure out how to
alter it. So I decided to just try to do a fresh install in Virtual Box.

Turns out the manual pretty much has you covered. It's not Arch Wiki level, but
good enough to set up a minimal install without any major issues.

At this point I was so fascinated by the core concepts that I just bit the
bullet - I decided to install it onto my physical machine. What could possibly
go wrong?

So I just created a quick subvolume (Btrfs again!) and rebooted. About half an
hour later I had a working NixOS install alongside my Arch. It was really
minimal but I was happy and I built on it quickly. A few evenings of playing
around later I just switched. I set up pretty much the same environment I had
before so my work productivity didn't suffer a bit. If anything it was better -
now I no longer feared system upgrades because rollback had my back.

# The Crash

And then it happened. My hard drive failed and there were issues with the
external backup as well. I managed to recover data but not the system. No
biggie - NixOS is trivial to rebuild from the config file. Or so I thought. I
had to locally rebuild the world, starting from `glibc`. (I think it was my
error for using an old boot image.) That took pretty much the whole day.

So I devised plan for re-installs which I tested out and I'm now documenting
here for future usage:

1. download fresh `nixos-unstable` image
1. boot it and install minimal system with network manager
1. boot into the new system
1. update channels
1. fetch personal config
1. rebuild

With this approach I'm up and running an exact replica of the system in about
half an hour. Yup, my hard drive crashed again. This time my backup was fine but
I decided to actually go this route to test out the procedure. Works like a
charm. And I also did it when a bought a new computer. It's kind of creepy when
after half an hour you are greeted by a familiar desktop on an unfamiliar
hardware. Creepy but very cool.
