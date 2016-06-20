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

func TestUserDefinedFunctionWithoutArguments(t *testing.T) {

	ParseAndTestBlock(t,
		`set fsauce (func {
       return "chipotle"
     })
     set sauce (fsauce)`, func(block Block) {

			global := NewGlobalContext()

			block.Run(global, []Argument{})

			result, found := global.Get("sauce")

			if !found {
				t.Error("expected sauce to be set")
			} else {
				if result.String() != "chipotle" {
					t.Errorf("expected sauce to be set to (chipotle), found %v", result.String())
				}
			}

		})
}

func TestUserDefinedFunctionWithOneArgument(t *testing.T) {

	ParseAndTestBlock(t,
		`set fsauce (func pepper {
       return (pepper)
     })
     set sauce (fsauce "chipotle")`, func(block Block) {

			global := NewGlobalContext()

			block.Run(global, []Argument{})

			result, found := global.Get("sauce")

			if !found {
				t.Error("expected sauce to be set")
			} else {
				if result.String() != "chipotle" {
					t.Errorf("expected sauce to be set to (chipotle), found %v", result.String())
				}
			}

		})
}

func TestIfWithoutElse(t *testing.T) {

	ParseTestAndRunBlock(t,
		`set pepper "galapeno"
     if (true) {
      set pepper "chipotle"
     }`, func(context RunContext, blockResult Value) {

			result, found := context.Get("pepper")

			if !found {
				t.Error("expected pepper to be set")
			} else {
				if result.String() != "chipotle" {
					t.Errorf("expected pepper to be set to (chipotle), found %v", result.String())
				}
			}

		})

	ParseTestAndRunBlock(t,
		`set pepper "galapeno"
     if (false) {
      set pepper "chipotle"
     }`, func(context RunContext, blockResult Value) {

			result, found := context.Get("pepper")

			if !found {
				t.Error("expected pepper to be set")
			} else {
				if result.String() != "galapeno" {
					t.Errorf("expected pepper to be set to (galapeno), found %v", result.String())
				}
			}

		})
}

func TestIfWithElse(t *testing.T) {

	ParseTestAndRunBlock(t,
		`set pepper "galapeno"
     if (true) {
      set pepper "chipotle"
     } else {
      set pepper "chilli"
     }`, func(context RunContext, blockResult Value) {

			result, found := context.Get("pepper")

			if !found {
				t.Error("expected pepper to be set")
			} else {
				if result.String() != "chipotle" {
					t.Errorf("expected pepper to be set to (chipotle), found %v", result.String())
				}
			}

		})

	ParseTestAndRunBlock(t,
		`set pepper "galapeno"
     if (false) {
      set pepper "chipotle"
     } else {
      set pepper "chilli"
     }`, func(context RunContext, blockResult Value) {

			result, found := context.Get("pepper")

			if !found {
				t.Error("expected pepper to be set")
			} else {
				if result.String() != "chilli" {
					t.Errorf("expected pepper to be set to (chilli), found %v", result.String())
				}
			}

		})

}
