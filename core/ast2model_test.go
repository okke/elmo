package elmo

import "testing"

func TestCreateBlockWithOneCall(t *testing.T) {
	ParseAndTestBlock(t, "chipotle", func(block Block) {

		if len(block.Calls()) != 1 {
			t.Error("exptected 1 call")
		}

		if block.Calls()[0].Name() != "chipotle" {
			t.Errorf("exptected call to chipotle, got call to %s", block.Calls()[0].Name())
		}

		if len(block.Calls()[0].Arguments()) != 0 {
			t.Error("call should have no arguments")
		}
	})
}

func TestCreateBlockWithTwoCalls(t *testing.T) {
	ParseAndTestBlock(t, "chipotle; sauce", func(block Block) {

		if len(block.Calls()) != 2 {
			t.Error("exptected 2 calls")
		} else {
			if block.Calls()[0].Name() != "chipotle" {
				t.Errorf("exptected call to chipotle, got call to %s", block.Calls()[0].Name())
			}

			if len(block.Calls()[0].Arguments()) != 0 {
				t.Error("call should have no arguments")
			}

			if block.Calls()[1].Name() != "sauce" {
				t.Errorf("exptected call to sauce, got call to %s", block.Calls()[0].Name())
			}

			if len(block.Calls()[1].Arguments()) != 0 {
				t.Error("call should have no arguments")
			}
		}

	})
}

func TestCreateBlockWithOneCallWithOneArgument(t *testing.T) {
	ParseAndTestBlock(t, "chipotle sauce", func(block Block) {

		if len(block.Calls()) != 1 {
			t.Error("exptected 1 call")
		}

		if block.Calls()[0].Name() != "chipotle" {
			t.Errorf("exptected call to chipotle, got call to %s", block.Calls()[0].Name())
		}

		if len(block.Calls()[0].Arguments()) != 1 {
			t.Error("call should one argument")
		}
	})
}

func TestCreateBlockWithOneCallWithOneArgumentAndAComment(t *testing.T) {
	ParseAndTestBlock(t, "chipotle sauce # njam", func(block Block) {

		if len(block.Calls()) != 1 {
			t.Error("exptected 1 call")
		}

		if block.Calls()[0].Name() != "chipotle" {
			t.Errorf("exptected call to chipotle, got call to %s", block.Calls()[0].Name())
		}

		if len(block.Calls()[0].Arguments()) != 1 {
			t.Error("call should one argument")
		}

		if block.Calls()[0].Arguments()[0].String() != "sauce" {
			t.Errorf("exptected call with sauce as argument, got call with %s", block.Calls()[0].Arguments()[0].String())
		}

	})
}

func TestCreateBlockWithOneCallWithTwoArguments(t *testing.T) {
	ParseAndTestBlock(t, "chipotle sauce in_a_jar", func(block Block) {

		if len(block.Calls()) != 1 {
			t.Error("exptected 1 call")
		}

		if block.Calls()[0].Name() != "chipotle" {
			t.Errorf("exptected call to chipotle, got call to %s", block.Calls()[0].Name())
		}

		if len(block.Calls()[0].Arguments()) != 2 {
			t.Error("call should have two arguments")
		}

		if block.Calls()[0].Arguments()[0].String() != "sauce" {
			t.Errorf("exptected argument (sauce), got (%s)", block.Calls()[0].Arguments()[0].String())
		}

		if block.Calls()[0].Arguments()[1].String() != "in_a_jar" {
			t.Errorf("exptected argument (in_a_jar), got (%s)", block.Calls()[0].Arguments()[1].String())
		}
	})
}

func TestCreateBlockWithOneShortcutSet(t *testing.T) {
	ParseAndTestBlock(t, "chipotle : sauce ", func(block Block) {

		if len(block.Calls()) != 1 {
			t.Error("exptected 1 call")
		}

		if block.Calls()[0].Name() != "set" {
			t.Errorf("exptected call to set, got call to %s", block.Calls()[0].Name())
		}

		if len(block.Calls()[0].Arguments()) != 2 {
			t.Error("call should have one arguments")
		}

		if block.Calls()[0].Arguments()[0].String() != "chipotle" {
			t.Errorf("exptected argument (chipotle), got (%s)", block.Calls()[0].Arguments()[0].String())
		}

		if block.Calls()[0].Arguments()[1].String() != "sauce" {
			t.Errorf("exptected argument (sauce), got (%s)", block.Calls()[0].Arguments()[1].String())
		}
	})
}

