package elmo

import "testing"

func TestSetValueIntoGlobalContext(t *testing.T) {

	ParseTestAndRunBlock(t, "set chipotle \"sauce\"", ExpectValueSetTo(t, "chipotle", "sauce"))
	ParseTestAndRunBlock(t, "set to_many_arguments chipotle \"sauce\"", ExpectErrorValueAt(t, 1))
}

func TestSetValueIntoGlobalContextUsingShortcut(t *testing.T) {

	ParseTestAndRunBlock(t, "chipotle: \"sauce\"", ExpectValueSetTo(t, "chipotle", "sauce"))

}

func TestSetValueIntoGlobalContextAndGetIt(t *testing.T) {

	ParseTestAndRunBlock(t,
		`set chipotle "sauce"
     set sauce (get chipotle)`, ExpectValueSetTo(t, "sauce", "sauce"))

	ParseTestAndRunBlock(t,
		`set chipotle "sauce"
     set sauce (chipotle)`, ExpectValueSetTo(t, "sauce", "sauce"))

	ParseTestAndRunBlock(t,
		`set chipotle "sauce"
     set sauce (get)`, ExpectErrorValueAt(t, 2))

}

func TestDynamicSetValueIntoGlobalContext(t *testing.T) {

	ParseTestAndRunBlock(t,
		`set chipotle "sauce"
     set (chipotle) 147`, ExpectValueSetTo(t, "sauce", "147"))
}

func TestUserDefinedFunctionWithoutArguments(t *testing.T) {

	ParseTestAndRunBlock(t,
		`set fsauce (func)`, ExpectErrorValueAt(t, 1))

	ParseTestAndRunBlock(t,
		`set fsauce (func {
       return "chipotle"
     })
     set sauce (fsauce)`, ExpectValueSetTo(t, "sauce", "chipotle"))

	ParseTestAndRunBlock(t,
		`set fsauce (func {
        return
      })
      set sauce (fsauce)`, ExpectErrorValueAt(t, 2))
}

func TestUserDefinedFunctionWithOneArgument(t *testing.T) {

	ParseTestAndRunBlock(t,
		`set fsauce (func pepper {
       return (pepper)
     })
     set sauce (fsauce "chipotle")`, ExpectValueSetTo(t, "sauce", "chipotle"))
}

func TestIfWithoutElse(t *testing.T) {

	ParseTestAndRunBlock(t,
		`if`, ExpectErrorValueAt(t, 1))

	ParseTestAndRunBlock(t,
		`if {}`, ExpectErrorValueAt(t, 1))

	ParseTestAndRunBlock(t,
		`if 33 {}`, ExpectErrorValueAt(t, 1))

	ParseTestAndRunBlock(t,
		`set pepper "galapeno"
     if (true) {
      set pepper "chipotle"
     }`, ExpectValueSetTo(t, "pepper", "chipotle"))

	ParseTestAndRunBlock(t,
		`set pepper "galapeno"
     if (false) {
      set pepper "chipotle"
     }`, ExpectValueSetTo(t, "pepper", "galapeno"))
}

func TestIfWithElse(t *testing.T) {

	ParseTestAndRunBlock(t,
		`if (false) {} else {} soep`, ExpectErrorValueAt(t, 1))

	ParseTestAndRunBlock(t,
		`if (false) {} ilse {}`, ExpectErrorValueAt(t, 1))

	ParseTestAndRunBlock(t,
		`set pepper "galapeno"
     if (true) {
      set pepper "chipotle"
     } else {
      set pepper "chilli"
     }`, ExpectValueSetTo(t, "pepper", "chipotle"))

	ParseTestAndRunBlock(t,
		`set pepper "galapeno"
     if (false) {
      set pepper "chipotle"
     } else {
      set pepper "chilli"
     }`, ExpectValueSetTo(t, "pepper", "chilli"))

	ParseTestAndRunBlock(t,
		`set pepper "galapeno"
     if (false) {
      set pepper "chilli"
     } {
      set pepper "chipotle"
     }`, ExpectValueSetTo(t, "pepper", "chipotle"))

}

