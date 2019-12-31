package elmo

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// SerializationEventType decribes the type of a serialization event
// (note, serialization is a sequence of events, inspired by the SAX way of dealing with XML)
//
type SerializationEventType uint8

const (
	// SEValue denotes the serialization of a regular value
	//
	SEValue = iota

	// SEOpenDict denotes the start of a dictionary
	//
	SEOpenDict

	// SECloseDict denotes the end of a dictionary
	//
	SECloseDict

	// SEDictKey denotes a dictionary key (followed by a dictionary value)
	//
	SEDictKey

	// SEDictRef denotes a reference to a previously serialized dictionary
	//
	SEDictRef

	// SEOpenList denotes the start of a list
	//
	SEOpenList

	// SECloseList denotes the end of a list
	//
	SECloseList

	// SEListRef denotes a reference to a previously serialized list
	//
	SEListRef
)

// SerializationEvent is a struct holding event specific data
//
type SerializationEvent struct {
	// Type
	//
	T SerializationEventType

	// Key
	//
	K string

	// Value
	//
	V uuid.UUID
}

// SerializationResult is a struct holding the result of serializing an elmo value
//
type SerializationResult struct {
	// Mapping of values (uuid > binary presentation)
	//
	M map[uuid.UUID][]byte

	// List of events
	L []SerializationEvent
}

// Serialize serializes complete value graphs. It can handle both
// plain values and complex structurers composed of lists and dictionaries
//
// Serialize will capture value's UUID's so when deserialized, the constructed
// value graph will contain the same value id's.
//
// Serialize will produce a struct that can be used to produce a binary value
//
func Serialize(value Value) *SerializationResult {

	result := &SerializationResult{
		M: make(map[uuid.UUID][]byte, 0),
		L: make([]SerializationEvent, 0, 0)}

	result.addValue(value)

	return result

}

func (result *SerializationResult) addValue(value Value) error {
	switch value.Type() {
	case TypeList:
		return result.addList(value)
	case TypeDictionary:
		return result.addDict(value)
	default:
		return result.addLiteral(value)
	}
}

func (result *SerializationResult) addList(value Value) error {

	_, alreadySerialized := result.M[value.UUID()]
	if alreadySerialized {
		result.L = append(result.L, SerializationEvent{T: SEListRef, V: value.UUID()})
		return nil
	}

	// lists do not have a binary representation
	//
	result.M[value.UUID()] = []byte{}

	result.L = append(result.L, SerializationEvent{T: SEOpenList, V: value.UUID()})

	for _, innerValue := range value.Internal().([]Value) {
		result.addValue(innerValue)
	}

	result.L = append(result.L, SerializationEvent{T: SECloseList})

	return nil
}

func (result *SerializationResult) addDict(value Value) error {
	_, alreadySerialized := result.M[value.UUID()]
	if alreadySerialized {
		result.L = append(result.L, SerializationEvent{T: SEDictRef, V: value.UUID()})
		return nil
	}

	// dictionaries do not have a binary representation
	//
	result.M[value.UUID()] = []byte{}

	result.L = append(result.L, SerializationEvent{T: SEOpenDict, V: value.UUID()})

	dictValue := value.(DictionaryValue)
	for _, key := range dictValue.Keys() {
		result.L = append(result.L, SerializationEvent{T: SEDictKey, K: key})
		foundValue, _ := dictValue.Resolve(key)
		result.addValue(foundValue)
	}
	result.L = append(result.L, SerializationEvent{T: SECloseDict})

	return nil
}

func (result *SerializationResult) addLiteral(value Value) error {

	serializable, isSerializable := value.(SerializableValue)

	if !isSerializable {
		return errors.Errorf("could not serialize %v", value)
	}

	_, alreadySerialized := result.M[value.UUID()]
	if !alreadySerialized {
		result.M[value.UUID()] = serializable.ToBinary().AsBytes()
	}
	result.L = append(result.L, SerializationEvent{T: SEValue, V: value.UUID()})

	return nil
}
