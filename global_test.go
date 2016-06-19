package elmo

import "testing"

func TestSetValueIntoGlobalContext(t *testing.T) {

	ParseAndTestBlock(t, "set chipotle \"sauce\"", func(block Block) {

		global := NewGlobalContext()

		block.Run(global, []Argument{})

		result, found := global.Get("chipotle")

		if !found {
			t.Error("expected chipotle to be set")
		} else {
			if result.String() != "sauce" {
				t.Errorf("expected chipotle to be set to (sauce), found %v", result.String())
			}
		}

	})
}

func TestSetValueIntoGlobalContextAndGetIt(t *testing.T) {

	ParseAndTestBlock(t,
		`set chipotle "sauce"
     set sauce (chipotle)`, func(block Block) {

			global := NewGlobalContext()

			block.Run(global, []Argument{})

			result, found := global.Get("sauce")

			if !found {
				t.Error("expected sauce to be set")
			} else {
				if result.String() != "sauce" {
					t.Errorf("expected sauce to be set to (sauce), found %v", result.String())
				}
			}

		})
}

func TestDynamicSetValueIntoGlobalContext(t *testing.T) {

	ParseAndTestBlock(t,
		`set chipotle "sauce"
     set (chipotle) 147`, func(block Block) {

			global := NewGlobalContext()

			block.Run(global, []Argument{})

			result, found := global.Get("sauce")

			if !found {
				t.Error("expected sauce to be set")
			} else {
				if result.String() != "147" {
					t.Errorf("expected sauce to be set to (147), found %v", result.String())
				}
			}

		})
}
