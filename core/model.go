package elmo

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"math"
	"strings"
	"sync"

	"github.com/google/uuid"
)

// Runnable is a type that can be interpreted
//
type Runnable interface {
	Run(RunContext, []Argument) Value
}

//
// ---[LITERALS]---------------------------------------------------------------
//

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
	// TypeFloat represents a type for a floating point value
	TypeFloat
	// TypeBoolean represents a type for a boolean value
	TypeBoolean
	// TypeList represents a type for an array value
	TypeList
	// TypeDictionary represents a type for a map value
	TypeDictionary
	// TypeError represents a type for an error value
	TypeError
	// TypeInternal represents an internal type
	TypeInternal
	// TypeBlock represents a type for a code block
	TypeBlock
	// TypeCall represent the type for a function call
	TypeCall
	// TypeGoFunction represents a type for an internal go function
	TypeGoFunction
	// TypeReturn represents a function result containing multiple values
	TypeReturn
	// TypeNil represents the type of a nil value
	TypeNil
	// TypeBinary represents the value of a byte array
	TypeBinary
)

var typeInfoIdentifier = NewTypeInfo("identifier")
var typeInfoString = NewTypeInfo("string")
var typeInfoInteger = NewTypeInfo("int")
var typeInfoFloat = NewTypeInfo("float")
var typeInfoBoolean = NewTypeInfo("bool")
var typeInfoList = NewTypeInfo("list")
var typeInfoDictionary = NewTypeInfo("dict")
var typeInfoError = NewTypeInfo("error")
var typeInfoBlock = NewTypeInfo("block")
var typeInfoCall = NewTypeInfo("call")
var typeInfoGoFunction = NewTypeInfo("func")
var typeInfoReturn = NewTypeInfo("return")
var typeInfoNil = NewTypeInfo("nil")
var typeInfoBinary = NewTypeInfo("binary")

// TypeInfo represents kinf of subType for TypeInternal values
//
type TypeInfo interface {
	ID() int64
	Name() Value
}

type typeInfo struct {
	id   int64
	name string
}

func (typeInfo *typeInfo) Name() Value {
	return NewIdentifier(typeInfo.name)
}

func (typeInfo *typeInfo) ID() int64 {
	return typeInfo.id
}

var typeCounter int64

// NewTypeInfo constructs a new type object
//
func NewTypeInfo(name string) TypeInfo {
	typeCounter = typeCounter + 1
	return &typeInfo{id: typeCounter, name: name}
}

type baseValue struct {
	info  TypeInfo
	id    uuid.UUID
	mutex sync.Mutex
}

type nothing struct {
	baseValue
}

// Nothing represents nil
//
var Nothing = &nothing{}

// True represents boolean value true
//
var True = NewBooleanLiteral(true)

// False represents boolean value false
//
var False = NewBooleanLiteral(false)

// Zero is 0
//
var Zero = NewIntegerLiteral(0)

// One is 1
//
var One = NewIntegerLiteral(1)

type identifier struct {
	baseValue
	value []string
}

type stringLiteral struct {
	baseValue
	value string
}

type integerLiteral struct {
	baseValue
	value int64
}

type floatLiteral struct {
	baseValue
	value float64
}

type booleanLiteral struct {
	baseValue
	value bool
}

type listValue struct {
	baseValue
	frozen bool
	values []Value
}

type returnValue struct {
	baseValue
	values []Value
}

type dictValue struct {
	baseValue
	frozen bool
	parent *dictValue
	values map[string]Value
}

type internalValue struct {
	baseValue
	value interface{}
}

type errorValue struct {
	baseValue
	meta   ScriptMetaData
	lineno int
	msg    string
	fatal  bool
	ignore bool
}

// BinaryData is a struct used to serialize values to binary data
// Note, property names are as short as possible and public
// so the gob package can easily serialize it
//
type BinaryData struct {
	// Type Id for core types or -1
	//
	I int64

	// Name for non core types
	//
	N string

	// Actual data
	//
	D []byte
}

