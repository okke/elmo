package main

import "github.com/okke/elmo/runner"

func main() {
	context := runner.NewMainContext()
	runner := runner.NewRunner(context)
	runner.Main()
}
