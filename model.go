package elmo

import "fmt"

// Type represents an internal value type
//
type Type uint8

const (
	typeIdentifier Type = iota
	typeString
)

type identifier struct {
	value string
}

type stringLiteral struct {
	value string
}

// Value represents data within elmo
//
type Value interface {
	String() string
	Type() Type
}

func (identifier *identifier) String() string {
	return identifier.value
}

func (identifier *identifier) Type() Type {
	return typeIdentifier
}

func (stringLiteral *stringLiteral) String() string {
	return fmt.Sprintf("\"%s\"", stringLiteral.value)
}

func (stringLiteral *stringLiteral) Type() Type {
	return typeString
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

type argument struct {
	value Value
}

// Argument represent a function call parameter
//
type Argument interface {
	String() string
}

type call struct {
	functionName string
	arguments    []Argument
}

func (argument *argument) String() string {
	return argument.value.String()
}

// Call is a function call
//
type Call interface {
	Name() string
	Arguments() []Argument
}

func (call *call) Name() string {
	return call.functionName
}

func (call *call) Arguments() []Argument {
	return call.arguments
}

// NewCall contstructs a new function call
//
func NewCall(name string, arguments []Argument) Call {
	return &call{functionName: name, arguments: arguments}
}

type block struct {
	calls []Call
}

// Block is a list of function calls
//
type Block interface {
	Calls() []Call
}

func (block *block) Calls() []Call {
	return block.calls
}

// NewBlock contsruct a new block of function calls
//
func NewBlock(calls []Call) Block {
	return &block{calls: calls}
}
