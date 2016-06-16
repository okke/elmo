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
	ParseAndTestBlock(t, "chipotle \"sauce\"", func(block Block) {

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
	ParseAndTestBlock(t, "chipotle 36", func(block Block) {

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
