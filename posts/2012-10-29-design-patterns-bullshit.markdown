---
title: Design patterns are bullshit!
--- 

It all started Monday
morning at college. Professor of [Information
Systems](http://en.wikipedia.org/wiki/Information_systems "Information systems")
was talking about evolution of programming techniques and at the end of
the chain was
[OOP](http://en.wikipedia.org/wiki/Object-oriented_programming "Object-oriented programming").
The pinnacle of program design. Well, no. Even professor admitted that
OOP failed to deliver. (No word on
[FP](http://en.wikipedia.org/wiki/Functional_programming "Functional programming")
though, I was kinda disappointing). This made me think about problems
of[Java](http://www.oracle.com/technetwork/java/ "Java (programming language)")
and the like.
Some hours later I'm sitting in cafe Metropol above
[Kiberpipa](http://maps.google.com/maps?ll=46.056184,14.503798&spn=0.005,0.005&q=46.056184,14.503798%20(Kiberpipa)&t=h "Kiberpipa")
having tea with some friends - freshmen from FRI. And one of them asks
me if there is a class on [design
patterns](http://en.wikipedia.org/wiki/Design_pattern_%28computer_science%29 "Design pattern (computer science)").
I give a puzzled look and say no and he goes on to explain he's
currently reading [Design
Patterns](http://en.wikipedia.org/wiki/Design_pattern_%28computer_science%29 "Design pattern (computer science)")
in C#. Apparently an awesome book that teaches you some general
approaches to problems.
But I have quite strong opinions on some topics. And one of them are
design patterns. I believe they are bullshit and I said that. I didn't
read the GoF book and I don't have any intention to do so in near
future. But I'm familiar with some patterns. Mostly workarounds for
shortcomings of OOP languages and contrived solutions to non-existent
problems. And I explained this.
Response was to paraphrase: "Well, MVC is a great design pattern". But I
argued that it isn't a design pattern at all. Since everybody had a
laptop, a quick wiki search cleared this up. It's an [architectural
pattern](http://en.wikipedia.org/wiki/Architectural_pattern "Architectural pattern").
Not a design pattern.

### What are design patterns

So I asked him out of curiosity for some patterns from the book(also
trying to prove my point).

#### Singleton

Well I used this one, so I cannot say it's bad. But I can still say it's
overused and when it lets you access [mutable
object](http://en.wikipedia.org/wiki/Immutable_object "Immutable object")
graph. Because this leaks mutable shared data all across your code base.
But it gets real messy with double locking and stuff like that. Luckily
Java solves this with enum. Now design pattern becomes just a keyword.
That doesn't qualify as a pattern in my book. C# has static classes
that are basically same thing, singleton instance held by the [class
loader](http://en.wikipedia.org/wiki/Java_Classloader "Java Classloader").
For some reason C# folks call this Monostate, but it's the **same
concept**. Concepts are important, because they're language agnostic.
Even though patterns try, they inherently aren't. And just for the last
nail in the coffin.
[Scala](http://www.scala-lang.org/ "Scala (programming language)") has
`object` keyword that creates a singleton. 

#### Factory, static factory...

This one is also abused. You even have metafactories that produce
factories. Srsly...wtf? That's exactly like currying but wrapped up in
objects and contrived to the point you no longer see it's currying. The
only useful use-case for them is abstraction over constructors. That is,
a single method that, depending on the arguments calls some constructor
for the type it constructs. But again, this is hardly a pattern. 

#### Decorator

Decorator in statically typed languages means just implementing an
interface or mixing in a trait(scala). I don't see the point here.
That's why interfaces(traits) exist. 

#### Proxy

My friend insisted that this is something like [Strategy
pattern](http://en.wikipedia.org/wiki/Strategy_pattern "Strategy pattern")(didn't
use the word) but from the other way. Proxy is the user of strategy
object. Anyway strategy is just a wrapper for function value, nothing to
see here. But not I googled it and I'm even more puzzled. It seems this
means object reuse. Like literally sharing data. You know, like you do
with any immutable values. Did I misread?

And about here we digressed to static vs [dynamic
typing](http://en.wikipedia.org/wiki/Type_system "Type system"). Sadly I
was the only one for static. Apparently nobody wants compiler to help
them. Now I'm sad. Compilers are one of the best programs out there. And
with type inference they do all the work for you. You must be a fool not
to use one and rather do a shitload of unit tests instead.

So... and anyone point to a useful design pattern or am I right and they
truly are bullshit?
