package runner

import (
	"fmt"
	"io/ioutil"
	"os"
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
	arguments             *runnerArgs
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
		running:               true,
		arguments:             newRunnerArgs(),
	}
}

func (parent *runner) New(context elmo.RunContext, prefix string) Runner {
	return &runner{
		context:               context,
		history:               parent.history,
		shouldMakeSuggestions: parent.shouldMakeSuggestions,
		running:               true,
		promptPrefix:          prefix,
		arguments:             newRunnerArgs(),
	}
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
		return elmo.TrueOrFalse(runner.shouldMakeSuggestions)
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

func (runner *runner) storeArgs() {

	flagValues := make(map[string]elmo.Value)
	for k, v := range runner.arguments.userFlags {
		flagValues[k] = elmo.ConvertAnyToValue(v)
	}

	runner.context.Set("args", elmo.NewDictionaryValue(nil, map[string]elmo.Value{
		"script": elmo.NewStringLiteral(runner.arguments.elmoFile),
		"flags":  elmo.NewDictionaryValue(nil, flagValues),
		"args":   elmo.NewListValueFromStrings(runner.arguments.userArgs),
		"raw":    elmo.NewListValueFromStrings(runner.arguments.rawUserArgs),
	}))

}

const debugFlag = "debug"
const autoreloadFlag = "autoreload"
const replFlag = "repl"
const versionFlag = "version"
const helpFlag = "help"

func help() {

	fmt.Println("usage: elmo <flags>? <script-file>? <script-flags>? <script-args>?")
	fmt.Printf("  %-15v  start elmo in debug mode\n", "-"+debugFlag)
	fmt.Printf("  %-15v  start elmo in auto reload mode\n", "-"+autoreloadFlag)
	fmt.Printf("  %-15v  open repl after script execution\n", "-"+replFlag)
	fmt.Printf("  %-15v  print elmo's version\n", "-"+versionFlag)
}

// Main starts the elmo runtime. Either in repl mode or by interpreting an elmo source file
//
func (runner *runner) Main() {

	parseArguments(os.Args[1:], runner.arguments)
	runner.storeArgs()

	if _, needHelp := runner.arguments.elmoFlags[helpFlag]; needHelp {
		help()
		return
	}

	if _, wantsVersion := runner.arguments.elmoFlags[versionFlag]; wantsVersion {
		fmt.Printf("%v\n", elmo.Version)
		return
	}

	_, elmo.GlobalSettings().Debug = runner.arguments.elmoFlags[debugFlag]
	_, elmo.GlobalSettings().HotReload = runner.arguments.elmoFlags[autoreloadFlag]
	_, elmo.GlobalSettings().StartRepl = runner.arguments.elmoFlags[replFlag]

	runner.context.RegisterModule(elmo.NewModule("debug", initDebugModule(runner, elmo.GlobalSettings().Debug)))

	if runner.arguments.elmoFile == "" {
		// no source specified so running elmo as a REPL
		//
		runner.Repl()
	} else {

		runner.read(runner.arguments.elmoFile)

		if elmo.GlobalSettings().StartRepl {
			runner.Repl()
		}
	}
}
