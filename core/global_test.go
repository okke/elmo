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

func TestIncrementValue(t *testing.T) {

	ParseTestAndRunBlock(t,
		`incr chipotle
		 chipotle`, ExpectValue(t, NewIntegerLiteral(1)))

	ParseTestAndRunBlock(t,
		`incr chipotle
		 incr chipotle
		 chipotle`, ExpectValue(t, NewIntegerLiteral(2)))

	ParseTestAndRunBlock(t,
		`set chipotle 3
	 	 set galapeno (incr chipotle)
	 	 galapeno`, ExpectValue(t, NewIntegerLiteral(4)))

	ParseTestAndRunBlock(t,
		`set chipotle 3.0
 	 	 set galapeno (incr chipotle)
 	 	 galapeno`, ExpectValue(t, NewFloatLiteral(4.0)))

	// increments returns incremented value but also changes
	// incremented variable
	//
	ParseTestAndRunBlock(t,
		`set chipotle 3
	 	 set galapeno (incr chipotle)
	 	 chipotle`, ExpectValue(t, NewIntegerLiteral(4)))

	ParseTestAndRunBlock(t,
		`incr chipotle 3
 		 chipotle`, ExpectValue(t, NewIntegerLiteral(3)))

	ParseTestAndRunBlock(t,
		`incr chipotle 3
 		 incr chipotle 5
 		 chipotle`, ExpectValue(t, NewIntegerLiteral(8)))

	ParseTestAndRunBlock(t,
		`incr chipotle "galapeno"
  	 chipotle`, ExpectErrorValueAt(t, 1))

	ParseTestAndRunBlock(t,
		`incr chipotle 3
  	 incr chipotle "galapeno"
  	 chipotle`, ExpectErrorValueAt(t, 2))

	ParseTestAndRunBlock(t,
		`set chipotle (incr 3)
	 	 set chipotle (incr 3)
	 	 chipotle`, ExpectValue(t, NewIntegerLiteral(4)))
}

func TestDynamicSetValueIntoGlobalContext(t *testing.T) {

	ParseTestAndRunBlock(t,
		`set chipotle "sauce"
     set (chipotle) 147`, ExpectValueSetTo(t, "sauce", "147"))
}

