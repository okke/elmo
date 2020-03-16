package elmo

import (
	"errors"
	"fmt"
)

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
	return fmt.Sprintf("{...%#v}", block)
}

func (block *block) Type() Type {
	return TypeBlock
}

func (block *block) Internal() interface{} {
	return errors.New("Internal() not implemented on block")
}

func (block *block) Enrich(dict DictionaryValue) {

}

func (b *block) CopyWithinContext(context RunContext) Block {
	return &block{astNode: b.astNode, baseValue: b.baseValue, calls: b.calls, capturedContext: context}
}

// NewBlock contsruct a new block of function calls
//
func NewBlock(meta ScriptMetaData, node *node32, calls []Call) Block {
	return &block{astNode: astNode{meta: meta, node: node}, baseValue: baseValue{info: typeInfoBlock}, calls: calls}
}
