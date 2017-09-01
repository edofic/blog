---
title: Implementing apply on tuples in Scala
---

One of the first things you learn as a newcomer to Scala is the difference
between a list and a tuple: a list must be homogeneous but a list can be
heterogeneous. That is all elements of a list must have the same types but a
tuple can contain things of different types.

A direct consequence of this is that a list can define a by-index accessor and
a tuple cannot. You can use `list(3)` but you can't do `tuple(3)` - you need to
do `tuple._4` (and there is that pesky off-by-one).

So let's use the awesome powers of Scala to negate this and implement `apply`
method on tuples.

## First steps

Let's start with with baby steps and not tackle full blown `apply` with integer
index but instead do an approximation with special access constants. Something
like

```scala
val t1 = ("1", 123)
val t2 = (false, 1234, "foobar")

println(t1(_2)) // prints 123
println(t2(_1)) // prints false
println(t2(_3)) // prints "foobar"
```

It's not a big step from regular tuple accessors but it's a big move since it
introduces a single polymorphic `apply`.

To pull this off we'll type classes in conjunction with singleton types. The
`apply` will take in one of the singletons and use implicit resolution to pull
in the function that does the proper projection. The important part is that the
implicit resolution also needs to compute the output type of the `apply`
method.

```scala
// this integer references can be avoided when/if SIP-23 is implemented
val _1 = 1: Integer
val _2 = 2: Integer
val _3 = 3: Integer
// ... up to 22
```

The type-class needs to correlate the tuple type (`T`), the output type (`A`)
and the index (`N`). The instances just implement the obvious connection
between indexes and tuple accessors.

```scala
trait TupleGet[T, A, N <: Integer] {
  def get(t: T): A
}
object TupleGet {
  implicit def tuple2get1[A, B] = new TupleGet[(A, B), A, _1.type] {
    def get(t: (A, B)): A = t._1
  }
  implicit def tuple2get2[A, B] = new TupleGet[(A, B), B, _2.type] {
    def get(t: (A, B)): B = t._2
  }
  implicit def tuple3get1[A, B, C] = new TupleGet[(A, B, C), A, _1.type] {
    def get(t: (A, B, C)): A = t._1
  }
  implicit def tuple3get2[A, B, C] = new TupleGet[(A, B, C), B, _2.type] {
    def get(t: (A, B, C)): B = t._2
  }
  implicit def tuple3get3[A, B, C] = new TupleGet[(A, B, C), C, _3.type] {
    def get(t: (A, B, C)): C = t._3
  }
  // ...
}
```

If you haven't seen types like `_1.type` before: this is precisely the type of
the value `_1`. Only its exact copies can inhabit it at any time.

Now we can write the extension `apply` method that uses this mechanism.

```scala
implicit class TupleAccessor[T](tuple: T) {
  def apply[A](i: Integer)(implicit getter: TupleGet[T, A, i.type]): A =
    getter.get(tuple)
}
```

And it will work. But you need to use `_1` instead of regular `1` to give the
implicit resolution mechanism sufficient clues. Not pretty enough.


## Enter Nats

To integrate deeper with the language we need to lift the integers into the
type level somehow.

Since I don't feel like writing 22 instances let's start off with unary
encoding.

```scala
sealed trait TNat
class TZero extends TNat
class TSucc[A <: TNat] extends TNat
```

And a corresponding value level representation that also keeps the type around
with an implicit conversion from integers so this is all auto-magical.

```scala
sealed trait Nat {
  type N <: TNat
}

case object Zero extends Nat {
  type N = TZero
}

case class Succ(prev: Nat) extends Nat {
  type N = TSucc[prev.N]
}

object Nat {
  implicit def fromInt(n: Int): Nat = {
    assert(n >= 0)
    if (n == 0) {
      Zero
    } else {
      Succ(fromInt(n - 1))
    }
  }
}
```

