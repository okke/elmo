package elmo

import "testing"

func TestSetValueIntoGlobalContext(t *testing.T) {

	ParseTestAndRunBlock(t, "set chipotle \"sauce\"", ExpectValueSetTo(t, "chipotle", "sauce"))
	ParseTestAndRunBlock(t, "set to_many_arguments chipotle \"sauce\"", ExpectErrorValueAt(t, 1))
}

func TestSetValueIntoGlobalContextUsingShortcut(t *testing.T) {

	ParseTestAndRunBlock(t, "chipotle: \"sauce\"", ExpectValueSetTo(t, "chipotle", "sauce"))

}

func TestType(t *testing.T) {

	ParseTestAndRunBlock(t, `type chipotle`, ExpectValue(t, NewIdentifier("identifier")))
	ParseTestAndRunBlock(t, `type "chipotle"`, ExpectValue(t, NewIdentifier("string")))
	ParseTestAndRunBlock(t, `type 3`, ExpectValue(t, NewIdentifier("int")))
	ParseTestAndRunBlock(t, `type 3.0`, ExpectValue(t, NewIdentifier("float")))
	ParseTestAndRunBlock(t, `type []`, ExpectValue(t, NewIdentifier("list")))
	ParseTestAndRunBlock(t, `type (dict [])`, ExpectValue(t, NewIdentifier("dict")))

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
			 return [(pepper)]
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

func TestWhile(t *testing.T) {

	ParseTestAndRunBlock(t,
		`pepper: 1
     while (lte (pepper) 10) {
       incr pepper
     }
		 pepper`, ExpectValue(t, NewIntegerLiteral(11)))
}

func TestDoWhile(t *testing.T) {

	ParseTestAndRunBlock(t,
		`pepper: 1
     do {
       incr pepper
     } while (lte (pepper) 5)
		 pepper`, ExpectValue(t, NewIntegerLiteral(6)))
}

func TestUntil(t *testing.T) {

	ParseTestAndRunBlock(t,
		`pepper: 1
     until (eq (pepper) 10) {
       incr pepper
     }
		 pepper`, ExpectValue(t, NewIntegerLiteral(10)))
}

func TestDoUntil(t *testing.T) {

	ParseTestAndRunBlock(t,
		`pepper: 1
     do {
       incr pepper
     } until (eq (pepper) 5)
		 pepper`, ExpectValue(t, NewIntegerLiteral(5)))
}

func TestListCreation(t *testing.T) {

	ParseTestAndRunBlock(t,
		`[3]`, ExpectValue(t, NewListValue([]Value{NewIntegerLiteral(3)})))

	ParseTestAndRunBlock(t,
		`[3 4]`, ExpectValue(t, NewListValue([]Value{NewIntegerLiteral(3), NewIntegerLiteral(4)})))

	ParseTestAndRunBlock(t,
		`[3 "chipotle"]`, ExpectValue(t, NewListValue([]Value{NewIntegerLiteral(3), NewStringLiteral("chipotle")})))

	ParseTestAndRunBlock(t,
		`[[3 "chipotle"]]`, ExpectValue(t, NewListValue([]Value{NewListValue([]Value{NewIntegerLiteral(3), NewStringLiteral("chipotle")})})))
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
		`set ll [1 2 3]
	 	 ll 1`, ExpectValue(t, NewIntegerLiteral(2)))

	ParseTestAndRunBlock(t,
		`set ll [1 2 3]
		 set idx 2
 	 	 ll (idx)`, ExpectValue(t, NewIntegerLiteral(3)))

	ParseTestAndRunBlock(t,
		`set ll [1 2 3]
 	 	 ll -1`, ExpectValue(t, NewIntegerLiteral(3)))

	ParseTestAndRunBlock(t,
		`set ll [1 2 3]
	 	 ll -2`, ExpectValue(t, NewIntegerLiteral(2)))

	ParseTestAndRunBlock(t,
		`set ll [1 2 3 4 5 6]
 	 	 ll 0 2`, ExpectValue(t, NewListValue([]Value{NewIntegerLiteral(1), NewIntegerLiteral(2), NewIntegerLiteral(3)})))

	ParseTestAndRunBlock(t,
		`set ll [1 2 3]
		 ll 0 -1`, ExpectValue(t, NewListValue([]Value{NewIntegerLiteral(1), NewIntegerLiteral(2), NewIntegerLiteral(3)})))

	ParseTestAndRunBlock(t,
		`set ll [1 2 3 4 5 6]
		 ll 3 -1`, ExpectValue(t, NewListValue([]Value{NewIntegerLiteral(4), NewIntegerLiteral(5), NewIntegerLiteral(6)})))

	ParseTestAndRunBlock(t,
		`set ll [1 2 3]
		 ll 0 -2`, ExpectValue(t, NewListValue([]Value{NewIntegerLiteral(1), NewIntegerLiteral(2)})))

	// when last index is smaller than first, reverse result
	//
	ParseTestAndRunBlock(t,
		`set ll [1 2 3 4 5 6]
	 	ll 3 1`, ExpectValue(t, NewListValue([]Value{NewIntegerLiteral(4), NewIntegerLiteral(3), NewIntegerLiteral(2)})))

	// when last index is smaller than first, reverse result
	//
	ParseTestAndRunBlock(t,
		`set ll [1 2 3 4 5 6]
	 	 ll -3 1`, ExpectValue(t, NewListValue([]Value{NewIntegerLiteral(4), NewIntegerLiteral(3), NewIntegerLiteral(2)})))

	// reverse list using index accessors
	//
	ParseTestAndRunBlock(t,
		`set ll [1 2 3]
 	 	 ll -1 0`, ExpectValue(t, NewListValue([]Value{NewIntegerLiteral(3), NewIntegerLiteral(2), NewIntegerLiteral(1)})))

	// index must be integer
	//
	ParseTestAndRunBlock(t,
		`set ll [1 2 3]
 		 ll "chipotle"`, ExpectErrorValueAt(t, 2))

	// index out of bounds
	//
	ParseTestAndRunBlock(t,
		`set ll [1 2 3]
  	 ll 3`, ExpectErrorValueAt(t, 2))

	// index out of bounds
	//
	ParseTestAndRunBlock(t,
		`set ll [1 2 3]
   	 ll -4`, ExpectErrorValueAt(t, 2))
}

