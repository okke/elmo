package elmo

import "testing"

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

func TestUserDefinedTemplate(t *testing.T) {

	ParseTestAndRunBlock(t,
		`pepper: "jalapeno"
		 t: (template &"\{$pepper}\{$pepper}")
		 t
		`, ExpectValue(t, NewStringLiteral("jalapenojalapeno")))

	ParseTestAndRunBlock(t,
		`t: (template pepper &"\{$pepper}\{$pepper}")
		 t "chipotle"
		`, ExpectValue(t, NewStringLiteral("chipotlechipotle")))

	ParseTestAndRunBlock(t,
		`
		 t: (template "help!" &"\{$pepper}\{$pepper}")
		 help t
	    `, ExpectValue(t, NewStringLiteral("help!")))

	ParseTestAndRunBlock(t,
		`t: (template "double pepper" pepper &"\{$pepper}\{$pepper}")
		 help t
		`, ExpectValue(t, NewStringLiteral("double pepper")))
}

func TestFuncWithOptionalArguments(t *testing.T) {
	ParseTestAndRunBlock(t,
		`f: (func i?5 {
	 		return (multiply $i $i)
	 	 })
		 (multiply (f 2) (f))`, ExpectValue(t, NewIntegerLiteral(100)))

	ParseTestAndRunBlock(t,
		`greet: (func name greeting?"Hello" { echo "\{$greeting} \{$name}"})
		(greet "chipotle")`, ExpectValue(t, NewStringLiteral("Hello chipotle")))
}

func TestFuncCallUsesCorrectScope(t *testing.T) {
	ParseTestAndRunBlock(t,
		`d: { f1: (func {return 1}); f2: (func f {return $f}) }
		 f1: 8
		 (d.f2 $f1)`, ExpectValue(t, NewIntegerLiteral(8)))

	ParseTestAndRunBlock(t,
		`d: { f1: (func {return 1}); f2: (func f {return $f}) }
		 (d.f2 $f1)`, ExpectErrorValueAt(t, 2))
}

func TestLoadedTemplatesUsesCorrectContext(t *testing.T) {
	TestMoFile(t, "loadtemplates", func(context RunContext) {})
}
