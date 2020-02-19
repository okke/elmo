package elmohttp

import (
	"testing"

	elmo "github.com/okke/elmo/core"
	dict "github.com/okke/elmo/modules/dictionary"
)

func initTestContext(context elmo.RunContext) {
	context.RegisterModule(Module)
	context.RegisterModule(dict.Module)
}

func TestClient(t *testing.T) {
	elmo.TestMoFile(t, "client", initTestContext)
}
