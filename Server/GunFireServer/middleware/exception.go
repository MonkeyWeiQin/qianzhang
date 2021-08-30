package middleware

import (
	"com.xv.admin.server/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"runtime"
)

func Exception() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			err := recover()
			if err != nil {
				message := utils.GetMessage(err)
				if message !=""{
					utils.FailWithMessage(c, message)
				}else{
					PrintStack()
					utils.FailWithMessage(c, fmt.Sprintf("%s", err))
				}
			}
		}()
		c.Next()
	}
}
func PrintStack() {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	fmt.Printf("==> %s\n", string(buf[:n]))
}