---
title: Coroutines for free
---

### Motivation

My first run-in with coroutines was with Python's generators. 

```python
def ints():
    i = 0
    while True:
        yield i
        i += 1
```
 This is a function that never returns. Instead it runs in an infinite loop over all natural numbers. But it *yield*s every number. This means it actually stops and "returns" the number to the caller but is available for more execution. Kind-of like implementing an `Iterator` but the control flow is reversed. Now you don't produce a value on demand but rather you run and output values at will while the runtime system stops and resumes you. 

 Now if you would happen to invent `yield` how would you implement it? I think I would hack up the runtime. I would add support for suspending a function. Effectively creating delimited continuations. A function could request that it current scope(and "program counter") be stored on the heap. I personally find this a bit nasty. Even more so in a statically compiled language.

 
### Monads to the rescue

This `yield` is in its essence an effect. Monads are a way to implement effects. Let's try to use them to create *suspendable computations*. But let's start with something simpler and work our way up.


### Trampolines for free

A [trampoline](http://en.wikipedia.org/wiki/Trampoline_(computing)) is a system to implement function calls effectively trading stack for heap. This allows for some nifty tricks like implementing a tail recursion in a system that doesn't optimize for it or doing deep recursive calls that would otherwise result in a stack overflow. In Haskell it just provides explicit boundary between results and thunks.

This is quite trivial to implement. Consider this definition.

```haskell
data Trampoline a = Return a | Bounce (Trampoline a)

runTrampoline :: Trampoline a -> a
runTrampoline (Return a) = a
runTrampoline (Bounce t) = runTrampoline t
```

Function `a -> b` now becomes `a -> Trampoline b`. It can wrap up its result in a `Return` or bounce by wrapping up a thunk(Haskell is non-strict!) in a `Bounce`. `runTrampoline` is implemented with tail recursion but could be done with a `while` loop in languages that don't optimize it.

But now comes the important trick! Consider the `Free` monad from package `free`. 
```haskell
data Free f a = Pure a | Free (f (Free f a))
```
This looks very similar but with an additional layer `f` - a functor. Using the `Identity` functor we can very simply implement the trampoline. 

```haskell
type Trampoline = Free Identity

bounce :: Trampoline a -> Trampoline a
bounce = Free . Identity

runTrampoline :: Trampoline a -> a
runTrampoline = runIdentity . retract
```

And we get a trampoline for free! Also a `Monad` and `Applicative` instances(and more) that make composition of trampoline-using code very easy. This would else require some boilerplate code. I've also included a utility function `bounce` that is now analogous to the `Bounce` constructor from before. 

If the underlying functor is a monad too(which `Identity` is) `Free` knows how to collapse itself - we get `runTrampoline` for free.


### Generators

Finally! Keeping trampolines in mind let think about generators. We could write one as a normal function but with explicit notion of *generator* in the return type. A generator function return a signal that it's done or a value and the next generator function - continuation. 

```haskell
data Producer a = Done | Yield a (Producer a)
```

Again we see structure similar to `Free`. And in fact we can implement this with `Free` if use use the `((,) a)` functor. What is this cryptic type signature? It's a tuple with one missing type parameter. Think of it like `(a,)`. We fix the first member and our `Functor` instance maps over the second member. 

But there's an additional catch. When we're done we don't return anything but `Pure` has a parameter `a`. We just fix this to `()`.

```haskell
type Producer a = Free ((,) a) ()

yield :: a -> Producer a 
yield a = Free (a, return ())
```

The `yield` function only takes a value, it doesn't know about the continuation, it always users `return ()` which is essentially "I'm done". But if we use monadic(or applicative) style we can substitute it with our continuation - the function passed to bind. 

But before looking at an example let's look at the original definition(the one without `Free`) again. It looks a lot like list! In fact it's exactly the same. This means we can simply convert to list just by substituting constructors.

```haskell
consumeProducer :: Producer a -> [a]
consumeProducer (Pure ()) = []
consumeProducer (Free (a, c)) = a : consumeProducer c
```
This is exactly the thing `Free` is meant to do. We described our computation in a pure way and then written an interpreter.

Now the promised example.

```haskell
producerExample = consumeProducer $ do
  yield 1
  yield 2 
  yield 3
```

The semantics is quite trivial and you should immediately see this produces `[1,2,3]`. More importantly it inverts the control allowing us to do arbitrary computations(like recursion) between yields. 

In fact we even did better than Python! Inside a `Producer` you can call to other producers. And they can yield directly! No need for nasty loops like

```python
for e in helper_function():
    yield e
```

You just write

```haskell
helper_function
```

But we also did worse. This approach does not allow for any other effects inside producers. This can be easily mitigated using transformers. But I'll show this in another post.


### Consumers

What if we don't want to convert our producers to list but just consume them immediately? We need some sort of consumers. 

I'll skip a few steps of reasoning and jump straight to the conclusion. We can use the same approach. A consumer is just a function that returns a result or a continuation that takes a value and returns another consumer. This means it can consume an arbitrary amount of values and stop when it wants. It's very similar to folding but with the added control of stopping. Again we inverted the control. It's usually the fold that controls the "looping" not the function we pass into the fold. This concept is called an *iteratee*.

```haskell
data Consumer a b = Result b | Await (a -> Consumer a b)
```

And yet again we see the familiar structure. And yet again our functor is a bit funky: `((->) a)`. It's a function that takes a known parameter but we map over it's result. Think of `a -> `.

```haskell
type Consumer a = Free ((->) a) 

await :: Consumer a a
await = Free return
```

The `await` function might me a bit unintuitive at first. But take a close look at it's type signature. It consumes `a`s and produces an `a`. And it seems reasonable to say it should consume exactly one `a` and then return it. And that's exactly what it does. Function `return` takes one element and packs it up into a `Pure` completing our consumer. 

Then we use monadic(or applicative as I'll show) style to consume more than one.

```haskell
(+) <$> await <*> await 
```

This consumes two values and returns their sum. 

Now we could feed it a list(being dual to our producers) but I'll skip this and write a more interesting function. One that feeds a producer into a consumer. The code for feeding a list is almost exactly the same and you should try to write it as an exercise.  

```haskell
pipe :: Producer a -> Consumer a r -> Maybe r
pipe _ (Pure r) = Just r
pipe (Free (a, c)) (Free f) = pipe c $ f a
pipe _ _ = Nothing
```

What's going here? 

* If our consumer is done we're done too.
* If producer is yielding and consumer is awaiting we feed that value into consumers continuation and recurse.
* Else it's not possible to continue.

And now we can directly consume. But one problem remains. We cannot compose `Producer` and `Consumer` together and get something composable again. We get a `Maybe`. I argue that is is a smell of bad design. We should be able to compose again and again as this will greatly simplify the design of the system using these components. It turns out we can do this but since this is getting a bit long I'll cover this topic in another post.