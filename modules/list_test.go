package el

import (
	"testing"

	"github.com/okke/elmo/core"
)

func listContext() elmo.RunContext {
	context := elmo.NewGlobalContext()
	context.RegisterModule(Module)
	return context
}

func TestAppend(t *testing.T) {

	elmo.ParseTestAndRunBlockWithinContext(t, listContext(),
		`el: (load "el")
    l: (list 1 2 3)
    l: (el.append (l) 4)
		l`, elmo.ExpectValue(t, elmo.ParseAndRun(elmo.NewGlobalContext(), "list 1 2 3 4")))
}
