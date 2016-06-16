package elmo

import "fmt"

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
	// TypeGoFunction represents a type for an internal go function
	TypeGoFunction
	// TypeNil represents the type of a nil value
	TypeNil
)

type nothing struct {
}

// Nothing represents nil
//
var Nothing = &nothing{}

type identifier struct {
	value string
}

type stringLiteral struct {
	value string
}

type integerLiteral struct {
	value int64
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

func (identifier *identifier) Print() string {
	return identifier.value
}

func (identifier *identifier) String() string {
	return identifier.value
}

func (identifier *identifier) Type() Type {
	return TypeIdentifier
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

func (integerLiteral *integerLiteral) Print() string {
	return fmt.Sprintf("%d", integerLiteral.value)
}

func (integerLiteral *integerLiteral) String() string {
	return fmt.Sprintf("%d", integerLiteral.value)
}

func (integerLiteral *integerLiteral) Type() Type {
	return TypeInteger
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

// NewGoFunction creates a new go function
//
func NewGoFunction(name string, value GoFunction) NamedValue {
	return &goFunction{name: name, value: value}
}

//
// ---[ARGUMENT]---------------------------------------------------------------
//

type argument struct {
	value Value
}

// Argument represent a function call parameter
//
type Argument interface {
	String() string
	Type() Type
}

func (argument *argument) String() string {
	return argument.value.String()
}

func (argument *argument) Type() Type {
	return argument.value.Type()
}

//
// ---[CALL]-------------------------------------------------------------------
//

type call struct {
	functionName string
	arguments    []Argument
}

// Call is a function call
//
type Call interface {
	Name() string
	Arguments() []Argument
	Run(context RunContext) Value
}

func (call *call) Name() string {
	return call.functionName
}

func (call *call) Arguments() []Argument {
	return call.arguments
}

func (call *call) Run(context RunContext) Value {
	value, found := context.Get(call.functionName)

	if found {
		if value.Type() == TypeGoFunction {
			// TODO How to handle arguments
			//
			return value.(Runnable).Run(context, call.arguments)
		}
		return value
	}

	return Nothing
}

// NewCall contstructs a new function call
//
func NewCall(name string, arguments []Argument) Call {
	return &call{functionName: name, arguments: arguments}
}

//
// ---[BLOCK]------------------------------------------------------------------
//

type block struct {
	calls []Call
}

// Block is a list of function calls
//
type Block interface {
	Calls() []Call
	Run(context RunContext) Value
}

func (block *block) Calls() []Call {
	return block.calls
}

func (block *block) Run(context RunContext) Value {
	var result Value = Nothing

	for _, call := range block.calls {
		result = call.Run(context)
	}

	return result
}

// NewBlock contsruct a new block of function calls
//
func NewBlock(calls []Call) Block {
	return &block{calls: calls}
}