package runner

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"

	elmo "github.com/okke/elmo/core"
	"github.com/okke/elmo/modules/actor"
	bin "github.com/okke/elmo/modules/binary"
	"github.com/okke/elmo/modules/data"
	dict "github.com/okke/elmo/modules/dictionary"
	http "github.com/okke/elmo/modules/elmohttp"
	"github.com/okke/elmo/modules/inspect"
	"github.com/okke/elmo/modules/list"
	"github.com/okke/elmo/modules/str"
	"github.com/okke/elmo/modules/sys"

	prompt "github.com/c-bata/go-prompt"
)

type runner struct {
	context               elmo.RunContext
	history               []string
	shouldMakeSuggestions bool
	running               bool
	promptPrefix          string
}

// Runner represents the commandline usage of elmo
//
type Runner interface {
	Main()
	Repl()
	RegisterReplExit(func(elmo.RunContext, []elmo.Argument) elmo.Value)
	New(elmo.RunContext, string) Runner
	Stop()
}

// NewRunner constructs a new CommandLine
//
func NewRunner(context elmo.RunContext) Runner {
	return &runner{context: context,
		history:               make([]string, 0, 0),
		shouldMakeSuggestions: true,
		running:               true}
}

func (parent *runner) New(context elmo.RunContext, prefix string) Runner {
	return &runner{
		context:               context,
		history:               parent.history,
		shouldMakeSuggestions: parent.shouldMakeSuggestions,
		running:               true,
		promptPrefix:          prefix}
}

func (runner *runner) Stop() {
	runner.running = false
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
	context.RegisterModule(data.Module)
	context.RegisterModule(http.Module)
	context.RegisterModule(inspect.Module)

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

func (runner *runner) getCommandsForCompleter(word string) ([]string, string, elmo.DictionaryValue) {

	parts := strings.Split(word, ".")

	if len(parts) > 1 {
		identifier := elmo.NewNameSpacedIdentifier(parts[:len(parts)-1]).(elmo.IdentifierValue)
		_, dict, found := identifier.LookUp(runner.context)

		if found && dict != nil && dict.Type() == elmo.TypeDictionary {
			return dict.(elmo.DictionaryValue).Keys(), parts[len(parts)-1], dict.(elmo.DictionaryValue)
		}
	}

	return runner.context.Keys(), word, nil
}

func (runner *runner) findCommand(cmd string, inDictionary elmo.DictionaryValue) (elmo.Value, bool) {
	if inDictionary != nil {
		value, found := inDictionary.Resolve(cmd)
		return value, found
	}
	value, found := runner.context.Get(cmd)
	return value, found
}

const seperatorForCompleter = "(){}[]$,;"

const seperatorForCompletion = "(){}[]$,;."

func (runner *runner) completer(in prompt.Document) []prompt.Suggest {

	s := []prompt.Suggest{}
	if !runner.shouldMakeSuggestions {
		return s
	}

	word := strings.TrimLeft(in.GetWordBeforeCursor(), seperatorForCompleter)
	if word == "" {
		return s
	}

	commands := make([]string, 0, 0)
	possibleCommands, word, inDictionary := runner.getCommandsForCompleter(word)
	for _, cmd := range possibleCommands {
		if strings.HasPrefix(cmd, word) {
			commands = append(commands, cmd)
		}
	}

	for _, cmd := range commands {
		description := ""
		if value, found := runner.findCommand(cmd, inDictionary); found {
			if help, ok := value.(elmo.HelpValue); ok {
				description = help.Help().String()
			}
		}
		newline := strings.IndexRune(description, '\n')
		if newline > 0 {
			description = description[:newline]
		}
		s = append(s, prompt.Suggest{Text: cmd, Description: description})
	}

	return s
}

func (runner *runner) input(displayPrompt string, morePrompt string) string {
	needText := true
	in := ""

	usePrompt := displayPrompt
	for needText {
		in = in + prompt.Input(usePrompt, runner.completer,
			prompt.OptionCompletionWordSeparator(seperatorForCompletion),
			prompt.OptionTitle("elmo"),
			prompt.OptionHistory(runner.history),
			prompt.OptionPrefixTextColor(prompt.Yellow),
			prompt.OptionPreviewSuggestionTextColor(prompt.Blue),
			prompt.OptionSelectedSuggestionBGColor(prompt.LightGray),
			prompt.OptionSuggestionBGColor(prompt.DarkGray))

		if !strings.HasSuffix(strings.TrimRight(in, " \t"), "\\") {
			needText = false
		}
		usePrompt = morePrompt
	}

	return in
}

func (runner *runner) RegisterReplExit(f func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value) {
	runner.context.SetNamed(elmo.NewGoFunctionWithHelp("exit", `quit elmo`, f))
}

func (runner *runner) Repl() {

	runner.RegisterReplExit(func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {
		runner.Stop()
		return elmo.Nothing
	})

	// provide a function to change autocomplete behaviour
	//
	runner.context.SetNamed(elmo.NewGoFunctionWithHelp("autoComplete", "set auto complete on or off", func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {
		argLen, err := elmo.CheckArguments(arguments, 0, 1, "autoComplete", "<true|false>?")
		if err != nil {
			return err
		}

		if argLen == 1 {
			value := elmo.EvalArgument(context, arguments[0])
			if value == nil || value.Type() != elmo.TypeBoolean {
				return elmo.NewErrorValue("autoComplete expects a boolean value")
			}
			runner.shouldMakeSuggestions = value.Internal().(bool)
		}
		return elmo.NewBooleanLiteral(runner.shouldMakeSuggestions)
	}))

	prompt := "e>mo: "
	if runner.promptPrefix != "" {
		prompt = fmt.Sprintf("(%s) %s", runner.promptPrefix, prompt)
	}
	for runner.running {
		command := runner.input(prompt, "    : ")
		runner.history = append(runner.history, command)
		value := elmo.ParseAndRun(runner.context, command)

		if value != nil && value != elmo.Nothing {
			fmt.Printf("%v\n", value)
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
	debugPtr := flag.Bool("debug", false, "enforce debug mode")
	versionPtr := flag.Bool("version", false, "print version info and quit")
	helpPtr := flag.Bool("help", false, "print help text and quit")

	flag.Parse()

	runner.context.RegisterModule(elmo.NewModule("debug", initDebugModule(runner, *debugPtr)))

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
		runner.Repl()
	} else {
		runner.read(flag.Args()[0])

		if *replPtr {
			runner.Repl()
		}

	}

}
