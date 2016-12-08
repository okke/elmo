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
