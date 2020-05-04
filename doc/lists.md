# Working with lists

## The basics

Lists are one of Elmo's core type. They can be defined as literals using brackets. Everything between two brackets is part of the same list.

```elmo

# create a list using spaces to separate list items

l: [1 2 3]

# or use comma's

l: [1,2,3]

# or use newlines

l: [
  1
  2
  3
]

```

In Elmo, lists can contain anything and list items can be of different types.

```elmo
l: [1 "chipotle" for me]
```

Lists can contain other lists.

```elmo
l: [[1 2 3] [4 5 6]]
```

### lists are functions

A list itself acts like a function that can be used to access list content.

```elmo
# get the first item of a list

l: [0 1 2 3 4 5 6 7 8 9]
first: (l 0)

# get the last item of a list

last: (l -1)

# get a range

range: (l 0 5)

# get a range from a position to the end

lastfew: (l 7 -1)

# get the list in reverse

reverse: (l -1 0)
```

### Get the lenght of a list

To get the length of a list, elmo's build-in ``len`` function can be used

```
len [1 2 3] |eq 3 |assert
```

## The list module

Elmo comes a with a build in module with some handy functions that operate on lists.

### list.new

Instead of using the bracket syntax to construct a list, it's also possble to do it prgrammatically using the ``list.new`` function.

```elmo
list: (load list)

l: (list.new 0 1 2 3 4 5 6 7 8 9)
```

list.new simply uses all passed arguments to construct a list. This can be handy when you want to convert the result of a function call which returns multiple values, into a list.

```elmo
list: (load list)
f: (func { return 1 2 3 })
l: (f |list.new)
```

### list.tuple

When you want to pipe all items of a list into a function as function arguments, you can use the ``list.tuple`` function.

```elmo
list: (load list)

# function accepting 3 arguments

acceptXYZ: (func x y z { return (multiply $x $y | plus $z) })

# function producing a list with 3 items

produceList: produceList: (func { return [2 4 6] })

# pipe both functions using list.tuple

result: (produceList | list.tuple | acceptXYZ)
```

### list.at

Instead of using a list itself as a function to retrieve items from a list, ``list.at`` can be used as well.

```elmo
list: (load list)

l: [0 1 2 3 4 5 6 7 8 9]

# get one item

first: (list.at $l 0)

# get multiple items (in this case from the 4th to (and including) the 9th)

more: (list.at $l 3 8)
```

### list.append and list.append!

The list module contains two functions to add items to a list. ``list.append`` and ``list.append!``. The first does not change the original list but creates a new list with the appended items. The second (note, the difference is the ``!`` at the end), will actually change the content of a list.

```elmo
list: (load list)

l1: [0 1 2]
l2: (list.append $l 3 4 5)
eq $l1 [0 1 2] |assert
eq $l2 [0 1 2 3 4 5] |assert

l3: (list.append! $l 3 4 5)
eq $l1 [0 1 2 3 4 5] |assert
eq $l2 [0 1 2 3 4 5] |assert
```

### list.prepend and list.prepend!

Just like append, there are two variants of prepend: ``list.prepend`` and ``list.prepend!``. Prepend adds items to the start of the list. And only ``list.prepend!`` actually changes the content of the list.

```elmo
list: (load list)

l1: [0 1 2]
l2: (list.prepend $l 3 4 5)
eq $l1 [0 1 2] |assert
eq $l2 [5 4 3 0 1 2 |assert

l3: (list.prepend! $l 3 4 5)
eq $l1 [5 4 3 0 1 2] |assert
eq $l2 [5 4 3 0 1 2] |assert
```

### list.each

To iterate over all items the ``list.each`` function can be used. ``list.each`` accept a block of code and an identifier which is used to pass an items's value to the block as variable.

```elmo
list: (load list)
list.each [1 2 3] item { puts $item }
```

It's also possible to pass the item's position in the list to the block of code.

```elmo
list: (load list)
list.each [1 2 3] item index { puts "\{$index}:\{$item}" }
```

### list.map

To map all values of a list to a different value, ``list.map`` can be used. ``list.map`` works the same way as ``list.each`` but constructs a list with all values which are returned by the given block of code.

```elmo
list: (load list)
list.map [1 2 3] item { return (multiply $item 2) }
```

And just like ``list.each``, it's possible to pass the item's position in the list to the block of code.

```elmo
list: (load list)
list.map [1 2 3] item index { return (multiply $item (incr $index)) }
```

### list.where

To filter a list and to select only the items which comply to a given 'where' block, ``list.where`` can be used. Like ``list.map`` and ``list.each``, ``list.where`` accepts a block of code and identifiers to pass values and positions to the block of code.

```elmo
list: (load list)
list.where [1 2 3 4 5 6] item { return (modulo $item 2 |eq 0) }

# and using an index indentifier

list.where [1 2 3 4 5 6] item index { return (modulo $index 2 |eq 0) }
```


### list.sort

Lists can be sorted using either ``list.sort`` (sort the list and return a new list with sorted items) or ``list.sort!`` (sort and change the list)

```elmo
list: (load list)
l1: [c h i p o t l e]
l2: (list.sort $l1)
eq $l1 [c h i p o t l e] |assert
eq $l2 [c e h i l o p t] |assert


l3: (list.sort! $l1)
eq $l1 [c e h i l o p t] |assert
eq $l3 [c e h i l o p t] |assert
```

### list.flatten

When lists contain other lists, the ``list.flatten`` function can be used to flatten everything (or until a given depth) into one list.

```elmo
list: (load list)
list.flatten [1 2 [1 [a b] 3] 4] |eq [1 2 1 a b 3 4] |assert
```

same example but now with a specified depth:

```elmo
list: (load list)
list.flatten [1 2 [1 [a b] 3] 4] 1 |eq [1 2 1 [a b] 3 4] |assert
```
