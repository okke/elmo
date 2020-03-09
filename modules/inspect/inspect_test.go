package inspect

import (
	"testing"

	elmo "github.com/okke/elmo/core"
)

func initTestContext(context elmo.RunContext) {
	context.RegisterModule(Module)
}

func TestMeta(t *testing.T) {
	elmo.TestMoFile(t, "meta", initTestContext)
}
