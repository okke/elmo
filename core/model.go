package elmo

import (
	"errors"
	"fmt"
	"strings"
)

// Runnable is a type that can be interpreted
//
type Runnable interface {
	Run(RunContext, []Argument) Value
}

//
// ---[LITERALS]---------------------------------------------------------------
//

// Type represents an internal value type
//
type Type uint8

const (
	// TypeIdentifier represents a type for an identifier value
	TypeIdentifier Type = iota
	// TypeString represents a type for a string value
	TypeString
	// TypeInteger represents a type for an integer value
	TypeInteger
	// TypeBoolean represents a type for a boolean value
	TypeBoolean
	// TypeList represents a type for an array value
	TypeList
	// TypeDictionary represents a type for a map value
	TypeDictionary
	// TypeError represents a type for an error value
	TypeError
	// TypeBlock represents a type for a code block
	TypeBlock
	// TypeCall represent the type for a function call
	TypeCall
	// TypeGoFunction represents a type for an internal go function
	TypeGoFunction
	// TypeReturn represents a function result containing multiple values
	TypeReturn
	// TypeNil represents the type of a nil value
	TypeNil
)

type nothing struct {
}

// Nothing represents nil
//
var Nothing = &nothing{}

// True represents boolean value true
//
var True = NewBooleanLiteral(true)

// False represents boolean value false
//
var False = NewBooleanLiteral(false)

type identifier struct {
	value string
}

type stringLiteral struct {
	value string
}

type integerLiteral struct {
	value int64
}

type booleanLiteral struct {
	value bool
}

type listValue struct {
	values []Value
}

type returnValue struct {
	values []Value
}

type dictValue struct {
	parent *dictValue
	values map[string]Value
}

type errorValue struct {
	meta   ScriptMetaData
	lineno int
	msg    string
}

// GoFunction is a native go function that takes an array of input values
// and returns an output value
//
type GoFunction func(RunContext, []Argument) Value

type goFunction struct {
	name  string
	value GoFunction
}

// Value represents data within elmo
//
type Value interface {
	Print() string
	String() string
	Type() Type
	Internal() interface{}
}

// ErrorValue represents an Error
//
type ErrorValue interface {
	Value
	Error() string
	SetAt(meta ScriptMetaData, lineno int)
	At() (ScriptMetaData, int)
	IsTraced() bool
}

// RunnableValue represents a value that can evaluated to another value
//
type RunnableValue interface {
	Value
	Runnable
}

// NamedValue represent data with a name
//
type NamedValue interface {
	Value
	Name() string
}

func (nothing *nothing) Print() string {
	return "nil"
}

func (nothing *nothing) String() string {
	return "nil"
}

func (nothing *nothing) Type() Type {
	return TypeNil
}

func (nothing *nothing) Internal() interface{} {
	return nil
}

func (identifier *identifier) Print() string {
	return identifier.value
}

func (identifier *identifier) String() string {
	return identifier.value
}

func (identifier *identifier) Type() Type {
	return TypeIdentifier
}

func (identifier *identifier) Internal() interface{} {
	return identifier.value
}

func (stringLiteral *stringLiteral) Print() string {
	return fmt.Sprintf("\"%s\"", stringLiteral.value)
}

func (stringLiteral *stringLiteral) String() string {
	return fmt.Sprintf("%s", stringLiteral.value)
}

func (stringLiteral *stringLiteral) Type() Type {
	return TypeString
}

func (stringLiteral *stringLiteral) Internal() interface{} {
	return stringLiteral.value
}

func (integerLiteral *integerLiteral) Print() string {
	return fmt.Sprintf("%d", integerLiteral.value)
}

func (integerLiteral *integerLiteral) String() string {
	return fmt.Sprintf("%d", integerLiteral.value)
}

func (integerLiteral *integerLiteral) Type() Type {
	return TypeInteger
}

func (integerLiteral *integerLiteral) Internal() interface{} {
	return integerLiteral.value
}

func (booleanLiteral *booleanLiteral) Print() string {
	return booleanLiteral.String()
}

func (booleanLiteral *booleanLiteral) String() string {
	return fmt.Sprintf("%v", booleanLiteral.value)
}

func (booleanLiteral *booleanLiteral) Type() Type {
	return TypeBoolean
}

