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

	elmo.ParseTestAndRunBlockWithinContext(t, listContext(),
		`el: (load "el")
    l1: (list 1 2 3)
    l2: (el.append (l1) 4)
		l1`, elmo.ExpectValue(t, elmo.ParseAndRun(elmo.NewGlobalContext(), "list 1 2 3")))

	elmo.ParseTestAndRunBlockWithinContext(t, listContext(),
		`el: (load "el")
    l: (list 1 2 3)
    l: (el.append (l) 4 5 6)
		l`, elmo.ExpectValue(t, elmo.ParseAndRun(elmo.NewGlobalContext(), "list 1 2 3 4 5 6")))
}
