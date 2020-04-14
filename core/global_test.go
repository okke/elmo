package elmo

import "testing"

func TestSetValueIntoGlobalContext(t *testing.T) {
	ParseTestAndRunBlock(t, "set chipotle \"sauce\"", ExpectValueSetTo(t, "chipotle", "sauce"))
	ParseTestAndRunBlock(t, "set to_many_arguments chipotle \"sauce\"", ExpectErrorValueAt(t, 1))
}

func TestSetValueIntoGlobalContextUsingShortcut(t *testing.T) {
	ParseTestAndRunBlock(t, "chipotle: \"sauce\"", ExpectValueSetTo(t, "chipotle", "sauce"))
}

func TestAssignBlockWithSet(t *testing.T) {
	ParseTestAndRunBlock(t, "set peppers {jalapeno:1; habanero: 2} |type", ExpectValue(t, NewIdentifier("dict")))
}

func TestAssignBlockWithLet(t *testing.T) {
	ParseTestAndRunBlock(t, "let peppers {jalapeno:1; habanero: 2} |type", ExpectValue(t, NewIdentifier("block")))
}

func TestEcho(t *testing.T) {
	ParseTestAndRunBlock(t, `echo chipotle`, ExpectValue(t, NewIdentifier("chipotle")))
	ParseTestAndRunBlock(t, `echo chipotle # njam`, ExpectValue(t, NewIdentifier("chipotle")))
}

func TestToS(t *testing.T) {
	ParseTestAndRunBlock(t, `to_s chipotle`, ExpectValue(t, NewStringLiteral("chipotle")))
}

func TestFirst(t *testing.T) {
	ParseTestAndRunBlock(t, `first`, ExpectValue(t, Nothing))
}

