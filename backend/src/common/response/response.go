package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type SuccessCode int

const (
	GetSuccess SuccessCode = 20000 + iota
	CreateSuccess
	UpdateSuccess
	DeleteSuccess
)

type FailedCode int

const (
	GetFailed FailedCode = 20010 + iota
	CreateFailed
	UpdateFailed
	DeleteFailed
)

type ErrorCode int

const (
	ValidationError ErrorCode = 20020 + iota
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Success(c *gin.Context, code SuccessCode, msg string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: int(code),
		Msg:  msg,
		Data: data,
	})
}

func Fail(c *gin.Context, code FailedCode, msg string) {
	c.JSON(http.StatusOK, Response{
		Code: int(code),
		Msg:  msg,
		Data: map[string]interface{}{},
	})
}

func Error(c *gin.Context, code ErrorCode, msg string) {
	c.JSON(http.StatusOK, Response{
		Code: int(code),
		Msg:  msg,
		Data: map[string]interface{}{},
	})
}
