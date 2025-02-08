package client

type ResponseCode interface {
	Type() string
	Message() string
	// Child returns additional hierarchical response code if any. In most cases it is nil.
	Child() ResponseCode
}

type Tag struct {
	msg   string
	child ResponseCode
}

func (rc Tag) Type() string {
	return "TAG"
}

func (rc Tag) Message() string {
	return rc.msg
}

func (rc Tag) Child() ResponseCode {
	return rc.child
}
