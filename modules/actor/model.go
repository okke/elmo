package actor

import (
	"fmt"

	"github.com/okke/elmo/core"
)

var typeInfoActor = elmo.NewTypeInfo("actor")

type actor struct {
	channel chan elmo.Value
}

// Actor interacts with elmo actors
//
type Actor interface {
	Send(elmo.Value)
	Receive() elmo.Value
}

func (actor *actor) Send(value elmo.Value) {
	if freezable, ok := value.(elmo.FreezableValue); ok {
		freezable.Freeze()
	}
	actor.channel <- value
}

func (actor *actor) Receive() elmo.Value {
	return <-actor.channel
}

func (actor *actor) String() string {
	return fmt.Sprintf("actor(%p)", actor)
}

// NewActor constructs a new actor
//
func NewActor() Actor {
	return &actor{channel: make(chan elmo.Value)}
}
