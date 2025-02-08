package proto

type ResponseCode interface {
	Type() string
	Message() string
	// Child returns additional hierarchical response code if any. In most cases it is nil.
	Child() ResponseCode
}

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
