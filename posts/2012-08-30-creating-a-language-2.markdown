---
title: Making a programming language Part 2 - something that kinda works
---

[Table of contents](/posts/2012-08-29-creating-a-language-1.html),
[Whole project on github](https://github.com/edofic/scrat-lang),
[relevant version on github](https://github.com/edofic/scrat-lang/tree/blogpost1and2)

In the [Part 1](/posts/2012-08-29-creating-a-language-1.html) I posted a working
repl([read-eval-print-loop](http://en.wikipedia.org/wiki/Read–eval–print_loop "Read–eval–print loop"))
for simple math expressions but I kinda cheated and only explained how I
built the [AST](http://en.wikipedia.org/wiki/Abstract_syntax_tree "Abstract syntax tree").

### AST elements

Just [scala](http://www.scala-lang.org/ "Scala (programming language)") case classes
```scala
sealed trait Expression
case class Number(n: Double) extends Expression
case class Add(left: Expression, right: Expression) extends Expression
case class Subtract(left: Expression, right: Expression) extends Expression
case class Multiply(left: Expression, right: Expression) extends Expression
case class Divide(left: Expression, right: Expression) extends Expression
```

### [Parser combinators](http://en.wikipedia.org/wiki/Parser_combinator "Parser combinator") revisited

I use power of scala library to cheat a bit and do lexing and
[parsing](http://en.wikipedia.org/wiki/Parsing "Parsing") in one
step.Basic parser combinators from scala api documentation, everything
you need to define productions in your grammar.
```scala
p1 ~ p2 // sequencing: must match p1 followed by p2
p1 | p2 // alternation: must match either p1 or p2, with preference given to p1
p1.?    // optionality: may match p1 or not
p1.*    // repetition: matches any number of repetitions of p1
```
However, to transform the matched string to an AST you need something
more
```scala
private def number: Parser[Expression] = """\d+\.?\d*""".r ^^ {
  s => Number(s.toDouble)
}
```

Firstly, in the `RegexParser` class is an [implicit conversion](http://en.wikipedia.org/wiki/Type_conversion "Type conversion")
from `Regex` to Parser. So I could write
```scala
private def number: Parser[String] = """\d+\.?\d*""".r
```
Notice the [type annotation](http://en.wikipedia.org/wiki/Type_signature "Type signature").
Inferred type would be Regex, since this function is private I can still
have implicit conversion, but I rather have all parsers be of type
`Parser[_]`.The `^^` part is an action combinator - a [map
function](http://en.wikipedia.org/wiki/Map_%28higher-order_function%29 "Map (higher-order function)").
But as `^^` is only available on Parser instances my regex has already
been implicitly converted. So in my lambda I already know(scala can
infer) the type of s to be String. One last example
```scala
private def term: Parser[Expression] = factor ~ rep(("*" | "/") ~ factor) ^^ {
  case head ~ tail =>
    var tree: Expression = head
    tail.foreach {
        case "*" ~ e => tree = Multiply(tree, e)
        case "/" ~ e => tree = Divide(tree, e)
    }
    tree
}
```
Function rep is also from Parsers class matches any number of
repetitions(including 0). Here's the type signature
```scala
def rep[T](p: ⇒ Parser[T]): Parser[List[T]]
```

The catch here is that `~` returns a single parser that matches both
sides, but fortunately it can be pattern matched to extract both sides.
And I can even use meaningful names since I am in fact matching a head
with an optional tail. Inside the [case
statement](http://en.wikipedia.org/wiki/Switch_statement "Switch statement")
I used more imperative style to build a tree, nothing fancy here.
Folding was a bit awkward in this case for me(maybe I'm just
incompetent) so I went with a for each loop.Apply the same pattern to
+/- part and you have yourself a tree.Oh, yeah...and the top parser
function. A bit changed from last time, to yield useful error messages
```scala
def apply(s: String): List[Expression] = parseAll(expr, s) match {
  case Success(tree, _) => Right(tree)
  case NoSuccess(msg, _) => Left(msg)
}
```

### Evaluator

Now what I promised in previous post - evaluation. At first I planned on
compiling the code for JVM but I just wanted to see some results first
so I decided to do a simple interpreter, no compiling whatsoever - for
now. My first approach was to modify the Expression
```scala
sealed trait Expression{
  def eval(): Double
}
```
and implement this on all case classes hardcoding the double type and
coupling AST representation and evaluation together. Yikes. Granted, it
worked, but what an ugly way to do it. So I did a hard reset(always use
git folks! or something similar) and went about doing a standalone
evaluator. Since scala's pattern matching skills are awesome and I'm
already using case classes why not just do that.
```scala
object Evaluator {
  import Tokens._
  def apply(e: Expression): Double = e match {
    case Number(n) => n
    case Add(l, r) => apply(l) + apply(r)
    case Subtract(l, r) => apply(l) - apply(r)
    case Multiply(l, r) => apply(l) * apply(r)
    case Divide(l, r) => apply(l) / apply(r)
  }
}
```
This is all the code. Just pattern matching and recursion. But yes,
still hardcoded double as the data-type. Looking back, not a great
decision...but hey I got something working and this is fun.

**next time: [Adding features(contants, exponents, function calls)](/posts/2012-08-31-creating-a-language3.html)**
