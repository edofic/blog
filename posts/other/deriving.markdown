# Deriving tricks

This day is focused on extensions of the `deriving` mechanism. This allows us to automatically derive some useful common instances for our own datatypes and thus leave us with less code to maintain.

## Functor

Functors are a very useful beast. There is the usual notion of the container like `List` or `Maybe` that you `fmap` over element-wise but you can do more powerful stuff with functors. 

The first motivating example I would like to show are free monads. You can get a monad for your "language" for free, given that the parameter you give to `Free` is a functor in itself. And herein lies the problem. You have to write an instance of `Functor` and `Free` is not so free anymore.

Lets look at a mini example.

    data MiniIoF a = Terminate
                   | PrintLine String a
                   | ReadLine (String -> a)

    type MiniIoF = Free MiniIoF

This describes a mini language for doing IO operations. You can run this by pattern matching but I digress. The important thing to notice here is the type parameter `a` describing the type of the "remaining" or continuation if you will. And if we want to do something with this remaining, say append an action to the program, we need to be able to `fmap` over `MiniIoF` to access it. This is why we need the `Functor` instance. 

But the instance itself turns out to be rather boring to write.

    instance Functor MiniIoF where
        _ `fmap` Terminate = Terminate
        f `fmap` PrintLine s a = PrintLine s (f a)
        f `fmap` ReadLine g = ReadLine (f . g)

We are just mechanically taking it apart, putting `f` at the right places, and putting it back together. Considering the functor laws this is **the only valid instance** (that is not throwing errors or doing something funny). So why can't compiler write it for us?

It turns out that GHC can write it for us. For this to work we just need to enable the `DeriveFunctor` extension.

    {-# LANGUAGE DeriveFunctor #-}
    data MiniIoF a = Terminate
                  | PrintLine String a
                  | ReadLine (String -> a)
                  deriving Functor

And as if by magic we have our functor.


## Foldable, Traversable

I already mentioned containers when talking about functors. Apart from mapping the two common general operations over containers are traversal and folding. So if we are writing custom containers (or some different data structure) it is nice to provide instances for `Traversable` and `Foldable`. Let's take a look at a reimplementation of lists and their instances

    data List a = Nil | Cons a (List a) deriving (Eq, Show, Functor)

    instance Foldable List where
      foldr _ z Nil = z
      foldr f z (Cons x xs) = f x $ foldr f z xs

    instance Traversable List where
      traverse _ Nil = pure Nil
      traverse f (Cons x xs) = Cons <$> f x <*> traverse f xs

So now we can fold over our lists

    λ: foldr (++) "" $ Cons "1" $ Cons "2" $ Cons "3" $ Nil
    "123"

But looking at these instances you might notice they are pretty mechanical definitions. If you look at the types they pretty much write themselves. And indeed they can write themselves using some GHC trickery. Namely extensions `DeriveFoldable` and `DeriveTraversable` respectively

    {-# LANGUAGE DeriveFoldable #-}
    {-# LANGUAGE DeriveTraversable #-}

    data List a = Nil | Cons a (List a) 
                  deriving (Eq, Show, Functor, Foldable, Traversable)


## DataTypeable

There are some type classes that aren't just boring to implement but you are even discouraged from writing your instances manually. These two culprits are `Data` and `Typeable`. `Data` is a type class for abstracting over "data shape" allowing you to write generic code for any user-defined ADT without template Haskeltresel. This is something you really don't want to write by hand since it is not only terribly boring but also error prone. The `Data` type class is also a subclass of `Typeable` so the two are often defined in pair. `Typeable` defines a way for generating a value representative of a type. This allows type-casing and other "evil" things at runtime.  But it is sometimes necessary for writing efficient or generic code. It generates 128 bit fingerprint for a given type and nests fingerprints for type parameters if there are any. Since these fingerprints are supposed to be unique you must be familiar with implementation details of other instances (e.g. those provided in `base`) in order not to clash. It might be a good idea to let generation of this code to GHC using the extension `DeriveDataTypeable`

Lets extend the `List` example from before with two more derived instances

    {-# LANGUAGE DeriveFoldable #-}
    {-# LANGUAGE DeriveTraversable #-}
    {-# LANGUAGE DeriveDataTypeable #-}

    data List a = Nil | Cons a (List a) 
                  deriving ( Eq, Show
                           , Functor, Foldable, Traversable
                           , Typeable, Data)

## Newtypes

Sometimes you want to treat some values of a given type differently. And this is where newtypes come in. Just wrap it up. It doesn't cost anything at runtime. But it might cost you something at code-writing-time. Let's say you have some monad transformer stack for your application and you wrap it up in a newtype. 

    newtype App a = App { unApp :: ReaderT Config (StateT AppState IO) a }

But `App` is not a monad anymore. Nor is it `MonadReader`, `MonadState` or `MonadIO` if you are using mtl. You have to code directly against the specific type. You lost the ability to run programs polymorphic in the monad type against your application monad stack. In order to get it back you have to write all the instances. But they will be extremely boring since all you have to do is lift the operations. And it gets worse! Now you are creating overhead and hoping that compiler is smart enough to optimize it away (newtypes are guaranteed to disappear). The alternative is to just use a type alias but this can means giving up some precision in our types.

Or you can just ask the compiler to generate all the instances for you. And the GHC doesn't have to play by its own rules. It can just take the existing instances and insert some (free at runtime) coercions to make them work with newtypes. Thi s is the `GeneralizedNewtypeDeriving` extension.

    {-# LANGUAGE GeneralizedNewtypeDeriving #-}
    newtype App a = App { unApp :: ReaderT Config (StateT AppState IO) a }
                    deriving (Monad, MonadReader Config, 
                              MonadState AppState, MonadIO)

And we can happily use out newtypes and still have all the polymorphism we want.