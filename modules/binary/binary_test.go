package bin

import (
	"testing"

	elmo "github.com/okke/elmo/core"
)

func binContext() elmo.RunContext {
	context := elmo.NewGlobalContext()
	context.RegisterModule(Module)
	return context
}

func TestNewWithoutParent(t *testing.T) {

	elmo.ParseTestAndRunBlockWithinContext(t, binContext(),
		`bin: (load bin)
		 bin.new`, elmo.ExpectErrorValueAt(t, 2))

	elmo.ParseTestAndRunBlockWithinContext(t, binContext(),
		`bin: (load bin)
     bin.new {}`, elmo.ExpectErrorValueAt(t, 2))

	elmo.ParseTestAndRunBlockWithinContext(t, binContext(),
		`bin: (load bin)
 		 b: (bin.new "chipotle")
     type $b`, elmo.ExpectValue(t, elmo.NewIdentifier("binary")))

}

func TestToValue(t *testing.T) {

	elmo.ParseTestAndRunBlockWithinContext(t, binContext(),
		`bin: (load bin)
		 assert (bin.new 5 | bin.toValue | eq 5)`, elmo.ExpectValue(t, elmo.True))

	elmo.ParseTestAndRunBlockWithinContext(t, binContext(),
		`bin: (load bin)
		 assert (bin.new 3.1415 | bin.toValue | eq 3.1415)`, elmo.ExpectValue(t, elmo.True))

	elmo.ParseTestAndRunBlockWithinContext(t, binContext(),
		`bin: (load bin)
		 assert (bin.new $true | bin.toValue | eq $true)`, elmo.ExpectValue(t, elmo.True))

	elmo.ParseTestAndRunBlockWithinContext(t, binContext(),
		`bin: (load bin)
		 assert (bin.new "chipotle" | bin.toValue | eq "chipotle")`, elmo.ExpectValue(t, elmo.True))

	elmo.ParseTestAndRunBlockWithinContext(t, binContext(),
		`bin: (load bin)
		 assert (bin.new jalapeno | bin.toValue | eq jalapeno)`, elmo.ExpectValue(t, elmo.True))

}

func TestBinaryLength(t *testing.T) {
	elmo.ParseTestAndRunBlockWithinContext(t, binContext(),
		`bin: (load bin)
		 bin.new jalapeno | len`, elmo.ExpectValue(t, elmo.NewIntegerLiteral(78)))
}
