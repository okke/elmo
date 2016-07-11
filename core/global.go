package elmo

import (
	"fmt"
	"reflect"
)

var noArguments = []Argument{}

// GlobalContext is the shared runtime of all elmo scripts
//
var GlobalContext = NewGlobalContext()

// NewGlobalContext constructs a new context and initializes it with
// all default global values
//
func NewGlobalContext() RunContext {
	context := NewRunContext(nil)

	context.Set("nil", Nothing)
	context.Set("true", True)
	context.Set("false", False)

	context.SetNamed(set())
	context.SetNamed(get())
	context.SetNamed(_return())
	context.SetNamed(_func())
	context.SetNamed(_if())
	context.SetNamed(list())
	context.SetNamed(dict())
	context.SetNamed(mixin())
	context.SetNamed(load())
	context.SetNamed(puts())
	context.SetNamed(eq())
	context.SetNamed(ne())

	return context
}

func set() NamedValue {
	return NewGoFunction("set", func(context RunContext, arguments []Argument) Value {

		// set expects exactly 2 arguments
		//
		if len(arguments) != 2 {
			return NewErrorValue("invalid call to set, expected 2 parameters: usage set <identifier> <value>")
		}

		name := EvalArgument2String(context, arguments[0])
		value := EvalArgument(context, arguments[1])

		context.Set(name, value)

		return value
	})
}

func get() NamedValue {
	return NewGoFunction("get", func(context RunContext, arguments []Argument) Value {
		// get expects exactly 1 argument
		//
		if len(arguments) != 1 {
			return NewErrorValue("invalid call to get, expected 1 parameter: usage get <identifier>")
		}

		result, found := context.Get(EvalArgument2String(context, arguments[0]))
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
			return NewErrorValue("invalid call to return, expected 1 parameter: usage return <value>")
		}

		return EvalArgument(context, arguments[0])
	})
}

func _func() NamedValue {
	return NewGoFunction("func", func(context RunContext, arguments []Argument) Value {

		// get expects at least 1 argument
		//
		if len(arguments) < 1 {
			return NewErrorValue("invalid call to func, expect at least 1 parameter: usage func <identifier>* {...}")
		}

		argNamesAsArgument := arguments[0 : len(arguments)-1]
		block := arguments[len(arguments)-1]
		argNames := make([]string, len(argNamesAsArgument))

		for i, v := range argNamesAsArgument {
			argNames[i] = EvalArgument2String(context, v)
		}

		if block.Type() != TypeBlock {
			return NewErrorValue("invalid call to func, last parameter must be a block: usage func <identifier> <identifier>* {...}")
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
				subContext.Set(argNames[i], EvalArgument(innerContext, v))
			}

			return block.Value().(Block).Run(subContext, noArguments)
		})

	})
}

func _if() NamedValue {
	return NewGoFunction("if", func(context RunContext, arguments []Argument) Value {
		// if expects at least 2 arguments
		//
		arglen := len(arguments)
		if arglen < 2 {
			return NewErrorValue("invalid call to if, expect at least 2 parameters: usage if <condition> {...}")
		}

		condition := EvalArgument(context, arguments[0])
		if condition.Type() != TypeBoolean {
			return NewErrorValue("if condition does not evaluate to a boolean value")
		}

		if condition.(*booleanLiteral).value {
			return EvalArgumentWithBlock(context, arguments[1])
		}

		// condition not true, check else part
		//
		switch arglen {
		case 2:
			return Nothing
		case 3:
			return EvalArgumentWithBlock(context, arguments[2])
		case 4:
			if arguments[2].Value().String() == "else" {
				return EvalArgumentWithBlock(context, arguments[3])
			}
			return NewErrorValue("invalid call to if, expected else as 3rd argument")
		default:
			return NewErrorValue("invalid call to if, too many arguments")
		}

	})
}

func list() NamedValue {
	return NewGoFunction("list", func(context RunContext, arguments []Argument) Value {
		values := make([]Value, len(arguments))
		for i, arg := range arguments {
			values[i] = EvalArgument(context, arg)
		}
		return NewListValue(values)
	})
}

