package elmo

import "fmt"

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
	// TypeBoolean represents a type for a boolean value
	TypeBoolean
	// TypeError represents a type for an error value
	TypeError
	// TypeBlock represents a type for a code block
	TypeBlock
	// TypeCall represent the type for a function call
	TypeCall
	// TypeGoFunction represents a type for an internal go function
	TypeGoFunction
	// TypeNil represents the type of a nil value
	TypeNil
)

type nothing struct {
}

// Nothing represents nil
//
var Nothing = &nothing{}

type identifier struct {
	value string
}

type stringLiteral struct {
	value string
}

type integerLiteral struct {
	value int64
}

type booleanLiteral struct {
	value bool
}

type errorValue struct {
	meta   ScriptMetaData
	lineno int
	msg    string
}

// GoFunction is a native go function that takes an array of input values
// and returns an output value
//
type GoFunction func(RunContext, []Argument) Value

type goFunction struct {
	name  string
	value GoFunction
}

// Value represents data within elmo
//
type Value interface {
	Print() string
	String() string
	Type() Type
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

// NamedValue represent data with a name
//
type NamedValue interface {
	Value
	Name() string
}

func (nothing *nothing) Print() string {
	return "nil"
}

func (nothing *nothing) String() string {
	return "nil"
}

func (nothing *nothing) Type() Type {
	return TypeNil
}

func (identifier *identifier) Print() string {
	return identifier.value
}

func (identifier *identifier) String() string {
	return identifier.value
}

func (identifier *identifier) Type() Type {
	return TypeIdentifier
}

func (stringLiteral *stringLiteral) Print() string {
	return fmt.Sprintf("\"%s\"", stringLiteral.value)
}

func (stringLiteral *stringLiteral) String() string {
	return fmt.Sprintf("%s", stringLiteral.value)
}

func (stringLiteral *stringLiteral) Type() Type {
	return TypeString
}

func (integerLiteral *integerLiteral) Print() string {
	return fmt.Sprintf("%d", integerLiteral.value)
}

func (integerLiteral *integerLiteral) String() string {
	return fmt.Sprintf("%d", integerLiteral.value)
}

func (integerLiteral *integerLiteral) Type() Type {
	return TypeInteger
}

func (booleanLiteral *booleanLiteral) Print() string {
	return booleanLiteral.String()
}

func (booleanLiteral *booleanLiteral) String() string {
	return fmt.Sprintf("%v", booleanLiteral.value)
}

func (booleanLiteral *booleanLiteral) Type() Type {
	return TypeBoolean
}

func (errorValue *errorValue) Print() string {
	return errorValue.String()
}

func (errorValue *errorValue) String() string {
	if errorValue.meta != nil {
		meta, lineno := errorValue.At()
		return fmt.Sprintf("error at %s at line %d: %s", meta.Name(), lineno, errorValue.msg)
	}
	return fmt.Sprintf("error: %s", errorValue.msg)
}

func (errorValue *errorValue) Type() Type {
	return TypeError
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

func (goFunction *goFunction) Print() string {
	return goFunction.String()
}

func (goFunction *goFunction) String() string {
	return fmt.Sprintf("go:%s", goFunction.name)
}

func (goFunction *goFunction) Type() Type {
	return TypeGoFunction
}

func (goFunction *goFunction) Name() string {
	return goFunction.name
}

func (goFunction *goFunction) Run(context RunContext, arguments []Argument) Value {
	return goFunction.value(context, arguments)
}

// NewIdentifier creates a new identifier value
//
func NewIdentifier(value string) Value {
	return &identifier{value: value}
}

// NewStringLiteral creates a new string literal value
//
func NewStringLiteral(value string) Value {
	return &stringLiteral{value: value}
}

// NewIntegerLiteral creates a new integer value
//
func NewIntegerLiteral(value int64) Value {
	return &integerLiteral{value: value}
}

// NewBooleanLiteral creates a new integer value
//
func NewBooleanLiteral(value bool) Value {
	return &booleanLiteral{value: value}
}

// NewErrorValue creates a new Error
//
func NewErrorValue(msg string) ErrorValue {
	return &errorValue{msg: msg}
}

// NewGoFunction creates a new go function
//
func NewGoFunction(name string, value GoFunction) NamedValue {
	return &goFunction{name: name, value: value}
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
	functionName string
	arguments    []Argument
}

// Call is a function call
//
type Call interface {
	// Call can be used
	Value
	// Call can be executed
	Runnable

	Name() string
	Arguments() []Argument
}

func (call *call) Name() string {
	return call.functionName
}

func (call *call) Arguments() []Argument {
	return call.arguments
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

func (call *call) Run(context RunContext, arguments []Argument) Value {
	value, found := context.Get(call.functionName)

	if found {
		if value.Type() == TypeGoFunction {
			// TODO How to handle arguments
			//
			return call.addInfoWhenError(value.(Runnable).Run(context, call.arguments))
		}
		return call.addInfoWhenError(value)
	}

	panic(fmt.Sprintf("call to undefined \"%s\"", call.functionName))
	// return Nothing
}

func (call *call) Print() string {
	return call.String()
}

func (call *call) String() string {
	return fmt.Sprintf("(%s ...)", call.functionName)
}

func (call *call) Type() Type {
	return TypeCall
}

// NewCall contstructs a new function call
//
func NewCall(meta ScriptMetaData, begin uint32, end uint32, name string, arguments []Argument) Call {
	return &call{astNode: astNode{meta: meta, begin: begin, end: end}, functionName: name, arguments: arguments}
}

//
// ---[BLOCK]------------------------------------------------------------------
//

type block struct {
	astNode
	calls []Call
}

// Block is a list of function calls
//
type Block interface {
	// Block can be used as a value
	Value
	// Block can be executed
	Runnable
	Calls() []Call
}

func (block *block) Calls() []Call {
	return block.calls
}

func (block *block) Run(context RunContext, arguments []Argument) Value {
	var result Value = Nothing

	for _, call := range block.calls {
		result = call.Run(context, []Argument{})
	}

	return result
}

func (block *block) Print() string {
	return block.String()
}

func (block *block) String() string {
	return fmt.Sprintf("{...}")
}

func (block *block) Type() Type {
	return TypeBlock
}

// NewBlock contsruct a new block of function calls
//
func NewBlock(meta ScriptMetaData, begin uint32, end uint32, calls []Call) Block {
	return &block{astNode: astNode{meta: meta, begin: begin, end: end}, calls: calls}
}
