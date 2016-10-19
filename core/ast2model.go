package elmo

import (
	"fmt"
	"strconv"
)

// Ast2Block converts an ast node to a code block
//
func Ast2Block(node *node32, meta ScriptMetaData) Block {

	calls := []Call{}

	for _, call := range nodeChildren(node) {
		calls = append(calls, Ast2Call(call, meta))
	}

	return NewBlock(meta, node.begin, node.end, calls)
}

// Ast2List converts an ast node to a call to list
//
func Ast2List(node *node32, meta ScriptMetaData, pipe Call) Call {
	var arguments = []Argument{}
	for _, argument := range nodeChildren(node) {
		arguments = append(arguments, Ast2Argument(argument.up, meta))
	}
	return NewCallWithFunction(meta, node.begin, node.end, ListConstructor, arguments, pipe)
}

// Ast2Call converts an ast node to a function call
//
func Ast2Call(node *node32, meta ScriptMetaData) Call {

	children := nodeChildren(node)
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
		pipeTo = Ast2Call(nodeChildren(children[childrenLength-1])[1], meta)
	}

	// when call does not start with an identifier, it can be a literal without arguments
	//
	if (childrenLength == 1 && pipeTo == nil) || (childrenLength == 2 && pipeTo != nil) {
		if children[0].up.pegRule == ruleList {
			return Ast2List(children[0].up, meta, pipeTo)
		}

		if children[0].up.pegRule != ruleIdentifier {
			panic(fmt.Sprintf("invalid call %v: %v", children[0].up, nodeText(children[0].up, meta.Content())))
		}
	}

	for idx, argument := range children {
		if idx == 0 {
			functionArg = argument.up
			if functionArg.pegRule == ruleIdentifier {
				functionName = append(functionName, nodeText(argument, meta.Content()))
			} else {
				panic(fmt.Sprintf("found non identifier %v: %v", argument, nodeText(argument, meta.Content())))
			}

		} else {
			if argument.pegRule == ruleArgument {
				if appendToFunctionName {
					functionName = append(functionName, nodeText(argument, meta.Content()))
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
		return NewArgument(meta, node.begin, node.end, NewIdentifier(nodeText(node, meta.Content())))
	case ruleStringLiteral:
		txt := nodeText(node, meta.Content())
		return NewArgument(meta, node.begin, node.end, NewStringLiteral(txt[1:len(txt)-1]))
	case ruleNumber:
		txt := nodeText(node, meta.Content())

		// first try if its an integer value
		//
		i, err := strconv.ParseInt(txt, 10, 64)
		if err != nil {

			// then try parsing as float
			//
			f, err := strconv.ParseFloat(txt, 64)
			if err != nil {
				panic(err)
			}

			return NewArgument(meta, node.begin, node.end, NewFloatLiteral(f))
		}
		return NewArgument(meta, node.begin, node.end, NewIntegerLiteral(i))
	case ruleFunctionCall:
		return NewArgument(meta, node.begin, node.end, Ast2Call(node, meta))
	case ruleBlock:
		return NewArgument(meta, node.begin, node.end, Ast2Block(node, meta))
	case ruleList:
		return NewArgument(meta, node.begin, node.end, Ast2List(node, meta, nil))
	default:
		panic(fmt.Sprintf("invalid argument node: %v", node))
	}
}
