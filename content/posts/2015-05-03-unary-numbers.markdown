---
title: Lazy unary numbers
date: 2015-05-03
---

We are used to encoding numbers on computers in binary. Binary is the "simplest" base that yields logarithmic length though it [may not be optimal](http://web.williams.edu/Mathematics/sjmiller/public_html/105Sp10/addcomments/Hayes_ThirdBase.htm). But can we do simpler? How about unary?

Unary is often used with Turing machines where we don't care for efficiency and I will assume this same stance. Let's forget about efficiency and explore what can unary numbers do that binary can't. Specifically lazy unary numbers as otherwise the systems are equivalent. I'll be using Haskell as it is lazy by default and thus a good fit.

Let's start off with a simple definition.

```haskell
data Nat = Zero | Succ Nat deriving Eq
```

This are natural  numbers as usually defined in mathematics. An example value (3) in this encoding would be `Succ $ Succ $ Succ Zero`.  For the sake of simplicity I will from now on assume that we do not use bottom values that correspond to errors.

## Infinity

The first interesting property of this representation is a simple encoding of infinity in finite space.

```haskell
nInf :: Nat
nInf = Succ nInf
```

This ties the knot and only constructs one `Succ` that points to itself. The interesting (and useless) bit is that we cannot observe this in a given value! At least not in pure functions. Given a `Nat` we cannot possibly know if it is infinite or just very large.

## Equality and ordering

I just derived `Eq` before not giving it much thought. But notice now that comparison between two `Nat`s may not return if both are infinite. However if one is finite we can detect that the two are not equal. All is not lost. We can even do more and define ordering

```haskell
instance Ord Nat where
  Zero `compare` Zero = EQ
  Zero `compare` Succ _ = LT
  Succ _ `compare` Zero = GT
  Succ n `compare` Succ m = n `compare` m
```

This does the obvious thing: it "zips" the two numbers together and finds a differing spot if there is any. Again it will loop forever if both numbers are infinite.

## More instances

This type is also enumerable which can be implemented by an isomorphisim with integer
via an `Enum` instance

```haskell
instance Enum Nat where
  toEnum 0 = Zero
  toEnum n | n > 0 = Succ $ toEnum $ n - 1

  fromEnum Zero = 0
  fromEnum (Succ n) = fromEnum n + 1
```

Given this we can now show `Nat` as a decimal number

```haskell
instance Show Nat where
  show n = show $ fromEnum n
```

And I saved the best for last: full `Num` instance

```haskell
instance Num Nat where
  Zero + m = m
  Succ n + m = Succ $ n + m

  n * Zero = Zero
  Zero * m = Zero
  n * Succ m = n + n * m

  Zero - _ = Zero
  n - Zero = n
  Succ n - Succ m  = n - m

  fromInteger 0 = Zero
  fromInteger n | n < 0     = fromInteger $ negate n
                | otherwise = Succ $ fromInteger $ n - 1

  negate = error "cannot represent negative Nat"
  abs = id
  signum _ = Succ Zero
```

Let's walk through it. Addition recurses on first argument. Note that the second argument is reused. It will only allocate `Succ`s matching the first argument.

Multiplication again recurses on one argument but does addition at each step. Since there is no way to represent negative numbers in this scheme I defined `0 - n = 0` which is a bit shady but works in most cases. Similarly `negate` throws an error. Anther consequence is that `abs` is just identity and `signum` always returns `1`.

But the most useful function is `fromInteger`. I extended `toEnum` by handling negative cases. This does not look like much but due to literal polymorphism we can now write decimal literals and they will be automatically converted to `Nat` where this is the expected type.

## Conclusion

There are two interesting things. First is encoding infinity. Not much on itself. But the second thing is partial evaluation. By traversing *n* `Succ`s we know that the number is greater or equal to *n*. This means we can compute even with infinite numbers as long as we don't need to look at the exact result.
