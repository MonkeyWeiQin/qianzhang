package network

import (
	"battle_rabbit/codec"
	"battle_rabbit/iface"
	"battle_rabbit/service/log"
	"context"
	"net"
	"time"
)

type (
	Handler func(session iface.ISession, msg *codec.Message)

	Gate struct {
		ctx         context.Context
		canal       context.CancelFunc
		addr        string
		server      iface.IServer
		connManager iface.IConnManager
		encryption  iface.IEncryption

		handles       map[uint32]Handler
		canPermission func(session *Session, msg *codec.Message) (*codec.Message, bool)
	}
)

func NewGate(addr string) *Gate {
	gate := new(Gate)
	gate.addr = addr
	gate.ctx, gate.canal = context.WithCancel(context.Background())
	gate.connManager = NewConnManager()
	gate.handles = make(map[uint32]Handler)
	gate.canPermission = func(session *Session, msg *codec.Message) (*codec.Message, bool) { return nil, true }
	return gate
}

func (gate *Gate) Start() {
	// init TCP server
	gate.server = NewTCPServer(gate.addr, gate.newAgent)
	go gate.server.Start()

	go func() {
		for {
			log.Debug("当前连接数量: ", gate.connManager.Len())
			time.Sleep(time.Second * 10)
		}
	}()
}

func (gate *Gate) OnStop() {
	gate.server.Stop()
	gate.connManager.ClearConn()
	log.Debug("Gate Service Stopped !!!")
}

func (gate *Gate) newAgent(conn net.Conn) iface.IAgent {
	return newAgent(gate.ctx, conn, gate)
}

func (gate *Gate) SetEncryptionComponent(handle iface.IEncryption) {
	gate.encryption = handle
}

func (gate *Gate) SetCheckPermissionFunc(fn func(session *Session, msg *codec.Message) (*codec.Message, bool)) {
	gate.canPermission = fn
}

func (gate *Gate) Register(msgId uint32, handler Handler) {
	if _, ok := gate.handles[msgId]; ok {
		log.Fatal("消息ID已经被占用!!")
	}
	gate.handles[msgId] = handler
}
