package proto

import (
	"battle_rabbit/define"
	"github.com/gin-gonic/gin"
	"net/http"
)

type HTTPResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}


func HTTPSend(c *gin.Context, code int, data interface{}, msg string) {
	// 开始时间
	c.JSON(http.StatusOK, HTTPResponse{
		code,
		data,
		msg,
	})
}


func ResponseOk(c *gin.Context) {
	HTTPSend(c, define.MsgCode200, map[string]interface{}{}, "ok")
}
func ResponseOkWithMsg(c *gin.Context, message string) {
	HTTPSend(c, define.MsgCode200, map[string]interface{}{}, message)
}
func ResponseOkWithData(c *gin.Context, data interface{}) {
	HTTPSend(c, define.MsgCode200, data, "")
}
func ResponseFail(c *gin.Context) {
	HTTPSend(c, define.MsgCode500, map[string]interface{}{}, "")
}
func ResponseFailWithCode(c *gin.Context, code int) {
	HTTPSend(c, code, map[string]interface{}{}, "")
}
func ResponseFailWithMsg(c *gin.Context, msg string) {
	HTTPSend(c, define.MsgCode500, map[string]interface{}{}, msg)
}
func ResponseFailWithData(c *gin.Context, code int, data interface{}) {
	HTTPSend(c, code, data, "")
}

