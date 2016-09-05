package elmo

type module struct {
	name   string
	loaded Value

	initializer ModuleInitializer
}

// ModuleInitializer is used to initialize a module
//
type ModuleInitializer func(RunContext) Value

// Module is a set of elmo functions
//
type Module interface {
	Name() string
	Content(context RunContext) Value
}

func (module *module) Name() string {
	return module.name
}

func (module *module) Content(context RunContext) Value {
	if module.loaded == nil {
		module.loaded = module.initializer(context)
	}
	return module.loaded
}

// NewModule creates a new module
//
func NewModule(name string, initializer ModuleInitializer) Module {
	return &module{name: name, initializer: initializer}
}

// NewMappingForModule creates a new dictionary holding a modules functions
//
func NewMappingForModule(context RunContext, namedValues []NamedValue) Value {
	mapping := make(map[string]Value)

	for _, v := range namedValues {
		mapping[v.Name()] = v
	}

	return NewDictionaryValue(nil, mapping)
}
