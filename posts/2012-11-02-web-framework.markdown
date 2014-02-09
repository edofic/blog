---
title: Hunt for a web framework that works
---

  -------------------------
  [![Never Internet Explorer](http://upload.wikimedia.org/wikipedia/en/e/ea/Never_Internet_Explorer.png)](http://en.wikipedia.org/wiki/File%3ANever_Internet_Explorer.png)
  Never Internet Explorer (Photo credit: [Wikipedia](http://en.wikipedia.org/wiki/File%3ANever_Internet_Explorer.png))
  -------------------------

I have this personal project I want to do that includes a web
application and I want to learn something. So I'm on the hunt for
language, environment and framework.

### Other stuff

I did some [PHP](http://www.php.net/ "PHP") a few years back
and definitely don't want to go there anymore. I also did some
[.NET](http://msdn.microsoft.com/netframework ".NET Framework") and it's
even part of curriculum here at FRI. But clicking on wizards in [Visual
Studio](http://www.microsoft.com/visualstudio/en-us "Microsoft Visual Studio")
feels weird to me. Not like development should be done. And I also use
GNU/Linux as my primary(and only)
[OS](http://en.wikipedia.org/wiki/Operating_system "Operating system"),
so that's out of the water. I did read about [java server
pages](http://en.wikipedia.org/wiki/JavaServer_Pages "JavaServer Pages")
and faces and even tried few things out. But luckily I didn't get to do
this project I was preparing for and I didn't need it. It looked ugly
anyway. I did some flirting with
[GWT](http://code.google.com/webtoolkit "Google Web Toolkit"), does that
even count as a web framework? 

### Node.js

I heard about [nodejs](http://nodejs.org/ "Node.js") quite some time ago
but I put off looking into it because my js was really rusty. But
recently I brushed up on my javascript skills(to do a "compiler" into
js) and gave it a shot. Node is good. It's fast, it's agile, it makes
you think in a different way. I was feeling empowered. I did some simple
stuff and I liked it.

### Static vs dynamic

Later I kinda got a job as [Ruby on
Rails](http://rubyonrails.org/ "Ruby on Rails") dev. And I hated it. So
I didn't take it. It would take up too much time anyway - I'm a student.
Ruby is okay. Rails is okay. But problem was the size of the problem.
Application we were buiding(a team of devs) was quite complex and I came
into existing(moderate size) ruby code base. Learning ruby and rails as
I go was fun, but navigating the project was pain in the ass. Of course
documentation was non existent and
[IDE](http://en.wikipedia.org/wiki/Integrated_development_environment "Integrated development environment")
couldn't help me because it didn't know. So a lot of regex searching and
walking around asking stuff. Also refactoring...Inevitable but hard. 

This cemented my opinion on static vs dynamic typing. (Static for
everything but a short script, more on that another time).

### Scala

Then I learned about the good parts of [static
typing](http://en.wikipedia.org/wiki/Type_system "Type system") through
scala and haskell. Doing web in haskell seems a bit intimidating(I will
give it a go eventually, I promise) so I roll with scala. I looks there
are two big names here. Play! and Lift. I watched a few talks and read
few blogs about both to see central points. 

Big difference seems to be their view on state. Lift goes for stateful,
Play for stateless. Play kinda seems like it has a bigger community, but
their documentation is stellar and they're now part of Typesafe stack.
No brainer then. Play it is.

### [Play! framework](http://www.playframework.org/ "Play Framework")

I dived into documentation. Reading samples and explanation about
infrastructure and
[APIs](http://en.wikipedia.org/wiki/Application_programming_interface "Application programming interface").
Samples really clicked with me - it felt like porn. No analogy, reading
elegant scala sources for a web app for my first time felt like I was
doing something naughty, like things shouldn't be that good.

Live reloading is great too. A friend of mine is a J2EE dev and he's
constantly nagging about build and deploy times. I get that near
instantaneous. And [compile
time](http://en.wikipedia.org/wiki/Compile_time "Compile time") checking
of routes and templates? Oh my god, yes. Bear in mind, compile time is
all the time. When I hit ctrl-s for "save all open files" I quickly see
if compiler has any complaints, even before I refresh the browser. 

I just did some experimenting with features...for a few hours.
Everything feels so simple but powerful. Why nobody told me about this
before?

Okay, it has to have some weaknesses but I didn't find them. Yet. And
that's what counts.

### [Heroku](http://www.heroku.com/ "Heroku")

Now this is just a cherry for the top of my cake. It took me two minutes
to deploy my hello world app, and that includes time needed to install
Heroku's tookkit. You just create an app and push to remote git repo.
Heroku detects it's a Play/Scala app and install dependencies. Rest is
done by SBT. And it just works. Hassle free deployment for developers.
Yay.

Now I have my stack and even a host. So I just need to write an awesome
service and generate traffic. [How hard can it
be?](http://www.google.si/url?sa=t&rct=j&q=&esrc=s&source=web&cd=6&cad=rja&ved=0CFEQtwIwBQ&url=http%3A%2F%2Fwww.youtube.com%2Fwatch%3Fv%3DnVE09yyznfc&ei=jAKUUIPEDdGN4gTR5YCwDg&usg=AFQjCNFEfy9P-ultVxu5ZkJgFIV4m1aODA&sig2=OX1AHn6wJRVV3cca7L1vnQ)
