---
title: Setting up for scala development on Android
---

  -------------------------
  [![Image representing Android as depicted in Crun...](http://www.crunchbase.com/assets/images/resized/0001/4601/14601v1-max-450x450.png)](http://www.crunchbase.com/product/android)
  Image via [CrunchBase](http://www.crunchbase.com/)
  -------------------------

  -------------------------
  [![Scala (programming language)](http://upload.wikimedia.org/wikipedia/en/thumb/8/85/Scala_logo.png/300px-Scala_logo.png)](http://en.wikipedia.org/wiki/File%3AScala_logo.png)
  Scala (programming language) (Photo credit: [Wikipedia](http://en.wikipedia.org/wiki/File%3AScala_logo.png))
  -------------------------

I've been developing for android more than a year now and a few months
in scala. So naturally I wanted to combine the two. But it's not dead
simple. This is kinda a tutorial an a reference if I ever forget how to
do this. It took me a few days to figure it all out. I tried maven, ant
with special config and sbt(I need to learn more about this one) but in
the end I just wanted fast solution integrated into my IDE.

So I use [IntelliJ
IDEA](http://www.jetbrains.com/idea/ "IntelliJ IDEA") community edition
for my IDE.  You should check it out, it's totally awesome. It's primary
a Java IDE but scala plugin rocks. It offers some more advanced text
editing capabilities, not like vim or emacs but enough for me. It also
brings up coloring and editing features that are language aware. So you
have a shortcut(Ctrl-W) to select semantically valid block. And press it
again to expand to next bigger valid piece of code. And stuff like
that. Real-time structure view is nice and there are some cool
refactorings. But scala
[REPL](http://en.wikipedia.org/wiki/Read–eval–print_loop "Read–eval–print loop")
is where fun begins. You get your module classpath pre-set and you get
**full editor capabilities in REPL**. Enough with advertisement(they
didn't pay me to do this) and let's get to work. 

### Prerequisites
-   JDK...duh!  I use
    [OpenJDK](http://openjdk.java.net/projects/jdk7/ "OpenJDK") 7, IDEA
    gives some warnings but it works like a charm
-   [Android
    SDK](http://en.wikipedia.org/wiki/Android_software_development "Android software development")
    and at least one platform
-   IntelliJ IDEA
-   scala distribution. I recommend you use latest stable release from
    [here](http://www.scala-lang.org/downloads) 

### Setting up

First install scala plugin. It's quite straightforward. Plugin
Manager->Browse repos->search for scala->select->ok.

Now actual setting up. I use global libraries for all my projects, you
can also put these into just Libraries and to that on per-project basis.

Open project structure(no project open) and go to Global Libraries. You
need to create two libraries containing jars from /lib/.

First scala-compiler
with scala-compiler.jar and scala-library.jar and then scala-library
with scala-library.jar and anything else you might need. Reason for
scala library in compiler is that compiler also relies on scala lib. I
needed quite some time to figure this out.

This whole process can be automated if you add scala to your project
when creating it but it's not possible with android so you need to know
how to do it by hand. 

### Creating a project
-   Project from scratch
-   add android module and configure it
-   now go to project structure. add scala facet to this module and go
    to its settings and set the compiler jar.
-   back to module and add dependency to global scala-library
-   set dependency to **provided**. This is important. Else it will try
    to dex whole library and you'll end up with "too  many methods
    error".

Now your project should compile. But not run.

### Running
Obviously not including scala library in the build means you need to
provide it in another way. For developing on emulator I customized it to
provide predexed scala library. [Excellent
tutorial.](http://zegoggl.es/2011/07/how-to-preinstall-scala-on-your-android-phone.html)

In a nutshell

    $ git clone git://github.com/jberkel/android-sdk-scala.git$ cd android-sdk-scala$ ./bin/createdexlibs
    $ bin/createramdisks
    $ emulator -avd ... -ramdisk /path/to/custom.img$ adb shell mkdir -p /data/framework$ for i in configs/framework/*.jar; do adb push $i /data/framework/; done

And reboot.

There is also (not so trivial) part about patching a device. However you
can try [Scala
Installer](https://play.google.com/store/apps/details?id=com.mobilemagic.scalainstaller&feature=search_result#?t=W251bGwsMSwxLDEsImNvbS5tb2JpbGVtYWdpYy5zY2FsYWluc3RhbGxlciJd) from
Play to do this. I had some success and some failures.

Now the app should run on your device.

### Deploying

Well it doesn't work on other devices right now. For export you need to
change scala-library dependency back to **compile**to include it into
the build. Trick now is to enable ProGuard to remove unnecessary methods
and classed to fit the jar through dexer. You do this in
Tools->[Android](http://code.google.com/android/ "Android")->Export.
Select ProGuard and your config. I got mine from jberkel's repo. That's
it. Sadly this export takes quite some time. Scala's standard library is
not a piece of cake afterall(actually *it is* a cake). Minute and a
half on my machine for small apps. So I only to this for testing on
other phones and deployment.

### Faster compilation

Compiling with scala-libray set to provided is much faster but not fast
enough for me. I want to be doing stuff not [waiting for it to
compile](http://xkcd.com/303/). 

Turns out compiler is the big time sucker(and I'm being Capt. Obvious).
Afterall scalac is not known for it's speed.

Enter **FSC**or Fast Scala Compiler. This is a scala compiler running in
the background having everything preloaded and just does incremental
compilation. It even comes with standard scala distribution and is
supported by IntelliJ IDEA. Great. 

To set it up just head over to Project Structure->Scala facet and
select Use FSC. And then immediately click Setting to access Project
Settings and set compiler jar for the compiler.

Success. Scala builds are now on par(or even faster!) than java ones. 

No more fencing for me.
