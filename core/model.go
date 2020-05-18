package elmo

import (
	"github.com/google/uuid"
)

// Runnable is a type that can be interpreted
//
type Runnable interface {
	Run(RunContext, []Argument) Value
}

// Value represents data within elmo
//
type Value interface {
	String() string
	Type() Type
	Internal() interface{}
	Info() TypeInfo
	IsType(TypeInfo) bool
	UUID() uuid.UUID
}

// IdentifierValue represents a value that can be lookedup
//
type IdentifierValue interface {
	LookUp(RunContext) (DictionaryValue, Value, bool)
}

// StringValue represents a value of a string with Dynamic blocks of content
//
type StringValue interface {
	Value
	ResolveBlocks(RunContext) Value
	CopyWithinContext(context RunContext) StringValue
}

// IncrementableValue represents a value that can be incremented
//
type IncrementableValue interface {
	Increment(Value) Value
}

// DictionaryValue represents a value that can be used as dictionary
//
type DictionaryValue interface {
	Value
	Keys() []string
	Resolve(string) (Value, bool)
	Merge([]DictionaryValue) Value
	Replace(DictionaryValue)
	Set(symbol Value, value Value) (Value, ErrorValue)
	Remove(symbol Value) (Value, ErrorValue)
}

// Listable type can convert a value to a list
//
type Listable interface {
	List() []Value
}

// ListValue represents a value that can be used as a list of values
//
type ListValue interface {
	Value
	Listable
	Append(Value)
}

// MathValue represents a value that knows how to apply basic arithmetics
//
type MathValue interface {
	Plus(Value) Value
	Minus(Value) Value
	Multiply(Value) Value
	Divide(Value) Value
	Modulo(Value) Value
}

// ComparableValue represents a value that can be compared
//
type ComparableValue interface {
	Compare(RunContext, Value) (int, ErrorValue)
}

// HelpValue represents a value with help
//
type HelpValue interface {
	Help() Value
}

// MutableValue represents a value that can be mutated
//
type MutableValue interface {
	Mutate(value interface{}) (Value, ErrorValue)
}

// FreezableValue represents a value that can be frozen
// (protected from modification)
//
type FreezableValue interface {
	Freeze() Value
	Frozen() bool
}

// CloseableValue represents a resource that can be closed
//
type CloseableValue interface {
	Close()
}

// SerializableValue represents a value that can be serialized to
// a binary representation
//
type SerializableValue interface {
	ToBinary() BinaryValue
}

// BinaryValue represents a value that can be deserialized to
// a regular value
//
type BinaryValue interface {
	ToRegular() Value
	AsBytes() []byte
}

// ErrorValue represents an Error
//
type ErrorValue interface {
	Value
	Error() string
	SetAt(meta ScriptMetaData, lineno int)
	At() (ScriptMetaData, int)
	AtAbs() (ScriptMetaData, int)
	IsTraced() bool
	Panic() ErrorValue
	IsFatal() bool
	Ignore() ErrorValue
	CanBeIgnored() bool
}

// RunnableValue represents a value that can evaluated to another value
//
type RunnableValue interface {
	Value
	Runnable
}

type Inspectable interface {
	Meta() ScriptMetaData
	BeginsAt() uint32
	EndsAt() uint32
	Enrich(DictionaryValue)
}

// UserDefinedFunction represents a function written in elmo
//
type UserDefinedFunction interface {
	Block() Block
}

// NamedValue represent data with a name
//
type NamedValue interface {
	Value
	Name() string
}

// ValueWithLength represents a value of which length can be determined
//
type ValueWithLength interface {
	Length() Value
}
