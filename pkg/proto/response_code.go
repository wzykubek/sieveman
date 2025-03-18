package proto

type Tag struct {
	Msg       string
	ChildCode ResponseCode
}

func (rc Tag) Type() string {
	return "TAG"
}

func (rc Tag) Message() string {
	return rc.Msg
}

func (rc Tag) Child() ResponseCode {
	return rc.ChildCode
}

type NonExistent struct {
	Msg       string
	ChildCode ResponseCode
}

func (rc NonExistent) Type() string {
	return "NONEXISTENT"
}

func (rc NonExistent) Message() string {
	return rc.Msg
}

func (rc NonExistent) Child() ResponseCode {
	return rc.ChildCode
}
