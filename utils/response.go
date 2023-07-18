package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"net/http"
)

const (
	StatusOK    = http.StatusOK
	StatusError = http.StatusInternalServerError

	MsgOK   = "OK"
	MsgFAIL = "FAIL"
)

type Response struct {
	Code int    `json:"code"`
	Data any    `json:"data"`
	Msg  string `json:"msg"`
}

func Render(c *gin.Context, code int, data any, msg string) {
	c.Render(http.StatusOK, render.JSON{Data: Response{code, data, msg}})
}

func OK(c *gin.Context) {
	Render(c, StatusOK, map[string]any{}, MsgOK)
}

func OkWithMessage(c *gin.Context, msg string) {
	Render(c, StatusOK, map[string]any{}, msg)
}

func OkWithData(c *gin.Context, data any, msgs ...string) {
	Render(c, StatusOK, data, Msg(MsgOK, msgs...))
}

func Fail(c *gin.Context) {
	Render(c, StatusError, map[string]any{}, MsgFAIL)
}

func FailWithMessage(c *gin.Context, msg string) {
	Render(c, StatusError, map[string]any{}, msg)
}

func FailWithData(c *gin.Context, data any, msgs ...string) {
	Render(c, StatusError, data, Msg(MsgFAIL, msgs...))
}

func Msg(msg string, msgs ...string) string {
	if len(msgs) > 0 {
		msg = msgs[0]
	}
	return msg
}
