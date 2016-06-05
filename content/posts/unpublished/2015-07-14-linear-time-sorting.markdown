---
title: Linear time sorting is cool
---

Did you know that you can sort data in time linear with respect to the lenght of said data?
Some people dismiss linear time sorts because they've learned that `n log n` is the bottom bound for sorting an arbitrary input.

# But *n log n* is optimal!?

Let's even sketch out a proof. Any function that sorts an array of elements will in fact figure out a permutation of elements into a sorted order. Even if your algorithm does not work in terms of permutations it can be viewed as a black box that computes a permutation.

There are `n!` permutations of an array of length `n`. Each time you compare two elements there are three possible possible results (less, equal, greater). You can now imagine a decision tree that has a comparison at each inner node and a permutation at each leaf.
An optimal tree will be balanced with a branching factor of 3. Because we know how many leaves we have (`n!`) we can compute the depth of the tree. Using Stirling's approximation for factorial immediately yields the result `n log n`. This means you have to do at least `n log n` comparisons to figure out the permutation.

But this is a theorem about **comparison** sorts. You cannot do less than `n log n` *comparisons* if all you have is pair-wise comparison. Assuming you can do comparisons in `O(1)` this makes your sort `O(n log n)`. But beware. If your comparison is not constant-time this is not true any more. If you are comparing arbitrary precision integers or strings your sort will be `O(n k log n)` where `k` is the length of a single element.

This still leaves place for improvement if you have more structure on your data. I'm going to talk about the algorithms that have access to the underlying binary representation of data. Sorting will mean sorting the bit-strings that represent the data. This is enough for most sorts since you can usually build up a surrogate key that conforms to this.

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

And we didn't use comparisons! I will even argue this is linear with respect to the length of the input array. We allocated an array of constant size. This is `O(1)`. Then we did one operation per input element. And another when outputting it. It's a hand-wavy argument but it's only for illustration.
A slightly less naive version of this is known as counting sort and actually performs really well on some types of input.

# A generalization

What if we want to use different size of integers? If we just used this algorithm for 64 bit integers we would run out of space before even starting.

We can tackle this problem by looking at a constant number of bits of each element at a time. The simplest is just a single bit, but a byte works out quite well in practice. So now we sort our array by the least significant byte. We now cannot use the naive counting because index does not tell us anything about the remaining bytes. We have to replace counters by buckets that hold elements.

```python
def get_byte(n, i):
  return

def one_pass(ar, i):
  buckets = [[] for _ in range(256)]
  for a in ar:
    buckets[get_byte(a, i)].append(a)
  return sum(buckets, [])
```

And the sorted result is just the concatenation of the buckets. The good thing about this is that it's a stable sort. It will preserve the order of elements that it considers equal (at one byte that is). So we can repeat this with other bytes (going from right to left) and we will end up with a sorted array. This was not very intuitive to me but if you work through an example by hand you quickly figure it out.

```python
def sort(ar):
  for i in range(math.ceil(max(map(math.log2, ar)) / 8)):
    ar = one_pass(ar, i)
  return ar
```

We can compact it all into one simple function

```python
def sort(ar):
  for i in range(math.ceil(max(map(math.log2, ar)) / 8)):
    buckets = [[] for _ in range(256)]
    for a in ar:
      buckets[(a >> (8 * i)) & 0xff].append(a)
    ar = sum(buckets, [])
  return ar
```

This is known as Least Significant Byte (LSB) Radix Sort. Let's take a quick look at it's complexity. It does a constant amount of work (one append) per each byte of input. So it's linear with respect to number of input bytes. For fixed precision numbers this means it's even linear with respect to number of input elements.


# TODO

I think linear time sorting is under-appreciated and could be used more problems more efficiently but it isn't easily accessible enough.
I only became interested in the topic after hearing Edward Kmett's awesome talk on the topic. He also wrote an amazing library for Haskell that makes linear time sorting and grouping trivial to use when you need it. But what are we poor developers that cannot use Haskell forced left with? Not much really.

Last few days I've been optimizing the click-to-sort functionality in a web app. Nothing special. I got to the needed performance using `Array.prototype.sort` and by memoized surrogate key. Computation of the key is still the bottle neck but it's fast enough for our use-case.
But it got me thinking. The price of sorting will at some data size factor in and the actual sorting will become the slow part. Can we do better? I sought out a library for doing a linear time sort and only found one specialized for 32 bit numbers. But I have strings.
I guess I'll have to write my own. I plan to explore the different approaches over the coming days and write about the results. Hopefully I manage to write something usefully fast and package it up for others to use.
