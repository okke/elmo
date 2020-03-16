package elmo

type nothing struct {
	baseValue
}

// Nothing represents nil
//
var Nothing = &nothing{}

func (nothing *nothing) String() string {
	return "nil"
}

func (nothing *nothing) Type() Type {
	return TypeNil
}

func (nothing *nothing) Internal() interface{} {
	return nil
}
