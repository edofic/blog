---
title: Cheap tagged types in Scala
date: 2015-05-02
---

Sometimes you want to distinguish between different types that have the same underlying representation. For example both `UserId` and `ProductId` could be represented by `Long`. The usual solution is to introduce wrappers in order to make the distinction safe.

```scala
case class UserId(id: Long)
case class ProductId(id: Long)
```

But this introduces runtime overhead of boxing and unboxing over and over which may add up in some cases. Luckily Scala 2.10 introduced value classes. We can ensure no runtime overhead by extending `AnyVal` (this can only be done with classes with one field).

```scala
case class UserId(id: Long) extends AnyVal
case class ProductId(id: Long) extends AnyVal
```

### Inheritance

Let's say that we also want a general `Id` type

```scala
case class Id(id: Long) extends AnyVal
```

So far so good. But we cannot make a value class extend this `Id`! [A value class may only extend universal traits](http://docs.scala-lang.org/overviews/core/value-classes.html). This means we could define a trait that represents the notion of `Id` but we could not make it into a concrete type with values. And even more problems occur when we want to play around with variance. What now?

### Tagging with types
A possible solution is to define an "empty" higher order type and store tags into type parameters.

```scala
trait Tagged[+V, +T]
type @@[+V, +T] = V with Tagged[V, T]

```
`@@` represents a union of two types so values of type `@@[V,T]` will be equal to values of `V` at runtime (as `T` is empty) but we have access to `T` at compile time.

We can also write a simple implicit class helper for easy tagging

```scala
implicit class Taggable[V](val value: V) extends AnyVal {
  def tag[T] = value.asInstanceOf[V @@ T]
  def tagAs[T <: V @@ _] = value.asInstanceOf[T]
}
```

We create values simply by casting as the runtime representation will stay the same.

Usage may look like

```scala
trait TId
trait TUser

type Id = Long @@ TId
type UserId = Long @@ TUser

1234.tag[TUser] // inferred type: Long @@ TUser
5678.tagAs[UserId]
```

A good thing about this approach is that unboxing is automatical since `V <:< @@[V,T]`. But sometimes you may want to untag your values in order to pass them somewhere where you don't want to keep the tags. For this we just need a function that uses the automatic unboxing
```scala
implicit class Untaggable[V, T](val tagged: V with Tagged[V, T]) extends AnyVal {
  def untag: V = tagged
}
```
The trick is just to "pattern match" on the type we implicitly convert.

### Collections
Sometimes you want to tag a collection of something. You could `xs.map(_.tag[Foo])` but this would actually create a new collection at runtime. We can get away just with casting (thus in constant time)!. Notice that collections are nothing special, we may just as well cast a json printer instead of creating a wrapper.
```scala
implicit class TaggableM[M[_], V](val value: M[V]) extends AnyVal {
  def tagM[T] = value.asInstanceOf[M[V @@ T]]
  def tagAsM[T <: V @@ _] = value.asInstanceOf[M[T]]
}
implicit class UntaggableM[M[+_], V, T](val tagged: M[V with Tagged[V, T]]) extends AnyVal {
  def untagM: M[V] = tagged
}
```
This is an abstraction over any `M[_]`. You could write abstractions for other shapes but in practice I never needed anything other since this covers collections and most typeclass instances.

### Variance
An observant reader noticed I defined `@@` to be covariant in both arguments. You probably should leave the value type covariant since it is by nature covariant at runtime but you may change it to invariant if you want to "disable" automatic upcasting. However the tag type may also be contravariant although I found that covarinace is what you naturally expect and covers most cases. Sadly I haven't found a way to abstract over variance.

### Disclaimer

I got this idea from ScalaZ but implemented it in my own way a while ago.
