package data

import (
	"testing"

	elmo "github.com/okke/elmo/core"
)


func initTestContext(context elmo.RunContext) {
	context.RegisterModule(Module)
}

func TestConvertCSV(t *testing.T) {
	elmo.TestMoFile(t, "peppers", initTestContext)
}

func TestConvertJSON(t *testing.T) {
	elmo.TestMoFile(t, "habanero", initTestContext)
}
