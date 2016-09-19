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

// Ast2List converts an ast node to a call to list
//
func Ast2List(node *node32, meta ScriptMetaData) Call {
	var arguments = []Argument{}
	for _, argument := range Children(node) {
		arguments = append(arguments, Ast2Argument(argument.up, meta))
	}
	return NewCall(meta, node.begin, node.end, []string{"list"}, arguments, nil)
}

// Ast2Call converts an ast node to a function call
//
func Ast2Call(node *node32, meta ScriptMetaData) Call {

	children := Children(node)
	childrenLength := len(children)

	if childrenLength == 1 && children[0].pegRule == ruleLine {
		return Ast2Call(children[0], meta)
	}

	var functionName = []string{}
	var functionArg *node32
	var arguments = []Argument{}
	var appendToFunctionName = false

	var pipeTo Call

	if children[childrenLength-1].pegRule == rulePipedOutput {
		pipeTo = Ast2Call(Children(children[childrenLength-1])[1], meta)
	}

	for idx, argument := range children {
		if idx == 0 {
			functionArg = argument
			functionName = append(functionName, Text(argument, meta.Content()))
		} else {
			if argument.pegRule == ruleArgument {
				if appendToFunctionName {
					functionName = append(functionName, Text(argument, meta.Content()))
					appendToFunctionName = false
				} else {
					arguments = append(arguments, Ast2Argument(argument.up, meta))
				}
			} else if argument.pegRule == ruleShortcut {

				cut := argument.up

				if cut.pegRule == ruleCOLON {
					// convert identifier : value => set identifier value
					//
					functionName = []string{"set"}
					arguments = append(arguments, Ast2Argument(functionArg, meta))
				} else if cut.pegRule == ruleDOT {
					// convert identifier . value => (identifier value)
					appendToFunctionName = true
				} else {
					panic(fmt.Sprintf("could not create shortcut call for %v", cut))
				}

			}
		}
	}

	return NewCall(meta, node.begin, node.end, functionName, arguments, pipeTo)
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
	case ruleList:
		return NewArgument(meta, node.begin, node.end, Ast2List(node, meta))
	default:
		panic(fmt.Sprintf("invalid argument node: %v", node))
	}
}
