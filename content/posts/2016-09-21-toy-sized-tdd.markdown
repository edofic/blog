---
title: TDD-ing a toy sized project
date: 2016-09-21
date: 2016-09-21
---

Just recently I was porting a toy sized parser combinator library (a proof of
concept) from Haskell to Python. You know, for educational purposes. It turns
out I'm not smart enough to keep the complicated types (and explicit laziness)
in my head even for such a small project. So my solution was to do TDD. To
clarify: I wanted to test happy paths through my functions to make sure at
least types fit together.

So tests I will write. But that means I need another file for tests (yes, it
was I single file project, that's what I mean with "toy-sized"), some actual
structure, a test runner, probably a package description for dependencies...
How about no? I remembered seeing a thing called `doctest` once. Turns out it
fit perfectly!.

This is how using `doctest` looks in Python:

```python
def parseString(target):
  """
  >>> parseString('foo')('foobar')
  [('foo', 'bar')]
  """
  def p(s):
    if s.startswith(target):
      return [(target, s[len(target):])]
    else:
      return []
  return p


if __name__ == '__main__':
  import doctest
  doctest.testmod()
```

And then run with `python myfile.py`. It is that simple (and built-int). You
just put examples in the docstring and they will be machine checked.

But I'm lazy and don't want to run tests every time by hand. I want a
poor-man's runner with `--watch` capability. And I can have it as a bash
one-liner (given that `inotify-tools` package is installed on my system).

```bash
while true; do
  inotifywait myfile.py
  clear
  python myfile.py
done
```

This will automatically clear the screen and re-run tests every time I change
the file - save it from my editor. Now I can finally develop my toy projects in
split screen with terminal and vim using only my editor's save function to run
tests :)


## Bonus section

Doctest is ported to Haskell, so I can use it there as well. Just need to
`stack install` it globally and it will be available for my one-file projects.

```haskell
-- |
-- >>> parseNumber "123"
-- [(Number 123,""),(Number 12,"3"),(Number 1,"23")]
parseNumber :: Parser Expr
parseNumber = mapParser (Number . read . concat) $ plus parseDigit
```

And run with

```bash
while true; do
  inotifywait myfile.py
  clear
  doctest myfile.py
done
```
