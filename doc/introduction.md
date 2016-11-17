# background

Learn a new language every year. It's the Pragmatic programmers advice every programmer should take under serious consideration. Don't get stuck in eating mac and cheese every day of the week just because it smells comfy. And given that advice, why not go one step further. Write your own language, maybe not every year, but once in a while. Just because it's fun, because it can be brainteasing, because it's a very good learning experience and just because other programmers can learn your language when looking for a new language to learn. I started to learn Go a few years ago but never found the opportunity to really learn the real tricks and trades. So looking for a reason to do serious software development with Go, I decided to create Elmo. It's combining my passion for software engineering, programming languages and my passion to explore new things.

From a design point of view, Elmo is a pretty simple language. On purpose. The first scripting language that came across my path, I'm talking years ago, was Tcl. Pronounce as 'tickle'. A simple language with a huge potential. A lovely glue between hardcore c and c++ components. After Tcl many other languages have crossed my path and a few of them, including Go, really made me smile. That's why simplicity, glueing code and creating big smiles are the main objectives of Elmo. And note, Elmo is not designed for speed. It's glue, not nitro.

# getting to know Elmo

Everything in Elmo is a function that can be called like calling commands. Every call start with a symbol, a function name, followed by zero or more arguments that are separated by spaces. A very simple example:

```elmo
puts "chipotles in a jar"
puts 3 " chipotles in a jar"
```

Arguments can, for readability purpose, be separated by comma's. And comma's can be followed by a new line so function calls can use more than one line of text. Like this:

```elmo
puts "the battle between ",
     "the bold and brave chipotle ",
     "and the nasty jalapeno"
```

Arguments can be function calls themselves using a Lisp like polish notation.

```elmo
puts (multiply 11 11) " chipotles"
```

Elmo, like most languages, has a notion of variables. That can be assigned using the set function or that can be set using a special shortcut construction

```elmo
# using set
#
set pepper "chipotle"

# or using the assignment shortcut
#
pepper: "chipotle"
```

And, of course, variables can be used as function arguments. Just like they are functions themselves using the (...) construction or by using a resolve shortcut, the dollar sign, well known from many other scripting languages.

```elmo
pepper: "chipotle"

puts (pepper)
puts $pepper
```

Or, welcome to Elmo's dynamic nature, even as function names.

```elmo
say: puts

$say "don't do this at home"

# or like

(say) "chipotles for president!"
```

Elmo's core types are symbols, strings, numbers, functions, lists, dictionaries and code blocks.

```elmo
symbol: pepper
string: "chipotles"
integer: 38
float: 3.14
f: (func i { return (multiply 2 $i) } )
list: [1 2 3]
dict: {
  name: "jalapeno"
  hotness: 3
}
puts $dict.name, " eats ", (f $integer), " ", $string, " for breakfast"
```
