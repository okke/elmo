# Working with strings

## The basics

Strings are one of Elmo's core type. They can be defined as literals using double quotes or backticks.

```elmo
# regular string
#
s: "chipotle"

# multiline string
#
txt: `chipotle
jalapeno
chilli`
```

Double quoted strings can contain contain special escape charaters. The following characters are supported:

| escape | character |
| --- | --- |
| \t | tab |
| \n | newline |
| \" | quote |
| \\ | backslash |

A simple example:

```elmo
escaped: "bla\n\t\\\"bla"
```

Multiline strings declared using backticks do not support escape characters. Everything between the backticks is taken literal. The only exception is a double backtick that can be used to include a backtick character in a multiline string.

```elmo
multi: `bla``bla`
```

Strings, like most things in Elmo, act as functions. Sound strange? Well, it's not. The arguments passed to a string are used to access a strings content.

```elmo
# get the 4th character, (note, counting starts at 0) of chipotle
#
char1: ("chipotle" 3)

# do the same using a variable
#
pepper: "chipotle"
char2: (pepper 3)

# should be the same and should contain the 4th character
#
assert (eq $char1 $char2)
assert (eq $char1 "p")
```

Instead of accessing single characters, string functions can be used to extract substrings.

```elmo
pepper: "chipotle"

# extract from the first (0) to the 3rd characters
#
chip: (pepper 0 3)

# extract from the 5th until the end
#
otle: (pepper 4 -1)

assert (eq $chip "chip" |and (eq $otle "otle"))
```

Note the usage of negative indexes. A negative index is counting from the end of a string backward. So -1 is the last character, -2 the second to last etc.

There is a small little trick that uses a string function to reverse a string. By extract from the last (-1) to the first (0) character.

```elmo
pepper: ("chipotle" -1 0)
assert (eq $pepper "eltopihc")
```


# The String module

Elmo comes a with a build in module with some handy functions that operate on strings.

## string.at

Functions exactly like accessing a string using itself.

```elmo
string: (load string)
chip: (string.at "chipotle" 0 3)
assert (eq $chip "chip")
```

## string.len

Get the length of a string

```elmo
string: (load string)
len: (string.len "chipotle")
assert (eq $len 8)
```

## string.concat

Concat multiple strings into one new string

```elmo
string: (load string)
pepper: (string concat "chi" "po" "tle")
assert (eq $pepper "chipotle")
```

## string.trim

Remove characters from the beginning and ending of a trim.  This function operates
in different modes. By default it trims from both the beginning and the ending of a string.
And by default it trims spaces and tabs

```elmo
string: (load string)
pepper: (string trim "   chipotle\t")
assert (eq $pepper "chipotle")
```

It's possible to specify the characters that need to be trimmed.

```elmo
string: (load string)
pepper: (string.trim "...chipotle..." ".ec")
assert (eq $pepper "hipotl")
```

It's possible to trim only from the beginning.

```elmo
string: (load string)
pepper: (string.trim left "...chipotle..." ".ec")
assert (eq $pepper "hipotle...")
```

And only from the end.

```elmo
string: (load string)
pepper: (string.trim right "...chipotle..." ".ec")
assert (eq $pepper "...chipotl")
```

Instead of specifying a set of characters, it's also possible to specify a prefix
or suffix that need to be trimmed.

```elmo
string: (load string)

pepper: (string.trim suffix "chipotle.json" ".json")
assert (eq $pepper "chipotle")

pepper: (string.trim prefix "http://chipotle.com" "http://")
assert (eq $pepper "chipotle.com")
```

## string.replace

Replace values in string with something else.

```elmo
string: (load string)
pepper: (string.replace "chipotle" "o" "a")
assert (eq $pepper "chipatle")
```


By default it will only replace the first occurrence of the specified value.
It's also possible to replace all values.

```elmo
string: (load string)
pepper: (string.replace all "jalapeno" "a" "o")
assert (eq $pepper "jolopeno")
```

And it's possible to relace the last occurrence only.

```elmo
string: (load string)
pepper: (string.replace last "jalapeno" "a" "o")
assert (eq $pepper "jalopeno")
```

## string.find

Scans a string to search for a specified value. By default it returns the
index of the first occurrence.

```elmo
string: (load string)
index: (string.find "chipotle in a jar" "in")
assert (eq $index 9)
```

It's possible to scan for the last occurrence of a value

```elmo
string: (load string)
index: (string.find last "chipotle in a jar in a jar" "in")
assert (eq $index 18)
```

It's also possible  to scan for all occurrences. In that case, a list of indexes in returned.

```elmo
string: (load string)
all: (string.find all "chipotle in a jar in a jar" "in")
assert (eq $all [9 18])
```

## string count

Count the occurrences of a specific value in a string.

```elmo
string: (load string)
count: (string.count "chipotle in a jar in a jar" "in")
assert (eq $count 2)
```

Note, this is the same as using string.find in combination with the len function
of the list module

```elmo
string: (load string)
list: (load list)
count: string.find all "chipotle in a jar in a jar" "in" | list.len
assert (eq $count 2)
```

## string.split

Splits a string into substrings. When no divider is given, the string will be
split into a list of characters.

```elmo
string: (load string)
split: (string.split "chipotle")
assert (eq $split ["c" "h" "i" "p" "o" "t" "l" "e"])
```

A string can also be split by specifying a splitting value.

```elmo
string: (load string)
split: (string.split "chipotle,jalapeno,chilli" ",")
assert (eq $split ["chipotle" "jalapeno" "chilli"])
```

## string.startsWith

Checks if a string starts with a given value.

```elmo
string: (load string)
assert (string.startsWith "chipotle" "chip")
```

## string.endsWith

Checks if a string ends with a given value.

```elmo
string: (load string)
assert (string.endsWith "jalapeno" "peno")
```