func TestListCreation(t *testing.T) {

	ParseTestAndRunBlock(t,
		`list 3`, ExpectValue(t, NewListValue([]Value{NewIntegerLiteral(3)})))

	ParseTestAndRunBlock(t,
		`list 3 4`, ExpectValue(t, NewListValue([]Value{NewIntegerLiteral(3), NewIntegerLiteral(4)})))

	ParseTestAndRunBlock(t,
		`list 3 "chipotle"`, ExpectValue(t, NewListValue([]Value{NewIntegerLiteral(3), NewStringLiteral("chipotle")})))

	ParseTestAndRunBlock(t,
		`list (list 3 "chipotle")`, ExpectValue(t, NewListValue([]Value{NewListValue([]Value{NewIntegerLiteral(3), NewStringLiteral("chipotle")})})))
}

func TestListAccess(t *testing.T) {

	ParseTestAndRunBlock(t,
		`set ll (list 1 2 3)
		 ll 0`, ExpectValue(t, NewIntegerLiteral(1)))

	ParseTestAndRunBlock(t,
		`set ll (list 1 2 3)
	 	 ll 1`, ExpectValue(t, NewIntegerLiteral(2)))

	ParseTestAndRunBlock(t,
		`set ll (list 1 2 3)
		 set idx 2
 	 	 ll (idx)`, ExpectValue(t, NewIntegerLiteral(3)))

	ParseTestAndRunBlock(t,
		`set ll (list 1 2 3)
 	 	 ll -1`, ExpectValue(t, NewIntegerLiteral(3)))

	ParseTestAndRunBlock(t,
		`set ll (list 1 2 3)
	 	 ll -2`, ExpectValue(t, NewIntegerLiteral(2)))

	ParseTestAndRunBlock(t,
		`set ll (list 1 2 3 4 5 6)
 	 	 ll 0 2`, ExpectValue(t, NewListValue([]Value{NewIntegerLiteral(1), NewIntegerLiteral(2), NewIntegerLiteral(3)})))

	ParseTestAndRunBlock(t,
		`set ll (list 1 2 3)
		 ll 0 -1`, ExpectValue(t, NewListValue([]Value{NewIntegerLiteral(1), NewIntegerLiteral(2), NewIntegerLiteral(3)})))

	ParseTestAndRunBlock(t,
		`set ll (list 1 2 3 4 5 6)
		 ll 3 -1`, ExpectValue(t, NewListValue([]Value{NewIntegerLiteral(4), NewIntegerLiteral(5), NewIntegerLiteral(6)})))

	ParseTestAndRunBlock(t,
		`set ll (list 1 2 3)
		 ll 0 -2`, ExpectValue(t, NewListValue([]Value{NewIntegerLiteral(1), NewIntegerLiteral(2)})))
	// index must be integer
	//
	ParseTestAndRunBlock(t,
		`set ll (list 1 2 3)
 		 ll "chipotle"`, ExpectErrorValueAt(t, 2))

	// index out of bounds
	//
	ParseTestAndRunBlock(t,
		`set ll (list 1 2 3)
  	 ll 3`, ExpectErrorValueAt(t, 2))

	// index out of bounds
	//
	ParseTestAndRunBlock(t,
		`set ll (list 1 2 3)
   	 ll -4`, ExpectErrorValueAt(t, 2))
}

func TestDictionary(t *testing.T) {

	ParseTestAndRunBlock(t,
		`set d (dict 1 2 3 4)
		 d 1`, ExpectValue(t, NewIntegerLiteral(2)))

	ParseTestAndRunBlock(t,
		`set d (dict (list 1 2 3 4))
 		 d 1`, ExpectValue(t, NewIntegerLiteral(2)))

	ParseTestAndRunBlock(t,
		`set d (dict 1 2 3 4)
 		 d 3`, ExpectValue(t, NewIntegerLiteral(4)))

	ParseTestAndRunBlock(t,
		`set d (dict (list 1 2 3 4))
  		 d 3`, ExpectValue(t, NewIntegerLiteral(4)))

	ParseTestAndRunBlock(t,
		`set d (dict 1 2 3 4)
  	 d 5`, ExpectNothing(t))

	ParseTestAndRunBlock(t,
		`set d (dict "1" 2 3 4)
	   d 1`, ExpectValue(t, NewIntegerLiteral(2)))

	ParseTestAndRunBlock(t,
		`set d (dict a 2 b 4)
	 	 d b`, ExpectValue(t, NewIntegerLiteral(4)))

	ParseTestAndRunBlock(t,
		`set l (list a 2 b 4)
		 set d (dict (l))
 	 	 d b`, ExpectValue(t, NewIntegerLiteral(4)))
}

