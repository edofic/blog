---
title: "GraphQL ❤️  SQLite"
date: 2023-05-06T13:20:31+02:00
---

If you've done anything nontrivial with GraphQL you're probably familiar with
how "N+1 select problem" sneaks up on you. If not, this is how [gqlgen
docs](https://gqlgen.com/reference/dataloaders/) explain it:

> Imagine your graph has query that lists todos…

```graphql
query { todos { user { name } } }
```

> and the todo.user resolver reads the User from a database…

```go
func (r *todoResolver) User(ctx context.Context, obj *model.Todo) (*model.User, error) {
	res := db.LogAndQuery(
		r.Conn,
		"SELECT id, name FROM users WHERE id = ?",
		obj.UserID,
	)
	defer res.Close()

	if !res.Next() {
		return nil, nil
	}
	var user model.User
	if err := res.Scan(&user.ID, &user.Name); err != nil {
		panic(err)
	}
	return &user, nil
}
```

> The query executor will call the Query.Todos resolver which does a select *
> from todo and returns N todos. If the nested User is selected, the above
> UserRaw resolver will run a separate query for each user, resulting in N+1
> database queries. e.g.

```sql
SELECT id, todo, user_id FROM todo
SELECT id, name FROM users WHERE id = ?
SELECT id, name FROM users WHERE id = ?
SELECT id, name FROM users WHERE id = ?
SELECT id, name FROM users WHERE id = ?
SELECT id, name FROM users WHERE id = ?
SELECT id, name FROM users WHERE id = ?
```

> Whats even worse? most of those todos are all owned by the same user! We can
> do better than this.

And then goes on to explain how to use a [dataloader](https://github.com/graph-gophers/dataloader) to solve this by effectively batching and deduplicating queries. 
The solution is neat and the damage to your resolver logic minimal. Heck, I use this in production at work. So what's the issue exactly? 

# Enter SQLite

![where we're going we don't need batching](/images/we-dont-need-batching.jpeg)

The above wisdom only applies when you use an external database. The performance characteristics of using an embedded database (e.g. SQLite) are different. In fact the SQLite homepage explicitly claims that [Many Small Queries Are Efficient In SQLite
](https://www.sqlite.org/np1queryprob.html).

So how does that apply? We'll you can just skip dataloders and naively query. And the performance will be similar. See there is a minimal overhead to doing a query but since in the above case where you're repeatedly hitting the same user it's pretty much guaranteed to already be in memory most of the time. So the overhead is minimal compared to a map lookup - which you would be doing otherwise. And there is a cost to using a dataloader too, including (but not limited to) having to wait a bit before firing a query because you don't know if more are coming or not. And this may actually be the limiting factor for getting your response latency down for cheap requests. 

So GraphQL loves to make a lot of small queries and SQLite is very good at handling that. Truly a match made in heaven.
