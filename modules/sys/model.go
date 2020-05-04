package sys

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"

	elmo "github.com/okke/elmo/core"
)

var typeInfoCommand = elmo.NewTypeInfo("command")

type command struct {
	pipeFrom Command
	cmd      *exec.Cmd
}

type running struct {
	stdout io.ReadCloser
	stdin  io.WriteCloser
	stderr io.ReadCloser
}

// Command represents all data to execute an os command
//
type Command interface {
	elmo.Listable
	Execute() elmo.Value
	Pipe() (Running, error)
}

// Running represents io pipes of a running command
//
type Running interface {
	Stdout() io.ReadCloser
	Stdin() io.WriteCloser
	Stderr() io.ReadCloser
	CloseAll()
}

func (running *running) Stdout() io.ReadCloser {
	return running.stdout
}

func (running *running) Stdin() io.WriteCloser {
	return running.stdin
}

func (running *running) Stderr() io.ReadCloser {
	return running.stderr
}

func (running *running) CloseAll() {
	defer running.stderr.Close()
	defer running.stdin.Close()
	defer running.stdout.Close()
}

func (command *command) String() string {
	if command.pipeFrom != nil {
		return fmt.Sprintf("%v | command(%s %v)", command.pipeFrom, command.cmd.Path, command.cmd.Args)
	}
	return fmt.Sprintf("command(%s %v)", command.cmd.Path, command.cmd.Args)
}

func (command *command) Pipe() (Running, error) {
	// TODO do not forget to pipe

	stdout, err := command.cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	stdin, err := command.cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	stderr, err := command.cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	err = command.cmd.Start()
	if err != nil {
		return nil, err
	}

	// pipe in background when there is a command piping into this command
	//
	if command.pipeFrom != nil {
		go func() {
			pipe, piperr := command.pipeFrom.Pipe()
			if piperr != nil {
				return
			}
			scanner := bufio.NewScanner(pipe.Stdout())
			w := bufio.NewWriter(stdin)
			for scanner.Scan() {
				txt := scanner.Text()
				w.WriteString(txt)
				w.WriteString("\n")
			}
			w.Flush()
			stdin.Close()
		}()
	}

	return &running{stdout: stdout, stdin: stdin, stderr: stderr}, nil
}

func (command *command) Execute() elmo.Value {

	// run piped commands first
	//
	run, err := command.Pipe()
	if err != nil {
		return elmo.NewErrorValue(err.Error())
	}

	// ensure pipes are closed when ready
	//
	defer run.CloseAll()

	// read from Stdout and convert to elmo values
	//
	scanner := bufio.NewScanner(run.Stdout())
	var result = []elmo.Value{}
	for scanner.Scan() {
		result = append(result, elmo.NewStringLiteral(scanner.Text()))
	}

	// return values as list
	//
	return elmo.NewListValue(result)
}

func (command *command) List() []elmo.Value {
	return command.Execute().Internal().([]elmo.Value)
}

// NewCommand creates a new Command
//
func NewCommand(pipeFrom Command, name string, args []string) Command {
	return &command{pipeFrom: pipeFrom, cmd: exec.Command(name, args...)}
}

// NewCommandValue constructs a new command as elmo value
//
func NewCommandValue(pipeFrom Command, name string, args []string) elmo.Value {
	return elmo.NewInternalValue(typeInfoCommand, NewCommand(pipeFrom, name, args))
}
