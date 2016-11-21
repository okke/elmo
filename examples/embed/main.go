package main

import "github.com/okke/elmo/runner"

func main() {
	// create a context with all elmo's default modules
	//
	context := runner.NewMainContext()

	// add our own lottle modules
	//
	context.RegisterModule(Module)

	// and run!
	//
	runner := runner.NewRunner(context)
	runner.Main()
}
