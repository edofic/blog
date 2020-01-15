---
title: Union types in scala
date: 2013-01-27
---

I've done some research  a while
ago on union types and found a nice [implementation by Miles
Sabin](http://www.chuusai.com/2011/06/09/scala-union-types-curry-howard/) but
it only works for declaring types of function parameters. And you can
also do this with with [type
classes](http://en.wikipedia.org/wiki/Type_class "Type class"). What do
I mean with "only function parameters"? In "everything is a function"
kind of view there are three places to put types

1.  function parameters
2.  value(val or let binding in haskell and the like)
3.  function [return
    type](http://en.wikipedia.org/wiki/Return_type "Return type")

Even though Miles' encoding with [Curry-Howard
isomorphism](http://en.wikipedia.org/wiki/Curry%E2%80%93Howard_correspondence "Curry–Howard correspondence")
is ingenious it only applies to point 1. Let's fix that! Oh yeah you
could also use Either but that adds up boilerplate(even with implicits!)
and packing. And I want my union types unboxed.

### Enter type tags

I first saw this idea in Scalaz and immediately clicked with me. The
idea is to combine structural types with existing types. You don't touch
the value, just enhance the compile-time type with a tag. Structural
types are, well types that care about structure not name(JVM has nominal
type system) and (sadly) use reflection to do stuff at runtime. You can
even use it do do clean-ish method invocation with reflection.
```scala
type HasFoo = { def foo(): Bar }
(a:Any).asInstanceOf[HasFoo].foo()
//or even
(a:Any).asInstanceOf[{def bar: Int}].bar
```

But at compile-time they play along very nicely. Only concern if you
don't want reflection is not to use anything at runtime. But you can use
type members! They only appear at compilation so this works out perfect.
From scalaz

```scala
type Tagged[U] = {type Tag=U}
type @@[V,T] = V with Tagged[T]

type ExampleTaggedType = Foo @@ Bar

implicit class Taggable[A](value: A){
  def tag[T] = value.asInstanceOf[A @@ T]
}

val a = SomeType @@ Baz = someValue.tag[Baz]
```

Tagged is just a conversion between a regular generic and a structural
type with Tag member Type. This allows you to define @@ - that's simple
too: it takes the type and mixes in the structural tag. Foo @@ Bar now
means type Foo but with tag Bar, the only difference between them is
that values of type Foo @@ Bar now have a member type Tag that will
equal Bar. And I threw in an implicit for easier tagging.

### Subtyping with existentials

We need one more tool to get everything in place - subtyping. I mean
subtyping as in "if A is subtype of B and C is subtype of D then A @@ C
is also subtype of B @@ D". Technical term for this is covariance. And
with classes is done by simply adding pluses to [type
signature](http://en.wikipedia.org/wiki/Type_signature "Type signature")
like type @@[+V,+T]=... but this doesn't work for types. You get a
compile error as type parameters are used in non-covariant positions.
Luckily we can work around this with use-site subtyping. Instead of
using type A @@ B you can use X @@ Y for some types X and Y. Quite
literaly with this syntax

```scala
type Something = X @@ Y forSome {type X; type Y}
```

And these are [existential
types](http://en.wikipedia.org/wiki/Type_system "Type system"). Because
X and Y have to exist. Simple. There is some sugar with _(as usual) but
it doesn't work at all places; let's not get into detail.

### Either, or should I say Union

Now that we have the necessary foundation lets take a look at making
unboxed version of Either from scala standard library. Conversion is
pretty simple: instead of boxing into Left and Right, tag with Left and
Right. And instead of being of type Either a value will be some type
tagged with Either. I first implemented this with my type hierarchy only
to realize I reimplemented(without real functionality) Either.

```scala
type Or[A,B] = @@[_, _ <: Either[A,B]]

implicit def any2taggedLeft[A](a:A): A Or Nothing = a.tag[Left[A,Nothing]]

implicit def any2taggedRight[A](a:A): Nothing Or A= a.tag[Right[Nothing,A]]
```

There is Or type(intended for inline use) that equals something tagged
with something that's a subtype of Either[A,B] - this is anonymous use
site subtyping, a workaround for not being able to make tags covariant.
And that's pretty much all there is to it. The two implicits are just
for automatic tagging so you can use regular types and compiler will tag
them for you to type check union types. And when this code compiled and
worked I started cheering. Let me replace my current lack of enthusiasm
by use examples to let you fully gasp the implications on your own.

```scala
def id[A](a:A):A=a //a helper for testing parameter types

val a: Int Or String = 1
val b: Int Or String = "hi"
val c: Int Or String = 1.2 //does not compile

id[Int Or String](1)
id[Int Or String]("hi")
id[Int Or String[(1.2) //does not compile

def f(n: Int): Int Or String =
  if(n>0)
    n: Int Or String
  else
    "n is negative": Int Or String
```

As you can see all three requirements from the intro are satisfied while
keeping values unboxed. There are some rough edges unfortunately.
Compiler refuses to insert tags without explicit type annotations that
inform it to insert implicit conversions. Which is kinda weird because
it work for regular Either. Still looking into that. And the real
letdown is pattern matching. You can't match against types...compiler
just complains about these types not being possible. And it's kinda
right since String isn't in fact a subtype of Or[String,Int]. Luckily
there's a workaround that's not too ugly.

```scala
(a:Any) match {
  case _:String => "string"
  case _:Int => "int"
}
```

I have an idea how to solve both issues but it involves forking scala
standard library and is thus (much) less portable. Of course if anyone
has an idea how to convince the compiler to accept (fake) subtypes or
how to steer [type
inference](http://en.wikipedia.org/wiki/Type_inference "Type inference")
please let me know.
