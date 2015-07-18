---
title: Linear time sorting is cool
---

I think linear time sorting is under-appreciated and could be used more problems more efficiently but it isn't easily accessible enough.
I only became interested in the topic after hearing Edward Kmett's awesome talk on the topic. He also wrote an amazing library for Haskell that makes linear time sorting and grouping trivial to use when you need it. But what are we poor developers that cannot use Haskell forced left with? Not much really.

Last few days I've been optimizing the click-to-sort functionality in a web app. Nothing special. I got to the needed performance using `Array.prototype.sort` and by memoized surrogate key. Computation of the key is still the bottle neck but it's fast enough for our use-case.
But it got me thinking. The price of sorting will at some data size factor in and the actual sorting will become the slow part. Can we do better? I sought out a library for doing a linear time sort and only found one specialized for 32 bit numbers. But I have strings.
I guess I'll have to write my own. I plan to explore the different approaches over the coming days and write about the results. Hopefully I manage to write something usefully fast and package it up for others to use.

# But *n log n* is optimal!?

Some people dismiss linear time sorts because they've learned that `n log n` is the bottom bound for sorting an arbitrary input. Yes this is of course true. But this is a theorem about **comparison** sorts. You cannot do less than `n log n` *comparisons* if all you have is pair-wise comparison. Assuming you can do comparisons in `O(1)` this makes your sort `O(n log n)`. This still leaves place for improvement if you have more structure on your data. I'm going to talk about the algorithms that have access to the underlying binary representation of data. Sorting will mean sorting the bit-strings that represent the data. This is enough for most sorts since you can usually build up a surrogate key that conforms to this. I will talk more about time complexity later.

# A naive approach

Let `A` be an input array of 16 bit integers. We allocate another array of size `2 ^ 16` - 16k. Then we loop over the input

    var tmp[2 ^ 16]
    for a in A:
      tmp[a] += 1

This is enough to spit out the sorted array

   var sorted = [];
   for i, c in enumerate(A):
     for _ in range(c):
       sorted.append(i)

And we didn't use comparisons! I will even argue this is linear with respect to the length of the input array. We allocated an array of constant size. This is `O(1)`. Then we did one operation per input element. And another when outputting it. It's a handwavy argument but it's only for illustration.
A slightly less naive version of this is known as counting sort and actually performs really well on some types of input.


