---
title: Making a programming language Part 6 - functions
date: 2012-09-25
---

  --------------------
  [![Illustration of Function (mathematics).](http://upload.wikimedia.org/wikipedia/commons/thumb/8/8a/Function_illustration.svg/200px-Function_illustration.svg.png)](http://commons.wikipedia.org/wiki/File%3AFunction_illustration.svg)
  Illustration of Function (mathematics). (Photo credit: [Wikipedia](http://commons.wikipedia.org/wiki/File%3AFunction_illustration.svg))
  --------------------

[Table of contents](/posts/2012-08-29-creating-a-language-1.html),
[Whole project on github](https://github.com/edofic/scrat-lang)

Long overdue, I'm finally writing about the most interesting part -
[user defined functions](http://en.wikipedia.org/wiki/User-defined_function "User-defined function").
Objects should be in the next post as they are a natural extension of
what I'm about to do. And because I'm to lazy to write a post that
long.

### What's a function?

A function is something that takes parameters and returns a result. And
I'm opting for side-effects as this is simpler to get something working
that doing the whole [referential
transparency](http://en.wikipedia.org/wiki/Referential_transparency_%28computer_science%29 "Referential transparency (computer science)")
and
[IO](http://en.wikipedia.org/wiki/Io_%28programming_language%29 "Io (programming language)")
[monad(Haskell)](http://en.wikipedia.org/wiki/Monad_%28functional_programming%29 "Monad (functional programming)").
Although it would be interesting to have sort of [dynamically
typed](http://en.wikipedia.org/wiki/Type_system "Type system")
[Haskell](http://haskell.org/ "Haskell (programming language)")
thingy...maybe in the future. 

So -
[side-effect](http://en.wikipedia.org/wiki/Side_effect_%28computer_science%29 "Side effect (computer science)")-ful
functions mean function body is a block of expressions, not just an
expression, if I'm to do useful side-effects. This also means I want
[lexical
scope](http://en.wikipedia.org/wiki/Scope_%28computer_science%29 "Scope (computer science)"). 

#### Body scope

Let's think about that. A function body needs a scope. It's parent scope
is kinda obvious - it's the scope in which the function was defined. I
want closures! It's [local
scope](http://en.wikipedia.org/wiki/Local_variable "Local variable")
gets kinda tricky. I implemented read-only cascading and shadowing on
write. If expressions can also have side-effects. So in different
executions a parent variable may be shadowed or not. This means I cannot
reuse the body scope, as the shadows need to be cleaned out(I'm
considering an implementation that reuses it currently, but that's
another story). As I'm not after performance I can simply create a new
scope for each invocation of the function. 

Parameters can be implemented as regular variables pre-put into the
local scope before executing the function body. 

### Parsing

That was the hardest part. I had quite some problems in my grammar as I
tried to introduce blocks. Mostly ambiguity and infinite recursion. I'll
just post the interesting bits here - see [the full
commit](https://github.com/edofic/scrat-lang/commit/181d513801567cb51e7ebc5637d1a64913290b13) if
you're interested in details.

The function definition boils down to:
```scala
private def exprList: Parser[List[Expression]] = repsep(expr, "\\n+".
private def block: Parser[List[Expression]] =
  """\{\n*""".r ~> exprList <~ """\n*\}""".
private def functionDef: Parser[FunctionDef] =
  "func" ~> identifier ~ ("(" ~> repsep(identifier, ",") <~ ")") ~ block ^^ {
    case id ~ args ~ body => FunctionDef(id, args, body)
  }
```

Oh yes, and since I use newlines as separators now, they aren't
whitespace and I have to handle them explicitly.
```scala
override protected val whiteSpace = """[ \t\x0B\f\r]""".r
```

And later on I added lambdas which have optional identifier - so the
only shcange is opt(identifier) in the parser.

### Evaluation

It's just another node in the AST - a case class.
```scala
case class FunctionDef(name: Identifier, args: List[Identifier],
                       body: List[Expression]) extends Expression
```

Now I needed something to execute the definition and create a [function
value](http://en.wikipedia.org/wiki/Function_%28mathematics%29 "Function (mathematics)")
in the enclosing scope. (this is in Evaluator class)
```scala
def createFunFromAst(arglist: List[Identifier],
      body: List[Expression], scope: SScope): FunctionVarArg =
  (args: Any) => args match {
    case lst: List[Any] => {
      if (lst.length != arglist.length) {
        throw new ScratInvalidTypeError(
          "expected " + arglist.length + " arguments, but got " + lst.length)
      } else {
        val closure = new SScope(Some(scope))
        (arglist zip lst) foreach {
        t => closure.put(t._1.id, t._2)
        }
        apply(body)(closure)        }
    }
    case other =>
      throw new ScratInvalidTypeError("expected list of arguments but got" + other)
  }
```
So, what am I doing here? I'm taking a part of the function
definition(all except name - which is optional in the definition and not
needed here) and the parent scope and then returning "FunctionVarArg"
which is the same as native functions in standard library. This new
function relies heavily on scala's closures(this would not be possible
in this way in java!).  First it checks if got a list of arguments(case
clauses) or it throws an exception. Then it checks the arity. Scrat is
dynamically typed, but not sooo dynamically typed. If everything matches
up it creates a new scope("closure"), and inserts key-value pairs for
arguments(zip+foreach). And then it evaluates it's body -
apply(body)(closure). Mind you, this happens on every execution as
createFunFromAst return a function value that, upon execution, does
this.
Oh yes, there is also a case clause in Evaluator's apply that invokes
createFunFromAst, again trivial.
Such functions are indistinguishable to native functions from scrat's
point of view and are invoked by same syntax and by same code in
Evaluator.

### A sample

First thing I tried to implement(literaly) was fibonaci's sequence

  func fib(n) { if n==0 then 1 else if n==1 then 1 else fib(n-1) + fib(n-2) }
  println("20th fibbonacci number is", fib(20))

Excuse me for the ugly nested if's, but this was neccessary as I have
not implemented < yet. But hey, it works.

#### Sneak peak

At this point I realised an awesome way to implement objects. With
constructors like:

  func create(n){this}

**next [objects and costructors](/posts/2012-09-27-creating-a-language-7a.html)**
