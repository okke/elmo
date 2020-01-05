package elmo

import "testing"

func reconstruct(value Value) Value {
	return Serialize(value).ToValue()
}

func TestSerializeAndReconstructList(t *testing.T) {

	reconstructed := reconstruct(NewListValue([]Value{}))

	if reconstructed == nil {
		t.Error("expected to reconstruct something")
	}

	if reconstructed.Type() != TypeList {
		t.Error("expected list")
	}

	reconstructed.(ListValue).Append(NewStringLiteral("chipotle"))
	reconstructed.(ListValue).Append(NewStringLiteral("jalapeno"))

	reconstructed = reconstruct(reconstructed)

	if len(reconstructed.Internal().([]Value)) != 2 {
		t.Error("expected two literals")
	}

	if reconstructed.Internal().([]Value)[0].String() != "chipotle" {
		t.Error("expected chipotle as first value")
	}

	if reconstructed.Internal().([]Value)[1].String() != "jalapeno" {
		t.Error("expected jalapeno as first value")
	}
}

func TestSerializeAndReconstructListInList(t *testing.T) {
	l1 := NewListValue([]Value{NewStringLiteral("chipotle"), NewStringLiteral("jalapeno")})
	l2 := NewListValue([]Value{NewStringLiteral("habanero"), NewStringLiteral("piri piri")})

	reconstructed := reconstruct(NewListValue([]Value{l1, l2}))

	if len(reconstructed.Internal().([]Value)) != 2 {
		t.Error("expected two elements not", len(reconstructed.Internal().([]Value)), reconstructed)
	}

	if reconstructed.Internal().([]Value)[0].UUID().String() != l1.UUID().String() {
		t.Error("expected uuid's to be the same for first element")
	}

	if reconstructed.Internal().([]Value)[1].UUID().String() != l2.UUID().String() {
		t.Error("expected uuid's to be the same for second element")
	}
}

func TestSerializeAndReconstructDictionary(t *testing.T) {
	keyValues := map[string]Value{}
	dict := NewDictionaryValue(nil, keyValues)
	serialized := Serialize(dict)
	reconstructed := serialized.ToValue()

	if reconstructed == nil {
		t.Error("expected to reconstruct something")
	}

	if reconstructed.Type() != TypeDictionary {
		t.Error("expected dictionary")
	}

}

func TestSerializeAndReconstructListInDict(t *testing.T) {
	l1 := NewListValue([]Value{NewStringLiteral("chipotle"), NewStringLiteral("jalapeno")})
	l2 := NewListValue([]Value{NewStringLiteral("habanero"), NewStringLiteral("piri piri")})

	keyValues := map[string]Value{"hot": l1, "hotter": l2}

	reconstructed := reconstruct(NewDictionaryValue(nil, keyValues))

	l1r, f1 := reconstructed.(DictionaryValue).Resolve("hot")
	l2r, f2 := reconstructed.(DictionaryValue).Resolve("hotter")

	if !f1 {
		t.Error("expected to find hot list")
	}

	if l1r.UUID().String() != l1.UUID().String() {
		t.Error("expected uuid's to be the same for hot list")
	}

	if !f2 {
		t.Error("expected to find hotter list")
	}

	if l2r.UUID().String() != l2.UUID().String() {
		t.Error("expected uuid's to be the same for hotter list")
	}

}
