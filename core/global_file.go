package elmo

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

func newFileDictionary(path string) Value {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return addFunctionsToFile(NewDictionaryValue(nil, map[string]Value{
				"exists": False,
				"path":   NewStringLiteral(path)}))

		}
		return NewErrorValue(err.Error())
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return NewErrorValue(err.Error())
	}

	return addFunctionsToFile(NewDictionaryValue(nil, map[string]Value{
		"exists":  True,
		"name":    NewStringLiteral(info.Name()),
		"path":    NewStringLiteral(path),
		"absPath": NewStringLiteral(absPath),
		"mode":    NewStringLiteral(info.Mode().String()),
		"size":    NewIntegerLiteral(info.Size()),
		"isDir":   TrueOrFalse(info.IsDir())}).(DictionaryValue))
}

// file return the file command which creates a dictionary with file info and
// functions to read the content of the file.
//
func file() NamedValue {
	return NewGoFunctionWithHelp("file", `Returns a dictionary with file information based on a given path
		Usage: file <path>
		Returns: dictionary {
		  exists (boolean)
		  name (string)
		  path (string)
		  absPath (string)
		  mode (string)
		  isDir (boolean)
		  size (integer)
		  binary (function)
		  string (function)
		  write (function)
		  append (function)
		}
		when the file does not exists, only given path and exists indicator are set in returned dictionary`,

		func(context RunContext, arguments []Argument) Value {
			if _, err := CheckArguments(arguments, 1, 1, "file", "<path>"); err != nil {
				return err
			}

			path := EvalArgument(context, arguments[0])
			if path.Type() == TypeError {
				return path
			}

			return newFileDictionary(path.String())
		})
}

func tempFile() NamedValue {
	return NewGoFunctionWithHelp("tempFile", `creates a temporary file, run some code and remove it
		Usage: tempFile <identifier> <block>
		Returns: resulting value of code execution
		
		example usage:

		tempFile tmp {
			tmp.append "some content"
			return tmp.string
		}

		another example which shows temporary files are deleted after temFile as executed the block of code:

		f: (file (tempFile tmp { return $tmp.absPath }))
		not $f.exists |assert
		`,

		func(context RunContext, arguments []Argument) Value {
			if _, err := CheckArguments(arguments, 2, 2, "tempFile", "<identifier> <code>"); err != nil {
				return err
			}

			name := EvalArgument2String(context, arguments[0])
			block := EvalArgument(context, arguments[1])
			if block.Type() != TypeBlock {
				return NewErrorValue("tmpFile expects a block of elmo code as last parameter")
			}

			tmpFile, err := ioutil.TempFile("", "elmo")
			if err != nil {
				return NewErrorValue(err.Error())
			}

			defer os.Remove(tmpFile.Name())

			file := newFileDictionary(tmpFile.Name())

			subContext := context.CreateSubContext()
			subContext.Set(name, file)

			return block.(Block).Run(subContext, []Argument{})

		})
}

func addFunctionsToFile(file DictionaryValue) DictionaryValue {
	file.Set(NewIdentifier("binary"), fileBinaryContent(file))
	file.Set(NewIdentifier("string"), fileStringContent(file))
	file.Set(NewIdentifier("write"), fileWrite(file))
	file.Set(NewIdentifier("append"), fileAppend(file))
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

			buf := EvalArguments2Buffer(context, arguments)

			if err := ioutil.WriteFile(path.String(), buf.Bytes(), 0644); err != nil {
				return NewErrorValue(err.Error())
			}

			return newFileDictionary(path.String())
		})
}

func fileAppend(file DictionaryValue) NamedValue {
	return NewGoFunctionWithHelp("append", `Append content to a file
		Usage: file.append <value> 
		Returns: the file itself`,

		func(context RunContext, arguments []Argument) Value {

			path, found := file.Resolve("path")
			if !found {
				return NewErrorValue("missing path in file dictionary")
			}

			buf := EvalArguments2Buffer(context, arguments)

			f, err := os.OpenFile(path.String(), os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				return NewErrorValue(err.Error())
			}

			defer f.Close()

			if _, err = f.Write(buf.Bytes()); err != nil {
				panic(err)
			}

			return newFileDictionary(path.String())
		})
}
