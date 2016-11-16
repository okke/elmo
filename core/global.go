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
	context.SetNamed(defined())
	context.SetNamed(once())
	context.SetNamed(incr())
	context.SetNamed(_return())
	context.SetNamed(ampersand())
	context.SetNamed(_func())
	context.SetNamed(_if())
	context.SetNamed(while())
	context.SetNamed(until())
	context.SetNamed(do())
	context.SetNamed(mixin())
	context.SetNamed(load())
	context.SetNamed(eval())
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
	return NewGoFunction(`type/Get type information of a runtime value.
		Usage: type value

		Examples:

		> a:3
		> type a
		will result in identifier

		> type $a
		will result in int`,

		func(context RunContext, arguments []Argument) Value {
			_, ok, err := CheckArguments(arguments, 1, 1, "type", "<value>")
			if !ok {
				return err
			}
			return EvalArgument(context, arguments[0]).Info().Name()
		})
}

func set() NamedValue {
	return NewGoFunction(`set/Set a variable
		Usage: set <symbol> <value>
		Alternative usage: set <symbol>* value
		Returns: value that has been assigned to the denoted variable

		Examples:

		> set a 3
		> a
		will result in 3
		> set b (set c 3)
		> b
		will result in 3
		> c
		will result in 3

    Alternative example using a function that returns multiple values:
		> set f (func { return 1 2 })
		> set a b $f
		> a
		will result in 1
		> b
		will result in 2

		Note, instead of using set, it's possible to use the ':' shortcut like:
		> a: 3
		or
		> f: (func {....})`,

		func(context RunContext, arguments []Argument) Value {

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
	return NewGoFunction(`get/Gets a the value of a variable
		Usage: get <symbol>
		Returns: content of variable denoted by symbol

		Examples:

		> a: 3
		> get a
		will result in 3

		> b: (get a)
		> b
		will result in 3

		Note, the usage of get is most of the times unnecessary. using $ or (...)
		will do the same
		> a: 3
		> a
		> b: $a
		> c: (b)`,

		func(context RunContext, arguments []Argument) Value {

			_, ok, err := CheckArguments(arguments, 1, 1, "get", "<identifier>")
			if !ok {
				return err
			}

			identifier := EvalArgument(context, arguments[0])

			if identifier.Type() != TypeIdentifier {
				return NewErrorValue(fmt.Sprintf("can not get non identifier %v", arguments[0]))
			}

			_, result, found := identifier.(IdentifierValue).LookUp(context)

			if found {
				return result
			}

			return Nothing

		})
}

func defined() NamedValue {
	return NewGoFunction(`defined/Check if a variable is defined
		Usage: defined <symbol>
		Returns: true or false

		Examples:

		> a: 3
		> defined a
		will result in true
		> defined b
		will result in false

		Example combining it with assert:
		> assert (defined a)`,

		func(context RunContext, arguments []Argument) Value {

			_, ok, err := CheckArguments(arguments, 1, 1, "defined", "<identifier>")
			if !ok {
				return err
			}

			identifier := EvalArgument(context, arguments[0])

			if identifier.Type() != TypeIdentifier {
				return False
			}

			_, _, found := identifier.(IdentifierValue).LookUp(context)

			if found {
				return True
			}

			return False

		})
}

func once() NamedValue {
	return NewGoFunction(`once/Sets a variable only once
		Usage: once <symbol> <value>
		Returns value that was set

		Examples:

		> once a 1
		> once a 2
		> a
		will result in 1`,

		func(context RunContext, arguments []Argument) Value {

			_, ok, err := CheckArguments(arguments, 2, 2, "once", "<identifier> <value>")
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
	return NewGoFunction(`incr/Increments variable with 1 or given value
		Usage: incr <symbol> <value>?
		Returns: incremented value

		Examples:

		> a: 1
		> incr a
		> a
		will result in 2
		> incr a 3
		> a
		will result in 5

		Note, symbol must be an integer or floating point variable.
		> s: "chipotle"
		> incr s
		will result in an error`,

		func(context RunContext, arguments []Argument) Value {

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
	return NewGoFunction(`return/Stops processing and returns a value to caller
		Usage return <value>*
		Returns: returned value as value or as a multiple return value

		Examples:
		> f: (func {return 1})
		> f
		will result in 1
		> f: (func {return 1 2})
		> f
		will result in a multiple return value <[1 2]> that can be assigned like
		> set a b (f)
		> a
		will result in 1
		> b
		will result in 2`,

		func(context RunContext, arguments []Argument) Value {

			argLen := len(arguments)

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

func ampersand() NamedValue {
	return NewGoFunction("&/Internal function, can't be used",
		func(context RunContext, arguments []Argument) Value {
			// no argument checks are needed, ampersand (&) is an internal function
			// part of elmo's syntax
			//
			name := EvalArgument(context, arguments[0])

			if name.Type() == TypeIdentifier {
				_, value, found := name.(IdentifierValue).LookUp(context)
				if found {
					return value
				}
			}

			return NewErrorValue(fmt.Sprintf("could not resolve &%v", name))
		})
}

func _func() NamedValue {
	return NewGoFunction(`func/Create a new function
		Usage: func <symbol>* {...}
		Returns: a new function

		Given symbols denote function parameter names.

		Examples:

		> func a {...}
		will create a function that accepts one parameter called 'a'
		> func a { return $a }
		will create an echo function`,

		func(context RunContext, arguments []Argument) Value {

			argLen, ok, err := CheckArguments(arguments, 1, math.MaxInt16, "func", "<identifier>* {...}")
			if !ok {
				return err
			}

			argNamesAsArgument := arguments[0 : len(arguments)-1]
			block := EvalArgument(context, arguments[len(arguments)-1])
			argNames := make([]string, len(argNamesAsArgument))

			if argLen == 1 && block.Type() != TypeBlock {
				// block is not a block, maybe its an identifier that can be used
				// to lookup a function insted of creating on
				//
				result, found := context.Get(EvalArgument2String(context, arguments[len(arguments)-1]))
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

				return block.(Block).Run(subContext, NoArguments)
			})

		})
}

func _if() NamedValue {
	return NewGoFunction(`if/Conditionally execute (a block of) code
	  Usage: if <condition> {...} (else {...})?
		Returns: value of executed (block of) code

		Examples:

		> if (eq $a $b) "equal"
		> if (eq $a $b) {
		>   "equals"
		> }

		> if (eq $a $b) "equal" else "different"
		> if (eq $a $b) {
		>  ...
		> } else {
		>  ...
		> }

		Note, the result of a call to if can be assigned to a variable

		> e: (if (eq $a $b) "equal" else "different")
		> puts $e`,

		func(context RunContext, arguments []Argument) Value {

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
	return NewGoFunction(`while/Repeat (a block of) code while a given condition is true
	  Usage: while <condition> {...}
		Returns: result of given code or nil when code is not executed

		Examples:

		> a: 1
		> while (lt $a 10) { incr a }

		Note, the result of while can be assigned to a variable.

		>  c: (while (lt $a 10) { incr a; incr b })`,

		createLoop("while", true))
}

func until() NamedValue {
	return NewGoFunction(`until/Repeat (a block of) code until a given condition is true
	  Usage: until <condition> {...}
		Returns: result of given code  or nil when code is not executed

		Examples:

		> a: 1
		> until (eq $a 10) { incr a }

		Note, the result of while can be assigned to a variable.

		>  c: (until (eq $a 10) { incr a; incr b })`, createLoop("until", false))
}

func do() NamedValue {
	return NewGoFunction(`do/Special variant of while and until
		Usage: do {...} (while|until) <condition>
		Returns: value of executed code

		Example:

		> a: 1
		> do {puts $a; incr a} while (lt $a 10)
		> a: 1
		> do {puts $a; incr a} until (eq $a 10)`,

		func(context RunContext, arguments []Argument) Value {

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
	return NewGoFunction(`mixin/Mixin all key value pairs of a given dictionary as variables
		Usage: mixin <dictionary>*

		Examples:

		> d: {
		>  key: "value"
		>  value: "key"
		> }
		> mixin $d
		> key
		will result in value
		> value
		will result in key

		Mixin can be used in combination with load to load functions directly into
		current scope

		> mixin (load sys)

		Mixin can also be used to populate dictionaries

		> e: {
		>  mixin $d
		> }
		> e.key
		will result in "value"`,

		func(context RunContext, arguments []Argument) Value {

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

			}

			if argLen == 1 {
				return dict
			}

			return Nothing

		})
}

func puts() NamedValue {
	return NewGoFunction(`puts/Write values to stdout
		Usage puts <value>*
		Returns: nil

		Examples:
		> puts "chipotle"
		will print chipotle
		> puts 1 2 3
		will print 123`,

		func(context RunContext, arguments []Argument) Value {
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
	return NewGoFunction(`echo/Returns given value
		Usage: echo <value>
		Returns: given value

		Examples:

		> a: 1
		> b: (echo a)
		is the same as
		> b: $a
		and also the same as
		> b: (a)`,

		func(context RunContext, arguments []Argument) Value {
			_, ok, err := CheckArguments(arguments, 1, 1, "echo", "<value>")
			if !ok {
				return err
			}

			return EvalArgument(context, arguments[0])
		})
}

func sleep() NamedValue {
	return NewGoFunction(`sleep/Pause for given number of milliseconds
		Usage: sleep <number>
		Returns: nil

		Example:

		> sleep 1000
		will pause for one second`,

		func(context RunContext, arguments []Argument) Value {

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
	return NewGoFunction(`load/Loads given module or script
		Usage: load <module|script>
		Returns: A dictionary containing loaded variables

		Examples:

		> str: (load string)
		> str.len "chipotle"

		> helper: (load "include/functions")

		Last example will load the script 'incude/functions.mo' that should be
		located relatively from the current script

		Note, the ".mo" extension is implied and should not be specified`,

		func(context RunContext, arguments []Argument) Value {

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

func eval() NamedValue {
	return NewGoFunction(`eval/Evaluate a block of code
		Usage: eval <dict>? <block>
		Returns: evaluation result

		Example:

		> eval {
		>   puts "in block"
	  > }

		Example with additional variables:

		> context: {a:1; b:2}
		> eval $context {
		>   puts $a "," $b
	  > }

		Eval is escpecially handy to execute blocks that are passed to a function

		> pepper: (func block {
	 	>   return (eval $block)
	 	> })

		> pepper {
		>  "chipotle"
		> }
		`,

		func(context RunContext, arguments []Argument) Value {

			argLen, ok, err := CheckArguments(arguments, 1, 2, "eval", "<dict>? <block>")
			if !ok {
				return err
			}

			var blockContext = context.CreateSubContext()
			var blockArg = 0

			if argLen == 2 {
				blockArg = 1
				dict := EvalArgument(context, arguments[0])
				if dict.Type() == TypeBlock {
					dict = NewDictionaryWithBlock(context, dict.(Block))
				}
				blockContext.Mixin(dict)

				// ensure mixed in key/value pairs are not overriding
				// local context
				//
				blockContext = blockContext.CreateSubContext()
			}

			result := EvalArgumentWithBlock(blockContext, arguments[blockArg])
			if result.Type() == TypeBlock {
				return result.(Block).Run(blockContext, []Argument{})
			}
			return result

		})
}

func eq() NamedValue {
	return NewGoFunction(`eq/Checks if two arguments are the same
		Usage: eq <value> <value>
		Returns: true (equal) or false (differ)

		Examples:

		> eq 1 1
		will result in true
		> eq 1 2
		will result in false

		> sauce: (func {return nil})
		> eq $nil $sauce
		will result in true`,

		func(context RunContext, arguments []Argument) Value {

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
	return NewGoFunction(`ne/Checks if two arguments are different
		Usage: ne <value> <value>
		Returns: true (differ) or false (same)

		Examples:

		> ne 1 2
		will result in true
		> ne 1 1
		will result in false

		> sauce: (func {return $nil})
		> ne $nil $sauce
		will result in false`, func(context RunContext, arguments []Argument) Value {

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
	return NewGoFunction(`gt/Checks if first argument is greater than second argument
		Usage: gt <value> <value>
		Returns: true (greater than) or false (less or equal)

		Examples:

		> gt 2 1
		will result in true
		> gt 2 2
		will result in false`,

		func(context RunContext, arguments []Argument) Value {

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
	return NewGoFunction(`gte/Checks if first argument is greater than or equal to second argument
		Usage: gte <value> <value>
		Returns: true (greater than or equal) or false (less than)

		Examples:

		> gte 2 1
		will result in true
		> gte 2 2
		will result in true`, func(context RunContext, arguments []Argument) Value {

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
	return NewGoFunction(`lt/Checks if first argument is less than second argument
		Usage: lt <value> <value>
		Returns: true (less than) or false (greater or equal)

		Examples:

		> lt 1 2
		will result in true
		> lt 2 2
		will result in false`, func(context RunContext, arguments []Argument) Value {

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
	return NewGoFunction(`lte/Checks if first argument is greater than or equal to second argument
		Usage: lte <value> <value>
		Returns: true (less than or equal) or false (greater than)

		Examples:

		> lte 2 1
		will result in true
		> lte 2 2
		will result in true`, func(context RunContext, arguments []Argument) Value {

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
	return NewGoFunction(`and/Logical and operation on multiple boolean values
		Usage: and <boolean>*
		Returns: true (all arguments are true) or false (at least one argument is false)

		Examples:

		> and (eq 1 1) (eq 2 2)
		will result in true
		> and (eq 1 1) (eq 1 2)
		will result in false
		> and $true $true $false
		will result in false

		Note, 'and' uses lazy evaluation: as soon as an argument is false, it will
		stop evaluating and it will return false.

		> and $false (eval {puts soep; true})
		will do nothing and will result in false
		> and $true (eval {puts "chipotle"; true})
		will evaluate the second argument and will write chipotle to stdout`,

		func(context RunContext, arguments []Argument) Value {

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
	return NewGoFunction(`or/Logical or operation on multiple boolean values
		Usage: or <boolean>*
		Returns: true (at least one arguments is true) or false (all arguments are false)

		Examples:

		> or (eq 1 2) (eq 2 2)
		will result in true
		> or (eq 1 3) (eq 1 2)
		will result in false
		> or $true $true $false
		will result in true

		Note, 'or' uses lazy evaluation: as soon as an argument is true, it will
		stop evaluating and it will return true.

		> or $true (eval {puts soep; true})
		will do nothing and will result in true
		> or $false (eval {puts "chipotle"; true})
		will evaluate the second argument and will write chipotle to stdout`,

		func(context RunContext, arguments []Argument) Value {

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
	return NewGoFunction(`not/reverses a boolean value
		Usage: not <boolean>
		Returns inverted boolean value

		Examples:

		> not $true
		will result in false
		> not (eq 1 2)
		will result in true`, func(context RunContext, arguments []Argument) Value {

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
	return NewGoFunction(`plus/Add two numbers
		Usage: plus <number> <number>
		Returns: the sum of the two numbers

		Examples:

		> plus 1 2
		will result in 3
		> plus 1.0 2.0
		will result in 3.000000`,

		func(context RunContext, arguments []Argument) Value {

			return arithmeticOperation(context, arguments, "plus", func(v1 Value, v2 Value) Value {
				return v1.(MathValue).Plus(v2)
			})

		})
}

func minus() NamedValue {
	return NewGoFunction(`minus/Subtracts two numbers
		Usage: minus <number> <number>
		Returns: the subtraction of the two numbers

		Examples:

		> minus 2 1
		will result in 1
		> minus 2.0 1.0
		will result in 1.000000`,

		func(context RunContext, arguments []Argument) Value {

			return arithmeticOperation(context, arguments, "minus", func(v1 Value, v2 Value) Value {
				return v1.(MathValue).Minus(v2)
			})

		})
}

func multiply() NamedValue {
	return NewGoFunction(`multiply/Multiplies two numbers
		Usage: multiply <number> <number>
		Returns: the product of the two numbers

		Examples:

		> multiply 3 4
		will result in 12
		> multiply 0.5 10
		will result in 5.000000`,

		func(context RunContext, arguments []Argument) Value {

			return arithmeticOperation(context, arguments, "multiply", func(v1 Value, v2 Value) Value {
				return v1.(MathValue).Multiply(v2)
			})

		})
}

func divide() NamedValue {
	return NewGoFunction(`divide/divides two numbers
		Usage: devide <number> <number>
		Returns: the division result

		Examples:

		> divide 6 3
		will result in 2
		> divide 7 3
		will result in 2
		> divide 7.0 3
		will result in 2.333333

		Note, dividing by zero will result in an error`,

		func(context RunContext, arguments []Argument) Value {

			return arithmeticOperation(context, arguments, "divide", func(v1 Value, v2 Value) Value {
				return v1.(MathValue).Divide(v2)
			})

		})
}

func modulo() NamedValue {
	return NewGoFunction(`modulo/calculates the remainder of a division
		Usage: modulo <number> <integer>
		Returns: the division result

		Examples:

		> modulo 6 3
		will result in 0
		> modulo 7 3
		will result in 1
		> divide 7.5 3
		will result in 1.500000

		Note the second argument must be a non 0 integer, otherwise an error will be returned`,

		func(context RunContext, arguments []Argument) Value {

			return arithmeticOperation(context, arguments, "modulo", func(v1 Value, v2 Value) Value {
				return v1.(MathValue).Modulo(v2)
			})

		})
}

func assert() NamedValue {
	return NewGoFunction(`assert/Evaluate a boolean value and return an error when false
		Usage: assert <boolean> <error>?
		Returns: true or an error

		Examples:

		> assert $false
		will result in an error
		> assert $false "Told you so"
		will result in an error with "Told you so" as error message
		> assert (defined a)
		will result in an error when a is not defined`,

		func(context RunContext, arguments []Argument) Value {

			argLen, ok, err := CheckArguments(arguments, 1, 2, "assert", "<boolean> <error>?")
			if !ok {
				return err
			}

			check := EvalArgument(context, arguments[0])
			if check != nil && check.Type() == TypeBoolean {
				if check.(*booleanLiteral).value {
					return True
				}

				if argLen == 2 {
					return NewErrorValue(EvalArgument2String(context, arguments[1]))
				}

				return NewErrorValue("assertion failed")

			}

			return NewErrorValue("assert: first argument does not evaluate to a boolean value")
		})
}

func _error() NamedValue {
	return NewGoFunction(`error/constructs an error with a user defined message
		Usage: error <message>
		Returns: a user defined error

		Example:
		> error "chipotle not hot enought error"
		will result in an error

		Note, givens message value is evaluated to a string`,

		func(context RunContext, arguments []Argument) Value {

			_, ok, err := CheckArguments(arguments, 1, 1, "error", "<message>")
			if !ok {
				return err
			}

			return NewErrorValue(EvalArgument2String(context, arguments[0]))

		})
}

func help() NamedValue {
	return NewGoFunction("help/Get help. Usage 'help' or 'help identifier'", func(context RunContext, arguments []Argument) Value {

		argLen, ok, err := CheckArguments(arguments, 0, 1, "help", "")
		if !ok {
			return err
		}

		// get help for a specific function
		//
		if argLen == 1 {

			identifier := EvalArgument(context, arguments[0])

			if identifier.Type() == TypeIdentifier {

				_, result, found := identifier.(IdentifierValue).LookUp(context)

				if found {
					help, ok := result.(HelpValue)
					if ok {
						return help.Help()
					}
					return result
				}
			}

			// no help available
			//
			return Nothing

		}

		keys := context.Keys()
		sort.Strings(keys)

		result := []Value{}
		for _, key := range keys {
			value, _ := context.Get(key)
			if value.Type() == TypeGoFunction {
				result = append(result, NewIdentifier(key))
			} else if value.Type() == TypeDictionary {

				subkeys := []string{}
				for k := range value.Internal().(map[string]Value) {
					subkeys = append(subkeys, k)
				}

				sort.Strings(subkeys)

				for _, subkey := range subkeys {
					subValue, _ := value.Internal().(map[string]Value)[subkey]
					if subValue.Type() == TypeGoFunction {
						result = append(result, NewNameSpacedIdentifier([]string{key, subkey}))
					}
				}
			}
		}

		return NewListValue(result)

	})
}