The trick here is the type member that is defined to match the number itself.
This way the type checker gets access to the "value. Now we just need to
rewrite the instances to use this.

Except it doesn't work. This mechanism will generate the `Nats` but it's too
weak - `N` will just be `n.N` and this is not enough information to power
implicit resolution. If only there was a way to fully generate the types at
compile time...


## Macros to the rescue

We can keep the `Nat` type and the implicit conversion into it but implement it
as a macro that just generates `TNat`s directly (without recursion). This will
give strong enough types to power implicit resolution.


```scala
trait Nat {
  type N <: TNat
}

object Nat {
  def fromIntImpl(c: Context)(n: c.Expr[Int]): c.Expr[Nat] = {
    import c.universe._
    val value = c.eval(n)
    val type_ = (1 to value).foldLeft("TZero")( (t, _) => s"TSucc[$t]")
    c.Expr[Nat](c.parse(s"new Nat { type N = $type_ }"))
  }

  implicit def fromInt(n: Int): Nat = macro fromIntImpl
}
```

Now we can rewrite the instances.

```scala
trait At[T, A, N <: TNat] {
  def apply(t: T): A
}
object At {
  implicit def tuple2_at0[A, B] = new At[(A, B), A, TZero] {
    def apply(t: (A, B)): A = t._1
  }
  implicit def tuple2_at1[A, B] = new At[(A, B), B, TSucc[TZero]] {
    def apply(t: (A, B)): B = t._2
  }
  implicit def tuple3_at0[A, B, C] = new At[(A, B, C), A, TZero] {
    def apply(t: (A, B, C)): A = t._1
  }
  implicit def tuple3_at1[A, B, C] = new At[(A, B, C), B, TSucc[TZero]] {
    def apply(t: (A, B, C)): B = t._2
  }
  implicit def tuple3_at2[A, B, C] = new At[(A, B, C), C, TSucc[TSucc[TZero]]] {
    def apply(t: (A, B, C)): C = t._3
  }
}
```

And it works!

```scala
object App {
  implicit class TupleGet[T](t: T) {
    def apply[A](n: Nat)(implicit at: At[T, A, n.N]): A = at(t)
  }
  def main(args: Array[String]) {
      val t1 = ("1", 123)
      val t2 = (false, 1234, "foobar")

      println(t1(0))
      println(t2(1))
      println(t2(2))
  }
}
```


## Simplifying

Since we already have a macro sitting in there why not have it do all the heavy
lifting? We can cut bulk of the code and remove a lot of the complexity of the
implicits if we just use a macro that transforms `a.apply(n)` into `a._${n +
1}`.

```scala
def tupleApplyImpl(c: Context)(n: c.Tree): c.Tree = {
  import c.universe._
  val value = c.eval(c.Expr[Int](n))
  val q"$_($tuple)" = c.prefix.tree
  c.parse(s"${tuple.toString}._${value + 1}")
}

implicit class TupleOps[T](tuple: T) {
  def apply(n: Int): Any = macro tupleApplyImpl
}
```

This works but I think it's not the best idea. See, macros don't compose. You
cannot really use this from another function, you always need to statically
know the index. As where with the previous implementation you are good to go as
long as you have a good `Nat` and the `At` instance. Which you can pass in
programmatically. And you just push the implicit conversion to `Nat` a layer
out. This way you can re-use the indexing mechanism for other operations and
you just have a single macro sitting in the background powering things instead
of having a bunch of one-offs.

In fact this is exactly what
[shapeless](https://github.com/milessabin/shapeless/wiki) already does! It uses
very similar ideas and pushes them much further. But most importantly, it
already includes our `apply` so if you actually want to do this just import
shapeless.

```scala
import shapeless.syntax.std.tuple._

val t1 = ("1", 123)
val t2 = (false, 1234, "foobar")

println(t1.apply(0))
println(t2.apply(1))
println(t2.apply(2))
```
