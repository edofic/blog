---
title: Static linking with GHC on ArchLinux
---

There are many reasons why to prefer dynamic linking to static but I'll not go through them. Sometimes you just want static linking, period. In my case it was to show that [Go](http://golang.org/)'s static executables without dependencies are not something special and other languages can do it as good as well - Haskell included. 
My compiler of choice is GHC and I'm running ArchLinux. More on why this is important later.

Firstly I'll need I file to test - `Main.hs`
```haskell 
module Main where

main = putStrLn "Hello static world"
```
Then a quick search supplied me with the command line to compile this statically

```bash
ghc -O -static -threaded -optl-static Main.hs
```
`-O` is for optimizations, `-static` instructs GHC to do static compilation, `-threaded` includes *pthread* and `-optl-static` pushes `-static` flag to `ld`.

But it didn't work. Instead I got a bunch of errors from `ld` telling me I'm missing `librt` and `libgmp`. Running `locate librt` turned up results as well as `locate libgmp`. I was flabbergasted. 

Then I tried running the same thing on Ubuntu 12.04 LTS and it worked. The resulting binary also run on my Arch without problems. Now I was just sad. I tried searching online for my problem but apparently my google-fu is insufficient. I also tried setting *gold* as preferred linker but to no avail. 

### Few weeks later

Today I was playing around with C and when I got something working I decided to link it statically so I can send it off to a colleague who doesn't have all these obscure libraries installed. And I hit into a similar problem. Now it couldn't find `libgc` - a library I was using that worked like a charm when using dynamic linking. 

Apparently the problem didn't lie in GHC but in my linker. Time to put on my Sherlock Holmes hat and investigate. 

Turns out I'm a bloody ignorant idiot. There are dynamic libraries(with `.so` extension) and there are static libraries(with `.a` extension). I remember knowing this once. And I had all dynamic libraries installed but not static. This was the root of my problems with GHC and now with GCC. 

More researching turned up that Arch shies away from providing static libraries in order to encourage dynamic linking. If you want static objects you'll have to build them from source. 

### Solution

I build `libgmp` and then also `libc` in order to get `librt` out. It wasn't that long. But for your convenience here are the resulting files if you want them [libgmp.a](../files/libgmp.a) [librt.a](../files/librt.a)

I dumped those into `/usr/local/lib` because I didn't want to pollute my global libraries. Now I just need to convince GHC to use them. Easy. Just set `LD_LIBRARY_PATH` to that path.

```bash
LD_LIBRARY_PATH=/usr/local/lib ghc -O -static -threaded -optl-static Main.hs
```

And it works. Now I'm happy.
