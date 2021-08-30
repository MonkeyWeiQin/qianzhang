package game

import (
	"battle_rabbit/codec"
	"battle_rabbit/define"
	"battle_rabbit/excel"
	"battle_rabbit/iface"
	"battle_rabbit/model"
	"battle_rabbit/protocol"
	"battle_rabbit/service/log"
	"battle_rabbit/service/mgoDB"
	jsoniter "github.com/json-iterator/go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// 获取背包卡片数据
func (g *Game) GetBackpackData(sess iface.ISession, msg *codec.Message) {
	p := g.GetPlayer(sess.GetUid())
	if p == nil {
		sess.Send(protocol.ErrCode(msg.Id, define.MsgCode401))
		return
	}

	var (
		back = p.Package
		data []interface{}
	)

	for _, pack := range back {
		data = append(data, pack)
	}
	sess.Send(protocol.SuccessData(msg.Id,data))
}

// 使用经验卡(给角色加经验)
// TODO 客户端没有一张可用的卡还在往上发,后期优化
func (g *Game) UseExpCard(sess iface.ISession, msg *codec.Message) {
	var req map[string]int
	err := jsoniter.Unmarshal(msg.Data, &req)
	if err != nil {
		log.Error(err)
		sess.Send(protocol.Err(msg.Id))
		return
	}
	player := g.GetPlayer(sess.GetUid())
	if player == nil {
		sess.Send(protocol.ErrCode(msg.Id, define.MsgCode401))
		return
	}
	exp := 0
	var cards = player.packIndex
	var update []mongo.WriteModel
	for cardId, n := range req {
		if icard, ok := cards[cardId]; ok {
			card := icard.(*model.UpgradeCard)
			if n > card.Num {
				// 出错了, 数据库中没有这个数据或者数量不足
				sess.Send(protocol.ErrCode(msg.Id, define.MsgCode606))
				return
			}
			exp += excel.CardDataConf[card.TabId].Exp * n
			card.Num -= n
			req[card.TabId] = card.Num
			// 将req中的数据改为剩下的数据丢给前端
			update = append(update, mongo.NewUpdateOneModel().SetFilter(bson.M{"uid": player.Uid, "tabId": card.TabId}).SetUpdate(bson.M{"$set": bson.M{"num": card.Num}}))
		} else {
			// 出错了, 数据库中没有这个数据或者数量不足
			sess.Send(protocol.ErrCode(msg.Id, define.MsgCode404))
			return
		}
	}

	br, err := mgoDB.GetMgo(define.MgoDBNameBattle).GetCol(model.GetBackPackColl().CollName).BulkWrite(nil, update)
	if err != nil || br.ModifiedCount == 0 {
		log.Error(err, br)
		sess.Send(protocol.Err(msg.Id))
		return
	}
	// 给角色加经验

	r, ok := g.playAddExp(player, exp, false)
	if !ok {
		sess.Send(protocol.Err(msg.Id))
		return
	}
	protocol.NoticeRoleLevelUpgrade(sess, r)

	sess.Send(protocol.SuccessData(msg.Id, map[string]interface{} {
		"rlv":   r.RLv,
		"rExp":  r.RExp,
		"cards": req,
	}))
}
