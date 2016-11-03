package elmo

import (
	"fmt"
	"math"
	"reflect"
	"sort"
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

	context.SetNamed(_type())
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
	context.SetNamed(mixin())
	context.SetNamed(load())
	context.SetNamed(puts())
	context.SetNamed(echo())
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
	context.SetNamed(plus())
	context.SetNamed(minus())
	context.SetNamed(multiply())
	context.SetNamed(divide())
	context.SetNamed(modulo())
	context.SetNamed(assert())
	context.SetNamed(_error())
	context.SetNamed(help())

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
			value = NewDictionaryWithBlock(context, value.(Block))
		}

		values[i] = value
	}
	return NewListValue(values)
}

// CheckArguments checks number of arguments
//
func CheckArguments(arguments []Argument, min int, max int, fname string, usage string) (int, bool, ErrorValue) {
	argLen := len(arguments)
	if argLen < min || argLen > max {
		return argLen, false, NewErrorValue(fmt.Sprintf("Invalid call to %s. Usage: %s %s", fname, fname, usage))
	}
	return argLen, true, nil
}

func _type() NamedValue {
	return NewGoFunction("type", func(context RunContext, arguments []Argument) Value {
		_, ok, err := CheckArguments(arguments, 1, 1, "type", "<value>")
		if !ok {
			return err
		}
		return EvalArgument(context, arguments[0]).Info().Name()
	})
}

func set() NamedValue {
	return NewGoFunction("set", func(context RunContext, arguments []Argument) Value {

		argLen, ok, err := CheckArguments(arguments, 2, math.MaxInt16, "set", "<identifier>* value")
		if !ok {
			return err
		}

		var value = EvalArgument(context, arguments[argLen-1])

		// value can evaluate to a multiple return so will result
		// in multiple assignments
		if value.Type() == TypeReturn {
			returnedValues := value.(*returnValue).values
			returnedLength := len(returnedValues)

			for i := 0; i < (argLen-1) && i < returnedLength; i++ {
				name := EvalArgument2String(context, arguments[i])
				context.Set(name, returnedValues[i])
			}
		} else {

			_, ok, err := CheckArguments(arguments, 2, 2, "set", "<identifier> value")
			if !ok {
				return err
			}

			// convert block to dictionary
			//
			if value.Type() == TypeBlock {
				value = NewDictionaryWithBlock(context, value.(Block))
			}

			name := EvalArgument2String(context, arguments[0])
			context.Set(name, value)
		}

		return value
	})
}

