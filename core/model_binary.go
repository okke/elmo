package elmo

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

// BinaryData is a struct used to serialize values to binary data
// Note, property names are as short as possible and public
// so the gob package can easily serialize it
//
type BinaryData struct {
	// Type Id for core types or -1
	//
	I int64

	// Name for non core types
	//
	N string

	// Actual data
	//
	D []byte
}

type binaryValue struct {
	baseValue
	data []byte
}

func (binaryValue *binaryValue) String() string {
	return fmt.Sprintf("binary:%v", binaryValue.data)
}

func (binaryValue *binaryValue) Type() Type {
	return TypeBinary
}

func (binaryValue *binaryValue) Internal() interface{} {
	return binaryValue.data
}

func (binaryValue *binaryValue) AsBytes() []byte {
	return binaryValue.data
}

func (binaryValue *binaryValue) ToRegular() Value {

	buf := bytes.NewBuffer(binaryValue.data)
	decoder := gob.NewDecoder(buf)
	var bdata BinaryData
	if err := decoder.Decode(&bdata); err != nil {
		return NewErrorValue(err.Error())
	}

	dataBuffer := bytes.NewBuffer(bdata.D)
	decoder = gob.NewDecoder(dataBuffer)

	switch bdata.I {
	case typeInfoIdentifier.ID():
		actualData := []string{}
		if err := decoder.Decode(&actualData); err != nil {
			return NewErrorValue(err.Error())
		}
		return NewNameSpacedIdentifier(actualData)
	case typeInfoString.ID():
		actualData := []rune{}
		if err := decoder.Decode(&actualData); err != nil {
			return NewErrorValue(err.Error())
		}
		return NewStringLiteralFromRunes(actualData)
	case typeInfoInteger.ID():
		actualData := int64(0)
		if err := decoder.Decode(&actualData); err != nil {
			return NewErrorValue(err.Error())
		}
		return NewIntegerLiteral(actualData)
	case typeInfoFloat.ID():
		actualData := float64(0)
		if err := decoder.Decode(&actualData); err != nil {
			return NewErrorValue(err.Error())
		}
		return NewFloatLiteral(actualData)
	case typeInfoBoolean.ID():
		actualData := true
		if err := decoder.Decode(&actualData); err != nil {
			return NewErrorValue(err.Error())
		}
		return TrueOrFalse(actualData)
	case typeInfoList.ID(), typeInfoDictionary.ID():
		actualData := SerializationResult{}
		if err := decoder.Decode(&actualData); err != nil {
			return NewErrorValue(err.Error())
		}
		return actualData.ToValue()
	default:
		return NewErrorValue("binaryValue.AsRegular: operation unsupported yet")
	}

}

func (binaryValue *binaryValue) Length() Value {
	return NewIntegerLiteral(int64(len(binaryValue.data)))
}

// NewBinaryValue creates a new Binary
//
func NewBinaryValue(data []byte) Value {
	return &binaryValue{baseValue: baseValue{info: typeInfoBinary}, data: data}
}

// NewBinaryValueFromInternal constructs a binary value from regular data
// id: TypeId (only for core types)
// typeName: Name of non core typeName
//
func NewBinaryValueFromInternal(id int64, typeName string, value interface{}) BinaryValue {
	var namesBuffer bytes.Buffer
	gob.NewEncoder(&namesBuffer).Encode(value)

	data := &BinaryData{id, typeName, namesBuffer.Bytes()}

	var structBuffer bytes.Buffer
	gob.NewEncoder(&structBuffer).Encode(data)

	return NewBinaryValue(structBuffer.Bytes()).(BinaryValue)
}
