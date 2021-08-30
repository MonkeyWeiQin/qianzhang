package router

import (
	v1 "battle_rabbit/server/login/api/v1"
	"battle_rabbit/server/login/middleware"
	"github.com/gin-gonic/gin"
)

func InitAccountRouter(gin *gin.Engine) {


	router := gin.Group("v1").Use(middleware.JwtAuth())
	{
		router.GET("/checkToken", v1.CheckToken)
	}

	public := gin.Group("v1")
	{
		public.POST("/login_dev_id", v1.LoginByDevId)
	}
}