func dictWithBlock(context RunContext, arguments []Argument) Value {
	block := arguments[0]

	// use NewRunContext so block will be evaluated within same scope
	//
	subContext := NewRunContext(context)

	block.Value().(Block).Run(subContext, noArguments)

	return NewDictionaryValue(subContext.Mapping())
}

func dict() NamedValue {
	return NewGoFunction("dict", func(context RunContext, arguments []Argument) Value {
		mapping := make(map[string]Value)

		if len(arguments) == 1 {
			evaluated := EvalArgument(context, arguments[0])
			if evaluated.Type() == TypeBlock {
				return dictWithBlock(context, arguments)
			}

			if evaluated.Type() != TypeList {
				return NewErrorValue(fmt.Sprintf("dict needs a list as argument. Can not create dictionary from %v", evaluated))
			}

			values := evaluated.Internal().([]Value)

			if (len(values) % 2) != 0 {
				return NewErrorValue("dict can not create a dictionary from an odd number of elements")
			}

			var key Value

			for i, val := range values {
				if i%2 == 0 {
					key = val
				} else {
					mapping[key.String()] = val
				}
			}

		} else {

			if (len(arguments) % 2) != 0 {
				return NewErrorValue("dict can not create a dictionary from an odd number of elements")
			}

			var key Value

			for i, arg := range arguments {
				if i%2 == 0 {
					key = EvalArgument(context, arg)
				} else {
					mapping[key.String()] = EvalArgument(context, arg)
				}
			}
		}

		return NewDictionaryValue(mapping)
	})
}

func mixin() NamedValue {
	return NewGoFunction("mixin", func(context RunContext, arguments []Argument) Value {

		arglen := len(arguments)

		if arglen == 0 {
			return NewErrorValue("mixin expects something to mix in")
		}

		var dict Value
		for _, arg := range arguments {
			dict = EvalArgument(context, arg)
			if dict.Type() != TypeDictionary {
				return NewErrorValue(fmt.Sprintf("mixin can only mix in dictionaries, not %s", dict.String()))
			}

			for k, v := range dict.Internal().(map[string]Value) {
				context.Set(k, v)
			}
		}

		if arglen == 1 {
			return dict
		}

		return Nothing

	})
}

func puts() NamedValue {
	return NewGoFunction("puts", func(context RunContext, arguments []Argument) Value {
		for _, arg := range arguments {
			fmt.Printf("%s", EvalArgument(context, arg))
		}
		fmt.Printf("\n")
		return Nothing
	})
}

func load() NamedValue {
	return NewGoFunction("load", func(context RunContext, arguments []Argument) Value {

		if len(arguments) != 1 {
			return NewErrorValue("invalid call to load, expected 1 parameter: usage load <package name>")
		}

		name := EvalArgument2String(context, arguments[0])

		module, found := context.Module(name)

		if found {
			content := module.Content(context)
			return content
		}

		return NewErrorValue(fmt.Sprintf("could not find module %s", name))
	})
}

func eq() NamedValue {
	return NewGoFunction("eq", func(context RunContext, arguments []Argument) Value {

		argLen := len(arguments)

		if argLen != 2 {
			return NewErrorValue("invalid call to eq, expected exactly 2 parameters: usage eq <value> <value>")
		}

		if reflect.DeepEqual(EvalArgument(context, arguments[0]), EvalArgument(context, arguments[1])) {
			return True
		}

		return False

	})
}

func ne() NamedValue {
	return NewGoFunction("ne", func(context RunContext, arguments []Argument) Value {

		argLen := len(arguments)

		if argLen != 2 {
			return NewErrorValue("invalid call to eq, expected exactly 2 parameters: usage eq <value> <value>")
		}

		if reflect.DeepEqual(EvalArgument(context, arguments[0]), EvalArgument(context, arguments[1])) {
			return False
		}

		return True

	})
}
