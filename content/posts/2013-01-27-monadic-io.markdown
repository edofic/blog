---
title: Monadic IO with ScalaZ
date: 2013-01-27
---

I just recently scratched the
surface with [scalaz](http://code.google.com/p/scalaz/). Think of it as
an additional standard library for scala that's
[FP](http://en.wikipedia.org/wiki/Functional_programming "Functional programming")
oriented. It provides a bunch of [type
classes](http://en.wikipedia.org/wiki/Type_class "Type class"),
instances for pretty much everything, some fancy data types, pimps(Pimp
My Library) for standard library collections, actor implementation and
probably some stuff I'm not aware of. I could really use a "map of
scalaz" - but I'll probably dive into source and scaladoc anyway. One
fancy feature that's not noted on their[Google Code
page](http://code.google.com/p/scalaz/) is
[IO](http://en.wikipedia.org/wiki/Io_%28programming_language%29 "Io (programming language)")
monad implementation.

---------------------
![Denali Landscape](http://farm7.static.flickr.com/6157/6183470322_cbbf4881d2_m.jpg "Denali Landscape")
Real World Tm - the thing your programs have to interact with (Photo credit: blmiers2)
---------------------

I've written a bit about monadic IO but let's recap. IO monad is a data structure that represents a tiny
language which let's you describe and compose IO actions without
actually performing them - allowing you to keep your functions pure and
compose/reuse better.

### Monadic [HelloWorld](http://en.wikipedia.org/wiki/Hello_world_program "Hello world program") with scalaz

I'll be using SBT to manage dependencies(scala and scalaz) so let's
create a project. We just need build.sbt file with this content(or use
my template: g8 edofic/scalaz-empty)

```scala
scalaVersion := "2.9.2"

resolvers += "Scala Tools Snapshots" at "http://scala-tools.org/repo-snapshots/"

libraryDependencies += "org.scalaz" %% "scalaz-core" % "6.0.4"
```

and then run "sbt console" to get a
[REPL](http://en.wikipedia.org/wiki/Read%E2%80%93eval%E2%80%93print_loop "Read–eval–print loop")
with scala 2.9.2 and scalaz 6.0.4. Pretty useful. Let's get to code now.

```scala
import scalaz._
import Scalaz._
import scalaz.effects._

val greeter = println("hello world").pure[IO]

//scalaz provides a helper for that so we can write
val greeter2 = putStrLn("hello from scalaz")

greeter.unsafePerformIO
```

First a bunch of imports to get scalaz magic and then the crucial line.
Scalaz provides an implicit conversion which allows us to call .pure[A]
on every value. This gets
[Monad](http://en.wikipedia.org/wiki/Monad_%28functional_programming%29 "Monad (functional programming)")[A]
instance from implicit scope and lifts the value into the monad. It's a
type class way (think dependency injection) to invoke monad
"constructor", a simple sample

```scala
scala> "hi".pure[Option]
res0: Option[java.lang.String] = Some(hi)
```

Scalaz also provides a helper putStrLn(and many more, including readLn)
for more succinct code. Note that lifting our println(by hand) does call
by name([lazy
evaluation](http://en.wikipedia.org/wiki/Lazy_evaluation "Lazy evaluation"))
so println is not actually invoked yet! To perform the IO action you
need to explicitly call unsafePerformIO. Intentionally verbose and scary
name to make you think twice where you perform your side effects. But
printline is not that interesting. Let't take a look at input. Type of
readLn is IO[String] and it's perform method returns a String. Notice
that IO is a monad so you can do map and flatMap on it. For example
readLn.map(_.toInt) gives you and action that will read and parse an
integer returning you something of type Int. FlatMap is more interesting
because it gives you composition. FlatMapping another IO action will
execute actions in sequence giving last action access to previous
actions' return values via closures. That's referred to as monadic
style. It's basically a pure way of doing imperative programming. A bit
more involved example

```scala
for{
  _    <- putStrLn("who are you?")
  name <- readLn
  _    <- putStrLn("hello " + name)
} yield ()
```

And a tiny calculator showing off composition

```scala
val inputInt = for{
  _ <- putStrLn("enter an integer")
  raw <- readLn
} yield raw.toInt

val adderIO = for {
  a <- inputInt
  b <- inputInt
  _ <- putStrLn((a+b).toString)
} yield ()

//and to run it
adderIO.unsafePerformIO
```
