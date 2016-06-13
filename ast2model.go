package elmo

// Ast2Block converts an ast node to a code block
//
func Ast2Block(node *node32, buf string) Block {

	calls := []Call{}

	for _, call := range Children(node) {
		calls = append(calls, Ast2Call(call, buf))
	}

	return NewBlock(calls)
}

// Ast2Call converts an ast node to a function call
//
func Ast2Call(node *node32, buf string) Call {

	functionName := ""
	arguments := []Argument{}

	for idx, argument := range Children(node) {
		if idx == 0 {
			functionName = buf[argument.begin : argument.begin+(argument.end-argument.begin)]
		} else {
			arguments = append(arguments, Ast2Argument(argument, buf))
		}
	}

	return NewCall(functionName, arguments)
}

// Ast2Argument converts an ast node to a function argument
//
func Ast2Argument(node *node32, buf string) Argument {
	return nil
}
