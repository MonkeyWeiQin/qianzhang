package activity

import (
	"github.com/gin-gonic/gin"
)

func InitActivityRouter(gin *gin.Engine) {
	//sevenDayRouter := gin.Group("/v1/activity/sevenDay").Use(middleware.JwtAuth())
	//{
	//	sevenDayRouter.POST("/", activity.SetSevenDayRule)  //设置活动规则
	//	sevenDayRouter.GET("/", activity.GetSevenDay)		//用户获取活动内容
	//	sevenDayRouter.POST("/signIn", activity.SevenDaySignIn)  //用户参加活动
	//}
}