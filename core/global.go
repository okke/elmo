package elmo

import (
	"bytes"
	"fmt"
	"math"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
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
	context.SetNamed(let())
	context.SetNamed(get())
	context.SetNamed(first())
	context.SetNamed(defined())
	context.SetNamed(once())
	context.SetNamed(incr())
	context.SetNamed(_return())
	context.SetNamed(ampersand())
	context.SetNamed(_func())
	context.SetNamed(template())
	context.SetNamed(_if())
	context.SetNamed(while())
	context.SetNamed(until())
	context.SetNamed(do())
	context.SetNamed(mixin())
	context.SetNamed(load())
	context.SetNamed(eval())
	context.SetNamed(parse())
	context.SetNamed(puts())
	context.SetNamed(echo())
	context.SetNamed(toS())
	context.SetNamed(sleep())
	context.SetNamed(eq())
	context.SetNamed(ne())
	context.SetNamed(gt())
	context.SetNamed(gte())
	context.SetNamed(lt())
	context.SetNamed(lte())
	context.SetNamed(compare())
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
	context.SetNamed(_panic())
	context.SetNamed(help())
	context.SetNamed(_close())
	context.SetNamed(_len())
	context.SetNamed(freeze())
	context.SetNamed(frozen())
	context.SetNamed(_uuid())
	context.SetNamed(_time())
	context.SetNamed(file())
	context.SetNamed(tempFile())
	context.SetNamed(test())
	context.SetNamed(globalSettings())
	context.SetNamed(elmoVersion())

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
func CheckArguments(arguments []Argument, min int, max int, fname string, usage string) (int, ErrorValue) {
	argLen := len(arguments)
	if argLen < min || argLen > max {
		return argLen, NewErrorValue(fmt.Sprintf("Invalid call to %s. Usage: %s %s", fname, fname, usage))
	}
	return argLen, nil
}

func _type() NamedValue {
	return NewGoFunctionWithHelp("type", `Get type information of a runtime value.
		Usage: type value

		Examples:

		> a:3
		> type a
		will result in identifier

		> type $a
		will result in int`,

		func(context RunContext, arguments []Argument) Value {
			_, err := CheckArguments(arguments, 1, 1, "type", "<value>")
			if err != nil {
				return err
			}
			value := EvalArgument(context, arguments[0])
			info := value.Info()
			if info == nil {
				return NewStringLiteral("?")
			}
			return info.Name()
		})
}

func setOrLet(convertBlockToDictionary bool, name, help string) NamedValue {
	return NewGoFunctionWithHelp(name, help, func(context RunContext, arguments []Argument) Value {

		argLen, err := CheckArguments(arguments, 2, math.MaxInt16, name, "<identifier>* value")
		if err != nil {
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

			_, err := CheckArguments(arguments, 2, 2, "set", "<identifier> value")
			if err != nil {
				return err
			}

			// convert block to dictionary
			//
			if convertBlockToDictionary && value.Type() == TypeBlock {
				value = NewDictionaryWithBlock(context, value.(Block))
			}

			name := EvalArgument2String(context, arguments[0])
			context.Set(name, value)
		}

		if value.Type() == TypeError {
			if !value.(ErrorValue).IsFatal() {
				// can ignore non fatal errors in assignments
				//
				return value.(ErrorValue).Ignore()
			}
		}
		return value
	})
}

func set() NamedValue {
	return setOrLet(true, "set", `set/Set a variable
		Usage: set <symbol> <value>
		Alternative usage: set <symbol>* value
		Returns: value that has been assigned to the denoted variable

		When evaluation of the assigned value results in a non fatal error,
		the assignment will result in an error that won't break the execution flow.

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
		> f: (func {....})
		
		When assigning a block of code to a variable, the block of code will be executed and
		the result will be a dictionary with values.
		> set peppers {jalapeno:"hot"; habanero: "hotter"}
		`)
}

func let() NamedValue {
	return setOrLet(false, "let", `Set a variable including assigning blocks of code as block of code instead of a dictionary value
		Usage: let <symbol> <value>
		Alternative usage: let <symbol>* value
		Returns: value that has been assigned to the denoted variable`)
}

func get() NamedValue {
	return NewGoFunctionWithHelp("get", `Gets a the value of a variable
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

			_, err := CheckArguments(arguments, 1, 1, "get", "<identifier>")
			if err != nil {
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

func first() NamedValue {
	return NewGoFunctionWithHelp("first", `Gets a the first value of given arguments
		Usage: first <symbol>*
		Returns: first value and ignore other parameters
		

		Examples:

		> get 1 2 3
		will result in 1`,

		func(context RunContext, arguments []Argument) Value {
			if len(arguments) == 0 {
				return Nothing
			}
			return EvalArgument(context, arguments[0])
		})
}

func defined() NamedValue {
	return NewGoFunctionWithHelp("defined", `Check if a variable is defined
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

			_, err := CheckArguments(arguments, 1, 1, "defined", "<identifier>")
			if err != nil {
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
	return NewGoFunctionWithHelp("once", `Sets a variable only once
		Usage: once <symbol> <value>
		Returns value that was set

		Examples:

		> once a 1
		> once a 2
		> a
		will result in 1`,

		func(context RunContext, arguments []Argument) Value {

			_, err := CheckArguments(arguments, 2, 2, "once", "<identifier> <value>")
			if err != nil {
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
	return NewGoFunctionWithHelp("incr", `Increments variable with 1 or given value
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

			argLen, err := CheckArguments(arguments, 1, 2, "incr", "<identifier> <value>?")
			if err != nil {
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

			_, ok := incrValue.(IncrementableValue)
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
	return NewGoFunctionWithHelp("return", `Stops processing and returns a value to caller
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
	return NewGoFunctionWithHelp("&", "Internal function, can't be used",
		func(context RunContext, arguments []Argument) Value {
			// no argument checks are needed, ampersand (&) is an internal function
			// part of elmo's syntax

			// when argument is a string, return the raw string
			// including possible blocks which are not evaluated
			// so resulting value can be used as template
			//
			if arguments[0].Value().Type() == TypeString {

				// the only issue is the current context should be captured
				// so blocks inside the string are evaluted in the right context
				//
				str := arguments[0].Value().(StringValue)
				return str.CopyWithinContext(context)
			}

			// Otherwise evaluate to identifier which resolved and returned
			// without further evaluation so functions can be passed as arguments
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

func _if() NamedValue {
	return NewGoFunctionWithHelp("if", `Conditionally execute (a block of) code
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

			argLen, err := CheckArguments(arguments, 2, math.MaxInt16, "if", "<condition> {...} (else {...})?")
			if err != nil {
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

		_, err := CheckArguments(arguments, 2, 2, name, "<condition> {...}")
		if err != nil {
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
	return NewGoFunctionWithHelp("while", `Repeat (a block of) code while a given condition is true
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
	return NewGoFunctionWithHelp("until", `Repeat (a block of) code until a given condition is true
	  Usage: until <condition> {...}
		Returns: result of given code  or nil when code is not executed

		Examples:

		> a: 1
		> until (eq $a 10) { incr a }

		Note, the result of while can be assigned to a variable.

		>  c: (until (eq $a 10) { incr a; incr b })`, createLoop("until", false))
}

func do() NamedValue {
	return NewGoFunctionWithHelp("do", `Special variant of while and until
		Usage: do {...} (while|until) <condition>
		Returns: value of executed code

		Example:

		> a: 1
		> do {puts $a; incr a} while (lt $a 10)
		> a: 1
		> do {puts $a; incr a} until (eq $a 10)`,

		func(context RunContext, arguments []Argument) Value {

			_, err := CheckArguments(arguments, 3, 3, "do", "{} while|until <condition>")
			if err != nil {
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
	return NewGoFunctionWithHelp("mixin", `Mixin all key value pairs of a given dictionary as variables
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

			argLen, err := CheckArguments(arguments, 0, math.MaxInt16, "mixin", "<value>*")
			if err != nil {
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
	return NewGoFunctionWithHelp("puts", `Write values to stdout
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
	return NewGoFunctionWithHelp("echo", `Returns given value
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
			_, err := CheckArguments(arguments, 1, 1, "echo", "<value>")
			if err != nil {
				return err
			}

			return EvalArgument(context, arguments[0])
		})
}

func toS() NamedValue {
	return NewGoFunctionWithHelp("to_s", `Converts given value to a string
		Usage: to_s <value>
		Returns: string representation of value

		Examples:

		> a: 1
		> b: (to_s a)
		`,

		func(context RunContext, arguments []Argument) Value {
			_, err := CheckArguments(arguments, 1, 1, "to_s", "<value>")
			if err != nil {
				return err
			}

			return NewStringLiteral(EvalArgument(context, arguments[0]).String())
		})
}

func sleep() NamedValue {
	return NewGoFunctionWithHelp("sleep", `Pause for given number of milliseconds
		Usage: sleep <number>
		Returns: nil

		Example:

		> sleep 1000
		will pause for one second`,

		func(context RunContext, arguments []Argument) Value {

			_, err := CheckArguments(arguments, 1, 1, "sleep", "<number>")
			if err != nil {
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
	return NewGoFunctionWithHelp("load", `Loads given module or script
		Usage: load <module|script>
		Returns: A dictionary containing loaded variables

		When loading results in an error, it will return it as a fatal error

		Examples:

		> str: (load string)
		> str.len "chipotle"

		> helper: (load "include/functions")

		Last example will load the script 'incude/functions.mo' that should be
		located relatively from the current script

		Note, the ".mo" extension is implied and should not be specified`,

		func(context RunContext, arguments []Argument) Value {

			_, err := CheckArguments(arguments, 1, 1, "load", "<package name>")
			if err != nil {
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

			if loaded.Type() == TypeError {
				return loaded.(ErrorValue).Panic()
			}

			return loaded
		})
}

func eval() NamedValue {
	return NewGoFunctionWithHelp("eval", `Evaluate a block of code
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

			argLen, err := CheckArguments(arguments, 1, 2, "eval", "<dict>? <block>")
			if err != nil {
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
			} else if result.Type() == TypeString {
				return result.(StringValue).ResolveBlocks(blockContext)
			}
			return result

		})
}

func parse() NamedValue {
	return NewGoFunctionWithHelp("parse", `Parses a string of elmo code and returns a block of code
		Usage: parse string
		Returns: elmo code block
		`,

		func(context RunContext, arguments []Argument) Value {

			if _, err := CheckArguments(arguments, 1, 1, "parse", "<string>"); err != nil {
				return err
			}

			block, err := Parse2Block(EvalArgument2String(context, arguments[0]), "parse")
			if err != nil {
				return NewErrorValue(err.Error())
			}
			return block
		})
}

func compareValues(context RunContext, v1 Value, v2 Value, f func(int) Value) Value {
	c1, comparable := v1.(ComparableValue)
	if !comparable {
		return NewErrorValue(fmt.Sprintf("invalid comparison, expected comparable values instead of %v and %v", v1, v2))
	}

	result, err := c1.Compare(context, v2)
	if err != nil {
		return err
	}
	return f(result)

}

func eq() NamedValue {
	return NewGoFunctionWithHelp("eq", `Checks if two arguments are the same
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

			_, err := CheckArguments(arguments, 2, 2, "eq", "<value> <value>")
			if err != nil {
				return err
			}

			v1 := EvalArgument(context, arguments[0])
			v2 := EvalArgument(context, arguments[1])

			// first try to compare the two values
			//
			if result := compareValues(context, v1, v2, func(result int) Value {
				if result == 0 {
					return True
				}
				return False
			}); result.Type() != TypeError {
				return result
			}

			// if that did not work, simple do a deep equal
			//
			if reflect.DeepEqual(EvalArgument(context, arguments[0]), EvalArgument(context, arguments[1])) {
				return True
			}

			return False

		})
}

func ne() NamedValue {
	return NewGoFunctionWithHelp("ne", `Checks if two arguments are different
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

		_, err := CheckArguments(arguments, 2, 2, "ne", "<value> <value>")
		if err != nil {
			return err
		}

		v1 := EvalArgument(context, arguments[0])
		v2 := EvalArgument(context, arguments[1])

		// first try to compare the two values
		//
		if result := compareValues(context, v1, v2, func(result int) Value {
			if result == 0 {
				return False
			}
			return True
		}); result.Type() != TypeError {
			return result
		}

		if reflect.DeepEqual(EvalArgument(context, arguments[0]), EvalArgument(context, arguments[1])) {
			return False
		}

		return True

	})
}

func gt() NamedValue {
	return NewGoFunctionWithHelp("gt", `Checks if first argument is greater than second argument
		Usage: gt <value> <value>
		Returns: true (greater than) or false (less or equal)

		Examples:

		> gt 2 1
		will result in true
		> gt 2 2
		will result in false`,

		func(context RunContext, arguments []Argument) Value {

			_, err := CheckArguments(arguments, 2, 2, "gt", "<value> <value>")
			if err != nil {
				return err
			}

			v1 := EvalArgument(context, arguments[0])
			v2 := EvalArgument(context, arguments[1])

			return compareValues(context, v1, v2, func(result int) Value {
				if result == 1 {
					return True
				}
				return False
			})

		})
}

func gte() NamedValue {
	return NewGoFunctionWithHelp("gte", `Checks if first argument is greater than or equal to second argument
		Usage: gte <value> <value>
		Returns: true (greater than or equal) or false (less than)

		Examples:

		> gte 2 1
		will result in true
		> gte 2 2
		will result in true`, func(context RunContext, arguments []Argument) Value {

		_, err := CheckArguments(arguments, 2, 2, "gte", "<value> <value>")
		if err != nil {
			return err
		}

		v1 := EvalArgument(context, arguments[0])
		v2 := EvalArgument(context, arguments[1])

		return compareValues(context, v1, v2, func(result int) Value {
			if result == -1 {
				return False
			}
			return True
		})

	})
}

func lt() NamedValue {
	return NewGoFunctionWithHelp("lt", `Checks if first argument is less than second argument
		Usage: lt <value> <value>
		Returns: true (less than) or false (greater or equal)

		Examples:

		> lt 1 2
		will result in true
		> lt 2 2
		will result in false`, func(context RunContext, arguments []Argument) Value {

		_, err := CheckArguments(arguments, 2, 2, "lt", "<value> <value>")
		if err != nil {
			return err
		}

		v1 := EvalArgument(context, arguments[0])
		v2 := EvalArgument(context, arguments[1])

		return compareValues(context, v1, v2, func(result int) Value {
			if result == -1 {
				return True
			}
			return False
		})

	})
}

func lte() NamedValue {
	return NewGoFunctionWithHelp("lte", `Checks if first argument is greater than or equal to second argument
		Usage: lte <value> <value>
		Returns: true (less than or equal) or false (greater than)

		Examples:

		> lte 2 1
		will result in true
		> lte 2 2
		will result in true`, func(context RunContext, arguments []Argument) Value {

		_, err := CheckArguments(arguments, 2, 2, "lte", "<value> <value>")
		if err != nil {
			return err
		}

		v1 := EvalArgument(context, arguments[0])
		v2 := EvalArgument(context, arguments[1])

		return compareValues(context, v1, v2, func(result int) Value {
			if result == 1 {
				return False
			}
			return True
		})

	})
}

func compare() NamedValue {
	return NewGoFunctionWithHelp("compare", `Compares two values 
		Usage: comnpare <value> <value>
		Returns: true -1 (first value is less then second), 0 (values are equal) or 1 (first value is greater than second value)

		Examples:

		> compare 2 1
		will result in 1
		> compare [1 2 3] [1 2 3 4]
		will result in -1`, func(context RunContext, arguments []Argument) Value {

		_, err := CheckArguments(arguments, 2, 2, "lte", "<value> <value>")
		if err != nil {
			return err
		}

		v1 := EvalArgument(context, arguments[0])
		v2 := EvalArgument(context, arguments[1])

		return compareValues(context, v1, v2, func(result int) Value {
			return NewIntegerLiteral(int64(result))
		})

	})
}

func and() NamedValue {
	return NewGoFunctionWithHelp("and", `Logical and operation on multiple boolean values
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
	return NewGoFunctionWithHelp("or", `Logical or operation on multiple boolean values
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
	return NewGoFunctionWithHelp("not", `reverses a boolean value
		Usage: not <boolean>
		Returns inverted boolean value

		Examples:

		> not $true
		will result in false
		> not (eq 1 2)
		will result in true`, func(context RunContext, arguments []Argument) Value {

		_, err := CheckArguments(arguments, 1, 1, "not", "<boolean>")
		if err != nil {
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
	_, err := CheckArguments(arguments, 2, 2, name, "<value> <value>")
	if err != nil {
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
	return NewGoFunctionWithHelp("plus", `Add two numbers
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
	return NewGoFunctionWithHelp("minus", `Subtracts two numbers
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
	return NewGoFunctionWithHelp("multiply", `Multiplies two numbers
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
	return NewGoFunctionWithHelp("divide", `divides two numbers
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
	return NewGoFunctionWithHelp("modulo", `calculates the remainder of a division
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
	return NewGoFunctionWithHelp("assert", `Evaluate a boolean value and return an error when false
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

			argLen, err := CheckArguments(arguments, 1, 2, "assert", "<boolean> <error>?")
			if err != nil {
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
	return NewGoFunctionWithHelp("error", `constructs an error with a user defined message
		Usage: error <message>
		Returns: a user defined error

		Example:
		> error "chipotle not hot enought error"
		will result in an error

		Note, given message value is evaluated to a string`,

		func(context RunContext, arguments []Argument) Value {

			_, err := CheckArguments(arguments, 1, 1, "error", "<message>")
			if err != nil {
				return err
			}

			return NewErrorValue(EvalArgument2String(context, arguments[0]))

		})
}

func _panic() NamedValue {
	return NewGoFunctionWithHelp("panic", `constructs a fatal error with a user defined message
		Usage: panic <message>
		Returns: a user defined fatal error

		Example:
		> fatal "chipotle not hot enought error"
		will result in an error

		Note, given message value is evaluated to a string

		Also note, fatal error can't be assigned to variables`,

		func(context RunContext, arguments []Argument) Value {

			_, err := CheckArguments(arguments, 1, 1, "panic", "<message>")
			if err != nil {
				return err
			}

			return NewErrorValue(EvalArgument2String(context, arguments[0])).Panic()

		})
}

func formatHelp(s string) string {
	splitted := strings.Split(s, "\n")
	var buf bytes.Buffer
	for i, v := range splitted {
		buf.WriteString(strings.Trim(v, " \t"))
		if i < (len(splitted) - 1) {
			buf.WriteString("\n")
		}
	}
	return buf.String()
}

func help() NamedValue {
	return NewGoFunctionWithHelp("help", "Get help. Usage 'help' or 'help symbol'", func(context RunContext, arguments []Argument) Value {

		argLen, err := CheckArguments(arguments, 0, 1, "help", "<symbol>?")
		if err != nil {
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
						return NewStringLiteral(formatHelp(help.Help().String()))
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

func _close() NamedValue {
	return NewGoFunctionWithHelp("close", `Closes a value (in case this is supported by given value)
		Usage: close <value>
		Returns: value`,

		func(context RunContext, arguments []Argument) Value {

			_, err := CheckArguments(arguments, 1, 1, "close", "<value>")
			if err != nil {
				return err
			}

			value := EvalArgument(context, arguments[0])
			if closeable, ok := value.(CloseableValue); ok {
				closeable.Close()
			} else {

				// also check if internal value implements the CloseableValue interface
				//
				if closeable, ok := value.Internal().(CloseableValue); ok {
					closeable.Close()
				} else {
					return NewErrorValue("value is not closable")
				}
			}

			return value
		})
}

func _len() NamedValue {
	return NewGoFunctionWithHelp("len",
		`Determine length of given value
		Usage: len <value>
		Returns: int`,

		func(context RunContext, arguments []Argument) Value {

			_, err := CheckArguments(arguments, 1, 1, "len", "<value>")
			if err != nil {
				return err
			}

			value := EvalArgument(context, arguments[0])
			if withLength, ok := value.(ValueWithLength); ok {
				return withLength.Length()
			}

			return NewErrorValue(fmt.Sprintf("can not determine length of variable with type %s", value.Info().Name()))

		})
}

func freeze() NamedValue {
	return NewGoFunctionWithHelp("freeze!", `Freezes a value (makes a value immutable)
		Usage: freeze <value>
		Returns: frozen value`,

		func(context RunContext, arguments []Argument) Value {

			_, err := CheckArguments(arguments, 1, 1, "freeze!", "<value>")
			if err != nil {
				return err
			}

			value := EvalArgument(context, arguments[0])
			if freezable, ok := value.(FreezableValue); ok {
				return freezable.Freeze()
			}

			return value
		})
}

func frozen() NamedValue {
	return NewGoFunctionWithHelp("frozen", `Checks if a value is frozen (immutable)
		Usage: frozen <value>
		Returns: boolean (true when frozen, false when not)`,

		func(context RunContext, arguments []Argument) Value {

			_, err := CheckArguments(arguments, 1, 1, "frozen", "<value>")
			if err != nil {
				return err
			}

			value := EvalArgument(context, arguments[0])
			if freezable, ok := value.(FreezableValue); ok {
				return TrueOrFalse(freezable.Frozen())
			}

			// by default, all non freezable values are frozen
			//
			return True
		})
}

func _uuid() NamedValue {
	return NewGoFunctionWithHelp("uuid", `Returns a unique id for given value or a new id when no value is given
		Usage: uuid <value>?
		Returns: uuid`,

		func(context RunContext, arguments []Argument) Value {
			argLen, err := CheckArguments(arguments, 0, 1, "uuid", "<value>?")
			if err != nil {
				return err
			}

			if argLen == 1 {
				return NewStringLiteral(EvalArgument(context, arguments[0]).UUID().String())
			}

			return NewStringLiteral(uuid.New().String())
		})
}

// TimeDictionary takes a go time and convert into an elmo dictionary
//
func TimeDictionary(t time.Time) Value {
	_, zoneOffset := t.Zone()
	return NewDictionaryValue(nil, map[string]Value{
		"zoneOffset": NewIntegerLiteral(int64(zoneOffset)),
		"year":       NewIntegerLiteral(int64(t.Year())),
		"month":      NewIntegerLiteral(int64(t.Month())),
		"day":        NewIntegerLiteral(int64(t.Day())),
		"hour":       NewIntegerLiteral(int64(t.Hour())),
		"minute":     NewIntegerLiteral(int64(t.Minute())),
		"second":     NewIntegerLiteral(int64(t.Second())),
		"nano":       NewIntegerLiteral(int64(t.Nanosecond())),
		"timestamp":  NewIntegerLiteral(t.UnixNano())})
}

func _time() NamedValue {
	return NewGoFunctionWithHelp("time", `Generate a time dictionary based on given input
		Usage: 
		> time // without arguments time will return the current time
		> time <int> // with one integer argument, time will convert given timestamp to time dictionary
		> time <string> // with one string argument, time will convert given timestamp according to RFC3339 to time dictionary
		> time <format> <string> // with two arguments, time will convert given string according to given format to time dictionary

		Supported formats: ANSIC UnixDate RubyDate RFC822 RFC822Z RFC850 RFC1123 RFC1123Z RFC3339 RFC3339Nano Kitchen

		Returns: dictionary with time values: {
			zoneOffset (integer, nr of seconds )
			year (integer)
			month (integer)
			day (integer)
			hour (integer)
			minute (integer)
			second (integer)
			nano (integer)
			timestamp (integer, unix timestamp)
		}
		
		Example:
		
		parsedTime: (time Kitchen "3:04PM")
		anotherTime: (time RFC1123Z "Mon, 02 Jan 2006 15:04:05 -0700")
		currentTime: $time
		fromTimestamp: (time 0) 
		`,

		func(context RunContext, arguments []Argument) Value {
			argLen, err := CheckArguments(arguments, 0, 2, "time", "<format>? <string>?")
			if err != nil {
				return err
			}

			if argLen == 0 {
				// currrent time
				return TimeDictionary(time.Now())
			}

			var format string = time.RFC3339
			var timestr string = ""

			if argLen == 2 {

				formatstr := EvalArgument2String(context, arguments[0])
				switch formatstr {
				case "ANSIC":
					format = time.ANSIC
				case "UnixDate":
					format = time.UnixDate
				case "RubyDate":
					format = time.RubyDate
				case "RFC822":
					format = time.RFC822
				case "RFC822Z":
					format = time.RFC822Z
				case "RFC850":
					format = time.RFC850
				case "RFC1123":
					format = time.RFC1123
				case "RFC1123Z":
					format = time.RFC1123Z
				case "RFC3339":
					format = time.RFC3339
				case "RFC3339Nano":
					format = time.RFC3339Nano
				case "Kitchen":
					format = time.Kitchen
				default:
					format = formatstr
				}
				timestr = EvalArgument2String(context, arguments[1])

			} else {

				timearg := EvalArgument(context, arguments[0])
				if timearg.Type() == TypeInteger {
					// timestamp
					return TimeDictionary(time.Unix(0, timearg.Internal().(int64)))
				}
				timestr = EvalArgument2String(context, arguments[0])
			}

			time, timeerr := time.Parse(format, timestr)
			if timeerr != nil {
				return NewErrorValue(timeerr.Error())
			}
			return TimeDictionary(time)
		})
}

func test() NamedValue {
	return NewGoFunctionWithHelp("test", `Runs all test functions in a given dictionary
	`, func(context RunContext, arguments []Argument) Value {
		if _, err := CheckArguments(arguments, 1, 1, "test", "<suite>"); err != nil {
			return err
		}

		// first argument of a dictionary function can be an identifier with the name of the dictionary
		//
		dict, ok := EvalArgumentOrSolveIdentifier(context, arguments[0]).(DictionaryValue)

		if !ok {
			return NewErrorValue("invalid call to test, expect a dictionary with test functions as first argument: usage test <suite>")
		}

		keyNames := dict.Keys()

		results := make(map[string]Value, 0)
		failed := false

		for _, key := range keyNames {
			if !strings.HasPrefix(key, "test") {
				continue
			}

			value, _ := dict.Resolve(key)
			if value.Type() != TypeGoFunction {
				continue
			}

			result := value.(Runnable).Run(context, []Argument{})
			if result.Type() == TypeError {
				results[key] = result
				failed = true
			}
		}

		if failed {
			return NewErrorValue(fmt.Sprintf("test suite failed: %v", results))
		}

		return True
	})
}
