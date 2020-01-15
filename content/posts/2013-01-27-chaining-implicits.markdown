---
title: Chaining implicit conversion in scala
date: 2013-01-27
---

Today I was hanging
in the
#[scala](http://www.scala-lang.org/ "Scala (programming language)")
[IRC
channel](http://en.wikipedia.org/wiki/Internet_Relay_Chat "Internet Relay Chat")
and somebody came along(forgot the nick, sorry) and asked about some
[compilation
error](http://en.wikipedia.org/wiki/Compilation_error "Compilation error").
I deduced he was trying to chain implicit conversions. And this doesn't
work. Else compilation would take forever and would also compile some
wrong code by inserting long strings of implicits. But then somebody
else responded(I think nick started with d) and gave a solution to
implicit chaining. But I'm not giving it away yet, you'll have to read a
bit more.

### The problem

Let's say I have some
[classes](http://en.wikipedia.org/wiki/Class_%28computer_programming%29 "Class (computer programming)")
that just add semantics to values(think labeled boxes)

```scala
case class Foo(a: Int)
case class Bar(b: Int)
case class Baz(c: Int)
```

And because they look similar we may want to do some [implicit
conversion](http://en.wikipedia.org/wiki/Type_conversion "Type conversion").
Let's say that that's needed is conversion from Bar to Foo and from Baz
to Bar.
```scala
implicit def bar2foo(bar: Bar) = Foo(bar.b)
implicit def baz2bar(baz: Baz) = Bar(baz.c)
```

On a level this also implies that Baz can be seen as Foo. But when you
try it
```scala
println(foo.a)
println(bar.a)
println(baz.a)
```

you get a compile error

    Implicits.scala:18: error: value a is not a member of Implicits.Baz
    println(baz.a)
    ^
    one error found

### The solution

As said, scala
[compiler](http://en.wikipedia.org/wiki/Compiler "Compiler") doesn't try
to traverse the graph of implicit conversions and therefore doesn't
figure out how to insert needed conversions here. But here lies a handy
catch. What if the conversion wasn't from Bar to Foo but from something
Bar-like to Foo. Then the compiler would know that this chaining is
intended and is (probably) not a dead end. Good news: scala lets you
express that. It's called a view bound. You make the method generic and
limit input type to something that can be seen(implicitly converted to)
as a Bar. Here's the code.
```scala
implicit def bar2foo[T <% Bar](bar: T) = Foo(bar.b)
```

All the magic lies in the `<%` operator. That's the view bound. When you
type `baz.a` the compiler sees that a property is available on the Foo
class and that there's an implicit conversion to Foo from something
Bar-like. Fortunately there's a conversion from Baz to Bar in scope so
it can invoke the method. Inside the method it tries to get the b
property from baz but it's not available. So it checks the conversions
and sees that a conversion from Baz to Bar exists in the scope and it
satisfies the need for b. It inserts this conversions and compilation
happily chucks along.

![Nest all the implicit conversions!](/images/chaining-implicits/all_the.jpg)

Please don't do that.
Implicits can make code hard to read even without nesting. So use with
care. Scala is a very powerful language but great power comes with great
responsibility.
