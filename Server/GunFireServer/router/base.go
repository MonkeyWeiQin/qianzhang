package router

import (
	v1 "com.xv.admin.server/api/v1"
	"com.xv.admin.server/config"
	"fmt"
	g "github.com/gin-gonic/gin"
	"net/http"
)

func InitBaseRouter(gin *g.Engine)  {
	gin.GET("/", func(context *g.Context) {
		fmt.Println(config.ENV.Server.Resource)
		// 指明html加载文件目录
		gin.LoadHTMLGlob(config.ENV.Server.Resource + "/index.html")
		context.HTML(http.StatusOK, "index.html", nil)
	})

	// 简单的路由组: v1
	router := gin.Group("base")
	{
		router.POST("/login", v1.AdminLogin)
		//router.POST("/submit", submitEndpoint)
		//router.POST("/read", readEndpoint)
	}
}


