package elmo

type scriptMetaData struct {
	name    string
	content []rune
}

// ScriptMetaData contains accessor function for a script's meta DecimalConstant
//
type ScriptMetaData interface {
	Name() string
	Content() []rune
	PositionOf(absolutePosition int) (int, int)
}

func (scriptMetaData *scriptMetaData) Name() string {
	return scriptMetaData.name
}

func (scriptMetaData *scriptMetaData) Content() []rune {
	return scriptMetaData.content
}

func (scriptMetaData *scriptMetaData) PositionOf(absolutePosition int) (int, int) {
	found := translatePositions([]rune(scriptMetaData.content), []int{absolutePosition})
	return found[absolutePosition].line, found[absolutePosition].symbol
}

// NewScriptMetaData constructs a meta data object for scripts
//
func NewScriptMetaData(name string, content string) ScriptMetaData {
	return &scriptMetaData{name: name, content: []rune(content)}
}
