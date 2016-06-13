package elmo

type argument struct {
}

// Argument represent a function call parameter
//
type Argument interface {
}

type call struct {
	functionName string
	arguments    []Argument
}

// Call is a function call
//
type Call interface {
	Name() string
	Arguments() []Argument
}

func (call *call) Name() string {
	return call.functionName
}

func (call *call) Arguments() []Argument {
	return call.arguments
}

// NewCall contstructs a new function call
//
func NewCall(name string, arguments []Argument) Call {
	return &call{functionName: name, arguments: arguments}
}

type block struct {
	calls []Call
}

// Block is a list of function calls
//
type Block interface {
	Calls() []Call
}

func (block *block) Calls() []Call {
	return block.calls
}

// NewBlock contsruct a new block of function calls
//
func NewBlock(calls []Call) Block {
	return &block{calls: calls}
}
