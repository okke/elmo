package elmo

import (
	"fmt"
	"strconv"
)

// Ast2Block converts an ast node to a code block
//
func Ast2Block(node *node32, buf string) Block {

	calls := []Call{}

	for _, call := range Children(node) {
		calls = append(calls, Ast2Call(call, buf))
	}

	return NewBlock(calls)
}

// Ast2Call converts an ast node to a function call
//
func Ast2Call(node *node32, buf string) Call {

	children := Children(node)
	if len(children) == 1 && children[0].pegRule == ruleLine {
		return Ast2Call(children[0], buf)
	}

	functionName := ""
	arguments := []Argument{}

	for idx, argument := range children {
		if idx == 0 {
			functionName = Text(argument, buf)
		} else {
			if argument.pegRule == ruleArgument {
				arguments = append(arguments, Ast2Argument(argument.up, buf))
			}
		}
	}

	return NewCall(functionName, arguments)
}

// Ast2Argument converts an ast node to a function argument
//
func Ast2Argument(node *node32, buf string) Argument {
	switch node.pegRule {
	case ruleIdentifier:
		return NewArgument(NewIdentifier(Text(node, buf)))
	case ruleStringLiteral:
		txt := Text(node, buf)
		return NewArgument(NewStringLiteral(txt[1 : len(txt)-1]))
	case ruleDecimalConstant:
		txt := Text(node, buf)
		i, err := strconv.ParseInt(txt, 10, 64)
		if err != nil {
			panic(err)
		}
		return NewArgument(NewIntegerLiteral(i))
	case ruleFunctionCall:
		return NewArgument(Ast2Call(node, buf))
	case ruleBlock:
		return NewArgument(Ast2Block(node, buf))
	default:
		panic(fmt.Sprintf("invalid argument node: %v", node))
	}
}
