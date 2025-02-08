package client

type Response interface {
	Type() string
	Code() string // TODO: create types for codes
	Message() string
}

type ResponseOK struct {
	code string
	msg  string
}

func (r ResponseOK) Type() string {
	return "OK"
}

func (r ResponseOK) Code() string {
	return r.code
}

func (r ResponseOK) Message() string {
	return r.msg
}

type ResponseNO struct {
	code string
	msg  string
}

func (r ResponseNO) Type() string {
	return "NO"
}

func (r ResponseNO) Code() string {
	return r.code
}

func (r ResponseNO) Message() string {
	return r.msg
}

type ResponseBYE struct {
	code string
	msg  string
}

func (r ResponseBYE) Type() string {
	return "BYE"
}

func (r ResponseBYE) Code() string {
	return r.code
}

func (r ResponseBYE) Message() string {
	return r.msg
}

type Message struct {
	Key   string
	Value string
}