func (booleanLiteral *booleanLiteral) Internal() interface{} {
	return booleanLiteral.value
}

func (listValue *listValue) Print() string {
	return listValue.String()
}

func (listValue *listValue) String() string {
	return fmt.Sprintf("%v", listValue.values)
}

func (listValue *listValue) Type() Type {
	return TypeList
}

func (listValue *listValue) Internal() interface{} {
	return listValue.values
}

func (listValue *listValue) index(context RunContext, argument Argument) (int, ErrorValue) {
	indexValue := EvalArgument(context, argument)

	if indexValue.Type() != TypeInteger {
		return 0, NewErrorValue("list accessor must be an integer")
	}

	i := (int)(indexValue.Internal().(int64))

	// negative index will be used to get elemnts from the end of the list
	//
	if i < 0 {
		i = len(listValue.values) + i
	}

	if i < 0 || i >= len(listValue.values) {
		return 0, NewErrorValue("list accessor out of bounds")
	}

	return i, nil
}

func (listValue *listValue) Run(context RunContext, arguments []Argument) Value {
	arglen := len(arguments)

	if arglen == 1 {
		i, err := listValue.index(context, arguments[0])

		if err != nil {
			return err
		}

		return listValue.values[i]
	}

	if arglen == 2 {
		i1, err := listValue.index(context, arguments[0])
		if err != nil {
			return err
		}
		i2, err := listValue.index(context, arguments[1])
		if err != nil {
			return err
		}

		if i1 > i2 {
			// return a reversed version of the sub list

			list := listValue.values[i2 : i1+1]
			length := len(list)
			reversed := make([]Value, length)
			copy(reversed, list)
			for i, j := 0, length-1; i < j; i, j = i+1, j-1 {
				reversed[i], reversed[j] = reversed[j], reversed[i]
			}
			return NewListValue(reversed)
		}

		return NewListValue(listValue.values[i1 : i2+1])

	}

	return NewErrorValue("too many arguments for list access")
}

func (returnValue *returnValue) Print() string {
	return returnValue.String()
}

func (returnValue *returnValue) String() string {
	return fmt.Sprintf("<%v>", returnValue.values)
}

func (returnValue *returnValue) Type() Type {
	return TypeReturn
}

func (returnValue *returnValue) Internal() interface{} {
	return returnValue.values
}

func (dictValue *dictValue) Print() string {
	return dictValue.String()
}

func (dictValue *dictValue) String() string {
	return fmt.Sprintf("%v", dictValue.values)
}

func (dictValue *dictValue) Type() Type {
	return TypeDictionary
}

func (dictValue *dictValue) Internal() interface{} {
	return dictValue.values
}

func (dictValue *dictValue) Resolve(key string) Value {
	value, found := dictValue.values[key]

	if found {
		return value
	}

	if dictValue.parent != nil {
		return dictValue.parent.Resolve(key)
	}

	return Nothing
}

func (dictValue *dictValue) Run(context RunContext, arguments []Argument) Value {

	key := EvalArgument(context, arguments[0])
	return dictValue.Resolve(key.String())
}

func (errorValue *errorValue) Print() string {
	return errorValue.String()
}

func (errorValue *errorValue) String() string {
	if errorValue.meta != nil {
		meta, lineno := errorValue.At()
		return fmt.Sprintf("error at %s at line %d: %s", meta.Name(), lineno, errorValue.msg)
	}
	return fmt.Sprintf("error: %s", errorValue.msg)
}

func (errorValue *errorValue) Type() Type {
	return TypeError
}

func (errorValue *errorValue) Internal() interface{} {
	return errorValue.msg
}

func (errorValue *errorValue) Error() string {
	return errorValue.String()
}

func (errorValue *errorValue) SetAt(meta ScriptMetaData, lineno int) {
	errorValue.meta = meta
	errorValue.lineno = lineno
}

func (errorValue *errorValue) At() (meta ScriptMetaData, lineno int) {
	return errorValue.meta, errorValue.lineno
}

func (errorValue *errorValue) IsTraced() bool {
	meta, lineno := errorValue.At()
	return meta != nil && lineno > 0
}

func (goFunction *goFunction) Print() string {
	return goFunction.String()
}

