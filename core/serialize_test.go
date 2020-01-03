package elmo

import "testing"

func TestSerializeAndReconstructList(t *testing.T) {
	list := NewListValue([]Value{})
	serialized := Serialize(list)
	reconstructed := serialized.ToValue()

	if reconstructed == nil {
		t.Error("expected to reconstruct something")
	}

}
