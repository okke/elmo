package elmo

import (
	"bytes"
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

	var firstArg Argument
	var arguments = []Argument{}
	var pipeTo Call

	if children[childrenLength-1].pegRule == rulePipedOutput {
		pipeTo = Ast2Call(nodeChildren(children[childrenLength-1])[1], meta)
		children = children[:len(children)-1]
	}

	for idx, argument := range children {
		if idx == 0 {
			if argument.pegRule == ruleAMPERSAND {
				firstArg = NewArgument(meta, argument.begin, argument.end, NewIdentifier(nodeText(argument, meta.Content())))
			} else {
				firstArg = Ast2Argument(argument.up, meta)
			}

		} else {
			if argument.pegRule == ruleArgument {
				arguments = append(arguments, Ast2Argument(argument.up, meta))
			} else if argument.pegRule == ruleCOLON {
				// convert identifier : value => set identifier value
				//
				arguments = append(arguments, firstArg)
				firstArg = NewArgument(meta, argument.begin, argument.end, NewIdentifier("set"))
			} else {
				panic(fmt.Sprintf("unexpected argument %v", argument))
			}
		}
	}

	return NewCall(meta, node.begin, node.end, firstArg, arguments, pipeTo)
}

func escapeString(c rune) rune {
	switch c {
	case 't':
		return '\t'
	case 'n':
		return '\n'
	default:
		return c
	}
}

func escapeLongString(c rune) rune {
	if c == '`' {
		return '`'
	}

	panic(fmt.Sprintf("can not escape multi line string using %c", c))
}

func replaceEscapes(s string, esc rune, escape func(c rune) rune) string {
	var buffer bytes.Buffer
	escaped := false
	for _, c := range s {
		// do something with c
		if c == esc && !escaped {
			escaped = true
		} else {
			if escaped {
				buffer.WriteRune(escape(c))
				escaped = false
			} else {
				buffer.WriteRune(c)
			}
		}
	}
	return buffer.String()
}

// Ast2Argument converts an ast node to a function argument
//
func Ast2Argument(node *node32, meta ScriptMetaData) Argument {
	switch node.pegRule {
	case ruleIdentifier:
		txt := nodeText(node, meta.Content())

		// can be followed by (DOT Identifier)
		//
		dot := node.next
		if dot != nil && dot.pegRule == ruleDOT {
			nextNode := dot.next
			if nextNode != nil && nextNode.pegRule == ruleIdentifier {
				return NewArgument(meta, node.begin, nextNode.end,
					NewNameSpacedIdentifier([]string{txt, nodeText(nextNode, meta.Content())}))
			}
			panic(fmt.Sprintf("expect identifier after namespace separator: %v", node))
		}

		return NewArgument(meta, node.begin, node.end, NewIdentifier(txt))
	case ruleStringLiteral:
		txt := nodeText(node, meta.Content())
		return NewArgument(meta, node.begin, node.end, NewStringLiteral(replaceEscapes(txt[1:len(txt)-1], '\\', escapeString)))
	case ruleLongStringLiteral:
		txt := nodeText(node, meta.Content())
		return NewArgument(meta, node.begin, node.end, NewStringLiteral(replaceEscapes(txt[1:len(txt)-1], '`', escapeLongString)))
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
