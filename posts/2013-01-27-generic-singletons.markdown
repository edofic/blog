---
title: Generic singletons through dependent method types
---

Ever
tried to write a
[generic](http://en.wikipedia.org/wiki/Generic_programming "Generic programming")
singleton? It's an oxymoron of a sort. But sometimes my brain dreams up
funny concepts to solve the problem at hand. Sadly I cannot remember
what I wanted to use them for. Anyway I think I just made all the
[methods](http://en.wikipedia.org/wiki/Method_%28computer_programming%29 "Method (computer programming)")
generic and solved it this way. But this doesn't really express the
notion of one entity that's agnostic to type of the parameters. With
generic methods you get a bunch of disconnected units - at least that's
the picture in my head.

### General

 --------------
 ![Challenge Accepted!](http://farm8.static.flickr.com/7192/6857158741_a4e3d23649_m.jpg "Challenge Accepted!")
 Challenge Accepted! (Photo credit: pierreee)
 --------------

So we've
established that generic singletons are just a silly idea, now let's try
to implement it - you know "Challenge accepted!" style.

```scala
object Singleton[A] { //error: ';' expected but '[' found.
  def id(a:A) = a
}
```

Luckily this naive approach does not compile. Try it.

```scala
class Singleton[A]{
  def id(a: A) = a 
}

def Singleton[A] = new Singleton[A]
```

But this is not a singleton anymore as there is a new instance created
everytime.

```scala
object Singleton{
  class Instance[A] private[Singleton] (){
    def id(a: A) = a 
  }

  private[Singleton] val cache = 
    collection.mutable.HashMap[Class[_],Instance[_]]()
  def apply[A]()(implicit m: Manifest[A]) = 
    cache.getOrElseUpdate(m.erasure, new Instance[A]).asInstanceOf[Instance[A]]
}

//use as
Singleton[Int].id(1)

//and some sugar
def singleton[A:Manifest] = Singleton[A]()

//now use like
singleton[String].id("hi")
```

Didn't bother with [thread
safety](http://en.wikipedia.org/wiki/Thread_safety "Thread safety"). And
dragging around Manifest just for caching is a bit of a pain, luckily
elevated by context bounds.

### Type specific

That worked for classes that don't bother about the generic types...what
if we want specific behavior per type? [Type
classes](http://en.wikipedia.org/wiki/Type_class "Type class"). 
A trait for general interface and objects for specific types. Magic lies
in making these objects implicit and having a method that gives you the
right one!

```scala
trait Singleton[A]{
  def f(a: A): A
}

implicit object IntSingleton extends Singleton[Int]{
  def f(n: Int) = n + 1
}

implicit object StringSingleton extends Singleton[String]{
  def f(s: String) = "Hello " + s
}

def Singleton[A](implicit s: Singleton[A]) = s

Singleton[Int].f(1) // == 2
```

### Specific functionality

Now let's add a little twist.

```scala
implicit object IntSingleton extends Singleton[Int]{
  def f(n: Int) = n + 1
  def g(n: Int) = n - 1
}

//this line doesn't compile
Singleton[Int].g(1)
```

and it doesn't compile for a reason. [Static
type](http://en.wikipedia.org/wiki/Type_system "Type system") of
returned from Singleton[Int] is Singleton[Int]. And there is no method g
defined in trait Singleton. But we are sure it exists. We could call
IntSingleton directly but this makes whole Singleton abstraction
worthless. Luckily for us scala 2.10 has a feature called method
[dependent
types](http://en.wikipedia.org/wiki/Dependent_type "Dependent type"). To
my best understanding this means support for return types that depend on
the type of parameters. And here we want the type to reflect which
instance is being returned. Nothing suits the job better as singleton
types. Singleton type is a type that has only one value. So an object A
will be of type A.type. This type is worthless for inference but great
when you're doing the type tags. We just need a little fix

```scala
def Singleton[A](implicit s: Singleton[A]): s.type = s
```

This method will return the correct instance along with it's type while
giving you a nice abstraction of [type
constructor](http://en.wikipedia.org/wiki/Type_constructor "Type constructor").
 
