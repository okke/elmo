package elmo

// ExpectValueSetTo expects a given variable is set to a given value
import (
	"reflect"
	"testing"
)

// ExpectValueSetTo returns a function that checks if a value is set in context
//
func ExpectValueSetTo(t *testing.T, key string, value string) func(RunContext, Value) {

	return func(context RunContext, blockResult Value) {

		if blockResult.Type() == TypeError {
			t.Error(blockResult.(ErrorValue).Error())
			return
		}

		result, found := context.Get(key)

		if !found {
			t.Errorf("expected %s to be set", key)
		} else {
			if result.String() != value {
				t.Errorf("expected %s to be set to (%s), found %s", key, value, result.String())
			}
		}
	}
}

// ExpectErrorValueAt returns a function that expects an error on a given line number
//
func ExpectErrorValueAt(t *testing.T, lineno int) func(RunContext, Value) {

	return func(context RunContext, blockResult Value) {

		if blockResult.Type() != TypeError {
			t.Errorf("expected error but found %v", blockResult)
			return
		}

		_, l := blockResult.(ErrorValue).At()

		if l != lineno {
			t.Errorf("expected error at line %d but found (%v) on line %d", lineno, blockResult.String(), l)
		}

	}
}

// ExpectNothing returns a function that expects evaluation returns Nothing
//
func ExpectNothing(t *testing.T) func(RunContext, Value) {

	return func(context RunContext, blockResult Value) {
		if blockResult != Nothing {
			t.Errorf("expected nothing but found %v", blockResult)
		}
	}
}

// ExpectValue returns a function that expects evaluation returns a given value
//
func ExpectValue(t *testing.T, value Value) func(RunContext, Value) {

	return func(context RunContext, blockResult Value) {
		if !reflect.DeepEqual(blockResult, value) {
			t.Errorf("expected (%v) but found (%v)", value, blockResult)
		}
	}
}
