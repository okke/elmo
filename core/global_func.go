package elmo

// func arguments can have default values. This is done by giving the
// argument a name that ends with an '?' and by using the next argument
// as default value. So it's possible to create a function like:
//
// greet: (func name greeting?"Hello" { echo "\{$greeting} \{$name}"})
//
// in this case ["name", "greeting"] and {"greeting":"Hello"} are returned
//
func extractArgNamesAndDefaultValues(argNameValues []Value) ([]string, map[string]Value) {
	useArgNames := make([]string, 0, len(argNameValues))
	defaultValues := make(map[string]Value, 0)

	useNextAsDefault := false
	name := ""

	for _, nameValue := range argNameValues {
		if !useNextAsDefault {
			name = nameValue.String()
			if name[len(name)-1] == '?' {
				useNextAsDefault = true
				name = name[:len(name)-1]
			}

			useArgNames = append(useArgNames, name)
		} else {
			// argument is default value for previous arg
			defaultValues[name] = nameValue
			useNextAsDefault = false
		}
	}

	if name != "" && useNextAsDefault {
		useArgNames = append(useArgNames, name)
	}

	return useArgNames, defaultValues
}

func createGoFunc(argNameValues []Value, evaluator func(evalContext RunContext) Value) (func(innerContext RunContext, innerArguments []Argument) Value, []string) {

	useArgNames, defaultValues := extractArgNamesAndDefaultValues(argNameValues)
	injectDefaultValues := len(defaultValues) > 0

	return func(innerContext RunContext, innerArguments []Argument) Value {

		cloneFrom := innerContext
		if cloneFrom.Parent() != nil {
			cloneFrom = cloneFrom.Parent()
		}
		subContext := cloneFrom.CreateSubContext()

		if innerContext.This() != nil {
			subContext.Set("this", innerContext.This())
		}

		if injectDefaultValues {
			for k, v := range defaultValues {
				subContext.Set(k, v)
			}
		}

		maxArgs := len(useArgNames)
		for i, v := range innerArguments {
			if i < maxArgs {
				subContext.Set(useArgNames[i], EvalArgument(innerContext, v))
			}
		}

		return evaluator(subContext)

	}, useArgNames
}

func splitArgumentsForFunc(context RunContext, argStart int, arguments []Argument) ([]Value, Value) {
	argNamesAsArgument := arguments[argStart : len(arguments)-1]
	code := EvalArgument(context, arguments[len(arguments)-1])
	argNames := make([]Value, len(argNamesAsArgument))

	for i, v := range argNamesAsArgument {
		argNames[i] = EvalArgument(context, v)
	}
	return argNames, code
}

func _func() NamedValue {
	return NewGoFunctionWithHelp("func", `Create a new function
		Usage: func <help>? <symbol>* {...}
		Returns: a new function

		When first argument is a string, this value will be used as help text

		Given symbols denote function parameter names.

		Examples:

		> func a {}
		will create a function that accepts one parameter called 'a' (and does nothing)
		> func a { return $a }
		will create an echo function

		Note, function can be used once they are assigned to a variable

		> f: (func a { return $a })

		Example with help text:

		> f: (func a "we need more chipotles" {})
		> help f
		will result in "we need more chipotles"`,

		func(context RunContext, arguments []Argument) Value {

			argLen := len(arguments)
			if argLen == 0 {
				return NewErrorValue("func expects <help>? <identifier>* {...}")
			}

			argStart := 0

			help := ""
			if arguments[0].Type() == TypeString {
				if argLen == 1 {
					return NewErrorValue("func with help should at least have a body also")
				}
				help = EvalArgument2String(context, arguments[0])
				argStart = 1
			}

			argNames, code := splitArgumentsForFunc(context, argStart, arguments)

			if code.Type() != TypeBlock {
				return NewErrorValue("invalid call to func, last parameter must be a block")
			}

			block := code.(Block)

			evaluator := func(evalContext RunContext) Value {
				return block.Run(evalContext, NoArguments)
			}

			f, useArgNames := createGoFunc(argNames, evaluator)
			return NewGoFunctionWithBlock("anonymous", help, f, useArgNames, code.(Block))

		})
}

func template() NamedValue {
	return NewGoFunctionWithHelp("template", `Create a new template function
		Usage: template <help>? <symbol>* &"..."
		Returns: a new function

		When first argument is a string, this value will be used as help text

		Given symbols denote function parameter names.

		Examples:

		> echo: (template a &"\{$a}")
		will create a function that accepts one parameter actually returns it

		> goFileGen: (template packageName code &"package \{$packageName}\n\{$code}")
		will create a function which can generate a Go file
		`,

		func(context RunContext, arguments []Argument) Value {

			argLen := len(arguments)
			if argLen == 0 {
				return NewErrorValue("template expects <help>? <identifier>* &\"...\"")
			}

			argStart := 0

			help := ""
			if arguments[0].Type() == TypeString {
				if argLen > 1 {
					help = EvalArgument2String(context, arguments[0])
					argStart = 1
				}
			}

			argNames, code := splitArgumentsForFunc(context, argStart, arguments)

			if code.Type() != TypeString {
				return NewErrorValue("invalid call to template, last parameter must be a string")
			}

			str := code.(StringValue)

			evaluator := func(evalContext RunContext) Value {
				return str.ResolveBlocks(evalContext)
			}

			f, _ := createGoFunc(argNames, evaluator)
			return NewGoFunctionWithHelp("template", help, f)

		})
}
