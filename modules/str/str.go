package str

import (
	"bytes"
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
		concat(),
		trim(),
		trimLeft(),
		trimRight(),
		trimPrefix(),
		trimSuffix(),
		replaceAll(),
		replaceFirst(),
		replaceLast(),
		findAll(),
		findFirst(),
		findLast(),
		count(),
		split(),
		endsWith(),
		startsWith(),
		upper(),
		lower(),
		padLeft(),
		padRight(),
		padBoth(),
	})
}

func at() elmo.NamedValue {
	return elmo.NewGoFunctionWithHelp("at", `get character at position`, func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

		_, err := elmo.CheckArguments(arguments, 2, 3, "at", "<string> <from> <to>?")
		if err != nil {
			return err
		}

		str := elmo.EvalArgument(context, arguments[0])
		if str.Type() == elmo.TypeError {
			return str
		}
		if str.Type() == elmo.TypeString {
			return str.(elmo.Runnable).Run(context, arguments[1:])
		}

		return elmo.NewStringLiteral(str.String()).(elmo.Runnable).Run(context, arguments[1:])
	})
}

func concat() elmo.NamedValue {
	return elmo.NewGoFunctionWithHelp("concat", `concatenate/join strings`, func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

		argLen, err := elmo.CheckArguments(arguments, 2, math.MaxInt16, "concat", "<string> <string>+")
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

func applyTrim(context elmo.RunContext, arguments []elmo.Argument, trimName string, trimFunc func(string, string) string) elmo.Value {
	argLen, err := elmo.CheckArguments(arguments, 1, 2, trimName, "<string> <cutset>?")
	if err != nil {
		return err
	}

	cutset := " \t\n\r"

	if argLen == 2 {
		cutset = elmo.EvalArgument2String(context, arguments[1])
	}

	return elmo.NewStringLiteral(trimFunc(elmo.EvalArgument2String(context, arguments[0]), cutset))
}

func trim() elmo.NamedValue {
	return elmo.NewGoFunctionWithHelp("trim", `trim a string from both sides`, func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {
		return applyTrim(context, arguments, "trim", strings.Trim)
	})
}

func trimLeft() elmo.NamedValue {
	return elmo.NewGoFunctionWithHelp("trimLeft", `trim a string from the left side`, func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {
		return applyTrim(context, arguments, "trimLeft", strings.TrimLeft)
	})
}

func trimRight() elmo.NamedValue {
	return elmo.NewGoFunctionWithHelp("trimRight", `trim a string from the right side`, func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {
		return applyTrim(context, arguments, "trimRight", strings.TrimRight)
	})
}

func trimPrefix() elmo.NamedValue {
	return elmo.NewGoFunctionWithHelp("trimPrefix", `remove a prefix of a string`, func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {
		return applyTrim(context, arguments, "trimPrefix", strings.TrimPrefix)
	})
}

func trimSuffix() elmo.NamedValue {
	return elmo.NewGoFunctionWithHelp("trimSuffix", `remove a suffix of a string`, func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {
		return applyTrim(context, arguments, "trimSuffix", strings.TrimSuffix)
	})
}

func getReplaceArgs(context elmo.RunContext, arguments []elmo.Argument) (string, string, string) {
	return elmo.EvalArgument2String(context, arguments[0]),
		elmo.EvalArgument2String(context, arguments[1]),
		elmo.EvalArgument2String(context, arguments[2])
}

func replaceAll() elmo.NamedValue {
	return elmo.NewGoFunctionWithHelp("replaceAll", `replace all occurences of a given text within a string`, func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

		_, err := elmo.CheckArguments(arguments, 3, 3, "replaceAll", "<string> <old> <new>")
		if err != nil {
			return err
		}

		value, oldValue, newValue := getReplaceArgs(context, arguments)

		return elmo.NewStringLiteral(strings.Replace(value, oldValue, newValue, -1))

	})
}

func replaceFirst() elmo.NamedValue {
	return elmo.NewGoFunctionWithHelp("replaceFirst", `replace the first occurences of a given text within a string`, func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

		_, err := elmo.CheckArguments(arguments, 3, 3, "replaceFirst", "<string> <old> <new>")
		if err != nil {
			return err
		}

		value, oldValue, newValue := getReplaceArgs(context, arguments)

		return elmo.NewStringLiteral(strings.Replace(value, oldValue, newValue, 1))

	})
}

func replaceLast() elmo.NamedValue {
	return elmo.NewGoFunctionWithHelp("replaceLast", `replace the last occurences of a given text within a string`, func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

		_, err := elmo.CheckArguments(arguments, 3, 3, "replaceLast", "<string> <old> <new>")
		if err != nil {
			return err
		}

		value, oldValue, newValue := getReplaceArgs(context, arguments)

		lastIndex := strings.LastIndex(value, oldValue)
		if lastIndex < 0 {
			return elmo.NewStringLiteral(value)
		}

		var buffer bytes.Buffer
		buffer.WriteString(value[0:lastIndex])
		buffer.WriteString(newValue)
		buffer.WriteString(value[lastIndex+len(oldValue):])

		return elmo.NewStringLiteral(buffer.String())

	})
}

func getFindArgs(context elmo.RunContext, arguments []elmo.Argument) (string, string) {
	return elmo.EvalArgument2String(context, arguments[0]),
		elmo.EvalArgument2String(context, arguments[1])
}

func findFirst() elmo.NamedValue {
	return elmo.NewGoFunctionWithHelp("findFirst", `find the index of the first occurences of a given text within a string`, func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

		_, err := elmo.CheckArguments(arguments, 2, 2, "findFirst", "<string> <value>")
		if err != nil {
			return err
		}

		value, what := getFindArgs(context, arguments)

		return elmo.NewIntegerLiteral(int64(strings.Index(value, what)))

	})
}

