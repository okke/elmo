package elmo

import (
	"fmt"
	"reflect"
	"testing"
)

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

func expectErrorValueAt(t *testing.T, lineno int) func(RunContext, Value) {

	return func(context RunContext, blockResult Value) {

		if blockResult.Type() != TypeError {
			t.Errorf("expected error but found %v", blockResult)
			return
		}

		_, l := blockResult.(ErrorValue).At()

		if l != lineno {
			fmt.Printf("%s\n", blockResult.String())
			t.Errorf("expected error at line %d but found it on line %d", lineno, l)
		}

	}
}

func expectNothing(t *testing.T) func(RunContext, Value) {

	return func(context RunContext, blockResult Value) {
		if blockResult != Nothing {
			t.Errorf("expected nothing but found %v", blockResult)
		}
	}
}

func expectValue(t *testing.T, value Value) func(RunContext, Value) {

	return func(context RunContext, blockResult Value) {
		if !reflect.DeepEqual(blockResult, value) {
			t.Errorf("expected (%v) but found (%v)", value, blockResult)
		}
	}
}

func TestSetValueIntoGlobalContext(t *testing.T) {

	ParseTestAndRunBlock(t, "set chipotle \"sauce\"", expectValueSetTo(t, "chipotle", "sauce"))
	ParseTestAndRunBlock(t, "set to_many_arguments chipotle \"sauce\"", expectErrorValueAt(t, 1))
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
     set sauce (get)`, expectErrorValueAt(t, 2))

}

func TestDynamicSetValueIntoGlobalContext(t *testing.T) {

	ParseTestAndRunBlock(t,
		`set chipotle "sauce"
     set (chipotle) 147`, expectValueSetTo(t, "sauce", "147"))
}

func TestUserDefinedFunctionWithoutArguments(t *testing.T) {

	ParseTestAndRunBlock(t,
		`set fsauce (func)`, expectErrorValueAt(t, 1))

	ParseTestAndRunBlock(t,
		`set fsauce (func {
       return "chipotle"
     })
     set sauce (fsauce)`, expectValueSetTo(t, "sauce", "chipotle"))

	ParseTestAndRunBlock(t,
		`set fsauce (func {
        return
      })
      set sauce (fsauce)`, expectErrorValueAt(t, 2))
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
		`if`, expectErrorValueAt(t, 1))

	ParseTestAndRunBlock(t,
		`if {}`, expectErrorValueAt(t, 1))

	ParseTestAndRunBlock(t,
		`if 33 {}`, expectErrorValueAt(t, 1))

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
		`if (false) {} else {} soep`, expectErrorValueAt(t, 1))

	ParseTestAndRunBlock(t,
		`if (false) {} ilse {}`, expectErrorValueAt(t, 1))

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

func TestListCreation(t *testing.T) {

	ParseTestAndRunBlock(t,
		`list 3`, expectValue(t, NewListValue([]Value{NewIntegerLiteral(3)})))

	ParseTestAndRunBlock(t,
		`list 3 4`, expectValue(t, NewListValue([]Value{NewIntegerLiteral(3), NewIntegerLiteral(4)})))

	ParseTestAndRunBlock(t,
		`list 3 "chipotle"`, expectValue(t, NewListValue([]Value{NewIntegerLiteral(3), NewStringLiteral("chipotle")})))

	ParseTestAndRunBlock(t,
		`list (list 3 "chipotle")`, expectValue(t, NewListValue([]Value{NewListValue([]Value{NewIntegerLiteral(3), NewStringLiteral("chipotle")})})))
}

func TestListAccess(t *testing.T) {

	ParseTestAndRunBlock(t,
		`set ll (list 1 2 3)
		 ll 0`, expectValue(t, NewIntegerLiteral(1)))

	ParseTestAndRunBlock(t,
		`set ll (list 1 2 3)
	 	 ll 1`, expectValue(t, NewIntegerLiteral(2)))

	ParseTestAndRunBlock(t,
		`set ll (list 1 2 3)
		 set idx 2
 	 	 ll (idx)`, expectValue(t, NewIntegerLiteral(3)))

	ParseTestAndRunBlock(t,
		`set ll (list 1 2 3)
 	 	 ll -1`, expectValue(t, NewIntegerLiteral(3)))

	ParseTestAndRunBlock(t,
		`set ll (list 1 2 3)
	 	 ll -2`, expectValue(t, NewIntegerLiteral(2)))

	// index must be integer
	//
	ParseTestAndRunBlock(t,
		`set ll (list 1 2 3)
 		 ll "chipotle"`, expectErrorValueAt(t, 2))

	// index out of bounds
	//
	ParseTestAndRunBlock(t,
		`set ll (list 1 2 3)
  	 ll 3`, expectErrorValueAt(t, 2))

	// index out of bounds
	//
	ParseTestAndRunBlock(t,
		`set ll (list 1 2 3)
   	 ll -4`, expectErrorValueAt(t, 2))
}

func TestPuts(t *testing.T) {

	ParseTestAndRunBlock(t,
		`puts 3`, expectNothing(t))

	ParseTestAndRunBlock(t,
		`puts "3 4 5"`, expectNothing(t))

	ParseTestAndRunBlock(t,
		`puts 3 4 5`, expectNothing(t))

	ParseTestAndRunBlock(t,
		`puts (list 3 4 5)`, expectNothing(t))

	ParseTestAndRunBlock(t,
		`set ll (list 3 4 5)
		 puts (ll 1)`, expectNothing(t))
}
