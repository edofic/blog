---
title: Making a programming language: Part 1 - how to start
---

### Table of contents

-   Part 1 - how to start (this article)
-   [Part 2 - something that kinda
    works](/posts/2012-08-30-creating-a-language-2.html)
-   [Part 3 - adding
    features](/posts/2012-08-31-creating-a-language3.html)
-   [Part 4 - Hello World](/posts/2012-09-01-creating-a-language4.html)
-   [Part 5 - variables and
    decisions](/posts/2012-09-02-creating-a-language-5.html)
-   [Part 6 - user defined
    functions](/posts/2012-09-25-creating-a-language-6.html)
-   [Part 7a - constructors and
    objects](/posts/2012-09-27-creating-a-language-7a.html)
-   [Part 7b - using
    objects](/posts/2012-10-08-creating-a-language-7b.html)
-   [Part 8 - going
    faster](/posts/2012-09-29-creating-a-language-8.html)

#### [Source version for this post](https://github.com/edofic/scrat-lang/tree/blogpost1and2)

Lately I gained some interest in [programming
languages](http://en.wikipedia.org/wiki/Programming_language "Programming language")
and [compilers](http://en.wikipedia.org/wiki/Compiler "Compiler"). Those
seem like quite some daunting monsters - just consider the amount of
different features and more importantly, the vast infinity of possible
programs.So where to start? I have a course about compilers at my
college, but I have to wait another year to enroll into that. So welcome
[Coursera](https://www.coursera.org/course/compilers). It's currently
available only in self-study but that's enough for me. I should mention
the [Dragon
Book](http://www.amazon.com/Compilers-Principles-Techniques-Alfred-Aho/dp/0201100886%3FSubscriptionId%3D0G81C5DAZ03ZR9WH9X82%26tag%3Dzem-20%26linkCode%3Dxm2%26camp%3D2025%26creative%3D165953%26creativeASIN%3D0201100886 "Compilers: Principles, Techniques, and Tools"),
but I didn't read that(yet) so I can't comment. Compiler is basicaly
`lexer -> parser -> optimiser -> code generator`.

### The regex approach

I made it through introduction and lesson on [lexical
analysis](http://en.wikipedia.org/wiki/Lexical_analysis "Lexical analysis")
and recognized [finite
automata](http://en.wikipedia.org/wiki/Finite-state_machine "Finite-state machine")
as something from college(thank your professor!) and finnaly understood
the connection from them to regex in java and the like(professor
mentioned [regular
expressions](http://en.wikipedia.org/wiki/Regular_expression "Regular expression")
but no-one figured out what was the connection). 

Feeling empowered by the newly obtained knowledge I set out to make a
lexer for my language.  I had no clear specification in mind since this
is supposed to be a fun and creative project...I intended to just invent
as I go.

My lexer was something like that(won't post the real code since
I'm embarrassed)
-   create a bunch of java regex objects that all match to the start of the string
-   take the string and try to match it against all regexes
-   return an identifier object corresponding to the match
-   remove matched part of the string
-   loop

Yeah, pretty horrible.I kinda abandoned this idea.

### The [recursive descent parser](http://en.wikipedia.org/wiki/Recursive_descent_parser "Recursive descent parser") approach

By now I was into [functional
programming](http://en.wikipedia.org/wiki/Functional_programming "Functional programming")
and [scala
language](http://www.scala-lang.org/ "Scala (programming language)").  I
also watched lesson on recursive descent on coursera. The code was
fugly, bunch of c++ with
[pointer arithmetic](http://en.wikipedia.org/wiki/Pointer_%28computing%29 "Pointer (computing)") and
side effects. I wan't pretty functions :(

I considered doing a framework for this in scala or perhaps java but...

  ----------------------------
  [![Scala (programming language)](http://upload.wikimedia.org/wikipedia/en/thumb/8/85/Scala_logo.png/300px-Scala_logo.png)](http://en.wikipedia.org/wiki/File%3AScala_logo.png)
  Scala (programming language) (Photo credit: [Wikipedia](http://en.wikipedia.org/wiki/File%3AScala_logo.png))
  ----------------------------

Enter scala's [parser
combinators](http://en.wikipedia.org/wiki/Parser_combinator "Parser combinator").
I suggest reading
[this](http://www.codecommit.com/blog/scala/the-magic-behind-parser-combinators)
if you arent familiar. One of the reasons scala is awesome. You get
parser "generation" baked into the standard library. Granted, it's a bit
slow, but who cares - this is a project for fun not for general usage.
On the PLUS side you get productions in the laguage itself and parsing
is almost like magic.

And in scaladoc you even get a nice [code
snippet](http://www.scala-lang.org/api/current/index.html#scala.util.parsing.combinator.RegexParsers).
What more can a geek ask for.

Well....an
[AST](http://en.wikipedia.org/wiki/Abstract_syntax_tree "Abstract syntax tree")
would be nice. Fortunately scala also provides action combinators that
allow to [pattern
match](http://en.wikipedia.org/wiki/Pattern_matching "Pattern matching")
the [parsers](http://en.wikipedia.org/wiki/Parsing "Parsing") themself
and map them into an AST. And it even isn't that complicated. 

I recreated the grammar from the scaladoc example and added the action
combinators. The code block is kinda long, but I promise there is a
short explanation at the bottom.

Notice the `^^` symbol. It sends the output from the parser  combinator
into my mapping function(anonymous). And I just create the case
class.There is also an evaluator at the and...but that's in part 2.

**[next -\> Part 2: something that kinda works](/posts/2012-08-30-creating-a-language-2.html)**
