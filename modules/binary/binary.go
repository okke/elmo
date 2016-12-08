package bin

import (
	"fmt"

	elmo "github.com/okke/elmo/core"
)

// Module contains functions that operate on binary data
//
var Module = elmo.NewModule("bin", initModule)

func initModule(context elmo.RunContext) elmo.Value {
	return elmo.NewMappingForModule(context, []elmo.NamedValue{
		_new()})
}

func _new() elmo.NamedValue {
	return elmo.NewGoFunction(`new/converts a regular elmo value to a binary representation
    Usage: new <value>
    Returns: binary value
    `,
		func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {
			_, err := elmo.CheckArguments(arguments, 1, 1, "new", "<value>")
			if err != nil {
				return err
			}

			value := elmo.EvalArgument(context, arguments[0])
			serializable, ok := value.(elmo.SerializableValue)
			if !ok {
				return elmo.NewErrorValue(fmt.Sprintf("new expects serializable value, not %v", value))
			}

			result, ok := serializable.ToBinary().(elmo.Value)
			if !ok {
				return elmo.NewErrorValue(fmt.Sprintf("could not create an elmo value from serialized value %v", value))
			}

			return result

		})
}
