package str

import (
	"bytes"
	"fmt"
	"math"
	"strings"

	elmo "github.com/okke/elmo/core"
)

// Module contains functions that operate on lists
//
var Module = elmo.NewModule("string", initModule)

func initModule(context elmo.RunContext) elmo.Value {
	return elmo.NewMappingForModule(context, []elmo.NamedValue{
		at(),
		_len(),
		concat(),
		trim(),
		replace()})
}

func _len() elmo.NamedValue {
	return elmo.NewGoFunction("len", func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

		_, ok, err := elmo.CheckArguments(arguments, 1, 1, "len", "<string>")
		if !ok {
			return err
		}

		str := elmo.EvalArgument(context, arguments[0])
		if str.Type() == elmo.TypeError {
			return str
		}

		return elmo.NewIntegerLiteral(int64(len(str.String())))
	})
}

func at() elmo.NamedValue {
	return elmo.NewGoFunction("at", func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

		_, ok, err := elmo.CheckArguments(arguments, 2, 3, "at", "<string> <from> <to>?")
		if !ok {
			return err
		}

		str := elmo.EvalArgument(context, arguments[0])
		if str.Type() == elmo.TypeError {
			return str
		}

		return elmo.NewStringLiteral(str.String()).(elmo.Runnable).Run(context, arguments[1:])
	})
}

func concat() elmo.NamedValue {
	return elmo.NewGoFunction("concat", func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

		argLen, ok, err := elmo.CheckArguments(arguments, 2, math.MaxInt16, "join", "<string> <string>+")
		if !ok {
			return err
		}

		var buffer bytes.Buffer

		for i := 0; i < argLen; i++ {
			value := elmo.EvalArgument(context, arguments[i])
			if value.Type() == elmo.TypeError {
				return value
			}
			buffer.WriteString(value.String())
		}

		return elmo.NewStringLiteral(buffer.String())

	})
}

func trim() elmo.NamedValue {
	return elmo.NewGoFunction("trim", func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

		argLen, ok, err := elmo.CheckArguments(arguments, 1, 3, "trim", "(left|right|prefix|suffix)? <string> <cutset>?")
		if !ok {
			return err
		}

		cIdx := 1

		left := false
		right := false
		suffix := false
		prefix := false

		value := elmo.EvalArgument(context, arguments[0])
		if value.Type() == elmo.TypeError {
			return value
		}

		if value.Type() == elmo.TypeIdentifier {
			switch value.String() {
			case "left":
				left = true
			case "right":
				right = true
			case "suffix":
				suffix = true
			case "prefix":
				prefix = true
			default:
				return elmo.NewErrorValue(fmt.Sprintf("trim left, right, prefix or suffix, not %v", value))
			}
		}

		if left || right || prefix || suffix {
			value = elmo.EvalArgument(context, arguments[1])
			if value.Type() == elmo.TypeError {
				return value
			}
			cIdx = 2
		}

		cutset := " \t"

		if (cIdx == 1 && argLen == 2) || (cIdx == 2 && argLen == 3) {
			csValue := elmo.EvalArgument(context, arguments[cIdx])
			if csValue.Type() == elmo.TypeError {
				return csValue
			}
			cutset = csValue.String()
		}

		if left {
			return elmo.NewStringLiteral(strings.TrimLeft(value.String(), cutset))
		}

		if right {
			return elmo.NewStringLiteral(strings.TrimRight(value.String(), cutset))
		}

		if prefix {
			return elmo.NewStringLiteral(strings.TrimPrefix(value.String(), cutset))
		}

		if suffix {
			return elmo.NewStringLiteral(strings.TrimSuffix(value.String(), cutset))
		}

		return elmo.NewStringLiteral(strings.Trim(value.String(), cutset))

	})
}

func replace() elmo.NamedValue {
	return elmo.NewGoFunction("replace", func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

		_, ok, err := elmo.CheckArguments(arguments, 3, 4, "replace", "(first|last|all)? <string> <old> <new>")
		if !ok {
			return err
		}

		idx := 1

		value := elmo.EvalArgument(context, arguments[0])
		if value.Type() == elmo.TypeError {
			return value
		}

		all := false
		last := false
		first := false

		if value.Type() == elmo.TypeIdentifier {
			switch value.String() {
			case "all":
				all = true
			case "first":
				first = true
			case "last":
				last = true
			default:
				return elmo.NewErrorValue(fmt.Sprintf("replace first, last or all, not %v", value))
			}
		}

		if all || last || first {
			value = elmo.EvalArgument(context, arguments[1])
			if value.Type() == elmo.TypeError {
				return value
			}
			idx = 2
		}

		oldValue := elmo.EvalArgument(context, arguments[idx])
		if oldValue.Type() == elmo.TypeError {
			return oldValue
		}

		newValue := elmo.EvalArgument(context, arguments[idx+1])
		if newValue.Type() == elmo.TypeError {
			return newValue
		}

		if all {
			return elmo.NewStringLiteral(strings.Replace(value.String(), oldValue.String(), newValue.String(), -1))
		}

		if last {
			lastIndex := strings.LastIndex(value.String(), oldValue.String())
			if lastIndex < 0 {
				return value
			}

			var buffer bytes.Buffer
			buffer.WriteString(value.String()[0:lastIndex])
			buffer.WriteString(newValue.String())
			buffer.WriteString(value.String()[lastIndex+len(oldValue.String()):])

			return elmo.NewStringLiteral(buffer.String())
		}

		return elmo.NewStringLiteral(strings.Replace(value.String(), oldValue.String(), newValue.String(), 1))

	})
}
