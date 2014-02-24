---
title: Comparing Haskell Web Frameworks
---

## Intro
Lately I got sucked into Haskell. When I first saw it it looked like it might be a bit impractical for real-life projects but this prejudice faded away slowly. Now I'm at the point where I have an idea for a web application and I'd like to do it in Haskell. After a quick search I found many frameworks and libraries that I might use. So I decided to do some exploration and implement a bunch of stuff with different technologies. This way I don't have to believe others and I can decide for myself what feels right to me.

All the code is available at [github](https://github.com/edofic/haskell-web). What follows is my personal opinion.

## Setting Up
First up is the set up of all used frameworks. I created a hello world project for each to get a sense of things. Here are some comments. 

Each project first got a cabal sandbox and `.cabal` file. 

### Happstack
Tutorial sure is short on how to get going. Let's see if I can figure this out. I added `happstack` and `happstack-server` to my dependencies. And in my `Main.hs` I've written the imports and one line of code that says return "hello world" for all requests. `cabal run` and voila, hello world page at `localhost:8000`. Things sure are simple. And I got out a single binary that I can run - or deploy. But this means no autorecompile and reload I grew to like from PlayFramework(or any dynamic framework really).

### Snap
Snap says in their tutorial to use their tool to generate the project skeleton. Ok I install the `snap` package(into sandbox). While it's installing I might stress that I prefer convention over configuration and this implies minimal initial set up, preferably something I can do by hand in a minute - like Happstack. I guess this might divide up libraries from frameworks.
OK. Snap installed. `snap init barebones` now created a skeleton. And there isn't much. It created `snap.cabal` that has some useful dependencies already in place and some GHC options. There is also main source that has an example how to do routing. Great. But I can see that a minimal example would still be just a few lines. Something I could do by hand. Great.

How about that automatic recompilation? Tutorial says I have to initialize a full fledged project instead of bare bones one - `snap init`. Whoa this reated a bunch of stuff. There are some templates in a special templating languagem bunch of Haskell, even some css. A bit much. I just wanted reloading. I want to build up to here myself. I don't want to understand all this just now. So back to bare bones and I'll try to figure out reloading later. Just good to know it's there.

### Scotty
Very similar to Happstack. Just a dependency and a oneliner to serve on localhost. Also meaning no autoreload. 

### Yesod.
Magic. You install `yesod-bin` package that gives you CLI tools. And then `yesod init` takes you through a wizard to setup the project - kind of cabal init but yesod specific. Then it proceeds to generate a huge project for you(huge as in starting template). It already has styles, templates and even tests. Great. Examples to learn from. But most of this is not needed. Consulting the Yesod book(Yesod has great documentation!) I found out that the minimal HelloWorld example is just using Yesod as a library and quite minimal in terms of lines of code. The major difference with others being liberal use of quasiquoters and template haskell. 

Yesod also comes with automatic recompilation and reloading. It's set up by the scaffolding generator and not trivial. But it's there for me to figure it out some day. 

## Routing
The process of going from a URL to something that eventually produces a response. Quite important in my opinion and a taste of the whole library as this is the place where everything comes together. 

### Happstack
Happstack gives you a monad transformer stack. That's all. Said monad is also `MonadPlus`. When you append together two instances the first one will be tried and if t fails it will try the other. To do routing you just write instances that filter on desired route and then perform the required action and then `msum` all of them together. Of course you could do it a bit smarter and expand this into a tree. On each level you just match one part of the path and maybe do some action that is required by all children actions. 

### Snap

Snap has *snaplets* to pack together functionality which is a nice idea. There is the `Snap` monad that provides basic behaviour, for more you use custom snaplets. But for the basic routing testing(what I did) `Snap` was enough. 

Again there is the notion of composing actions and trying something else on failure. But instead of `MonadPlus` you use `Applicative`'s  `<|>` to say *try A and if it fails try B*. But there is also `route` combinator that allows you to take a bunch of actions and compose them filtering by path for separate action. And out comes a single action that knows what to do on which path(and/or method). 

### Scotty

There are two monads(monad transformer stacks to be precise) that you use to define a Scotty app. First there is `ScottyM` that describes configuration, does set up and declares routes. You get nice helper methods like `get` and `post` that take the URL and an action in monad `ActionM`. This is where you go from request to response. I like this clean design that forces you to separate general set up and routing from implementation of each handler.

### Yesod

There's quasiquoter that takes custom(nice and terse) syntax for routes and binds them to *resources*. Which are automatically declared types. You then define handlers following some convention for their name and they will be automagically picked up and used for correct routes. And any screw-ups are compile time since everything gets type-checked. The price you have to pay is a bit more boilerplate when you want your routes to be more dynamic. 

For example if you create an abstract REST data source and apply it to several models you need to do some trickery. In others you would just define a function that takes partial path as parameter and then defines new routes. In Yesod there are probably more options but I went the safe route and declared a subsite that does all this and then added it to master routing. There is some machinery to write but all in all not that bad. For me it seems worth it as you now get type-safe routing(and rendering of URLs) to all your resources. Neat.

## Documentation 

Documentation is good to have so you don't need to read whole source and you can also learn some good practices this way. Especially important is when getting started. So here are some notes on documentation. 
Please note that I didn't try very hard and may have missed something(if so please leave a comment).

### Happstack
Has great [online docs](http://happstack.com/page/view-page-slug/3/documentation). Very nice is the [Happstack Book](http://happstack.com/docs/crashcourse/index.html) also available in pdf and some other formats(follow the previous link and check introduction).

### Snap
Is a bit poorer in this regard. [Here's](http://snapframework.com/docs) the docs. There's no comprehensive document(book) but separate modules are nicely covered. 

### Scotty
I dare say Scotty is so tiny it doesn't need a book. [Haddock docs](http://hackage.haskell.org/package/scotty-0.6.2/docs/Web-Scotty.html) and some playing around in GHCi was enough for me to grok it. It's just a simple layer on top of WAI and Warp. 

### Yesod

Has a [book](http://www.yesodweb.com/book). You can read it online or buy a paper edition. And it's *very* comprehensive. Contained just about everything I needed and more.


## Impressions

I did a bit of development and skimmed through official documentation reading a fair share(most?) of it. So I think I have seen enough to do  some judgement though it might change after more usage.

Happstack rubs me the wrong way. Sorry Happstack devs, you did a great job but it just doesn't feel right to me. I think all four candidates are very capable and have many interchangeable parts, so it comes down to style and personal preference. And I don't like Happstack. **UPDATE** Some readers expressed concern that this paragraph is too subjective and vague. In retrospective I concur. So let me explain what I don't like. I don't like their the architecture of Happstack applications. It lets you write bad code due to not enforcing separation between routing, controllers and models. No I am not disciplined and will try to cheat so I like tools that make this very hard(it's always possible to cheat). Also some impressions from its crash course book. A lot of documentation(this give an impression on framework focus) is geared towards templating and processing forms. Even a whole chapter on generating JavaScript - a concept I disagree with. On the other hand it was harder to find documentation on how to do RESTful content - which is more aligned with my interests. This is why it doesn't feel right to me. YMMV.

Scotty looks like it's on a similar level but even leaner. However this means you get to pick the best library(to your taste) for each functionality(templating, persistence etc.). And you get yourself a tailor made "framework" if this is an appropriate word. Of course this is a double-edged sword; you **have** to pick a library(or write it) for everything as there is almost nothing included. As an ArchLinux user I appreciate the minimality and I believe some would build even big applications with this. Also Scotty is *cool* with it's Star Trek references. I feel like Scotty is a great choice if you just want to throw up an API for something you wrote or render it in a few views but not for more complex apps. 

Snap mitigates this with snaplets. This way you can encapsulate and compose functionality. You could of course write something like this with Scotty yourself but why bother if you can just use it. And Snap also has an ecosystem of ready made snaplets for you to use. Stuff like Fay support or MySQL or MongoDb persistence layer. Bunch of integrations you don't need to do manually. Neat. Sounds similar to PlayFramework's plugins or Rails' gems. Definitely have to play with this more.

Yesod is opinionated and quite big. There is support for *a lot* of stuff. Like subsites, widgeds and persistence with multiple backend. And there's even more in other packages on Hackage. You get a familiar MVC model with *Front Controller* - routes in a single file. This is sometimes a pain and sounds clumsy(and sometimes it is) but it think it's worth it. You have control over your entry points. Period. You can quickly trace requests to handlers and models. It's harder to generate repetitive routes but it can be done. There are a lot of differences between the four but I think this is the biggest one. And I think this is the better way to do it for larger code bases. When there's too much to hold it all in your brain it's very nice to have your application's API laid out in one file. 


## (Lack of) Conclusion

I'll continue implementing and comparing stuff in [my github repository](https://github.com/edofic/haskell-web). Feel free to help or point me in the right direction.

If my opinion changes I'll post part two. Until then I'll stick with Yesod and play around with Snap some more. And I do think Scotty is a great micro framework.

I'll definitely implement a larger application with Yesod and then compare it to frameworks in other languages I worked with. 
