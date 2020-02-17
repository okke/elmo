package elmohttp

import (
	"testing"

	elmo "github.com/okke/elmo/core"
)


func initTestContext(context elmo.RunContext) {
	context.RegisterModule(Module)
}

func TestClient(t *testing.T) {
	elmo.TestMoFile(t, "client", initTestContext)
}
