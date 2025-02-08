package proto

type Ok struct {
	ResponseCode ResponseCode
	Msg          string
}

func (r Ok) Type() string {
	return "OK"
}

func (r Ok) Code() ResponseCode {
	return r.ResponseCode
}

func (r Ok) Message() string {
	return r.Msg
}

type No struct {
	ResponseCode ResponseCode
	Msg          string
}

func (r No) Type() string {
	return "NO"
}

func (r No) Code() ResponseCode {
	return r.ResponseCode
}

func (r No) Message() string {
	return r.Msg
}

type Bye struct {
	ResponseCode ResponseCode
	Msg          string
}

func (r Bye) Type() string {
	return "BYE"
}

func (r Bye) Code() ResponseCode {
	return r.ResponseCode
}

func (r Bye) Message() string {
	return r.Msg
}
