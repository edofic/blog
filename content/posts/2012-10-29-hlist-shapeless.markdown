---
title: Cool Monday - HList and Shapeless
---

HList as in
[heterogenous](http://en.wikipedia.org/wiki/Homogeneity_and_heterogeneity "Homogeneity and heterogeneity")
lists. This means every element is of different type. Yeah sure, just
list List in Java, but that is in no way typesafe. I want compiler to know the
type of every element and stop me if I try to do something silly.

### Linked lists to the rescue

So what's a [linked
list](http://en.wikipedia.org/wiki/Linked_list "Linked list") anyway? A
sequence of nodes with pointers to next. And a nice implementation(still
talking
[Java](http://www.oracle.com/technetwork/java/ "Java (programming language)")
here) would be generic to allow
[type-safety](http://en.wikipedia.org/wiki/Type_safety "Type safety")
for homogeneous lists. It turns out generics are solution for HLists too.
Just introduce additional [type
parameter](http://en.wikipedia.org/wiki/TypeParameter "TypeParameter").
Apocalisp has a [great
post](http://apocalisp.wordpress.com/2008/10/23/heterogeneous-lists-and-the-limits-of-the-java-type-system/)
on implementing them in Java.
Java requires A LOT of [type
annotation](http://en.wikipedia.org/wiki/Type_signature "Type signature").
It works but it's just painful and it doesn't pay off.

### [Type inference](http://en.wikipedia.org/wiki/Type_inference "Type inference") to the rescue

Type inference gets rid of this problem entirely. Let's implement whole
working HList in scala.
```scala
abstract class HList[H,T<:HList[_,_]] {
    def head: H
    def tail: T  def ::[A](a: A) = Hcons(a, this)
}

object HNil extends HList[Nothing, Nothing]{
    def head = throw new IllegalAccessException("head of empty hlist")
    def tail = throw new IllegalAccessException("tail of empty hlist")
}

case class Hcons[H,T<:HList[_,_]](
    private val head: H, private val tail: T) extends HList[H, T]
```
So this list can be instantiated like this
```scala
scala> val myHList = 1 :: "hi" :: 2.0 :: HNil
myHList: Hcons[Int,HList[java.lang.String,HList[Double,HList[Nothing,Nothing]]]] =
    Hcons(1,Hcons(hi,Hcons(2.0,HNil$@dbb62c)))
```

And it just works. Scala compiler does all the heavy lifting with type
annotations. This implementation bare bones and doesn't provide any
useful methods(even random access!). Check out Miles Sabin's [shapeless
project](https://github.com/milessabin/shapeless) for a useful
implementation and much more. I provides indexing, map, fold,
concatenation, type-safe casts, conversions to tuples(and abstracting
over arities!) and back. And even conversions with case classes. Just
click the link above and read the readme. It's awesome.

