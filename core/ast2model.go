package elmo

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

// Ast2Block converts an ast node to a code block
//
func Ast2Block(node *node32, meta ScriptMetaData) Block {

	calls := []Call{}

	for _, call := range nodeChildren(node) {
		calls = append(calls, Ast2Call(call, meta))
	}

	return NewBlock(meta, node, calls)
}

// Ast2List converts an ast node to a call to list
//
func Ast2List(node *node32, meta ScriptMetaData, pipe Call) Call {
	var arguments = []Argument{}
	for _, argument := range nodeChildren(node) {
		arguments = append(arguments, Ast2Argument(argument.up, meta))
	}

	return NewCallWithFunction(meta, node, ListConstructor, arguments, pipe)
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
	var pipeTo Runnable

	if children[childrenLength-1].pegRule == rulePipedOutput {
		pipeTo = Ast2Call(nodeChildren(children[childrenLength-1])[1], meta)
		children = children[:len(children)-1]
	}

	for idx, argument := range children {
		if idx == 0 {
			if argument.pegRule == ruleAMPERSAND {
				firstArg = NewArgument(meta, argument, NewIdentifier(nodeText(argument, meta.Content())))
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
				firstArg = NewArgument(meta, argument, NewIdentifier("set"))
			} else {
				panic(fmt.Sprintf("unexpected argument %v", argument))
			}
		}
	}

	return NewCall(meta, node, firstArg, arguments, pipeTo)
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

type ast2StringDriver struct {
	quoteRule  pegRule
	stringRule pegRule
	escapeRule pegRule
	escapeFunc func(rune) rune
}

var stringDriver = &ast2StringDriver{
	quoteRule:  ruleQuote,
	stringRule: ruleStringChar,
	escapeRule: ruleEscape,
	escapeFunc: escapeString}

var longStringDriver = &ast2StringDriver{
	quoteRule:  ruleBackTick,
	stringRule: ruleLongStringChar,
	escapeRule: ruleLongEscape,
	escapeFunc: escapeLongString}

// Ast2StringLiteral converts an ast node to a srtring value
//
func Ast2StringLiteral(node *node32, meta ScriptMetaData, driver *ast2StringDriver) Value {

	content := meta.Content()

	var sb strings.Builder

	// keep track of blocks at positions
	//
	var blocks []*blockAtPositionInString

	for _, child := range nodeChildren(node) {
		switch child.pegRule {
		case driver.quoteRule:
			// ignore
		case driver.stringRule:
			grandChildren := nodeChildren(child)
			if grandChildren != nil && len(grandChildren) > 0 && grandChildren[0].pegRule == driver.escapeRule {
				cursor := grandChildren[0].up
				if cursor == nil {
					sb.WriteRune(driver.escapeFunc(rune(nodeText(child, content)[1])))
				} else if cursor.pegRule == ruleBlockWithoutSpacing {
					block := Ast2Block(cursor, meta)
					if blocks == nil {
						blocks = make([]*blockAtPositionInString, 0, 0)
					}

					blocks = append(blocks, &blockAtPositionInString{at: sb.Len(), block: block})

				} else {
					panic("string parsing failed while escaping")
				}
			} else {
				sb.WriteRune(content[child.begin])
			}
		}
	}

	return newStringLiteralWithBlocks(sb.String(), blocks)
}

// Ast2Argument converts an ast node to a function argument
//
func Ast2Argument(node *node32, meta ScriptMetaData) Argument {
	switch node.pegRule {
	case ruleIdentifier:

		parts := []string{}
		begin := node
		end := node
		current := node

		for current != nil {
			end = current
			parts = append(parts, nodeText(current, meta.Content()))

			next := current.next
			if next != nil && next.pegRule == ruleDOT {
				current = next.next
			} else {
				current = nil
			}
		}

		return NewArgumentWithDots(meta, begin, end, NewNameSpacedIdentifier(parts))

	case ruleStringLiteral:
		return NewArgument(meta, node, Ast2StringLiteral(node, meta, stringDriver))
	case ruleLongStringLiteral:
		return NewArgument(meta, node, Ast2StringLiteral(node, meta, longStringDriver))
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

			return NewArgument(meta, node, NewFloatLiteral(f))
		}
		return NewArgument(meta, node, NewIntegerLiteral(i))
	case ruleFunctionCall:
		return NewArgument(meta, node, Ast2Call(node, meta))
	case ruleBlock:
		return NewArgument(meta, node, Ast2Block(node, meta))
	case ruleList:
		return NewArgument(meta, node, Ast2List(node, meta, nil))
	default:
		panic(fmt.Sprintf("invalid argument node: %v", node))
	}
}
