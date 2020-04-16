package elmo

import (
	"fmt"
)

// GoFunction is a native go function that takes an array of input values
// and returns an output value
//
type GoFunction func(RunContext, []Argument) Value

type goFunction struct {
	baseValue
	name  string
	help  Value
	value GoFunction
	block Block
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

type inspectableGoFunction struct {
	goFunction
	argNames []string
}

func (inspectableGoFunction *inspectableGoFunction) Block() Block {
	return inspectableGoFunction.goFunction.block
}

func (inspectableGoFunction *inspectableGoFunction) Meta() ScriptMetaData {
	return inspectableGoFunction.goFunction.block.Meta()
}

func (inspectableGoFunction *inspectableGoFunction) BeginsAt() uint32 {
	return inspectableGoFunction.goFunction.block.BeginsAt()
}

func (inspectableGoFunction *inspectableGoFunction) EndsAt() uint32 {
	return inspectableGoFunction.goFunction.block.EndsAt()
}

func (inspectableGoFunction *inspectableGoFunction) Enrich(dict DictionaryValue) {
	dict.Set(NewStringLiteral("arguments"), NewListValueFromStrings(inspectableGoFunction.argNames))
}

// NewGoFunctionWithBlock creates a new go function and stores the block of code for later inspection
//
func NewGoFunctionWithBlock(name string, help string, value GoFunction, argNames []string, block Block) NamedValue {
	return &inspectableGoFunction{goFunction: goFunction{
		baseValue: baseValue{info: typeInfoGoFunction},
		name:      name,
		help:      NewStringLiteral(help),
		value:     value,
		block:     block}, argNames: argNames}
}

// NewGoFunctionWithHelp creates a new go function
//
func NewGoFunctionWithHelp(name string, help string, value GoFunction) NamedValue {

	return &goFunction{baseValue: baseValue{info: typeInfoGoFunction}, name: name, help: NewStringLiteral(help), value: value}
}
