---
title: Cool Monday - Hindley-Milner on a dynamic language
date: 2012-11-05
---

So I'm getting into type theory. Slowly. Note to self: read a proper
book on this topic. I'm getting familiar with it through some practical
applications. Namely scala and haskell. 

That same [discussion about design
patterns](/posts/2012-10-29-design-patterns-bullshit.html) also
included dynamic vs [static
typing](http://en.wikipedia.org/wiki/Type_system "Type system"). And I
asked twitter about it. [HairyFotr](https://twitter.com/HairyFotr) linked
this [amazing talk about type
inference](http://screencasts.chariotsolutions.com/uncovering-the-unknown-principles-of-type-inference-)
to me. Basically there are two conclusions to be drawn
-   Every static typed language should have at least limited type
    [inference](http://en.wikipedia.org/wiki/Inference "Inference").
    It's compiler's job to do so and quite trivial to implement.
-   Properly done static typed language provides all features the that
    dynamic typed languages can. Safely.

As I'm (still) [implementing a
language](/posts/2012-08-29-creating-a-language-1.html) that
happens to be dynamic(because I was too lazy to look-up how to do type
checking) second point interests me more. 

Can I turn my language into a static one without changing syntax(adding
type annotations) and losing features? That would be awesome!

After a day of thinking, answer seems to be **YES!**

### Global type inference in a nutshell

So I want a dynamic-like syntax(no types anywhere). Good news is I don't
have nominal types, so I can use inference to get structural types. 

Java, C# and many other mainstream languages use [nominal
typing](http://en.wikipedia.org/wiki/Nominative_type_system "Nominative type system")
at the level of the vm. This means that type A is a subtype of type B
precisely when name of A is a subtype of name ob B. For example 
```java
interface Foo{
    void method(int a);
}

interface Bar{
    void method(int b);
}

If you have a method that takes in an instance of Foo, you cannot pass
an instance of Bar. Because Bar isn't subtype of Foo. Even though they
are [structurally](http://en.wikipedia.org/wiki/Structure "Structure")
the same. Fun fact: because
[JVM](http://en.wikipedia.org/wiki/Java_Virtual_Machine "Java Virtual Machine")
is designed like this, scala cannot have global type inference.
See, global inference(or "[type
reconstruction](http://en.wikipedia.org/wiki/Type_inference "Type inference")")
looks at usages and reconstructs properties and structure.  Another
example

    a = b.c + 1

That would be legal scrat code(or python, ruby or many other things).
You take b's property named c, add one to it and assign the result to a.
So b must have property c. That's the first structural requirement. From
the + operator you can see that this c must be numeric. And from this
follows that a is also numeric as it's the result of this computation.
If this were a body of a method and b it's parameter, type requirement
for be would be(in made up syntax) `b: { c: Number }` - object with member
of Number type. But that doesn't give you no class names.
So why is this **global**inference? It goes through the whole block of
code(usually a function body) and puts in stub types where it doesn't
have enough info and then solves the system of
[requirements](http://en.wikipedia.org/wiki/Requirement "Requirement")
and substitutes back.

### Possible problems

First of all, I have objects. This means I have to reconstruct object
structure. This shouldn't be too bad. And you can quite easily figure
out that A is subtype of B if set of requirements for A is a superset of
requirements for B. 

Then there's "mutable types". I concluded that following code should be
illegal

    a = 1
    a = "two"

as type of a should remain the same as in first assignment. This lets
you reason about the code much more. But there's a hidden mutability. I
use regular functions as object constructors returning a keyword "this"
that evaluates to current scope. But throughout the body you still
access this scope and it's type(it's an object after all) changes with
every new (first) assignment and function definition. But this should be
tracked through all possible code paths. Only problem is an if
expression. Type of an if(and it's side-effects) can only be common
super type of both then and else branch - an intersection of requirement
sets. In a dynamic language you can reason about conditions and conclude
when something should definitely be in scope and use it. Automatic
reasoning about conditions? This could turn out tricky. Perhaps in later
implementation.
An there's third an final problem(that I can see). Infinite types. I
have not seen a practical usage but it doesn't work and that bothers me.
Dynamic code in scrat just works, but same code translated to haskell
yields a compiler error - "can't instantiate infinite type". But
apparently infinite types can be detected, so maybe I can find a way to
present them and it will compile.

### Conclusion

There are some problems but I believe I can make it work. It would be
super awesome to have a language that feels dynamic but gives you all
benefits of static typing. Compilers should do the hard work after all!
And they should be capable of inferring general enough types that all
correct programs type check.
Or am I missing some important aspect that works only with dynamic
types?
