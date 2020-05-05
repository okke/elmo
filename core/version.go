//go:generate go run github.com/pointlander/peg elmo.peg
package elmo

import "fmt"

// versionData info, got the idea from: https://ariejan.net/2015/10/12/building-golang-cli-tools-update/
//
type versionData struct {
	Major  int64
	Minor  int64
	Patch  int64
	Branch string
	Name   string
	Commit string
}

// Version of elmo
var Version = versionData{0, 1, 0, BranchName, "Chipotle", CommitHash}

// VersionDictionary contains version data as provided to the elmo runtime
//
var VersionDictionary = NewDictionaryValue(nil, map[string]Value{
	"major":  NewIntegerLiteral(Version.Major),
	"minor":  NewIntegerLiteral(Version.Minor),
	"patch":  NewIntegerLiteral(Version.Patch),
	"branch": NewStringLiteral(Version.Branch),
	"name":   NewStringLiteral(Version.Name),
	"commit": NewStringLiteral(Version.Commit),
})

// CommitHash from git
//
var CommitHash string

// BranchName from git
//
var BranchName string

func (v versionData) String() string {
	return fmt.Sprintf("Elmo %d.%d.%d-%s \"%s\" (git commit hash %s)", v.Major, v.Minor, v.Patch, v.Branch, v.Name, v.Commit)
}

func elmoVersion() NamedValue {

	return NewGoFunctionWithHelp("elmoVersion", `returns a dictionary with info about this version of elmo`,
		func(context RunContext, arguments []Argument) Value {
			return VersionDictionary
		})
}
