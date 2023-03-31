---
title: "(Re)Learning AWK with ChatGPT"
date: 2023-03-31T16:12:00+02:00
---

[awk](https://en.wikipedia.org/wiki/AWK) is the Swiss army knife of the command
line. If you have a text-shaped problem then awk can probably be a solution.
[Here’s](https://www.youtube.com/watch?v=Sg4U4r_AgJU) a fabulous talk by Brian
Kernighan (the k in awk) on why this may be the case.

But it’s also a tool that I use just rarely enough that I keep forgetting
details. So everytime I reach for awk it’s a painful experience of not
remembering and having to cobble together 3 different StackOverflow answers to
get something resembling the script I want. Only to then realize I’ve been
“holding it wrong” and there’s a better approach 🤦.

I still believe it’s supremely useful as it can do everything and is available
everywhere so I keep inflicting this pain on myself.

# This time it’s different

[ChatGPT](https://chat.openai.com/) changes this. It’s training data
(presumably) already contains all of the StackOverflow answers that I was
combing through and it can cobble them together for me. So let’s try to use it
to cobble them together for me and try to learn more efficiently. 

## Tangent: chatgpt-cli

Since I like to stay in my terminal while I do deep-focus work I want to also
access ChatGPT from command line. I couldn’t find a client that I would like so
I wrote my own:
[github.com/edofic/chatgpt-cli](https://github.com/edofic/chatgpt-cli). It’s a
one shot command that fires off you query and streams the result real-time back
to your cli. But also allows you to ask follow up questions. 

## So... awk?

Motivating example this time was multi-line search. Specifically I wanted to
extract a section of a document (multiple lines) with known patterns for the
first and the last line. 

I remember doing this before... You create what amounts to a state machine. If a
line matches pattern A you go to `print` state, in `print` state you match for
pattern B, if match then exit, otherwise print the line. 

But staring at a blank prompt and trying to transcribe this to awk syntax by
heart got me thinking. I reached for ChatGPT

```text
$ chatgpt-cli 'write an awk script that prints all lines 
              from a match of "tags:" to next empty line'

Assuming the input file has lines with "tags:" and empty
lines separating the paragraphs, the following awk script
can be used to print all lines between a match of "tags:"
and the next empty line:

awk '/tags:/{p=1}p;/^$/{p=0}' input_file

Explanation:
- `/tags:/{p=1}`: If a line contains "tags:", set the 
  variable "p" to 1
- `p;`: If the variable "p" is 1 (i.e., if we have found a 
   match of "tags:"), print the line
- `/^$/{p=0}`: If a line is empty (i.e., it contains only 
  the beginning-of-line and end-of-line markers), set the
  variable "p" to 0 to stop printing

This script will print all lines between the first
occurrence of "tags:" and the next empty line. If there are
multiple paragraphs with "tags:", this script will print all
lines between each occurrence of "tags:" and the next empty
line in the paragraph.
```

This is exactly what I wanted. But it includes the explanation for all the
syntax. 

And it can do one think that an amalgamation of StackOverflow answers cant:
follow up questions.

```text
$ chatgpt-cli -c 'why does p; print the line? there is no 
                  print command'

In awk, the `print` command is often used to output the
results of a script. However, awk also has an implied print
statement that activates when no explicit print statement is
used.

In the script `p;`, `p` is a variable that can have two
values: 0 or 1. When `p` is 0, nothing is printed. When `p`
is 1 (i.e., the regex search matches), the entire line is
printed as part of the implied print statement.

So, while there is no explicit `print` statement, the script
`p;` is equivalent to `print $0;`. The `$0` variable
represents the entire line, and since there is no explicit
print command to specify a different output, the default
`print` command is used to output the entire line.
```

# Takeaway

Key message here is not that you can use ChatGPT to generate awk scripts. It’s
that in doing that you can also learn how these scripts work and how to do them
yourself. 

And at least for me that’s a much faster and frictionless solution that
searching for answers and fitting them together.
