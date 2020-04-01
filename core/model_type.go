package elmo

// Type represents an internal value type
//
type Type uint8

const (
	// TypeIdentifier represents a type for an identifier value
	TypeIdentifier Type = iota
	// TypeString represents a type for a string value
	TypeString
	// TypeInteger represents a type for an integer value
	TypeInteger
	// TypeFloat represents a type for a floating point value
	TypeFloat
	// TypeBoolean represents a type for a boolean value
	TypeBoolean
	// TypeList represents a type for an array value
	TypeList
	// TypeDictionary represents a type for a map value
	TypeDictionary
	// TypeError represents a type for an error value
	TypeError
	// TypeInternal represents an internal type
	TypeInternal
	// TypeBlock represents a type for a code block
	TypeBlock
	// TypeCall represent the type for a function call
	TypeCall
	// TypeGoFunction represents a type for an internal go function
	TypeGoFunction
	// TypeReturn represents a function result containing multiple values
	TypeReturn
	// TypeNil represents the type of a nil value
	TypeNil
	// TypeBinary represents the value of a byte array
	TypeBinary
)

var typeInfoIdentifier = NewTypeInfo("identifier")
var typeInfoString = NewTypeInfo("string")
var typeInfoInteger = NewTypeInfo("int")
var typeInfoFloat = NewTypeInfo("float")
var typeInfoBoolean = NewTypeInfo("bool")
var typeInfoList = NewTypeInfo("list")
var typeInfoDictionary = NewTypeInfo("dict")
var typeInfoError = NewTypeInfo("error")
var typeInfoBlock = NewTypeInfo("block")
var typeInfoCall = NewTypeInfo("call")
var typeInfoGoFunction = NewTypeInfo("func")
var typeInfoReturn = NewTypeInfo("return")
var typeInfoNil = NewTypeInfo("nil")
var typeInfoBinary = NewTypeInfo("binary")

// TypeInfo represents kinf of subType for TypeInternal values
//
type TypeInfo interface {
	ID() int64
	Name() Value
}

type typeInfo struct {
	id   int64
	name string
}

func (typeInfo *typeInfo) Name() Value {
	return NewIdentifier(typeInfo.name)
}

func (typeInfo *typeInfo) ID() int64 {
	return typeInfo.id
}

var typeCounter int64

// NewTypeInfo constructs a new type object
//
func NewTypeInfo(name string) TypeInfo {
	typeCounter = typeCounter + 1
	return &typeInfo{id: typeCounter, name: name}
}

func TypeMap(types ...Type) map[Type]bool {
	if types == nil {
		return map[Type]bool{}
	}

	mapping := make(map[Type]bool, len(types))
	for _, t := range types {
		mapping[t] = true
	}
	return mapping
}
