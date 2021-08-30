package router

import (
	"com.xv.admin.server/game/player"
	"com.xv.admin.server/game/task"
	"com.xv.admin.server/middleware"
	"github.com/gin-gonic/gin"
)

func InitPlayerRouter(gin *gin.Engine) {

	//
	router := gin.Group("players").
		Use(middleware.JwtAuth())
	{
		//--------------------------------任务相关操作--------------------------------
		router.GET("/player_maintasklist", task.GetPlayerMainTaskList)
		router.GET("/player_daytasklist", task.GetPlayerDayTaskList)


		//--------------------------------玩家相关操作--------------------------------
		router.POST("/player_login", player.PlayerLogin)
		router.POST("/player_register", player.PlayerRegister)

		router.GET("/player_baseinfo", player.GetPlayerBaseInfo)
		router.GET("/player_maillist", player.GetPlayerMailList)
	}
}
