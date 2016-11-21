# Working with functions

## Declaring functions

Elmo treats functions as first class citizens of the language. Almost Everything
is a function. And functions can be assigned to variables and can be passed to
other functions.

Functions are created using the build in 'func' function. It takes the names
of function arguments as arguments followed by a block of code. A simple examples
looks like:

```elmo
f: (func {
    # do nothing
    return $nil
})

square: (func x {
    return (multiply $x $x)
})
```

As soon as a return statement is found, execution will halt and the value specified
by the return statement is returned to the caller.

## Pipe function results

Elmo supports piping of function calls. This means that the result of a function
is used as an argument for the next function. One of the best examples of this feature building
complex boolean expression without having function calls in functions calls in functions calls.

```elmo
f: (func s {return (eq $s "chipotle")})

b: (func a b c {

  # check if a or b or c equals to "chipotle"
  #
  return (f $a |or (f $b) |or (f $c))
}
```

## Use multiple return values

Elmo supports multiple return values. So a function can return multiple values
to its caller. Multiple return values can be used to assign multiple variables
using one call to 'set'.

```elmo
# create a function thta returns two values
#
f: (func x {return (plus $x 1) (plus $x 2)})

# call it and assign the result to multiple variables
#
set a b (f 3)

assert (eq $a 4 |and (eq $b 5))  
```

Note, when dealing with multiple return values, the ":" shortcut does not work.

Multiple return values can also be used in pipes to other functions. All returned
values are used as arguments to the next function.

```elmo
string: (load string)

f: (func { return "chipotle" "jalapeno"})

both: (f | string.concat)
assert (eq $both "chipotlejalapeno")

more: (f | string.concat "chilli")
assert (eq $more "chipotlejalapenochilli")
```

## Function arguments

Since functions are Elmo's first class citizens, functions can be passed to other functions.
And of course called from within that other function.

```elmo

apply: (func v f {
  return (f $v)
})

val: (apply 5 (func x {
  return (multiply $x $x)
}))

```

In this case, a function is created at the moment it is needed. But there also
the situation where the function that need to be passed, already exists and assigned
to a variable. In that case, one would like to use this variable and does not want to recreate
the function. Using the variables symbol won't work. This will pass the symbol,
not the content of the variable. And using the $ shortcut won't work either. This
will evaluate the function and pass the result.

So following example won't work:

```elmo
square: (func x {
    return (multiply $x $x)
})

apply: (func v f {
  return (f $v)
})

# will return the symbol square
#
apply 5 square

# will result in an error, since square expects a parameter
#
apply 5 $square
```

In order to make this work, Elmo uses a special construction to point to
functions that are stored in a variable. Instead of using a '$', one can use
a '&' character. So &square refers to the function stored in a variable named
square.

```elmo
result: apply 5 &square
assert (eq $result 25)
```

And of course, these '&' constructions can be used in (multiple) return values
and can be piped to other functions.

```elmo
create: (func x { return $x &square})
result: create 5 |apply
assert (eq $result 25)
```
