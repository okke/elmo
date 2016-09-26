package elmoActor

import "github.com/okke/elmo/core"

// ActorModule contains functions that operate on actors
//
var Module = elmo.NewModule("actor", initModule)

func initModule(context elmo.RunContext) elmo.Value {
	return elmo.NewMappingForModule(context, []elmo.NamedValue{
		_new(),
		send(),
		receive(),
		current()})
}

const currentActorKey = "-actor"

func _new() elmo.NamedValue {
	return elmo.NewGoFunction("new", func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

		// get expects exactly 1 argument
		//
		if len(arguments) != 1 {
			return elmo.NewErrorValue("invalid call to actor.new, expect exactly 1 parameter: usage new {...}")
		}

		block := arguments[len(arguments)-1]
		if block.Type() != elmo.TypeBlock {
			return elmo.NewErrorValue("invalid call to actor.new, last parameter must be a block: usage new {...}")
		}

		// create a handle that can be used to communicate with the concurent 'actor'
		//
		actor := elmo.NewInternalValue(NewActor())

		// create a new context for the actor so the actor can set its own variables
		//
		subContext := context.CreateSubContext()

		// make the actor's handle available
		// (OV 26/9/2016 not sure yet if this is the desired method)
		//
		subContext.Set(currentActorKey, actor)

		// run block in its own context as a go routine
		// so we get concurrent execution
		//
		go block.Value().(elmo.Block).Run(subContext, elmo.NoArguments)

		return actor
	})
}

func send() elmo.NamedValue {
	return elmo.NewGoFunction("send", func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

		argLen := len(arguments)

		if argLen < 1 || argLen > 2 {
			return elmo.NewErrorValue("invalid call to actor.send, expected one or two parameters. usage: send <actor> <message>?")
		}

		// first argument of a list function can be an identifier with the name of the list
		//
		resolvedActor := elmo.EvalArgumentOrSolveIdentifier(context, arguments[0])

		if resolvedActor.Type() != elmo.TypeInternal {
			return elmo.NewErrorValue("invalid call to actor.send, expected an actor as first parameter. usage: send <actor> <message>")
		}

		actualActor := resolvedActor.Internal().(Actor)
		if argLen == 1 {
			actualActor.Send(elmo.True)
		} else {
			message := elmo.EvalArgument(context, arguments[1])
			actualActor.Send(message)
		}

		return elmo.Nothing
	})
}

func receive() elmo.NamedValue {
	return elmo.NewGoFunction("receive", func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

		if len(arguments) != 0 {
			return elmo.NewErrorValue("invalid call to actor.receive, did not expect any parameter. usage: receive")
		}

		actor, found := context.Get(currentActorKey)

		if !found {
			return elmo.NewErrorValue("invalid call to actor.receive, not in an actor context. usage: receive")
		}

		return actor.Internal().(Actor).Receive()
	})
}

func current() elmo.NamedValue {
	return elmo.NewGoFunction("current", func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {
		if len(arguments) != 0 {
			return elmo.NewErrorValue("invalid call to actor.receive, did not expect any parameter. usage: receive")
		}

		actor, found := context.Get(currentActorKey)

		if found {
			return actor
		}

		return elmo.Nothing
	})
}
