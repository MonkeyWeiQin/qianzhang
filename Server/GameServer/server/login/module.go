package login

import (
	"battle_rabbit/config"
	"battle_rabbit/iface"
	"battle_rabbit/server/login/middleware"
	"battle_rabbit/server/login/router"
	"battle_rabbit/service/log"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type Login struct {
	app        iface.IApp
	Config     *config.LoginConfig
	NodeId     string
	stopped    bool
	engine     *gin.Engine
	httpServer *http.Server
}

func (login *Login) InitRouters() {
	// 静态文件路由
	login.engine.StaticFS("/static", http.Dir(login.Config.Resource+"/static/"))
	router.InitAccountRouter(login.engine)
}

func (login *Login) InitMiddleware() {
	middleware.InitJwtAuth(login.Config.JwtSecret)
	login.engine.Use(middleware.CORSMiddleware())
	login.engine.Use(middleware.Exception())
	login.engine.Use(middleware.JwtAuth())
}

func GetLoginModule() *Login {
	return new(Login)
}

func (login *Login) GetModuleType() string {
	return "Login"
}

func (login *Login) OnInit(app iface.IApp, nodeId string) {
	conf, ok := app.GetConfig().LoginConf[nodeId]
	if !ok {
		log.Fatal("login server load config err : nodeId: %s", nodeId)
	}
	login.app = app
	login.Config = conf
	login.NodeId = nodeId
	// create engine
	gin.SetMode(gin.DebugMode)
	login.engine = gin.Default()

	login.InitRouters()
	login.InitMiddleware()

	login.httpServer = &http.Server{
		Addr:         login.Config.Host,
		Handler:      login.engine,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

}
func (login *Login) OnRun() {
	go func() {
		err := login.httpServer.ListenAndServe()
		if err != nil {
			if login.stopped {
				return
			}
			log.Fatal(err)
		}
	}()
}

func (login *Login) OnStop() {
	login.stopped = true
	_ = login.httpServer.Shutdown(context.Background())
	log.Debug("Login server Stopped !!!")
}
