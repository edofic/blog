---
title: Practical Future[Option[A]] in Scala
---

### Motivation

In real world concurrent code you often come across the `Future[Option[A]]` type(where `A` is usually some concrete type). And then you need to compose these things. This is not straightforward in Scala.

If you've done some Haskell just `import scalaz._` and you can skip the rest of this article. Scalaz library defines a monad typeclass(and many others) that formally specifies what it means to be a monad(not "has a flatMap-ish thingy"). Then it's easy to build abstractions upon this. 

But what if you don't want to add another dependency and just want to make this tiny bit more practical? I can be done and is not that complicated. You will also learn something and maybe even become motivated to bite the bullet and start using Scalaz. 

### Meaning

The `Future[Option[A]]` type combines the notion of concurrency(from `Future`) and the notion of failure(from `Option`) giving you a *concurrent computation that might fail*. What we want is a monad instance for this. With only standard library code you probably would use `Future#flatMap` in combination with `Option#map` and `Option#getOrElse`(or just `Option#fold`). This gets messy and unreadable quite quickly due to loads of boilerplate. Let's fix this!

### Newtyping

A monad in scala means a `flatMap` method, so it's safe to assume we need to define a `flatMap` with this special semantics somewhere. 
You might want to define an implicit class that has the new method *but* it wouldn't work as `Future` already has a `flatMap` method.
I took a page from Haskell's book. When you want to override semantics there you wrap up the type into a *newtype*. This is just compile-time type information that is completely free at runtime. Luckily Scala has this in form of `AnyVal`.

```scala
import concurrent.Future

case class FutureO[+A](future: Future[Option[A]]) extends AnyVal
``` 

I've also made it covariant for ease of use.

### The Monad

It's actually quite easy to implement

import concurrent.{Future, ExecutionContext}

```scala
case class FutureO[+A](future: Future[Option[A]]) extends AnyVal {
    def flatMap[B](f: A => FutureO[B])
                  (implicit ec: ExecutionContext): FutureO[B] = {
        FutureO {
                future.flatMap { optA => 
                optA.map { a =>
                    f(a).future
                } getOrElse Future.successful(None)
            }
        }
    }

    def map[B](f: A => B)
              (implicit ec: ExecutionContext): FutureO[B] = {
        FutureO(future.map(_ map f))
    }
}
```

You need to pull in an execution context and then do the usual boilerplate thing ending up with a wrap again. I've also added `map` method which is trivial but we need it because in Scala for comprehensions are desugared into `flatMap` and `map` to avoid using `point`(abstract constructor) since it's harder to express OO-style.. 

### Usage 

What good is this `FutureO`? Let's do a contrived example. We need a function that might fail - `divideEven` that only divides even numbers. And we need concurrency - we'll divide two numbers concurrently. 

```scala
def divideEven(n: Int): Option[Int] = 
    if (n % 2 == 0) Some(n/2) else None

//first spawn both computations
val f1 = Future(divideEven(14))
val f2 = Future(divideEven(16))

//and combine them
val fc = for {
    a <- FutureO(f1)
    b <- FutureO(f2)
} yield a + b

//prints out Success(Some(15))
fc.future onComplete println
```

It works. How would this look without `FutureO`? 
```scala
val fc = for {
    oa <- f1
    ob <- f2
} yield for {
    a <- oa
    b <- ob
} yield a + b

fc onComplete println
```

Two layers of for comprehensions. It gets even hairier if you have data dependencies between your futures. Consider `divideEven` again but this time we want to divide a number twice in a row. And we'll be keeping the futures around just to prove a point. Let's imagine that `divideEven` does some blocking IO and we want to push it into another tread-pool.

```scala
def divideTwiceF(n: Int): Future[Option[Int]] = {
    val fo = for {
        n1 <- FutureO(Future(divideEven(n)))
        n2 <- FutureO(Future(divideEven(n1)))
    } yield n2
    fo.future
}
```

And it works as expected. The `FutureO` part there is just to alter the monadic semantics. As an exercise try to rewrite this without `FutureO` and squirm in disgust.

### Usability

Inside for comprehension(or manual flatMaps) you can still construct failed futures, throw or return `None`(in a future). However putting in a default value(getOrElse) is a bit trickier and going back to regular `Future` inside same comprehension is impossible. But you can fix this. You can define methods like `orElse` on the `FutureO`. You can also overload the `flatMap` to enable interop with regular futures. However this screws up type inference and I would advise against it as it could introduce some nasty bugs. 

Try to implement combinators you need and leave some comments. Especially if you come across something nice or find a case where `FutureO` is more awkward to use than regular futures.

### Theory 

What we defined is actually a specialized monad transformer. Monads(such as `Future` and `Option`) have this nasty property that they don't compose. You cannot write a function that takes a two monads and outputs a composed one. 

However you can write such a function if you fix one of the monads and only take one as a parameter. This is called a monad transformer. In this case I fixed `Option`. Take a closer look, we are only using `flatMap` and a constructor(`point`) from `Future`. This means we could abstract over the whole monad class. And this is what Scalaz does with `OptionT`. But to do this it needs to define a monad typeclass and instances for each and every monad they find. 
What I did is fix the other monad too and this produces a concrete instance you can use without any typeclasses. There is a downside of course. This is a one-off hack. If you want other transformers you'll have to write them as well. At that point I think would be a good idea to start using Scalaz. 