package player

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

// 获得玩家的邮件列表
func GetPlayerMailList(c *gin.Context) {
	fmt.Println("GetPlayerMailList")
}

// 设置邮件的阅读状态
func SetPlayerMailState(c *gin.Context) {
	fmt.Println("SetPlayerMailState")
}

// 删除邮件
func DelPlayerMail(c *gin.Context) {
	fmt.Println("DelPlayerMail")
}
