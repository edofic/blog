---
title: Tailwind CSS vs Pico CSS
date: 2022-01-18
---

tl;dr it's a tradeoff so it depends. 

With that out of the ways let's start with the fact that this is not an article about technical merits and which is best. It about my experience and my beliefs about which kinds of projects are a good fit for which tool. 

So what did I do? I have a (useless) pet project I keep fiddling with to try out new things. At one point I decided my eyes hurt and I want to make the frontend a bit prettier. It is paramount to explain that I have indeed been living under a rock for a few years and haven't really touched any CSS. So I first did some exploration to see what has the world moved on to. 

Apparently nothing. To my surprise(?) [Bootstrap](https://getbootstrap.com/) is very much still a thing (it being the "thing" I used apart from in-house component libraries). To the point of there being [jokes about it's ubiquity](https://www.dagusa.com/). [Foundation](https://get.foundation/) is now big too. I vaguely remember it being the new kid on the block. Apparently that's [Bulma](https://bulma.io/) - makes me chuckle every time since I grew up with Dragon Ball :joy:

# Enter Tailwind

But the real hot piece of tech is not really a CSS framework as I remember them. Everybody is now talking about [Tailwind CSS](https://tailwindcss.com/) - a "utility-first CSS framework". In case you've been living under a rock (like me): it's a library that instead of giving you `btn-primary` gives you `bg-blue-600`. It pretty much looks like inline styles but with a tad more structure. An example from their homepage:

```html
<figure class="md:flex bg-slate-100 rounded-xl p-8 md:p-0 dark:bg-slate-800">
  <img class="w-24 h-24 md:w-48 md:h-auto md:rounded-none rounded-full mx-auto" src="/sarah-dayan.jpg" alt="" width="384" height="512">
  <div class="pt-6 md:p-8 text-center md:text-left space-y-4">
    <blockquote>
```

My first reaction was somewhere along the lines "this is beyond idiotic - WHY WOULD YOU DO THIS?!?!?" (I'm told this is a common reaction). So naturally I had to try it. I am a bit of a contrarian when it comes to software and I do believe a lot of good ideas are mis-labeled as stupid just because they are not understood when framed within current "best practices". Of course there still are actually bad ideas but given the buzz around Tailwind there may actually be something here. 

So how bad is it really? To style a button I did this

```html
class="bg-fuchsia-500 hover:bg-fuchsia-400 text-white font-bold rounded py-1 px-2"
```

So quite bad. I really don't want to repeat this all over the place. Unfortunately one important piece does not come across when just reading the examples: I don't need to.  The preferred way to do e.g. buttons is to create a button component in my framework-of-choice and just stick all the classes in there. Or if I really need it (e.g. multiple kinds of buttons?) I can in fact create a class that "inherits" all the properties like this

```
.btn {
    @apply bg-fuchsia-500 hover:bg-fuchsia-400 text-white font-bold rounded py-1 px-2;
}
```

and it's expanded in css processing. 

## The experience

The thing that surprised me a bit is that Tailwind will reset all styles by default. So e.g. `h1` and `h4` and `p` all look the same. You **have** to style each and every element. 

And there is no real _look_ out of the box. It can be what ever you need it to be. Super flexible abut this also means a lot more work to get started.

Tooling is very nice as well. Unused class pruning works out of the box and there is a neat VSC plugin with autocomplete that can show definitions (and colors!) inline. Though I'm slightly irritated that it's a asset pipeline step and that it requires Node as this makes it the most complicated part of my build.

All in all the workflow feels very productive as it's easy to get in flow - there is no switching between files, you just iterate until you get there. You do get used to the mess of classes and you can read it well (tooling also helps) but deep down I still think it's a bit messy. But since it's also productive...is this even a problem or just left-over thinking from a different time?

# In the other corner: Pico CSS

One day on Twitter I randomly see an ex-coworker [complaining about the same build step thing](https://twitter.com/anze3db/status/1480602792259694593?cxt=HHwWgsC54aOZlIwpAAAA)

> I'm removing  Tailwind CSS from a few of my projects.
>
> Not the best fit for small projects where it ended up being the most complex step in the build/release process 😅
> 
> I replaced it with picoCSS and even this feels like a bit of an overkill for what I'm using it.
>
> -- @anze3db

This resonated with me so I needed to look into it. [It's](https://picocss.com/) completely opposite approach:

> Elegant styles for all natives HTML elements without .classes

Where to do a button you write `<button>`. And it's 10kb and requires no build step? Next thing I know I'm restyling my toy project to see the difference. 

And...there wasn't much code to look at :open_mouth: I guess this makes sense, this is the whole point after all. I did need to tweak a few things to make my html match Pico's idea of _semantic_ and i did tweak the colors via css variables. But apart from that there isn't really much to code in sight (a tiny amount of classes sprinkled around).

At around this point I also realized that I have been living under a rock much larger than I realized. Minimal css frameworks are a such a numerous bunch now that there is even a [kitchen sink demo](https://github.com/dohliam/dropin-minimal-css) to preview different ones. A whole 102 of them at the time of my writing!

## The experience

Not much to say...it mostly gets out of the way. I did keep the kitchen sink open in open in one tab for reference on how to structure things but most of the time things just worked. 

So where is the downside? Well the tailwind version did look better. And this is to be expected as I styled it exactly how I wanted it to looks where with Pico (or any other minimal framework) I get a canned look (modulo some tweaking). So we're back to "Every Bootstrap Page" problem? I guess...but it's much less pronounced since there isn't just one big library that everybody uses and there are less "components" baked in. I can even imagine using different frameworks for different projects to get a different looks as the lock-in is truly minimal (as there is barely any code).

The problem does arise when you have a look you're going for (or a designer that gives you one). Then grabbing a canned looks and trying to make it fit starts to sound like a bad idea.

# Who wins?

We all do since we have options. Cheekiness aside I do believe both approaches have use cases. 

If you're a solo dev doing an mvp or even a team doing an internal tool and you just want to deliver ASAP and "just make it decent" is the design brief then I think a minimal framework is the right approach. By focusing on the content intead of the look you can get more done faster.

On the other hand if you have a concrete design (or even a desired look) then Tailwind starts to sounds like a much better approach. In my opinion it's not actually an alternative to minimal frameworks but instead an alternative to writing CSS from scratch. 
