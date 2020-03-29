package elmo

import (
	"fmt"
	"path/filepath"
	"testing"
)

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

// Parse2Block will parse given script and return it as block or error
//
func Parse2Block(s string, fileName string) (Block, error) {
	grammar := &ElmoGrammar{Buffer: s}

	grammar.Init()

	if err := grammar.Parse(); err != nil {
		return nil, err
	}

	return Ast2Block(grammar.AST(), NewScriptMetaData(fileName, s)), nil
}

// ParseAndRun will parse given script and execute test function on its result
//
func ParseAndRun(context RunContext, s string) Value {
	return ParseAndRunWithFile(context, s, "stdin")
}

// ParseAndRunWithFile will parse given script and execute test function on its result
//
func ParseAndRunWithFile(context RunContext, s string, fileName string) (val Value) {

	defer func() {
		if r := recover(); r != nil {
			val = NewErrorValue(fmt.Sprintf("%v", r))
		}
	}()

	absPath, err := filepath.Abs(fileName)

	if err != nil {
		return NewErrorValue(fmt.Sprintf("could not get absolute path of %s:%v", fileName, err))
	}

	currentScript := context.ScriptName()

	context.SetScriptName(NewStringLiteral(absPath))

	grammar := &ElmoGrammar{Buffer: s}

	grammar.Init()

	if err := grammar.Parse(); err != nil {
		parseErr, isParseError := err.(*parseError)
		if isParseError {
			return NewErrorValueWithToken(err.Error(), &parseErr.max)
		}
		return NewErrorValue(err.Error())
	}

	block := Ast2Block(grammar.AST(), NewScriptMetaData(fileName, s))
	result := block.Run(context, NoArguments)
	context.SetScriptName(currentScript)
	return result

}

// ParseAndTestBlock will parse given script to block and execute test function on its result
//
func ParseAndTestBlock(t *testing.T, s string, testfunc ...func(Block)) {
	ParseAndTest(t, s, func(ast *node32) {
		block := Ast2Block(ast, NewScriptMetaData("test", s))
		if block == nil {
			t.Error("no block constructed")
		} else {
			for _, f := range testfunc {
				f(block)
			}
		}
	})
}

// ParseTestAndRunBlock will parse given script to block, run it and execute test function on its result
//
func ParseTestAndRunBlock(t *testing.T, s string, testfunc ...func(RunContext, Value)) {

	ParseAndTestBlock(t, s, func(block Block) {
		global := NewGlobalContext()
		result := block.Run(global, []Argument{})

		for _, f := range testfunc {
			f(global, result)
		}

	})

}

// ParseTestAndRunBlockWithinContext will parse given script to block, run it within given context and execute test function on its result
//
func ParseTestAndRunBlockWithinContext(t *testing.T, context RunContext, s string, testfunc ...func(RunContext, Value)) {

	ParseAndTestBlock(t, s, func(block Block) {
		result := block.Run(context, []Argument{})

		for _, f := range testfunc {
			f(context, result)
		}

	})

}

// run an elmo test file located in the ./test directory
//
func TestMoFile(t *testing.T, mo string, contextInitializer func(RunContext)) {
	context := NewGlobalContext()
	contextInitializer(context)

	loader := NewLoader(context, []string{"./test"})

	loaded := loader.Load(mo)

	if loaded.Type() == TypeError {
		t.Error("error while loading", mo, ":", loaded.String())
	}
}
