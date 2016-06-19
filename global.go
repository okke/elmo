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
	context.SetNamed(_return())
	context.SetNamed(_func())

	return context
}

func set() NamedValue {
	return NewGoFunction("set", func(context RunContext, arguments []Argument) Value {

		// set expects exactly 2 arguments
		//
		if len(arguments) != 2 {
			panic("invalid call to set, expected 2 parameters: usage set <identifier> <value>")
		}

		name := evalArgument2String(context, arguments[0])
		value := evalArgument(context, arguments[1])

		context.Set(name, value)

		return value
	})
}

func get() NamedValue {
	return NewGoFunction("get", func(context RunContext, arguments []Argument) Value {
		// get expects exactly 1 argument
		//
		if len(arguments) != 1 {
			panic("invalid call to get, expected 1 parameter: usage get <identifier>")
		}

		result, found := context.Get(evalArgument2String(context, arguments[0]))
		if found {
			return result
		}

		return Nothing

	})
}

func _return() NamedValue {
	return NewGoFunction("return", func(context RunContext, arguments []Argument) Value {
		// return expects exactly 1 argument
		//
		if len(arguments) != 1 {
			panic("invalid call to return, expected 1 parameter: usage return <value>")
		}

		return evalArgument(context, arguments[0])
	})
}

func _func() NamedValue {
	return NewGoFunction("func", func(context RunContext, arguments []Argument) Value {

		// get expects at least 1 argument
		//
		if len(arguments) < 1 {
			panic("invalid call to func, expect at least 2 parameters: usage func <identifier> <identifier>* {...}")
		}

		argNamesAsArgument := arguments[0 : len(arguments)-1]
		block := arguments[len(arguments)-1]
		argNames := make([]string, len(argNamesAsArgument))

		for i, v := range argNamesAsArgument {
			argNames[i] = evalArgument2String(context, v)
		}

		if block.Type() != TypeBlock {
			panic("invalid call to func, last parameter must be a block: usage func <identifier> <identifier>* {...}")
		}

		return NewGoFunction("user_defined_function", func(innerContext RunContext, innerArguments []Argument) Value {

			if len(argNames) != len(innerArguments) {
				// TODO better error handling
				panic("argument mismatch")
			}

			cloneFrom := innerContext
			if cloneFrom.Parent() != nil {
				cloneFrom = cloneFrom.Parent()
			}
			subContext := cloneFrom.CreateSubContext()

			for i, v := range innerArguments {
				subContext.Set(argNames[i], evalArgument(innerContext, v))
			}

			return block.Value().(Block).Run(subContext, noArguments)
		})

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
