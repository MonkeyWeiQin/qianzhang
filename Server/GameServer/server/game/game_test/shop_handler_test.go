package game_test

import (
	"battle_rabbit/codec"
	"battle_rabbit/define"
	"battle_rabbit/service/log"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"testing"
)

func TestGame_GetShopList(t *testing.T) {
	cli := login()
	defer cli.Stop()

	msg := codec.NewMsgPackage(define.GetShopListMsgId,nil)

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