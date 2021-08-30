package game_test

import (
	"battle_rabbit/codec"
	"battle_rabbit/define"
	"fmt"
	"testing"
	"time"
)

func TestGame_FreeOpenChest(t *testing.T) {

	cli := login()
	defer cli.Stop()

	tick := time.After(time.Second * 10)

	msg := codec.NewMsgPackage(define.OneTimesOpenChestMsgId,[]byte("{\"id\":\"treasurebox004\"}"))

	ReqChan <- msg


	for {
		select {
		case msg =  <- RespChan:
			fmt.Println("\r\n 客户端收到数据: ", string(msg.Data))
		case <- tick:
			fmt.Println("退出...........!!! ")
			return
		}
	}



}


// {"code":200,"data":{"Attachments":[{"itemId":"diamond","count":100,"itemType":1},{"itemId":"card009_1","count":2,"itemType":5},{"itemId":"card001_2","count":2,"itemType":5},{"itemId":"card011_1","count":2,"itemType":5},{"itemId":"card010_4","count":2,"itemType":5},{"itemId":"card001_8","count":1,"itemType":5},{"itemId":"card002_1","count":2,"itemType":5},{"itemId":"card005_1","count":1,"itemType":5},{"itemId":"card008_1","count":1,"itemType":5},{"itemId":"card010_2","count":1,"itemType":5},{"itemId":"strength","count":1,"itemType":2},{"itemId":"card004_1","count":2,"itemType":5},{"itemId":"card006_3","count":2,"itemType":5},{"itemId":"card010_3","count":2,"itemType":5},{"itemId":"card010_1","count":1,"itemType":5},{"itemId":"card006_4","count":1,"itemType":5},{"itemId":"card014_1","count":2,"itemType":5},{"itemId":"card007_1","count":2,"itemType":5},{"itemId":"card013_1","count":1,"itemType":5},{"itemId":"gold","count":58,"itemType":0},{"itemId":"card001_1","count":1,"itemType":5},{"itemId":"card001_4","count":2,"itemType":5},{"itemId":"card001_5","count":2,"itemType":5},{"itemId":"card001_7","count":1,"itemType":5},{"itemId":"card006_2","count":2,"itemType":5},{"itemId":"card010_5","count":2,"itemType":5},{"itemId":"card012_1","count":2,"itemType":5},{"itemId":"card001_3","count":2,"itemType":5},{"itemId":"card001_6","count":1,"itemType":5},{"itemId":"card001_9","count":2,"itemType":5},{"itemId":"card003_1","count":1,"itemType":5},{"itemId":"card006_1","count":1,"itemType":5}],"Replace":{}}}
// {"code":200,"data":{"Attachments":[{"itemId":"card003_1","count":1,"itemType":5},{"itemId":"card006_1","count":2,"itemType":5},{"itemId":"card006_2","count":1,"itemType":5},{"itemId":"card010_5","count":2,"itemType":5},{"itemId":"card012_1","count":1,"itemType":5},{"itemId":"card001_3","count":2,"itemType":5},{"itemId":"card001_6","count":2,"itemType":5},{"itemId":"card001_9","count":2,"itemType":5},{"itemId":"diamond","count":100,"itemType":1},{"itemId":"card009_1","count":1,"itemType":5},{"itemId":"card001_2","count":2,"itemType":5},{"itemId":"card011_1","count":2,"itemType":5},{"itemId":"card008_1","count":2,"itemType":5},{"itemId":"card010_2","count":2,"itemType":5},{"itemId":"card010_4","count":2,"itemType":5},{"itemId":"card001_8","count":2,"itemType":5},{"itemId":"card002_1","count":1,"itemType":5},{"itemId":"card005_1","count":1,"itemType":5},{"itemId":"card010_3","count":2,"itemType":5},{"itemId":"strength","count":1,"itemType":2},{"itemId":"card004_1","count":1,"itemType":5},{"itemId":"card006_3","count":1,"itemType":5},{"itemId":"card010_1","count":1,"itemType":5},{"itemId":"card006_4","count":2,"itemType":5},{"itemId":"card014_1","count":2,"itemType":5},{"itemId":"card001_5","count":1,"itemType":5},{"itemId":"card001_7","count":1,"itemType":5},{"itemId":"card007_1","count":2,"itemType":5},{"itemId":"card013_1","count":2,"itemType":5},{"itemId":"gold","count":240,"itemType":0},{"itemId":"card001_1","count":1,"itemType":5},{"itemId":"card001_4","count":2,"itemType":5}],"Replace":{}}}


//{"code":200,"data":{"Attachments":[{"itemId":"weapon_02_1","count":1,"itemType":7}],"Replace":{}}}
//{"code":200,"data":{"Attachments":[{"itemId":"card011_1","count":1,"itemType":5}],"Replace":{}}}

//{"code":200,"data":{"Attachments":[{"itemId":"hero_05_1","count":1,"itemType":4},{"itemId":"hero_01_1","count":1,"itemType":4},{"itemId":"hero_02_1","count":1,"itemType":4}],"Replace":{}}}




