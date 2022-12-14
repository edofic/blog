---
title: Hypermedia driven micro frontends
date: 2022-06-29
---

Micro frontends? Huh?

> extending the microservice idea to frontend development
> -- https://micro-frontends.org/

or: how to scale frontend development once you have too many engineers and they start stepping on each other's toes. Assuming of course your frontend is an SPA ([grug](https://grugbrain.dev/) might have something to say about both points).

But of course splitting an SPA into multiple smaller applications and having it all play together brings along a ton of additional complexity. Just like splitting a monolith into micro services :wink:.

The idea I want to share with this post is that it does not have to be this way.

> What are you trying to tell me? That I can dodge bullets?
>
> No, Neo. I'm trying to tell you that when you're ready, you won't have to.

# Hypermedia-driven application architecture

In case you're not familiar just read [the original essay](https://htmx.org/essays/hypermedia-driven-applications/). Here's the intro to get the flavour:

> The Hypermedia Driven Application (HDA) architecture is a new/old approach to building web applications. It combines the simplicity & flexibility of traditional Multi-Page Applications (MPAs) with the better user experience of Single-Page Applications (SPAs).

The particular implementation I like (see also [How I fell in love with low-js](/posts/2022-01-28-low-js/)) is [htmx](https://htmx.org/)

# Micro HDAs?

How do you split up an HDA(hypermedia-driven application)? Into multiple "micro HDA"s? No, in fact there is no need to split anything. Just like an SPA can call an API implemented by multiple services, an HDA can load hypermedia (think HTML over HTTPS) from a backend implemented by multiple services.

Just like you can put an API gateway in front of your services you can put....I don't know, we used to just call this reverse proxies. The point being: you can render your html wherever you want. It's all hypermedia. The **browser** already does the integration. In fact you can argue the original purpose of the web browser was/is to seamlessly integrate hypermedia from distinct services.

# But.....problems!

Surely this cannot really work? I mean support the complexities of a real world application. I'll admit I haven't tried it out on a real use case but pretty much whatever I throw at it seems to work. Browsers are really great nowadays.

## Styles

Ha, I still need to compile all my CSS to include in the main page! Not really. You can add both `link` and `style` tags to "components" (meaning async loaded HTML) and it will just work. Even get unloaded when you remove that part from the DOM. You only need to take care of scoping - making sure that component styles don't affect other components. But that's pretty easy by wrapping into a container and then using that in all selectors.

## Scripting

This approach does not really work for JS. Yes you can add `script` tags but scoping is harder and there's no real way to unload.

Solution is to embrace [Locality of Behavior](https://htmx.org/essays/locality-of-behaviour/). Just put your behavior next to the UI and the problem goes away. You might want to include e.g. [_hyperscript](https://hyperscript.org/) or [Alpine.js](https://alpinejs.dev/) in the main page - effectively standardising on the tooling - to get the most mileage out of this approach.

I'll admit that "standardising on tooling" kinda defeats the purpose of micro frontends but you can still keep to [pure JS](https://youmightnotneedjquery.com/) and still get to implement the server any way you see fit.

You can still pull in full fledged JS frameworks but at this point are you still doing HDA?

# An example

Now let me share a small (but complete) example of how this can work in practice to make things a bit more concrete.
Below is a short video of two instances of my app side by side to demo the real time functionality.

![demo video](/images/micro-frontends/demo.gif)

It's a "landing page" with a "subscribe" functionality. But there is also a (very simplistic) chat on the bottom right. You might have guessed it - chat is implemented separately.

## Implementation

Let's jump right in. For my ~~API gateway~~ reverse proxy I'm using [Caddy](https://caddyserver.com/) with the [following config](/files/micro-frontends/Caddyfile)

```Caddyfile
{
	http_port 3000
	auto_https off
}

localhost:3000 {
	reverse_proxy /chat* localhost:8080
	reverse_proxy * localhost:8000
}
```

This will listen on port 3000 and proxy requests that have prefix `/chat` to port 8080 and everything else to port 8000. Yes I'll be implementing two services: one main page (my "legacy monolith" :wink:) and a chat service (extracted microservice).

If you want to follow along you only need [Go](https://go.dev/). To run caddy:

```sh
go run github.com/caddyserver/caddy/v2/cmd/caddy@v2.5.1 \
  run -config=Caddyfile
```

### Landing page

Full source is [landing.go](/files/micro-frontends/landing.go) - run with `go run landing.go`. This is stdlib http server serving the main template and a the "subscribe" form POST with the partial response. I'm including [htmx](https://htmx.org/) for interactivity

```html
<script src="https://unpkg.com/htmx.org@1.7.0"></script>
```

and some basic styles

```html
<link rel="stylesheet" href="https://the.missing.style">
```

Subscribe form is pretty straightforward

```html
<form hx-get="/subscribe">
  <label>Email address
    <input name="email" type="email">
  </label>
  <button type="submit">Subscribe</button>
</form>
```

But the real party is the link that loads the chat


*NOTE* adding `hx-trigger="load"` here will automatically load chat when the page loads, effectively lazy loading this component.

```html
<a hx-get="/chat" hx-swap="outerHTML">Live chat</a>
```
Huh? Plain old htmx. Nothing fancy. That's right. Plain old HDA. But since my reverse proxy is routing this request to the chat service the magic can happen.

# Chat service

Full source: [chat.go](/files/micro-frontends/chat.go). Run with `go run chat.go`. There's `/` which provides a basic wrapper for development so `/chat` can return a partial and be used by the landing page service. That's it. This is all the magic. The rest of this section will explain the inner workings of the chat itself - skip ahead if you're not interested. o

Posting messages is a regular form

```html
<form hx-post="/chat/msg">
<input type="text" name="msg" autocomplete="off">
<button type="submit">Send</button>
<form>
```

The handler will just return an empty form back. But it will also dispatch the message to all listeners. This oneliner sets up server-sent-events listener and wires the events into the DOM

```html
<div hx-sse="connect:/chat/events
     swap:message"
     hx-swap="afterbegin">
</div>
```

On the server side this is is simple enough to implement without any libraries, just nee do set the content type

```go
w.Header().Set("Content-Type", "text/event-stream")
```

And write events as lines (and make sure to flush)

```go
fmt.Fprintf(w, "data: <div>%v</div>\n\n", msg)
if f, ok := w.(http.Flusher); ok {
    f.Flush()
}
```

# Conclusion

I hope that a full example in 3 files is evidence enough that there really is nothing to see here and things just work (and the title is a bit clickbaity - whoops).

That said I hope I still provided some useful insight. I do in fact see this pattern as useful if you're starting out with something simple, say a Rails/Django app, or maybe Wordpress or even a statically generated page and then figure out you want a piece of it that would really better be implemented as a separate service. HDA gives you an easy way out without implementing APIs and pulling in SPA frameworks. I do believe there is some value in here. But I would not go splitting up majestic HDA monoliths [for the sake of splitting them up](https://grugbrain.dev/#grug-on-microservices).
