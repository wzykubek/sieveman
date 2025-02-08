package proto

type Response interface {
	Type() string
	Code() ResponseCode
	Message() string
}

type ResponseOK struct {
	ResponseCode ResponseCode
	Msg          string
}

func (r ResponseOK) Type() string {
	return "OK"
}

func (r ResponseOK) Code() ResponseCode {
	return r.ResponseCode
}

func (r ResponseOK) Message() string {
	return r.Msg
}

type ResponseNO struct {
	ResponseCode ResponseCode
	Msg          string
}

func (r ResponseNO) Type() string {
	return "NO"
}

func (r ResponseNO) Code() ResponseCode {
	return r.ResponseCode
}

func (r ResponseNO) Message() string {
	return r.Msg
}

type ResponseBYE struct {
	ResponseCode ResponseCode
	Msg          string
}

func (r ResponseBYE) Type() string {
	return "BYE"
}

func (r ResponseBYE) Code() ResponseCode {
	return r.ResponseCode
}

func (r ResponseBYE) Message() string {
	return r.Msg
}

type Message struct {
	Key   string
	Value string
}
