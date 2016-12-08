package runner

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/okke/elmo/core"
	"github.com/okke/elmo/modules/actor"
	bin "github.com/okke/elmo/modules/binary"
	"github.com/okke/elmo/modules/dictionary"
	"github.com/okke/elmo/modules/list"
	"github.com/okke/elmo/modules/str"
	"github.com/okke/elmo/modules/sys"
	"github.com/peterh/liner"
)

type runner struct {
	context elmo.RunContext
}

// Runner represents the commandline usage of elmo
//
type Runner interface {
	Main()
}

// NewRunner constructs a new CommandLine
//
func NewRunner(context elmo.RunContext) Runner {
	return &runner{context: context}
}

// NewMainContext constructs a context with all elmo's default modules
//
func NewMainContext() elmo.RunContext {
	context := elmo.NewGlobalContext()

	context.RegisterModule(str.Module)
	context.RegisterModule(list.Module)
	context.RegisterModule(dict.Module)
	context.RegisterModule(actor.Module)
	context.RegisterModule(sys.Module)
	context.RegisterModule(bin.Module)

	return context
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

func (runner *runner) createCommandLine() *liner.State {
	commandLine := liner.NewLiner()

	commandLine.SetCompleter(func(line string) (possibilities []string) {

		for cmd := range runner.context.Mapping() {
			if strings.HasPrefix(cmd, strings.ToLower(line)) {
				possibilities = append(possibilities, cmd)
			}
		}
		return
	})

	return commandLine
}

func replReadMore(commandLine *liner.State, command string) string {
	trimmed := strings.TrimSpace(command)
	if len(trimmed) == 0 {
		return trimmed
	}
	last := string(trimmed[len(trimmed)-1:])

	var inMultiLine = false
	var current = trimmed

	for strings.Index("{}()[],;`", last) != -1 || strings.Count(trimmed, "`") == 1 || inMultiLine {

		// TODO: 18okt2016 should check if character not with a string or a comment
		//
		fDepth := strings.Count(trimmed, "(") - strings.Count(trimmed, ")")
		bDepth := strings.Count(trimmed, "{") - strings.Count(trimmed, "}")
		lDepth := strings.Count(trimmed, "[") - strings.Count(trimmed, "]")

		wantMore := strings.Index(",;", last) != -1

		// poor mans multi line parsing
		//
		if inMultiLine {
			if (strings.Count(current, "`") % 2) == 1 {
				inMultiLine = false
			} else {
				wantMore = true
			}
		} else {
			if (strings.Count(current, "`") % 2) == 1 {
				inMultiLine = true
				wantMore = true
			}
		}

		if fDepth > 0 || bDepth > 0 || lDepth > 0 || wantMore {
			if next, err := commandLine.Prompt("    : " + strings.Repeat("\t", bDepth)); err == nil {
				if next == "--" {
					return trimmed
				}
				trimmed = trimmed + next
				last = string(trimmed[len(trimmed)-1:])
				current = next
			}
		} else {
			return trimmed
		}

	}

	return trimmed
}

func (runner *runner) repl() {

	commandLine := runner.createCommandLine()

	// provide an exit function so the repl can be stoppped
	// (TODO 18oct2016 is an exit hook not better?)
	//
	runner.context.SetNamed(elmo.NewGoFunction("exit", func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {
		commandLine.Close()
		os.Exit(0)
		return elmo.Nothing
	}))

	for {
		if command, err := commandLine.Prompt("e>mo: "); err == nil {

			command = replReadMore(commandLine, command)
			value := elmo.ParseAndRun(runner.context, command)

			if value != nil {
				commandLine.AppendHistory(command)
				if value != elmo.Nothing {
					fmt.Printf("%v\n", value)
				}
			}

		}
	}

}

func (runner *runner) read(source string) {
	b, err := ioutil.ReadFile(source)
	if err != nil {
		fmt.Print(err)
	}
	result := elmo.ParseAndRunWithFile(runner.context, string(b), source)
	if result.Type() == elmo.TypeError {
		fmt.Printf("error: %v\n", result)
	}
}

func (runner *runner) Main() {

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
		runner.repl()
	} else {
		runner.read(flag.Args()[0])

		if *replPtr {
			runner.repl()
		}

	}

}
