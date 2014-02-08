---
title: Unfolding Abstract Datatypes
---

# Unfolding Abstract Datatypes by Jeremy Gibbons

Let's compare apples and oranges! Consider object oriented languages versus languages with algebraic datatypes and pattern matching(which are usually functionally oriented). How do we model one in the other?

Let's first consider the easier way - modeling ADTs in OO. You create an interface for the type and inheriting classes for constructors. For pattern matching you check if a value is an instance of a particular class. Scala does this quite well with sealed traits for type and case classes for constructors and deconstructors. 

But the point of the article is the modeling the other way around - objects with ADTs. There are usually three points that define object orientation:

 * polymorphism
 * encapsulation 
 * inheritance

Polymorphism is a given when using algebraic datatypes, since this is how we express sum. Never mind inheritance let's focus on encapsulation, as does the article. 

How can we define an algebraic datatype that we cannot look into but only has a "known public interface"? In other words we want to abstract over the implementation type. We need a big word. We need existential quantification. Simply put, `Complex` is a type inside which exists type of it's implementation but we may not know what it is. 

Let's try a concrete example. Take complex numbers. You could represent them using Cartesian or polar coordinates(or something else). A complex number has no type parameters and is a pair of a representation and a bunch of methods that operate on this representation - the interface.

```haskell
{-# LANGUAGE ExistentialQuantification #-}

data ComplexInterface s = ComplexInterface {
  _re :: s -> Double,
  _im :: s -> Double  
}

data Complex = forall s . Complex (ComplexInterface s) s 
```

The `forall` is the keyword that abstracts over representation type `s`. The crucial part is that `s` effectively seals off the representation from outer world and we cannot peek into it. However we are guaranteed that the interface has compatible methods and this is the only way to operate on the actual representation.

Now let's create a Cartesian representation, starting out with actual data and then the methods. 

```haskell
data CartesianComplex = CartesianComplex Double Double

cartesianComplex = ComplexInterface {
  _re = \(CartesianComplex re _) -> re,
  _im = \(CartesianComplex _ im) -> im
} 
```

We now have our interface and one concrete implementation, but it's cumbersome to use so let's add some helpers.

```haskell
re :: Complex -> Double
re (Complex interface value) = _re interface value

im :: Complex -> Double
im (Complex interface value) = _im interface value

toCartesian c = (re c, im c)
```

Great. But we can extract real and imaginary part and also convert it to a tuple. But how do we construct such a value? Just like in OO. We call the actual constructor and then "cast" it to Complex - the interface - by binding in together method implementations. Implementation type is now forgotten for all the compiler cares. 

```haskell
newCartesianComplex :: Double -> Double -> Complex
newCartesianComplex x y = Complex cartesianComplex $ CartesianComplex x y
```

Imagine you now try to implement add or some other binary function of Complex. You can add it to the specification and even implement it for cartesians, but you cannot write the helper function, or even call the damn thing on your own. Try it. Compiler will complain it cannot prove that both representation types are the same. And he's correct, nobody guarantees they are the same, this is the point of polymorphism. Of course same goes for OO implementation. So how would you implement add? We could add a method `add :: s -> Complex -> s` that uses public interface of `Complex` to extract needed data. Or we can just write a standalone function that does this.

```haskell
add :: Complex -> Complex -> Complex
add c1 c2 = newCartesianComplex r i where
  r = re c1 + re c2
  i = im c1 + im c2
```

This puts less burden on the implementer of the `ComplexInterface` but it also mean the sum will always be represented by CartesianComplex. 

## Data and codata