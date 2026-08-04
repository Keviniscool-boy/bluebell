package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code ResCode     `json:"code"`
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data"`
}

func ResponseError(c *gin.Context, code ResCode) {
	rd := &Response{
		Code: code,
		Msg:  code.Msg(),
		Data: nil,
	}
	c.JSON(http.StatusOK, rd)
}
func ResponseErrorWithMsg(c *gin.Context, code ResCode, msg interface{}) {
	rd := &Response{
		Code: code,
		Msg:  msg,
		Data: nil,
	}
	c.JSON(http.StatusOK, rd)
}
func ResponseSuccess(c *gin.Context, data interface{}) {
	rd := &Response{
		Code: CodeSuccess,
		Msg:  CodeSuccess.Msg(),
		Data: data,
	}
	c.JSON(http.StatusOK, rd)
}
