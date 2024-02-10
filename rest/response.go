package rest

type Response struct {
	Status int    `json:"status"`
	Data   any    `json:"data"`
	Msg    string `json:"msg"`
}

var (
	StatusOK    = 0
	StatusError = 100

	MsgOK     = "OK"
	MsgFAIL   = "FAIL"
	EmptyData = map[string]any{}
)

func OK() Response {
	return Response{StatusOK, EmptyData, MsgOK}
}

func OkWithMessage(msg string) Response {
	return Response{StatusOK, EmptyData, msg}
}

func OkWithData(data any, msgs ...string) Response {
	return Response{StatusOK, data, Msg(MsgOK, msgs...)}
}

func Fail() Response {
	return Response{StatusError, EmptyData, MsgFAIL}
}

func FailWithMessage(msg string) Response {
	return Response{StatusError, EmptyData, msg}
}

func FailWithData(data any, msgs ...string) Response {
	return Response{StatusError, data, Msg(MsgFAIL, msgs...)}
}

func Msg(msg string, msgs ...string) string {
	if len(msgs) > 0 {
		msg = msgs[0]
	}
	return msg
}
