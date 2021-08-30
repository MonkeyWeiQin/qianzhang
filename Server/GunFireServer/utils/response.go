package utils

import (
	"com.xv.admin.server/config"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

const (
	UNAUTHORIZED = 401 // token 过期或者失效
	SUCCESS      = 200 // 请求成功
	ERROR        = 500
)

type Pagination struct {
	Data      interface{} `json:"data"`
	PageTotal int64       `json:"page_total"`
	PageSize  int64       `json:"page_size"`
	Current   int64       `json:"current"`
	Count     int64       `json:"count"`
}

func Send(c *gin.Context, code int, data interface{}, msg string) {
	// 开始时间
	c.JSON(http.StatusOK, Response{
		code,
		data,
		msg,
	})
}

func Ok(c *gin.Context) {
	Send(c, SUCCESS, map[string]interface{}{}, GetMessage(00000))
}

func OkWithMessage(c *gin.Context, message string) {
	Send(c, SUCCESS, map[string]interface{}{}, message)
}

func OkWithData(c *gin.Context, data interface{}) {
	Send(c, SUCCESS, data, GetMessage(00000))
}
func OkWithPagination(c *gin.Context, data interface{}, count int64, pageSize int64, current int64) {
	Pagination := new(Pagination)
	Pagination.Data = data
	Pagination.Current = current
	Pagination.Count = count
	Pagination.PageSize = pageSize
	Pagination.PageTotal = int64(math.Ceil(float64(count) / float64(pageSize)))
	Send(c, SUCCESS, Pagination, GetMessage(00000))
}

func OkDetailed(c *gin.Context, data interface{}, message string) {
	Send(c, SUCCESS, data, message)
}

func Fail(c *gin.Context) {
	Send(c, ERROR, map[string]interface{}{}, GetMessage(10002))
}

func FailWithMessage(c *gin.Context, message string) {
	Send(c, ERROR, nil, message)
}

func FailWithDetailed(c *gin.Context, code int, data interface{}, message string) {
	Send(c, code, data, message)
}

func GetMessage(temp interface{}) string {
	message := ""
	if config.ENV.Server.Language != "en_us" {
		message = language(&En{}).strategy.message(temp)
	} else {
		message = language(&Zh{}).strategy.message(temp)
	}
	return message
}
