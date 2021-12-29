---
title: Speed is a feature; thoughts on modern web performance
date: 2021-12-29
---

I'm pretty sure that as computers got faster they got slower. As chips got
faster the software slowed down. I'm acutely aware that we do much more now but
just the user interfaces really have no reason for slowing down.

Why do I care? Because a fast web is a delight to browse. A fast app is a
delight to use. The tool feels like an extension of your body. Where as a
sluggish tool get's in the way. Feels like you're battling it instead of using
it. This is one of the main reasons I still prefer vim over more "modern" IDEs.

Turns out it's not just a feeling. Dan Luu [measured input lag on different devices](https://danluu.com/input-lag/) and some old ones are much faster.

But that's not exactly what I have in mind. For a while I've been working on
backends for predominantly web applications, so the web is where my nagging
feeling points at. And there it roughly breaks down into 3 problems

- we're serving too much
- we're serving too slow
- we're doing too much on the client


# We're serving too much

According to a [random SEO site](https://www.seoptimer.com/blog/webpage-size/)

> Web page sizes have been growing steadily over the years. The first web page on the internet was only 4 KB in size. This was particularly because browsers back then did not support a lot of things that they do today. For example, it was not until 1993 that browsers began to support images. According to httparchive the average web page size in 2017 was 3 MB. This is a huge increase from the 1.6 MB average of 2014. It is predicted that the average page size will be 4 MB by 2019.

The irony: in 2021 (when I' writing this) that very post weighs in at 5.3 MB.
And takes about 5s to load on my (temporary) "not great, not terrible"
connection.

Pingdom did [an
analysis](https://www.pingdom.com/blog/webpages-are-getting-larger-every-year-and-heres-why-it-matters/)
of the top 1000 sites. The mean is just over 2MB. The catch: it's from 2018, my
bet is everything blew up significantly since then.

And my hands on experience: I've been annoyed at slow page loads at work and
I've looked at what we're loading. Turns out that our SPA comes in at about 7MB
(yikes :grimacing:) and loading this over 10mbps connection takes a while. Turns
out most of it is cacheable. And by setting the right headers you can get that
down to ~500kb and load times improve 10x. Great, right? I guess this works if
you have a heavy site with repeat users where caching can actually kick in.

But there is still the problem of heavy sites. The above "random SEO site"
shrinks for a whole whopping megabyte with adblocking. That is also a second to
load for me, or about 20% improvement in page load. Let's look at the rest.

- 127 requests
    - 83 of which to the same domain
    - 44 to other domains
- 1.6Mb of CSS
- 2.1Mb of JS
- 1.2Mb of images

And I suppose this is considered normal nowadays. I'm the weirdo who optimizes
his blog for size (last post comes in at ~900kb, most of it images, which
reminds me I should probably set up image compression pipeline).

```sh
566K Jan  4 22:03 angular2.min.js
563K Jan  4 22:05 angular2.0.0-beta.0-all.umd.min.js
486K Jan  4 21:50 ember.1.13.8.min.js
435K Jan  4 21:48 ember.2.2.0.min.js
205K Jan  4 22:06 angular2.0.0-beta.0-Rx.min.js
144K Jan  4 21:59 react-with-addons-0.14.5.min.js
143K Jan  4 21:46 angular.1.4.5.min.js
132K Jan  4 21:56 react-0.14.5.min.js
121K Jan  4 21:35 angular.1.3.2.min.js
5.3K Jan  4 22:00 redux-3.0.5.min.js
706B Jan  4 21:57 react-dom-0.14.5.min.js
63K  Oct 13 03:02 vue-2.0.3.min.js
```

Add in other utility libraries and you're probably loading a megabyte of JS (if
you're not focusing on the bundle size) before you even start writing your
application.

I have some thoughts on technical solutions but first repeat after me: you
cannot solve people problems with technology. The problem is not the technology
the problem is that speed (and bundle sizes) are ignored and never a priority.
Developers having super fast machines and speedy connections don't help. But
we're probably so far gone that nobody bats an eye at 5s load times for simple
pages anymore. Even more so with SPAs with the "it's only on the first load"
apologism.

# We're serving too slow

How long does it take your server to respond to a request? For static assets?
For an API call (if you're doing an SPA)?

More importantly: do you even care? Do you have SLOs in place? Or at least some
metrics?

If you're working at a FAANG on some microservices then you're confused why
someone would not care about this...but the majority of the internet is not
FAANGs. And a random website does quite poorly. Let's look at our random
website again (I promise I'm not picking on them I just think it's quite
representative).

- For some reference first
  - I'm on a 10mbps connection
  - about 30ms away from any servers (local network + ISP)
  - about 210ms away (ping) from the server of the site of interest
  - initial connection needs two roundtrips so the baseline is about ~420ms
- base html
  - 22kb, so transfer is not a real bottleneck
  - 2.1s with curl
  - 1.8s in a browser that reuses old connections
  - so around 1.5s is actually spent by the server
- static assets
  - around 200-300ms spent by the server
  - even for assets <1kb, higher for larger

Before criticising too much, am I constructing a strawman here? The site I'm
exploring is based on Wordpress so I'd say this is actually representative of a
big part of the internet.

Let's point the loupe inward - how do I fare? Numbers for an article from this
blog (static content)

- Reference
  - Same network connection (10mbps, 30ms away from civilization)
  - blog is behind Cloudflare CDN which is also about 30ms away
- base html
  - 9kb
  - ~230ms with curl
  - 150ms from a cold browser
  - 70ms with existing connections
  - ~50ms for small assets
  - ~200ms for images (transfer time becomes a factor)
  - caveat: uncached resources hit the origin server and incur about 100ms
    penalty so 250/170ms in the browser

I'm not too happy with the absolute numbers - for some contrast a raw nginx
takes <5ms to serve a small asset. So 35ms ought to be doable. But since I don't
wan to build and maintain my own CDN I'm settling for current numbers.

But I still think going into seconds makes for a poor experience.

A common approach nowadays is to build an SPA that loads data via async API
requests. This is now you're time to useful information on screen. Good numbers
are hard to get by here but from what first hand experience I have I can say
response times can quickly creep up into seconds and nobody bats and eye since
there is a spinner there and it's normal that things take a while to
load....right?

Again there are a plethora of technical solutions but first you need to care. If
you don't continuously measure response times and act when they creep up you're
agreeing to slow responses and thus slow user experience. Again it's a question
of priority.

I even think you can deliver a decent experience on pretty much any stack (no
need to rewrite in rust :wink:) it just comes down to paying some attention and
taking performance into account when making day-to-day decisions.

# We're doing too much on the client

I'll admit this section is most hand-wavy because I don't do much frontend work
but I still think it's (somewhat) valuable.

This is specifically a problem of fat clients - single page applications. A
typical flow there goes something like this

- you click a button
- javascript makes a request to the server and shows a spinner
- server fetches the data and responds with json
- browser parses the json
- the client application manipulates this response somehow
- javascript renders templates/views/components based on this manipulated
  response, often into a virtual dom
- view is applied to actual DOM and displayed to the user

The problem here is that all of these steps take some time and they all need to
run sequentially. So small inefficiencies pile up. Yes js is "fast". And so is
the browser parsing the json and so is react rendering. But if you need to do
multiple requests and aggregate data (common) or you return more data than
actually ends up on the screen (very common) you're wasting precious time at
every step.

But premature optimization is the root of all evil!!! Well the full quote from
Knuth is

> Programmers waste enormous amounts of time thinking about, or worrying about, the speed of noncritical parts of their programs, and these attempts at efficiency actually have a strong negative impact when debugging and maintenance are considered. We should forget about small efficiencies, say about 97% of the time: premature optimization is the root of all evil. Yet we should not pass up our opportunities in that critical 3%.

And I don't dispute. My only claim is that the number is much larger than 3% in
"modern" layered architectures. I think fat clients are still subject to all
problems form above but the situation is exacerbated by the fact there is
slowness on all layers so the effects compound.

# Conclusion

Speed is a feature. You have to work to get it. But it's usually neglected and
downprioritized into oblivion. It's usually not the technology. Usually it's
culture, it's accepting that this is fine.

What are we to do? As I said, I don't think there is a technological silver
bullet that can solve this. And I'm also not a proponent of "more discipline
will solve things" - I don't believe this can scale to larger organisations.

Should we abandon all hope? I don't think so. I think the right approach is to
be the voice of technology, a cheerleader for delightful snappy applications.
Advocate for performance not for performance's sake but for users. A product
that is snappier is more delightful to use and makes for happier customers (you
know NPS and all that jazz).

Just care a little bit and don't be shy to talk about it. With some luck you may
infect one or two more people. This way we can nudge the needle just a little
bit. It's not too late.
