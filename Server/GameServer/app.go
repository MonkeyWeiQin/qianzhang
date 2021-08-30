package main

import (
	"battle_rabbit/config"
	"battle_rabbit/iface"
	"battle_rabbit/utils/file"
	"context"
	"errors"
	"fmt"
)

const (
	BinDirName     = "bin"
	CsvDirName     = "csv"
	configFileName = "server.json"
)

type App struct {
	ctx    context.Context
	cancel context.CancelFunc
	config *config.Config

	AppPath       string
	AppBinPath    string
	Stop          chan int
	Serialization iface.ISerialization
}

func NewApp() *App {
	AppPath, confPath, err := getApplicationPath()
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	app := new(App)
	app.ctx,app.cancel = context.WithCancel(ctx)
	app.AppPath = AppPath
	app.AppBinPath = file.PathConversionByOs(AppPath + "/" + BinDirName)
	app.config = config.LoadConfigFile(confPath)
	app.Stop = make(chan int)
	return app
}

func (app *App) Run(modules ...iface.Module) {
	var startModules []iface.Module
	if app.config.Debug {
		for _, module := range modules {
			module.OnInit(app,app.config.Dev[module.GetModuleType()])
			startModules = append(startModules, module)
		}
	} else {
		// 从配置命令行启动server
		for _, module := range modules {
			if module.GetModuleType() == *server {
				module.OnInit(app,*id)
				startModules = append(startModules, module)
			}
		}
	}
	if len(startModules) == 0 {
		panic("No service was started, or the corresponding service instance could not be found! ")
	}

	for _, module := range startModules {
		module.OnRun()
	}

	go func() {
		<-app.Stop
		app.cancel()
		for _, module := range startModules {
			module.OnStop()
		}
		fmt.Println("App Shutdown !!!")
	}()

}
func (app *App) OnStop() {
	app.Stop <- 1
}

// prc 客户端
//func (app *App) GetRPCClient() iface.IRPCClient {
//	return nil
//}

func (app *App) GetConfig() *config.Config {
	return app.config
}

func (app *App) GetAppPath() string {
	return app.AppPath
}

func (app *App) GetBinPath() string {
	return app.AppBinPath
}


//// 创建一个RPC客户端
//func (app *App) newRPCClient() iface.IRPCClient {
//	rpcCli := rpc.NewRPCClient(app.ctx,
//		rpc.SetCliBasePath(app.GetConfig().RPC.BasePath),
//		rpc.SetCliConsulAddr(app.GetConfig().RPC.RegisterAddr),
//		rpc.SetCliFailMode(client.Failfast),
//	)
//	rpcCli.OnStart()
//
//	app.rpcCli = rpcCli
//	return rpcCli
//}
//
//// 创建一个rpc服务端
//func (app *App) NewRPCServer(addr string) iface.IRPCServer {
//	return rpc.NewRPCServer(addr, app.GetConfig().RPC.BasePath, app.GetConfig().RPC.RegisterAddr)
//	//return rpcSvr
//}

func getApplicationPath() (string, string, error) {
	path, err := file.GetExecutePath()
	if err != nil {
		return "", "", err
	}
	conf := fmt.Sprintf("%s/%s/%s", path, BinDirName, configFileName)
	if file.IsExist(file.PathConversionByOs(conf)) {
		return path, conf, nil
	}
	path, err = file.GetWorkPath()
	if err != nil {
		return "", "", err
	}
	conf = fmt.Sprintf("%s/%s/%s", path, BinDirName, configFileName)
	if file.IsExist(file.PathConversionByOs(conf)) {
		return path, conf, nil
	}
	return "", "", errors.New("not find app path ! please check 'bin' path and 'config.json' file")
}