func TestType(t *testing.T) {

	ParseTestAndRunBlock(t, `type chipotle`, ExpectValue(t, NewIdentifier("identifier")))
	ParseTestAndRunBlock(t, `type "chipotle"`, ExpectValue(t, NewIdentifier("string")))
	ParseTestAndRunBlock(t, `type 3`, ExpectValue(t, NewIdentifier("int")))
	ParseTestAndRunBlock(t, `type 3.0`, ExpectValue(t, NewIdentifier("float")))
	ParseTestAndRunBlock(t, `type []`, ExpectValue(t, NewIdentifier("list")))

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
     get`, ExpectErrorValueAt(t, 2))

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
		`result: (incr chipotle "galapeno")
   	 type $result`, ExpectValue(t, NewIdentifier("error")))

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

func TestDefined(t *testing.T) {

	ParseTestAndRunBlock(t,
		`defined`, ExpectErrorValueAt(t, 1))

	ParseTestAndRunBlock(t,
		`defined chipotle`, ExpectValue(t, False))

	ParseTestAndRunBlock(t,
		`defined 3`, ExpectValue(t, False))

	ParseTestAndRunBlock(t,
		`set chipotle "sauce"
     defined chipotle`, ExpectValue(t, True))
}

func TestUserDefinedFunctionWithoutArguments(t *testing.T) {

	ParseTestAndRunBlock(t,
		`func`, ExpectErrorValueAt(t, 1))

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

func TestUserDefinedFunctionReUse(t *testing.T) {

	ParseTestAndRunBlock(t,
		`func chipotle`, ExpectErrorValueAt(t, 1))

	ParseTestAndRunBlock(t,
		`fsauce: (func {
			 return "chipotle"
		 })
		 fsoup: (func fsauce)
		 set soup (fsoup)`, ExpectValueSetTo(t, "soup", "chipotle"))

	ParseTestAndRunBlock(t,
		`fsauce: (func {
 			 return "chipotle"
 		 })
 		 fsoup: &fsauce
 		 set soup (fsoup)`, ExpectValueSetTo(t, "soup", "chipotle"))

	ParseTestAndRunBlock(t,
		`fsauce: { a: (func {
  			 return "chipotle"
  	 })}
  	 fsoup: &fsauce.a
  	 set soup (fsoup)`, ExpectValueSetTo(t, "soup", "chipotle"))

	ParseTestAndRunBlock(t,
		`fsauce: { a: {b: (func {
	 			 return "chipotle"
	 	 })}}
	 	 fsoup: &fsauce.a.b
	 	 set soup (fsoup)`, ExpectValueSetTo(t, "soup", "chipotle"))

	ParseTestAndRunBlock(t,
		`&fsauce.a`, ExpectErrorValueAt(t, 1))
}

func TestUserDefinedFunctionInFunction(t *testing.T) {

	ParseTestAndRunBlock(t,
		`sauce: (func {
			 pepper: "chipotle"
			 return (func {
				 return $pepper
			 })
		 })
		 soup: (sauce)
		 soup`, ExpectValue(t, NewStringLiteral("chipotle")))

	ParseTestAndRunBlock(t,
		`sauce: (func pepper {
			 pepper: "jalapeno"
 			 return (func {
 				 return $pepper
 			 })
 		 })
 		 soup: (sauce "chipotle")
 		 soup`, ExpectValue(t, NewStringLiteral("jalapeno")))
}

func TestUserDefinedFunctionWithOneArgument(t *testing.T) {

	ParseTestAndRunBlock(t,
		`func 3`, ExpectErrorValueAt(t, 1))

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

func TestUserDefinedFunctionWithHelp(t *testing.T) {

	ParseTestAndRunBlock(t,
		`func "sauce from heaven"`, ExpectErrorValueAt(t, 1))

	ParseTestAndRunBlock(t,
		`fsauce: (func "sauce from heaven" {

		 })
		 help fsauce`, ExpectValue(t, NewStringLiteral("sauce from heaven")))

	ParseTestAndRunBlock(t,
		`fsauce: (func "sauce from heaven" chipotle {
			return $chipotle
		 })
		 help fsauce`, ExpectValue(t, NewStringLiteral("sauce from heaven")))

	ParseTestAndRunBlock(t,
		`fsauce: (func "sauce from heaven" chipotle {
	 		return $chipotle
	 	 })
	 	 fsauce "jalapeno"`, ExpectValue(t, NewStringLiteral("jalapeno")))

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
		`d: {a: 2; b: 4}
		 mixin $d
	 	 b`, ExpectValue(t, NewIntegerLiteral(4)))

	ParseTestAndRunBlock(t,
		`d: {b: "chipotle"}
		 mixin $d
 	 	 b`, ExpectValue(t, NewStringLiteral("chipotle")))

	ParseTestAndRunBlock(t,
		`hot_or_not: {
			set chipotle (false)
			set galapeno (true)
		 }
		 well_tell_me: {
		   mixin (hot_or_not)
	   }
		 mixin $well_tell_me
  	 chipotle`, ExpectValue(t, NewBooleanLiteral(false)))
}

func TestDictionaryAccessShortcut(t *testing.T) {
	ParseTestAndRunBlock(t,
		`io: {
		  read: (func {
		  	return "chipotle"
		  })
		}
		io.read`, ExpectValue(t, NewStringLiteral("chipotle")))

	ParseTestAndRunBlock(t,
		`io: {
			read: (func name {
				return (name)
			})
		}
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
		`sauce: {
		  hot: (func {
		  	return "chipotle"
		  })
			same: (func {
				return (this.hot)
			})
		}
		sauce.same`, ExpectValue(t, NewStringLiteral("chipotle")))

	ParseTestAndRunBlock(t,
		`sauce: {
		  hot: (func {
		  	return "chipotle"
		  })
			same: (func {
				return (this.hot)
			})
		}
		sauce.same
		this`, ExpectErrorValueAt(t, 10))

	ParseTestAndRunBlock(t,
		` soup: (func {
				return (this.hot) # will fail, this is not defined
		  })
		  sauce: {
			  hot: (func {
				  return "chipotle"
			  })
			  same: (func {
				  return (soup)
			  })
		  }
		  sauce.same`, ExpectErrorValueAt(t, 2))
}

func TestLoad(t *testing.T) {

	context := NewGlobalContext()

	context.RegisterModule(NewModule("yippie", func(context RunContext) Value {

		return NewMappingForModule(context, []NamedValue{NewGoFunctionWithHelp("nop", "", func(context RunContext, arguments []Argument) Value {
			return Nothing
		})})

	}))

	ParseTestAndRunBlockWithinContext(t, context,
		`yy: (load "yippie")
		 yy.nop`, ExpectNothing(t))
}

