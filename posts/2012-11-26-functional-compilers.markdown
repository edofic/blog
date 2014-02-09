---
title: Cool Monday: Functional compilers and atoms 
--- 

I've seen
[this great
talk](http://skillsmatter.com/podcast/scala/functional-compilers-from-cfg-to-exe/ac-5896)
by Daniel Spiewak on Functional
[Compilers](http://en.wikipedia.org/wiki/Compiler "Compiler"). He talks
about lexical and [semantic
analysis](http://en.wikipedia.org/wiki/Compiler "Compiler")
in particular.
First, problems with traditional
[lexing](http://en.wikipedia.org/wiki/Lexical_analysis "Lexical analysis")
with scanner. You can only have regular tokens or you have do do some
dirty hacking and
[backtrace](http://en.wikipedia.org/wiki/Stack_trace "Stack trace") the
scanner therefore losing linearity. And you can solve this with
[scannerless
parsing](http://en.wikipedia.org/wiki/Scannerless_parsing "Scannerless parsing")
- putting [regular
expressions](http://en.wikipedia.org/wiki/Regular_expression "Regular expression")
into your grammar. In fact this approach seems simpler to me, as the
[only proper parser I've
done](/posts/2012-08-29-creating-a-language-1.html)
works this way. But this is not the interesting part.


### Semantic analysis

This is where fun kicks in. After you parse the [source
code](http://en.wikipedia.org/wiki/Source_code "Source code") into an
[AST](http://en.wikipedia.org/wiki/Abstract_syntax_tree "Abstract syntax tree")
you need to do a bunch of operations on it. Naming, typing(even
optimization in later phases). If I want to stay functional(which
usually I do) my instinct tells me to do recursive traversal and rewrite
the tree. And that's exactly what my language does. But there is one
huge problem. AST is not a tree. It's huge misnomer. AST is just a
[spanning
tree](http://en.wikipedia.org/wiki/Spanning_tree "Spanning tree") in the
program graph. See, when you add stuff like let expressions, or
types(what I'm doing currently) you get problems
```haskell
let a = 1in f a
```
There are edges between the siblings. Or going back up. Or skipping
levels. Definitely not trees. These edges may be implicit but you still
have to store the information. Traditional solution to this is to
compute look-up tables(maps) and carry them along with the tree. So the
AST remains a tree but it has some additional stuff that implicitly
makes it into a graph. Problem is this gets nasty when you carry along a
lot of information and you have to be careful with you updates.
There is one more solution. Vars. Works like a charm. Except that it's
terrible to reason about quite the opposite of functional. But there
exists a fix.

### Atoms

Think write-once vars. But not quite. The idea is to have containers
that can be written to but are only ever seen in a one state. Problem
with vars is that they can be seen  in multiple states and you have to
keep track of these states. Vals solve this by not letting you mutate
state. And lazy vals provide machinery to delay initialization(great for
solving circular dependencies). But they don't let you escape the scope.
Or deal with a situation when you need to init them when you have data
not when you need to read them. And this is the problem in a compiler.
You compute data coming from out of the scope and you need to store it.
And some time later you need to read it. And you use atoms.  First some
code, then explanation.
```scala
class Atom[A] {
  private var value: A = _
  private var isSet = false
  private var isForced = false
  protected def populate(){
    sys.error("cannont self-populate atom")
  }

  def update(a: A) {
    if (!isSet || !isForced){
      value = a      
      isSet = true
    }
  }  

  def apply(): A {
    isForced = true
    if (isSet) {
      value
    } else {
      populate()
      if (!isSet) {
        sys.error("value not set")
      }
      value
    }
  }
}
```

Here is the workflow... you create an atom, you can write(update) to
it-in fact writes are indempotent and you can do many successive writes
as long as you don't read the value. Once written to isSet flag is set
and atom can be read(apply method) setting the isForced flag. If the
atom isn't set when you try to read it it will try to populate itself.
Populate method is intended to be overwritten and may contain data in
it's closure or even perform some side-effects. And you can safely
assume it will only execute once. And if everything fails and atom isn't
set you get an error. Yay no bothering with nulls any more.
You can quickly see how atoms are great containers for storing computed
information in the AST for passing it on to the later stages.
