package elmo

import (
	"fmt"
	"math"
	"strings"
)

type floatLiteral struct {
	baseValue
	value float64
}

func (floatLiteral *floatLiteral) String() string {
	return strings.TrimRight(fmt.Sprintf("%f", floatLiteral.value), "0")
}

func (floatLiteral *floatLiteral) Type() Type {
	return TypeFloat
}

func (floatLiteral *floatLiteral) Internal() interface{} {
	return floatLiteral.value
}

func (floatLiteral *floatLiteral) Increment(value Value) Value {
	if value.Type() == TypeFloat {
		return NewFloatLiteral(floatLiteral.value + value.Internal().(float64))
	}
	if value.Type() == TypeInteger {
		return NewFloatLiteral(floatLiteral.value + float64(value.Internal().(int64)))
	}
	return NewErrorValue("can not add non number to float")
}

func (floatLiteral *floatLiteral) Plus(value Value) Value {
	return floatLiteral.Increment(value)
}

func (floatLiteral *floatLiteral) Minus(value Value) Value {
	if value.Type() == TypeFloat {
		return NewFloatLiteral(floatLiteral.value - value.Internal().(float64))
	}
	if value.Type() == TypeInteger {
		return NewFloatLiteral(floatLiteral.value - float64(value.Internal().(int64)))
	}
	return NewErrorValue("can not subtract non number from float")
}

func (floatLiteral *floatLiteral) Multiply(value Value) Value {
	if value.Type() == TypeFloat {
		return NewFloatLiteral(floatLiteral.value * value.Internal().(float64))
	}
	if value.Type() == TypeInteger {
		return NewFloatLiteral(floatLiteral.value * float64(value.Internal().(int64)))
	}
	return NewErrorValue("can not multiply float by non number")
}

func (floatLiteral *floatLiteral) Divide(value Value) Value {
	if value.Type() == TypeFloat {
		if value.Internal().(float64) == 0.0 {
			return NewErrorValue("can not divide float by 0.0")
		}
		return NewFloatLiteral(floatLiteral.value / value.Internal().(float64))
	}
	if value.Type() == TypeInteger {
		if value.Internal().(int64) == 0 {
			return NewErrorValue("can not divide float by 0")
		}
		return NewFloatLiteral(floatLiteral.value / float64(value.Internal().(int64)))
	}
	return NewErrorValue("can not multiply float by non number")
}

func (floatLiteral *floatLiteral) Modulo(value Value) Value {
	if value.Type() == TypeInteger {
		if value.Internal().(int64) == 0 {
			return NewErrorValue("can not divide integer by 0")
		}

		div := floatLiteral.value / float64(value.Internal().(int64))
		total := math.Floor(div) * float64(value.Internal().(int64))

		return NewFloatLiteral(floatLiteral.value - total)
	}
	if value.Type() == TypeFloat {
		if value.Internal().(float64) == 0.0 {
			return NewErrorValue("can not divide integer by 0")
		}

		div := floatLiteral.value / float64(value.Internal().(float64))
		total := math.Floor(div) * float64(value.Internal().(float64))

		return NewFloatLiteral(floatLiteral.value - total)
	}
	return NewErrorValue("can not divide float by non number to calculate a modulo")
}

func (floatLiteral *floatLiteral) Compare(context RunContext, value Value) (int, ErrorValue) {
	if value.Type() == TypeFloat {
		v1 := floatLiteral.value
		v2 := value.Internal().(float64)

		if v1 > v2 {
			return 1, nil
		}
		if v2 > v1 {
			return -1, nil
		}
		return 0, nil
	}
	return 0, NewErrorValue("can not compare float with non float")
}

func (floatLiteral *floatLiteral) ToBinary() BinaryValue {
	return NewBinaryValueFromInternal(typeInfoFloat.ID(), "", floatLiteral.value)
}

// NewFloatLiteral creates a new integer value
//
func NewFloatLiteral(value float64) Value {
	return &floatLiteral{baseValue: baseValue{info: typeInfoFloat}, value: value}
}
