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
		trim()})
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

		argLen, ok, err := elmo.CheckArguments(arguments, 1, 3, "trim", "(left|right)? <string> <cutset>?")
		if !ok {
			return err
		}

		cIdx := 1

		left := false
		right := false

		value := elmo.EvalArgument(context, arguments[0])
		if value.Type() == elmo.TypeError {
			return value
		}

		if value.Type() == elmo.TypeIdentifier {
			left = (value.String() == "left")
			if !left {
				right = (value.String() == "right")
				if !right {
					return elmo.NewErrorValue(fmt.Sprintf("trim left or trim right, not %v", value))
				}
			}

		}

		if left || right {
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
		return elmo.NewStringLiteral(strings.Trim(value.String(), cutset))

	})
}
