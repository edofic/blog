---
title: Making a programming language: Part 8 - going faster
---

[Table of contents](/posts/2012-08-29-creating-a-language-1.html), 
[Whole project on github](https://github.com/edofic/scrat-lang)

First of all, I wrote some tests for scrat. That was a bit challenging
to get started. How do you test a language? I decided to write a bunch
of programs that exercise and combine different language features. And
then stare into code until I was absolutely sure they are correct.
That's the problem with implementing a new language, nobody can tell you
if your code is correct but your software, but you don't even know if
software is correct.
So I wrote some test programs and
[parsed](http://en.wikipedia.org/wiki/Parsing "Parsing") and evaluated
them in my mind to produce final result. And then I wrote a class using
ScalaTest that generates tests by iterating over the array of these
[tuples](http://en.wikipedia.org/wiki/Tuple "Tuple"). Quite cool, oh and
I also included some descriptions into the tuples. So I get nice output.
I pondered this idea quite some time ago but finally implemented it a
bit after functions and objects.

So now I can do crazy
[refactorings](http://www.techopedia.com/definition/3865/refactoring "Refactoring")
while still maintaining the language - I trust my tests! But first of
all, some benchmarks. I decided to measure execution time of parsing and
evaluation on the [linked-list
consturctor](https://raw.github.com/edofic/scrat-lang/master/run/llist.scrat).
To my surprise, parsing was quite slow. It started on about 800ms and
dropped down to 200ms after a few executions(first time [object
creation](http://en.wikipedia.org/wiki/Object_lifetime "Object lifetime")
and
[JIT](http://en.wikipedia.org/wiki/Just-in-time_compilation "Just-in-time compilation")
I suppose). 200ms for 67 [lines of
code](http://en.wikipedia.org/wiki/Source_lines_of_code "Source lines of code")?
That would mean about 3 seconds on a 1000 line file **IF** complexity is
linear(which I later leart isn't!). And that's 3 seconds on fifth run
or so, first run(which is the only one when doing real stuff) would be
10 seconds+, unacceptable. (evaluation takes a few ms)

### Research

At first I just gave it some thought. Well it's a **recursive** descent
parser, and from what I remember it back traces on failure. So
efficiency has a lot to do with grammar structure. And it won't be
linear because you have to do (usually) more back tracing when dealing
with longer input.

Internets here I come. 

My thoughts were confirmed. I found some guys on forums complaining over
speed and then Mr. Odersky himself commented something like this(from
memory): 

> Well the parsers in scala standard library are more like an example
> how to do recusive descent, functional style. They are perfectly 
> usable for parsing command lines but not for long files. You should 
> use a parser generator for that.  

I was bummed. The reason I didn't use a parser generator was this close
integration that parser as a library could provide me. By the way
implementation of
[Parsers](http://en.wikipedia.org/wiki/Parsing "Parsing") is remarkably
short, ~800 lines but most of them are comments. But it has quite some
problems. A lot of object creation - every time a parsing function is
invoked **many object are created**. Each parser is an object and each
combinator is an object too. This in itself is not a big deal, but no
[memoization](http://en.wikipedia.org/wiki/Memoization "Memoization") is
performed, so it becomes a big deal. Now many objects are created on
each try, so when backtracing you have to
[GC](http://www.techopedia.com/definition/27271/automatic-memory-management-amm "Automatic Memory Management")
all these object and then recreate them. It's clean but it's no wonder
it's slow.
Some time passed by and I accidentally found out about [Packrat
Parsing](http://en.wikipedia.org/wiki/Parsing_expression_grammar "Parsing expression grammar").
[This paper](http://scala-programming-language.1934581.n4.nabble.com/attachment/1956909/0/packrat_parsers.pdf)
provides details but the gist of it is to use [lazy
evaluation](http://en.wikipedia.org/wiki/Lazy_evaluation "Lazy evaluation")
and memoization to reduce object creation and speed things up.


### Conversion and results

Conversion is dead easy. It's fully described in scala api
documentation. Basically you mix in PackratParsers and change 
`def : Parser[] = ...` to `lazy val : PackratParser[] = ...` and that's it. Mixed
in trait provides the necessary implicit conversions, lazy makes sure
the creation is only done once and new implementation of parseAll does
some clever parsing. Oh and you needn't convert all parsers, the paper
says the optimal perfomance is achieved with the right mix of standard
[recursive descent
parsers](http://en.wikipedia.org/wiki/Recursive_descent_parser "Recursive descent parser")
and packrat(packrat does include some overhead on specific grammars and
inputs). But I just converted all and run the benchmark again. And
behold...11ms. On the best run. More like 15 on average. But that's
still more than a whole order of magnitude faster. And it should scale
better. 


**next:** don't know yet. I caught up with my implementation(finally!). I
thinking about making functions stronger by relaxing the rules of
invocation and doing some syntax sugar for lambdas. That and java
interop. Or possibly compiling to bytecode. 
