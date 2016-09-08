package elmo

type runContext struct {
	properties map[string]Value
	this       Value
	modules    map[string]Module
	parent     RunContext
	stopped    bool
}

// RunContext provides a runtime environment for script execution
//
type RunContext interface {
	Set(key string, value Value)
	Remove(key string)
	SetNamed(value NamedValue)
	SetThis(this Value)
	This() Value
	Get(key string) (Value, bool)
	CreateSubContext() RunContext
	Parent() RunContext
	Mapping() map[string]Value
	RegisterModule(module Module)
	Module(name string) (Module, bool)
	Stop()
	isStopped() bool
}

func (runContext *runContext) Set(key string, value Value) {
	runContext.properties[key] = value
}

func (runContext *runContext) Remove(key string) {
	delete(runContext.properties, key)
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

	if runContext.parent != nil {
		return runContext.parent.Get(key)
	}

	return nil, false
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

// NewRunContext constructs a new run context
//
func NewRunContext(parent RunContext) RunContext {
	return &runContext{parent: parent, properties: make(map[string]Value), this: Nothing, modules: make(map[string]Module)}
}
