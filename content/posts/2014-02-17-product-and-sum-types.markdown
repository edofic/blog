---
title: Product and sum types (Go and Haskell)
---

Let's take a look at Go-language functions. They let you return more than one result. Like in
```go
func foo(n int) (int, int) {
    return n+1, n+2
}
```
This function `foo` returns two numbers. But does it really? Pairs(and tuples in general) aren't a first class citizen in Go. Yes you can return them from function and you can even return a result of a function that returns the same signature as you but then you must immediately bind each result to separate variable.
```go
a, b := foo(0) //this is legal
pair := foo(0) //this is not
```

I haven't checked but it seems pretty obvious that this multiple-results is just syntax sugar that passes in pointers to the variables you eventually bind to. So this gives you speed because there is little overhead but it abstracts poorly - just like Go in general.


### Product types

Product type is a fancy name for something like a tuple. Usually just a tuple that has a "tag" or a "name" so we can pattern match on that. And where does the product name come from? Let's for a moment consider a naive representation of types with with finite sets. This doesn't really work but it's enough for this analogy and generalizes nicely to real world. 

Consider type `A`(modelled by set A) and type `B`(and set B). Now you can get type `(A,B)` it's model being set `A x B` - the Cartesian **product**. Well this is just the intuition. 

More exactly: the cardinality of `(A,B)` is product of cardinalities of `A` and `B`.


### Sum types

The other thing from the title. We got out product. Now we want sums. Going back to the set model, we want a type that would be modelled by a set of cardinality that is the sum of cardinalities of both sets. Now you might be thinking union... But there is a problem with union as sets `A` and `B` might overlap. So we must tag each element to remember from which set it came. And then do union. We'll see an example later. 

This can be done in Go. Kinda. But it's very cumbersome. We can model this with OO-like polymorphism. Which we have to model in Go with interface in structs. 
```go
type AorB interface {}

type A struct {
    Value int
}

type B struct {
    Value int
}
```

But when we have a value of type `AorB` we have to try to cast and check for success. Urhg boilerplate.


### Enter Haskell

Can we do better? Sure. We can use Haskell and it's ADTs.

ADT stands for *Algebraic Data Types* which just means we can express sum and product types. Great. Let's take a look at the syntax

```haskell
type Product1 a b = (a, b) -- tuples are build in

data MyProduct a b = MyProduct a b -- or a custom tagged tuple

data AorB a b = A a | B b -- generic product
```

Much simpler. 


### Go errors

A typical Go function has some side-effects and computes a value or returns an error. For example

```go
func f(parameter string) (Result, error) {
    err := launchMissiles()
    if err != nil {
        return nil, err
    }
    res := ....
    return res, nil
}
```
You can clearly see that `f` is returning a product type. But is that really what we mean? No. The semantics is either `Result` either `error`. Error being set usually implies that `Result` isn't well defined. 

Let's translate this into Haskell just to ease up on syntax

```haskell
f :: String -> IO (Option Result, Option Error) 
f parameter = do 
    err <- launchMissiles
    case err of 
        | Just error -> return (Nothing, Just error)
        | Nothing    -> do
            ...
            return (Just res, Nothing)
```

I said *ease up* but this is definitely more noisy. Why? Because we're doing it wrong. We're modelling sum types with product types. 

`IO (Option Result, Option Error)` is a  nasty type signature. What we really want is `IO (Either Error Result)`. `Either` is built-in generic product of two elements. It's usually used for this case: result or error. 

```haskell
f :: String -> IO (Either Error Result)
f parameter = launchMissiles >>= next where
    next (Left err) = return (Left error)
    next (Right ()) = do
        ...
        return $ Right res
```

### Monad transformers

Maybe it's just me, but this is much better. But that part `next (Left err) = return (Left error)` feels weird. We're just repackaging. This is boilerplate. And what should you do with boilerplate? Extract it into a function. Luckily someone else already did it. What we want is *Either monad transformer*. A construct that lets us wrap up `m (Either a b)` for any monad `m` and handles errors automatically. If we fill in our types we get `EitherT IO Error Result`.
Notice the lack of parentheses in the type. This makes me happy. Now let's write our function. 

```haskell
f :: String -> EitherT IO Error Result
f parameter = do
    launchMissiles
    ...
    return $ Right res
```

*Much* cleaner. As I promised. In fact this is a bit java-esque. Looks like checked exceptions - we declared out error type in the signature and it looks like it's throwing because error handling is not to be seen in code. Well it isn't, it's more akin to the starting Go code just with all boilerplate extracted out. And thanks to the *do notation* the invocations of functions that handle errors are now invisible. How cool is that?

### What's happening?

In fact the *do notation* is just syntax sugar for `>>=` operator(bind). Kind of overloaded semicolon - however weird this sounds. In between every line we have an implicit semicolon(or you can type it, it's optional) that gets translated into call of bind. 

Implementation of bind(for `Either`) evaluates it's first argument(first line) and pattern matches it. If it's a left(in our case Error) it returns left immediately(without even evaluating further lines as Haskell is lazy), if it's right it extracts the value "puts it in scope" and passes it on to further lines. In fact this is just `Either`, `EitherT` threads another monad(`m`) in between to enable other effects - in our case `IO`. 