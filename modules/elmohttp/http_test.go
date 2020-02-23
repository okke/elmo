package elmohttp

import (
	"testing"

	elmo "github.com/okke/elmo/core"
	dict "github.com/okke/elmo/modules/dictionary"
	"github.com/okke/elmo/modules/str"
)

func initTestContext(context elmo.RunContext) {
	context.RegisterModule(Module)
	context.RegisterModule(dict.Module)
	context.RegisterModule(str.Module)
}

func TestTestServer(t *testing.T) {
	elmo.TestMoFile(t, "testserver", initTestContext)
}

func TestClient(t *testing.T) {
	elmo.TestMoFile(t, "client", initTestContext)
}
