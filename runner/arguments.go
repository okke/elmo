package runner

import (
	"fmt"
	"strings"
)

type cliArgs []string

func (args *cliArgs) next() (string, cliArgs) {
	arr := *args

	if len(arr) == 0 {
		return "", []string{}
	}
	return arr[0], arr[1:]
}

func (args *cliArgs) putBack(value string) cliArgs {
	arr := *args
	return append([]string{value}, arr...)
}

func parseFlags(args cliArgs, flags func(name, value string)) cliArgs {
	if len(args) == 0 {
		return args
	}

	arg, more := args.next()
	if arg == "" {
		return more
	}

	if arg[0] == '-' {
		pair := strings.SplitN(arg, "=", 2)
		if len(pair) == 1 {
			// -flag
			//
			flags(arg[1:], "true")
		} else {
			// -flag=value
			//
			flags(pair[0][1:], pair[1])
		}

		return parseFlags(more, flags)
	}

	return more.putBack(arg)

}

func parseArguments(parse []string, setter runnerArgsSetter) {

	// elmo <elmo-flag>* ( <elmo-file>? <user-flag>* <user-arg>* )

	fmt.Println("parse args:", parse)

	var args cliArgs = parse

	rest := parseFlags(args, setter.SetElmoFlag)
	if len(rest) == 0 {
		return
	}
	setter.SetRawUserArgs(rest)

	elmoFile, userPart := rest.next()

	if elmoFile != "" {
		fmt.Println("set elmofile:", elmoFile)
		setter.SetElmoFile(elmoFile)
	}
	if len(userPart) == 0 {
		return
	}

	userArgs := parseFlags(userPart, setter.SetUserFlag)
	setter.SetUserArgs(userArgs)

	return
}

type runnerArgs struct {
	elmoFlags   map[string]string
	elmoFile    string
	userFlags   map[string]string
	userArgs    []string
	rawUserArgs []string
}

// runnerArgsSetter captures the result of a runner's argument parsing
//
type runnerArgsSetter interface {
	SetElmoFlag(name, value string)
	SetElmoFile(name string)
	SetUserFlag(name, value string)
	SetUserArgs(args []string)
	SetRawUserArgs(args []string)
}

func (runnerArgs *runnerArgs) SetElmoFlag(name, value string) {
	runnerArgs.elmoFlags[name] = value
}
func (runnerArgs *runnerArgs) SetElmoFile(name string) {
	runnerArgs.elmoFile = name
}
func (runnerArgs *runnerArgs) SetUserFlag(name, value string) {
	runnerArgs.userFlags[name] = value
}
func (runnerArgs *runnerArgs) SetUserArgs(args []string) {
	runnerArgs.userArgs = args
}
func (runnerArgs *runnerArgs) SetRawUserArgs(args []string) {
	runnerArgs.rawUserArgs = args
}

func newRunnerArgs() *runnerArgs {
	return &runnerArgs{
		elmoFlags:   make(map[string]string, 0),
		elmoFile:    "",
		userFlags:   make(map[string]string, 0),
		userArgs:    []string{},
		rawUserArgs: []string{},
	}
}
