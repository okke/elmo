package sys

import (
	"testing"

	elmo "github.com/okke/elmo/core"
	"github.com/okke/elmo/modules/list"
)

func sysContext() elmo.RunContext {
	context := elmo.NewGlobalContext()
	context.RegisterModule(Module)
	context.RegisterModule(list.Module)
	return context
}

func TestExec(t *testing.T) {

	elmo.ParseTestAndRunBlockWithinContext(t, sysContext(),
		`mixin (load sys)
     ls "./testdata" |exec`, elmo.ExpectValue(t, elmo.ParseAndRun(elmo.NewGlobalContext(), "[\"chipotle.txt\" \"jalapeno.txt\"]")))

	elmo.ParseTestAndRunBlockWithinContext(t, sysContext(),
		`mixin (load sys)
	   ls "./testdata" |wc "-l"|exec`, elmo.ExpectValue(t, elmo.ParseAndRun(elmo.NewGlobalContext(), "[\"       2\"]")))

	elmo.ParseTestAndRunBlockWithinContext(t, sysContext(),
		`mixin (load sys)
     chipotle |exec`, elmo.ExpectErrorValueAt(t, 2))

	elmo.ParseTestAndRunBlockWithinContext(t, sysContext(),
		`mixin (load sys)
     ls | chipotle |exec`, elmo.ExpectErrorValueAt(t, 2))
}

func TestAsList(t *testing.T) {
	elmo.ParseTestAndRunBlockWithinContext(t, sysContext(),
		`mixin (load sys)
     mixin (load list)
     ls "./testdata" |append "espelette.txt"`, elmo.ExpectValue(t, elmo.ParseAndRun(elmo.NewGlobalContext(), "[\"chipotle.txt\" \"jalapeno.txt\" \"espelette.txt\"]")))
}
