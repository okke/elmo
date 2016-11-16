package elmo

import (
	"errors"
	"fmt"
	"math"
	"strings"
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
	info TypeInfo
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
	values []Value
}

type returnValue struct {
	baseValue
	values []Value
}

type dictValue struct {
	baseValue
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
	Compare(Value) (int, ErrorValue)
}

// HelpValue represents a value with help
//
type HelpValue interface {
	Help() Value
}

// ErrorValue represents an Error
//
type ErrorValue interface {
	Value
	Error() string
	SetAt(meta ScriptMetaData, lineno int)
	At() (ScriptMetaData, int)
	IsTraced() bool
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

func (nothing *nothing) String() string {
	return "nil"
}

func (nothing *nothing) Type() Type {
	return TypeNil
}

func (nothing *nothing) Internal() interface{} {
	return nil
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

func (integerLiteral *integerLiteral) Compare(value Value) (int, ErrorValue) {
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

func (floatLiteral *floatLiteral) String() string {
	return fmt.Sprintf("%f", floatLiteral.value)
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

func (floatLiteral *floatLiteral) Compare(value Value) (int, ErrorValue) {
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

func (booleanLiteral *booleanLiteral) String() string {
	return fmt.Sprintf("%v", booleanLiteral.value)
}

func (booleanLiteral *booleanLiteral) Type() Type {
	return TypeBoolean
}

func (booleanLiteral *booleanLiteral) Internal() interface{} {
	return booleanLiteral.value
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

func (dictValue *dictValue) Run(context RunContext, arguments []Argument) Value {

	key := EvalArgument(context, arguments[0])
	result, _ := dictValue.Resolve(key.String())
	return result
}

func (errorValue *errorValue) String() string {
	if errorValue.meta != nil {
		meta, lineno := errorValue.At()
		return fmt.Sprintf("error at %s at line %d: %s", meta.Name(), lineno, errorValue.msg)
	}
	return fmt.Sprintf("error: %s", errorValue.msg)
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

// NewErrorValue creates a new Error
//
func NewErrorValue(msg string) ErrorValue {
	return &errorValue{baseValue: baseValue{info: typeInfoError}, msg: msg}
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
			return result
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
