# Embedding Elmo

Elmo is designed to be embedded in your own Go programs. In the examples directory
of Elmo, there is very simple example showing the basic usage. Testing it is easy.

```bash
# in the directory where elmo's sources are cloned
cd examples/embed
go install
embed
```

Like running Elmo, a REPL is started showing the Elmo prompt. But now, an additional
module is available that contains two additional functions.

```elmo
example: (load example)
example.chipotle
example.jalapeno
```

Let's explain how this works and how you can create your own Elmo extensions.

In this example, chipotle and jalapeno are both Go functions that are run in
an Elmo context (needed to interact with Elmo) and take arguments. A simple function that
does nothing, just like jalapeno, looks like this:

```go
func jalapeno() elmo.NamedValue {
	return elmo.NewGoFunction("jalapeno", func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

		// Do it your self
		//
		return elmo.Nothing
	})
}
```

Jalapeno is actually a function that returns a so called 'named' value that Elmo
can use to assign to a variable. And this 'named' value, is just a function with a name. As you can
see the actual function does nothing. So it returns elmo.Nothing. Which is Elmo's build-in
variant of nil.

A little bit more exiting is chipotle. This functions actually does something. It
will return a string. And if an integer argument is passed, it will even repeat that
string multiple times. Chipotle looks like:

```go
func chipotle() elmo.NamedValue {
	return elmo.NewGoFunction("chipotle", func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

		argLen, ok, err := elmo.CheckArguments(arguments, 0, 1, "chipotle", "<string>")
		if !ok {
			return err
		}

		if argLen == 0 {
			return elmo.NewStringLiteral("love them!")
		}

		value := elmo.EvalArgument(context, arguments[0])
		if value.Type() == elmo.TypeInteger {
			return elmo.NewStringLiteral(strings.Repeat("love them!", int(value.Internal().(int64))))
		}

		return elmo.NewErrorValue("please use nothing or an integer value as first argument")
	})
}
```

It first checks the arguments and it will accept minimal 0 arguments and maximal 1 argument.
If this check fails, an error object (elmo.Value) is returned. Then it checks the
number of arguments and if an argument is specified, it will check if it's integer.

Now the two functions are there, we can create the module that can be loaded from
within Elmo.

```go
var Module = elmo.NewModule("example", func(context elmo.RunContext) elmo.Value {
	return elmo.NewMappingForModule(context, []elmo.NamedValue{
		chipotle(),
		jalapeno()})
})
```

Our module is called example (so it can be loaded using 'load example') and is
initialized by the NewModule constructor through the given initializer function.
Which will register both chipotle and jalapeno.

The last part is creating the main routine that will run the extended version of Elmo.
Either in REPL mode or in regular mode.

```go
package main

import "github.com/okke/elmo/runner"

func main() {
	// create a context with all elmo's default modules
	//
	context := runner.NewMainContext()

	// add our own little module
	//
	context.RegisterModule(Module)

	// and run!
	//
	runner := runner.NewRunner(context)
	runner.Main()
}
```

That's it.
