package elmo

import "strings"

type stringLiteral struct {
	baseValue
	value  []rune
	blocks []*blockAtPositionInString
}

type blockAtPositionInString struct {
	at    int
	block Block
}

func (stringLiteral *stringLiteral) String() string {
	if stringLiteral.blocks == nil {
		return string(stringLiteral.value)
	}
	return stringLiteral.resolveBlocksToString(nil, func(Block) string { return "\\{...}" })
}

func (stringLiteral *stringLiteral) Type() Type {
	return TypeString
}

func (stringLiteral *stringLiteral) Internal() interface{} {
	return stringLiteral.value
}

func (stringLiteral *stringLiteral) index(context RunContext, argument Argument) (int, ErrorValue) {
	indexValue := EvalArgument(context, argument)

	if indexValue.Type() != TypeInteger {
		return 0, NewErrorValue("string accessor must be an integer")
	}

	i := (int)(indexValue.Internal().(int64))

	// negative index will be used to get elemnts from the end of the list
	//
	if i < 0 {
		i = len(stringLiteral.value) + i
	}

	if i < 0 || i >= len(stringLiteral.value) {
		return 0, NewErrorValue("string accessor out of bounds")
	}

	return i, nil
}

func (stringLiteral *stringLiteral) Run(context RunContext, arguments []Argument) Value {

	arglen := len(arguments)

	if arglen == 1 {
		i, err := stringLiteral.index(context, arguments[0])

		if err != nil {
			return err
		}

		return NewStringLiteralFromRunes(stringLiteral.value[i : i+1])
	}

	if arglen == 2 {
		i1, err := stringLiteral.index(context, arguments[0])
		if err != nil {
			return err
		}
		i2, err := stringLiteral.index(context, arguments[1])
		if err != nil {
			return err
		}

		if i1 > i2 {
			// return a reversed version of the sub list

			sub := stringLiteral.value[i2 : i1+1]
			runes := []rune(sub)
			for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
				runes[i], runes[j] = runes[j], runes[i]
			}
			return NewStringLiteral(string(runes))
		}

		return NewStringLiteralFromRunes(stringLiteral.value[i1 : i2+1])

	}

	return NewErrorValue("too many arguments for string access")
}

func (stringLiteral *stringLiteral) ToBinary() BinaryValue {
	return NewBinaryValueFromInternal(typeInfoString.ID(), "", stringLiteral.value)
}

func (stringLiteral *stringLiteral) Compare(context RunContext, value Value) (int, ErrorValue) {
	return strings.Compare(stringLiteral.String(), value.String()), nil
}

func (stringLiteral *stringLiteral) resolveBlocksToString(context RunContext, withBlock func(Block) string) string {

	var sb strings.Builder
	value := stringLiteral.value
	at := 0

	for _, blockPosition := range stringLiteral.blocks {
		if at < blockPosition.at {
			sb.WriteString(string(value[at:blockPosition.at]))
		}

		sb.WriteString(withBlock(blockPosition.block))

		at = blockPosition.at
	}

	if at < len(value) {
		sb.WriteString(string(value[at:]))
	}

	return sb.String()
}

func (stringLiteral *stringLiteral) ResolveBlocks(context RunContext) Value {

	if stringLiteral.blocks == nil {
		return stringLiteral
	}

	return NewStringLiteral(stringLiteral.resolveBlocksToString(context, func(block Block) string {
		if insertValue := block.Run(context, []Argument{}); insertValue != nil && insertValue != Nothing {
			return insertValue.String()
		}

		return ""
	}))
}

func (stringLiteral *stringLiteral) Length() Value {
	return NewIntegerLiteral(int64(len(stringLiteral.value)))
}

// NewStringLiteral creates a new string literal value
//
func NewStringLiteral(value string) Value {
	return &stringLiteral{baseValue: baseValue{info: typeInfoString}, value: []rune(value)}
}

// NewStringLiteralFromRunes creates a new string literal value
//
func NewStringLiteralFromRunes(value []rune) Value {
	return &stringLiteral{baseValue: baseValue{info: typeInfoString}, value: value}
}

// newStringLiteralWithBlocks creates a new string literal value and registers
// at which positions in the string dynamic content must be added
//
func newStringLiteralWithBlocks(value string, blocks []*blockAtPositionInString) Value {
	return &stringLiteral{baseValue: baseValue{info: typeInfoString}, value: []rune(value), blocks: blocks}
}
