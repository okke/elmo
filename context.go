package elmo

type runContext struct {
	properties map[string]Value
	parent     RunContext
}

// RunContext provides a runtime environment for script execution
//
type RunContext interface {
	Set(key string, value Value)
	SetNamed(value NamedValue)
	Get(key string) (Value, bool)
	CreateSubContext() RunContext
}

func (runContext *runContext) Set(key string, value Value) {
	runContext.properties[key] = value
}

func (runContext *runContext) SetNamed(value NamedValue) {
	runContext.Set(value.Name(), value)
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

func (runContext *runContext) CreateSubContext() RunContext {
	return NewRunContext(runContext)
}

// NewRunContext constructs a new run context
//
func NewRunContext(parent RunContext) RunContext {
	return &runContext{parent: parent, properties: make(map[string]Value)}
}
