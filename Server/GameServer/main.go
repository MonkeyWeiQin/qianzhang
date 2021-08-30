package main

import (
	"battle_rabbit/excel"
	"battle_rabbit/server/game"
	"battle_rabbit/server/login"
	"battle_rabbit/service/log"
	"battle_rabbit/service/mgoDB"
	"battle_rabbit/service/redisDB"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

var (
	signalChan      = make(chan os.Signal, 1)
	shutdownSignals = []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGKILL}

	server = flag.String("type", "", "server type")
	id     = flag.String("id", "", "server id")
)

/*
必须有 type 和 id 两个参数,启动时之启动了其中某个服务器,比如游戏服务器 game_server.exe -type=Gate -id=Gate_001 启动游戏服务器

测试模式 debug = true 时, game login 每个模块只启动一个实例
*/

// 主程序
func main() {
	flag.Parse()
	// debug
	go func() {
		http.ListenAndServe("0.0.0.0:8155", nil)
	}()
	runtime.GOMAXPROCS(runtime.NumCPU())

	app := NewApp()
	// init Log setting
	err := log.SetLogger(app.GetConfig().Log)
	if err != nil {
		fmt.Print(err)
		return
	}
	// Init Xls config
	excel.InitXlsxConfig(app.AppBinPath)

	// init DB
	mgoDB.OnInit(app.GetConfig().MongoConf)
	redisDB.OnInit(app.GetConfig().RedisConf)

	// start modules service
	app.Run(
		game.GetGameModel(),
		login.GetLoginModule(),
	)

	signal.Notify(signalChan, shutdownSignals...)
	<-signalChan
	app.OnStop()
}
