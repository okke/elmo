package elmo

import "fmt"

type booleanLiteral struct {
	baseValue
	value bool
}

// True represents boolean value true
//
var True = NewBooleanLiteral(true)

// False represents boolean value false
//
var False = NewBooleanLiteral(false)

func (booleanLiteral *booleanLiteral) String() string {
	return fmt.Sprintf("%v", booleanLiteral.value)
}

func (booleanLiteral *booleanLiteral) Type() Type {
	return TypeBoolean
}

func (booleanLiteral *booleanLiteral) Internal() interface{} {
	return booleanLiteral.value
}

func (booleanLiteral *booleanLiteral) ToBinary() BinaryValue {
	return NewBinaryValueFromInternal(typeInfoBoolean.ID(), "", booleanLiteral.value)
}

func (booleanLiteral *booleanLiteral) Compare(context RunContext, value Value) (int, ErrorValue) {
	if value.Type() != TypeBoolean {
		return -1, NewErrorValue("can not compare boolean with non boolean")
	}
	with := value.Internal().(bool)
	if booleanLiteral.value {
		if with {
			return 0, nil
		}
		return 1, nil
	}
	if !with {
		return 0, nil
	}
	return -1, nil
}

// NewBooleanLiteral creates a new integer value
//
func NewBooleanLiteral(value bool) Value {
	return &booleanLiteral{baseValue: baseValue{info: typeInfoBoolean}, value: value}
}
