package list

import elmo "github.com/okke/elmo/core"

// Listable type can convert a value to a list
//
type Listable interface {
	List() []elmo.Value
}
