package data

import (
	"testing"

	elmo "github.com/okke/elmo/core"
)

func testMoFile(t *testing.T, mo string, contextInitializer func(elmo.RunContext)) {
	context := elmo.NewGlobalContext()
	contextInitializer(context)

	loader := elmo.NewLoader(context, []string{"./test"})

	loaded := loader.Load(mo)

	if loaded.Type() == elmo.TypeError {
		t.Error("error while loading", mo, ":", loaded.String())
	}
}

func initTestContext(context elmo.RunContext) {
	context.RegisterModule(Module)
}

func TestConvertCSV(t *testing.T) {
	testMoFile(t, "peppers", initTestContext)
}

func TestConvertJSON(t *testing.T) {
	testMoFile(t, "habanero", initTestContext)
}
