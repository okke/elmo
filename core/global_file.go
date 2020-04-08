package elmo

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
)

// file return the file command which creates a dictionary with file info and
// functions to read the content of the file.
//
func file() NamedValue {
	return NewGoFunction(`file/Returns a dictionary with file information based on a given path
		Usage: file <path>
		Returns: dictionary {
		  exists (boolean)
		  name (string)
		  path (string)
		  absPath (string)
		  mode (string)
		  isDir (boolean)
		  binary (function)
		  string (function)
		}
		when the file does not exists, only given path and exists indicator are set in returned dictionary`,

		func(context RunContext, arguments []Argument) Value {
			if _, err := CheckArguments(arguments, 1, 1, "file", "<path>?"); err != nil {
				return err
			}

			path := EvalArgument(context, arguments[0])
			if path.Type() == TypeError {
				return path
			}

			info, err := os.Stat(path.String())
			if err != nil {
				if os.IsNotExist(err) {
					return addFunctionsToFile(NewDictionaryValue(nil, map[string]Value{
						"exists": NewBooleanLiteral(false),
						"path":   path}))

				}
				return NewErrorValue(err.Error())
			}

			absPath, err := filepath.Abs(path.String())
			if err != nil {
				return NewErrorValue(err.Error())
			}

			return addFunctionsToFile(NewDictionaryValue(nil, map[string]Value{
				"exists":  NewBooleanLiteral(true),
				"name":    NewStringLiteral(info.Name()),
				"path":    NewStringLiteral(path.String()),
				"absPath": NewStringLiteral(absPath),
				"mode":    NewStringLiteral(info.Mode().String()),
				"isDir":   NewBooleanLiteral(info.IsDir())}).(DictionaryValue))
		})
}

func addFunctionsToFile(file DictionaryValue) DictionaryValue {
	file.Set(NewIdentifier("binary"), fileBinaryContent(file))
	file.Set(NewIdentifier("string"), fileStringContent(file))
	file.Set(NewIdentifier("write"), fileWrite(file))
	return file
}

func getFileContent(file DictionaryValue, transform func([]byte) Value) Value {
	isDir, found := file.Resolve("isDir")
	if found && isDir.Internal().(bool) {
		return NewErrorValue("can not read the content of a directory")
	}

	path, found := file.Resolve("path")
	content, err := ioutil.ReadFile(path.String())
	if err != nil {
		return NewErrorValue(err.Error())
	}
	return transform(content)

}

func fileBinaryContent(file DictionaryValue) NamedValue {
	return NewGoFunctionWithHelp("binary", `Returns the binary content of a file
		Usage: file.binary 
		Returns: file content as a binary value`,

		func(context RunContext, arguments []Argument) Value {
			return getFileContent(file, func(content []byte) Value {
				return NewBinaryValue(content)
			})
		})
}

func fileStringContent(file DictionaryValue) NamedValue {
	return NewGoFunctionWithHelp("string", `Returns the content of a file as string
		Usage: file.string
		Returns: file content as a string value`,

		func(context RunContext, arguments []Argument) Value {
			return getFileContent(file, func(content []byte) Value {
				return NewStringLiteral(string(content))
			})
		})
}

func fileWrite(file DictionaryValue) NamedValue {
	return NewGoFunctionWithHelp("write", `Writes content to a file
		Usage: file.write <value> 
		Returns: the file itself`,

		func(context RunContext, arguments []Argument) Value {

			path, found := file.Resolve("path")
			if !found {
				return NewErrorValue("missing path in file dictionary")
			}

			var buf bytes.Buffer

			for _, arg := range arguments {
				data := EvalArgument(context, arg)
				if data.Type() == TypeBinary {
					buf.Write(data.(BinaryValue).AsBytes())
				} else {
					buf.WriteString(data.String())
				}
			}

			if err := ioutil.WriteFile(path.String(), buf.Bytes(), 0644); err != nil {
				return NewErrorValue(err.Error())
			}

			return file
		})
}
