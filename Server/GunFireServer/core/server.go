package core

import (
	"com.xv.admin.server/config"
	"com.xv.admin.server/global"
	"com.xv.admin.server/middleware"
	"com.xv.admin.server/router"
	"com.xv.admin.server/router/v1/user"
	"com.xv.admin.server/service/mgoDB"
	"com.xv.admin.server/service/redisDB"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Start(engine *gin.Engine) *gin.Engine {
	// 全局使用使用跨域中间件，方便测试
	engine.Use(middleware.CORSMiddleware())
	engine.Use(middleware.Exception())

	// 加载路由
	InitRouters(engine)

	// mongo连接初始化
	mgoDB.NewClient(nil)
	// redis 连接初始化
	redisDB.ConnectRedis()
	fmt.Println(config.ENV.AppPath)
	global.InitXlsxConfig(config.ENV.AppPath + "/config")

	return engine
}

func InitRouters(gin *gin.Engine) {
	// 静态文件路由
	gin.StaticFS("/static", http.Dir(config.ENV.Server.Resource+"/static/"))

	router.InitAdminRouter(gin)
	router.InitSystemRouter(gin)
	user.InitUserRouter(gin)


	//router.InitBaseRouter(gin)
	//router.InitPlayerRouter(gin)
	//router.InitAdminRouter(gin)
	//router.InitSystemRouter(gin)
	//player.InitAccountRouter(gin)
	//public.InitAdsVideoConfig(gin)
	//activity.InitActivityRouter(gin)
}
