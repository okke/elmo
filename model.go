package elmo

import "fmt"

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
)

type identifier struct {
	value string
}

type stringLiteral struct {
	value string
}

type integerLiteral struct {
	value int64
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
	return TypeIdentifier
}

func (stringLiteral *stringLiteral) String() string {
	return fmt.Sprintf("\"%s\"", stringLiteral.value)
}

func (stringLiteral *stringLiteral) Type() Type {
	return TypeString
}

func (integerLiteral *integerLiteral) String() string {
	return fmt.Sprintf("%d", integerLiteral.value)
}

func (integerLiteral *integerLiteral) Type() Type {
	return TypeInteger
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

type argument struct {
	value Value
}

// Argument represent a function call parameter
//
type Argument interface {
	String() string
	Type() Type
}

type call struct {
	functionName string
	arguments    []Argument
}

func (argument *argument) String() string {
	return argument.value.String()
}

func (argument *argument) Type() Type {
	return argument.value.Type()
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
