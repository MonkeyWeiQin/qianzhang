package network

import (
	"battle_rabbit/codec"
	"battle_rabbit/iface"
	"battle_rabbit/service/log"
	"context"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"net"
	"runtime/debug"
	"time"
)

type (
	Agent struct {
		ctx        context.Context
		cancel     context.CancelFunc
		session    *Session
		gate       *Gate
		conn       net.Conn
		stopped    bool
		writerChan chan iface.IMessage
		SessionStorage
		SessionLearner
	}
)

func newAgent(ctx context.Context, conn net.Conn, gate *Gate) *Agent {
	ctx, cancel := context.WithCancel(ctx)
	agent := &Agent{
		ctx:        ctx,
		cancel:     cancel,
		conn:       conn,
		gate:       gate,
		writerChan: make(chan iface.IMessage, 1),
	}
	return agent
}

// 启动读写
func (agt *Agent) Run() error {
	defer func() {
		if e := recover(); e != nil {
			log.Error("[(agt *Agent) Run()]===>  :" + e.(error).Error())
			log.Error(string(debug.Stack()))
			return
		}
	}()
	// 创建session
	agt.session = NewSession(agt)
	agt.gate.connManager.Add(agt)
	agt.Connect(agt.session)

	go agt.WriterLoop()
	agt.ReadLoop()
	return nil
}

// 关闭读写,并退出
func (agt *Agent) OnClose() (err error) {
	defer func() {
		if e := recover(); e != nil {
			log.Alert("Agent OnClose recover ===: ", e)
			err = e.(error)
		}
	}()
	if agt.stopped {
		return nil
	}
	agt.DisConnect(agt.session)
	agt.stopped = true
	close(agt.writerChan)
	_ = agt.conn.Close()

	agt.gate.connManager.Remove(agt.session.GetSessionId())

	return nil
}

func (agt *Agent) ReadLoop() {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println("[Error]===>  ReadLoop() :", e)
			fmt.Println(string(debug.Stack()))
			return
		}
		agt.cancel()
	}()
	for {
		select {
		case <-agt.ctx.Done():
			return
		default:
			msg := new(codec.Message)
			err := msg.ReadPack(agt.conn)
			if err != nil {
				log.Warn("ReadLoop :: ", err)
				return
			}
			// 解密
			if agt.gate.encryption != nil {
				data, err := agt.gate.encryption.Decode(msg.GetData())
				if err != nil {
					continue
				}
				msg.Data = data
			}
			// ======================= Debug start=================================================================================================
			data := make(map[string]interface{})
			_ = jsoniter.Unmarshal(msg.Data, &data)
			str, _ := jsoniter.Marshal(data)
			log.Debug("读取到数据 ID：%d, data: %s  , 数据长度：%d ", msg.Id, string(str), msg.GetDataLen())
			// ======================= Debug end=================================================================================================

			agt.doMsgHandler(agt.session, msg)
		}
	}
}

func (agt *Agent) WriterLoop() {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println("[Error]===>  WriterLoop() :" + e.(error).Error())
			fmt.Println(string(debug.Stack()))
		}
		agt.cancel()
	}()

	for {
		select {
		case <-agt.ctx.Done():
			return
		case data, ok := <-agt.writerChan:
			if ok {
				// ======================= Debug start=================================================================================================
				tmp := make(map[string]interface{})
				_ = jsoniter.Unmarshal(data.GetData(), &tmp)
				str, _ := jsoniter.Marshal(tmp)
				log.Debug("发送数据 ID: %d, data: %s , 数据长度：%d ", data.GetMsgId(), string(str), data.GetDataLen())
				// ======================= Debug end=================================================================================================

				// 加密
				if agt.gate.encryption != nil {
					by, err := agt.gate.encryption.Decode(data.GetData())
					if err != nil {
						continue
					}
					data.SetData(by)
				}

				// 发送数据
				err := data.WriterPack(agt.conn)
				if err != nil {
					fmt.Println("Send Data error:, ", err, " Conn Writer exit")
					return
				}
			} else {
				log.Debug("msgBuffChan is Closed")
				return
			}
		}
	}
}

func (agt *Agent) SendMessage(message iface.IMessage) {
	defer func() {
		if e := recover(); e != nil {
			log.Error(e.(error).Error())
		}
	}()
	// 超时控制2秒,以免channal装满引起阻塞
	ti := time.NewTimer(time.Second * 5)
	select {
	case <-ti.C:
		return
	case agt.writerChan <- message:
		ti.Stop()
		return
	}

}

func (agt *Agent) GetSession() iface.ISession {
	return agt.session
}

func (agt *Agent) GetUid() int {
	return agt.session.GetUid()
}

func (agt *Agent) Stopped() bool {
	return agt.stopped
}

// msg handler
func (agt *Agent) doMsgHandler(session *Session, msg *codec.Message) {
	defer func() {
		if e := recover(); e != nil {
			log.Error("DoMsgHandler =====:", e)
			log.Error(string(debug.Stack()))
		}
	}()

	// 检验是否已经通过认证了
	if handler, ok := agt.gate.handles[msg.GetMsgId()]; ok {
		if m, ok := agt.gate.canPermission(session, msg); !ok {
			agt.writerChan <- m
			return
		}
		handler(agt.session, msg)
	} else {
		log.Warn("收到 [未注册消息] ID: %d", msg.GetMsgId())
	}
}
