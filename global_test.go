package elmo

import "testing"

func expectValueSetTo(t *testing.T, key string, value string) func(RunContext, Value) {

	return func(context RunContext, blockResult Value) {

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

func TestSetValueIntoGlobalContext(t *testing.T) {

	ParseTestAndRunBlock(t, "set chipotle \"sauce\"", expectValueSetTo(t, "chipotle", "sauce"))
}

func TestSetValueIntoGlobalContextAndGetIt(t *testing.T) {

	ParseTestAndRunBlock(t,
		`set chipotle "sauce"
     set sauce (chipotle)`, expectValueSetTo(t, "sauce", "sauce"))
}

func TestDynamicSetValueIntoGlobalContext(t *testing.T) {

	ParseTestAndRunBlock(t,
		`set chipotle "sauce"
     set (chipotle) 147`, expectValueSetTo(t, "sauce", "147"))
}

func TestUserDefinedFunctionWithoutArguments(t *testing.T) {

	ParseTestAndRunBlock(t,
		`set fsauce (func {
       return "chipotle"
     })
     set sauce (fsauce)`, expectValueSetTo(t, "sauce", "chipotle"))
}

func TestUserDefinedFunctionWithOneArgument(t *testing.T) {

	ParseTestAndRunBlock(t,
		`set fsauce (func pepper {
       return (pepper)
     })
     set sauce (fsauce "chipotle")`, expectValueSetTo(t, "sauce", "chipotle"))
}

func TestIfWithoutElse(t *testing.T) {

	ParseTestAndRunBlock(t,
		`set pepper "galapeno"
     if (true) {
      set pepper "chipotle"
     }`, expectValueSetTo(t, "pepper", "chipotle"))

	ParseTestAndRunBlock(t,
		`set pepper "galapeno"
     if (false) {
      set pepper "chipotle"
     }`, expectValueSetTo(t, "pepper", "galapeno"))
}

func TestIfWithElse(t *testing.T) {

	ParseTestAndRunBlock(t,
		`set pepper "galapeno"
     if (true) {
      set pepper "chipotle"
     } else {
      set pepper "chilli"
     }`, expectValueSetTo(t, "pepper", "chipotle"))

	ParseTestAndRunBlock(t,
		`set pepper "galapeno"
     if (false) {
      set pepper "chipotle"
     } else {
      set pepper "chilli"
     }`, expectValueSetTo(t, "pepper", "chilli"))

	ParseTestAndRunBlock(t,
		`set pepper "galapeno"
     if (false) {
      set pepper "chilli"
     } {
      set pepper "chipotle"
     }`, expectValueSetTo(t, "pepper", "chipotle"))

}
