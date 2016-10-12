package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/okke/elmo/core"
	"github.com/okke/elmo/modules/actor"
	"github.com/okke/elmo/modules/dictionary"
	"github.com/okke/elmo/modules/list"
	"github.com/okke/elmo/modules/sys"
)

func createMainContext() elmo.RunContext {
	context := elmo.NewGlobalContext()

	context.RegisterModule(list.Module)
	context.RegisterModule(dict.Module)
	context.RegisterModule(actor.Module)
	context.RegisterModule(sys.Module)

	// provide an exit function so the repl can be stoppped
	// (TODO 12sep2016: should it be here?)
	//
	context.SetNamed(exit())

	return context
}

var mainContext = createMainContext()

func exit() elmo.NamedValue {
	return elmo.NewGoFunction("exit", func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {
		os.Exit(0)
		return elmo.Nothing
	})
}

func help() {

	// flag.printDefault also prints test flags so
	// let's print something useful ourselves
	//
	fmt.Printf("usage: elmo <flags>? <source>?\n")
	flag.VisitAll(func(f *flag.Flag) {
		if strings.HasPrefix(f.Name, "test.") {
			return
		}
		fmt.Printf(" -%s : %s\n", f.Name, f.Usage)
	})
}

func repl() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("e>mo: ")
	for scanner.Scan() {
		fmt.Printf("%v\ne>mo: ", elmo.ParseAndRun(mainContext, scanner.Text()))
	}
}

func read(source string) {
	b, err := ioutil.ReadFile(source)
	if err != nil {
		fmt.Print(err)
	}
	result := elmo.ParseAndRunWithFile(mainContext, string(b), source)
	if result.Type() == elmo.TypeError {
		fmt.Printf("error: %v\n", result)
	}
}

func main() {

	replPtr := flag.Bool("repl", false, "enforce REPL mode, even after reading from file")
	versionPtr := flag.Bool("version", false, "print version info and quit")
	helpPtr := flag.Bool("help", false, "print help text and quit")

	flag.Parse()

	if *helpPtr {
		help()
		return
	}

	if *versionPtr {
		fmt.Printf("%v\n", elmo.Version)
		return
	}

	if flag.NArg() == 0 {
		// no source specified so running elmo as a REPL
		//
		repl()
	} else {
		read(flag.Args()[0])

		if *replPtr {
			repl()
		}

	}

}
