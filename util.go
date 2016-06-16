package elmo

import "testing"

// ParseAndTest will parse given script and execute test function on its result
//
func ParseAndTest(t *testing.T, s string, testfunc func(*node32, string)) {
	grammar := &ElmoGrammar{Buffer: s}

	grammar.Init()

	if err := grammar.Parse(); err != nil {
		t.Error(err)
	} else {
		testfunc(grammar.AST(), s)
	}

}

// ParseAndTestBlock will parse given script to block and execute test function on its result
//
func ParseAndTestBlock(t *testing.T, s string, testfunc func(Block)) {
	ParseAndTest(t, s, func(ast *node32, buf string) {
		block := Ast2Block(ast, buf)
		if block == nil {
			t.Error("no block constructed")
		} else {
			testfunc(block)
		}
	})
}
