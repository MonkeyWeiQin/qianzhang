package router

import (
	v1 "com.xv.admin.server/api/v1"
	"com.xv.admin.server/middleware"
	"github.com/gin-gonic/gin"
)

func InitAdminRouter(gin *gin.Engine) {
	loginRouter := gin.Group("/admin").Use(middleware.RouterAuth())
	{
		loginRouter.POST("/login", v1.AdminLogin)
	}
	router := gin.Group("/admin").Use(middleware.JwtAuth()).Use(middleware.RouterAuth())
	{
		router.GET("/info", v1.GetAdminInfo)
		router.GET("/list", v1.GetAdminList)
		router.POST("/create", v1.CreateAdmin)
		router.POST("/updatePassword",v1.UpdatePassword)
		router.POST("/updateRole",v1.UpdateRole)

	}
}
