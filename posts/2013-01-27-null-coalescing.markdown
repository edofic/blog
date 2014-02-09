---
title: Null-coalescing(??) in scala
--- 
 
 -----------------
 ![The supermassive black holes](http://upload.wikimedia.org/wikipedia/commons/thumb/d/d4/BlackHole.jpg/300px-BlackHole.jpg)
 Null reminds me of black holes.
 -----------------

I was doing my homework
today(yes I am aware I should be enjoying myself on 30th December) and had some problems
with [concatenating](http://en.wikipedia.org/wiki/Concatenation "Concatenation") possibly
null strings in
[LINQ](http://en.wikipedia.org/wiki/Language_Integrated_Query "Language Integrated Query").
Quick trip to
[StackOverflow](http://stackoverflow.com/ "Stack Overflow") and I find
out
[C#](http://msdn2.microsoft.com/en-us/vcsharp/aa336809.aspx "C Sharp (programming language)")
has some funky operators that solve this in a (sort-of) clean way.
```csharp
var outputString = input1 ?? "" + input2 ?? "";
```

I like [type
inference](http://en.wikipedia.org/wiki/Type_inference "Type inference")
so I use var's extensively - please don't judge me. What this does is
concatenate input1 and input2 substituting null values with [empty
string](http://en.wikipedia.org/wiki/Empty_string "Empty string"). In
scala you would write something like

```scala
val outputString = Option(input1).getOrElse("") + Option(input2).getOrElse("")
```

but that's verbose and ugly. Wrapping and unwrapping to get some added
semantics of the Option factory(returns None for null). But you
shouldn't have nulls in the first place if you're using scala. That's
what Option is for! Enough ranting, let's implement this [??
operator](http://en.wikipedia.org/wiki/Null_coalescing_operator "Null coalescing operator")
in scala.

```scala
implicit class NullCoalescent[A](a: A){
  def ??(b: => A) = if(a!=null) a else b
}
```

Yup. That's it! You only need this in scope and of you go writing code
C#-style.

```scala
scala> "hi" ?? "there"
res0: String = hi

scala> (null:String) ?? "empty"
res2: String = empty
```

Note the type of b parameter in ?? method. It's => A. By name
evaluation. This means our new operator behaves properly and only
evaluates right hand side if value is in fact null. This lets you for
example log unexpected nulls while substituting for default value

```scala
val key = valueFromUser ?? {
  log("key is null!"}
  defaultKey
}
```

This works because scala let's you do side effects like that. I just
wanted to add one reason why I love scala is this easy way of defining
structures that feel like they're part of the language while being
nothing more than just libraries.