type binaryValue struct {
	baseValue
	data []byte
}

// GoFunction is a native go function that takes an array of input values
// and returns an output value
//
type GoFunction func(RunContext, []Argument) Value

type goFunction struct {
	baseValue
	name  string
	help  Value
	value GoFunction
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

// IncrementableValue represents a value that can be incremented
//
type IncrementableValue interface {
	Increment(Value) Value
}

// DictionaryValue represents a value that can be used as dictionary
//
type DictionaryValue interface {
	Keys() []string
	Resolve(string) (Value, bool)
	Merge([]DictionaryValue) Value
	Set(symbol Value, value Value) (Value, ErrorValue)
	Remove(symbol Value) (Value, ErrorValue)
}

// ListValue represents a value that can be used as a list of values
//
type ListValue interface {
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

// NamedValue represent data with a name
//
type NamedValue interface {
	Value
	Name() string
}

func (baseValue *baseValue) Info() TypeInfo {
	return baseValue.info
}

func (baseValue *baseValue) IsType(typeInfo TypeInfo) bool {
	if baseValue.info == nil {
		return false
	}

	return baseValue.info.ID() == typeInfo.ID()
}

func (baseValue *baseValue) Type() Type {
	fmt.Printf("check type of %v\n", baseValue)
	panic("baseValue does not support type")
}

func (baseValue *baseValue) Internal() interface{} {
	panic("baseValue does not support internal")
}

func (baseValue *baseValue) String() string {
	return "baseValue[?]"
}

func (baseValue *baseValue) UUID() uuid.UUID {
	baseValue.mutex.Lock()
	defer baseValue.mutex.Unlock()

	if baseValue.id[0] == 0 {
		baseValue.id = uuid.New()
	}

	return baseValue.id
}

func (nothing *nothing) String() string {
	return "nil"
}

func (nothing *nothing) Type() Type {
	return TypeNil
}

func (nothing *nothing) Internal() interface{} {
	return nil
}

func (binaryValue *binaryValue) String() string {
	return fmt.Sprintf("binary:%v", binaryValue.data)
}

func (binaryValue *binaryValue) Type() Type {
	return TypeBinary
}

func (binaryValue *binaryValue) Internal() interface{} {
	return binaryValue.data
}

func (binaryValue *binaryValue) AsBytes() []byte {
	return binaryValue.data
}

func (binaryValue *binaryValue) ToRegular() Value {

	buf := bytes.NewBuffer(binaryValue.data)
	decoder := gob.NewDecoder(buf)
	var bdata BinaryData
	if err := decoder.Decode(&bdata); err != nil {
		return NewErrorValue(err.Error())
	}

	dataBuffer := bytes.NewBuffer(bdata.D)
	decoder = gob.NewDecoder(dataBuffer)

	switch bdata.I {
	case typeInfoIdentifier.ID():
		actualData := []string{}
		if err := decoder.Decode(&actualData); err != nil {
			return NewErrorValue(err.Error())
		}
		return NewNameSpacedIdentifier(actualData)
	case typeInfoString.ID():
		actualData := ""
		if err := decoder.Decode(&actualData); err != nil {
			return NewErrorValue(err.Error())
		}
		return NewStringLiteral(actualData)
	case typeInfoInteger.ID():
		actualData := int64(0)
		if err := decoder.Decode(&actualData); err != nil {
			return NewErrorValue(err.Error())
		}
		return NewIntegerLiteral(actualData)
	case typeInfoFloat.ID():
		actualData := float64(0)
		if err := decoder.Decode(&actualData); err != nil {
			return NewErrorValue(err.Error())
		}
		return NewFloatLiteral(actualData)
	case typeInfoBoolean.ID():
		actualData := true
		if err := decoder.Decode(&actualData); err != nil {
			return NewErrorValue(err.Error())
		}
		return NewBooleanLiteral(actualData)
	case typeInfoList.ID(), typeInfoDictionary.ID():
		actualData := SerializationResult{}
		if err := decoder.Decode(&actualData); err != nil {
			return NewErrorValue(err.Error())
		}
		return actualData.ToValue()
	default:
		return NewErrorValue("binaryValue.AsRegular: operation unsupported yet")
	}

}

func (identifier *identifier) String() string {
	if len(identifier.value) == 1 {
		return identifier.value[0]
	}

	return strings.Join(identifier.value, ".")
}

func (identifier *identifier) Type() Type {
	return TypeIdentifier
}

func (identifier *identifier) Internal() interface{} {
	return identifier.value
}

func (identifier *identifier) LookUp(context RunContext) (DictionaryValue, Value, bool) {

	result, found := context.Get(identifier.value[0])
	if !found {
		return nil, NewErrorValue(fmt.Sprintf("could not resolve %v", identifier)), false
	}

	if len(identifier.value) == 1 {
		return nil, result, true
	}

	if result.Type() != TypeDictionary {
		return nil, NewErrorValue(fmt.Sprintf("%s is not a dictionary", identifier.value[0])), false
	}

	var dict = result.(DictionaryValue)
	var lookup Value

	for _, name := range identifier.value[1:] {
		lookup, found = dict.Resolve(name)

		if found {
			if lookup.Type() != TypeDictionary {
				return dict, lookup, true
			}

			dict = lookup.(DictionaryValue)
		} else {

			return dict, NewErrorValue(fmt.Sprintf("could not resolve %v", identifier)), false

		}
	}
	return dict, lookup, true
}

func (identifier *identifier) ToBinary() BinaryValue {
	return NewBinaryValueFromInternal(typeInfoIdentifier.ID(), "", identifier.value)
}

func (stringLiteral *stringLiteral) String() string {
	return fmt.Sprintf("%s", stringLiteral.value)
}

func (stringLiteral *stringLiteral) Type() Type {
	return TypeString
}

func (stringLiteral *stringLiteral) Internal() interface{} {
	return stringLiteral.value
}

func (stringLiteral *stringLiteral) index(context RunContext, argument Argument) (int, ErrorValue) {
	indexValue := EvalArgument(context, argument)

	if indexValue.Type() != TypeInteger {
		return 0, NewErrorValue("string accessor must be an integer")
	}

	i := (int)(indexValue.Internal().(int64))

	// negative index will be used to get elemnts from the end of the list
	//
	if i < 0 {
		i = len(stringLiteral.value) + i
	}

	if i < 0 || i >= len(stringLiteral.value) {
		return 0, NewErrorValue("string accessor out of bounds")
	}

	return i, nil
}

func (stringLiteral *stringLiteral) Run(context RunContext, arguments []Argument) Value {

	arglen := len(arguments)

	if arglen == 1 {
		i, err := stringLiteral.index(context, arguments[0])

		if err != nil {
			return err
		}

		return NewStringLiteral(stringLiteral.value[i : i+1])
	}

	if arglen == 2 {
		i1, err := stringLiteral.index(context, arguments[0])
		if err != nil {
			return err
		}
		i2, err := stringLiteral.index(context, arguments[1])
		if err != nil {
			return err
		}

		if i1 > i2 {
			// return a reversed version of the sub list

			sub := stringLiteral.value[i2 : i1+1]
			runes := []rune(sub)
			for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
				runes[i], runes[j] = runes[j], runes[i]
			}
			return NewStringLiteral(string(runes))
		}

		return NewStringLiteral(stringLiteral.value[i1 : i2+1])

	}

	return NewErrorValue("too many arguments for string access")
}

func (stringLiteral *stringLiteral) ToBinary() BinaryValue {
	return NewBinaryValueFromInternal(typeInfoString.ID(), "", stringLiteral.value)
}

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

func (returnValue *returnValue) String() string {
	return fmt.Sprintf("<%v>", returnValue.values)
}

func (returnValue *returnValue) Type() Type {
	return TypeReturn
}

func (returnValue *returnValue) Internal() interface{} {
	return returnValue.values
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

func (errorValue *errorValue) String() string {
	kind := "error"
	if errorValue.IsFatal() {
		kind = "fatal error"
	}
	if errorValue.meta != nil {
		meta, lineno := errorValue.At()
		return fmt.Sprintf("%s(at %s at line %d: %s)", kind, meta.Name(), lineno, errorValue.msg)
	}
	return fmt.Sprintf("%s(%s)", kind, errorValue.msg)
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

func (errorValue *errorValue) Type() Type {
	return TypeError
}

func (errorValue *errorValue) Internal() interface{} {
	return errorValue.msg
}

func (errorValue *errorValue) Error() string {
	return errorValue.String()
}

func (errorValue *errorValue) SetAt(meta ScriptMetaData, lineno int) {
	errorValue.meta = meta
	errorValue.lineno = lineno
}

func (errorValue *errorValue) At() (meta ScriptMetaData, lineno int) {
	return errorValue.meta, errorValue.lineno
}

func (errorValue *errorValue) IsTraced() bool {
	meta, lineno := errorValue.At()
	return meta != nil && lineno > 0
}

func (errorValue *errorValue) Panic() ErrorValue {
	errorValue.fatal = true
	return errorValue
}

func (errorValue *errorValue) IsFatal() bool {
	return errorValue.fatal
}

func (errorValue *errorValue) Ignore() ErrorValue {
	errorValue.ignore = true
	errorValue.fatal = false
	return errorValue
}

func (errorValue *errorValue) CanBeIgnored() bool {
	return errorValue.ignore
}

func (goFunction *goFunction) String() string {
	return fmt.Sprintf("func(%s)", goFunction.name)
}

func (goFunction *goFunction) Type() Type {
	return TypeGoFunction
}

func (goFunction *goFunction) Internal() interface{} {
	return goFunction.value
}

func (goFunction *goFunction) Name() string {
	return goFunction.name
}

func (goFunction *goFunction) Run(context RunContext, arguments []Argument) Value {
	return goFunction.value(context, arguments)
}

func (goFunction *goFunction) Help() Value {
	if goFunction.help == nil {
		return Nothing
	}
	return goFunction.help
}

// NewIdentifier creates a new identifier value
//
func NewIdentifier(value string) Value {
	return &identifier{baseValue: baseValue{info: typeInfoIdentifier}, value: []string{value}}
}

// NewNameSpacedIdentifier creates a new identifier value
//
func NewNameSpacedIdentifier(value []string) Value {
	return &identifier{baseValue: baseValue{info: typeInfoIdentifier}, value: value}
}

// NewStringLiteral creates a new string literal value
//
func NewStringLiteral(value string) Value {
	return &stringLiteral{baseValue: baseValue{info: typeInfoString}, value: value}
}

// NewIntegerLiteral creates a new integer value
//
func NewIntegerLiteral(value int64) Value {
	return &integerLiteral{baseValue: baseValue{info: typeInfoInteger}, value: value}
}

// NewFloatLiteral creates a new integer value
//
func NewFloatLiteral(value float64) Value {
	return &floatLiteral{baseValue: baseValue{info: typeInfoFloat}, value: value}
}

// NewBooleanLiteral creates a new integer value
//
func NewBooleanLiteral(value bool) Value {
	return &booleanLiteral{baseValue: baseValue{info: typeInfoBoolean}, value: value}
}

// NewListValue creates a new list of values
//
func NewListValue(values []Value) Value {
	return &listValue{baseValue: baseValue{info: typeInfoList}, values: values}
}

// NewDictionaryValue creates a new map of values
// TODO: 31okt2016 introduce interface for map parents
//
func NewDictionaryValue(parent interface{}, values map[string]Value) Value {
	if parent == nil {
		return &dictValue{baseValue: baseValue{info: typeInfoDictionary}, parent: nil, values: values}
	}
	return &dictValue{baseValue: baseValue{info: typeInfoDictionary}, parent: parent.(*dictValue), values: values}
}

// NewDictionaryWithBlock constructs a new dictionary by evaluating given block
//
func NewDictionaryWithBlock(context RunContext, block Block) Value {

	// use NewRunContext so block will be evaluated within same scope
	//
	subContext := NewRunContext(context)

	block.Run(subContext, NoArguments)

	return NewDictionaryValue(nil, subContext.Mapping())
}

// NewDictionaryFromList constructs a dictionary from a list of values
// where the list is of the form [key value key value]
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

// NewInternalValue wraps a go value into an elmo value
//
func NewInternalValue(info TypeInfo, value interface{}) Value {
	return &internalValue{baseValue: baseValue{info: info}, value: value}
}

// NewBinaryValueFromInternal constructs a binary value from regular data
// id: TypeId (only for core types)
// typeName: Name of non core typeName
//
func NewBinaryValueFromInternal(id int64, typeName string, value interface{}) BinaryValue {
	var namesBuffer bytes.Buffer
	gob.NewEncoder(&namesBuffer).Encode(value)

	data := &BinaryData{id, typeName, namesBuffer.Bytes()}

	var structBuffer bytes.Buffer
	gob.NewEncoder(&structBuffer).Encode(data)

	return NewBinaryValue(structBuffer.Bytes()).(BinaryValue)
}

// NewErrorValue creates a new Error
//
func NewErrorValue(msg string) ErrorValue {
	return &errorValue{baseValue: baseValue{info: typeInfoError}, msg: msg}
}

// NewBinaryValue creates a new Binary
//
func NewBinaryValue(data []byte) Value {
	return &binaryValue{baseValue: baseValue{info: typeInfoBinary}, data: data}
}

// NewGoFunction creates a new go function
//
func NewGoFunction(name string, value GoFunction) NamedValue {

	splitted := strings.Split(name, "/")
	actualName := splitted[0]
	var help Value = Nothing
	if len(splitted) > 1 {
		help = NewStringLiteral(splitted[1])
	}

	return &goFunction{baseValue: baseValue{info: typeInfoGoFunction}, name: actualName, help: help, value: value}
}

// NewReturnValue creates a new list of values
//
func NewReturnValue(values []Value) Value {
	return &returnValue{baseValue: baseValue{info: typeInfoReturn}, values: values}
}

//
// ---[ASTNode]---------------------------------------------------------------
//
type astNode struct {
	meta  ScriptMetaData
	begin uint32
	end   uint32
}

//
// ---[ARGUMENT]---------------------------------------------------------------
//

type argument struct {
	astNode
	value Value
}

// Argument represent a function call parameter
//
type Argument interface {
	String() string
	Type() Type
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

// NewArgument constructs a new function argument
//
func NewArgument(meta ScriptMetaData, begin uint32, end uint32, value Value) Argument {
	return &argument{astNode: astNode{meta: meta, begin: begin, end: end}, value: value}
}

// NewDynamicArgument constructs a new function argument without script info
//
func NewDynamicArgument(value Value) Argument {
	return &argument{value: value}
}

//
// ---[CALL]-------------------------------------------------------------------
//

type call struct {
	astNode
	baseValue
	firstArgument Argument
	function      GoFunction
	arguments     []Argument
	pipe          Call
}

// Call is a function call
//
type Call interface {
	RunnableValue

	Name() string
	Arguments() []Argument
	WillPipe() bool
}

func (call *call) Name() string {
	if call.function != nil {
		return fmt.Sprintf("%v", call.function)
	}
	return call.firstArgument.String()
}

func (call *call) Arguments() []Argument {
	return call.arguments
}

func (call *call) WillPipe() bool {
	return call.pipe != nil
}

func (call *call) addInfoWhenError(value Value) Value {
	if value.Type() == TypeError {
		if value.(ErrorValue).IsTraced() {
			// TODO: add trace??
			//
			return value
		}

		lineno, _ := call.meta.PositionOf(int(call.begin))
		value.(ErrorValue).SetAt(call.meta, lineno)
	}
	return value
}

func (call *call) pipeResult(context RunContext, value Value) Value {
	if !call.WillPipe() {
		return value
	}

	if value.Type() == TypeReturn {
		values := value.(*returnValue).values
		arguments := make([]Argument, len(values))
		for i, v := range values {
			arguments[i] = &argument{value: v}
		}
		return call.pipe.Run(context, arguments)
	}

	return call.pipe.Run(context, []Argument{&argument{value: value}})
}

func createArgumentsForMissingFunc(context RunContext, call *call, arguments []Argument) []Argument {
	// pass evaluated arguments to the 'func missing' function
	// as a list of values
	//
	values := make([]Value, len(arguments))
	for i, value := range arguments {
		values[i] = EvalArgument(context, value)
	}

	// and pass the original function name as first argument
	//
	return []Argument{
		NewArgument(call.meta, call.astNode.begin, call.astNode.end, NewIdentifier(call.firstArgument.Value().(*identifier).value[len(call.firstArgument.Value().(*identifier).value)-1])),
		NewArgument(call.meta, call.astNode.begin, call.astNode.end, NewListValue(values))}
}

func (call *call) Run(context RunContext, additionalArguments []Argument) Value {

	if call.function != nil {
		return call.pipeResult(context, call.addInfoWhenError(call.function(context, call.Arguments())))
	}

	var inDict DictionaryValue
	var value Value
	var found bool
	var useArguments []Argument

	var function IdentifierValue

	switch call.firstArgument.Type() {
	case TypeCall:
		value = call.firstArgument.Value().(Runnable).Run(context, []Argument{})
		if value.Type() == TypeIdentifier {
			function = value.(IdentifierValue)
			inDict, value, found = function.LookUp(context)
		}
		found = true
	case TypeIdentifier:
		function = call.firstArgument.Value().(IdentifierValue)
		inDict, value, found = function.LookUp(context)
	default:
		value = call.firstArgument.Value()
		found = true
	}

	if additionalArguments != nil && len(additionalArguments) > 0 {
		useArguments = append([]Argument{}, additionalArguments...)
		useArguments = append(useArguments, call.arguments...)
	} else {
		useArguments = call.arguments
	}

	// when call can not be resolved, try to find the 'func missing' function
	//
	if !found {
		if inDict == nil {
			value, found = context.Get("?")
		} else {
			value, found = inDict.Resolve("?")
		}

		if found {
			useArguments = createArgumentsForMissingFunc(context, call, useArguments)
		}
	}

	if found {

		if value == nil {
			return call.pipeResult(context, call.addInfoWhenError(NewErrorValue(fmt.Sprintf("call to %s results in invalid nil value", call.Name()))))
		}

		if inDict != nil {
			this := context.This()
			context.SetThis(inDict.(Value))
			defer func() {
				if this == nil {
					context.SetThis(nil)
				} else {
					context.SetThis(this.(Value))
				}
			}()
		}

		if value.Type() == TypeGoFunction {
			return call.pipeResult(context, call.addInfoWhenError(value.(Runnable).Run(context, useArguments)))
		}

		// runnable values can be used as functions to access their content
		//
		runnable, isRunnable := value.(Runnable)
		if (isRunnable) && (len(call.arguments) > 0) {
			return call.pipeResult(context, call.addInfoWhenError(runnable.Run(context, useArguments)))
		}

		return call.pipeResult(context, call.addInfoWhenError(value))
	}

	return call.pipeResult(context, call.addInfoWhenError(NewErrorValue(fmt.Sprintf("call to undefined \"%s\"", call.firstArgument))))
}

func (call *call) String() string {
	return fmt.Sprintf("(%s ...)", call.Name())
}

func (call *call) Type() Type {
	return TypeCall
}

func (call *call) Internal() interface{} {
	return errors.New("Internal() not implemented on call")
}

// NewCall contstructs a new function call
//
func NewCall(meta ScriptMetaData, begin uint32, end uint32, firstArg Argument, arguments []Argument, pipeTo Call) Call {
	return &call{astNode: astNode{meta: meta, begin: begin, end: end}, baseValue: baseValue{info: typeInfoCall},
		firstArgument: firstArg, arguments: arguments, pipe: pipeTo}
}

// NewCallWithFunction constructs a call that does not need to be resolved
//
func NewCallWithFunction(meta ScriptMetaData, begin uint32, end uint32, function GoFunction, arguments []Argument, pipeTo Call) Call {
	return &call{astNode: astNode{meta: meta, begin: begin, end: end}, baseValue: baseValue{info: typeInfoCall},
		function: function, arguments: arguments, pipe: pipeTo}
}

//
// ---[BLOCK]------------------------------------------------------------------
//

type block struct {
	astNode
	baseValue
	capturedContext RunContext
	calls           []Call
}

// Block is a list of function calls
//
type Block interface {
	// Block can be used as a value
	Value
	// Block can be executed
	Runnable

	Calls() []Call
	CopyWithinContext(RunContext) Block
}

func (block *block) Calls() []Call {
	return block.calls
}

func (block *block) Run(context RunContext, arguments []Argument) Value {
	var result Value = Nothing

	joined := context
	if block.capturedContext != nil {
		joined = joined.Join(block.capturedContext)
	}

	for _, call := range block.calls {
		result = call.Run(joined, []Argument{})
		if joined.isStopped() {
			context.Stop()
			break
		}
		if result.Type() == TypeError {
			if !result.(ErrorValue).CanBeIgnored() {
				context.Stop()
				return result
			}
		}
	}

	return result
}

func (block *block) String() string {
	return fmt.Sprintf("{...}")
}

func (block *block) Type() Type {
	return TypeBlock
}

func (block *block) Internal() interface{} {
	return errors.New("Internal() not implemented on block")
}

func (b *block) CopyWithinContext(context RunContext) Block {
	return &block{astNode: b.astNode, baseValue: b.baseValue, calls: b.calls, capturedContext: context}
}

// NewBlock contsruct a new block of function calls
//
func NewBlock(meta ScriptMetaData, begin uint32, end uint32, calls []Call) Block {
	return &block{astNode: astNode{meta: meta, begin: begin, end: end}, baseValue: baseValue{info: typeInfoBlock}, calls: calls}
}

// EvalArgument evaluates given argument
//
func EvalArgument(context RunContext, argument Argument) Value {

	if argument.Type() == TypeCall {
		return argument.Value().(Runnable).Run(context, NoArguments)
	}

	if argument.Type() == TypeBlock {
		return argument.Value().(Block).CopyWithinContext(context)
	}

	return argument.Value()

}

// EvalArgumentOrSolveIdentifier evaluates given argument
//
func EvalArgumentOrSolveIdentifier(context RunContext, argument Argument) Value {

	if argument.Type() == TypeIdentifier {
		value, found := context.Get(argument.String())
		if found {
			return value
		}
		return NewErrorValue(fmt.Sprintf("could not find %v", argument.String()))
	}

	return EvalArgument(context, argument)

}

// EvalArgumentWithBlock evaluates given argument and if argument is a block
// it will evaluate block content
//
func EvalArgumentWithBlock(context RunContext, argument Argument) Value {

	if argument.Type() == TypeCall || argument.Type() == TypeBlock {
		return argument.Value().(Runnable).Run(context, NoArguments)
	}

	return argument.Value()

}

// EvalArgument2String evaluates given argument and returns it String presentation
//
func EvalArgument2String(context RunContext, argument Argument) string {

	return EvalArgument(context, argument).String()

}
