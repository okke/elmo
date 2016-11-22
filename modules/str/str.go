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
		replace(),
		find(),
		count(),
		split(),
		endsWith(),
		startsWith()})
}

func _len() elmo.NamedValue {
	return elmo.NewGoFunction("len", func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

		_, err := elmo.CheckArguments(arguments, 1, 1, "len", "<string>")
		if err != nil {
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

		_, err := elmo.CheckArguments(arguments, 2, 3, "at", "<string> <from> <to>?")
		if err != nil {
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

		argLen, err := elmo.CheckArguments(arguments, 2, math.MaxInt16, "join", "<string> <string>+")
		if err != nil {
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

		argLen, err := elmo.CheckArguments(arguments, 1, 3, "trim", "(left|right|prefix|suffix)? <string> <cutset>?")
		if err != nil {
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

func allLastFirst(cmd string, value elmo.Value) (bool, bool, bool, elmo.ErrorValue) {
	if value.Type() == elmo.TypeIdentifier {
		switch value.String() {
		case "all":
			return true, false, false, nil
		case "first":
			return false, false, true, nil
		case "last":
			return false, true, false, nil
		default:
			return false, false, false, elmo.NewErrorValue(fmt.Sprintf("%s first, last or all, not %v", cmd, value))
		}
	}
	return false, false, false, nil
}

func replace() elmo.NamedValue {
	return elmo.NewGoFunction("replace", func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

		_, err := elmo.CheckArguments(arguments, 3, 4, "replace", "(first|last|all)? <string> <old> <new>")
		if err != nil {
			return err
		}

		idx := 1

		value := elmo.EvalArgument(context, arguments[0])
		if value.Type() == elmo.TypeError {
			return value
		}

		all, last, first, err := allLastFirst("replace", value)

		if err != nil {
			return err
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

func find() elmo.NamedValue {
	return elmo.NewGoFunction("find", func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

		_, err := elmo.CheckArguments(arguments, 2, 3, "find", "(first|last|all)? <string> <value>")
		if err != nil {
			return err
		}

		value := elmo.EvalArgument(context, arguments[0])
		if value.Type() == elmo.TypeError {
			return value
		}

		all, last, first, err := allLastFirst("find", value)

		if err != nil {
			return err
		}

		idx := 1
		if all || last || first {
			value = elmo.EvalArgument(context, arguments[1])
			if value.Type() == elmo.TypeError {
				return value
			}
			idx = 2
		}

		whatValue := elmo.EvalArgument(context, arguments[idx])
		if whatValue.Type() == elmo.TypeError {
			return whatValue
		}

		if last {
			return elmo.NewIntegerLiteral(int64(strings.LastIndex(value.String(), whatValue.String())))
		}

		if all {
			whatStr := whatValue.String()
			whatLen := len(whatStr)
			result := []elmo.Value{}
			findIn := value.String()
			foundAt := strings.Index(findIn, whatValue.String())
			at := 0
			for foundAt >= 0 {
				result = append(result, elmo.NewIntegerLiteral(int64(at+foundAt)))
				at = at + foundAt + whatLen
				findIn = findIn[foundAt+whatLen:]
				foundAt = strings.Index(findIn, whatStr)
			}
			return elmo.NewListValue(result)

		}

		return elmo.NewIntegerLiteral(int64(strings.Index(value.String(), whatValue.String())))
	})
}

func count() elmo.NamedValue {
	return elmo.NewGoFunction("count", func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

		_, err := elmo.CheckArguments(arguments, 2, 2, "count", "<string> <value>")
		if err != nil {
			return err
		}

		str := elmo.EvalArgument(context, arguments[0])
		if str.Type() == elmo.TypeError {
			return str
		}

		sep := elmo.EvalArgument(context, arguments[1])
		if sep.Type() == elmo.TypeError {
			return sep
		}

		return elmo.NewIntegerLiteral(int64(strings.Count(str.String(), sep.String())))

	})
}

func split() elmo.NamedValue {
	return elmo.NewGoFunction("split", func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

		argLen, err := elmo.CheckArguments(arguments, 1, 2, "split", "<string> <value>?")
		if err != nil {
			return err
		}

		str := elmo.EvalArgument(context, arguments[0])
		if str.Type() == elmo.TypeError {
			return str
		}

		splitBy := ""
		if argLen == 2 {
			by := elmo.EvalArgument(context, arguments[1])
			if by.Type() == elmo.TypeError {
				return by
			}
			splitBy = by.String()
		}

		splitted := strings.Split(str.String(), splitBy)
		values := make([]elmo.Value, len(splitted))
		for i, v := range splitted {
			values[i] = elmo.NewStringLiteral(v)
		}

		return elmo.NewListValue(values)
	})
}

func endsWith() elmo.NamedValue {
	return elmo.NewGoFunction("endsWith", func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

		_, err := elmo.CheckArguments(arguments, 2, 2, "endsWith", "<string> <value>")
		if err != nil {
			return err
		}

		str := elmo.EvalArgument(context, arguments[0])
		if str.Type() == elmo.TypeError {
			return str
		}

		suffix := elmo.EvalArgument(context, arguments[1])
		if suffix.Type() == elmo.TypeError {
			return suffix
		}

		return elmo.NewBooleanLiteral(strings.HasSuffix(str.String(), suffix.String()))

	})
}

func startsWith() elmo.NamedValue {
	return elmo.NewGoFunction("startsWith", func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

		_, err := elmo.CheckArguments(arguments, 2, 2, "startsWith", "<string> <value>")
		if err != nil {
			return err
		}

		str := elmo.EvalArgument(context, arguments[0])
		if str.Type() == elmo.TypeError {
			return str
		}

		prefix := elmo.EvalArgument(context, arguments[1])
		if prefix.Type() == elmo.TypeError {
			return prefix
		}

		return elmo.NewBooleanLiteral(strings.HasPrefix(str.String(), prefix.String()))

	})
}
