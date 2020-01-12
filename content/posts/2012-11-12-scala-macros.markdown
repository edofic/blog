---
title: Cool Monday - Scala Macros
date: 2012-11-12
---

  ----------------------
  [![Garden flower](http://upload.wikimedia.org/wikipedia/commons/thumb/b/b6/Garden_flower_.jpg/300px-Garden_flower_.jpg)](http://commons.wikipedia.org/wiki/File%3AGarden_flower_.jpg)
  Macro shot (Photo credit: [Wikipedia](http://commons.wikipedia.org/wiki/File%3AGarden_flower_.jpg))
  ----------------------

For me the highlight of this week was discovering Bootstrap. I heard of
it before but never looked into it. Probably because I wasn't doing web
stuff. The thing is bloody awesome. Back on topic.

[Scala](http://www.scala-lang.org/ "Scala (programming language)") 2.10
RC2 was released this Friday. Considering 2.9 had 3 RC releases, 2.10
final is probably quite near. And it brings some awesome features. One2
of them are
[macros](http://en.wikipedia.org/wiki/Macro_%28computer_science%29 "Macro (computer science)")


### Macros

So what are macros basically? Code executed at [compile
time](http://en.wikipedia.org/wiki/Compile_time "Compile time"). And
that's about it. 

So what is so great about that? Well you can do AST modifications and
stuff that gets quite into
[compiler](http://en.wikipedia.org/wiki/Compiler "Compiler")-plugin
territory in your regular code. That means you can do pretty advanced
stuff and do abstraction with performance. Oh yeah, you can also do
[type checking](http://en.wikipedia.org/wiki/Type_system "Type system")
and emit compile-time errors. Safety first kids!

### Usages

SLICK uses macros(in the experimental
[API](http://en.wikipedia.org/wiki/Application_programming_interface "Application programming interface"))
to transform scala expressions into
[SQL](http://www.iso.org/iso/catalogue_detail.htm?csnumber=45498 "SQL").
At compile time! ScalaMock uses it to provide more natural API for
testing. As said, you can use it for [code
generation](http://en.wikipedia.org/wiki/Code_generation_%28compiler%29 "Code generation (compiler)")
or validation at compile time. Good library design will be able to
minimize [boilerplate
code](http://en.wikipedia.org/wiki/Boilerplate_code "Boilerplate code")
even further now. And some people will argue that macros make scala even
harder language.

### Type Macros

This is the best part for me. But unfortuntely it's not implemented yet.
Or at least not dcumented. There are some methods with suspisious names
in the API but no useful documentation. In all presentations this is
referred to as "future work" but I still have my fingers crossed it
makes it into final release.

So what's the fuss? Automatically generated types. Large scale code-gen. 

As in ability to programatically create types at compile time. As a
consequence you can create whole classes with bunch of methods. And I
already have a [use
case](http://en.wikipedia.org/wiki/Use_case "Use case") of my own. I
want to make a typesafe ORM for Android that's super fast. I did
[YodaLib](https://github.com/edofic/YodaLib) ORM while back. It uses
reflection(although it's fast enough usually) and provides a Cursor that
lazily wraps rows into classes. And you need to make sure by hand that
your class coresponds to columns of your result set. Not very safe. I
had an idea to make static-typed safe inferface for the database when I
first heard about HList. You would do projection as a
[HList](/posts/2012-10-29-hlist-shapeless.html) and
result rows would be lazily wrapped into HLists. But using them for
wrapping every row(possibly filling data in with reflection) would be a
performance penalty. Not to mention a mess to implement. Now consider
using a macro to generate code for wrapping. It would be no slower than
accessing columns by hand. And a type macro would automatically create a
case class for given projection. Heavens. I'm just waiting for official
documentation on macros...this is a tempting project.


### Documentation

Here's [scalamacros.org](http://scalamacros.org/) which gives some
information. Also s[ome quite useful
slides](http://scalamacros.org/talks/2012-04-28-MetaprogrammingInScala210.pdf).
I hope now that 2.10 is in RC things stabilize, because in the milestone
releases api was changing constantly. [Nightly API
scaladoc](http://www.scala-lang.org/archives/downloads/distrib/files/nightly/docs/library/index.html)....
Proper documentation is apparently [coming
soon](http://docs.scala-lang.org/sips/pending/self-cleaning-macros.html),


### Le Code

A use case for macros, loop unrolling.

Below is a trivial sample of repetitive code.
```scala
class Manual{
  def show(n: Int){
    println("number "+n)
  }

  def code(){
    show(1)
    show(2)
    show(3)
    show(4)
    show(5)
  }
}
```

We can deal with repetition writing a loop. ([Higher order
function](http://en.wikipedia.org/wiki/Higher-order_function "Higher-order function")
really)
```scala
for( i <- 1 to 5 ) show(i)
```
But this doesnt generate the same AST(and bytecode)!
Protip: use

    scalac -Xprint:parser -Yshow-trees Manual.scala

to see AST after parsing.
Sometimes(rarely!) you want to unroll the loop to produce same [byte
code](http://en.wikipedia.org/wiki/Bytecode "Bytecode") as typing all
the iterations by hand.
```scala
Macros.unroll(1,5,1)(show)
```

With a proper unroll macro defined. I spent an hour to come up with this
implementation...and then scalac started crashing on me... Is there
something terrible in my code?
I gave up and went on to do useful stuff...But macros hear me! I'll be
back.
```scala
import reflect.macros.Context
import scala.language.experimental.macros

object Macros {
  def unroll(start: Int, end: Int, step: Int)(body: Int => Unit) = macro unrollImpl
  def unrollImpl(c: Context)(start: c.Expr[Int], end: c.Expr[Int], step: c.Expr[Int])
                (body: c.Expr[Int => Unit]): c.Expr[Any] = {
    import c.universe._
    val Literal(Constant(start_value: Int)) = start.tree
    val Literal(Constant(end_value: Int)) = end.tree
    val Literal(Constant(step_value: Int)) = step.tree
    val invocations = Range(start_value, end_value, step_value) map { n =>
      val n_exp = c.Expr(Literal(Constant(n)))
      reify{((body.splice)(n_exp.splice))}.tree
    }
    c.Expr(Block(invocations:_*))
  }
}
```
