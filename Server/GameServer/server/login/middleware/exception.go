package middleware

import (
	"battle_rabbit/server/login/proto"
	"fmt"
	"github.com/gin-gonic/gin"
	"runtime"
)

func Exception() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover();err != nil {
				PrintStack()
				proto.ResponseFail(c)
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