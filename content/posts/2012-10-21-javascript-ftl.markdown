---
title: Javascript faster than light! (well C actually)
date: 2012-10-21
---

  -----------------------
  [![157/365. Acorn - Oak Nut - The Scrat Problem.](http://farm8.static.flickr.com/7156/6772086623_646ee6ab31_m.jpg)](http://www.flickr.com/photos/42149364@N03/6772086623)
  157/365. Acorn - Oak Nut - The Scrat Problem. (Photo credit: [Anant N S (www.thelensor.tumblr.com)](http://www.flickr.com/photos/42149364@N03/6772086623))
  -----------------------

Disclaimer: I never was a fan of js, but I've come to think it's quite
AWESOME!

Anyway I [invented my own toy language
scrat](/posts/2012-08-29-creating-a-language-1.html)
recently. And I now I want it to go fast and do cool stuff. So I went on
to compile it. Well more appropriate term would be
"[translate](http://en.wikipedia.org/wiki/Translation "Translation")"(as
[zidarsk8](https://twitter.com/zidarsk8) pointed out) since my target is
[JavaScript](http://en.wikipedia.org/wiki/JavaScript "JavaScript"). And
then I use node.js to run it - browser test sometime in the future.
Enough about that, I'll be doing a post when I get everything to run
under js.

My original purpuse for translation was speed as node uses V8 and that's
quite speedy. So I did a quick test. I wrote a simple
recursive [Fibonacci sequence](http://en.wikipedia.org/wiki/Fibonacci_number "Fibonacci number")
generator. The cool thing about this is that it takes `fib(n)` steps to
calculate `fib(n)` but call-stack depth is just n - I don't have loops and
I haven't implemented [tail call optimization](http://en.wikipedia.org/wiki/Tail_call "Tail call") yet.
And then I wrote same thing in js an noticed it's quite a bit faster.
Great. Now halfway through
[implementation](http://en.wikipedia.org/wiki/Implementation "Implementation")(more
like 80%) I decided to do a real benchmark.

### Scrat

Here's the source

    func fib(n) if n<2 then 1 else fib(n-1) + fib(n-2)
    println(fib(30))

Neat huh?

And then I timed this repeatedly and all results were about the same:

    time scrat fib.scrat
    1346269.0
    real 0m3.254s
    user 0m3.652s
    sys 0m0.096s

Of course there is some startup overhead that must be taken into account
so I ran an empty file

    time scrat empty
    real 0m0.420s
    user 0m0.448s
    sys 0m0.032s

To obtain total [running
time](http://en.wikipedia.org/wiki/Time_complexity "Time complexity") of
3254 - 420 = 2830ms

### Javascript

Then I translated my source into js. Below is the untouched(apart from
whitespace) result
```javascript
function fib(n){
  return (function(){
    if(n<2.0){
      return 1.0;
    } else {
      return (fib((n)-(1.0)))+(fib((n)-(2.0)));
    }
  }());
}
```

In scrat ifs are expressions too, so the if is wrapped in an [anonymous
function](http://en.wikipedia.org/wiki/Anonymous_function "Anonymous function").
In spite of additional invocations, running time decreased dramaticaly:
128ms.

Real reason for this test was my wory of if overhead so I did a by-hand
implementation
```javascript
function fib(n){
  if(n<2){
    return 1;
  } else {
    return fib(n-1)+fib(n-2);
  }
}
```

Running time: 35ms.
Auch! Wrapping the if statement into an if expression multiplies running
time by almost 4!! But it's still 22 times faster than my interpreter.
(My code sucks I guess)

### C

At this point you should be wondering what does this has to do with C.
Not much. I tried to do an implementation in C just for kicks. To see
how much overhead my by-head function still has. I was assuming C
program will go in something like 10ms.

My best attempt(in the same style: recursion, if expression)
```c
#include

int fib(int n){
  return (n<2)?1:fib(n-1)+fib(n-2);
}

int main(){
  printf("%d", fib(39));
  return 0;
}
```
Startup time is neglectible here, since it doesn't load an interpreter
or a framework. So here's the full running time..ready?
**634 fricking miliseconds!**
That's only 4 times faster than my interpreted code. And 18 times slower
than javascript. I'm not sure how is this even possible. It's probably
just my bad implementation. But rules were: keep the style.
So I hereby declare: js is faster than C. (in this microbenchmark)

UPDATE:
-------

I did something terribly wrong. Look at the C code closely. Its `fib(39)`
where in scrat and js I called `fib(30)`. I just compared apples and
oranges. 

Fixing the C code I got average 20ms. A bit faster than node. So it
turns out javascript isn't faster than light(c) but it's pretty damn
close. 

I guess this whole post is now wrong, but it was fun to do nonetheless. 