func TestDictionary(t *testing.T) {

	ParseTestAndRunBlock(t,
		`set d (dict 1 2 3 4)
		 d 1`, ExpectValue(t, NewIntegerLiteral(2)))

	ParseTestAndRunBlock(t,
		`set d (dict [1 2 3 4])
 		 d 1`, ExpectValue(t, NewIntegerLiteral(2)))

	ParseTestAndRunBlock(t,
		`set d (dict 1 2 3 4)
 		 d 3`, ExpectValue(t, NewIntegerLiteral(4)))

	ParseTestAndRunBlock(t,
		`set d (dict [1 2 3 4])
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
		`set l [a 2 b 4]
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
		`io: [1 2 3]
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
	ParseTestAndRunBlock(t, `eq "1" [1]`, ExpectValue(t, False))
	ParseTestAndRunBlock(t, `eq [1] [1]`, ExpectValue(t, True))
	ParseTestAndRunBlock(t, `eq [1] [0]`, ExpectValue(t, False))
	ParseTestAndRunBlock(t, `eq`, ExpectErrorValueAt(t, 1))
	ParseTestAndRunBlock(t, `eq 1`, ExpectErrorValueAt(t, 1))
	ParseTestAndRunBlock(t, `eq 1 1 1`, ExpectErrorValueAt(t, 1))
}

func TestNe(t *testing.T) {
	ParseTestAndRunBlock(t, `ne 1 1`, ExpectValue(t, False))
	ParseTestAndRunBlock(t, `ne 1 0`, ExpectValue(t, True))
	ParseTestAndRunBlock(t, `ne 1 "1"`, ExpectValue(t, True))
	ParseTestAndRunBlock(t, `ne "1" "1"`, ExpectValue(t, False))
	ParseTestAndRunBlock(t, `ne "1" [1]`, ExpectValue(t, True))
	ParseTestAndRunBlock(t, `ne [1] [1]`, ExpectValue(t, False))
	ParseTestAndRunBlock(t, `ne [1] [0]`, ExpectValue(t, True))
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

	ParseTestAndRunBlock(t, `gt 1 1.0`, ExpectErrorValueAt(t, 1))
	ParseTestAndRunBlock(t, `gt 1.0 1`, ExpectErrorValueAt(t, 1))
	ParseTestAndRunBlock(t, `gt "chipotle" "galapeno"`, ExpectErrorValueAt(t, 1))

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

func TestPlus(t *testing.T) {
	ParseTestAndRunBlock(t, `plus 1 1`, ExpectValue(t, NewIntegerLiteral(2)))
	ParseTestAndRunBlock(t, `plus 1 -1`, ExpectValue(t, NewIntegerLiteral(0)))
	ParseTestAndRunBlock(t, `plus -1 1`, ExpectValue(t, NewIntegerLiteral(0)))
	ParseTestAndRunBlock(t, `plus 1 1.0`, ExpectValue(t, NewFloatLiteral(2.0)))
	ParseTestAndRunBlock(t, `plus 1.0 1`, ExpectValue(t, NewFloatLiteral(2.0)))
	ParseTestAndRunBlock(t, `plus 88 99`, ExpectValue(t, NewIntegerLiteral(187)))
	ParseTestAndRunBlock(t, `plus 88 "chipotle"`, ExpectErrorValueAt(t, 1))
	ParseTestAndRunBlock(t, `plus 1.0 "chipotle"`, ExpectErrorValueAt(t, 1))
	ParseTestAndRunBlock(t, `plus "galapeno" 88`, ExpectErrorValueAt(t, 1))
	ParseTestAndRunBlock(t, `plus "galapeno" 1.0`, ExpectErrorValueAt(t, 1))
}

func TestMinus(t *testing.T) {
	ParseTestAndRunBlock(t, `minus 1 1`, ExpectValue(t, NewIntegerLiteral(0)))
	ParseTestAndRunBlock(t, `minus 1 -1`, ExpectValue(t, NewIntegerLiteral(2)))
	ParseTestAndRunBlock(t, `minus -1 1`, ExpectValue(t, NewIntegerLiteral(-2)))
	ParseTestAndRunBlock(t, `minus 1 1.0`, ExpectValue(t, NewFloatLiteral(0.0)))
	ParseTestAndRunBlock(t, `minus 1.0 1`, ExpectValue(t, NewFloatLiteral(0.0)))
	ParseTestAndRunBlock(t, `minus 88 99`, ExpectValue(t, NewIntegerLiteral(-11)))
	ParseTestAndRunBlock(t, `minus 88 "chipotle"`, ExpectErrorValueAt(t, 1))
	ParseTestAndRunBlock(t, `minus 1.0 "chipotle"`, ExpectErrorValueAt(t, 1))
	ParseTestAndRunBlock(t, `minus "galapeno" 88`, ExpectErrorValueAt(t, 1))
	ParseTestAndRunBlock(t, `minus "galapeno" 1.0`, ExpectErrorValueAt(t, 1))
}

func TestMultiply(t *testing.T) {
	ParseTestAndRunBlock(t, `multiply 3 4`, ExpectValue(t, NewIntegerLiteral(12)))
	ParseTestAndRunBlock(t, `multiply 3 -4`, ExpectValue(t, NewIntegerLiteral(-12)))
	ParseTestAndRunBlock(t, `multiply -1 2`, ExpectValue(t, NewIntegerLiteral(-2)))
	ParseTestAndRunBlock(t, `multiply 3 2.0`, ExpectValue(t, NewFloatLiteral(6.0)))
	ParseTestAndRunBlock(t, `multiply 2.0 3`, ExpectValue(t, NewFloatLiteral(6.0)))
	ParseTestAndRunBlock(t, `multiply 88 99`, ExpectValue(t, NewIntegerLiteral(8712)))
	ParseTestAndRunBlock(t, `multiply 88 "chipotle"`, ExpectErrorValueAt(t, 1))
	ParseTestAndRunBlock(t, `multiply 1.0 "chipotle"`, ExpectErrorValueAt(t, 1))
	ParseTestAndRunBlock(t, `multiply "galapeno" 88`, ExpectErrorValueAt(t, 1))
	ParseTestAndRunBlock(t, `multiply "galapeno" 1.0`, ExpectErrorValueAt(t, 1))
}

func TestDivide(t *testing.T) {
	ParseTestAndRunBlock(t, `divide 8 3`, ExpectValue(t, NewIntegerLiteral(2)))
	ParseTestAndRunBlock(t, `divide 8 -3`, ExpectValue(t, NewIntegerLiteral(-2)))
	ParseTestAndRunBlock(t, `divide -8 2`, ExpectValue(t, NewIntegerLiteral(-4)))
	ParseTestAndRunBlock(t, `divide 3 2.0`, ExpectValue(t, NewFloatLiteral(3/2.0)))
	ParseTestAndRunBlock(t, `divide 2.0 3`, ExpectValue(t, NewFloatLiteral(2.0/3)))
	ParseTestAndRunBlock(t, `divide 88 99`, ExpectValue(t, NewIntegerLiteral(88/99)))
	ParseTestAndRunBlock(t, `divide 88 0`, ExpectErrorValueAt(t, 1))
	ParseTestAndRunBlock(t, `divide 88 0.0`, ExpectErrorValueAt(t, 1))
	ParseTestAndRunBlock(t, `divide 88.0 0`, ExpectErrorValueAt(t, 1))
	ParseTestAndRunBlock(t, `divide 88.0 0.0`, ExpectErrorValueAt(t, 1))
	ParseTestAndRunBlock(t, `divide 88 "chipotle"`, ExpectErrorValueAt(t, 1))
	ParseTestAndRunBlock(t, `divide 1.0 "chipotle"`, ExpectErrorValueAt(t, 1))
	ParseTestAndRunBlock(t, `divide "galapeno" 88`, ExpectErrorValueAt(t, 1))
	ParseTestAndRunBlock(t, `divide "galapeno" 1.0`, ExpectErrorValueAt(t, 1))
}

func TestModulo(t *testing.T) {
	ParseTestAndRunBlock(t, `modulo 10 3`, ExpectValue(t, NewIntegerLiteral(1)))
	ParseTestAndRunBlock(t, `modulo 10.0 3`, ExpectValue(t, NewFloatLiteral(1.0)))
	ParseTestAndRunBlock(t, `modulo 11.5 3`, ExpectValue(t, NewFloatLiteral(2.5)))
	ParseTestAndRunBlock(t, `modulo 11.5 3.5`, ExpectValue(t, NewFloatLiteral(1.0)))
	ParseTestAndRunBlock(t, `modulo 11.5 0`, ExpectErrorValueAt(t, 1))
	ParseTestAndRunBlock(t, `modulo 11.5 0.0`, ExpectErrorValueAt(t, 1))
	ParseTestAndRunBlock(t, `modulo 11 0`, ExpectErrorValueAt(t, 1))
	ParseTestAndRunBlock(t, `modulo 11 0.0`, ExpectErrorValueAt(t, 1))
	ParseTestAndRunBlock(t, `modulo 88 "chipotle"`, ExpectErrorValueAt(t, 1))
	ParseTestAndRunBlock(t, `modulo 1.0 "chipotle"`, ExpectErrorValueAt(t, 1))
	ParseTestAndRunBlock(t, `modulo "galapeno" 88`, ExpectErrorValueAt(t, 1))
	ParseTestAndRunBlock(t, `modulo "galapeno" 1.0`, ExpectErrorValueAt(t, 1))
}
