---
title: Approaches to designing a Haskell API
---

Recently I've been thinking about the design of programming interfaces, especially in Haskell. But don't let the title misguide you; this is not supposed to be a tutorial or a guide but simply an showcase of different styles. Feel free to tell me I'm wrong or missed something.

## The problem
Let's say we are writing an interface to RESTful web service. Our goal is to create type safe functions and descriptive models but all in all easy to use. 

The service should be simple so our examples are kept small. So let's have a single resource that supports `POST` and `GET` on single items. In pure Haskell it would look like

```haskell
newtype Id = Id Integer
data Record = Record { ... }

get :: Id -> Record
create :: Record -> Id
```

But this functions couldn't possibly be pure since they should be talking to our service, so we have to adapt this interface.


## The OO way
Our functions should at the very least take some sort of a client as an input. Let's define a client that does raw JSON requests to our service. This is just a rough sketch so we have something to work with

```haskell
type Client = Method -> Path -> Maybe RequestBody -> IO Response
```

Our API would now look like

```haskell
get :: Client -> Id -> IO Record
get client id = fromJSON <$> client "GET" (show id) Nothing

create :: Client -> Record -> IO Id
create client record = 
  fromJSON <$> client "POST" "" $ Just $ toJSON record
```

Which is actually not bad. We presume our service (or the network) will never fail, we push the burden of managing the client to the user and we clutter every function's signature with `Client` but we could do worse. On the up side we can now use this in a object-oriented-looking way.

```haskell
client `get` someRecordId
```

## Monads
Everything is better with monads right? So we can define a monad to replace the client. What's the essence of the client? It carries some configuration and has the ability to perform IO. So we build a monad that can do these two things. And since I'm lazy I'll just do

```haskell
type ClientM = ReaderT ClientConfig IO
```

But there is a problem. It is impossible to write an IO transformer so this new monad will have to be the base of the monad transformer stack. And this hurts modularity. Can we do better? Yes! Let's just do a monad class

```haskell
class Monad m => MonadClient m where
  rawRequest :: Method -> Path -> Maybe RequestBody -> m Response
```

and we can provide a way to run this again with `Reader` and `IO`.

```haskell
instance (MonadReader ClientConfig m, MonadIO m) => MonadClient m where
  rawRequest method path body = do
    conf <- ask
    liftIO $ ...
```

But now we can do better. We can write an instance for running tests that doesn't perform IO but instead simulates the service locally.

This is somewhat lighter for the user since he doesn't need to manage the client any more neither does he need to use the client for every request. He just need to put configuration and IO into his monad stack. 

## mtl?
Can we push the monad class approach further? Let's take a look at the [mtl](http://hackage.haskell.org/package/mtl) library. It provides classes of operations for every monad. We can write a class for our whole service.

```haskell
class MonadClient m where
  get :: Id -> m Record
  create :: Record -> m Id
```
And now implement this in terms of our first (naive) implementation which is not directly exported any more. User only sees the class and gets an instance that will work with the configuration and ability to perform IO. This now makes it very easy to write an instance that uses `State` (or something similar) to simulate the service for testing without bothering with JSON and other implementation details of the actual service. 

But it also has a downside. Imagine a bigger service with multiple resources. The class will explode and become unwieldy. Also making the implementation harder. If the implementor wants to support another backend he now has a big instance to write. If he wants to add a method he now has several instances to fix. Sound a bit like the [expression problem](http://en.wikipedia.org/wiki/Expression_problem).   


## Purity to the rescue
A good way to structure your code is to separate pure functions from impure actions and minimize the latter. 

One way to do this is to make a pure *description* of requests and responses. Then define a uniform intermediate representation that works well with our protocol and the client that actually does the requests. Only the client needs to perform any effects, all other code can now be pure.

```haskell
data Request = RequestGet Id | RequestCreate Record
data Resp = ResponseGet Record | ResponseCreate Id

data Raw = Raw Method Path (Maybe RequestBody)

toRaw :: Request -> Raw
toRaw (RequestGet id) = Raw "GET" (show id) Nothing
toRaw (RequestCreate record) = Raw "POST" "" $ Just $ toJSON record

fromResponse :: Response -> Resp
...

client :: Raw -> IO Response
...

request :: Request -> IO Resp
request r = fromResponse <$> client (toRaw r)
```

But this is terrible. First of all the `fromResponse` function is hard to write. But most importantly the `request` function is horribly unsafe. We cannot ensure the response will match the request. 

## GADTs
To ensure this we need generalised algebraic datatypes. 

```haskell
data Request resp where
  RequestGet :: Id -> Request Record
  RequestCreate :: Record -> Request Id

data Raw a = Raw Method Path (Maybe RequestBody) (JSON -> a)

toRaw :: Request a -> Raw a
toRaw (RequestGet id) = Raw "GET" (show id) Nothing fromJSON
toRaw (RequestCreate record) = Raw "POST" "" (Just $ toJSON record ) fromJSON

client :: Raw -> IO Response
...

request :: Request a -> IO a
request r = parse <$> client raw where
  raw@(Raw _ _ _ parse) = toRaw 
          
```

There are two things going on here. First we encode the type of response into the request so `request` can be safe. Second we use that type information that is locally available inside `toRaw` to pick the right instance for decoding JSON and put the specialised function into the raw representation. 

We now have it all: safety, modularity (we can write tests in terms of pure requests and responses), we can simply plug a new backend and even explicitly talk about requests since they are just plain old values. 

But we cannot statically determine the type of the request nor can we simply add a new type of request. Former being a philosophical remark and latter a real world requirement. We've again hit the expression problem. If we add a request we need to modify existing code in all functions creating an intermediate representation (or working directly with requests). At least adding a new backed is very simple since it only depends on the intermediate representation. 

## Type classes revisited
We want to be able to statically enforce types of requests. This is simply achieved if we define a single constructor type for every request instead of a sum type of all requests. But now we cannot have a function to convert it into an intermediate form. But we can have a type class. Moreover using multi param type classes and functional dependencies we can encode the type of result for each request and require the instances to parse the result from the intermediate form. Functional dependencies will ensure we can always compute the result type from the request type.

```haskell
data Get = Get Id
data Create = Create Record

data Raw = Raw Method Path (Maybe RequestBody)

class RequestRaw req resp where
  toRaw :: req -> Raw
  fromResponse :: Response -> resp

instance RequestRaw Get Record where
  toRaw (Get id) = Raw "GET" (show id) Nothing 
  fromResponse = fromJSON

instance RequestRaw Create Id where
  toRaw (Create record) = Raw "POST" "" (Just $ toJSON record )
  fromResponse = fromJSON

client :: Raw -> IO Response
...

request :: RequestRaw req resp => req -> IO resop
request r = fromResponse <$> client (toRaw r)
```

I believe we achieved our goal. We can add a new request without modifying existing code by simply adding new instances. And we can still add new backends that only rely on the intermediate form. We still can have pure tests and as a bonus big APIs will not require giant functions anymore, we can even break them up into several modules.

## Conclusion
I would argue that each of these approaches (except unsafe ADTs) has its pros and cons and therefore its place in some implementation. If I missed anything or made an error please let me know - I'll be happy to update the post.