package elmo

import "fmt"

type internalValue struct {
	baseValue
	value interface{}
}

func (internalValue *internalValue) String() string {
	return fmt.Sprintf("%v", internalValue.value)
}

func (internalValue *internalValue) Type() Type {
	return TypeInternal
}

func (internalValue *internalValue) Internal() interface{} {
	return internalValue.value
}

// NewInternalValue wraps a go value into an elmo value
//
func NewInternalValue(info TypeInfo, value interface{}) Value {
	return &internalValue{baseValue: baseValue{info: info}, value: value}
}
