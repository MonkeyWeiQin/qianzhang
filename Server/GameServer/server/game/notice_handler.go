package game

import (
	"battle_rabbit/codec"
	"battle_rabbit/iface"
	"battle_rabbit/model"
	"battle_rabbit/protocol"
	"battle_rabbit/service/log"
)

// GetNotice 获取有效的公告
func (g *Game) GetNotice(sess iface.ISession, msg *codec.Message) {
	list, err := model.GetNoticeCollection().GetNotice(nil)
	if err != nil {
		log.Error("获取失败", err.Error())
		sess.Send(protocol.Err(msg.Id))
		return
	}
	sess.Send(protocol.SuccessData(msg.Id, list))
}

// GetBroadCast 获取有效的公告
func (g *Game) GetBroadCast(sess iface.ISession, msg *codec.Message) {
	list, err := model.GetBroadcastCollection().GetBroadcast(nil)
	if err != nil {
		log.Error("获取失败", err.Error())
		sess.Send(protocol.Err(msg.Id))
		return
	}
	sess.Send(protocol.SuccessData(msg.Id, list))
}
