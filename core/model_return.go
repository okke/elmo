package elmo

import "fmt"

type returnValue struct {
	baseValue
	values []Value
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

// NewReturnValue creates a new list of values
//
func NewReturnValue(values []Value) Value {
	return &returnValue{baseValue: baseValue{info: typeInfoReturn}, values: values}
}
