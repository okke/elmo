package ed

import (
	"testing"

	"github.com/okke/elmo/core"
)

func dictContext() elmo.RunContext {
	context := elmo.NewGlobalContext()
	context.RegisterModule(Module)
	return context
}

func TestKeys(t *testing.T) {

	elmo.ParseTestAndRunBlockWithinContext(t, dictContext(),
		`ed: (load "ed")
    peppers: {
      chipotle: {
        heat: 2
      }
      galapeno: {
        heat: 3
      }
    }
		ed.keys peppers`, elmo.ExpectValue(t, elmo.ParseAndRun(elmo.NewGlobalContext(), "list chipotle galapeno")))

	// keys function should return keys in sorted order
	//
	elmo.ParseTestAndRunBlockWithinContext(t, dictContext(),
		`ed: (load "ed")
    peppers: {
      galapeno: {
        heat: 3
      }
      chipotle: {
        heat: 2
      }
    }
    ed.keys peppers`, elmo.ExpectValue(t, elmo.ParseAndRun(elmo.NewGlobalContext(), "list chipotle galapeno")))

}
