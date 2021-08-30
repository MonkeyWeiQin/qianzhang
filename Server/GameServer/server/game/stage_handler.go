package game

import (
	"battle_rabbit/codec"
	"battle_rabbit/define"
	"battle_rabbit/excel"
	"battle_rabbit/global"
	"battle_rabbit/iface"
	"battle_rabbit/model"
	"battle_rabbit/protocol"
	"battle_rabbit/protocol/request"
	"battle_rabbit/protocol/response"
	"battle_rabbit/service/log"
	"battle_rabbit/utils"
	jsoniter "github.com/json-iterator/go"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

// 获取活动关卡配置,开关等
func (g *Game) GetActiveStageConf(sess iface.ISession, msg *codec.Message) {
	activeStages := model.GetActiveStageConf()
	sess.Send(protocol.SuccessData(msg.Id, activeStages))
	return
}

// 玩家关卡记录(活动和普通)
func (g *Game) GetUserStageRecord(sess iface.ISession, msg *codec.Message) {
	player := g.GetPlayer(sess.GetUid())
	if player == nil {
		sess.Send(protocol.ErrCode(msg.Id, define.MsgCode401))
		return
	}

	save := false
	t := int(utils.GetMidnightTime().Unix())

	if player.Stage.GoldStage.Time != t {
		player.Stage.GoldStage.Time = t
		player.Stage.GoldStage.Num = 3
		save = true
	}

	if player.Stage.ChallengeStage.Time != t {
		player.Stage.ChallengeStage.Time = t
		player.Stage.ChallengeStage.Num = 3
		save = true
	}

	if player.Stage.DefenseStage.Time != t {
		player.Stage.DefenseStage.Time = t
		player.Stage.DefenseStage.Num = 3
		save = true
	}

	if player.Stage.ResourceStage.Time != t {
		player.Stage.ResourceStage.Time = t
		player.Stage.ResourceStage.Num = 3
		save = true
	}

	// 无尽关卡
	//TODO 跨月处理
	if player.Stage.EndlessStage.Time != t {
		player.Stage.EndlessStage.Time = t
		player.Stage.EndlessStage.Num = 3
		save = true
	}

	if save {
		model.GetStageColl().UpdateOne(nil, bson.M{"_id": player.Stage.Id}, player.Stage)
	}
	sess.Send(protocol.SuccessData(msg.Id, player.Stage))
}

// 玩家开始进入某个关卡
func (g *Game) StartPlayStage(sess iface.ISession, msg *codec.Message) {
	var code = define.MsgCode200
	defer func() {
		if code != define.MsgCode200 {
			sess.Send(protocol.ErrCode(msg.Id, code))
		}
	}()
	//stageId 普通关卡是int, 活动关卡是string
	req := new(request.StartPlayStageReq)
	err := jsoniter.Unmarshal(msg.Data, req)
	if err != nil {
		log.Error("解析请求出错 :", err)
		code = define.MsgCode400
		return
	}
	player := g.GetPlayer(sess.GetUid())
	if player == nil {
		sess.Send(protocol.ErrCode(msg.Id, define.MsgCode401))
		return
	}

	if req.StageType == define.GeneralStageType || req.StageType == define.DifficultStageType {
		strength := 0
		stageId := int(req.StageId.(float64))
		if req.StageType == define.GeneralStageType {
			strength = excel.GeneralStageConf[stageId].EnergyConsume
		} else {
			strength = excel.DifficultStageConf[stageId].EnergyConsume
		}
		ok := player.Account.FlushStrength( - strength )
		if  !ok {
			log.Error(err, ok)
			code = define.MsgCode500
			return
		}
		// 推送体力
		protocol.NoticeUserStrengthChange(sess, player.Account.Strength,utils.GetNowTimeStamp())
	} else {
		if !player.Stage.CanActiveStageUseNumber(player.Uid, req.StageType) {
			code = define.MsgCode602
			return
		}
	}

	resp := new(response.PlayerJoinGameResp)
	role := player.Role

	if player.mainAttr == nil {
		player.TotalAttribute()
	}

	for _, weapon := range player.useEquip {
		switch weapon.Type {
		case define.MainWeaponType:
			resp.WeaponMain = &response.JoinGameWeaponMainResp{
				Level:      weapon.Level,
				RelationId: weapon.RelationId,
				Attribute:  player.mainAttr,
			}

		case define.SubWeaponType:
			resp.WeaponSub = &response.JoinGameWeaponSubResp{
				Level:      weapon.Level,
				RelationId: weapon.RelationId,
				Attribute:  player.subAttr,
			}
		}
	}
	resp.Role = &response.JoinGameRoleResp{
		RLv:        role.RLv,
		RExp:       role.RExp,
		Quality:    role.Quality,
		RelationId: role.RelationId,
		SkinId:     role.SkinId,
		Attribute:  player.roleAttr,
	}

	// 英雄自己的属性 + 英雄武器属性 + 被动技能属性 + 羁绊技能(若存在) + 英雄装备(若存在)
	if len(player.useHero) > 0 {
		fetterSkill := map[string][]*response.JoinGameHeroResp{}
		for _, hero := range player.useHero {
			// 装备属性
			var attr = new(global.Attribute)
			//英雄自己的属性
			xlsHero := excel.HeroDataConf[excel.GetConfigId(hero.RelationId, hero.Lv)]
			//英雄武器属性
			heroWeapon := excel.HeroWeaponDataConf[xlsHero.WeaponId]
			*attr = *heroWeapon.Attribute
			attr.AttackA += xlsHero.Attack
			attr.CriticalB += xlsHero.Critical
			attr.MoveSpeedA += xlsHero.MoveSpeed
			//英雄装备(若存在)
			equipment, ok := excel.HeroEquipmentDataConf[hero.EquipmentId]
			if ok {
				attr = attr.Add(equipment.Attribute)
			}
			// 被动技能属性
			for _, psId := range hero.PsIds {
				if psId == "0" || psId == "" {
					continue
				}
				attr = attr.Add(excel.HeroPassiveSkillConf[psId].Attribute)
			}

			heroResp := &response.JoinGameHeroResp{
				Mid:           hero.Id,
				RelationId:    hero.RelationId,
				Lv:            hero.Lv,
				StarLv:        hero.StarLv,
				Attribute:     attr,
				ActiveSkillLv: excel.HeroActiveSkillConf[hero.ASId].Level,
				PassSkillIds:  hero.PsIds,
			}
			// 羁绊技能(若存在)
			if xlsHero.FSId != "0" && xlsHero.FSId != "" {
				fetterSkill[xlsHero.FSId] = append(fetterSkill[xlsHero.FSId], heroResp)
			}
			resp.HeroList = append(resp.HeroList, heroResp)
		}

		// 这里判断是否有羁绊关系
		// 羁绊技能: 两个以上的英雄同时拥有一样的基本技能ID,那么这些拥有相同基本技能的英雄就可以增加这个基本技能对应的属性
		for psId, resps := range fetterSkill {
			if len(resps) >= 2 {
				for _, resp := range resps {
					resp.Attribute = resp.Attribute.Add(excel.HeroPassiveSkillConf[psId].Attribute)
				}
			}
		}
	}
	sess.Send(protocol.SuccessData(msg.Id, resp))
}

// 关卡结算 104010
func (g *Game) SettlementStage(sess iface.ISession, msg *codec.Message) {
	var code = define.MsgCode200
	defer func() {
		if code != define.MsgCode200 {
			sess.Send(protocol.ErrCode(msg.Id, code))
		}
	}()
	// 获取当前关卡的记录 第一次通关和重复通关奖励不一样
	req := new(request.SettlementStageReq)
	err := jsoniter.Unmarshal(msg.GetData(), &req)
	if err != nil {
		log.Error("解析请求出错 :", err)
		code = define.MsgCode400
		return
	}
	player := g.GetPlayer(sess.GetUid())
	if player == nil {
		sess.Send(protocol.ErrCode(msg.Id, define.MsgCode401))
		return
	}

	var r *response.AttachmentsResp
	switch req.StageType {
	case define.GeneralStageType, define.DifficultStageType: // 普通,困难

		r, code = g.mainStageSettlement(player, sess, req)
		if code != define.MsgCode200 {
			return
		}
		resp := &response.SettlementStageResp{
			StageId:     req.StageId,
			StageType:   req.StageType,
			StagePass:   req.StagePass,
			Attachments: r.Attachments,
			Replace:     r.Replace,
		}
		sess.Send(protocol.SuccessData(msg.Id, resp))
		return

	case define.ActiveStageTypeChallenge, // 挑战
		//define.ActiveStageTypeDefense, // 防守 TODO 还未开通
		define.ActiveStageTypeResource: // 资源

		resp := &response.SettlementStageResp{
			StageId:   req.StageId,
			StageType: req.StageType,
			StagePass: req.StagePass,
		}
		if !req.StagePass {
			sess.Send(protocol.SuccessData(msg.Id, resp))
			return
		}

		r, code = g.activeStageSettlement(player, sess, req)
		if code != define.MsgCode200 {
			return
		}

		resp.Attachments = r.Attachments
		resp.Replace = r.Replace
		sess.Send(protocol.SuccessData(msg.Id, resp))
		return
	case define.ActiveStageTypeGold: // 金币
		gold := req.Items[0].Count
		_, err = model.GetUserCollection().UpdateOne(nil, bson.M{"uid": player.Uid}, bson.M{"gold": player.Account.Gold + gold})
		if err != nil {
			code = define.MsgCode500
			return
		}
		player.Account.Gold += gold
		resp := &response.SettlementStageResp{
			StageId:   req.StageId,
			StageType: req.StageType,
			StagePass: req.StagePass,
		}
		resp.Attachments = append(resp.Attachments, &global.Item{
			ItemID:   "gold",
			Count:    gold,
			ItemType: define.ItemGoldType,
		})
		sess.Send(protocol.SuccessData(msg.Id, resp))
		return

	case define.ActiveStageTypeEndless: // 无尽
		_, err = model.GetStageColl().UpdateOne(nil, bson.M{"_id": player.Stage.Id}, bson.M{"endlessStage.monster": req.KillPeople, "endlessStage.lastStageId": req.StageId})
		if err != nil {
			code = define.MsgCode500
			return
		}
		player.Stage.EndlessStage.LastStageId = req.StageId.(string)
		player.Stage.EndlessStage.Monster += req.KillPeople
		sess.Send(protocol.SuccessData(msg.Id, map[string]int{"KillPeople": player.Stage.EndlessStage.Monster}))
		return
	default:
		code = define.MsgCode500
		return
	}
}

// define.GeneralStageType || define.DifficultStageType
func (g *Game) mainStageSettlement(player *Player, sess iface.ISession, req *request.SettlementStageReq) (r *response.AttachmentsResp, code int) {
	var (
		reward    *excel.RewardDataConfig
		stageId   = int(req.StageId.(float64))
		updater   = bson.M{}
		stageMain = player.Stage.StageMain
	)
	if req.StagePass {
		if req.StageType == define.GeneralStageType {
			if stageId < player.Stage.StageMain.GeneralStage { // 二次刷
				reward = excel.RewardDataConf[excel.GeneralStageConf[stageId].RepeatReward]
			} else {
				reward = excel.RewardDataConf[excel.GeneralStageConf[stageId].SucceedReward]
				stageMain.GeneralStage = excel.GeneralStageConf[stageId].NextLevelId
				updater["stageMain.generalStage"] = stageMain.GeneralStage
				if excel.GeneralStageConf[stageId].ChapterReward != "0" {
					stageMain.RewardGeneral[excel.GeneralStageConf[stageId].MapId] = excel.GeneralStageConf[stageId].ChapterReward
					updater["stageMain.rewardGeneral"] = stageMain.RewardGeneral
				}
			}
		} else {
			if stageId < player.Stage.StageMain.DifficultStage { // 二次刷
				reward = excel.RewardDataConf[excel.DifficultStageConf[stageId].RepeatReward]
			} else {
				reward = excel.RewardDataConf[excel.DifficultStageConf[stageId].SucceedReward]
				stageMain.DifficultStage = excel.DifficultStageConf[stageId].NextLevelId
				updater["stageMain.difficultStage"] = stageMain.DifficultStage
				if excel.DifficultStageConf[stageId].ChapterReward != "0" {
					stageMain.RewardDifficult[excel.DifficultStageConf[stageId].MapId] = excel.DifficultStageConf[stageId].ChapterReward
					updater["stageMain.rewardDifficult"] = stageMain.RewardDifficult
				}
			}
		}
		// 组合item
		req.Items = append(req.Items, rewardToItems(reward)...)
	}

	// 主角
	// 账号
	stageLog := &model.PlayerStageLogModel{
		Uid:        player.Uid,
		StageId:    stageId,
		StagePass:  req.StagePass,
		CreateTime: time.Now().Unix(),
		StageType:  req.StageType,
		Items:      req.Items,
	}

	if InsertErr := model.GetStageLogCollection().Insert(stageLog); InsertErr != nil {
		log.Error("新增出错 :", InsertErr)
		return nil, define.MsgCode500
	}

	resp, code := saveItem(g, player, sess, req.Items)
	if code != 200 {
		return nil, code
	}

	if len(updater) > 0 {
		if _, err := model.GetStageColl().UpdateOne(nil, bson.M{"uid": player.Uid}, updater); err != nil {
			log.Error("更新用户当前最新关卡数失败")
			return nil, define.MsgCode500
		}
	}

	return resp, define.MsgCode200
}

func rewardToItems(reward *excel.RewardDataConfig) (items []*global.Item) {
	if reward.Gold > 0 {
		items = append(items, &global.Item{
			Count:    reward.Gold,
			ItemType: define.ItemGoldType,
		})
	}

	if reward.Diamond > 0 {
		items = append(items, &global.Item{
			Count:    reward.Diamond,
			ItemType: define.ItemDiamondType,
		})
	}

	if reward.PhysicalPower > 0 {
		items = append(items, &global.Item{
			Count:    reward.PhysicalPower,
			ItemType: define.ItemStrengthType,
		})
	}

	if reward.Exp > 0 {
		items = append(items, &global.Item{
			Count:    reward.Exp,
			ItemType: define.ItemExpType,
		})

	}

	if reward.Hero != "0" {
		items = append(items, global.DisassembleToItem(reward.Hero, define.ItemHeroType)...)
	}

	if reward.Card != "0" {
		items = append(items, global.DisassembleToItem(reward.Card, define.ItemCardType)...)
	}
	if reward.MainWeapon != "0" {
		items = append(items, global.DisassembleToItem(reward.MainWeapon, define.ItemMainWeaponType)...)
	}

	if reward.SubWeapon != "0" {
		items = append(items, global.DisassembleToItem(reward.SubWeapon, define.ItemSubWeaponType)...)
	}

	if reward.Armor != "0" {
		items = append(items, global.DisassembleToItem(reward.Armor, define.ItemArmorType)...)
	}

	if reward.Ornament != "0" {
		items = append(items, global.DisassembleToItem(reward.Ornament, define.ItemOrnamentsType)...)
	}

	if reward.HeroEquip != "0" {
		items = append(items, global.DisassembleToItem(reward.HeroEquip, define.ItemHeroEquipType)...)
	}
	return
}

func (g *Game) activeStageSettlement(player *Player, sess iface.ISession, req *request.SettlementStageReq) (r *response.AttachmentsResp, code int) {
	var (
		items   []*global.Item
		stageId = req.StageId.(string)
	)
	if req.StageType == define.ActiveStageTypeResource {
		conf, ok := excel.ResourceStageConf[stageId]
		if !ok {
			log.Error("资源关卡没有找到奖励配置: id: %s", stageId)
			code = define.MsgCode500
			return
		}
		items = global.DisassembleToItem(conf.RewardIds, define.ItemCardType)
		if len(items) == 0 {
			log.Warn("没有拆分出奖励的物品: 源数据 : %s", conf.RewardIds)
			code = define.MsgCode200
			return
		}
	} else { // 挑战关卡
		conf, ok := excel.ChallengeStageConf[stageId]
		if !ok {
			log.Error("资源关卡没有找到奖励配置: id: %s", stageId)
			code = define.MsgCode500
			return
		}
		reward, ok := excel.ActiveStageRewardConf[conf.RewardId]
		if !ok {
			log.Error("活动关卡没有找到奖励配置: id: %s", stageId)
			code = define.MsgCode500
			return
		}

		items = rewardToItems(reward)
		if len(items) == 0 {
			log.Warn("没有拆分出奖励的物品: 源数据 : %s", conf.RewardId)
			code = define.MsgCode200
			return
		}
	}
	resp, code := saveItem(g, player, sess, items)
	if code != 200 {
		return nil, code
	}

	return resp, define.MsgCode200
}

/*章节通关奖励*/
func (g *Game) ReceiveChapterReward(sess iface.ISession, msg *codec.Message) {
	req := new(request.ReceiveChapterRewardReq)
	err := jsoniter.Unmarshal(msg.Data, &req)
	if err != nil {
		sess.Send(protocol.ErrCode(msg.Id, define.MsgCode400))
		return
	}

	player := g.GetPlayer(sess.GetUid())
	if player == nil {
		sess.Send(protocol.ErrCode(msg.Id, define.MsgCode401))
		return
	}
	var (
		mainStage = player.Stage.StageMain
		rewardId  string
		ok        bool
		updater = bson.M{}
	)

	if req.Type == define.GeneralStageType {
		rewardId, ok = mainStage.RewardGeneral[req.MapId]
		if !ok {
			sess.Send(protocol.ErrCode(msg.Id, define.MsgCode404))
			return
		}
		delete(mainStage.RewardGeneral,req.MapId)
		updater["stageMain.rewardGeneral"] = mainStage.RewardGeneral

	} else {
		rewardId, ok = mainStage.RewardDifficult[req.MapId]
		if !ok {
			sess.Send(protocol.ErrCode(msg.Id, define.MsgCode404))
			return
		}
		delete(mainStage.RewardDifficult,req.MapId)
		updater["stageMain.rewardDifficult"] = mainStage.RewardDifficult
	}

	if reward, _ok := excel.RewardDataConf[rewardId]; _ok && ok {
		items := rewardToItems(reward)
		resp, code := saveItem(g, player, sess, items)
		if code != 200 {
			sess.Send(protocol.ErrCode(msg.Id, code))
			return
		}
		n, err := model.GetStageColl().UpdateOne(nil, bson.M{"uid": sess.GetUid()}, updater)
		if n == 0 || err != nil {
			log.Error("更新章节奖领取记录失败 : ", err, n)
			sess.Send(protocol.ErrCode(msg.Id, define.MsgCode500))
			return
		}
		sess.Send(protocol.SuccessData(msg.Id, resp))
		return
	}

	sess.Send(protocol.ErrCode(msg.Id, define.MsgCode404))
}
