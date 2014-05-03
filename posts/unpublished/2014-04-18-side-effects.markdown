--- 
title: On the evil of side-effects
---

Consider this Scala function

```scala
def loopup(datasetName: Name, key: Key): Value = {
    val dataset = Dataset.load(datasetName)
    dataset.find(key)
}
```

It's just a utility function for loading some value from some dataset so we don't have to manually load it every time if we just want a single value. It's not that great, it could have a more descriptive name and handle failure with `Option` or `Either` but this is not my point here. 

My point is this function will eventually be used in a different setting because someone(maybe even you) will find it handy and won't bother with implementation details. I'm talking about something like

```scala
val keys: Seq[Key] = ...
val datasetName: Name = ...
val values = keys.map(loopup(datasetName, _)) //the problem!
```
This will load the dataset every time! And this might be an expensive operation hitting the disk or even network.