package elmo

import "fmt"

type listValue struct {
	baseValue
	frozen bool
	values []Value
}

func (listValue *listValue) String() string {
	return fmt.Sprintf("%v", listValue.values)
}

func (listValue *listValue) Type() Type {
	return TypeList
}

func (listValue *listValue) Internal() interface{} {
	return listValue.values
}

func (listValue *listValue) index(context RunContext, argument Argument) (int, ErrorValue) {
	indexValue := EvalArgument(context, argument)

	if indexValue.Type() != TypeInteger {
		return 0, NewErrorValue("list accessor must be an integer")
	}

	i := (int)(indexValue.Internal().(int64))

	// negative index will be used to get elemnts from the end of the list
	//
	if i < 0 {
		i = len(listValue.values) + i
	}

	if i < 0 || i >= len(listValue.values) {
		return 0, NewErrorValue("list accessor out of bounds")
	}

	return i, nil
}

func (listValue *listValue) Run(context RunContext, arguments []Argument) Value {
	arglen := len(arguments)

	if arglen == 1 {
		i, err := listValue.index(context, arguments[0])

		if err != nil {
			return err
		}

		return listValue.values[i]
	}

	if arglen == 2 {
		i1, err := listValue.index(context, arguments[0])
		if err != nil {
			return err
		}
		i2, err := listValue.index(context, arguments[1])
		if err != nil {
			return err
		}

		if i1 > i2 {
			// return a reversed version of the sub list

			list := listValue.values[i2 : i1+1]
			length := len(list)
			reversed := make([]Value, length)
			copy(reversed, list)
			for i, j := 0, length-1; i < j; i, j = i+1, j-1 {
				reversed[i], reversed[j] = reversed[j], reversed[i]
			}
			return NewListValue(reversed)
		}

		return NewListValue(listValue.values[i1 : i2+1])

	}

	return NewErrorValue("too many arguments for list access")
}

func (listValue *listValue) Compare(context RunContext, value Value) (int, ErrorValue) {
	if value.Type() != TypeList {
		return 0, NewErrorValue("can not compare list with non list")
	}

	v1 := listValue.values
	v2 := value.Internal().([]Value)

	for i := range v1 {

		if i >= len(v2) {
			return 1, nil
		}

		c1, comparable := v1[i].(ComparableValue)
		if !comparable {
			return 0, NewErrorValue("list contains uncomparable type")
		}
		result, err := c1.Compare(context, v2[i])
		if err != nil {
			return 0, err
		}
		if result != 0 {
			return result, nil
		}

		if i == len(v1)-1 && len(v2) > len(v1) {
			return -1, nil
		}
	}

	return 0, nil
}

func (listValue *listValue) Append(value Value) {
	listValue.values = append(listValue.values, value)
}

func (listValue *listValue) Mutate(value interface{}) (Value, ErrorValue) {
	if listValue.Frozen() {
		return listValue, NewErrorValue("can not mutate frozen value")
	}

	listValue.values = value.([]Value)
	return listValue, nil
}

func (listValue *listValue) Freeze() Value {
	listValue.frozen = true
	for _, value := range listValue.values {
		if freezable, ok := value.(FreezableValue); ok && !freezable.Frozen() {
			freezable.Freeze()
		}
	}
	return listValue
}

func (listValue *listValue) Frozen() bool {
	return listValue.frozen
}

func (listValue *listValue) ToBinary() BinaryValue {
	return NewBinaryValueFromInternal(typeInfoList.ID(), "", Serialize(listValue))
}

// NewListValue creates a new list of values
//
func NewListValue(values []Value) Value {
	return &listValue{baseValue: baseValue{info: typeInfoList}, values: values}
}
