package elmo

import (
	"fmt"
	"reflect"
	"time"
)

// NoArguments is an array of arguments with no arguments
//
var NoArguments = []Argument{}

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
	context.SetNamed(once())
	context.SetNamed(incr())
	context.SetNamed(_return())
	context.SetNamed(_func())
	context.SetNamed(_if())
	context.SetNamed(while())
	context.SetNamed(until())
	context.SetNamed(do())
	context.SetNamed(dict())
	context.SetNamed(mixin())
	context.SetNamed(new())
	context.SetNamed(load())
	context.SetNamed(puts())
	context.SetNamed(sleep())
	context.SetNamed(eq())
	context.SetNamed(ne())
	context.SetNamed(gt())
	context.SetNamed(gte())
	context.SetNamed(lt())
	context.SetNamed(lte())
	context.SetNamed(and())
	context.SetNamed(or())
	context.SetNamed(not())

	return context
}

// ListConstructor constructs a list at runtime
//
func ListConstructor(context RunContext, arguments []Argument) Value {
	values := make([]Value, len(arguments))
	for i, arg := range arguments {
		var value = EvalArgument(context, arg)

		// accept blocks within a list
		// as dictionaries in order to support
		// [{...} {...} ...] constructions
		//
		if value.Type() == TypeBlock {
			value = dictWithBlock(context, value.(Block))
		}

		values[i] = value
	}
	return NewListValue(values)
}

