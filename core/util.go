package elmo

import (
	"fmt"
	"reflect"
	"testing"
)

// ExpectValueSetTo expects a given variable is set to a given value
//
func ExpectValueSetTo(t *testing.T, key string, value string) func(RunContext, Value) {

	return func(context RunContext, blockResult Value) {

		if blockResult.Type() == TypeError {
			t.Error(blockResult.(ErrorValue).Error())
			return
		}

		result, found := context.Get(key)

		if !found {
			t.Errorf("expected %s to be set", key)
		} else {
			if result.String() != value {
				t.Errorf("expected %s to be set to (%s), found %s", key, value, result.String())
			}
		}
	}
}

// ExpectErrorValueAt ecpects an error on a given line number
//
func ExpectErrorValueAt(t *testing.T, lineno int) func(RunContext, Value) {

	return func(context RunContext, blockResult Value) {

		if blockResult.Type() != TypeError {
			t.Errorf("expected error but found %v", blockResult)
			return
		}

		_, l := blockResult.(ErrorValue).At()

		if l != lineno {
			t.Errorf("expected error at line %d but found (%v) on line %d", lineno, blockResult.String(), l)
		}

	}
}

// ExpectNothing expects evauation returns Nothing
//
func ExpectNothing(t *testing.T) func(RunContext, Value) {

	return func(context RunContext, blockResult Value) {
		if blockResult != Nothing {
			t.Errorf("expected nothing but found %v", blockResult)
		}
	}
}

// ExpectValue expects evaluation returns a given value
//
func ExpectValue(t *testing.T, value Value) func(RunContext, Value) {

	return func(context RunContext, blockResult Value) {
		if !reflect.DeepEqual(blockResult, value) {
			t.Errorf("expected (%v) but found (%v)", value, blockResult)
		}
	}
}

// ParseAndTest will parse given script and execute test function on its result
//
func ParseAndTest(t *testing.T, s string, testfunc func(*node32)) {
	grammar := &ElmoGrammar{Buffer: s}

	grammar.Init()

	if err := grammar.Parse(); err != nil {
		t.Error(err)
	} else {
		testfunc(grammar.AST())
	}

}

// ParseAndRun will parse given script and execute test function on its result
//
func ParseAndRun(context RunContext, s string) Value {
	grammar := &ElmoGrammar{Buffer: s}

	grammar.Init()

	if err := grammar.Parse(); err != nil {
		panic(fmt.Sprintf("could not parse: %v", err))
	} else {
		block := Ast2Block(grammar.AST(), NewScriptMetaData("", s))
		return block.Run(context, NoArguments)
	}

}

// ParseAndTestBlock will parse given script to block and execute test function on its result
//
func ParseAndTestBlock(t *testing.T, s string, testfunc func(Block)) {
	ParseAndTest(t, s, func(ast *node32) {
		block := Ast2Block(ast, NewScriptMetaData("test", s))
		if block == nil {
			t.Error("no block constructed")
		} else {
			testfunc(block)
		}
	})
}

// ParseTestAndRunBlock will parse given script to block, run it and execute test function on its result
//
func ParseTestAndRunBlock(t *testing.T, s string, testfunc func(RunContext, Value)) {

	ParseAndTestBlock(t, s, func(block Block) {
		global := NewGlobalContext()
		result := block.Run(global, []Argument{})

		testfunc(global, result)
	})

}

// ParseTestAndRunBlockWithinContext will parse given script to block, run it within given context and execute test function on its result
//
func ParseTestAndRunBlockWithinContext(t *testing.T, context RunContext, s string, testfunc func(RunContext, Value)) {

	ParseAndTestBlock(t, s, func(block Block) {
		result := block.Run(context, []Argument{})
		testfunc(context, result)
	})

}
