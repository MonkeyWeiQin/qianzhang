package game

import (
	"battle_rabbit/codec"
	"battle_rabbit/excel"
	"battle_rabbit/iface"
	"battle_rabbit/model"
	"battle_rabbit/network"
	"battle_rabbit/protocol"
	"battle_rabbit/service/log"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

// 修改账号资料 TODO

// 心跳
func (g *Game) Heartbeat(sess iface.ISession, msg *codec.Message) {
	sess.Send(protocol.SuccessData(msg.Id, map[string]int64{"time": time.Now().Unix()}))
}

//更新并推送
func (g *Game) playAddExpAndPushClient(sess *network.Session, p *Player, addExp int, accountAdd bool) (r *protocol.LevelUpgradeNoticeData, ok bool) {
	r, ok = g.playAddExp(p, addExp, accountAdd)
	if !ok {
		return
	}
	protocol.NoticeRoleLevelUpgrade(sess, r)
	return
}

// 更新不推送
func (g *Game) playAddExp(p *Player, addExp int, accountAdd bool) (r *protocol.LevelUpgradeNoticeData, ok bool) {
	var (
		role = p.Role
		acc  = p.Account
		push = &protocol.LevelUpgradeNoticeData{}
	)

	if accountAdd {
		acc.Exp += addExp
	Account:
		// 处理账号加经验
		accountConf, ok := excel.AccountUpgradeConf[acc.Level]
		if !ok {
			return nil, false
		}

		if acc.Level < accountConf.MaxLv && acc.Exp > accountConf.Exp {
			acc.Exp -= accountConf.Exp
			acc.Level++
			push.UpAcc = true
			goto Account
		}

		if acc.Level == accountConf.MaxLv && acc.Exp > accountConf.Exp {
			acc.Exp = accountConf.Exp
		}
		_, err := model.GetUserCollection().UpdateOne(nil, bson.M{"uid": p.Uid}, bson.M{"lv": acc.Level, "exp": acc.Exp})
		if err != nil {
			log.Error(err)
			return nil, false
		}
	}

	role.RExp += addExp
Role:
	conf, ok := excel.RoleDataConf[excel.GetConfigId(role.RelationId, role.RLv)]
	if !ok {
		fmt.Printf("Role upgrade Lv warning: config with nil. relationId: %s lv : %d, global RoleDataConf data len: %d ", role.RelationId, role.RLv, len(excel.RoleDataConf))
		return nil, false
	}

	if role.RExp >= conf.Exp && conf.NextId != "0" {
		if excel.RoleDataConf[conf.NextId].StarLevel != conf.StarLevel { // 需要提升星级才能继续升级
			if role.RExp > conf.Exp {
				role.RExp = conf.Exp
			}
		} else {
			role.RExp -= conf.Exp
			role.RLv++
			push.UpRole = true
			goto Role
		}
	}
	_, err := model.RoleCollection().UpdateOne(nil, bson.M{"uid": p.Uid}, bson.M{"rlv": role.RLv, "rExp": role.RExp})
	if err != nil {
		log.Error(err)
		return nil, false
	}
	push.Lv = acc.Level
	push.Exp = acc.Exp
	push.RLv = role.RLv
	push.RExp = role.RExp
	return push, true
}
