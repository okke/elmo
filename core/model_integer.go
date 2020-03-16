package elmo

import "fmt"

type integerLiteral struct {
	baseValue
	value int64
}

// Zero is 0
//
var Zero = NewIntegerLiteral(0)

// One is 1
//
var One = NewIntegerLiteral(1)

func (integerLiteral *integerLiteral) String() string {
	return fmt.Sprintf("%d", integerLiteral.value)
}

func (integerLiteral *integerLiteral) Type() Type {
	return TypeInteger
}

func (integerLiteral *integerLiteral) Internal() interface{} {
	return integerLiteral.value
}

func (integerLiteral *integerLiteral) Increment(value Value) Value {
	if value.Type() == TypeInteger {
		return NewIntegerLiteral(integerLiteral.value + value.Internal().(int64))
	}
	return NewErrorValue("can not add non integer to integer")
}

func (integerLiteral *integerLiteral) Plus(value Value) Value {
	if value.Type() == TypeInteger {
		return NewIntegerLiteral(integerLiteral.value + value.Internal().(int64))
	}
	if value.Type() == TypeFloat {
		return NewFloatLiteral(float64(integerLiteral.value) + value.Internal().(float64))
	}
	return NewErrorValue("can not add non number to integer")
}

func (integerLiteral *integerLiteral) Minus(value Value) Value {
	if value.Type() == TypeInteger {
		return NewIntegerLiteral(integerLiteral.value - value.Internal().(int64))
	}
	if value.Type() == TypeFloat {
		return NewFloatLiteral(float64(integerLiteral.value) - value.Internal().(float64))
	}
	return NewErrorValue("can not subtract non number from integer")
}

func (integerLiteral *integerLiteral) Multiply(value Value) Value {
	if value.Type() == TypeInteger {
		return NewIntegerLiteral(integerLiteral.value * value.Internal().(int64))
	}
	if value.Type() == TypeFloat {
		return NewFloatLiteral(float64(integerLiteral.value) * value.Internal().(float64))
	}
	return NewErrorValue("can not multiply non number with integer")
}

func (integerLiteral *integerLiteral) Divide(value Value) Value {
	if value.Type() == TypeInteger {
		if value.Internal().(int64) == 0 {
			return NewErrorValue("can not divide integer by 0")
		}
		return NewIntegerLiteral(integerLiteral.value / value.Internal().(int64))
	}
	if value.Type() == TypeFloat {
		if value.Internal().(float64) == 0.0 {
			return NewErrorValue("can not divide integer by 0.0")
		}
		return NewFloatLiteral(float64(integerLiteral.value) / value.Internal().(float64))
	}
	return NewErrorValue("can not divide integer by non number")
}

func (integerLiteral *integerLiteral) Modulo(value Value) Value {
	if value.Type() == TypeInteger {
		if value.Internal().(int64) == 0 {
			return NewErrorValue("can not divide integer by 0")
		}
		return NewIntegerLiteral(integerLiteral.value % value.Internal().(int64))
	}
	return NewErrorValue("can not divide integer by non integer to calculate a modulo")
}

func (integerLiteral *integerLiteral) Compare(context RunContext, value Value) (int, ErrorValue) {
	if value.Type() == TypeInteger {
		v1 := integerLiteral.value
		v2 := value.Internal().(int64)

		if v1 > v2 {
			return 1, nil
		}
		if v2 > v1 {
			return -1, nil
		}
		return 0, nil
	}
	return 0, NewErrorValue("can not compare integer with non integer")
}

func (integerLiteral *integerLiteral) ToBinary() BinaryValue {
	return NewBinaryValueFromInternal(typeInfoInteger.ID(), "", integerLiteral.value)
}

// NewIntegerLiteral creates a new integer value
//
func NewIntegerLiteral(value int64) Value {
	return &integerLiteral{baseValue: baseValue{info: typeInfoInteger}, value: value}
}
