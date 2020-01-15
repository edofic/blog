---
title: Integrating Jira into Google Sheets
date: 2018-01-25
---

So I got a (recurring) task to compile some stats from Jira. You now, some
filtering and agregation. Possibly a pivot table and a chart. Sounds like a
perfect task for spreadsheets. But Jira doesn't have spreadsheets - at least
without addons, but that is not an options for me.

So I decided to use Google Sheets. And figure out a way to automatically fetch
data from Jira because there's no way in hell I'm doing that manually every
time. Turns out there are some scripts you are supposed to paste into sheets and
give them all the permissions. But I don't trust that, not without  properly
reviewing the thousands of lines of code.

Obviously I decided to roll my own. I figured a way that's actually really easy.
Just create a filter in Jira and click `Export -> Printable`. This will give you
an URL with all your tickets formatted nicely in a `<table>`.:

Why is this useful? Because Google Sheets knows how to scrape HTML tables. If
you append your credentials to the url. The whole formula will look something
like this:

```
=IMPORTHTML("https://my_instance.atlassian.net/sr/jira.issueviews:searchrequest-printable/14900/SearchRequest-14900.html?tempMax=1000&os_username=my_username&os_password=my_password", "table", 1)
```

This magically pulls your Jira table into a spreadsheet. Now you can use all
your spreadsheet tools to make reports out of this. Or even SQL! Google Sheets
has a `QUERY` function which takes a data range and an SQL-like expression to
run over it.

Now we just need a way to simply refresh the data (remember, this is a recurring
task). Sheet re-load the data when the URL changes and Jira ignores unknown
parameters. We can use this to add a dummy counter which we increment!

```
=IMPORTHTML(CONCATENATE("https://...my_passoword&dummy=", B1), "table", 1)
```

But this means we now need to edit a number. No fun. Let's write a script to do
it.

In `Tools` -> `Script Editor`


```
function increment() {
  SpreadsheetApp.getActiveSheet().getRange('B1').setValue(SpreadsheetApp.getActiveSheet().getRange('B1').getValue() + 1);
}
```

Now we just need a way to run it.

1. insert a drawing
1. right click
1. assign a script

Voila we have a button!
