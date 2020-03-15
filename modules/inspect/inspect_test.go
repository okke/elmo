package inspect

import (
	"testing"

	elmo "github.com/okke/elmo/core"
	"github.com/okke/elmo/modules/list"
	"github.com/okke/elmo/modules/str"
)

func initTestContext(context elmo.RunContext) {
	context.RegisterModule(Module)
	context.RegisterModule(list.Module)
	context.RegisterModule(str.Module)
}

func TestMeta(t *testing.T) {
	elmo.TestMoFile(t, "meta", initTestContext)
}

func TestCalls(t *testing.T) {
	elmo.TestMoFile(t, "calls", initTestContext)
}

func TestBlock(t *testing.T) {
	elmo.TestMoFile(t, "block", initTestContext)
}