func (goFunction *goFunction) String() string {
	return fmt.Sprintf("go:%s", goFunction.name)
}

func (goFunction *goFunction) Type() Type {
	return TypeGoFunction
}

func (goFunction *goFunction) Internal() interface{} {
	return goFunction.value
}

func (goFunction *goFunction) Name() string {
	return goFunction.name
}

func (goFunction *goFunction) Run(context RunContext, arguments []Argument) Value {
	return goFunction.value(context, arguments)
}

// NewIdentifier creates a new identifier value
//
func NewIdentifier(value string) Value {
	return &identifier{value: value}
}

// NewStringLiteral creates a new string literal value
//
func NewStringLiteral(value string) Value {
	return &stringLiteral{value: value}
}

// NewIntegerLiteral creates a new integer value
//
func NewIntegerLiteral(value int64) Value {
	return &integerLiteral{value: value}
}

// NewBooleanLiteral creates a new integer value
//
func NewBooleanLiteral(value bool) Value {
	return &booleanLiteral{value: value}
}

// NewListValue creates a new list of values
//
func NewListValue(values []Value) Value {
	return &listValue{values: values}
}

// NewDictionaryValue creates a new map of values
//
func NewDictionaryValue(parent *dictValue, values map[string]Value) Value {
	return &dictValue{parent: parent, values: values}
}

// NewErrorValue creates a new Error
//
func NewErrorValue(msg string) ErrorValue {
	return &errorValue{msg: msg}
}

// NewGoFunction creates a new go function
//
func NewGoFunction(name string, value GoFunction) NamedValue {
	return &goFunction{name: name, value: value}
}

// NewReturnValue creates a new list of values
//
func NewReturnValue(values []Value) Value {
	return &returnValue{values: values}
}

//
// ---[ASTNode]---------------------------------------------------------------
//
type astNode struct {
	meta  ScriptMetaData
	begin uint32
	end   uint32
}

//
// ---[ARGUMENT]---------------------------------------------------------------
//

type argument struct {
	astNode
	value Value
}

// Argument represent a function call parameter
//
type Argument interface {
	String() string
	Type() Type
	Value() Value
}

func (argument *argument) String() string {
	return argument.value.String()
}

func (argument *argument) Type() Type {
	return argument.value.Type()
}

func (argument *argument) Value() Value {
	return argument.value
}

// NewArgument constructs a new function argument
//
func NewArgument(meta ScriptMetaData, begin uint32, end uint32, value Value) Argument {
	return &argument{astNode: astNode{meta: meta, begin: begin, end: end}, value: value}
}

//
// ---[CALL]-------------------------------------------------------------------
//

type call struct {
	astNode
	functionName []string
	arguments    []Argument
	pipe         Call
}

// Call is a function call
//
type Call interface {
	RunnableValue

	Name() string
	Arguments() []Argument
	WillPipe() bool
}

func (call *call) Name() string {
	return strings.Join(call.functionName, ".")
}

func (call *call) Arguments() []Argument {
	return call.arguments
}

func (call *call) WillPipe() bool {
	return call.pipe != nil
}

func (call *call) addInfoWhenError(value Value) Value {
	if value.Type() == TypeError {
		if value.(ErrorValue).IsTraced() {
			// TODO: add trace??
			//
			return value
		}

		lineno, _ := call.meta.PositionOf(int(call.begin))
		value.(ErrorValue).SetAt(call.meta, lineno)
	}
	return value
}

func (call *call) pipeResult(context RunContext, value Value) Value {
	if !call.WillPipe() {
		return value
	}

	return call.pipe.Run(context, []Argument{&argument{value: value}})

}

