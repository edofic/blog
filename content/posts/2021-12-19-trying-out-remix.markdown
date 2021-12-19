---
title: "Trying Out Remix"
date: 2021-12-19
tags: ["trying-out"]
---

I came across [Supporting Remix with full stack Cloudflare Pages](
https://blog.cloudflare.com/remix-on-cloudflare-pages/) on Cloudflare blog (yes
I'm a CF fan) and [Remix](https://remix.run) piqued my interest. Mostly because
it addresses my concerns with "modern" web development: the wiring boilerplate
is mostly gone, there is no 100 file setup to get a workable dev environment. And
of course because it's outside my comfort zone - so maybe I'll learn something.

# Getting started

There's a code example right on the home page but I'm having hard time building
a mental model from just that. Luckily there is also a bit [Get started](https://remix.run/docs/en/v1/tutorials/blog) button which takes me to a tutorial. Great, let's follow that.

![following initial instructions](/images/remix/1.png)

Few more commands
```bash
cd my-remix-app
npm run dev
```
And we're up.

![localhost shows dmeo app](/images/remix/2.png)

Then there is text that really speaks to me

>  If you want, take a minute and poke around the starter template, there's a
>  lot of information in there.

Oh yes I intend to do that :smile:

Let's make a git repo to keep track of changes. Conveniently there is aready a
`gitignore`.

```bash
git init
git add .
git commit -m "Initial commit"
```

Let's now see what we have

```bash
 git ls-files
.gitignore
README.md
app/entry.client.tsx
app/entry.server.tsx
app/root.tsx
app/routes/index.tsx
package-lock.json
package.json
public/favicon.ico
remix.config.js
remix.env.d.ts
tsconfig.json
```

Honestly less than I expected. And I mean this in the best possible way. I like
environments which I can understand to the point of being able to write manually
from scratch. Yest I'm looking at you `create-react-app` and other humongous
frontend toolchains. I feel at home with Go where you have your dependencies
file and a single `go` command.

Back to source files; there's a typescript config but no webpack config (yay?)
or similar. This is in the typescript config

```bash
// Remix takes care of building everything in `remix build`.
```

Great!  Let's find the entrypoint now to better understand this.

Since I've been instructed to use `npm run dev` this means that `package.json`
will define the action.

```json
  "scripts": {
    "build": "remix build",
    "dev": "remix dev",
    "postinstall": "remix setup node",
    "start": "remix-serve build"
  },
```

And it does, no magic here.

There is something the looks like remix configuration (`remix.config.js`)

```js
module.exports = {
  appDirectory: "app",
  assetsBuildDirectory: "public/build",
  publicPath: "/build/",
  serverBuildDirectory: "build",
  devServerPort: 8002,
  ignoredRouteFiles: [".*"]
};
```

And apparently it points to `app` directory. Conveniently there I can find
`entry.client.tsx` and `entry.server.tsx` which I guess are my entrypoints. What
I found slightly strange is `.tsx` (which I believe is JSX for typescript) for
the *server*.


Maybe time to read some [more
docs](https://remix.run/docs/en/v1/tutorials/blog#your-first-route) :shrug:

# Baby's first code

![localhost shows dmeo app](/images/remix/3.png)

Confused for a moment...apparently docs are not up to date, the links are now in
`app/routes/index.tsx`. Oh well, I'll manage. I added the `li` as instructed

```html
<li>
  <Link to="/posts">Posts</Link>
</li>
```

and then of course it fails to compile as `Link` is not defined. Guessed the
import based on imports in `root.tsx` to be

```ts
import { Link } from "remix";
```
It works!

![link to posts](/images/remix/4.png)

Then created `app/routes/posts/index.tsx` as

```jsx
export default function Posts() {
  return (
    <div>
      <h1>Posts</h1>
    </div>
  );
}
```

And non-surprisingly it shows up now. Time for some remix magic now.

# Loaders

> If your web dev background is primarily in the last few years, you're probably used to creating two things here: an API route to provide data and a frontend component that consumes it. In Remix your frontend component is also its own API route and it already knows how to talk to itself on the server from the browser. That is, you don't have to fetch it.

```jsx
import { Link, useLoaderData } from "remix";

export const loader = () => {
  return [
    {
      slug: "my-first-post",
      title: "My First Post"
    },
    {
      slug: "90s-mixtape",
      title: "A Mixtape I Made Just For You"
    }
  ];
};

export default function Posts() {
  const posts = useLoaderData();
  return (
    <div>
      <h1>Posts</h1>
      <ul>
        {posts.map(post => (
          <li key={post.slug}>
            <Link to={post.slug}>{post.title}</Link>
          </li>
        ))}
      </ul>
    </div>
  );
}
```

And it show up as expected.

![link to posts](/images/remix/5.png)

Then the tutorial moves on to refactoring, extracting data and fetching from an
external data source but I'm more interested in this mechanics as it seems to me
this is the meat and the potatoes of remix. So what is going on?

Well apparently the `useLoaderData` hook is automatically tied to corresponding
exported `loader` in the same route. By convention instead of us manually wiring
up routing, neat. But I have questions now :smile: Can I pass parameters? Is
there `useQuery`-like magical caching and deduplication? Authentication? CORS?

More interestingly: if I looks at network in devtools I don't see an XHR
request, meaning our data is pre-rendered on the server. Super neat! But if i
navigate from index to posts I see a request to
`http://localhost:3000/posts?_data=routes%2Fposts%2Findex` which directly
returns my data.

# Routes

[Reading on](https://remix.run/docs/en/v1/tutorials/blog#dynamic-route-params)
the tutorial actually explains how to parametrize aka how to do `dynamic route
params`.

It's again done by convention, not explicit config: file name is a placeholder
for a parameter. Clever?

I'll admit my knee jerk reaction to this is "ugh, ugly. I want my router". But
it may work? I don't know, would actually need to do a larger project to
evaluate properly. In any case there is a way to get an overview of the routes:

```jsx
$ remix routes
<Routes>
  <Route file="root.tsx">
    <Route path="posts/:slug" file="routes/posts/$slug.tsx" />
    <Route path="posts" index file="routes/posts/index.tsx" />
    <Route index file="routes/index.tsx" />
  </Route>
</Routes>
```

So I can still have my cake (easily overview and discover routes) and eat it
(not write the router) too.

Ok, so how does this paramterized component look like?

```jsx
import { useLoaderData } from "remix";

export const loader = async ({ params }) => {
  return params.slug;
};

export default function PostSlug() {
  const slug = useLoaderData();
  return (
    <div>
      <h1>Some Post: {slug}</h1>
    </div>
  );
}
```

And sure enough I can now click on one of the posts and it renders

![rendered post](/images/remix/6.png)

So what is happening here?

* My url path `/posts/my-first-post` is matched against the `posts/:slug` route
  and my component is rendered
* then the `useLoaderData` hook fires and data for self url is requested -
  more concretely when navigating on frontend a request to
  `posts/my-first-post?_data=routes%2Fposts%2F%24slug` is dispatched (I'm
  assuming this is bypassed when server-side rendered and loader is called
  directly)
* server side matches this against routes again, picks up my page but since
  we're requesting data it now calls the `loader` with query path parameters
  passed in as named values in `params`. The basic example is just returning a
  string here but it could be any json, probably with data from some external
  data source as well
* response from the loader is wired (via fetch/XHR or server side rendering) to
  the hook and the string returned from the loader is bound to `slug` constant.

Now the analogy from the tutorial that loader is the controller and we're using
react as a view layer makes sense :grinning:

Then I followed the instructions to get markdown rendering up an running which
is well covered in the tutorial and standard javascript so I'll skip the
details.

# An experiment

But meanwhile I started wondering...can I put multiple components on screen and
each will automagically fetch it's data? I could peruse the docs some more...or
I can just try it out.

I created this `app/routes/footer.tsx`


```tsx
import { useLoaderData } from "remix";

export const loader = async () => {
  return new Date().getFullYear();
};

export default function Footer() {
  const year = useLoaderData();
  return (
    <div>
      All rights reserved © { year }
    </div>
  )
}
```

Yes, very silly, calling a server to get current year, but I just want to test
things out :sweat_smile:

One remark while I figure things out: tooling is quite responsive and helpful:
change detection, auto-recompilation and reloading out of the box.

```
Build failed with 1 error:
route-module:/Users/andrazbajt/personal/playground/my-remix-app/app/routes/posts/$slug.tsx:3:9: error: No matching export in "app/routes/footer.tsx" for import "Footer"
?? Rebuilt in 43ms
GET /posts/ 200 - - 19.486 ms
?? File changed: app/routes/posts/$slug.tsx
?? Rebuilding...
?? Rebuilt in 124ms
GET /posts/ 200 - - 11.870 ms
```

Then in `$slug.tsg`

```tsx
<div>
  <div dangerouslySetInnerHTML={{ __html: post.html }} />
  <Footer/>
</div>
```

And a cryptic error appears

![application error](/images/remix/7.png)

What is slug doing here? Sprinkling in some logging I notice that `year`
actually holds data from `$slug.txs` loader, not my footer loader. So I'm
holding it wrong - back to the docs it is.

# Styles

Next interesting bit is that components and also include styles

```tsx
import adminStyles from "~/styles/admin.css";

export const links = () => {
  return [{ rel: "stylesheet", href: adminStyles }];
};
```

and this gets picked up by `Links` component in `index.tsx`.

# Index routes

But the real fun begins with `index routes`. By adding this to
`app/routes/admin.tsx`

```tsx
import { Outlet } from "remix";
...
<main>
  <Outlet/>
</main>
```

things can now render *inside* this component. E.g. visiting `/admin/new` will
render first `app/routes/admin.tsx` but inside the `Outlet` there will be
`/app/routes/admin/new.tsx`. This time with working data fetching :wink: But
does this mean I can only structure my data-fetching components hierarchically?

Reading on I'm slightly surprised to see that remix includes its own forms. How
do they differ from plain old html `<form>`? I'm guessing some automagical
wiring to the backend, maybe shared validation logic :fingers_crossed:

And sure enough, there is wiring by convention

```tsx
import { redirect, Form } from "remix";
import { createPost } from "~/post";

export const action = async ({ request }) => {
  const formData = await request.formData();

  const title = formData.get("title");
  const slug = formData.get("slug");
  const markdown = formData.get("markdown");

  await createPost({ title, slug, markdown });

  return redirect("/admin");
};
```

So what happens when I submit this form? Feels very smooth, let's look under the
cover.

* a `POST` request is made to
  `/admin/new?_data=routes%2Fadmin%2Fnew` with regular form data payload.
* but the response is interesting: status is set to `204 No Content` and there
  is an interesting header: `X-Remix-Redirect: /admin`
* apparently the "redirect" is then actually just frontend navigation to new
  route as there are no more requests except one for new data
  `/admin?_data=routes%2Fadmin`

# Other stuff

Then there is [another much longer
tutorial](https://remix.run/docs/en/v1/tutorials/jokes) (also available as [a
video](https://www.youtube.com/watch?v=hsIWJpuxNj0)) that dives deeper and
covers more topics (like cookies, authentication, validation, error handling,
databases...).

How about deployment? I did start this with the intent of deploying to
Cloudflare but the post is already getting long and the app as-is is not really
suitable (direct filesystem access is not supported) for deploying, so maybe in
another post.

# Conclusion

We'll as much as I'm not a home at writing frontend code and I dislike the idea
of "full stack javascript" this actually did not hurt a bit and I can see how
time can be really useful for doing a very dynamic client app that as some
server side components.

Now, I would probably not pick this stack for a new project, mostly to me not
being comfortable with Node, but I might just use it for an opinionated React
toolchain with SSR and all other goodies.
