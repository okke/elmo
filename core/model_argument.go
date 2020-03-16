package elmo

import "github.com/google/uuid"

type argument struct {
	astNode
	value Value
}

// Argument represent a function call parameter
//
type Argument interface {
	Value

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

func (argument *argument) Info() TypeInfo {
	return argument.value.Info()
}

func (argument *argument) Internal() interface{} {
	return argument.value.Internal()
}

func (argument *argument) IsType(typeInfo TypeInfo) bool {
	return argument.value.IsType(typeInfo)
}

func (argument *argument) UUID() uuid.UUID {
	return argument.value.UUID()
}

func (argument *argument) Enrich(dict DictionaryValue) {
	dict.Set(NewStringLiteral("value"), argument.value)
	dict.Set(NewStringLiteral("type"), argument.value.Info().Name())
}

// NewArgument constructs a new function argument
//
func NewArgument(meta ScriptMetaData, node *node32, value Value) Argument {
	return &argument{astNode: astNode{meta: meta, node: node}, value: value}
}

// NewArgumentWithDots constructs a new function argument consisting of multiple ast nodes
//
func NewArgumentWithDots(meta ScriptMetaData, nodeBegin *node32, nodeEnd *node32, value Value) Argument {
	return &argument{astNode: astNode{meta: meta, node: nodeBegin, endNode: nodeEnd}, value: value}
}

// NewDynamicArgument constructs a new function argument without script info
//
func NewDynamicArgument(value Value) Argument {
	return &argument{value: value}
}
