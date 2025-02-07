package client

type Response interface {
	Type() string
	Message() string
}

type ResponseOK struct {
	msg string
}

func (r ResponseOK) Type() string {
	return "OK"
}

func (r ResponseOK) Message() string {
	return r.msg
}

type ResponseNO struct {
	msg string
}

func (r ResponseNO) Type() string {
	return "NO"
}

func (r ResponseNO) Message() string {
	return r.msg
}

type ResponseBYE struct {
	msg string
}

func (r ResponseBYE) Type() string {
	return "BYE"
}

func (r ResponseBYE) Message() string {
	return r.msg
}

type Message struct {
	Key   string
	Value string
}