func TestEvalWithBlock(t *testing.T) {
	ParseTestAndRunBlock(t,
		`eval`, ExpectErrorValueAt(t, 1))

	ParseTestAndRunBlock(t,
		`eval {
			"chipotle"
		 }`, ExpectValue(t, NewStringLiteral("chipotle")))

	ParseTestAndRunBlock(t,
		`pepper: (func block {
	 		return (eval $block)
	 	 })
		 pepper {
		   "chipotle"
		 }`, ExpectValue(t, NewStringLiteral("chipotle")))

	ParseTestAndRunBlock(t,
		`eval {
			pepper: "chipotle"
		 } {
		 	$pepper
		 }`, ExpectValue(t, NewStringLiteral("chipotle")))

	ParseTestAndRunBlock(t,
		`sauce: (func block {
		 		return (eval {pepper:"chipotle"} $block)
		 })
		 sauce {
		 	$pepper
		 }`, ExpectValue(t, NewStringLiteral("chipotle")))

	ParseTestAndRunBlock(t,
		`sauce: (func ingredient block {
	 		 return (eval {$ingredient:"chipotle"} $block)
	 	 })
	 	 sauce pepper {
		   assert (defined pepper)
	 		 $pepper
	 	 }`, ExpectValue(t, NewStringLiteral("chipotle")))

	ParseTestAndRunBlock(t,
		`sauce: (func block {
		 		 return (eval $block)
		 })
		 soup: (func {
		 	 pepper: "chipotle"
		 	 sauce {
		 		 $pepper
		 	 }
		 })
		 soup`, ExpectValue(t, NewStringLiteral("chipotle")))

	ParseTestAndRunBlock(t,
		`sauce: (func block {
		   return (eval {pepper:"jalapeno"} $block)
		 })
		 soup: (func {
			 pepper: "chipotle"
			 sauce {
	 			$pepper
	 		 }
		 })
		 soup`, ExpectValue(t, NewStringLiteral("chipotle")))

}

func TestEvalWithString(t *testing.T) {
	ParseTestAndRunBlock(t,
		`s:&"-\{a}-"; a:1; eval $s`, ExpectValue(t, NewStringLiteral("-1-")))
	ParseTestAndRunBlock(t,
		`s:&"-\{a}-"; b:{a:1}; eval $b $s`, ExpectValue(t, NewStringLiteral("-1-")))
}

func TestParseWithEval(t *testing.T) {
	ParseTestAndRunBlock(t,
		`s:"\"-\\{a}-\""; a:3; eval (parse $s)`, ExpectValue(t, NewStringLiteral("-3-")))
}

func TestEq(t *testing.T) {
	ParseTestAndRunBlock(t, `eq 1 1`, ExpectValue(t, True))
	ParseTestAndRunBlock(t, `eq 1 0`, ExpectValue(t, False))
	ParseTestAndRunBlock(t, `eq 1 "1"`, ExpectValue(t, False))
	ParseTestAndRunBlock(t, `eq "1" "1"`, ExpectValue(t, True))
	ParseTestAndRunBlock(t, `eq "1" [1]`, ExpectValue(t, False))
	ParseTestAndRunBlock(t, `eq [1] [1]`, ExpectValue(t, True))
	ParseTestAndRunBlock(t, `eq [1 2] [1 2]`, ExpectValue(t, True))
	ParseTestAndRunBlock(t, `eq [1] [0]`, ExpectValue(t, False))
	ParseTestAndRunBlock(t, `eq [1 2] [0 1]`, ExpectValue(t, False))
	ParseTestAndRunBlock(t, `eq`, ExpectErrorValueAt(t, 1))
	ParseTestAndRunBlock(t, `eq 1`, ExpectErrorValueAt(t, 1))
	ParseTestAndRunBlock(t, `eq 1 1 1`, ExpectErrorValueAt(t, 1))
	ParseTestAndRunBlock(t, `eval {m1: {}; m2: {}; eq $m1 $m2}`, ExpectValue(t, True))
	ParseTestAndRunBlock(t, `eval {m1: {a:1}; m2: {a:1}; eq $m1 $m2}`, ExpectValue(t, True))
	ParseTestAndRunBlock(t, `eval {m1: {a:1}; m2: {a:2}; eq $m1 $m2}`, ExpectValue(t, False))
	ParseTestAndRunBlock(t, `eval {m1: {a:1}; m2: {b:1}; eq $m1 $m2}`, ExpectValue(t, False))
	ParseTestAndRunBlock(t, `eval {m1: {a:1}; m2: {}; eq $m1 $m2}`, ExpectValue(t, False))
	ParseTestAndRunBlock(t, `eval {m1: {p:{x:1;y:1}}; m2: {p:{x:1;y:1}}; eq $m1 $m2}`, ExpectValue(t, True))
	ParseTestAndRunBlock(t, `eval {m1: {p:{x:1;y:2}}; m2: {p:{x:1;y:1}}; eq $m1 $m2}`, ExpectValue(t, False))
	ParseTestAndRunBlock(t, `eval {m1: {p:{x:1;z:1}}; m2: {p:{x:1;y:1}}; eq $m1 $m2}`, ExpectValue(t, False))
}

