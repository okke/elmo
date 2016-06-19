package elmo

var noArguments = []Argument{}

// NewGlobalContext constructs a new context and initializes it with
// all default global values
//
func NewGlobalContext() RunContext {
	context := NewRunContext(nil)

	context.Set("nil", Nothing)

	context.SetNamed(set())
	context.SetNamed(get())

	return context
}

func set() NamedValue {
	return NewGoFunction("set", func(context RunContext, arguments []Argument) Value {

		// set expects exactly 2 arguments
		//
		if len(arguments) != 2 {
			panic("invalid call to set, expected 2 parameters")
		}

		name := evalArgument2String(context, arguments[0])

		value := arguments[1].Value()
		if arguments[1].Type() == TypeCall {
			value = arguments[1].Value().(Runnable).Run(context, noArguments)
		}

		context.Set(name, value)

		return value
	})
}

func get() NamedValue {
	return NewGoFunction("get", func(context RunContext, arguments []Argument) Value {

		// get expects exactly 1 argument
		//
		if len(arguments) != 1 {
			panic("invalid call to get, expected 1 parameter")
		}

		result, found := context.Get(evalArgument2String(context, arguments[0]))
		if found {
			return result
		}

		return Nothing

	})
}

func evalArgument(context RunContext, argument Argument) Value {

	if argument.Type() == TypeCall {
		return argument.Value().(Runnable).Run(context, noArguments)
	}

	return argument.Value()

}

func evalArgument2String(context RunContext, argument Argument) string {

	return evalArgument(context, argument).String()

}
