---
title: Rename rebase with git
date: 2020-04-20
---

Want to perform a `git rebase` in which you rename a file without pesky
conflicts?

tl;dr

```bash
git filter-branch -f --tree-filter "git mv ORIGINAL_FILE_NAME NEW_FILE_NAME || true" -- $(git merge-base origin/master HEAD)..HEAD
git rebase origin/master
```

Ever wanted to rename a file you created/modified in a pull request but also
wanted to keep the pristine history you worked hard for? Well tough luck you
can either push a new commit to rename it (and keep old commits with the old
file name) or resolve a bunch of pointless conflicts. There must be a better
way...

Turns out you can have your commit history and eat it too! You can tell the git
to go through your history and do the rename on each commit _and then_ rebase.

Let's walk through the steps. First we need to figure out which commits are
affected. To start of we need to find what git calls a "merge base" - the
nearest common ancestor with the target branch. This is the point from which
commits will be considered by e.g. GitHub when you open a pull request.

```bash
git merge-base origin/master HEAD
```

Assuming the target repo is `origin` and the target branch is `master`. `hEAD`
just points to whatever is currently checked out this will give use the merge
base commit hash.

We get the full commit range with the `..` operator. We can put the merge-base
command into a subshell to build up our one-liner

```bash
$(git merge-base origin/master HEAD)..HEAD
```

_NONE_ this is not valid bash but a git construct to be used in further parameters.

And here comes the main star - git command to rewrite history: `git filter-branch`. It has many options and modes - read the docs for more info, here we'll just focus on our use case.

```bash
git filter-branch -f --tree-filter "git mv ORIGINAL_FILE_NAME NEW_FILE_NAME || true" -- $(git merge-base origin/master HEAD)..HEAD
```

Let's walk through each part.

- `filter-branch` is a command that rewrites a given set of commits using given commands
- `-f` is a _force_ so we don't need extra confirmations - feel free to skip this one
- `--tree-filter` - this tells git we want to rewrite _file trees_ meaning it will
  checkout each tree, run our command and commit it back (using something like and amend)
- `"git mv ORIGINAL_FILE_NAME NEW_FILE_NAME || true" ` this is our rename command.
   We use `git mv` instead of plain `mv` so the result gets added to index automatically.
   We add `|| true` so the exit code is 0 even if the file is missing - great for working
   with larger histories.
- `-- $(git merge-base origin/master HEAD)..HEAD` and here comes our commit range expression.
   Mind the space after `--` - this is `filter-branch` syntax to separate commits from
   everything else.

Then wait a bit (yes `filter-branch` may be quite slow).....And we're done.
Well, almost. Still need to do the actual rebase - thus far we've only done the
rename.

```bash
git rebase origin/master
```

My original use case was rebasing migrations that use sequence numbers as file names - you run into conflicts whenever somebody else merges any migration. So I've created an alias and put it into my `~/.gitconfig`


```bash
[alias]
  rebase-migration = "!f() { git filter-branch -f --tree-filter \"git mv migrations/$1.sql migrations/$2.sql || true\" -- $(git merge-base $3 HEAD)..HEAD; git rebase $3; }; f"
```

Now I can simply run

```bash
git rebase-migration 76 77 origin/develop
```
