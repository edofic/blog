---
title: Making a programming language Part 3 - adding features
date: 2012-08-31
---

[Table of
contents](/posts/2012-08-29-creating-a-language-1), 
[Whole project on github](https://github.com/edofic/scrat-lang)

So now I have a repl that can evaluate stuff like

    (2+3)*(7/2-1)

Not much of a [programming language](http://en.wikipedia.org/wiki/Programming_language "Programming language")
- more like a calculator, and not even a good one. Lets add some
features!


### Constants

Like pi, e and such. I have to change the grammar to match identifiers too.

Now I have
```scala
private def factor: Parser[Expression] = number | ("(" ~> expr <~ ")")
```

And I change that to
```scala
private def factor: Parser[Expression] = value | parenExpr
private def value: Parser[Expression] = number | identifier
private def identifier: Parser[Identifier] = "[a-zA-Z]\\w*".r ^^ {
  s => Identifier(s)
}
```

and I also added a new token type
```scala
case class Identifier(id: String) extends Expression
```

and enabled this in evaluator\
```scala
case Identifier(name) => {
  StdLib(name) match {
    case Some(v: Double) => v
    case _ => throw new SemanticError(name + " not found")
  }
}
```

StdLib is just a map object. I went the python way - variables(and
constants) are entries in a global dictionary. Just a decision I made
while implementing this. As I said, I don't have a plan, I don't know
what I'm doing and I don't know how stuff is done. How hard can it be?!
([Top Gear](http://www.bbc.co.uk/topgear/ "Top Gear (2002 TV series)")
anyone?)


### Exponentiation

Another math feature. It's more of a test for myself to see if I
understand grammars. Especially how to make a grammar not
[left-recursive](http://en.wikipedia.org/wiki/Left_recursion "Left recursion").
Because apparently such grammars don't work with
[RDP](http://en.wikipedia.org/wiki/Recursive_descent_parser "Recursive descent parser").
I turned out I don't understand grammars.
```scala
private def exponent: Parser[Expression] =
  (value | parenExpr) ~ "^" ~ (value | parenExpr) ^^ {
    case a ~ "^" ~ b => Exponent(a,b)
  }

private def factor: Parser[Expression] = (value ||| exponent) | parenExpr
```

The `|||` operator is an ugly hack. It tries both sides and returns the
longer match. By the time I was writing this I didn't know that order is
importat. If I just wrote `exponent | value` it would have worked, because
expoenent would match a value anyway and then failed on missing `^`.

Token and evaluation(uses `math.pow`) for this are quite trivial.


### Function calls
```scala
case class ExpList(lst: List[Expression]) extends Expression
case class FunctionCall(name: Identifier, args: ExpList) extends Expression
```

Simple: function call is a name and list of expressions to evaluate for
arguments(wrapped because even an expression list is an expression)

Parser:
```scala
  private def arglist: Parser[ExpList] = "(" ~> list <~ ")"
  private def functionCall: Parser[FunctionCall] = identifier ~ arglist ^^ {
    case id ~ args => FunctionCall(id, args)
  }
  private def value: Parser[Expression] = number | (identifier ||| functionCall)
```
Again, I was having trouble - parser just didn't work and resorted to `|||`.
`functionCall` should come before identifier.

Evaluating this is more interesting. I decided to make functions be
values too for obvious reasons -> [higher order
functions](http://en.wikipedia.org/wiki/Higher-order_function "Higher-order function")(I'm
into functional programming, remember?). So function values must be
stored in same "namespace". `StdLib`(the only "namespace") required to
become of type `Map[String,Any]`. I will have to do pattern matching
anyway since this will be dynamic-typed language. (Yes this is a plan, I
think it's easier to implement. [Static
typing](http://en.wikipedia.org/wiki/Type_system) ftw, but
that's next project). And I needed a type for function values to pattern
match on - I went with `Any=>Any` and sending in List(arg0,arg1,...)
doing more pattern matching inside the function. Will be slow but
hey...dynamic typing!

from evaluator
```scala
case FunctionCall(name, args) => {
  StdLib(name.id) match {
    case Some(f: FunctionVarArg) => f.apply(apply(args))
    case None => throw new ScratSemanticError("function " + name + "not found")
  }
}
```

and and example function in StdLib\
```scala
type FunctionVarArg = Any => Any
lazy val ln: FunctionVarArg = {
  case (d: Double) :: Nil => math.log(d)
  case other => throw ScratInvalidTypeError("expected single double but got " + other)
}
```

### Conclusion

As clearly illustrated above, not planning your grammar results in
constant changes in many places. So if you're doing something serious
just make the whole fricking grammar on a whiteboard beforehand.
Seriously. 


Anyway..now I still only have a calculator, but a much more powerful
one. I can write expressions like \

    e^piln(10+2)1+2*3/4^5-log(e)

But that's nearly not enough. I want to be Touring-complete an ideally
to be able to compile/interpret itself.\


**next: [Hello World(strings, printing and interpreter)](/posts/2012-09-01-creating-a-language4)**