func TestDictionaryWithBlock(t *testing.T) {

	ParseTestAndRunBlock(t,
		`set d (dict {
			set b 4
		 })
 	 	 d b`, ExpectValue(t, NewIntegerLiteral(4)))

	ParseTestAndRunBlock(t,
		`set val "chipotle"
		 set d (dict {
 			set b (val)
 		 })
  	 d b`, ExpectValue(t, NewStringLiteral("chipotle")))

	ParseTestAndRunBlock(t,
		`set val "chipotle"
 		 set d (dict {
  		set val "galapeno"
			set b (val)
  	 })
   	 d b`, ExpectValue(t, NewStringLiteral("galapeno")))
}

func TestMixin(t *testing.T) {
	ParseTestAndRunBlock(t,
		`mixin (dict a 2 b 4)
	 	 b`, ExpectValue(t, NewIntegerLiteral(4)))

	ParseTestAndRunBlock(t,
		`mixin (dict {
			set b "chipotle"
		 })
 	 	 b`, ExpectValue(t, NewStringLiteral("chipotle")))

	ParseTestAndRunBlock(t,
		`set hot_or_not (dict {
			set chipotle (false)
			set galapeno (true)
		 })
		 mixin (dict {
 			mixin (hot_or_not)
 		 })
  	 chipotle`, ExpectValue(t, NewBooleanLiteral(false)))
}

func TestDictionaryAccessShortcut(t *testing.T) {
	ParseTestAndRunBlock(t,
		`io: (dict {
		  read: (func {
		  	return "chipotle"
		  })
		})
		io.read`, ExpectValue(t, NewStringLiteral("chipotle")))

	ParseTestAndRunBlock(t,
		`io: (dict {
			read: (func name {
				return (name)
			})
		})
		io.read "chipotle"`, ExpectValue(t, NewStringLiteral("chipotle")))

	// Access shortcut can be on dictionaries only
	//
	ParseTestAndRunBlock(t,
		`io: (list 1 2 3)
		io.read`, ExpectErrorValueAt(t, 2))
}

func TestLoad(t *testing.T) {

	context := NewGlobalContext()

	context.RegisterModule(NewModule("yippie", func(context RunContext) Value {

		mapping := make(map[string]Value)

		mapping["nop"] = NewGoFunction("nop", func(context RunContext, arguments []Argument) Value {
			return Nothing
		})

		return NewDictionaryValue(mapping)

	}))

	ParseTestAndRunBlockWithinContext(t, context,
		`yy: (load "yippie")
		 yy.nop`, ExpectNothing(t))
}

func TestEq(t *testing.T) {
	ParseTestAndRunBlock(t, `eq 1 1`, ExpectValue(t, True))
	ParseTestAndRunBlock(t, `eq 1 0`, ExpectValue(t, False))
	ParseTestAndRunBlock(t, `eq 1 "1"`, ExpectValue(t, False))
	ParseTestAndRunBlock(t, `eq "1" "1"`, ExpectValue(t, True))
	ParseTestAndRunBlock(t, `eq "1" (list 1)`, ExpectValue(t, False))
	ParseTestAndRunBlock(t, `eq (list 1) (list 1)`, ExpectValue(t, True))
	ParseTestAndRunBlock(t, `eq (list 1) (list 0)`, ExpectValue(t, False))
	ParseTestAndRunBlock(t, `eq`, ExpectErrorValueAt(t, 1))
	ParseTestAndRunBlock(t, `eq 1`, ExpectErrorValueAt(t, 1))
	ParseTestAndRunBlock(t, `eq 1 1 1`, ExpectErrorValueAt(t, 1))
}