func TestCreateBlockWithTwoCallsWithTwoArguments(t *testing.T) {
	ParseAndTestBlock(t, "chipotle sauce in_a_jar\npimenton powder in_a_can", func(block Block) {

		if len(block.Calls()) != 2 {
			t.Error("exptected 2 calls")
		}

		if block.Calls()[1].Name() != "pimenton" {
			t.Errorf("exptected call to chipotle, got call to %s", block.Calls()[0].Name())
		}

		if len(block.Calls()[1].Arguments()) != 2 {
			t.Error("call should have two arguments")
		}

		if block.Calls()[1].Arguments()[0].String() != "powder" {
			t.Errorf("exptected argument (powder), got (%s)", block.Calls()[0].Arguments()[0].String())
		}

		if block.Calls()[1].Arguments()[1].String() != "in_a_can" {
			t.Errorf("exptected argument (in_a_can), got (%s)", block.Calls()[1].Arguments()[1].String())
		}
	})
}

func TestCreateBlockWithOneCallWithOneStringArgument(t *testing.T) {
	ParseAndTestBlock(t, "chipotle \"sauce\" #bla", func(block Block) {

		if len(block.Calls()) != 1 {
			t.Error("exptected 1 call")
		}

		if block.Calls()[0].Name() != "chipotle" {
			t.Errorf("exptected call to chipotle, got call to %s", block.Calls()[0].Name())
		}

		if len(block.Calls()[0].Arguments()) != 1 {
			t.Error("call should have one argument")
		}

		if block.Calls()[0].Arguments()[0].Type() != TypeString {
			t.Errorf("exptected a string argument")
		}

		if block.Calls()[0].Arguments()[0].String() != "sauce" {
			t.Errorf("exptected argument (sauce), got (%s)", block.Calls()[0].Arguments()[0].String())
		}

	})
}

func TestCreateBlockWithOneCallWithOneIntegerArgument(t *testing.T) {
	ParseAndTestBlock(t, "chipotle 36 # and more", func(block Block) {

		if len(block.Calls()) != 1 {
			t.Error("exptected 1 call")
		}

		if block.Calls()[0].Name() != "chipotle" {
			t.Errorf("exptected call to chipotle, got call to %s", block.Calls()[0].Name())
		}

		if len(block.Calls()[0].Arguments()) != 1 {
			t.Error("call should have one argument")
		}

		if block.Calls()[0].Arguments()[0].Type() != TypeInteger {
			t.Errorf("exptected an integer argument")
		}

		if block.Calls()[0].Arguments()[0].String() != "36" {
			t.Errorf("exptected argument (36), got (%s)", block.Calls()[0].Arguments()[0].String())
		}
	})
}

func TestCreateBlockWithOneCallWithOneBlockArgument(t *testing.T) {
	ParseAndTestBlock(t, "chipotle {}", func(block Block) {

		if len(block.Calls()) != 1 {
			t.Error("exptected 1 call")
		}

		if block.Calls()[0].Name() != "chipotle" {
			t.Errorf("exptected call to chipotle, got call to %s", block.Calls()[0].Name())
		}

		if len(block.Calls()[0].Arguments()) != 1 {
			t.Error("call should have one argument")
		}

		if block.Calls()[0].Arguments()[0].Type() != TypeBlock {
			t.Errorf("exptected a block argument")
		}

	})
}

func TestCreateBlockWithOneCallWithOneCallArgument(t *testing.T) {
	ParseAndTestBlock(t, "chipotle (sauce 33)", func(block Block) {

		if len(block.Calls()) != 1 {
			t.Error("exptected 1 call")
		}

		if block.Calls()[0].Name() != "chipotle" {
			t.Errorf("exptected call to chipotle, got call to %s", block.Calls()[0].Name())
		}

		if len(block.Calls()[0].Arguments()) != 1 {
			t.Error("call should have one argument")
		}

		if block.Calls()[0].Arguments()[0].Type() != TypeCall {
			t.Errorf("exptected a call argument")
		}

		call := block.Calls()[0].Arguments()[0].Value().(Call)
		if call.Name() != "sauce" {
			t.Errorf("exptected a call to (sauce) instead of %v", call.Name())
		}

		if len(call.Arguments()) != 1 {
			t.Errorf("exptected 1 argument for call to (sauce) instead of %d", len(call.Arguments()))
		}

		if call.Arguments()[0].Type() != TypeInteger {
			t.Errorf("exptected an integer argument for call to (sauce) instead of %d", call.Arguments()[0].Type())
		}

	})
}

func TestCreateBlockWithOneCallWithOneShortcutCallArgument(t *testing.T) {
	ParseAndTestBlock(t, "chipotle $sauce", func(block Block) {
	})
}

func TestCreateBlockWithOnePipedCall(t *testing.T) {
	ParseAndTestBlock(t, "chipotle | galapeno", func(block Block) {

		if len(block.Calls()) != 1 {
			t.Error("exptected 1 call")
		}

		call := block.Calls()[0]
		if call.Name() != "chipotle" {
			t.Errorf("exptected call to chipotle, got call to %s", block.Calls()[0].Name())
		}

		if !call.WillPipe() {
			t.Errorf("exptected the result of call to chipotle to be piped to next call")
		}
	})
}
