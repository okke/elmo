package elmo

import (
	"fmt"
	"strings"
)

type identifier struct {
	baseValue
	value []string
}

func (identifier *identifier) String() string {
	if len(identifier.value) == 1 {
		return identifier.value[0]
	}

	return strings.Join(identifier.value, ".")
}

func (identifier *identifier) Compare(context RunContext, value Value) (int, ErrorValue) {
	return strings.Compare(identifier.String(), value.String()), nil
}

func (identifier *identifier) Type() Type {
	return TypeIdentifier
}

func (identifier *identifier) Internal() interface{} {
	return identifier.value
}

func (identifier *identifier) LookUp(context RunContext) (DictionaryValue, Value, bool) {

	result, found := context.Get(identifier.value[0])
	if !found {
		return nil, NewErrorValue(fmt.Sprintf("could not resolve %v", identifier)), false
	}

	if len(identifier.value) == 1 {
		return nil, result, true
	}

	if result.Type() != TypeDictionary {
		return nil, NewErrorValue(fmt.Sprintf("%s is not a dictionary", identifier.value[0])), false
	}

	var dict = result.(DictionaryValue)
	var lookup Value

	for _, name := range identifier.value[1:] {
		lookup, found = dict.Resolve(name)

		if found {
			if lookup.Type() != TypeDictionary {
				return dict, lookup, true
			}

			dict = lookup.(DictionaryValue)
		} else {

			return dict, NewErrorValue(fmt.Sprintf("could not resolve %v", identifier)), false

		}
	}
	return dict, lookup, true
}

func (identifier *identifier) ToBinary() BinaryValue {
	return NewBinaryValueFromInternal(typeInfoIdentifier.ID(), "", identifier.value)
}

func (identifier *identifier) Length() Value {
	return NewIntegerLiteral(int64(len(identifier.String())))
}

// NewIdentifier creates a new identifier value
//
func NewIdentifier(value string) Value {
	return &identifier{baseValue: baseValue{info: typeInfoIdentifier}, value: []string{value}}
}

// NewNameSpacedIdentifier creates a new identifier value
//
func NewNameSpacedIdentifier(value []string) Value {
	return &identifier{baseValue: baseValue{info: typeInfoIdentifier}, value: value}
}
