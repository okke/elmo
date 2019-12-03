//go:generate go run github.com/pointlander/peg elmo.peg
package elmo

import "fmt"

// version info, got the idea from: https://ariejan.net/2015/10/12/building-golang-cli-tools-update/
//
type version struct {
	Major int
	Minor int
	Patch int
	Label string
	Name  string
}

// Version of elmo
var Version = version{0, 1, 0, "dev", "Chipotle"}

// Build info from git
var Build string

func (v version) String() string {
	if v.Label != "" {
		return fmt.Sprintf("Elmo version %d.%d.%d-%s \"%s\" (git commit hash %s)", v.Major, v.Minor, v.Patch, v.Label, v.Name, Build)
	}
	return fmt.Sprintf("Elmo version %d.%d.%d \"%s\" (git commit hash %s)", v.Major, v.Minor, v.Patch, v.Name, Build)
}