func set() NamedValue {
	return NewGoFunction("set", func(context RunContext, arguments []Argument) Value {

		argLen := len(arguments)

		if len(arguments) < 2 {
			return NewErrorValue("invalid call to set, expect at least 2 parameters: usage set <identifier> <value>")
		}

		var value = EvalArgument(context, arguments[argLen-1])

		if value.Type() == TypeReturn {
			returnedValues := value.(*returnValue).values
			returnedLength := len(returnedValues)

			for i := 0; i < (argLen-1) && i < returnedLength; i++ {
				name := EvalArgument2String(context, arguments[i])
				context.Set(name, returnedValues[i])
			}
		} else {

			// set expects exactly 2 arguments
			//
			if len(arguments) != 2 {
				return NewErrorValue("invalid call to set, expected 2 parameters: usage set <identifier> <value>")
			}

			// convert block to dictionary
			//
			if value.Type() == TypeBlock {
				value = dictWithBlock(context, value.(Block))
			}

			name := EvalArgument2String(context, arguments[0])
			context.Set(name, value)
		}

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

func once() NamedValue {
	return NewGoFunction("once", func(context RunContext, arguments []Argument) Value {

		argLen := len(arguments)

		if argLen != 2 {
			return NewErrorValue("invalid call to setOnce: usage new <identifier> <value>")
		}

		name := EvalArgument2String(context, arguments[0])

		existing, found := context.Get(name)
		if !found {
			existing = EvalArgument(context, arguments[1])
			context.Set(name, existing)
		}

		return existing
	})
}

func incr() NamedValue {
	return NewGoFunction("incr", func(context RunContext, arguments []Argument) Value {

		argLen := len(arguments)

		if argLen < 1 || argLen > 2 {
			return NewErrorValue("invalid call to incr, expected 1 or 2 parameters: usage incr <identifier> <value>?")
		}

		arg0 := EvalArgument(context, arguments[0])

		var incrValue = One
		if argLen == 2 {
			incrValue = EvalArgument(context, arguments[1])
		}

		var currentValue Value
		var found bool
		if arg0.Type() == TypeIdentifier {
			currentValue, found = context.Get(arg0.String())
		} else {
			currentValue, found = arg0, true
		}

		if found {

			if currentValue.Type() == TypeInteger || currentValue.Type() == TypeFloat {
				newValue := currentValue.(IncrementableValue).Increment(incrValue)

				if arg0.Type() == TypeIdentifier {
					context.Set(arg0.String(), newValue)
				}

				return newValue
			}

			return NewErrorValue("invalid call to incr, expected integer variable")

		}

		incrType := reflect.TypeOf(incrValue)
		shouldBe := reflect.TypeOf((*IncrementableValue)(nil)).Elem()

		if !incrType.Implements(shouldBe) {
			return NewErrorValue("invalid call to incr, expected a value that can be incremented")
		}

		// not found so set it to initial value
		//
		context.Set(arg0.String(), incrValue)
		return incrValue

	})
}

func _return() NamedValue {
	return NewGoFunction("return", func(context RunContext, arguments []Argument) Value {

		argLen := len(arguments)

		// return expects exactly 1 argument
		//
		var result Value

		switch argLen {
		case 0:
			result = Nothing
		case 1:
			result = EvalArgument(context, arguments[0])
		default:
			values := make([]Value, len(arguments))
			for i, arg := range arguments {
				values[i] = EvalArgument(context, arg)
			}
			result = NewReturnValue(values)
		}

		context.Stop()

		return result
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
				return NewErrorValue(fmt.Sprintf("invalid call to user defined function: expect %d parameters instead of %d", len(argNames), len(innerArguments)))
			}

			cloneFrom := innerContext
			if cloneFrom.Parent() != nil {
				cloneFrom = cloneFrom.Parent()
			}
			subContext := cloneFrom.CreateSubContext()

			if innerContext.This() != nil {
				subContext.Set("this", innerContext.This())
			}

			for i, v := range innerArguments {
				subContext.Set(argNames[i], EvalArgument(innerContext, v))
			}

			return block.Value().(Block).Run(subContext, NoArguments)
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

func createLoop(name string, stopCondition bool) func(context RunContext, arguments []Argument) Value {
	return func(context RunContext, arguments []Argument) Value {

		// if expects exactly 2 arguments
		//
		arglen := len(arguments)
		if arglen != 2 {
			return NewErrorValue(fmt.Sprintf("invalid call to %s, expect 2 parameters: usage %s <condition> {...}", name, name))
		}

		var result Value
		result = Nothing
		for {
			condition := EvalArgument(context, arguments[0])
			if condition.Type() != TypeBoolean {
				return NewErrorValue("condition does not evaluate to a boolean value")
			}
			if condition.(*booleanLiteral).value == stopCondition {
				result = EvalArgumentWithBlock(context, arguments[1])
			} else {
				return result
			}
		}
	}
}

func while() NamedValue {
	return NewGoFunction("while", createLoop("while", true))
}

func until() NamedValue {
	return NewGoFunction("until", createLoop("until", false))
}

func do() NamedValue {
	return NewGoFunction("do", func(context RunContext, arguments []Argument) Value {
		// do { } [while|until] condition

		// if expects exactly 2 arguments
		//
		arglen := len(arguments)
		if arglen != 3 {
			return NewErrorValue("invalid call to do, expect 3 parameters: usage do {...} while|until <condition>")
		}

		result := EvalArgumentWithBlock(context, arguments[0])

		// while => true, until => false
		//
		var stopCondition bool
		switch arguments[1].Value().String() {
		case "while":
			stopCondition = true
		case "until":
			stopCondition = false
		default:
			return NewErrorValue("expected while or until condition in do loop")
		}

		for {
			condition := EvalArgument(context, arguments[2])
			if condition.Type() != TypeBoolean {
				return NewErrorValue("condition does not evaluate to a boolean value")
			}
			if !(condition.(*booleanLiteral).value == stopCondition) {
				return result
			}
			result = EvalArgumentWithBlock(context, arguments[0])
		}

	})
}

func dictWithBlock(context RunContext, block Block) Value {

	// use NewRunContext so block will be evaluated within same scope
	//
	subContext := NewRunContext(context)

	block.Run(subContext, NoArguments)

	return NewDictionaryValue(nil, subContext.Mapping())
}

func dict() NamedValue {
	return NewGoFunction("dict", func(context RunContext, arguments []Argument) Value {
		mapping := make(map[string]Value)

		if len(arguments) == 1 {
			evaluated := EvalArgument(context, arguments[0])
			if evaluated.Type() == TypeBlock {
				return dictWithBlock(context, evaluated.(Block))
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

		return NewDictionaryValue(nil, mapping)
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

func new() NamedValue {
	return NewGoFunction("new", func(context RunContext, arguments []Argument) Value {
		argLen := len(arguments)

		if (argLen < 1) || (argLen > 2) {
			return NewErrorValue("new expect exactly one parameter. usage: new <dictionary> <block>?")
		}

		parent := EvalArgument(context, arguments[0])

		if parent.Type() != TypeDictionary {
			return NewErrorValue(fmt.Sprintf("new expects a dictionary, not %s", parent.String()))
		}

		var mapping map[string]Value
		if argLen == 2 {

			dict := EvalArgument(context, arguments[1])

			if dict.Type() != TypeBlock {
				return NewErrorValue(fmt.Sprintf("new can not construct dictionary from %s", dict.String()))
			}

			dict = dictWithBlock(context, dict.(Block))

			mapping = dict.(*dictValue).values
		} else {
			mapping = make(map[string]Value)
		}

		return NewDictionaryValue(parent.(*dictValue), mapping)
	})
}

func puts() NamedValue {
	return NewGoFunction("puts", func(context RunContext, arguments []Argument) Value {
		if len(arguments) == 1 {
			fmt.Printf("%s\n", EvalArgument(context, arguments[0]))
		} else {
			line := ""
			for _, arg := range arguments {
				line = fmt.Sprintf("%s%s", line, EvalArgument(context, arg))
			}
			fmt.Printf("%s\n", line)
		}

		return Nothing
	})
}

func sleep() NamedValue {
	return NewGoFunction("sleep", func(context RunContext, arguments []Argument) Value {

		if len(arguments) != 1 {
			return NewErrorValue("invalid call to sleep, expected 1 parameter: usage sleep <milliseconds>")
		}
		duration := EvalArgument(context, arguments[0])

		if duration.Type() != TypeInteger {
			return NewErrorValue("invalid call to sleep, expected integer parameter: usage sleep <milliseconds>")
		}

		sleepTime := time.Duration(duration.(*integerLiteral).value)
		time.Sleep(time.Millisecond * sleepTime)

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
			return NewErrorValue("invalid call to ne, expected exactly 2 parameters: usage ne <value> <value>")
		}

		if reflect.DeepEqual(EvalArgument(context, arguments[0]), EvalArgument(context, arguments[1])) {
			return False
		}

		return True

	})
}

func compareValues(v1 Value, v2 Value, f func(int) Value) Value {
	if v1.Type() == TypeInteger || v1.Type() == TypeFloat {
		result, err := v1.(ComparableValue).Compare(v2)
		if err != nil {
			return err
		}
		return f(result)
	}
	return NewErrorValue(fmt.Sprintf("invalid comparison, expected number values instead of %v and %v", v1, v2))
}

func gt() NamedValue {
	return NewGoFunction("gt", func(context RunContext, arguments []Argument) Value {

		argLen := len(arguments)

		if argLen != 2 {
			return NewErrorValue("invalid call to gt, expected exactly 2 parameters: usage gt <number> <number>")
		}

		v1 := EvalArgument(context, arguments[0])
		v2 := EvalArgument(context, arguments[1])

		return compareValues(v1, v2, func(result int) Value {
			if result == 1 {
				return True
			}
			return False
		})

	})
}

func gte() NamedValue {
	return NewGoFunction("gte", func(context RunContext, arguments []Argument) Value {

		argLen := len(arguments)

		if argLen != 2 {
			return NewErrorValue("invalid call to gte, expected exactly 2 parameters: usage gte <number> <number>")
		}

		v1 := EvalArgument(context, arguments[0])
		v2 := EvalArgument(context, arguments[1])

		return compareValues(v1, v2, func(result int) Value {
			if result == -1 {
				return False
			}
			return True
		})

	})
}

func lt() NamedValue {
	return NewGoFunction("lt", func(context RunContext, arguments []Argument) Value {

		argLen := len(arguments)

		if argLen != 2 {
			return NewErrorValue("invalid call to lt, expected exactly 2 parameters: usage lt <number> <number>")
		}

		v1 := EvalArgument(context, arguments[0])
		v2 := EvalArgument(context, arguments[1])

		return compareValues(v1, v2, func(result int) Value {
			if result == -1 {
				return True
			}
			return False
		})

	})
}

func lte() NamedValue {
	return NewGoFunction("lte", func(context RunContext, arguments []Argument) Value {

		argLen := len(arguments)

		if argLen != 2 {
			return NewErrorValue("invalid call to lte, expected exactly 2 parameters: usage lte <number> <number>")
		}

		v1 := EvalArgument(context, arguments[0])
		v2 := EvalArgument(context, arguments[1])

		return compareValues(v1, v2, func(result int) Value {
			if result == 1 {
				return False
			}
			return True
		})

	})
}

func and() NamedValue {
	return NewGoFunction("and", func(context RunContext, arguments []Argument) Value {

		argLen := len(arguments)

		for i := 0; i < argLen; i++ {
			condition := EvalArgument(context, arguments[i])
			if condition.Type() != TypeBoolean {
				return NewErrorValue("and condition does not evaluate to a boolean value")
			}

			if !condition.(*booleanLiteral).value {
				return False
			}
		}

		return True

	})
}

func or() NamedValue {
	return NewGoFunction("or", func(context RunContext, arguments []Argument) Value {

		argLen := len(arguments)

		for i := 0; i < argLen; i++ {
			condition := EvalArgument(context, arguments[i])
			if condition.Type() != TypeBoolean {
				return NewErrorValue("or condition does not evaluate to a boolean value")
			}

			if condition.(*booleanLiteral).value {
				return True
			}
		}

		return False

	})
}

func not() NamedValue {
	return NewGoFunction("not", func(context RunContext, arguments []Argument) Value {

		argLen := len(arguments)

		if argLen != 1 {
			return NewErrorValue("invalid call to not, expected exactly 1 parameters: usage not <boolean>")
		}

		condition := EvalArgument(context, arguments[0])
		if condition.Type() != TypeBoolean {
			return NewErrorValue("not condition does not evaluate to a boolean value")
		}

		if condition.(*booleanLiteral).value {
			return False
		}

		return True
	})
}
