package client

type Response interface {
	Type() string
	Code() ResponseCode
	Message() string
}

type ResponseOK struct {
	code ResponseCode
	msg  string
}

func (r ResponseOK) Type() string {
	return "OK"
}

func (r ResponseOK) Code() ResponseCode {
	return r.code
}

func (r ResponseOK) Message() string {
	return r.msg
}

type ResponseNO struct {
	code ResponseCode
	msg  string
}

func (r ResponseNO) Type() string {
	return "NO"
}

func (r ResponseNO) Code() ResponseCode {
	return r.code
}

func (r ResponseNO) Message() string {
	return r.msg
}

type ResponseBYE struct {
	code ResponseCode
	msg  string
}

func (r ResponseBYE) Type() string {
	return "BYE"
}

func (r ResponseBYE) Code() ResponseCode {
	return r.code
}

func (r ResponseBYE) Message() string {
	return r.msg
}

type Message struct {
	Key   string
	Value string
}
