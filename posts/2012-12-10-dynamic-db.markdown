---
title: Cool Monday: Exploration of dynamic db acces from scala 
--- 

I use [scala](http://www.scala-lang.org/ "Scala (programming language)") on
[Android](http://www.t-mobile.com/shop/phones/?capcode=AGE "Android phones")
and I don't like the integrated database
[API](http://en.wikipedia.org/wiki/Application_programming_interface "Application programming interface").
It's very verbose and very stateful. I had written my own
[ORM](http://en.wikipedia.org/wiki/Object-relational_mapping "Object-relational mapping")([DAO](http://en.wikipedia.org/wiki/Data_access_object "Data access object")
would be a more appropriate tag) a while back, before I used scala but
it's not enough anymore. So now I'm on a quest for a better database
API. My dream is something small that handles schema for me and is
[type-safe](http://en.wikipedia.org/wiki/Type_safety "Type safety"). A
nice
[DSL](http://en.wikipedia.org/wiki/Digital_subscriber_line "Digital subscriber line")
that is converted to
[SQL](http://www.iso.org/iso/catalogue_detail.htm?csnumber=45498 "SQL")
at [compile
time](http://en.wikipedia.org/wiki/Compile_time "Compile time")  and
does code generation. So it's fast like hand writing everything. But
reduces code footprint by an order of magnitude(at least).
[Scala](http://www.scala-lang.org/ "Scala (programming language)")
[SLICK](http://slick.typesafe.com/) looks promising. It fits most
requirements. But it's kinda big for android projects(you need scala
library too!) and has not yet hit a stable version so I wouldn't be
comfortable shipping it. Will definitely give it a thorough test when
scala 2.10 is stable and SLICK is released. Oh, and it needs a third
party [JDBC
 driver](http://en.wikipedia.org/wiki/JDBC_driver "JDBC driver") for
Android. This is another level of abstraction and therefore another
source of slowness. I contemplated writing my own clone targeted at
Android but   never came around to actually doing it(yet!). It seems
like a herculean task for single developer working in spare time.


### Meanwhile

Yesterday I stared thinking how [dynamic
languages](http://en.wikipedia.org/wiki/Dynamic_programming_language "Dynamic programming language")
handle databases. And I got an idea. Scala has type Dynamic that does
compilation magic to provide syntactic sugar for working with dynamic
languages or objects. Here's an idea: do queries in plain SQL and
perform extraction of data in a dynamic way. 

And how to do this? Just wrap up Cursor to provide necessary methods. 
```scala
class WrappedCursor(cursor: Cursor) implements Cursor{  //delegated methods go here}
```
Why I need this? Cake pattern of course, Dynamic cursor get's mixed in.
```scala
trait DynamicCursor extends Dynamic{ this: Cursor =>  
    def selectDynamic(name: String) = getColumn(getColumnIndex(name))
    def getColumn(index: Int) = getType(index) match {
        case Cursor.FIELD_TYPE_BLOB => getBlob(index)
        case Cursor.FIELD_TYPE_FLOAT => getDouble(index)
        case Cursor.FIELD_TYPE_INTEGER => getLong(index)
        case Cursor.FIELD_TYPE_NULL => null
        case Cursor.FIELD_TYPE_STRING => getString(index)  
    }  
    def toSeq = (0 until getColumnCount) map getColumn
}

```
I targeted API level 14(Ice Cream Sandwich) since getType(method on
Cursor) is available from 11 on.    Key method here is getColumn that
abstracts over types. So you can read a column and  do pattern matching
on it. Or you are evil and use implicit conversions from Any to String,
Long etc... Or use implicit conversion to "converter"
```scala
implicit class Converter(val value: Any) extends AnyVal{
    def blob = value.asInstanceOf[Array[Byte]]
    def double = value.asInstanceOf[Double]
    def long = value.asInstanceOf[Long]
    def string = value.asInstanceOf[String]
}
```
But the real deal is selectDynamic. This allows you to write code like
this
```scala
val c = new WrappedCursor(result) with DynamicCursorc.someColumn.long
```
This compiles down to selectDynamic("someColumn") that calls getColumn
and finally implicit conversion is inserted that allows for terse cast
to Long.
And I threw in a conversion from row to Seq that does a snapshot of
current row. This allows pattern matching on rows. Any you can now
construct a Stream that will handle Cursor state and lazily evaluate and
store these snapshots. Therefore you can abstract away all mutability
and handle cursor as immutable collection.

Said conversion to stream
```scala
def CursorStream(cursor: DynamicCursorRaw with Cursor) = {
    def loop(): Stream[Seq[Any]] = {
        if(cursor.isAfterLast)
            Stream.empty[Seq[Any]]
        else {
            val snapshot = cursor.toSeq
            cursor.moveToNext()
            snapshot #:: loop()    
        }  
    }  
    cursor.moveToFirst()  
    loop()
}
```

And some more implicits to help 
```scala
implicit class RichCursorRaw(cursor: Cursor) extends AnyVal{
    def dynamicRaw = new WrappedCursor(cursor) with DynamicCursorRaw  
    def toStream = CursorStream(dynamicRaw)
}
```

All the source is in the project on
github [https://github.com/edofic/dynamic-db-android](https://github.com/edofic/dynamic-db-android) (work
in progress).
