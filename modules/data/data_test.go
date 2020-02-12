package data

import (
	"testing"

	elmo "github.com/okke/elmo/core"
)

func TestConvertCSV(t *testing.T) {
	context := elmo.NewGlobalContext()
	context.RegisterModule(Module)

	loader := elmo.NewLoader(context, []string{"./test"})

	loaded := loader.Load("peppers")

	if loaded.Type() == elmo.TypeError {
		t.Error("error while loading peppers.mo:", loaded.String())
	}

}
