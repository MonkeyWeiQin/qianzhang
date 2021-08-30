package game_test

import (
	"battle_rabbit/codec"
	"battle_rabbit/define"
	"battle_rabbit/global"
	"battle_rabbit/protocol"
	"battle_rabbit/protocol/request"
	"battle_rabbit/service/log"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"testing"
)

func TestGame_SettlementStage(t *testing.T) {
	cli := login()
	defer cli.Stop()
	msg ,_ :=  protocol.MakeReqMsg(define.CreatePlayerStageMsgId,&request.SettlementStageReq{
		StageId:   100001,
		StageType: 1,
		StagePass: true,
		Items:     []*global.Item{
			{
				ItemID:   "",
				Count:    10,
				ItemType: 1,
			},
		},
	})

	msg = codec.NewMsgPackage(104010,[]byte("{\"stageId\":100003,\"stageType\":1,\"stagePass\":true,\"items\":[{\"itemId\":\"gold\",\"count\":408,\"itemType\":1}]}"))

	ReqChan <- msg

	n := 5
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

type A struct {
	ItemA int
}

type B struct {
	*A
	ItemB int
}

func TestGame_SettlementStage2(t *testing.T) {
	//b := B{
	//	A:&A{1}     ,
	//	ItemB: 2,
	//}
	//s ,_ := jsoniter.Marshal(b)




	//fmt.Println(string(s))
}