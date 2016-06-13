package elmo

import "testing"

func TestCreateBlockWithOneCall(t *testing.T) {
	ParseAndTest(t, "chipotle", func(ast *node32, buf string) {
		block := Ast2Block(ast, buf)
		if block == nil {
			t.Error("no block constructed")
		}

		if len(block.Calls()) != 1 {
			t.Error("exptected 1 call")
		}

		if block.Calls()[0].Name() != "chipotle" {
			t.Error("exptected call to chipotle")
		}
	})
}

func TestCreateBlockWithTwoCalls(t *testing.T) {
	ParseAndTest(t, "chipotle; sauce", func(ast *node32, buf string) {
		block := Ast2Block(ast, buf)
		if block == nil {
			t.Error("no block constructed")
		}

		if len(block.Calls()) != 2 {
			t.Error("exptected 2 calls")
		} else {
			if block.Calls()[0].Name() != "chipotle" {
				t.Error("exptected call to chipotle")
			}

			if block.Calls()[1].Name() != "sauce" {
				t.Error("exptected call to sauce")
			}
		}

	})
}
