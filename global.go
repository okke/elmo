package elmo

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

		name := ""

		// first argument should be an identifier or a call
		// that will result in an identifier
		//
		if arguments[0].Type() == TypeIdentifier {
			name = arguments[0].String()
		} else if arguments[0].Type() == TypeCall {
			value := arguments[0].Value().(Runnable).Run(context, []Argument{})
			name = value.String()
		} else {
			panic("first parameter of set must be an identifier or a function call that will evaluate to an identifier")
		}

		value := arguments[1].Value()
		if arguments[1].Type() == TypeCall {
			value = arguments[1].Value().(Runnable).Run(context, []Argument{})
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

		name := ""

		// first argument should be an identifier or a call
		// that will result in an identifier
		//
		if arguments[0].Type() == TypeIdentifier {
			name = arguments[0].String()
		} else if arguments[0].Type() == TypeCall {
			value := arguments[0].Value().(Runnable).Run(context, []Argument{})
			name = value.String()
		} else {
			panic("first parameter of set must be an identifier or a function call that will evaluate to an identifier")
		}

		result, found := context.Get(name)
		if found {
			return result
		}

		return Nothing

	})
}
