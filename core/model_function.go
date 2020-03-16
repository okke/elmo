package elmo

import (
	"fmt"
	"strings"
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

func (goFunction *goFunction) Block() Block {
	return goFunction.block
}

// NewGoFunction creates a new go function
//
func NewGoFunction(name string, value GoFunction) NamedValue {

	splitted := strings.SplitN(name, "/", 2)
	actualName := splitted[0]
	var help string = ""
	if len(splitted) > 1 {
		help = splitted[1]
	}

	return NewGoFunctionWithHelp(actualName, help, value)
}

// NewGoFunctionWithBlock creates a new go function and stores the block of code for later inspection
//
func NewGoFunctionWithBlock(name string, help string, value GoFunction, block Block) NamedValue {
	return &goFunction{
		baseValue: baseValue{info: typeInfoGoFunction},
		name:      name,
		help:      NewStringLiteral(help),
		value:     value,
		block:     block}
}

// NewGoFunctionWithHelp creates a new go function
//
func NewGoFunctionWithHelp(name string, help string, value GoFunction) NamedValue {

	return &goFunction{baseValue: baseValue{info: typeInfoGoFunction}, name: name, help: NewStringLiteral(help), value: value}
}
