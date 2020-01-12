---
title: Making a programming language Part 4 - Hello World
date: 2012-09-01
---

  ------------------
  [![printf(](http://farm8.static.flickr.com/7153/6691167811_440ed057ba_m.jpg)](http://www.flickr.com/photos/90209911@N00/6691167811)
  printf("hello, world\\n"); (Photo credit: [isipeoria](http://www.flickr.com/photos/90209911@N00/6691167811))
  ------------------
[Table of
contents](/posts/2012-08-29-creating-a-language-1.html),
[Whole project on github](https://github.com/edofic/scrat-lang)

What good is a language if you cannot do a [Hello World
program](http://en.wikipedia.org/wiki/Hello_world_program "Hello world program").
Every tutorial on every language I ever read has a Hello World in it
somewhere, even if it's a convoluted and sarcastic one.

So what do I need?
-   a way to print stuff to console 
-   strings

In that order. Since this is more of a math lang for now my first hello
world can just print 1 - arguably the simplest number.

### Println

This one is super easy. Just wrap scala's println. Here's the whole
code
```scala
lazy val sprint: FunctionVarArg = {
  case lst: List[_] =>
    lst --> mkString --> println
  case other =>
    throw ScratInvalidTypeError("expected a commaList but got " + other)
}
```

And then I put this as "println" into StdLib's map = global namespace.
Since I already have support for functions I just allowed them to do
side effects(semantics-no code).\
The `-->` operator might have confused you. It's in an [implicit
conversion](http://en.wikipedia.org/wiki/Type_conversion "Type conversion")
I defined:
```scala
implicit def any2applyFunc[A](a: A) = new AnyRef {
  def -->[B](f: A => B): B = f(a)
}
```
This is the so called [pimp my library pattern in
scala](http://www.artima.com/weblogs/viewpost.jsp?thread=179766). It
adds the `-->` operator to all objects. This operator takes another
single argument function, applies it and returns the result.
So I can have
```scala
lst --> mkString --> println
```
instead of
```scala
println(mkString(lst))
```

Think of it like
[Haskell](http://haskell.org/ "Haskell (programming language)")'s `$`
operator even if it works in a different way. Oh yes, mkString is
another function that I put into StdLib that takes a List(takes Any but
does pattern matching) and returns `List.mkString`

And now I have my Hello Math World
```scala
println(1)
```

### Strings

First I have to parse strings
```scala
private def string: Parser[SString] = "\".*?\"".r ^^ { s =>
  SString(s.substring(1, s.length - 1))
}
private def value: Parser[Expression] = number | string | identifier
```

SString is just a case class wrapper for string that extends Expression
so I can use it in other expressions.
When I was first writing this I forgot to add the action combinator to
strip of the quotes and was greatly mistified by all strings being
wrapped in "". I even spend half an hour debugging my evaluator before
it dawned on me.
I believe the evaluation of this is trivial.

Now I have a REAL hello world:
```scala
    println("Hello world")
```
Well...no. Hello world is a program you run. I need an interpreter to
run files.

### [Interpreter](http://en.wikipedia.org/wiki/Interpreter_%28computing%29 "Interpreter (computing)")

Not quite that different from
[REPL](http://en.wikipedia.org/wiki/Read–eval–print_loop "Read–eval–print loop").
In fact it's just REL: read, evaluate, loop. Printing is now only with
explicit println calls. And I don't need to catch exceptions and return
into the loop.  Whole code for the interpreter
```scala
object Interpreter {
  def main(args: Array[String]) = {
    if (args.length != 1) {
      println("parameters: filename to interpret")
    } else {
      interpretFile(new File(args(0)))
    }
  }
  val runtime = new ScratRuntime
  def interpretFile(file: File) {
    if (file.canRead) {
      val source = io.Source.fromFile(file)
      source.getLines().foreach(runtime.eval)
      source.close()
    } else {
      println("cannot open file " + file.getPath)
    }
  }
}
```
[![](http://cdn.memegenerator.net/instances/400x/26006910.jpg)](http://cdn.memegenerator.net/instances/400x/26006910.jpg)\
And now I can put my hello world in a file and run it. \
But I needed to decide on the extension. I know it's silly but I didn't
want to save the file until I had the extension in mind. And this mean
naming the language. Being in "logic mode" I asked my awsome girlfriend
who's more artsy type of a person and she immediately responded
"scrat"(she's huge fan of
[Scrat](http://en.wikipedia.org/wiki/List_of_Ice_Age_characters "List of Ice Age characters")
the squirrel from [Ice
Age](http://en.wikipedia.org/wiki/Ice_age "Ice age")). And them some
more funny names, but scrat stuck with me. So I named the file
hello.scrat.

**next: [variables and decisions](/posts/2012-09-02-creating-a-language-5.html)**

