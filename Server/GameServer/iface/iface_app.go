package iface

import (
	"battle_rabbit/config"
	"io"
	"net"
)

type (

	IServer interface {
		Init()
		//启动服务器方法
		Start()
		// 停止服务器方法
		Stop()
		SetNewAgentFn(fn func(c net.Conn) IAgent)
	}

	IConnManager interface {
		Add(agt IAgent)    //添加链接
		Remove(string)     //删除连接
		Get(string) IAgent //利用ConnID获取链接
		Len() int          //获取当前连接
		ClearConn()        //删除并停止所有链接
		Range(fn func(agt IAgent))
	}


	IMessage interface {
		GetDataLen() uint32
		GetMsgId() uint32
		GetData() []byte

		SetMsgId(uint32)
		SetData([]byte)

		WriterPack(io.Writer) error
		ReadPack(io.Reader) error
	}

	IEncryption interface {
		Encode(data []byte) ([]byte, error)
		Decode(data []byte) ([]byte, error)
	}

	// IAgent 客户端代理定义
	IAgent interface {
		Run() error
		OnClose() error
		ReadLoop()
		WriterLoop()
		SendMessage(IMessage)
		GetSession() ISession
		GetUid() int
		Stopped() bool
	}

	ISession interface {
		SetNodeId(string) error
		GetNodeId() string
		GetSessionId() string
		SetSessionId(string) error
		GetUid() int
		Bind(int,string) error
		Close() error
		IsConnect() bool
		Send(IMessage)
		GetUserCacheAll() map[string]interface{}
		GetUserCacheByKey(k string) interface{}
		SetUserCacheByKV(string, interface{}) error
		SetUserCacheByMap(map[string]interface{}) error
	}

	Module interface {
		GetModuleType() string
		OnInit(app IApp, nodeId string)
		OnRun()
		OnStop()
	}

	IRegister interface {
		GetFunctions() interface{}
		GetServerPath() string
	}

	// SessionLearner 客户端代理
	SessionLearner interface {
		Connect(a ISession)    //当连接建立  tcp协议握手成功
		DisConnect(a ISession) //当连接关闭	或者客户端主动发送close命令
	}

	//GetPushComponents() IPushComponents
	//GetSerializationComponents() ISerializationComponents
	IApp interface {
		//NewRPCClient() IRPCClient
		//NewRPCServer(addr string) IRPCServer
		//GetRPCClient() IRPCClient

		GetConfig() *config.Config
		GetAppPath() string
		GetBinPath() string
	}

	//IRPCClient interface {
	//	Call(servicePath string, serviceMethod string, args interface{}, reply interface{}) error
	//	Go(servicePath string, serviceMethod string, args interface{}, reply interface{}, chanRet chan *client.Call)
	//	CallAllGate(serviceMethod string, args interface{}, reply interface{}, chanRet chan *client.Call)
	//	OnStop()
	//}

	IRPCServer interface {
		OnStart()
		OnStop()
		RegisterHandler(rcvr interface{}, metadata string) error
		RegisterByName(name string, rcvr interface{}, metadata string) error
		RegisterByNodeId(nodeId string, rcvr interface{}, metadata string) error
		RegisterFunction(svrPath string, rcvr interface{}, metadata string) error
	}

	// StorageHandler Session信息持久化
	StorageHandler interface {
		/**
		  存储用户的Session信息
		  Session Bind Userid以后每次设置 settings都会调用一次Storage
		*/
		Storage(session ISession) (err error)
		/**
		  强制删除Session信息
		*/
		Delete(session ISession) (err error)
		/**
		  Bind user时会调用Query获取最新信息
		*/
		Query(uid int) (data []byte, err error)
	}

	ISerialization interface {
		Serialization(data interface{}) ([]byte, error)
		Deserialization(sess []byte, data interface{}) error
	}

	IGate interface {

	}
)
