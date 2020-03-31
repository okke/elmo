package data

import (
	"testing"

	elmo "github.com/okke/elmo/core"
)

func initTestContext(context elmo.RunContext) {
	context.RegisterModule(Module)
}

func TestConvertFromCSV(t *testing.T) {
	elmo.TestMoFile(t, "peppers", initTestContext)
}

func TestConvertFromJSON(t *testing.T) {
	elmo.TestMoFile(t, "habanero", initTestContext)
}

func TestConvertToJSON(t *testing.T) {
	elmo.TestMoFile(t, "to_json", initTestContext)
}
