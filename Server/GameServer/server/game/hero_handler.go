package game

import (
	"battle_rabbit/codec"
	"battle_rabbit/define"
	"battle_rabbit/excel"
	"battle_rabbit/iface"
	"battle_rabbit/model"
	"battle_rabbit/protocol"
	"battle_rabbit/protocol/request"
	"battle_rabbit/service/log"
	"battle_rabbit/utils"
	jsoniter "github.com/json-iterator/go"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	HeroUpgradeLv            int8 = 10 // 升级
	HeroUpgradeStar          int8 = 20 // 升星
	HeroUpgradeMainSkill     int8 = 30 // 主动技能升级
	HeroUpgradePassiveSkill1 int8 = 40 // 被动技能1升级
	HeroUpgradePassiveSkill2 int8 = 41 // 被动技能2升级
	HeroUpgradePassiveSkill3 int8 = 42 // 被动技能3升级
	HeroUpgradePassiveSkill4 int8 = 43 // 被动技能4升级
	HeroUpgradeEquipment     int8 = 50 // 装备升级
	HeroUpgradeStrengthen    int8 = 60 // 强化
)

// GetHeroList 获取英雄列表 110020
func (g *Game) GetHeroList(sess iface.ISession, msg *codec.Message) {
	player := g.GetPlayer(sess.GetUid())
	if player == nil {
		sess.Send(protocol.ErrCode(msg.Id, define.MsgCode401))
	}
	var list []*model.HeroModel
	for _, heroModel := range player.Hero {
		list = append(list, heroModel)
	}
	sess.Send(protocol.SuccessData(msg.Id, list))
}

// HeroGoToWar 英雄出战 110030
func (g *Game) HeroGoToWar(sess iface.ISession, msg *codec.Message) {
	var code = define.MsgCode200
	defer func() {
		if code != define.MsgCode200 {
			sess.Send(protocol.ErrCode(msg.Id, code))
		}
	}()
	req := new(request.HeroGoToWarRequest)
	err := jsoniter.Unmarshal(msg.GetData(), &req)
	if err != nil {
		log.Error("解析请求出错 :", err)
		code = define.MsgCode400
		return
	}

	player := g.GetPlayer(sess.GetUid())
	if player == nil {
		code = define.MsgCode401
		return
	}
	if len(player.useHero) == 3 {
		log.Error("出战失败，英雄数量上限 :")
		code = define.MsgCode1003
		return
	}

	if _, ok := player.useHero[req.HeroId]; ok {
		log.Error("英雄已出战 :")
		code = define.MsgCode1004
		return
	}

	hero, ok := player.Hero[req.HeroId]
	if !ok {
		log.Error("出战失败，没查询到英雄数据 :")
		code = define.MsgCode1005
		return
	}
	hero.Use = true
	player.useHero[hero.Id] = hero

	if n, err := model.GetHeroCollection().UpdateOne(nil, bson.M{"_id": hero.Id}, bson.M{"use": true}); err != nil || n == 0 {
		log.Error("出战失败 :", err, n)
		code = define.MsgCode500
		return
	}

	sess.Send(protocol.SuccessData(msg.Id, map[string]string{
		"relationId": hero.RelationId,
	}))
}

// HeroCancelGoToWar 英雄取消出战 110040
func (g *Game) HeroCancelGoToWar(sess iface.ISession, msg *codec.Message) {
	req := new(request.HeroGoToWarRequest)
	err := jsoniter.Unmarshal(msg.GetData(), &req)
	if err != nil {
		log.Error("解析请求出错 :", err)
		sess.Send(protocol.ErrCode(msg.Id, define.MsgCode400))
		return
	}
	player := g.GetPlayer(sess.GetUid())
	if player == nil {
		sess.Send(protocol.ErrCode(msg.Id, define.MsgCode401))
		return
	}

	if hero, ok := player.useHero[req.HeroId]; !ok {
		log.Error("没查询到英雄数据 id: %s ", req.HeroId)
		sess.Send(protocol.ErrCode(msg.Id, define.MsgCode404))
		return
	} else {
		hero.Use = false
		delete(player.useHero, hero.Id)
		// 更新
		if n, err := model.GetHeroCollection().UpdateOne(nil, bson.M{"_id": hero.Id}, bson.M{"use": false}); err != nil || n == 0 {
			log.Error("取消失败!! :", err, n)
			sess.Send(protocol.Err(msg.Id))
			return
		}
		sess.Send(protocol.SuccessData(msg.Id, map[string]string{
			"relationId": hero.RelationId,
		}))
	}
}

