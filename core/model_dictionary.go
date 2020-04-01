package elmo

import "fmt"

type dictValue struct {
	baseValue
	frozen bool
	parent *dictValue
	values map[string]Value
}

func (dictValue *dictValue) String() string {
	return fmt.Sprintf("%v", dictValue.values)
}

func (dictValue *dictValue) Type() Type {
	return TypeDictionary
}

func (dictValue *dictValue) Internal() interface{} {
	return dictValue.values
}

func (dictValue *dictValue) Keys() []string {
	keyNames := make([]string, len(dictValue.values))

	i := 0
	for k := range dictValue.values {
		keyNames[i] = k
		i++
	}

	return keyNames
}

func (dictValue *dictValue) Resolve(key string) (Value, bool) {
	value, found := dictValue.values[key]

	if found {
		return value, true
	}

	if dictValue.parent != nil {
		return dictValue.parent.Resolve(key)
	}

	return Nothing, false
}

func (dictValue *dictValue) Merge(withAll []DictionaryValue) Value {
	newMap := make(map[string]Value)

	for k, v := range dictValue.values {
		newMap[k] = v
	}

	for _, with := range withAll {
		for _, k := range with.Keys() {

			value, found := with.Resolve(k)
			if !found {
				return NewErrorValue(fmt.Sprintf("could not merge value %s", k))
			}
			newMap[k] = value
		}
	}

	return NewDictionaryValue(dictValue.parent, newMap)
}

func (dictValue *dictValue) Set(symbol Value, value Value) (Value, ErrorValue) {
	if dictValue.Frozen() {
		return dictValue, NewErrorValue("can not set value in frozen dictionary")
	}
	dictValue.values[symbol.String()] = value

	return dictValue, nil
}

func (dictValue *dictValue) Remove(symbol Value) (Value, ErrorValue) {
	if dictValue.Frozen() {
		return dictValue, NewErrorValue("can not remove value from frozen dictionary")
	}
	delete(dictValue.values, symbol.String())

	return dictValue, nil
}

func (dictValue *dictValue) runInternalCompareFunc(context RunContext, compareFunction Runnable, value Value) (int, ErrorValue) {
	if runnableCompare, compareIsRunnable := compareFunction.(Runnable); compareIsRunnable {

		subContext := context.CreateSubContext()
		subContext.SetThis(dictValue)

		value := runnableCompare.Run(subContext, []Argument{NewDynamicArgument(value)})
		if value.Type() == TypeError {
			return -1, value.(ErrorValue)
		}
		if value.Type() == TypeInteger {
			return int(value.Internal().(int64)), nil
		}
		return -1, NewErrorValue("found _compare did not return an integer")
	}
	return -1, NewErrorValue("found _compare is not runnable")
}

func (dictValue *dictValue) Compare(context RunContext, value Value) (int, ErrorValue) {

	if value.Type() != TypeDictionary {
		return 0, NewErrorValue("can not compare dictionary with non dictionary")
	}

	// when there is a compare function found in the dictionary, simply use that one
	//
	if compareFunction, hasCompareFunction := dictValue.Resolve("_compare"); hasCompareFunction {
		result, err := dictValue.runInternalCompareFunc(context, compareFunction.(Runnable), value)
		return result, err
	}

	keys1 := dictValue.Keys()
	keys2 := value.(DictionaryValue).Keys()

	if len(keys1) == len(keys2) {
		for _, key := range keys1 {
			kval2, found2 := value.(DictionaryValue).Resolve(key)
			if !found2 {
				return -1, NewErrorValue("can not compare asymetric dictionaries")
			}
			kval1, _ := dictValue.Resolve(key)

			comparable, isComparable := kval1.(ComparableValue)
			if !isComparable {
				return -1, NewErrorValue("dictionary contains uncomparable values")
			}

			compared, err := comparable.Compare(context, kval2)
			if err != nil {
				return -1, err
			}
			if compared != 0 {
				return compared, nil
			}
		}

		// everything equal
		return 0, nil
	}

	if len(keys1) < len(keys2) {
		return -1, nil
	}
	return 1, nil
}

func (dictValue *dictValue) Freeze() Value {
	dictValue.frozen = true
	for _, value := range dictValue.values {
		if freezable, ok := value.(FreezableValue); ok && !freezable.Frozen() {
			freezable.Freeze()
		}
	}
	return dictValue
}

func (dictValue *dictValue) Frozen() bool {
	return dictValue.frozen
}

func (dictValue *dictValue) Run(context RunContext, arguments []Argument) Value {

	key := EvalArgument(context, arguments[0])
	result, _ := dictValue.Resolve(key.String())
	return result
}

func (dictValue *dictValue) ToBinary() BinaryValue {
	return NewBinaryValueFromInternal(typeInfoDictionary.ID(), "", Serialize(dictValue))
}

// NewDictionaryValue creates a new map of values
// TODO: 31okt2016 introduce interface for map parents
//
func NewDictionaryValue(parent interface{}, values map[string]Value) DictionaryValue {
	if parent == nil {
		return &dictValue{baseValue: baseValue{info: typeInfoDictionary}, parent: nil, values: values}
	}
	return &dictValue{baseValue: baseValue{info: typeInfoDictionary}, parent: parent.(*dictValue), values: values}
}

// NewDictionaryWithBlock constructs a new dictionary by evaluating given block
//
func NewDictionaryWithBlock(context RunContext, block Block) DictionaryValue {

	// use NewRunContext so block will be evaluated within same scope
	//
	subContext := NewRunContext(context)

	block.Run(subContext, NoArguments)

	return NewDictionaryValue(nil, subContext.Mapping())
}

// NewDictionaryFromList constructs a dictionary from a list of values
// where the list is of the form [key value key value]
//
// note, NewDictionaryFromList can return an ErrorValue as well instead of a DictionaryValue
// when the list has off values
//
func NewDictionaryFromList(parent interface{}, values []Value) Value {

	if (len(values) % 2) != 0 {
		return NewErrorValue(fmt.Sprintf("can not create a dictionary from an odd number of elements using %v", values))
	}

	mapping := make(map[string]Value)

	var key Value

	for i, val := range values {
		if i%2 == 0 {
			key = val
		} else {
			mapping[key.String()] = val
		}
	}

	return NewDictionaryValue(parent, mapping)

}
