---
title: Making a programming language: Part 7a - objects
---

[Table of contents](/posts/2012-08-29-creating-a-language-1.html), 
[Whole project on github](https://github.com/edofic/scrat-lang)

My goal in this post is for this to compile

    func create(n){this}

and a call to it to return a reference to an
[object](http://en.wikipedia.org/wiki/Object_%28computer_science%29 "Object (computer science)")
that contains "n". Functions are a bit different from the rest of the
language. Not by
[implementation](http://en.wikipedia.org/wiki/Implementation "Implementation")
or usage but by [thought
process](http://en.wikipedia.org/wiki/Thought "Thought") behind
designing them. I actually thought about objects and
[functions](http://en.wikipedia.org/wiki/Function_%28mathematics%29 "Function (mathematics)")
before implementing any of them. Considering my implementation of
scopes(which I like) and shadowing I got this great idea that functions,
scopes and objects are just many faces of the same thing. Or "could be"
many faces of the same thing. Something a bit more powerful than a
function being the "thing". Let's see how that works. A function in this
language needs a new [local
scope](http://en.wikipedia.org/wiki/Local_variable "Local variable") for
every execution(not all functions, but this is a simplification because
I don't care about performance, see [previous
post](/posts/2012-09-25-creating-a-language-6.html)
for details). New scope. New something. New. Bells should be ringing
right now. I'm creating objects. I could just as well pass the reference
to that object. I even have the reference. It's current scope - "this"
in scala code. So I just need a [language
construct](http://en.wikipedia.org/wiki/Language_construct "Language construct")
to access that. How about "this".
```scala
//in SScope body
map.put("this", this)
```

Not even a language construct, just a magic variable. You can even
shadow it. Anticlimactic? I hope so. The whole point of object is that
it's just another view on function. I could end the post here but I want
to show why this is awesome(and therefore why I'm proud of it)

### Privates

You don't even need access modifiers, you can use shadowing and nesting

    func private(n) {
        that = this
        func (){ 
            func get() { that.n }
            func set(v) { that.n = v }
            this
        }()
    }

This is a constructor that returns an object with functions get and
set(but no n!). Outer object is available through its alias via nesting.
And returned value is the last expression - an immediately invoked
lambda(cumbersome syntax..gotta do something about that). However, this
is not bullet-proof

    o = private(1)
    o.that = func(){
               n=3
               this
             }
    o.get()

Returns 3. Though knowledge of implementation is needed to execute such
attack. **update:** I kinda sorta forgot to include how I made that
dot-access thingy-o.get() to work. Continuation below 

**next: [using objects](/posts/2012-10-08-creating-a-language-7b.html)**
