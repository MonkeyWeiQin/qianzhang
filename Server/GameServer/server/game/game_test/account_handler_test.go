package game_test

import (
	"battle_rabbit/define"
	"battle_rabbit/protocol"
	"battle_rabbit/protocol/request"
	"battle_rabbit/service/log"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"time"

	"testing"
)

func login() *Cli {
	cli := NewTcpCli()
	msg ,_ :=  protocol.MakeReqMsg(define.CheckTokenTestMsgId,&request.BindSessionReq{Url: "xxxxx4444"})
	ReqChan <- msg
	time.Sleep(time.Second)
	return cli

}

func Test_CheckTokenTest(t *testing.T) {
	cli := NewTcpCli()
	defer cli.Stop()

	msg ,_ :=  protocol.MakeReqMsg(define.CheckTokenTestMsgId,&request.BindSessionReq{})
	ReqChan <- msg

	n := 1
	for n > 0 {
		select {
		case msg =  <- RespChan:
			fmt.Println("\r\n 客户端收到数据: ", string(msg.Data))
			resp := new(RespStatus)
			err := jsoniter.Unmarshal(msg.Data,resp)
			if err != nil {
				log.Error(err)
			}

			if resp.Code == 200 {
				n --
			}
		}
	}
	fmt.Println("退出测试!!!! ")
}

func TestGame_GetPlayerInfo(t *testing.T) {
	cli := login()
	defer cli.Stop()

	msg ,_ :=  protocol.MakeReqMsg(define.GetPlayerInfoMsgId,nil)
	ReqChan <- msg

	n := 2
	for n > 0 {
		select {
		case msg =  <- RespChan:
			fmt.Println("\r\n 客户端收到数据: ", string(msg.Data))
			resp := new(RespStatus)
			err := jsoniter.Unmarshal(msg.Data,resp)
			if err != nil {
				log.Error(err)
			}

			if resp.Code == 200 {
				n --
			}
		}
	}
	fmt.Println("退出测试!!!! ")
}