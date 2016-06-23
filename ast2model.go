package elmo

import (
	"fmt"
	"strconv"
)

// Ast2Block converts an ast node to a code block
//
func Ast2Block(node *node32, meta ScriptMetaData) Block {

	calls := []Call{}

	for _, call := range Children(node) {
		calls = append(calls, Ast2Call(call, meta))
	}

	return NewBlock(meta, node.begin, node.end, calls)
}

// Ast2Call converts an ast node to a function call
//
func Ast2Call(node *node32, meta ScriptMetaData) Call {

	children := Children(node)
	if len(children) == 1 && children[0].pegRule == ruleLine {
		return Ast2Call(children[0], meta)
	}

	functionName := ""
	arguments := []Argument{}

	for idx, argument := range children {
		if idx == 0 {
			functionName = Text(argument, meta.Content())
		} else {
			if argument.pegRule == ruleArgument {
				arguments = append(arguments, Ast2Argument(argument.up, meta))
			}
		}
	}

	return NewCall(meta, node.begin, node.end, functionName, arguments)
}

// Ast2Argument converts an ast node to a function argument
//
func Ast2Argument(node *node32, meta ScriptMetaData) Argument {
	switch node.pegRule {
	case ruleIdentifier:
		return NewArgument(meta, node.begin, node.end, NewIdentifier(Text(node, meta.Content())))
	case ruleStringLiteral:
		txt := Text(node, meta.Content())
		return NewArgument(meta, node.begin, node.end, NewStringLiteral(txt[1:len(txt)-1]))
	case ruleDecimalConstant:
		txt := Text(node, meta.Content())
		i, err := strconv.ParseInt(txt, 10, 64)
		if err != nil {
			panic(err)
		}
		return NewArgument(meta, node.begin, node.end, NewIntegerLiteral(i))
	case ruleFunctionCall:
		return NewArgument(meta, node.begin, node.end, Ast2Call(node, meta))
	case ruleBlock:
		return NewArgument(meta, node.begin, node.end, Ast2Block(node, meta))
	default:
		panic(fmt.Sprintf("invalid argument node: %v", node))
	}
}
