package elmo

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type baseValue struct {
	info  TypeInfo
	id    uuid.UUID
	mutex sync.Mutex
}

func (baseValue *baseValue) Info() TypeInfo {
	return baseValue.info
}

func (baseValue *baseValue) IsType(typeInfo TypeInfo) bool {
	if baseValue.info == nil {
		return false
	}

	return baseValue.info.ID() == typeInfo.ID()
}

func (baseValue *baseValue) Type() Type {
	fmt.Printf("check type of %v\n", baseValue)
	panic("baseValue does not support type")
}

func (baseValue *baseValue) Internal() interface{} {
	panic("baseValue does not support internal")
}

func (baseValue *baseValue) String() string {
	return "baseValue[?]"
}

func (baseValue *baseValue) UUID() uuid.UUID {
	baseValue.mutex.Lock()
	defer baseValue.mutex.Unlock()

	if baseValue.id[0] == 0 {
		baseValue.id = uuid.New()
	}

	return baseValue.id
}
