---
title: Trying out Buffalo
date: 2021-02-07
tags: ["trying-out"]
---

I'm trying out a new concept: my "I should check this out" list of tools and technologies has grown quite a bit so I got
the idea that I could do short write ups of my experimentation and first impressions to give myself motivation to try
these things out. Will see how this turns out, hopefully this will turn out a series :)

Today I'll be trying out [Buffalo](https://gobuffalo.io/en/) - a Go web framework. I believe this one distinguishes
itself by focusing on developer productivity - a kind of "Rails for Go". But apparently it's mostly some glue and
tooling around a bunch of other established libraries. I very much like this approach. And yes, I agree that "just use
standard library" is not (always) the right approach. Sometimes you just want to bootstrap a small project quickly.

Let's see if it delivers on this promise.

_Disclaimer_ this is not intended to be a review or a tutorial, just my first impressions, trying to convey the flavor I
experience.

Final code is available [on Github](https://github.com/edofic/trying-out/tree/main/buffalo) if you want to read through
anything.

# Getting started

I start by working through the docs - [Installation](https://gobuffalo.io/en/docs/getting-started/installation/).

I need Go, or I need to update it at least (I use [NixOS](https://nixos.org/), following commands should work on most
Linux distributions if you use [Nix package manager](https://nixos.org/download.html).

```bash
nix-env -iA nixos.go
```

And some frontend stuff

```bash
nix-env -iA nixos.nodejs
nix-env -iA nixos.yarn
```

not bothering with sqlite at this point, I'll probably use a dockerized Postgres for development.

I'll install the buffalo tooling using the prebuilt Linux binary.

```bash
wget https://github.com/gobuffalo/buffalo/releases/download/v0.16.21/buffalo_0.16.21_Linux_x86_64.tar.gz
tar -xf buffalo_0.16.21_Linux_x86_64.tar.gz
mv ./buffalo ~/bin/
```

Does it work?

```bash
buffalo
```

I get some output. Yay!

Moving on to [Generating a new project](https://gobuffalo.io/en/docs/getting-started/new-project/)

```bash
buffalo new trying-out
```

Command suggests we open up the [generated readme](https://github.com/edofic/trying-out/blob/main/buffalo/README.md).
Apparently now we need to setup a database ourselves. let's go with dockerized Postgres

```bash
docker run --name trying_out  -e POSTGRES_PASSWORD=postgres -d 5432:5432 postgres
```

Following the readme I check `database.yml` and to my delight the defaults perfectly match my choices so I can just run
the tooling against my new Postgres server.

```bash
cd trying_out
vim database.yml
buffalo pop create -a
```

Let's see the schema

```
 $ docker exec -it trying_out psql -U postgres
psql (13.1 (Debian 13.1-1.pgdg100+1))
Type "help" for help.

postgres=# \l
                                       List of databases
          Name          |  Owner   | Encoding |  Collate   |   Ctype    |   Access privileges
------------------------+----------+----------+------------+------------+-----------------------
 postgres               | postgres | UTF8     | en_US.utf8 | en_US.utf8 |
 template0              | postgres | UTF8     | en_US.utf8 | en_US.utf8 | =c/postgres          +
                        |          |          |            |            | postgres=CTc/postgres
 template1              | postgres | UTF8     | en_US.utf8 | en_US.utf8 | =c/postgres          +
                        |          |          |            |            | postgres=CTc/postgres
 trying_out_development | postgres | UTF8     | en_US.utf8 | en_US.utf8 |
 trying_out_production  | postgres | UTF8     | en_US.utf8 | en_US.utf8 |
 trying_out_test        | postgres | UTF8     | en_US.utf8 | en_US.utf8 |
(6 rows)

postgres=# \c trying_out_development
You are now connected to database "trying_out_development" as user "postgres".
trying_out_development=# \dt
Did not find any relations.
```

No tables yet... Back to readme: start the dev server and open it up in a browser.

```bash
buffalo dev
open http://127.0.0.1:3000
```

Liftoff! This gives me a development server with auto re-compilation.

At this point I was a bit curious about the generated code: tried running `tree` - 1937 directories, 12700 files....
That's quite a bit, let's come back later since I have no good idea where to start.

# Generating resources

Let's try to generate a few more things instead. At this point I'm not following "Getting started" anymore and am trying
to emulate a demo I saw on YouTube - rails-like generation of CRUD resources.

The `buffalo` tool is pretty easy to use and discover commands, after a bit of faffing around with `-h` I get to

```bash
buffalo generate resource users
open http://localhost:3000/users
```

500, apparently i need to run migrations now to create the new tables

```bash
buffalo pop migrate up
```

Yay I now have CRUD! but no fields xD. Not that I didn't have to restart or recompile anything, dev server from before
picked up changes automatically.

Let's revert

```bash
buffalo pop migrate down # 1 step by default
```
need to also manually remove the migration in /migrations and clean up routes in actions/app.go (figured this out by
trial and error)

At this point I consulted the docs - they are quite nicely structured, I quickly found [Generating
resources](https://gobuffalo.io/en/docs/resources/#generating-resources) and now know how to generate resources with
some actual fields

```bash
buffalo generate resource users name email tier:int8
```

I added non-string column just for kicks so I can see how it looks when I need to specify the type.
And manually run the migrations.

```bash
buffalo pop migrate up
```

And we have functional crud!

## Create
![new](/images/buffalo/new_user.png)
## List
![list](/images/buffalo/list_users.png)
## Read
![details](/images/buffalo/user_details.png)
## Update
![edit](/images/buffalo/edit_user.png)


# Generated code

Now let's take a look at the generated code. First for our users. Apparently we generated action, model, and migrations.
Let's start at the rear.

## Migrations

Migrations seem quite straightforward (apparently a DSL of some sort to be able to run against
different SQL dialects).

up:

```
create_table("users") {
	t.Column("id", "uuid", {primary: true})
	t.Column("name", "string", {})
	t.Column("email", "string", {})
	t.Column("tier", "int8", {})
	t.Timestamps()
}
```

down
```
drop_table("users")
```

## Model

Apparently we also got test files for bot models and actions, complete with placeholder failing tests. NICE! There is
even `buffalo test` which does exactly what you expect.


```go
// User is used by pop to map your users database table to your go code.
type User struct {
    ID uuid.UUID `json:"id" db:"id"`
    Name string `json:"name" db:"name"`
    Email string `json:"email" db:"email"`
    Tier int8 `json:"tier" db:"tier"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// String is not required by pop and may be deleted
func (u User) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Users is not required by pop and may be deleted
type Users []User

// String is not required by pop and may be deleted
func (u Users) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (u *User) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: u.Name, Name: "Name"},
		&validators.StringIsPresent{Field: u.Email, Name: "Email"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (u *User) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (u *User) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
```

This is in fact the whole model - copied here since it's quite terse. I very much like the fact a regular Go struct and
not some smart active record as I prefer less magic.

## Actions
Action itself is actually quite hefty - comes in at 258 lines of code. [Here it
is](https://github.com/edofic/trying-out/blob/main/buffalo/actions/users.go) in it's entirety if you want to read
through, I'll just do a brief overview of my reading.


```go
type UsersResource struct{
  buffalo.Resource
}

func (v UsersResource) List(c buffalo.Context) error
func (v UsersResource) Show(c buffalo.Context) error
func (v UsersResource) New(c buffalo.Context) error
func (v UsersResource) Create(c buffalo.Context) error
func (v UsersResource) Edit(c buffalo.Context) error
func (v UsersResource) Update(c buffalo.Context) error
func (v UsersResource) Destroy(c buffalo.Context) error
```

There already is support for pagination, content negotiation (same content can be rendered as html/xml/json), model
validation on mutiation. Looks a bit verbose but if you read carefully the code is quite nice, idiomatic Go, with just
some patterns repeating for all endpoints.

I do have mixed feelings about it. On one hand it's nice to have everything spelled out like this - very easy to
scaffold and jump in to make customizations. On the other hand I'm afraid this will leave you with hard to maintain
copy-pasta over a longer project.

## Templates

There is [quite a bit going on here](https://github.com/edofic/trying-out/tree/main/buffalo/templates/users). I've never
used [Plush](https://github.com/gobuffalo/plush) but the syntax looks familiar and the templates look like regular
Bootstrap. Basically this follows the same pattern: no surprises, no magic, human-editable code, going for max
productivity.

## Other generated things

With slightly better understanding i can now try to read through the rest of it. Luckily we also got a generated
`.gitignore` so I can ask git (`git ls-files`) for a list of important files.

All in all (after I generated my users resource) I'm looking at 50 files, 1387 lines of code. Not too much to actually
read (or skim) through and get a feeling what's going on. And of course [docs cover the important
folders](https://gobuffalo.io/en/docs/getting-started/directory-structure/).

I'll just list a few things that have caught my eye
- sane multi-stage dockerfie
- sane gitignore
- go.mod
- main.go - very simple entry point that basically starts the server but leaves you with a place to customize. Again no
  magic, you can build this with regular `go build`
- public assets (complete with favicon & robots.txt)
- assets pipeline (yarn & webpack based, apparently "just works" with the `buffalo` tool)

# Conclusion

I think Buffalo delivers on the promise: you can hit the ground running, yet there is very little magic and it's built
upon other established libraries. But in order to achieve this it is opinionated and possibly restrictive. But this is
fine - it's a trade off you can make: I'll glacdly pick up Buffalo when I have a need for a quick crud-y site (bye
django...) but probably not for starting a new specialised backend/internal service. Then again, if it's curd-y
database-backed thing...

All in all I think it's a good option to have in one's toolbox and I'll probably be using it.
