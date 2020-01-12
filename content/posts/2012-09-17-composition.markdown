---
title: Pretty function composition in scala and asynchronous function composition on android
date: 2012-09-17
---

  ----------------------
  [![Surjective composition: the first function nee...](http://upload.wikimedia.org/wikipedia/commons/thumb/a/a2/Surjective_composition.svg/300px-Surjective_composition.svg.png)](http://commons.wikipedia.org/wiki/File%3ASurjective_composition.svg)
  Surjective composition: the first function need not be surjective. (Photo credit: [Wikipedia](http://commons.wikipedia.org/wiki/File%3ASurjective_composition.svg))
  ----------------------

Function composition is a nice way to sequence transformations on data.
For example in a
[compiler](http://en.wikipedia.org/wiki/Compiler "Compiler") you take
your source, [parse](http://en.wikipedia.org/wiki/Parsing "Parsing"),
check, optimize and generate code(in a nutshell). It's a linear series
of transformations -> perfect for [function
composition](http://en.wikipedia.org/wiki/Function_composition "Function composition")
In [Haskell](http://haskell.org/ "Haskell (programming language)") you
can use this beautiful
[syntax](http://en.wikipedia.org/wiki/Syntax "Syntax")
```haskell
compile = codegen . optimize .check . parse
```

leaving out the parameters and noting composition
as "." which kinda looks like ° used in math(if you squint a bit). In
scala you could do something like
```scala
val compile = codegen compose optimize compose check compose parse
```

Nearly there, I just want it to look a bit prettier. (Compose is a
method defined on Function1 trait) So I define an [implicit
conversion](http://en.wikipedia.org/wiki/Type_conversion "Type conversion")
from function to "composable"

    implicit def function2composable[A,B](f: A=>B) = new AnyRef{
      def -->[C](g: B=>C) = v => g(f(c))
    }

This creates an object and reimplements "compose" but I really like the
syntax:

    compile = parse --> check --> optimize --> codegen

[Functional programming](http://en.wikipedia.org/wiki/Functional_programming "Functional programming")
can be imagined as a [waterfall of data](http://swizec.com/blog/my-brain-cant-handle-oop-anymore/swizec/4320),
flowing from one function into the next, and the `-->` operator
represents this nicely. Let's go a step further. If the composition is
one-off, and this is a waterfall of
[DATA](http://en.wikipedia.org/wiki/Data "Data") it could be nice to
represent that. Something like
```scala
val result = source --> parse --> check --> optimize --> codegen
```

So now I'm taking a value and sending it through black boxes. Very nice.
Apart from the fact it doesn't work(yet!).
```scala
implicit def any2waterfall[A](a: A) = new AnyRef{
  def -->[B](f: A=>B) = f(a)
}
```

[Scala](http://www.scala-lang.org/ "Scala (programming language)")'s
awesome compiler can handle two implicit conversions with same method
names. Nice. You can even mix and match
```scala
val result = source --> (parse --> check --> optimize) --> codegen
```

This does the composition of parse, check and optimize into an anonymous
function, applies it to the source and then applies codegen to it's
result.

### Goin async

  ---------------------
  [![Image representing Android as depicted in Crun...](http://www.crunchbase.com/assets/images/resized/0001/4601/14601v1-max-450x450.png)](http://www.crunchbase.com/product/android)
  Image via [CrunchBase](http://www.crunchbase.com/)
  ---------------------

What about asynchronous calls? Can I compose those too? I think it's
possible with Lift actors(or with scalaz's?), but I needed to integrate
that into [Android](http://code.google.com/android/ "Android")'s
activities quite recently. Well I did not *need* to do it, but it was
quite a nice solution.

The usual way of doing things async in Android is with
the conveniently named AsyncTask. The problem is - you can't subclass it
in scala because of some compiler bug regarding varargs parameters.
Silly.

So let's do a lightweight(in terms of code) substitution. We can "spawn
a thread" using scala's actors. And activity can receive messages
through a Handler.
```scala
import android.os.{Message, Handler}

trait MessageHandler {
  private val handler = new Handler {
   override def handleMessage(msg: Message) {
     react(msg.obj)
   }
  }

  def react(msg: AnyRef)

  def !(message: AnyRef) {
   val msg = new Message
   msg.obj = message
   handler.sendMessage(msg)
  }
}
```

So in an activity that mixes in MessageHandler I can post messages to my
self. And I can do it async since Handler is thread safe
```scala
def react(msg: AnyRef) = msg match {
  case s: String = Toast(this, s, Toast.LENGTH_LONG).makeText()
  case _ => ()
}

...
//somewhere inside UI thread
import actors.Actor.actor
val that = this
actor{
  val msg = doExpensiveWork()
  that ! msg
}
...
```

Not the most concise way to write it, but I believe the most clear one.
Method doExpensiveWork is done in the background so it doesn't block UI
and it posts the result back as a message.

#### [Async](http://en.wikipedia.org/wiki/Asynchrony "Asynchrony") composition - finally

What I want to do now is use function composition to do something like
```scala
(input --> expensiveOne --> expensiveTwo) -!> displayResults
```

In other words, do "waterfall composition" in background using some
input I have now and post the result back into the UI thread to method
displayResults. That should be the magic of the -!> operator. Do the
left side in bg and post it to the right side async. I need a new trait
for that
```scala
trait MessageHandlerExec extends MessageHandler { outer =>
  override protected val handler = new Handler {
    override def handleMessage(msg: Message) = msg.obj match {
      case r: Runnable => r.run()
      case other: AnyRef => react(other)
    }
  }

  implicit def any2asyncComposable[A](a: => A) = new AnyRef{
    def -!>[B](f: A=>B) = outer ! new Runnable{
      def run() = f(a)
    }
  }
}
```

The trick here is using by-name parameters in the implicit conversion.
This delays the execution of a(which in example above would be a
waterfall and moves it into a worker thread.