// 英雄升级
//HeroUpgradeLv            int8 = 10 // 升级
//HeroUpgradeStar          int8 = 20 // 升星
//HeroUpgradeMainSkill     int8 = 30 // 主动技能升级
//HeroUpgradePassiveSkill1 int8 = 40 // 被动技能1升级
//HeroUpgradePassiveSkill2 int8 = 41 // 被动技能2升级
//HeroUpgradePassiveSkill3 int8 = 42 // 被动技能3升级
//HeroUpgradePassiveSkill4 int8 = 43 // 被动技能4升级
func (g *Game) HeroUpgrade(sess iface.ISession, msg *codec.Message) {
	var (
		err      error
		respCode = define.MsgCode200
	)
	defer func() {
		//if e := recover(); e != nil {
		//	log.Error(e)
		//	respCode = define.MsgCode500
		//}
		if err != nil || respCode != define.MsgCode200 {
			log.Error(err)
			sess.Send(protocol.ErrCode(msg.Id, respCode))
		}
	}()
	req := new(request.HeroUpgradeReq)
	err = jsoniter.Unmarshal(msg.Data, &req)
	if err != nil {
		respCode = define.MsgCode500
		return
	}
	player := g.GetPlayer(sess.GetUid())
	if player == nil {
		sess.Send(protocol.ErrCode(msg.Id, define.MsgCode401))
		return
	}

	heroModel, ok := player.heroIndex[req.Rid]
	if !ok {
		respCode = define.MsgCode1005
		return
	}

	var (
		gold            = 0
		diamond         = 0
		deplete, nextId string
		maxLv           bool
		cards           *model.UpgradeCard
		taskInfo        []*model.TaskProgress
	)
	switch req.Ty {
	case HeroUpgradeLv, HeroUpgradeStar: // 升级 升星
		conf := excel.HeroDataConf[excel.GetConfigId(heroModel.RelationId, heroModel.Lv)]
		nextConf := excel.HeroDataConf[excel.GetConfigId(heroModel.RelationId, heroModel.Lv+1)]

		if req.Ty == HeroUpgradeLv {
			// 提升等级,判断星级是否达到条件
			if nextConf != nil && conf.StarLevel != nextConf.StarLevel {
				respCode = define.MsgCode802
				return
			}
			gold = conf.UpGold
			diamond = conf.UpDiamond
			deplete = conf.UpDeplete
			maxLv = conf.Level == conf.MaxLevel
		} else {
			// 判断提升星级,判断等级是否达到条件
			if nextConf != nil && conf.StarLevel+1 != nextConf.StarLevel {
				respCode = define.MsgCode801
				return
			}
			gold = conf.StarGold
			diamond = conf.StarDiamond
			deplete = conf.StarDeplete
			maxLv = heroModel.StarLv == define.MaxStarLevel
		}

	case HeroUpgradeMainSkill,
		HeroUpgradePassiveSkill1,
		HeroUpgradePassiveSkill2,
		HeroUpgradePassiveSkill3,
		HeroUpgradePassiveSkill4: // 主动技能升级 被动动技能升级

		var conf *excel.SkillConfig
		if req.Ty == HeroUpgradeMainSkill {
			conf = excel.HeroActiveSkillConf[heroModel.ASId]
		} else {
			if int8(len(heroModel.PsIds)) < req.Ty-HeroUpgradePassiveSkill1+1 {
				respCode = define.MsgCode400
				return
			}
			conf = excel.HeroPassiveSkillConf[heroModel.PsIds[req.Ty-HeroUpgradePassiveSkill1]] // req.Ty-40 ==  0|1|2|3
		}
		gold = conf.UpgradeGold
		diamond = conf.UpgradeDiamond
		deplete = conf.UpgradeMaterials
		maxLv = conf.NextId == "0"
		nextId = conf.NextId

	default:
		respCode = define.MsgCode400
		return
	}

	if gold > player.Account.Gold {
		respCode = define.MsgCode604
		return
	}

	if diamond > player.Account.Diamond {
		respCode = define.MsgCode604
		return
	}

	if maxLv {
		respCode = define.MsgCode800
		return
	}

	packTabId, num := utils.UnmarshalItemsKV(deplete)
	if num > 0 {
		if icard, ok := player.packIndex[packTabId]; !ok {
			respCode = define.MsgCode404
			return
		} else {
			cards = icard.(*model.UpgradeCard)
			if num > cards.Num {
				respCode = define.MsgCode606
				return
			}
		}
	}

	switch req.Ty {
	case HeroUpgradeLv: // 升级
		heroModel.Lv++
		taskInfo = append(taskInfo, &model.TaskProgress{
			Type:      define.TaskHeroLevel,
			Condition: heroModel.RelationId,
			Lv:        heroModel.Lv,
		})
	case HeroUpgradeStar: // 升星
		heroModel.StarLv++
		heroModel.Lv++
	case HeroUpgradeMainSkill: // 主动技能升级
		heroModel.ASId = nextId
	case HeroUpgradePassiveSkill1,
		HeroUpgradePassiveSkill2,
		HeroUpgradePassiveSkill3,
		HeroUpgradePassiveSkill4: //被动动技能升级

		heroModel.PsIds[req.Ty-HeroUpgradePassiveSkill1] = nextId // 0,1,2,3
	}

	backpackId := ""
	var m map[string]int
	if num > 0 && cards != nil {
		cards.Num -= num
		num = cards.Num
		backpackId = cards.Id
		m = map[string]int{cards.TabId: cards.Num}
	}

	if gold != 0 {
		taskInfo = append(taskInfo, &model.TaskProgress{
			Type:      define.TaskConsumeGold,
			Num:       gold,
		})
		player.Account.Gold -= gold
		gold = player.Account.Gold
	}else{
		gold = -1
	}


	if  diamond != 0 {
		taskInfo = append(taskInfo, &model.TaskProgress{
			Type:      define.TaskConsumeDiamond,
			Num:       diamond,
		})
		player.Account.Diamond -= diamond
		diamond = player.Account.Diamond
	} else {
		diamond = -1
	}

	err = model.GetHeroCollection().Upgrade(nil, player.Uid, gold, diamond, heroModel, backpackId, num)
	if err != nil {
		respCode = define.MsgCode500
		return
	}
	// todo 刷新角色属性

	// 推送前端
	protocol.NoticeUserGoldAndDiamond(sess, player.Account.Gold, player.Account.Diamond)

	// TODO 推送背包,推送英雄数据

	sess.Send(protocol.SuccessData(msg.Id, map[string]interface{}{
		"hero": heroModel,
		"card": m,
	}))
	UpdateTask(player.Task,sess,taskInfo)
}
