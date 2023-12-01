package rsp

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

const (
	NO_GENERATE_CONFIG = 2000
	NO_CLIENT_CERT     = 2001

	ERROR   = 2998
	SUCCESS = 0
)

func Result(code int, msg string, data interface{}, c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		code,
		msg,
		data,
	})
}

func Ok(c *gin.Context) {
	Result(SUCCESS, "", nil, c)
}

func OkWithMessage(message string, c *gin.Context) {
	Result(SUCCESS, message, nil, c)
}

func OkWithData(data interface{}, c *gin.Context) {
	Result(SUCCESS, "", data, c)
}

func Fail(message string, c *gin.Context) {
	Result(ERROR, message, nil, c)
}

func FailWithCode(code int, message string, c *gin.Context) {
	Result(code, message, nil, c)
}
