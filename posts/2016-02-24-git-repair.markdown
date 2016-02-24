---
title: Repairing a corrupt Git repo using a clone
---

Quite recently I managed to make myself a corrupt git repository due to a file system failure.
See, git stores everything in content addressable blobs - the file name of something is it's hash. Which lends itself nicely to checking repository integrity - it keeps out malicious attackers as well as my file system problems.

I already hear you saying: Why not just make a new clone, git is distributed anyway?
Well, I wasn't diligent enough to push everything. I had local commits that were quite important, so I spent some time fixing it.

## fsck

Git has a command to manually check integrity of the repository: `git fsck`.
Running it lists all the errors.

```bash
 $ git fsck
error: garbage at end of loose object '3ce5b3af5d47179ff31a665ff267e1c7b6e4d8aa'
fatal: loose object 3ce5b3af5d47179ff31a665ff267e1c7b6e4d8aa (stored in .git/objects/3c/e5b3af5d47179ff31a665ff267e1c7b6e4d8aa) is corrupt
```

Luckily in my case the list was quite short so I went ahead and deleted all the objects that were listed as corrupted. So now my objects are fine, but I'm missing some. Luckily (again) corrupted objects did not contain any data pertaining to unpushed commits so I thought I can use a close to restore them.

## unpack

So I lied a bit, git doesn't store every blob in a separate file, that would become huge pretty quickly. Instead it uses **packfiles**. It packs several blobs into one file and does delta compression to reduce disk usage. So I cannot just copy over blobs from a clone.

Fortunately git has commands for dealing with packfiles as well. The one of interest is `git unpack-file` which takes a packfile, extracts all the blobs and dumps them into the repo. Potentially producing loose objects, but let's not care about that for a second.

So I made a bare clone from github

```bash
git clone --bare git@github.com/edofic/blog
```

And just unpacked everything

```bash
cd ~/actual-blog-repo
git unpack-file < /tmp/blog.git/objects/pack/*.pack
```

And it worked! `git fsck` did not complain anymore. Well at least not about garbage and corruption - just loose objects.

But that is easy to clean up: just prune them

```bash
git prune --expire now
```

And do a GC to re-compress.

```bash
git gc
```

Any my repo integrity is back!