func (call *call) Run(context RunContext, additionalArguments []Argument) Value {

	value, found := context.Get(call.functionName[0])

	var useArguments = call.arguments

	if additionalArguments != nil && len(additionalArguments) > 0 {
		useArguments = append([]Argument{}, useArguments...)
		useArguments = append(useArguments, additionalArguments...)
	}

	if found {

		if value == nil {
			return call.pipeResult(context, call.addInfoWhenError(NewErrorValue(fmt.Sprintf("call to %s results in invalid nil value", call.Name()))))
		}

		if len(call.functionName) > 1 {

			// call to a.b style of function name. a should be resolvable to a dictionary and b
			// can be something in that dictionary

			if value.Type() != TypeDictionary {
				return call.pipeResult(context, call.addInfoWhenError(NewErrorValue(fmt.Sprintf("%s does not resolve to dictionary. found %v", call.Name(), value))))
			}

			inDictValue := value.(*dictValue).Resolve(call.functionName[1])

			if inDictValue == nil {
				return call.pipeResult(context, call.addInfoWhenError(NewErrorValue(fmt.Sprintf("could not find %s.%s", call.functionName[0], call.functionName[1]))))
			}

			if inDictValue.Type() == TypeGoFunction {
				context.SetThis(value)
				result := call.addInfoWhenError(inDictValue.(Runnable).Run(context, useArguments))
				context.SetThis(Nothing)
				return call.pipeResult(context, result)
			}

			return call.pipeResult(context, call.addInfoWhenError(inDictValue))
		}

		if value.Type() == TypeGoFunction {
			return call.pipeResult(context, call.addInfoWhenError(value.(Runnable).Run(context, useArguments)))
		}

		// list and map values can be used as functions to access list content
		//
		if (value.Type() == TypeList || value.Type() == TypeDictionary) && (len(call.arguments) > 0) {
			return call.pipeResult(context, call.addInfoWhenError(value.(Runnable).Run(context, useArguments)))
		}

		return call.pipeResult(context, call.addInfoWhenError(value))
	}

	return call.pipeResult(context, call.addInfoWhenError(NewErrorValue(fmt.Sprintf("call to undefined \"%s\"", call.functionName))))
}

func (call *call) Print() string {
	return call.String()
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

// NewCall contstructs a new function call
//
func NewCall(meta ScriptMetaData, begin uint32, end uint32, name []string, arguments []Argument, pipeTo Call) Call {
	return &call{astNode: astNode{meta: meta, begin: begin, end: end}, functionName: name, arguments: arguments, pipe: pipeTo}
}

//
// ---[BLOCK]------------------------------------------------------------------
//

type block struct {
	astNode
	calls []Call
}

// Block is a list of function calls
//
type Block interface {
	// Block can be used as a value
	Value
	// Block can be executed
	Runnable
	Calls() []Call
}

func (block *block) Calls() []Call {
	return block.calls
}

func (block *block) Run(context RunContext, arguments []Argument) Value {
	var result Value = Nothing

	for _, call := range block.calls {
		result = call.Run(context, []Argument{})
		if context.isStopped() {
			break
		}
	}

	return result
}

func (block *block) Print() string {
	return block.String()
}

func (block *block) String() string {
	return fmt.Sprintf("{...}")
}

func (block *block) Type() Type {
	return TypeBlock
}

func (block *block) Internal() interface{} {
	return errors.New("Internal() not implemented on block")
}

// NewBlock contsruct a new block of function calls
//
func NewBlock(meta ScriptMetaData, begin uint32, end uint32, calls []Call) Block {
	return &block{astNode: astNode{meta: meta, begin: begin, end: end}, calls: calls}
}

// EvalArgument evaluates given argument
//
func EvalArgument(context RunContext, argument Argument) Value {

	if argument.Type() == TypeCall {
		return argument.Value().(Runnable).Run(context, NoArguments)
	}

	return argument.Value()

}

// EvalArgumentOrSolveIdentifier evaluates given argument
//
func EvalArgumentOrSolveIdentifier(context RunContext, argument Argument) Value {

	if argument.Type() == TypeIdentifier {
		value, found := context.Get(argument.String())
		if found {
			return value
		}
		return NewErrorValue(fmt.Sprintf("could not find %v", argument.String()))
	}

	return EvalArgument(context, argument)

}

// EvalArgumentWithBlock evaluates given argument and if argument is a block
// it will evaluate block content
//
func EvalArgumentWithBlock(context RunContext, argument Argument) Value {

	if argument.Type() == TypeCall || argument.Type() == TypeBlock {
		return argument.Value().(Runnable).Run(context, NoArguments)
	}

	return argument.Value()

}

// EvalArgument2String evaluates given argument and returns it String presentation
//
func EvalArgument2String(context RunContext, argument Argument) string {

	return EvalArgument(context, argument).String()

}
