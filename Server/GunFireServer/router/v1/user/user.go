package user

import (
	"com.xv.admin.server/api/v1/public"
	"com.xv.admin.server/api/v1/user"
	"com.xv.admin.server/middleware"
	"github.com/gin-gonic/gin"
)

func InitUserRouter(gin *gin.Engine) {

	// 简单的路由组: v1
	router := gin.Group("/user").Use(middleware.JwtAuth()).Use(middleware.RouterAuth())
	{
		router.GET("/list", user.GetUsersList)              // 获得玩家列表
		router.POST("/modifyDiamond", user.ModifyDiamond)   // 修改用户钻石
		router.POST("/modifyGold", user.ModifyGold)         // 修改用户金币
		router.POST("/modifyStrength", user.ModifyStrength) //修改用户体力
		router.POST("/modifyStatus", user.ModifyStatus)     // 修改用户状态
	}


	mailRouter :=  gin.Group("/user/mail").Use(middleware.JwtAuth()).Use(middleware.RouterAuth())
	{
		mailRouter.GET("/list", public.GetPlayerMailList) //获取邮件列表
		mailRouter.POST("/del", public.DelPlayerMail) //删除邮件
	}
}
