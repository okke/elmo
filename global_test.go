package elmo

import "testing"

func expectValueSetTo(t *testing.T, key string, value string) func(RunContext, Value) {

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

func expectErrorValue(t *testing.T) func(RunContext, Value) {

	return func(context RunContext, blockResult Value) {

		if blockResult.Type() != TypeError {
			t.Errorf("expected error but found %v", blockResult)
		}

	}
}

func TestSetValueIntoGlobalContext(t *testing.T) {

	ParseTestAndRunBlock(t, "set chipotle \"sauce\"", expectValueSetTo(t, "chipotle", "sauce"))
	ParseTestAndRunBlock(t, "set to_many_arguments chipotle \"sauce\"", expectErrorValue(t))
}

func TestSetValueIntoGlobalContextAndGetIt(t *testing.T) {

	ParseTestAndRunBlock(t,
		`set chipotle "sauce"
     set sauce (get chipotle)`, expectValueSetTo(t, "sauce", "sauce"))

	ParseTestAndRunBlock(t,
		`set chipotle "sauce"
     set sauce (chipotle)`, expectValueSetTo(t, "sauce", "sauce"))

	ParseTestAndRunBlock(t,
		`set chipotle "sauce"
     set sauce (get)`, expectErrorValue(t))

}

func TestDynamicSetValueIntoGlobalContext(t *testing.T) {

	ParseTestAndRunBlock(t,
		`set chipotle "sauce"
     set (chipotle) 147`, expectValueSetTo(t, "sauce", "147"))
}

func TestUserDefinedFunctionWithoutArguments(t *testing.T) {

	ParseTestAndRunBlock(t,
		`set fsauce (func)`, expectErrorValue(t))

	ParseTestAndRunBlock(t,
		`set fsauce (func {
       return "chipotle"
     })
     set sauce (fsauce)`, expectValueSetTo(t, "sauce", "chipotle"))

	ParseTestAndRunBlock(t,
		`set fsauce (func {
        return
      })
      set sauce (fsauce)`, expectErrorValue(t))
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
		`if`, expectErrorValue(t))

	ParseTestAndRunBlock(t,
		`if {}`, expectErrorValue(t))

	ParseTestAndRunBlock(t,
		`if 33 {}`, expectErrorValue(t))

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
		`if (false) {} else {} soep`, expectErrorValue(t))

	ParseTestAndRunBlock(t,
		`if (false) {} ilse {}`, expectErrorValue(t))

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