func TestNe(t *testing.T) {
	ParseTestAndRunBlock(t, `ne 1 1`, ExpectValue(t, False))
	ParseTestAndRunBlock(t, `ne 1 0`, ExpectValue(t, True))
	ParseTestAndRunBlock(t, `ne 1 "1"`, ExpectValue(t, True))
	ParseTestAndRunBlock(t, `ne "1" "1"`, ExpectValue(t, False))
	ParseTestAndRunBlock(t, `ne "1" (list 1)`, ExpectValue(t, True))
	ParseTestAndRunBlock(t, `ne (list 1) (list 1)`, ExpectValue(t, False))
	ParseTestAndRunBlock(t, `ne (list 1) (list 0)`, ExpectValue(t, True))
	ParseTestAndRunBlock(t, `ne`, ExpectErrorValueAt(t, 1))
	ParseTestAndRunBlock(t, `ne 1`, ExpectErrorValueAt(t, 1))
	ParseTestAndRunBlock(t, `ne 1 1 1`, ExpectErrorValueAt(t, 1))
}

func TestAnd(t *testing.T) {
	ParseTestAndRunBlock(t, `and (true) (true)`, ExpectValue(t, True))
	ParseTestAndRunBlock(t, `and (true) (true) (true)`, ExpectValue(t, True))
	ParseTestAndRunBlock(t, `and (false) (false)`, ExpectValue(t, False))
	ParseTestAndRunBlock(t, `and (false) (true)`, ExpectValue(t, False))
	ParseTestAndRunBlock(t, `and (true) (false)`, ExpectValue(t, False))
	ParseTestAndRunBlock(t, `and (true)`, ExpectValue(t, True))
	ParseTestAndRunBlock(t, `and (false)`, ExpectValue(t, False))
}

func TestOr(t *testing.T) {
	ParseTestAndRunBlock(t, `or (true) (true)`, ExpectValue(t, True))
	ParseTestAndRunBlock(t, `or (true) (true) (true)`, ExpectValue(t, True))
	ParseTestAndRunBlock(t, `or (false) (false)`, ExpectValue(t, False))
	ParseTestAndRunBlock(t, `or (false) (false) (false)`, ExpectValue(t, False))
	ParseTestAndRunBlock(t, `or (false) (false) (true)`, ExpectValue(t, True))
	ParseTestAndRunBlock(t, `or (false) (true)`, ExpectValue(t, True))
	ParseTestAndRunBlock(t, `or (true) (false)`, ExpectValue(t, True))
	ParseTestAndRunBlock(t, `or (true)`, ExpectValue(t, True))
	ParseTestAndRunBlock(t, `or (false)`, ExpectValue(t, False))
}

func TestNot(t *testing.T) {
	ParseTestAndRunBlock(t, `not (true)`, ExpectValue(t, False))
	ParseTestAndRunBlock(t, `not (false)`, ExpectValue(t, True))
	ParseTestAndRunBlock(t, `not (true) (true)`, ExpectErrorValueAt(t, 1))
	ParseTestAndRunBlock(t, `not true`, ExpectErrorValueAt(t, 1))
	ParseTestAndRunBlock(t, `not 1`, ExpectErrorValueAt(t, 1))
}

/*
func TestPuts(t *testing.T) {

	ParseTestAndRunBlock(t,
		`puts 3`, ExpectNothing(t))

	ParseTestAndRunBlock(t,
		`puts "3 4 5"`, ExpectNothing(t))

	ParseTestAndRunBlock(t,
		`puts 3 4 5`, ExpectNothing(t))

	ParseTestAndRunBlock(t,
		`puts (list 3 4 5)`, ExpectNothing(t))

	ParseTestAndRunBlock(t,
		`set ll (list 3 4 5)
		 puts (ll 1)`, ExpectNothing(t))
}
*/
