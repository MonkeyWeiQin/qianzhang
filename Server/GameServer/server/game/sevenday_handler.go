package game

import (
	"battle_rabbit/codec"
	"battle_rabbit/define"
	"battle_rabbit/excel"
	"battle_rabbit/global"
	"battle_rabbit/iface"
	"battle_rabbit/model"
	"battle_rabbit/protocol"
	"battle_rabbit/protocol/response"
	"battle_rabbit/service/log"
	"time"
)

// SevenDaySignIn 签到 101030
func (g *Game) SevenDaySignIn(sess iface.ISession, msg *codec.Message) {
	var code = define.MsgCode200
	defer func() {
		if code != define.MsgCode200 {
			sess.Send(protocol.ErrCode(msg.Id, code))
		}
	}()
	player := g.GetPlayer(sess.GetUid())
	if player == nil { // 重新登录
		code = define.MsgCode401
		return
	}
	signInfo := player.Account.Sign
	t := time.Now()
	zeroTm := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).Unix()

	if signInfo.LastTime > int(zeroTm) {
		log.Error("领取失败 : 当天已经领取了")
		code = define.MsgCode1303
		return
	}

	sevenDay, ok := excel.SevenDayConf[signInfo.Count+1]
	if ok {
		attas := []*global.Item{{
			ItemID:   sevenDay.ItemID,
			Count:    sevenDay.Count,
			ItemType: define.ItemType(sevenDay.ItemType),
		}}
		var resp *response.AttachmentsResp
		resp, code = saveItem(g, player, sess, attas)

		if code != define.MsgCode200 {
			log.Error("签到领取失败!!! ")
			return
		} else {
			_, err := model.GetUserCollection().SetSignInfo(player.Uid, signInfo.Count)
			if err != nil {
				log.Error("更新失败 :", err)
				code = define.MsgCode500
				return
			}
			signInfo.Count++
			signInfo.LastTime = int(time.Now().Unix())
			sess.Send(protocol.SuccessData(msg.Id,resp))
			return
		}
	} else {
		log.Error("配置未找到")
		code = define.MsgCode404
		return
	}
}
