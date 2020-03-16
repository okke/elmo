package elmo

import (
	"errors"
	"fmt"
)

type call struct {
	astNode
	baseValue
	firstArgument Argument
	function      GoFunction
	arguments     []Argument
	pipe          Call
}

// Call is a function call
//
type Call interface {
	Value
	Runnable

	Name() string
	Arguments() []Argument
	WillPipe() bool
}

func (call *call) Name() string {
	if call.function != nil {
		return fmt.Sprintf("%v", call.function)
	}
	return call.firstArgument.String()
}

func (call *call) Arguments() []Argument {
	return call.arguments
}

func (call *call) WillPipe() bool {
	return call.pipe != nil
}

func (call *call) addInfoWhenError(value Value) Value {
	if value == nil {
		return nil
	}
	if value.Type() == TypeError {
		if value.(ErrorValue).IsTraced() {
			// TODO: add trace??
			//
			return value
		}

		lineno, _ := call.meta.PositionOf(int(call.BeginsAt()))
		value.(ErrorValue).SetAt(call.meta, lineno)
	}
	return value
}

func (call *call) pipeResult(context RunContext, value Value) Value {
	if !call.WillPipe() {
		return value
	}

	if value.Type() == TypeReturn {
		values := value.(*returnValue).values
		arguments := make([]Argument, len(values))
		for i, v := range values {
			arguments[i] = &argument{value: v}
		}
		return call.pipe.Run(context, arguments)
	}

	return call.pipe.Run(context, []Argument{&argument{value: value}})
}

func createArgumentsForMissingFunc(context RunContext, call *call, arguments []Argument) []Argument {
	// pass evaluated arguments to the 'func missing' function
	// as a list of values
	//
	values := make([]Value, len(arguments))
	for i, value := range arguments {
		values[i] = EvalArgument(context, value)
	}

	// and pass the original function name as first argument
	//
	return []Argument{
		NewArgument(call.meta, call.astNode.node, NewIdentifier(call.firstArgument.Value().(*identifier).value[len(call.firstArgument.Value().(*identifier).value)-1])),
		NewArgument(call.meta, call.astNode.node, NewListValue(values))}
}

func (call *call) Run(context RunContext, additionalArguments []Argument) Value {

	if call.function != nil {
		return call.pipeResult(context, call.addInfoWhenError(call.function(context, call.Arguments())))
	}

	var inDict DictionaryValue
	var value Value
	var found bool
	var useArguments []Argument

	var function IdentifierValue

	switch call.firstArgument.Type() {
	case TypeCall:
		value = call.firstArgument.Value().(Runnable).Run(context, []Argument{})
		if value.Type() == TypeIdentifier {
			function = value.(IdentifierValue)
			inDict, value, found = function.LookUp(context)
		}
		found = true
	case TypeIdentifier:
		function = call.firstArgument.Value().(IdentifierValue)
		inDict, value, found = function.LookUp(context)
	case TypeString:
		value = call.firstArgument.Value().(StringValue).ResolveBlocks(context)
		found = true
	default:
		value = call.firstArgument.Value()
		found = true
	}

	if additionalArguments != nil && len(additionalArguments) > 0 {
		useArguments = append([]Argument{}, additionalArguments...)
		useArguments = append(useArguments, call.arguments...)
	} else {
		useArguments = call.arguments
	}

	// when call can not be resolved, try to find the 'func missing' function
	//
	if !found {
		if inDict == nil {
			value, found = context.Get("?")
		} else {
			value, found = inDict.Resolve("?")
		}

		if found {
			useArguments = createArgumentsForMissingFunc(context, call, useArguments)
		}
	}

	if found {

		if value == nil {
			return call.pipeResult(context, call.addInfoWhenError(NewErrorValue(fmt.Sprintf("call to %s results in invalid nil value", call.Name()))))
		}

		if inDict != nil {
			this := context.This()
			context.SetThis(inDict.(Value))
			defer func() {
				if this == nil {
					context.SetThis(nil)
				} else {
					context.SetThis(this.(Value))
				}
			}()
		}

		if value.Type() == TypeGoFunction {
			return call.pipeResult(context, call.addInfoWhenError(value.(Runnable).Run(context, useArguments)))
		}

		// runnable values can be used as functions to access their content
		//
		runnable, isRunnable := value.(Runnable)
		if (isRunnable) && (len(call.arguments) > 0) {
			return call.pipeResult(context, call.addInfoWhenError(runnable.Run(context, useArguments)))
		}

		return call.pipeResult(context, call.addInfoWhenError(value))
	}

	return call.pipeResult(context, call.addInfoWhenError(NewErrorValue(fmt.Sprintf("call to undefined \"%s\"", call.firstArgument))))
}

func (call *call) String() string {
	return fmt.Sprintf("(%s ...)", call.Name())
}

func (call *call) Type() Type {
	return TypeCall
}

func (call *call) Internal() interface{} {
	return errors.New("Internal() not implemented on call")
}

func (call *call) Enrich(dict DictionaryValue) {

	if call.pipe != nil {
		dict.Set(NewStringLiteral("pipe"), call.pipe)
	}

	dict.Set(NewStringLiteral("name"), NewStringLiteral(call.Name()))
}

// NewCall contstructs a new function call
//
func NewCall(meta ScriptMetaData, node *node32, firstArg Argument, arguments []Argument, pipeTo Call) Call {
	return &call{astNode: astNode{meta: meta, node: node}, baseValue: baseValue{info: typeInfoCall},
		firstArgument: firstArg, arguments: arguments, pipe: pipeTo}
}

// NewCallWithFunction constructs a call that does not need to be resolved
//
func NewCallWithFunction(meta ScriptMetaData, node *node32, function GoFunction, arguments []Argument, pipeTo Call) Call {
	return &call{astNode: astNode{meta: meta, node: node}, baseValue: baseValue{info: typeInfoCall},
		function: function, arguments: arguments, pipe: pipeTo}
}
