package elmo

// ExpectValueSetTo expects a given variable is set to a given value
import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func getCallingFunc() string {
	skip := 0

	_, fn, line, ok := runtime.Caller(skip)

	for ok {

		splitted := strings.Split(fn, "/")
		last := splitted[len(splitted)-1]
		if strings.HasSuffix(last, "_test.go") {
			return fmt.Sprintf("%v:%d", last, line)
		}

		skip = skip + 1
		_, fn, line, ok = runtime.Caller(skip)
	}

	return "unknown calling function"

}

// ExpectValueSetTo returns a function that checks if a value is set in context
//
func ExpectValueSetTo(t *testing.T, key string, value string) func(RunContext, Value) {

	return func(context RunContext, blockResult Value) {

		if blockResult.Type() == TypeError {
			t.Errorf("%v at %s", blockResult.(ErrorValue).Error(), getCallingFunc())
			return
		}

		result, found := context.Get(key)

		if !found {
			t.Errorf("expected %s to be set at %s", key, getCallingFunc())
		} else {
			if result.String() != value {
				t.Errorf("expected %s to be set to (%s), found %s at %s", key, value, result.String(), getCallingFunc())
			}
		}
	}
}

// ExpectErrorValueAt returns a function that expects an error on a given line number
//
func ExpectErrorValueAt(t *testing.T, lineno int) func(RunContext, Value) {

	return func(context RunContext, blockResult Value) {

		if blockResult.Type() != TypeError {
			t.Errorf("expected error but found %v at %s", blockResult, getCallingFunc())
			return
		}

		_, l := blockResult.(ErrorValue).At()

		if l != lineno {
			t.Errorf("expected error at line %d but found (%v) on line %d at %s", lineno, blockResult.String(), l, getCallingFunc())
		}

	}
}

// ExpectNothing returns a function that expects evaluation returns Nothing
//
func ExpectNothing(t *testing.T) func(RunContext, Value) {

	return func(context RunContext, blockResult Value) {
		if blockResult != Nothing {
			t.Errorf("expected nothing but found %v at %s", blockResult, getCallingFunc())
		}
	}
}

// ExpectValue returns a function that expects evaluation returns a given value
//
func ExpectValue(t *testing.T, value Value) func(RunContext, Value) {

	return func(context RunContext, blockResult Value) {
		if value == Nothing {
			if blockResult != Nothing {
				t.Errorf("expected value nothing but found %v of type %v at %s", blockResult, blockResult.Info().Name(), getCallingFunc())
			}
		}
		if !reflect.DeepEqual(blockResult, value) {
			t.Errorf("expected value %v of type %v but found %v of type %v at %s", value, value.Info().Name(), blockResult, blockResult.Info().Name(), getCallingFunc())
		}
	}
}

// ExpectValues returns a function that expects evaluation returns specified values
// as ReturnValue
//
func ExpectValues(t *testing.T, values ...Value) func(RunContext, Value) {

	return func(context RunContext, blockResult Value) {
		if !reflect.DeepEqual(blockResult, NewReturnValue(values)) {
			t.Errorf("expected value %v but found %v at %s", values, blockResult, getCallingFunc())
		}
	}
}
