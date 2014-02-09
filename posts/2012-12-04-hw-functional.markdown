---
title: Homework - functional style (outer sorting) 
--- 

I'm attending
Algorithms and data structures class this semester. Material it self is
quite interesting and one TA is pretty cool too. But I don't like
professor(makes whole experience very much worse) and I believe
homeworks could be much better. Oh, and we didn't even mention
functional approach...you know
[Haskell](http://haskell.org/ "Haskell (programming language)"),
[Scala](http://www.scala-lang.org/ "Scala (programming language)") and
the like. All we do is imperative, C-style code in Java. Enough ranting.
This is how it saw the bright side.


### Le problem

We were doing outer sorting. More specifically: balanced natural outer
merge sort. I hope I translated this right(probably not). In it's
essence the algorithm looks like this
-   you have multiple tracks you read from and write to
-   you have the current element of each track in memory
-   you write out squads(non-descending sub-sequence), this means you
    take the minimum element that is greater than last of if such
    element doesn't exist you take the [minimal
    element](http://en.wikipedia.org/wiki/Maximal_element "Maximal element"). 
-   every time a squad ends your write pointer hops to the next track.
-   repeat until all elements are on single track(hopefully sorted in
    non-descending order)

  ----------------------
  [![Reel of 1/2](http://upload.wikimedia.org/wikipedia/commons/thumb/a/ae/Tapesticker.jpg/300px-Tapesticker.jpg)](http://commons.wikipedia.org/wiki/File%3ATapesticker.jpg)
  Reel of 1/2" tape showing beginning-of-tape reflective marker. (Photo credit: [Wikipedia](http://commons.wikipedia.org/wiki/File%3ATapesticker.jpg))
  ----------------------  

Quite simple right? TA's even provided us classes(talking Java here)
InTrack and OutTrack to manage tracks and I/O. Well my problem is that I
grew to dislike imperative style. Surely it may matter for performance,
but since this is a homework, performance didn't matter - so I wrote
pretty code. I wanted my central code(the heart of the algorithm) to be
a few lines at most. 

This is my final product(bear in mind there was an additional twist:
code should also be capable of sorting in non-ascending order thus the
up variable).
Some additional explanation: all tracks should be in separate files(no
overwriting - for automatic checking). N is number of tracks, prefix is
track name prefix, and i is current iteration.
```java
int i = 0;
Iterable source = new InTrack(inName);
MultiSink sink;
do{
    sink = new MultiSink(prefix, i, N, up);
    for(int n : source) sink.write(n);
    source = new MultiSource(prefix, i, N, up);
    i++;
} while (sink.isMoreThanOneUsed());
```

I'm probably not allowed to share my full solution because of
university rules so I won't.
Now to comment on this. I have a source that's agnostic to the amount of
open files. And I have a similar sink. End-point switching is
implemented in MultiSink and element choosing is in MultiSource. Both
InTrack and MultiSource implement Iterable(and Iterator) so I can use
them in a for-each loop. And code is as pretty as I can get it(remaining
in java). All in all ~300 lines(with InTrack & stuff). After removing
uneeded utility methods and comments ~220 lines. Eww.. Thats way too
much.


Scala to the rescue
-------------------

Lets rewrite this in a functional matter using scala. And while I'm at
it, no vars or mutable collections. 

  ------------------------------
  [![Scala (programming language)](http://upload.wikimedia.org/wikipedia/en/thumb/8/85/Scala_logo.png/300px-Scala_logo.png)](http://en.wikipedia.org/wiki/File%3AScala_logo.png)
  Scala (programming language) (Photo credit: [Wikipedia](http://en.wikipedia.org/wiki/File%3AScala_logo.png))
  ------------------------------

Input can be a collection right? Just implement Traversable. Not really.
The whole point of tracks is they only hold one element in memory(or a
few for efficiency but that's currently not my concern). So a track can
be implemented as a
[Stream](http://en.wikipedia.org/wiki/Stream_%28computing%29 "Stream (computing)")(linked
list with a lazy val for tail).

```scala
def Reader(filename: String) = {
  val sc = new Scanner(new File(filename))
  def loop(): Stream[Int] = {
    if (sc.hasNextInt){
      sc.nextInt() #:: loop
    } else Stream.empty[Int]
  }  
  loop()
}
```

This is the [constructor
function](http://en.wikipedia.org/wiki/Constructor_%28object-oriented_programming%29 "Constructor (object-oriented programming)")
for the input stream. It just returns a recursive value that has a
Scanner in its closure. As stream elements are immutable you get an iron
clad guarantee that sc will stay in sync.   And you get all collections
stuff for free. Moving on, how to abstract over multiple streams? That
should be a stream again right? I kinda feel my code is too complicated
and that it could be done simpler but that's what I came up with
```scala
def MultiReader(prefix: String, phase: Int, N: Int, up: Boolean) = {
  def loop(last: Int, sources: Seq[Stream[Int]]): Stream[Int] = {
    val nonEmpty = sources.filterNot(_.isEmpty)
    if(nonEmpty.length==0)
      Stream.empty[Int]
    else {
      val (low,high) = nonEmpty
        .map(_.head)
        .zipWithIndex
        .partition(t => up && t._1 < last || !up && t._1 > last)
      val (e,i) = (if(high.length>0) high else low).minBy(_._1)      
      e #:: loop(e, nonEmpty.updated(i, nonEmpty(i).tail))    
    }
  }  
  loop(0, (0 until N).map(n => Reader(prefix + "-" + phase + "-" + n)))
}
```

Let's walk through. Again the stream is recursive. It starts with a
collection of Readers set to right files. Then in each step you filter
out empty stream(tracks with no more elements) and partition them
according to the last element(in the argument). If there are higher you
take their minimum else you take the minimum of lower. And loop with
passing on the read element and a new collection - non empty streams
with the read one advanced by one element.

Writer was a bit trickier. It needs internal state, but I prohibited
mutable state. Solution is to return a new Writer containing a new state
every time you write. Then the user must just be careful not to use
stale Writers - not that big a deal.
This is the Writer trait
```scala
trait Writer{
  def write(num: Int): Writer
  def moreThanOneUsed: Boolean
}
```

Very simple interface. And here's the recursive constructor function
```scala
def Writer(prefix: String, phase: Int, N: Int, up: Boolean): Writer = {
  val tracks = (0 until N).map { n => 
    new PrintWriter(new BufferedWriter(new FileWriter(prefix+"-"+phase+"-"+n)))
  }
  def mkWriter(i: Int, last: Int, used: Boolean): Writer = new Writer {
    def write(num: Int) = {
      val (ni,nu) = 
        if (up && num < last || !up && num > last)
          ((i + 1) % tracks.length, true)
        else 
          (i, used)      
      tracks(ni).print(num)      
      tracks(ni).print(' ')      
      tracks(ni).flush()      
      mkWriter(ni, num, nu)    
    }
    def moreThanOneUsed = used  
  }  

  mkWriter(0, if (up) Integer.MIN_VALUE else Integer.MAX_VALUE, used=false)
}
```

Creates all the tracks to be put in the closure. First writer has the
proper start value then every next is constructed like this: figure out
the new values for track number and 'used'(the long if) then actually
write out and return a new writer encapsulating track number and 'used'.
Since these writers are quite lightweight garbage collection pressure
shouldn't be a problem. Especially since the whole process is bound by
I/O. Anyway you could optimize by creating all possible states in
advance and just passing a new reference each time.

Putting it all together.
```scala
def loop(i: Int, source: Stream[Int]){
  val sink = source.foldLeft(Writer(prefix, i, N, up))(_ write _)
  if (sink.moreThanOneUsed) loop(i+1, MultiReader(prefix, i, N ,up))
}

loop(0, Reader(inName))
```

So you take an input stream, fold it over a writer writing in each step.
And if you used more than one track you repeat with a new input stream.
I find this solution to be MUCH MORE elegant. Not to mention it's just
65 lines of scala. But it makes me really sad they don't even mention
[functional
programming](http://en.wikipedia.org/wiki/Functional_programming "Functional programming")
at algorithms course. I'm probably gonna pay the professor and TA's a
visit in near future.
