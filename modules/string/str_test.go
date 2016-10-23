package str

import (
	"testing"

	elmo "github.com/okke/elmo/core"
)

func strContext() elmo.RunContext {
	context := elmo.NewGlobalContext()
	context.RegisterModule(Module)
	return context
}

func TestAt(t *testing.T) {

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
		 str.at`, elmo.ExpectErrorValueAt(t, 2))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
		 str.at "chipotle" 2`, elmo.ExpectValue(t, elmo.NewStringLiteral("i")))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
     pepper: "chipotle"
     str.at (pepper 0 99) 2`, elmo.ExpectErrorValueAt(t, 3))

}

func TestLen(t *testing.T) {

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
		 str.len`, elmo.ExpectErrorValueAt(t, 2))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
		 str.len "chipotle"`, elmo.ExpectValue(t, elmo.NewIntegerLiteral(8)))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
     pepper: "chipotle"
		 str.len (pepper 0 2)`, elmo.ExpectValue(t, elmo.NewIntegerLiteral(3)))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
     pepper: "chipotle"
     str.len (pepper 0 99)`, elmo.ExpectErrorValueAt(t, 3))

}

func TestConcat(t *testing.T) {

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
     str.concat`, elmo.ExpectErrorValueAt(t, 2))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
		 str.concat "chipotle" " " "sauce"`, elmo.ExpectValue(t, elmo.NewStringLiteral("chipotle sauce")))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
     pepper: "chipotle"
     str.concat "sauce" (pepper 0 99)`, elmo.ExpectErrorValueAt(t, 3))

}

func TestTrim(t *testing.T) {

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
     str.trim`, elmo.ExpectErrorValueAt(t, 2))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
     str.trim (error "chipotle")`, elmo.ExpectErrorValueAt(t, 2))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
     str.trim "chipotle" (error "chipotle")`, elmo.ExpectErrorValueAt(t, 2))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
     str.trim " chipotle "`, elmo.ExpectValue(t, elmo.NewStringLiteral("chipotle")))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
     str.trim left " chipotle "`, elmo.ExpectValue(t, elmo.NewStringLiteral("chipotle ")))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
	   str.trim left "!chipotle!" "!"`, elmo.ExpectValue(t, elmo.NewStringLiteral("chipotle!")))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
     str.trim left (error "!chipotle!") "!"`, elmo.ExpectErrorValueAt(t, 2))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
     str.trim right "!chipotle!" "!"`, elmo.ExpectValue(t, elmo.NewStringLiteral("!chipotle")))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
     str.trim right " chipotle "`, elmo.ExpectValue(t, elmo.NewStringLiteral(" chipotle")))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
     str.trim prefix "jar_with_chipotle" "jar_with_"`, elmo.ExpectValue(t, elmo.NewStringLiteral("chipotle")))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
     str.trim suffix "chipotle_in_a_jar" "_in_a_jar"`, elmo.ExpectValue(t, elmo.NewStringLiteral("chipotle")))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
      str.trim soup " chipotle "`, elmo.ExpectErrorValueAt(t, 2))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
      str.trim "<chipotle>" "<>"`, elmo.ExpectValue(t, elmo.NewStringLiteral("chipotle")))
}
