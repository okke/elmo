package elmo

import "fmt"

type runContext struct {
	properties map[string]Value
	this       Value
	scriptName Value
	modules    map[string]Module
	parent     RunContext
	joined     RunContext
	stopped    bool
}

// RunContext provides a runtime environment for script execution
//
type RunContext interface {
	Set(key string, value Value)
	Remove(key string)
	Mixin(value Value) Value
	SetNamed(value NamedValue)
	SetThis(this Value)
	This() Value
	SetScriptName(this Value)
	ScriptName() Value
	Get(key string) (Value, bool)
	Keys() []string
	CreateSubContext() RunContext
	Parent() RunContext
	Mapping() map[string]Value
	RegisterModule(module Module)
	Module(name string) (Module, bool)
	Stop()
	isStopped() bool
	Join(with RunContext) RunContext
}

func (runContext *runContext) Set(key string, value Value) {
	runContext.properties[key] = value
}

func (runContext *runContext) Remove(key string) {
	delete(runContext.properties, key)
}

func (runContext *runContext) Mixin(value Value) Value {
	if value.Type() != TypeDictionary {
		return NewErrorValue(fmt.Sprintf("mixin can only mix in dictionaries, not %s", value.String()))
	}

	for k, v := range value.Internal().(map[string]Value) {
		runContext.Set(k, v)
	}

	return value
}

func (runContext *runContext) SetNamed(value NamedValue) {
	runContext.Set(value.Name(), value)
}

func (runContext *runContext) This() Value {
	return runContext.this
}

func (runContext *runContext) SetThis(this Value) {
	runContext.this = this
}

func (runContext *runContext) ScriptName() Value {
	name := runContext.scriptName
	if name != nil {
		return name
	}
	if runContext.parent != nil {
		return runContext.parent.ScriptName()
	}
	return nil
}

func (runContext *runContext) SetScriptName(scriptName Value) {
	runContext.scriptName = scriptName
}

func (runContext *runContext) RegisterModule(module Module) {
	runContext.modules[module.Name()] = module
}

func (runContext *runContext) Module(name string) (Module, bool) {

	value, found := runContext.modules[name]

	if found {
		return value, true
	}

	if runContext.parent != nil {
		return runContext.parent.Module(name)
	}

	return nil, false
}

func (runContext *runContext) Get(key string) (Value, bool) {

	value, found := runContext.properties[key]

	if found {
		return value, true
	}

	if runContext.joined != nil {
		value, found := runContext.joined.Get(key)
		if found {
			return value, true
		}
	}

	if runContext.parent != nil {
		return runContext.parent.Get(key)
	}

	return nil, false
}

func (runContext *runContext) Keys() []string {
	keys := []string{}
	for k := range runContext.properties {
		keys = append(keys, k)
	}
	if runContext.parent != nil {
		keys = append(keys, runContext.parent.Keys()...)
	}
	if runContext.joined != nil {
		keys = append(keys, runContext.joined.Keys()...)
	}

	return keys
}

func (runContext *runContext) Parent() RunContext {
	return runContext.parent
}

func (runContext *runContext) CreateSubContext() RunContext {
	return NewRunContext(runContext)
}

func (runContext *runContext) Mapping() map[string]Value {
	return runContext.properties
}

func (runContext *runContext) Stop() {
	runContext.stopped = true
}

func (runContext *runContext) isStopped() bool {
	return runContext.stopped
}

func (rc *runContext) Join(with RunContext) RunContext {
	copy := &runContext{parent: rc.parent, properties: rc.properties, this: rc.this, scriptName: rc.scriptName, modules: rc.modules}
	copy.joined = with
	return copy
}

// NewRunContext constructs a new run context
//
func NewRunContext(parent RunContext) RunContext {
	return &runContext{parent: parent, properties: make(map[string]Value), this: Nothing, scriptName: nil, modules: make(map[string]Module)}
}