func findLast() elmo.NamedValue {
	return elmo.NewGoFunctionWithHelp("findLast", `find the index of the last occurences of a given text within a string`, func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

		_, err := elmo.CheckArguments(arguments, 2, 2, "findLast", "<string> <value>")
		if err != nil {
			return err
		}

		value, what := getFindArgs(context, arguments)

		return elmo.NewIntegerLiteral(int64(strings.LastIndex(value, what)))

	})
}

func findAll() elmo.NamedValue {
	return elmo.NewGoFunctionWithHelp("findAll", `find all indexes of the last occurences of a given text within a string`, func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

		_, err := elmo.CheckArguments(arguments, 2, 2, "findAll", "<string> <value>")
		if err != nil {
			return err
		}

		value, what := getFindArgs(context, arguments)

		whatLen := len(what)
		result := []elmo.Value{}
		foundAt := strings.Index(value, what)
		at := 0
		for foundAt >= 0 {
			result = append(result, elmo.NewIntegerLiteral(int64(at+foundAt)))
			at = at + foundAt + whatLen
			value = value[foundAt+whatLen:]
			foundAt = strings.Index(value, what)
		}
		return elmo.NewListValue(result)
	})
}

func count() elmo.NamedValue {
	return elmo.NewGoFunctionWithHelp("count", `count the number of occurences of a value inside a string`, func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

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
	return elmo.NewGoFunctionWithHelp("split", `split a string / reverse joins`, func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

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
	return elmo.NewGoFunctionWithHelp("endsWith", `checks if string ends with value`, func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

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

		return elmo.TrueOrFalse(strings.HasSuffix(str.String(), suffix.String()))

	})
}

func startsWith() elmo.NamedValue {
	return elmo.NewGoFunctionWithHelp("startsWith", `checks if string starts with value`, func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

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

		return elmo.TrueOrFalse(strings.HasPrefix(str.String(), prefix.String()))

	})
}

func upper() elmo.NamedValue {
	return elmo.NewGoFunctionWithHelp("upper", `converts a string to all uppercase characters
	usage str.upper <string>

	example:

	string (load string)
	string.upper "upper" |eq "UPPER" |assert
	`, func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

		_, err := elmo.CheckArguments(arguments, 1, 1, "upper", "<string>")
		if err != nil {
			return err
		}

		return elmo.NewStringLiteral(strings.ToUpper(elmo.EvalArgument2String(context, arguments[0])))

	})
}

func lower() elmo.NamedValue {
	return elmo.NewGoFunctionWithHelp("lower", `converts a string to all lowercase characters
	usage str.lower <string>

	example:

	string (load string)
	string.lower "LOWER" |eq "lower" |assert
	`, func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

		_, err := elmo.CheckArguments(arguments, 1, 1, "lower", "<string>")
		if err != nil {
			return err
		}

		return elmo.NewStringLiteral(strings.ToLower(elmo.EvalArgument2String(context, arguments[0])))

	})
}

func applyPadRight(str string, l int, pad string) string {
	if l <= 0 {
		return ""
	}
	for {
		str += pad
		if len(str) > l {
			return str[0:l]
		}
	}
}

func applyPadLeft(str string, l int, pad string) string {
	if l <= 0 {
		return ""
	}
	for {
		str = pad + str
		if len(str) > l {
			return str[len(str)-l:]
		}
	}
}

func applyPadBoth(str string, l int, pad string) string {
	if l <= 0 {
		return ""
	}
	for {
		str = pad + str + pad
		if len(str) > l {
			mid := len(str) / 2
			start := mid - (l / 2)
			return str[start : start+l]
		}
	}
}

func getPadArgs(context elmo.RunContext, arguments []elmo.Argument) (string, int, string, elmo.ErrorValue) {
	value := elmo.EvalArgument2String(context, arguments[0])
	length := elmo.EvalArgument(context, arguments[1])
	var lengthAsInt int
	if length.Type() == elmo.TypeInteger {
		lengthAsInt = int(length.Internal().(int64))
	} else {
		return "", -1, "", elmo.NewErrorValue("padding expects and integer length")
	}

	padding := " "
	if len(arguments) == 3 {
		padding = elmo.EvalArgument2String(context, arguments[2])
	}

	return value, lengthAsInt, padding, nil

}

func padLeft() elmo.NamedValue {
	return elmo.NewGoFunctionWithHelp("padLeft", `add padding to the beginning of a string`, func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

		_, err := elmo.CheckArguments(arguments, 2, 3, "padLeft", "<string> <length> <padding>?")
		if err != nil {
			return err
		}

		value, length, padding, err := getPadArgs(context, arguments)
		if err != nil {
			return err
		}

		return elmo.NewStringLiteral(applyPadLeft(value, length, padding))

	})
}

func padRight() elmo.NamedValue {
	return elmo.NewGoFunctionWithHelp("padRight", `add padding to the end of a string`, func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

		_, err := elmo.CheckArguments(arguments, 2, 3, "padRight", "<string> <length> <padding>?")
		if err != nil {
			return err
		}

		value, length, padding, err := getPadArgs(context, arguments)
		if err != nil {
			return err
		}

		return elmo.NewStringLiteral(applyPadRight(value, length, padding))

	})
}

func padBoth() elmo.NamedValue {
	return elmo.NewGoFunctionWithHelp("padBoth", `add padding to both the beginning and the end of a string`, func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

		_, err := elmo.CheckArguments(arguments, 2, 3, "padBoth", "<string> <length> <padding>?")
		if err != nil {
			return err
		}

		value, length, padding, err := getPadArgs(context, arguments)
		if err != nil {
			return err
		}

		return elmo.NewStringLiteral(applyPadBoth(value, length, padding))

	})
}
