---
title: Setting up for better Scala development on Android
---

![](http://www.crunchbase.com/assets/images/resized/0001/4601/14601v1-max-450x450.png)

I did [a tutorial how to set up
everything using
IntelliJ](/posts/2012-10-22-scala-android.html)
a while ago. I still think [IntelliJ
IDEA](http://www.jetbrains.com/idea/ "IntelliJ IDEA") is awesome and you
should use it(it has a free and open version) but I've found a better
way. No more clicking around in wizards...it's config file time. Relax -
it's quite simple. I recommend installing the Typesafe-stack. It gives
you the two tools listed below in nice packaged form with updates. And
that's about it. I'm quite fond of package managers => one place to
update everything.

### SBT

[Simple Build
Tool](https://github.com/harrah/xsbt/ "Simple Build Tool"). Exactly what
the name says. Simple. But boy it's powerful. And it has [a
plugin](https://github.com/jberkel/android-plugin)for doing android
development. Batteries included(managing emulator & packaging all from
one place). And it prefers convention over configuration. Great. All
this leads to quick results.

### Giter8


[Maven](http://maven.apache.org/ "Apache Maven") has archetypes and
that's a great feature since it enables you to just start developing
your application without messing with build settings to get the damn
thing to compile. SBT doesn't have that. But there's a third-party tool
that does something possibly even better. Yes, giter8 or g8 for
short(that's the command line too). A tool that fetches templates from
github(or any other git repo in recent versions) asks you for few
parameters and creates a project for you. Super simple to use and not
that hard to create your own templates.

### Let's get cracking

Use g8 to create an app and configure it

    g8 jberkel/android-app

    Template for Android apps in Scala 

    package [my.android.project]: com.edofic.demoapp
    name [My Android Project]: DemoApp
    main_activity [MainActivity]: 
    scala_version [2.9.2]: 2.9.1
    api_level [10]: 
    useProguard [true]: false
    scalatest_version [1.8]: 

    Applied jberkel/android-app.g8 in demoapp

This is from interactive session. Stuff after colons is my input...eh I
won't try to  explain. Just try it. Explanation on parameters is in
place though. Template is from the author of the plugin; I used default
activity name because I don't care. Scala version is 2.9.1 because
that's what's preinstalled on my phone(more on that later) and ProGuard
is off to exclude scala standard library from the dexing step. See SBT
is know that scala standard library doesn't fit through dexing and
excludes it. But it does some magic too. Whole packaging process is 2-3x
than ant/maven/eclipse/idea. And again twice faster if you use 2.10(RC5)
instead of 2.9.x. It gets under 10s on a good machine. Sadly my laptop
is not that fast so I use predexed scala library to get under that
magical 10s. Essentially I'm telling the
[compiler](http://en.wikipedia.org/wiki/Compiler "Compiler") not to
bother with the library and take care of it myself. There are two ways
of doing that. Note that you only need to do this once per device.

### Patching the device/emulator

Modify boot classpath and add libraries. Some nasty bootimage
manipulation for real devices. Simple script for emulator. What I use
for most development. See[here how to do
it](http://zegoggl.es/2011/07/how-to-preinstall-scala-on-your-android-phone.html).

### Shared libraries

This is how [Google](http://google.com "Google") provides libraries for
their maps
[API](http://en.wikipedia.org/wiki/Application_programming_interface "Application programming interface")
on android. There's even an[app on Play
store](https://play.google.com/store/apps/details?id=com.mobilemagic.scalainstaller&feature=search_result#?t=W251bGwsMSwxLDEsImNvbS5tb2JpbGVtYWdpYy5zY2FsYWluc3RhbGxlciJd)
that does that for you(requires root). More info on their [github
page](https://github.com/jbrechtel/Android-Scala-Installer). The
downside is you need to include a few lines into your manifest(just
development, you take it out for the release). I need to do more
research on this - maybe gonna post a guide when I figure out how to
install custom versions of scala like that.

### Compiling and running

Here comes a quirk. Plugin doesn't play nice with latest version of sbt.
Fortunately sbt is capable of  running different versions without any
hassle. You can provide "-sbt-version 0.12.0" every time or create a
project/build.properties with "sbt.version=0.12.0" in it(or use my
template - g8 edofic/android-app). This starts up sbt and does
compile-deploy-start operation.

    sbt
    ....
    android:start-debug

Checkout the plugin page for list of available commands.

### Why SBT?

I've only been converted recently but I've already found a bunch of nice
stuff I can do now. Firstly..dependency management. Now it's a breeze.
Typed resources via source generation - no more casting and exceptions
on findViewByID(see plugin github page for more). You can easily plug in
your own generation stuff too. And there's nifty feature I really
like: continuous compilation. In sbt console you write "~ compile" and
sbt will run incremental compilation triggered by file changes.

### Integrating with IDEA

Some people prefer programming in a text editor. I don't. So I like
[IDE](http://en.wikipedia.org/wiki/Integrated_development_environment "Integrated development environment")
integration. Luckily there's an[SBT plugin that generates IDEA
project](https://github.com/mpeltonen/sbt-idea) for seamless
interaction. You can even run SBT console inside idea for minimal window
switching. Usage(from plugin github page): Add the following line to
~/.sbt/plugins/build.sbt or PROJECT_DIR/project/plugins.sbt

    addSbtPlugin("com.github.mpeltonen" % "sbt-idea" % "1.2.0")

And then you can use "gen-idea" command in sbt console to generate idea
project files.

### Integrating with Eclipse

Same story goes for eclipse. Plugin can be found
[here](https://github.com/typesafehub/sbteclipse). To use it: Add
sbteclipse to your plugin definition file. You can use either the global
one at ~/.sbt/plugins/plugins.sbt or the project-specific one at
PROJECT_DIR/project/plugins.sbt:

    addSbtPlugin("com.typesafe.sbteclipse" % "sbteclipse-plugin" % "2.1.1")

and to generate project file you use "eclipse" command in sbt console.

### Links

Gathered for conveniance

* Typesafe stack [http://typesafe.com/stack/download-agreed](http://typesafe.com/stack/download-agreed)[](http://typesafe.com/)
* SBT [http://www.scala-sbt.org/](http://www.scala-sbt.org/)
* Giter8 [https://github.com/n8han/giter8](https://github.com/n8han/giter8)
* SBT anroid plugin [https://github.com/jberkel/android-plugin](https://github.com/jberkel/android-plugin)
* plugin author's template [https://github.com/jberkel/android-app.g8](https://github.com/jberkel/android-app.g8)
* sbt eclipse plugin [https://github.com/typesafehub/sbteclipse](https://github.com/typesafehub/sbteclipse)
* sbt idea plugin [https://github.com/mpeltonen/sbt-idea](https://github.com/mpeltonen/sbt-idea)
* scala installer for android [https://play.google.com/store/apps/details?id=com.mobilemagic.scalainstaller](https://play.google.com/store/apps/details?id=com.mobilemagic.scalainstaller)