func TestNe(t *testing.T) {
	ParseTestAndRunBlock(t, `ne 1 1`, ExpectValue(t, False))
	ParseTestAndRunBlock(t, `ne 1 0`, ExpectValue(t, True))
	ParseTestAndRunBlock(t, `ne 1 "1"`, ExpectValue(t, True))
	ParseTestAndRunBlock(t, `ne "1" "1"`, ExpectValue(t, False))
	ParseTestAndRunBlock(t, `ne "1" [1]`, ExpectValue(t, True))
	ParseTestAndRunBlock(t, `ne [1] [1]`, ExpectValue(t, False))
	ParseTestAndRunBlock(t, `ne [1 2] [1 2]`, ExpectValue(t, False))
	ParseTestAndRunBlock(t, `ne [1] [0]`, ExpectValue(t, True))
	ParseTestAndRunBlock(t, `ne [1 2] [0]`, ExpectValue(t, True))
	ParseTestAndRunBlock(t, `ne [0] [1 2]`, ExpectValue(t, True))
	ParseTestAndRunBlock(t, `ne [0 [1]] [0 [2]]`, ExpectValue(t, True))
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

	ParseTestAndRunBlock(t, `gt [1] [1]`, ExpectValue(t, False))
	ParseTestAndRunBlock(t, `gt [1 2] [1]`, ExpectValue(t, True))

	ParseTestAndRunBlock(t, `gt 1.0 1.0`, ExpectValue(t, False))
	ParseTestAndRunBlock(t, `gt 2.0 1.0`, ExpectValue(t, True))
	ParseTestAndRunBlock(t, `gt 0.0 1.0`, ExpectValue(t, False))
	ParseTestAndRunBlock(t, `gt -1.0 0.0`, ExpectValue(t, False))
	ParseTestAndRunBlock(t, `gt -1.0 -2.0`, ExpectValue(t, True))

	ParseTestAndRunBlock(t, `gt 1 1.0`, ExpectErrorValueAt(t, 1))
	ParseTestAndRunBlock(t, `gt 1.0 1`, ExpectErrorValueAt(t, 1))
	ParseTestAndRunBlock(t, `gt "chipotle" "galapeno"`, ExpectValue(t, False))

	ParseTestAndRunBlock(t, `b: { _compare: (func v { return 1 }) }; s: { _compare: (func v { return -1 }) }; (gt $b $s)`, ExpectValue(t, True))

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

	ParseTestAndRunBlock(t, `lt [1] [1]`, ExpectValue(t, False))
	ParseTestAndRunBlock(t, `lt [1] [1 2]`, ExpectValue(t, True))
	ParseTestAndRunBlock(t, `lt [[1 2] [3 4] [5 6]] [[1 2] [3 5]]`, ExpectValue(t, True))

	ParseTestAndRunBlock(t, `lt 1.0 1.0`, ExpectValue(t, False))
	ParseTestAndRunBlock(t, `lt 2.0 1.0`, ExpectValue(t, False))
	ParseTestAndRunBlock(t, `lt 0.0 1.0`, ExpectValue(t, True))
	ParseTestAndRunBlock(t, `lt -1.0 0.0`, ExpectValue(t, True))
	ParseTestAndRunBlock(t, `lt -1.0 -2.0`, ExpectValue(t, False))

	ParseTestAndRunBlock(t, `b: { _compare: (func v { return 1 }) }; s: { _compare: (func v { return -1 }) }; (lt $s $b)`, ExpectValue(t, True))
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

func TestCompare(t *testing.T) {
	ParseTestAndRunBlock(t, `compare 1 1`, ExpectValue(t, NewIntegerLiteral(0)))
	ParseTestAndRunBlock(t, `compare "chipotle" "jalapeno"`, ExpectValue(t, NewIntegerLiteral(-1)))
	ParseTestAndRunBlock(t, `compare "chipotle" 0`, ExpectValue(t, NewIntegerLiteral(1)))
	ParseTestAndRunBlock(t, `compare [1 2 3] [1 2]`, ExpectValue(t, NewIntegerLiteral(1)))
	ParseTestAndRunBlock(t, `compare $false $true`, ExpectValue(t, NewIntegerLiteral(-1)))
	ParseTestAndRunBlock(t, `compare $false $false`, ExpectValue(t, NewIntegerLiteral(0)))
	ParseTestAndRunBlock(t, `compare $true $false`, ExpectValue(t, NewIntegerLiteral(1)))
	ParseTestAndRunBlock(t, `compare $true $true`, ExpectValue(t, NewIntegerLiteral(0)))
	ParseTestAndRunBlock(t, `compare $true "chipotle"`, ExpectErrorValueAt(t, 1))
	ParseTestAndRunBlock(t, `e: { _compare: (func v { return 0 }) }; d: {}; compare $e $d`, ExpectValue(t, NewIntegerLiteral(0)))
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

func TestAssert(t *testing.T) {
	ParseTestAndRunBlock(t, `assert`, ExpectErrorValueAt(t, 1))
	ParseTestAndRunBlock(t, `assert (false)`, ExpectErrorValueAt(t, 1))
	ParseTestAndRunBlock(t, `assert (true) "everything fine"`, ExpectValue(t, True))
	ParseTestAndRunBlock(t, `assert (false) "things got really messy"`, ExpectErrorValueAt(t, 1))
	ParseTestAndRunBlock(t, `assert "chipotle" "things got really messy"`, ExpectErrorValueAt(t, 1))
}

func TestError(t *testing.T) {
	ParseTestAndRunBlock(t, `error`, ExpectErrorValueAt(t, 1))
	ParseTestAndRunBlock(t, `error "chipotle"`, ExpectErrorValueAt(t, 1))
	ParseTestAndRunBlock(t,
		`err: (error "chipotle")
		 "chipotle"`, ExpectValue(t, NewStringLiteral("chipotle")))
}

func TestPanic(t *testing.T) {
	ParseTestAndRunBlock(t, `panic`, ExpectErrorValueAt(t, 1))
	ParseTestAndRunBlock(t, `panic "chipotle"`, ExpectErrorValueAt(t, 1))
	ParseTestAndRunBlock(t,
		`err: (panic "chipotle")
		 "chipotle"`, ExpectErrorValueAt(t, 1))
}

func TestHelp(t *testing.T) {

	ParseTestAndRunBlock(t,
		`help help help`, ExpectErrorValueAt(t, 1))

	ParseTestAndRunBlock(t,
		`type (help)`, ExpectValue(t, NewIdentifier("list")))

	ParseTestAndRunBlock(t,
		`help help`, ExpectValue(t, NewStringLiteral("Get help. Usage 'help' or 'help symbol'")))

	ParseTestAndRunBlock(t,
		`chipotle: (func `+"`"+`not so hot pepper
			from mexico
			 love it`+"`"+` {})
		 help chipotle`, ExpectValue(t, NewStringLiteral("not so hot pepper\nfrom mexico\nlove it")))

}

func TestClosee(t *testing.T) {

	ParseTestAndRunBlock(t,
		`close`, ExpectErrorValueAt(t, 1))

	ParseTestAndRunBlock(t,
		`close 3`, ExpectErrorValueAt(t, 1))
}

func TestFreeze(t *testing.T) {

	ParseTestAndRunBlock(t,
		`freeze!`, ExpectErrorValueAt(t, 1))

	ParseTestAndRunBlock(t,
		`frozen`, ExpectErrorValueAt(t, 1))

	ParseTestAndRunBlock(t,
		`pepper: "chipotle"
		 frozen $pepper`, ExpectValue(t, True))

	ParseTestAndRunBlock(t,
		`pepper: ["chipotle"]
	 	 frozen $pepper`, ExpectValue(t, False))

	ParseTestAndRunBlock(t,
		`pepper: { name: "chipotle" }
 	 	 frozen $pepper`, ExpectValue(t, False))

	ParseTestAndRunBlock(t,
		`pepper: ["chipotle"]
		 freeze! $pepper
		 frozen $pepper`, ExpectValue(t, True))

	ParseTestAndRunBlock(t,
		`pepper: []
		 freeze! $pepper
		 frozen $pepper`, ExpectValue(t, True))

	ParseTestAndRunBlock(t,
		`pepper: {}
		 freeze! $pepper
		 frozen $pepper`, ExpectValue(t, True))

}

func TestUUID(t *testing.T) {

	ParseTestAndRunBlock(t,
		`pepper: {}
		 i: (uuid $pepper)
		 same: $pepper 
		 eq $i (uuid $same)`, ExpectValue(t, True))

	ParseTestAndRunBlock(t,
		`ne (uuid) (uuid)`, ExpectValue(t, True))
}

func TestTime(t *testing.T) {

	ParseTestAndRunBlock(t,
		`(time "2019-10-12T14:20:50.52+07:00") zoneOffset |divide 60 |divide 60 |eq 7`, ExpectValue(t, True))

	ParseTestAndRunBlock(t,
		`(time "ANSIC" "Mon Jan 2 15:04:05 2006") year |eq 2006`, ExpectValue(t, True))

	ParseTestAndRunBlock(t,
		`(time UnixDate "Mon Jan 2 15:04:05 MST 2006") month |eq 1`, ExpectValue(t, True))

	ParseTestAndRunBlock(t,
		`(time RubyDate "Mon Jan 02 15:04:05 -0700 2006") day |eq 2`, ExpectValue(t, True))

	ParseTestAndRunBlock(t,
		`(time RFC822 "02 Jan 06 15:04 MST") year |eq 2006`, ExpectValue(t, True))

	ParseTestAndRunBlock(t,
		`(time RFC822Z "02 Jan 06 15:04 -0700") hour |eq 15`, ExpectValue(t, True))

	ParseTestAndRunBlock(t,
		`(time RFC850 "Monday, 02-Jan-06 15:04:05 MST") minute |eq 4`, ExpectValue(t, True))

	ParseTestAndRunBlock(t,
		`(time RFC1123 "Mon, 02 Jan 2006 15:04:05 MST") second |eq 5`, ExpectValue(t, True))

	ParseTestAndRunBlock(t,
		`(time RFC1123Z "Mon, 02 Jan 2006 15:04:05 -0700") zoneOffset |divide 60 |divide 60 |eq -7`, ExpectValue(t, True))

	ParseTestAndRunBlock(t,
		`(time RFC3339 "2006-01-02T15:04:05+07:00") minute |eq 4`, ExpectValue(t, True))

	ParseTestAndRunBlock(t,
		`(time RFC3339Nano "2006-01-02T15:04:05.999999999Z") nano |eq 999999999`, ExpectValue(t, True))

	ParseTestAndRunBlock(t,
		`(time Kitchen "3:04PM") hour |eq 15`, ExpectValue(t, True))

	ParseTestAndRunBlock(t,
		`(time Kitchen "3:04PM") hour |eq 15`, ExpectValue(t, True))

	ParseTestAndRunBlock(t,
		`t: ((time) timestamp)
		 t2: ((time $t) timestamp)
		 eq $t $t2`, ExpectValue(t, True))
}

func TestLen(t *testing.T) {
	ParseTestAndRunBlock(t, `len`, ExpectErrorValueAt(t, 1))
	ParseTestAndRunBlock(t, `len (nil)`, ExpectValue(t, NewIntegerLiteral(0)))
	ParseTestAndRunBlock(t, `len $true`, ExpectErrorValueAt(t, 1))
	ParseTestAndRunBlock(t, `len chipotle`, ExpectValue(t, NewIntegerLiteral(8)))
	ParseTestAndRunBlock(t, `len "chipotle"`, ExpectValue(t, NewIntegerLiteral(8)))
	ParseTestAndRunBlock(t, `len ["c" "h" "i" "p" "o" "t" "l" "e"]`, ExpectValue(t, NewIntegerLiteral(8)))
}
