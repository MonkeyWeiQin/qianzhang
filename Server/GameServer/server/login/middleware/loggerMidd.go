package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

/*
	请求日志中间件
 */
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {

		// 设置 example 变量

		// 请求前
		fmt.Println("==========vvvvvvvvvvvvvvvvvvvvvvv==================")
		c.Next()

		// 请求后
		_,err := c.Writer.Write([]byte("1111111111111111111111111111111111"))
		if err != nil {
			fmt.Println("==========ddddddddddddd==================",err)
		}
		// 获取发送的 status
		status := c.Writer.Status()
		log.Println(status)
	}
}
