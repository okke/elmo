package elmo

import "fmt"

type errorValue struct {
	baseValue
	meta   ScriptMetaData
	lineno int
	msg    string
	fatal  bool
	ignore bool
}

func (errorValue *errorValue) String() string {
	kind := "error"
	if errorValue.IsFatal() {
		kind = "fatal error"
	}
	if errorValue.meta != nil {
		meta, lineno := errorValue.At()
		return fmt.Sprintf("%s(at %s at line %d: %s)", kind, meta.Name(), lineno, errorValue.msg)
	}
	return fmt.Sprintf("%s(%s)", kind, errorValue.msg)
}

func (errorValue *errorValue) Type() Type {
	return TypeError
}

func (errorValue *errorValue) Internal() interface{} {
	return errorValue.msg
}

func (errorValue *errorValue) Error() string {
	return errorValue.String()
}

func (errorValue *errorValue) SetAt(meta ScriptMetaData, lineno int) {
	errorValue.meta = meta
	errorValue.lineno = lineno
}

func (errorValue *errorValue) At() (meta ScriptMetaData, lineno int) {
	return errorValue.meta, errorValue.lineno
}

func (errorValue *errorValue) IsTraced() bool {
	meta, lineno := errorValue.At()
	return meta != nil && lineno > 0
}

func (errorValue *errorValue) Panic() ErrorValue {
	errorValue.fatal = true
	return errorValue
}

func (errorValue *errorValue) IsFatal() bool {
	return errorValue.fatal
}

func (errorValue *errorValue) Ignore() ErrorValue {
	errorValue.ignore = true
	errorValue.fatal = false
	return errorValue
}

func (errorValue *errorValue) CanBeIgnored() bool {
	return errorValue.ignore
}

// NewErrorValue creates a new Error
//
func NewErrorValue(msg string) ErrorValue {
	return &errorValue{baseValue: baseValue{info: typeInfoError}, msg: msg}
}
