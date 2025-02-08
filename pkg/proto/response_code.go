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
