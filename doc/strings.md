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

### Multiline strings

Multiline strings declared using backticks do not support escape characters. Everything between the backticks is taken literal. So you can dump a block of text including newlines into a variable like this:

```elmo
text: `
first line
second line
`
```

If you want to include a backtick in your text, simply use the double backtick:

```elmo
multi: `bla``bla`
```

### Strings are functions

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

### Code interpolation in strings

Elmo supports the interpolation of code execution results, mostly used to include variables inside strings, through the usage of a special escape sequence: ``\{}``

A simple example:

```elmo
pepper: "jalapeno"
sku: "2500"

puts "I love \{pepper} because it's not that hot (\{sku})"
```

Everything between ``\{`` and ``}`` will be evaluated as elmo code. An example

```elmo
i: 3
puts "incr \{$i} = \{incr i}"
puts "incr \{$i} = \{incr i}"
```

Multiline strings (backtick strings) support interpolation as well. Use a `` `{} `` escape sequence

```elmo
puts `
first line
second line
and say hello `{name}
`
```

### String templates

String with code interpolation are directly evaluated. But it's also possible to treat them as kind of templates by evaluating them later. A template is created by the use of the special ``&`` construction which is also used to refer to functions. ``&``, Followed by a string, constructs an unevaluated string template.

```elmo
template: &"I love \{pepper} but I really love \{hotterPepper}"
```

This wil create a template that looks like ``I love \{...}but I really love \{...}`` Which can be excuted like this:

```elmo
pepper: jalapeno
hotterPepper: habanero
puts (eval $template)
```

Elmo has a build in ``template`` function that can be used to make life easy when working with templates. The ``template`` function works like the ``func`` function but instead of accepting a code block, it accepts a template string.

```elmo
pepper: "jalapeno"
twolapeno: (template &"\{$pepper}-\{$pepper}")
puts (twolapeno)
```

Previous example used a global variable called pepper. But template, just like func, accepts parameter names and applies parameter values to evaluate the template string. Example:

```elmo
goFileGen: (template packageName code &"package \{$packageName}\n\{$code}")
goFile: (goFileGen "mypackage" "// lots of golang code")
puts $goFile
```

### String length

Elmo has a ``len`` function that will return the length of elmo values. ``len`` Supports strings. So the retrieve the length of a string simply use ``len "something"``

## The String module

Elmo comes a with a build in module with some handy functions that operate on strings.

### string.at

Functions exactly like accessing a string using itself.

```elmo
string: (load string)
chip: (string.at "chipotle" 0 3)
assert (eq $chip "chip")
```

### string.concat

Concat multiple strings into one new string

```elmo
string: (load string)
pepper: (string concat "chi" "po" "tle")
assert (eq $pepper "chipotle")
```

### string.trim

Remove characters from the beginning and ending of a trim.

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

### string.trimLeft

It's possible to trim only from the beginning.

```elmo
string: (load string)
pepper: (string.trimLeft "...chipotle..." ".ec")
assert (eq $pepper "hipotle...")
```

### string.trimRight

And only from the end.

```elmo
string: (load string)
pepper: (string.trimRight "...chipotle..." ".ec")
assert (eq $pepper "...chipotl")
```

### string.trimPrefix

Instead of specifying a set of characters to trim, it's also possible to specify a prefix

```elmo
string: (load string)

pepper: (string.trimPrefix "http://chipotle.com" "http://")
assert (eq $pepper "chipotle.com")
```

### string.trimSuffix

Or specifiy a suffix that need to be trimmed.

```elmo
string: (load string)

pepper: (string.trimSuffix "chipotle.json" ".json")
assert (eq $pepper "chipotle")
```


### string.replace

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

### string.find

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

### string count

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

### string.split

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

### string.startsWith

Checks if a string starts with a given value.

```elmo
string: (load string)
assert (string.startsWith "chipotle" "chip")
```

### string.endsWith

Checks if a string ends with a given value.

```elmo
string: (load string)
assert (string.endsWith "jalapeno" "peno")
```

### string.upper

Convert a string to uppercase characters

```elmo
string (load string)
string.upper "upper" |eq "UPPER" |assert
```

### string.lower

Convert a string to lowercase characters

```elmo
string (load string)
string.lower "LOWER" |eq "lower" |assert
```
