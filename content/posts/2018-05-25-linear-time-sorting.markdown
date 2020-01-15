---
title: Linear time sorting
date: 2018-05-25
description: "Sorting arrays in linear time"
---

Did you know that you can sort data in time linear with respect to the lenght
of said data?  Some people dismiss linear time sorts because they've learned
that `n log n` is the bottom bound for sorting an arbitrary input.

# But *n log n* is optimal!?

Let's even sketch out a proof. Any function that sorts an array of elements
will in fact figure out a permutation of elements into a sorted order. Even if
your algorithm does not work in terms of permutations it can be viewed as a
black box that computes a permutation.

There are `n!` permutations of an array of length `n`. Each time you compare
two elements there are three possible possible results (less, equal, greater).
You can now imagine a decision tree that has a comparison at each inner node
and a permutation at each leaf.  An optimal tree will be balanced with a
branching factor of 3. Because we know how many leaves we have (`n!`) we can
compute the depth of the tree. Using Stirling's approximation for factorial
immediately yields the result `n log n`. This means you have to do at least `n
log n` comparisons to figure out the permutation.

But this is a theorem about **comparison** sorts. You cannot do less than `n
log n` *comparisons* if all you have is pair-wise comparison. Assuming you can
do comparisons in `O(1)` this makes your sort `O(n log n)`. But beware: if your
comparison is not constant-time this is not true any more. If you are comparing
arbitrary precision integers or strings your sort will be `O(k n log n)` where
`k` is the length of a single element.

This still leaves place for improvement if you have more structure on your
data. I'm going to talk about the algorithms that have access to the underlying
binary representation of data. Sorting will mean sorting the bit-strings that
represent the data. This is enough for most sorts since you can usually build
up a surrogate key that conforms to this.

# A naive approach

Let `A` be an input array of 16 bit integers. We allocate another array of size
`2 ^ 16` - 64k. Then we loop over the input

```python
tmp = [0] * (2 ** 16)
for a in A:
  tmp[a] += 1
```

This is enough to spit out the sorted array

```python
var sorted = [];
for i, c in enumerate(A):
  for _ in range(c):
    sorted.append(i)
```

And we didn't use comparisons! I will even argue this is linear with respect to
the length of the input array. We allocated an array of constant size. This is
`O(1)`. Then we did one operation per input element. And another when
outputting it. It's a hand-wavy argument but it's only for illustration.  A
slightly less naive version of this is known as counting sort and actually
performs really well on some types of input.

# A generalization

What if we want to use different size of integers? If we just used this
algorithm for 64 bit integers we would run out of space before even starting.

We can tackle this problem by looking at a constant number of bits of each
element at a time. The simplest is just a single bit, but a byte works out
quite well in practice. So now we sort our array by the least significant byte.
We now cannot use the naive counting because index does not tell us anything
about the remaining bytes. We have to replace counters by buckets that hold
elements.

```python
def get_byte(a, i):
  return (a >> (8 * i)) & 0xff

def one_pass(ar, i):
  buckets = [[] for _ in range(256)]
  for a in ar:
    buckets[get_byte(a, i)].append(a)
  return sum(buckets, [])
```

And the sorted result is just the concatenation of the buckets. The good thing
about this is that it's a stable sort. It will preserve the order of elements
that it considers equal (at one byte that is). So we can repeat this with other
bytes (going from right to left) and we will end up with a sorted array. This
was not very intuitive to me but if you work through an example by hand you
quickly figure it out.

```python
def sort(ar):
  num_bytes = int(math.ceil(math.log(max(ar), 2)))
  for i in range(num_bytes):
    ar = one_pass(ar, i)
  return ar
```

We can compact it all into one simple function

```python
def sort(ar):
  num_bytes = int(math.ceil(math.log(max(ar), 2)))
  for i in range(num_bytes):
    buckets = [[] for _ in range(256)]
    for a in ar:
      buckets[(a >> (8 * i)) & 0xff].append(a)
    ar = sum(buckets, [])
  return ar
```

This is known as Least Significant Byte (LSB) Radix Sort. Let's take a quick
look at it's complexity. It does a constant amount of work (one append) per
each byte of input. So it's linear with respect to number of input bytes. For
fixed precision numbers this means it's even linear with respect to number of
input elements. Great! This is very simple. Well... not so quick. There's many
variations and tradeoffs to be made. One I quote like is the [American Flag
Sort](https://en.wikipedia.org/wiki/American_flag_sort) - which is an in-place
top down variation. This entails that you can use it to efficiently sort things
other than integers -> if you can discriminate them into equality buckets.

# Conclusion

I think linear time sorting is under-appreciated and could be used more
problems more efficiently but it isn't easily accessible enough.  I only became
interested in the topic after hearing [Edward Kmett's awesome talk on the
topic](https://www.youtube.com/watch?v=cB8DapKQz-I) He also wrote [an amazing
library for Haskell](https://hackage.haskell.org/package/discrimination) that
makes linear time sorting and grouping trivial to use when you need it. This is
the only library of the sort that I know of, buy as you've seen above, even
rolling your own is not that hard. So if that extra *log n* is bothering you,
remember the radix.