func get() NamedValue {
	return NewGoFunction("get", func(context RunContext, arguments []Argument) Value {

		_, ok, err := CheckArguments(arguments, 1, 1, "get", "<identifier>")
		if !ok {
			return err
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

		_, ok, err := CheckArguments(arguments, 2, 2, "once", "once <identifier>")
		if !ok {
			return err
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

		argLen, ok, err := CheckArguments(arguments, 1, 2, "incr", "<identifier> <value>?")
		if !ok {
			return err
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

			_, isIncrementable := currentValue.(IncrementableValue)
			if isIncrementable {

				newValue := currentValue.(IncrementableValue).Increment(incrValue)

				if arg0.Type() == TypeIdentifier {
					context.Set(arg0.String(), newValue)
				}

				return newValue
			}

			return NewErrorValue("invalid call to incr, expected variable that can be incremented")

		}

		_, ok = incrValue.(IncrementableValue)
		if !ok {
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

		argLen, ok, err := CheckArguments(arguments, 1, math.MaxInt16, "func", "<identifier>* {...}")
		if !ok {
			return err
		}

		argNamesAsArgument := arguments[0 : len(arguments)-1]
		block := arguments[len(arguments)-1]
		argNames := make([]string, len(argNamesAsArgument))

		if argLen == 1 && block.Type() != TypeBlock {
			// block is not a block, maybe its an identifier that can be used
			// to lookup a function insted of creating on
			//
			result, found := context.Get(EvalArgument2String(context, block))
			if found && result.Type() == TypeGoFunction {
				return result
			}
		}

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

		argLen, ok, err := CheckArguments(arguments, 2, math.MaxInt16, "if", "<condition> {...} (else {...})?")
		if !ok {
			return err
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
		switch argLen {
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

		_, ok, err := CheckArguments(arguments, 2, 2, name, "<condition> {...}")
		if !ok {
			return err
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
				if result.Type() == TypeError {
					return result
				}
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

		_, ok, err := CheckArguments(arguments, 3, 3, "do", "{} while|until <condition>")
		if !ok {
			return err
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
			if result.Type() == TypeError {
				return result
			}
		}

	})
}

func mixin() NamedValue {
	return NewGoFunction("mixin", func(context RunContext, arguments []Argument) Value {

		argLen, ok, err := CheckArguments(arguments, 0, math.MaxInt16, "mixin", "<value>*")
		if !ok {
			return err
		}

		var dict Value
		for _, arg := range arguments {
			dict = EvalArgument(context, arg)

			result := context.Mixin(dict)

			if result.Type() == TypeError {
				return result
			}

			for k, v := range dict.Internal().(map[string]Value) {
				context.Set(k, v)
			}
		}

		if argLen == 1 {
			return dict
		}

		return Nothing

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

func echo() NamedValue {
	return NewGoFunction("echo", func(context RunContext, arguments []Argument) Value {
		_, ok, err := CheckArguments(arguments, 1, 1, "echo", "<value>")
		if !ok {
			return err
		}

		return EvalArgument(context, arguments[0])
	})
}

func sleep() NamedValue {
	return NewGoFunction("sleep", func(context RunContext, arguments []Argument) Value {

		_, ok, err := CheckArguments(arguments, 1, 1, "sleep", "<number>")
		if !ok {
			return err
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

		_, ok, err := CheckArguments(arguments, 1, 1, "load", "<package name>")
		if !ok {
			return err
		}

		name := EvalArgument2String(context, arguments[0])

		module, found := context.Module(name)

		if found {
			content := module.Content(context)
			return content
		}

		loader := NewLoader(context, []string{})

		loaded := loader.Load(name)

		if loaded == nil {
			return NewErrorValue(fmt.Sprintf("could not find module %s", name))
		}

		return context.Mixin(loaded)
	})
}

func eq() NamedValue {
	return NewGoFunction("eq", func(context RunContext, arguments []Argument) Value {

		_, ok, err := CheckArguments(arguments, 2, 2, "eq", "<value> <value>")
		if !ok {
			return err
		}

		if reflect.DeepEqual(EvalArgument(context, arguments[0]), EvalArgument(context, arguments[1])) {
			return True
		}

		return False

	})
}

func ne() NamedValue {
	return NewGoFunction("ne", func(context RunContext, arguments []Argument) Value {

		_, ok, err := CheckArguments(arguments, 2, 2, "ne", "<value> <value>")
		if !ok {
			return err
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

		_, ok, err := CheckArguments(arguments, 2, 2, "gt", "<value> <value>")
		if !ok {
			return err
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

		_, ok, err := CheckArguments(arguments, 2, 2, "gte", "<value> <value>")
		if !ok {
			return err
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

		_, ok, err := CheckArguments(arguments, 2, 2, "lt", "<value> <value>")
		if !ok {
			return err
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

		_, ok, err := CheckArguments(arguments, 2, 2, "lte", "<value> <value>")
		if !ok {
			return err
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

		_, ok, err := CheckArguments(arguments, 1, 1, "not", "<boolean>")
		if !ok {
			return err
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

func arithmeticOperation(context RunContext, arguments []Argument, name string, f func(Value, Value) Value) Value {
	_, ok, err := CheckArguments(arguments, 2, 2, name, "<value> <value>")
	if !ok {
		return err
	}

	v1 := EvalArgument(context, arguments[0])
	_, ok1 := v1.(MathValue)
	if !ok1 {
		return NewErrorValue(fmt.Sprintf("%s does not support %v", name, v1))
	}

	v2 := EvalArgument(context, arguments[1])
	_, ok2 := v2.(MathValue)
	if !ok2 {
		return NewErrorValue(fmt.Sprintf("%s does not support %v", name, v2))
	}

	return f(v1, v2)
}

func plus() NamedValue {
	return NewGoFunction("plus", func(context RunContext, arguments []Argument) Value {

		return arithmeticOperation(context, arguments, "plus", func(v1 Value, v2 Value) Value {
			return v1.(MathValue).Plus(v2)
		})

	})
}

func minus() NamedValue {
	return NewGoFunction("minus", func(context RunContext, arguments []Argument) Value {

		return arithmeticOperation(context, arguments, "minus", func(v1 Value, v2 Value) Value {
			return v1.(MathValue).Minus(v2)
		})

	})
}

func multiply() NamedValue {
	return NewGoFunction("multiply", func(context RunContext, arguments []Argument) Value {

		return arithmeticOperation(context, arguments, "multiply", func(v1 Value, v2 Value) Value {
			return v1.(MathValue).Multiply(v2)
		})

	})
}

func divide() NamedValue {
	return NewGoFunction("divide", func(context RunContext, arguments []Argument) Value {

		return arithmeticOperation(context, arguments, "divide", func(v1 Value, v2 Value) Value {
			return v1.(MathValue).Divide(v2)
		})

	})
}

func modulo() NamedValue {
	return NewGoFunction("modulo", func(context RunContext, arguments []Argument) Value {

		return arithmeticOperation(context, arguments, "modulo", func(v1 Value, v2 Value) Value {
			return v1.(MathValue).Modulo(v2)
		})

	})
}

func assert() NamedValue {
	return NewGoFunction("assert", func(context RunContext, arguments []Argument) Value {

		_, ok, err := CheckArguments(arguments, 2, 2, "assert", "<boolean> <error>")
		if !ok {
			return err
		}

		check := EvalArgument(context, arguments[0])
		if check != nil && check.Type() == TypeBoolean {
			if check.(*booleanLiteral).value {
				return True
			}

			return NewErrorValue(EvalArgument2String(context, arguments[1]))
		}

		return NewErrorValue("assert: first argument does not evaluate to a boolean value")
	})
}

func _error() NamedValue {
	return NewGoFunction("error", func(context RunContext, arguments []Argument) Value {

		_, ok, err := CheckArguments(arguments, 1, 1, "error", "<message>")
		if !ok {
			return err
		}

		return NewErrorValue(EvalArgument2String(context, arguments[0]))

	})
}

func help() NamedValue {
	return NewGoFunction("help", func(context RunContext, arguments []Argument) Value {

		_, ok, err := CheckArguments(arguments, 0, 0, "help", "")
		if !ok {
			return err
		}

		keys := context.Keys()
		sort.Strings(keys)

		result := []Value{}
		for _, key := range keys {
			value, _ := context.Get(key)
			if value.Type() == TypeGoFunction {
				result = append(result, NewStringLiteral(key))
			} else if value.Type() == TypeDictionary {

				subkeys := []string{}
				for k := range value.Internal().(map[string]Value) {
					subkeys = append(subkeys, k)
				}

				sort.Strings(subkeys)

				for _, subkey := range subkeys {
					subValue, _ := value.Internal().(map[string]Value)[subkey]
					if subValue.Type() == TypeGoFunction {
						result = append(result, NewStringLiteral(key+"."+subkey))
					}
				}
			}
		}

		return NewListValue(result)

	})
}
