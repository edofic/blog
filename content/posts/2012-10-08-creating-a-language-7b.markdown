---
title: Making a programming language Part 7b - using objects
date: 2012-10-08
---

[Table of contents](/posts/2012-08-29-creating-a-language-1), 
[Whole project on github](https://github.com/edofic/scrat-lang)

Something like EPIC FAIL occured to me and I [published a post](/posts/2012-08-29-creating-a-language-1)
containing only half the content I intended to write. So I'm doing a
part b.

My intended usage of objects is something along the lines of

    objectName.someProperty
    objectName.someFunction()
    someFunction().someProperty
    someObject.someProperty.someFunction().someProperty.someFunction

Explanation

1.  getting a value from an
    [object](http://en.wikipedia.org/wiki/Object_%28computer_science%29 "Object (computer science)")
2.  invoking a function contained in an object
3.  getting a value from returned object of the invoked function
4.  a bit contrived example. Invoking a function contained inside a
    property(object) of an object and then getting a [function
    value](http://en.wikipedia.org/wiki/Function_%28mathematics%29 "Function (mathematics)")
    from a property of the returned value from the first function.
    That's a mouthful, just read the damn code instead

### Dot access

So everything bases on those little dots. First my thoughts were
something like "you just do expr <- expr | expr.expr". This is just
wrong. At least I should have reversed the order as this leads to
infinite [left
recursion](http://en.wikipedia.org/wiki/Left_recursion "Left recursion").
Then I might have got away. Then I realized I only need dots after
function calls and simple identifiers. Design choice(if you think it's a
bad one leave a comment). Notice the "simple
[identifier](http://en.wikipedia.org/wiki/Identifier "Identifier")".
That's what I did: Renamed identifier to simple identifier and put
something that handles dots under name identifier. And then fixed
everything. 
```scala
case class DotAccess(lst: List[Expression]) extends Expression

private def identifier: Parser[DotAccess] =
    rep1sep((functionCall | simpleIdentifier), ".") ^^ DotAccess.apply
```

That's about it. At least for parsing. Now the fun begins.

### Nesting scopes

Scopes were designed with nesting in mind. This is a double edged sword.
See, the "privates" can be  done if you rely on not being able to access
the parent scope. If dot access exposes full addressing functionality a
powerful feature ceases to exist. So some protection should be in place.
Something like strict get
```scala
class SScope
  ...
  def getStrict(key: String): Option[Any] = map.get(key)
  ...
```

And I also added an unlinked view to it just to ease usage. This is just
a method that returns new SScope with no parent overriding getters and
put to use map available in closure.
So now I can walk down the list in DotAccess recursively and explicitly
override the implicit scope parameter. And everything automagically
works. Well, not quite. If you have a function call, the arguments need
to be evaluated in top scope. Not in the nested one like the function
identifier. At first I didn't even think about this and only failing
attempts at more complex recursion brought up this quite obvious bug.\
So how to solve this? I could pre-evaluate all arguments, but I use
recursion to do this and it's two levels(at least) deeper from where
dots happen. So no go. I need to carry on the outer scope. I overloaded
the apply method from Evaluator so other code can still function(tests
ftw!) and all in all it looks like this:
```scala
def apply(e: List[Expression])(implicit scope: SScope): Any = {
  (e map apply).lastOption match {
    case Some(a) => a
    case None => ()
  }
}

def apply(e: Expression)(implicit scope: SScope): Any = apply(e, None)(scope)

def apply(e: Expression, auxScope: Option[SScope])
         (implicit scope: SScope): Any = e match {
  ...
  case DotAccess(list) =>
    val outerScope = scope
    def step(list: List[Expression])
            (implicit scope: SScope): Any = list match {
      case Nil =>
        throw new ScratInvalidTokenError("got empty list in DotAccess")
      case elem :: Nil => apply(elem, Some(scope))(outerScope)
      case head :: tail => apply(head) match {
        case s: SScope => step(tail)(s.unlinked)
        case other =>
          throw new ScratInvalidTypeError("expected scope, got " + other)
      }
    }
    step(list)
}
```
So an optional aux scope is the answer. It doesn't seem pretty to me,
but it does the job. 


**next: **[trying to go faster](/posts/2012-09-29-creating-a-language-8)
