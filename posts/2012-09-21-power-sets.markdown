---
title: Power sets
---

It all started out when a friend of mine([check him out](http://naspletu.org/))
 told me a story about someone having an
interview at [Google](http://google.com/ "Google"). He was [live
coding](http://en.wikipedia.org/wiki/Live_coding "Live coding") and
asked to implement a function that computes a [power
set](http://en.wikipedia.org/wiki/Power_set "Power set") of a given set.
He totaly over-engineered it and after an hour of fiddling came up with
4-liner in python. Ok. Big deal. How hard can it be? I immediately
started to brainstorm a solution of my own(and Matevž helped).
Paraphrased train of thought below. 
So what is a power set? Say a power set of set A. It's a set of all
[subsets](http://en.wikipedia.org/wiki/Subset "Subset") of A. It's all
possible combinations of including and excluding single elements. So
it's like applying filters to the original set. Wait a sec, a set of
size n has 2 to the n-th power subsets. It all matches up. N-bit number
are said filters. So you count from 0 to 2^n-1 you enumerate exactly
all the filters and therefore all the subsets. You just calculated the
power set. This took about 30sec in real time. Another minute to sketch
the implementation. You implement set as Lists. Let's say
[Java](http://www.oracle.com/technetwork/java/ "Java (programming language)")
is the language of choice. Any reasonable input would be shorter than
64 elements, at least on current machinery. So long can be used as
counter. And now a nested for loop to iterate over elements and
appending them if counter && (1 << n). Simple. Won't post the code
since I never wrote it.

### The original blogpost

This happened in a lab at Faculty for Computer and Information science.
On my way home I started thinkking....what is that 4-line
implementation? It can't be my algorithm, that won't fit - it just isn't
that elegant. Matevž also told me who was the original author and it was
someone whose blog I was subscribed to(I still am, it's pretty awesome,
[check it out](http://swizec.com/)). So I searched for the original
blogpost. And [found
it](http://swizec.com/blog/a-google-phone-interview/swizec/3813). It
wasn't four lines.
```scala
def powerset(set):
      binaries = [bin(a) for a in range(2^len(set))]
      power = []
      for yeses in binaries:
        subset = []
 yeses = str(yeses)
 for i in range(len(yeses)):
                 if yeses[i] == “1”:
          subset.append(set[i])
 power.append(subset)
      return power
```
But more importantly, it's essentially my algorithm. (And it didn't took
him an hour, read the post, it's interesting) And then I realised. This
is ugly. Quite ugly. It may be efficient and allow for doing iteration
and evaluating it lazily, but oh boy it's ugly. Can't it be done in a
simpler way?

#### Recursive solution

So by asking how I **compute** a power set I obtained the imperative
algorithm above.

  -----------------------------
  [![Balanced tree](http://upload.wikimedia.org/wikipedia/commons/thumb/3/33/Balanced_tree.png/300px-Balanced_tree.png)](http://commons.wikipedia.org/wiki/File%3ABalanced_tree.png)
  Balanced tree (Photo credit: [Wikipedia](http://commons.wikipedia.org/wiki/File%3ABalanced_tree.png))
  -----------------------------

Time to change the question. What **is** a power set?

It's a set of leaves of a balanced binary decision tree. Every inner
node represents a decision to either take or leave out the element at
respective index(I assume elements are indexed). That is root node
responds to first element, it's children to second element, etc. Every
leaf can be computed by following the decision path.  It's
a **recursive** structure. So the tree can be decomposed to left and
right subtree and recombined. An this is repeated until you get to
leaves. So a power set of A is an union of the power set of "A minus the
first element" and the same power set again but with every subset added
the first element. Subtrees are represented by the two powersets and the
recombination by adding head element to one set of sets and then doing
the union. And to stop the recursion: power set of an empty set is a set
containing just an empty set.

Now let's do that in code, my language of choice now is scala
```scala
def power[A](set: List[A]): List[List[A]] = set match {
    case Nil => List(Nil)
    case head::tail => power(tail) flatMap (set => List(set, head::set))
}
```

Yay just four lines. and the last one barely matters. But I can do even
better. The same code in haskell
```haskell
power [] = [[]]
power (head:tail) = power tail >>= \set ->[set, head:set]
```

So you apply the lambda that creates both
sides of the union to your list and then flatten it to produce the same
result. This code may look a bit cryptic at first sight but once you
parse the syntax you can grasp the intention much more quickly than the
imperative version. At least I can. Recursion is quite nice once you get
used to it.
