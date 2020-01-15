---
title: Virtual machine in C++
date: 2013-10-18
---

This is not a tutorial. This post
is a flashback I had today. It might be a bit fiction as my memory about
events tends to be fuzzy at times. But I promise it at least resembles
the real story.
I was in [elementary
school](http://en.wikipedia.org/wiki/Elementary_school "Elementary school")
and just found out about programming and was learning about c++. After
reading "[C++](http://en.wikipedia.org/wiki/C%2B%2B "C++") na kolenih"
by Goran Bervar I was empowered by knowledge and tried to do all sorts
of projects. Mostly in console. Stuff like subtitle format converter -
[NIH
syndrome](http://en.wikipedia.org/wiki/Not_Invented_Here "Not Invented Here").
I was a bit frustrated because I couldn't find any books about windows
programming in the library. Yes, library was may primary source of
information, because
my [English](http://en.wikipedia.org/wiki/English_language "English language") was
not nearly good enough for technical stuff.
I might add here I worked on
[Windows](http://www.microsoft.com/WINDOWS "Windows") 98(and later XP)
with [DevC++](http://www.bloodshed.net%2c/ "Dev-C++"). I found out about
[Visual
Studio](http://www.microsoft.com/visualstudio/en-us "Microsoft Visual Studio")
in a few years and did some Windows development.
I digressed a bit. Then came the most optimistic idea. A [virtual
machine](http://www.symantec.com/theme.jsp?themeid=protect-virtual-environments "Virtual machine").
Something quite high level(instruction to print) an eventually an
assembler. I now realize I was always [into language
stuff](/posts/2012-08-29-creating-a-language-1).
So a designed a [machine
language](http://en.wikipedia.org/wiki/Machine_code "Machine code") with
just enough instructions to do [Hello
World](http://en.wikipedia.org/wiki/Hello_world_program "Hello world program"),
that is PRINT and END.

### Implementation

At first I thought about doing a monolithic structure - [switch
case](http://en.wikipedia.org/wiki/Switch_statement "Switch statement")(in
fact what I've [done with scrat
recently](/posts/2012-08-29-creating-a-language-1)).
But I had some considerations. What if number of of instruction rises a
lot? I'll be left maintaining [spaghetti
code](http://en.wikipedia.org/wiki/Spaghetti_code "Spaghetti code"). Or
at least I thought that's what spaghetti code looks like, but in
retrospective I believe I had a good taste anyway. 

But I tried that anyway. Just for kicks. Did whole machine as one class
that had an array for memory and a single point of entry - boot. It run
a loop a while loop with
PC<-PC+1,
fetched instruction from memory, switched on them, called appropriate
method to implement that instruction and looped. Even had registers. I
think my current professor of Computer Architecture(this course brought
back the memory) might actually be proud if he heard what I did back
then. 

### Pointers

I was always quite comfortable with pointers. I don't now, they're
mathematicky concept. I like such stuff. Or perhaps it was because I was
young when I was introduced into the matter and wasn't spoiled with
automatic memory management(which I quite like nowadays). 

So I tried with [function
pointers](http://en.wikipedia.org/wiki/Function_pointer "Function pointer").
C is cool enough to let you have pointers to functions! And that means
[higher order
functions](http://en.wikipedia.org/wiki/Higher-order_function "Higher-order function").
But I didn't know about math enough to appreciate the concept as I do
now. But still - I thought it's extremely cool. So I did a function that
halted execution and printed out "no such instruction". Why you ask?
Well I did a 256-cell table(8-bit instruction) of pointers to functions.
Now I didn't have to switch - just a look-up and invocation. Great.
Apart from the fact it doesn't work. 

Compiler said something along the lines of "You cannot make a table of
pointers to functions!". I was puzzled. Skip 10 years into the future.
Today I was rethinking this and thought about casting the pointer. All
the functions would be void->void so I can cast back no problem. A
table of [void
pointers](http://en.wikipedia.org/wiki/Pointer_%28computing%29 "Pointer (computing)")
and casting. Yay!

Now 10 years back. I didn't think about casting the pointer. Type info
was sacred to me!

So I "invented" [function
objects](http://en.wikipedia.org/wiki/Function_object "Function object").

### Objects

I swear to god I have not heard about function objects back then. It
wasn't until this year reading Bloch's Efficient Java where he talks
about strategy objects. I immediately recognized my idea. So now I had
many classes, every one implementing execute method. And I had an array
of these objects. Now I did a look-up and invocation on an object.
Sweet. And it even worked. But sadly I dropped the project and went on
to graphics. Learnt SDL and did Tic-Tac-Toe. And dreamed about doing a
vector 3D engine(curves baby!). Which until this day I didn't try to
implement. Maybe I'll try in near future. 
