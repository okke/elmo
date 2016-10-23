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

func TestJoin(t *testing.T) {

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
     str.join`, elmo.ExpectErrorValueAt(t, 2))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
		 str.join "chipotle" " " "sauce"`, elmo.ExpectValue(t, elmo.NewStringLiteral("chipotle sauce")))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
     pepper: "chipotle"
     str.join "sauce" (pepper 0 99)`, elmo.ExpectErrorValueAt(t, 3))

}
