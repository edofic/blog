---
title: Haskell Web Development
---

### Intro
Lately I got sucked into Haskell. When I first saw it it looked like it might be a bit impractical for real-life projects but this prejudice faded away slowly. Now I'm at the point where I have an idea for a web application and I'd like to do it in Haskell. After a quick search I found many frameworks and libraries that I might use. 

### Setting Up
First up is the set up of all used frameworks. I created a hello world project for each to get a sense of things. Here are some comments. 

Each project first got a cabal sandbox and `.cabal` file. 

#### Happstack
http://www.happstack.com/docs/crashcourse/index.html
Tutorial sure is short on how to get going. Let's see if I can figure this out. I added `happstack` and `happstack-server` to my dependencies. And in my `Main.hs` I'vre written the imports and one line of code that says return "hello world" for all requests. `cabal run` and voila, hello world page at `localhost:8000`. Things sure are simple. And I got out a single binary that I can run - or deploy. But this means no autorecompile and reload I grew to like from Playframework(or any dynamic framework really).

####
http://snapframework.com/docs/quickstart
Snap says in their tutorial to use their tool to generate the project skeleton. Ok I install the `snap` package(into sandbox). While it's installing I might stress that I prefer convention over configuration and this implies minimal initial set up, preferably something I can do by hand in a minute - like Happstack. I guess this might divide up libraries from frameworks.
OK. Snap installed. `snap init barebones` now created a skeleton. And there isn't much. It created `snap.cabal` that has some useful dependencies already in place and some GHC options. There is also main source that has an example how to do routing. Great. But I can see that a minimal example would still be just a few lines. Something I could do by hand. Great.

How about that automatic recompilation? Tutorial says I have to initialize a full fledged project instead of bare bones one - `snap init`. Whoa this reated a bunch of stuff. There are some templates in a special templating languagem bunch of Haskell, even some css. A bit much. I just wanted reloading. I want to build up to here myself. I don't want to understand all this just now. So back to bare bones and I'll try to figure out reloading later. Just good to know it's there.