func TestSetValueOnceIntoGlobalContext(t *testing.T) {

	ParseTestAndRunBlock(t,
		`once chipotle "sauce"
		 once chipotle "jar"`, ExpectValueSetTo(t, "chipotle", "sauce"))
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
      set sauce (fsauce)`, ExpectNothing(t))

	ParseTestAndRunBlock(t,
		`set fsauce (func {
			 return "chipotle"
			 return "galapeno"
		 })
		 set sauce (fsauce)`, ExpectValueSetTo(t, "sauce", "chipotle"))

}

func TestUserDefinedFunctionWithOneArgument(t *testing.T) {

	ParseTestAndRunBlock(t,
		`set fsauce (func pepper {
       return (pepper)
     })
     set sauce (fsauce "chipotle")`, ExpectValueSetTo(t, "sauce", "chipotle"))
}

func TestUserDefinedFunctionWithMultipleReturnValues(t *testing.T) {

	ParseTestAndRunBlock(t,
		`set fsauce (func  {
       return "chipotle" "galapeno"
     })
     set sauce (fsauce)`, ExpectValueSetTo(t, "sauce", "chipotle"))

	ParseTestAndRunBlock(t,
		`set fsauce (func  {
	 		 return "chipotle" "galapeno"
	 	 })
	 	 set hot hotter (fsauce)`,
		ExpectValueSetTo(t, "hot", "chipotle"),
		ExpectValueSetTo(t, "hotter", "galapeno"))

	ParseTestAndRunBlock(t,
		`set fsauce (func  {
	 		 return "chipotle" "galapeno"
	 	 })
	 	 set also_hot also_hotter (set hot hotter (fsauce))`,
		ExpectValueSetTo(t, "hot", "chipotle"),
		ExpectValueSetTo(t, "hotter", "galapeno"),
		ExpectValueSetTo(t, "also_hot", "chipotle"),
		ExpectValueSetTo(t, "also_hotter", "galapeno"))
}

func TestPipeToUserDefinedFunction(t *testing.T) {

	ParseTestAndRunBlock(t,
		`set fsauce (func pepper {
       return (pepper)
     })
		 set injar (func pepper {
			 return (list (pepper))
		 })
     fsauce "chipotle" | injar`, ExpectValue(t, NewListValue([]Value{NewStringLiteral("chipotle")})))
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

	ParseTestAndRunBlock(t,
		`set fsauce (func test {
			 if (test) {
			   return "chipotle"
			 }
	 		 return "galapeno"
	 	 })
	 	 set sauce (fsauce (true))`, ExpectValueSetTo(t, "sauce", "chipotle"))

	ParseTestAndRunBlock(t,
		`set fsauce (func test {
	 		 if (test) {
	 			 return "chipotle"
	 		 }
	 		 return "galapeno"
	 	 })
	 	 set sauce (fsauce (false))`, ExpectValueSetTo(t, "sauce", "galapeno"))

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
		`set ll [
		  1
		  2
		  3
		 ]
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

	// when last index is smaller than first, reverse result
	//
	ParseTestAndRunBlock(t,
		`set ll (list 1 2 3 4 5 6)
	 	ll 3 1`, ExpectValue(t, NewListValue([]Value{NewIntegerLiteral(4), NewIntegerLiteral(3), NewIntegerLiteral(2)})))

	// when last index is smaller than first, reverse result
	//
	ParseTestAndRunBlock(t,
		`set ll (list 1 2 3 4 5 6)
	 	 ll -3 1`, ExpectValue(t, NewListValue([]Value{NewIntegerLiteral(4), NewIntegerLiteral(3), NewIntegerLiteral(2)})))

	// reverse list using index accessors
	//
	ParseTestAndRunBlock(t,
		`set ll (list 1 2 3)
 	 	 ll -1 0`, ExpectValue(t, NewListValue([]Value{NewIntegerLiteral(3), NewIntegerLiteral(2), NewIntegerLiteral(1)})))

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

func TestDictionaryInListAsBlock(t *testing.T) {

	ParseTestAndRunBlock(t,
		`set l [
  		{
     		a:3
	  	}
    	{
      	a:4
     	}
     ]
		 first: (l 0)
		 first a`, ExpectValue(t, NewIntegerLiteral(3)))
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

func TestSetWithDictionaryAsBlock(t *testing.T) {
	ParseTestAndRunBlock(t,
		`set d {
		   set b 4
	   }
	   d b`, ExpectValue(t, NewIntegerLiteral(4)))
}

func TestDictionaryFunctionsKnowDictionary(t *testing.T) {
	ParseTestAndRunBlock(t,
		`sauce: (dict {
		  hot: (func {
		  	return "chipotle"
		  })
			same: (func {
				return (this.hot)
			})
		})
		sauce.same`, ExpectValue(t, NewStringLiteral("chipotle")))

	ParseTestAndRunBlock(t,
		`sauce: (dict {
		  hot: (func {
		  	return "chipotle"
		  })
			same: (func {
				return (this.hot)
			})
		})
		sauce.same
		this`, ExpectErrorValueAt(t, 10))

	ParseTestAndRunBlock(t,
		` soup: (func {
				return (this.hot) # will fail, this is not defined
		  })
		  sauce: (dict {
			hot: (func {
				return "chipotle"
			})
			same: (func {
				return (soup)
			})
		})
		sauce.same`, ExpectErrorValueAt(t, 2))

}

func TestNewConstructsDictionary(t *testing.T) {
	ParseTestAndRunBlock(t,
		`peppers: {
		  hot: (func {
		  	return "chipotle"
		  })
		}
		sauce: (new (peppers))
		sauce.hot`, ExpectValue(t, NewStringLiteral("chipotle")))

	ParseTestAndRunBlock(t,
		`peppers: {
		  hot: (func {
		  	return "chipotle"
		  })
		}
		sauce: (new (peppers) {
		  same: (func {
			  return (this.hot)
		  })
	  })
		sauce.same`, ExpectValue(t, NewStringLiteral("chipotle")))

	ParseTestAndRunBlock(t,
		`peppers: {
		  hot: (func {
		  	return "chipotle"
		  })
		}
		sauce: (new (peppers) {
		  hot: (func {
			  return "galapeno"
		  })
	  })
		sauce.hot`, ExpectValue(t, NewStringLiteral("galapeno")))
}

func TestLoad(t *testing.T) {

	context := NewGlobalContext()

	context.RegisterModule(NewModule("yippie", func(context RunContext) Value {

		return NewMappingForModule(context, []NamedValue{NewGoFunction("nop", func(context RunContext, arguments []Argument) Value {
			return Nothing
		})})

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

func TestGt(t *testing.T) {
	ParseTestAndRunBlock(t, `gt 1 1`, ExpectValue(t, False))
	ParseTestAndRunBlock(t, `gt 2 1`, ExpectValue(t, True))
	ParseTestAndRunBlock(t, `gt 0 1`, ExpectValue(t, False))
	ParseTestAndRunBlock(t, `gt -1 0`, ExpectValue(t, False))
	ParseTestAndRunBlock(t, `gt -1 -2`, ExpectValue(t, True))

	ParseTestAndRunBlock(t, `gt 1.0 1.0`, ExpectValue(t, False))
	ParseTestAndRunBlock(t, `gt 2.0 1.0`, ExpectValue(t, True))
	ParseTestAndRunBlock(t, `gt 0.0 1.0`, ExpectValue(t, False))
	ParseTestAndRunBlock(t, `gt -1.0 0.0`, ExpectValue(t, False))
	ParseTestAndRunBlock(t, `gt -1.0 -2.0`, ExpectValue(t, True))
}

func TestGte(t *testing.T) {
	ParseTestAndRunBlock(t, `gte 1 1`, ExpectValue(t, True))
	ParseTestAndRunBlock(t, `gte 2 1`, ExpectValue(t, True))
	ParseTestAndRunBlock(t, `gte 0 1`, ExpectValue(t, False))
	ParseTestAndRunBlock(t, `gte -1 0`, ExpectValue(t, False))
	ParseTestAndRunBlock(t, `gte -1 -2`, ExpectValue(t, True))

	ParseTestAndRunBlock(t, `gte 1.0 1.0`, ExpectValue(t, True))
	ParseTestAndRunBlock(t, `gte 2.0 1.0`, ExpectValue(t, True))
	ParseTestAndRunBlock(t, `gte 0.0 1.0`, ExpectValue(t, False))
	ParseTestAndRunBlock(t, `gte -1.0 0.0`, ExpectValue(t, False))
	ParseTestAndRunBlock(t, `gte -1.0 -2.0`, ExpectValue(t, True))
}

func TestLt(t *testing.T) {
	ParseTestAndRunBlock(t, `lt 1 1`, ExpectValue(t, False))
	ParseTestAndRunBlock(t, `lt 2 1`, ExpectValue(t, False))
	ParseTestAndRunBlock(t, `lt 0 1`, ExpectValue(t, True))
	ParseTestAndRunBlock(t, `lt -1 0`, ExpectValue(t, True))
	ParseTestAndRunBlock(t, `lt -1 -2`, ExpectValue(t, False))

	ParseTestAndRunBlock(t, `lt 1.0 1.0`, ExpectValue(t, False))
	ParseTestAndRunBlock(t, `lt 2.0 1.0`, ExpectValue(t, False))
	ParseTestAndRunBlock(t, `lt 0.0 1.0`, ExpectValue(t, True))
	ParseTestAndRunBlock(t, `lt -1.0 0.0`, ExpectValue(t, True))
	ParseTestAndRunBlock(t, `lt -1.0 -2.0`, ExpectValue(t, False))
}

func TestLte(t *testing.T) {
	ParseTestAndRunBlock(t, `lte 1 1`, ExpectValue(t, True))
	ParseTestAndRunBlock(t, `lte 2 1`, ExpectValue(t, False))
	ParseTestAndRunBlock(t, `lte 0 1`, ExpectValue(t, True))
	ParseTestAndRunBlock(t, `lte -1 0`, ExpectValue(t, True))
	ParseTestAndRunBlock(t, `lte -1 -2`, ExpectValue(t, False))

	ParseTestAndRunBlock(t, `lte 1.0 1.0`, ExpectValue(t, True))
	ParseTestAndRunBlock(t, `lte 2.0 1.0`, ExpectValue(t, False))
	ParseTestAndRunBlock(t, `lte 0.0 1.0`, ExpectValue(t, True))
	ParseTestAndRunBlock(t, `lte -1.0 0.0`, ExpectValue(t, True))
	ParseTestAndRunBlock(t, `lte -1.0 -2.0`, ExpectValue(t, False))
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
