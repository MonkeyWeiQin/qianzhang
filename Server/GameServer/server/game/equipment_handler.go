package game

import (
	"battle_rabbit/codec"
	"battle_rabbit/define"
	"battle_rabbit/excel"
	"battle_rabbit/iface"
	"battle_rabbit/model"
	"battle_rabbit/protocol"
	"battle_rabbit/protocol/request"
	"battle_rabbit/protocol/response"
	"battle_rabbit/service/log"
	"battle_rabbit/utils"
	jsoniter "github.com/json-iterator/go"
	"go.mongodb.org/mongo-driver/bson"
)

//GetequipmentList 获取装备列表 107010
func (g *Game) GetEquipmentList(sess iface.ISession, msg *codec.Message) {
	req := new(request.GetEquipmentListRequest)
	if req.Page == 0 {
		req.Page = 1
	}
	err := jsoniter.Unmarshal(msg.GetData(), &req)
	if err != nil {
		log.Error("解析请求出错 :", err)
		sess.Send(protocol.ErrCode(msg.Id, define.MsgCode400))
		return
	}

	p := g.GetPlayer(sess.GetUid())
	if p == nil {
		sess.Send(protocol.ErrCode(msg.Id, define.MsgCode401))
		return
	}

	resp := &response.GetEquipmentListResponse{
		EquipmentType: req.EquipmentType,
	}

	for _, equipmentModel := range p.Equip {
		if equipmentModel.Type == req.EquipmentType {
			resp.List = append(resp.List, equipmentModel)
		}
	}
	sess.Send(protocol.SuccessData(msg.Id, resp))
}

// ChangeEquipment 切换装备 107020
func (g *Game) ChangeEquipment(sess iface.ISession, msg *codec.Message) {
	req := new(request.ChangeEquipmentRequest)
	err := jsoniter.Unmarshal(msg.GetData(), &req)
	if err != nil {
		log.Error("解析请求出错 :", err)
		sess.Send(protocol.ErrCode(msg.Id, define.MsgCode400))
		return
	}
	// {"equipmentType":1,"equipmentid":"c42ncnnt9l82rm07k61g"}
	var (
		player = g.GetPlayer(sess.GetUid())
		ne     *model.EquipmentModel
	)
	oldId := ""
	switch req.EquipmentType {
	case define.MainWeaponType,
		define.SubWeaponType,
		define.ArmorType,
		define.OrnamentType:
		var ok = false

		for _, equipmentModel := range player.useEquip {
			if equipmentModel.Type == req.EquipmentType {
				equipmentModel.Use = false
				delete(player.useEquip, equipmentModel.Id)
				oldId = equipmentModel.Id
				break
			}
		}

		ne, ok = player.Equip[req.EquipmentId]
		if !ok {
			sess.Send(protocol.ErrCode(msg.Id, define.MsgCode404))
			return
		}

		ne.Use = true
		player.useEquip[ne.Id] = ne

	default:
		sess.Send(protocol.ErrCode(msg.Id, define.MsgCode400))
		return
	}

	err = model.GetEquipmentCollection().ChangeEquipment(oldId, ne.Id)
	if err != nil {
		log.Error(err)
		sess.Send(protocol.ErrCode(msg.Id, define.MsgCode500))
		return
	}
	sess.Send(protocol.SuccessData(msg.Id, req))
	// 刷新武器属性
	g.flushAttr(player, sess)
}

// CancelEquipment 取消装备 107030
func (g *Game) CancelEquipment(sess iface.ISession, msg *codec.Message) {
	req := new(request.ChangeEquipmentRequest)
	err := jsoniter.Unmarshal(msg.GetData(), &req)
	if err != nil {
		log.Error("解析请求出错 :", err)
		sess.Send(protocol.ErrCode(msg.Id, define.MsgCode400))
		return
	}

	switch req.EquipmentType {
	case define.ArmorType, define.OrnamentType:
		err = model.GetEquipmentCollection().CancelEquipment(req.EquipmentId)
		if err != nil {
			sess.Send(protocol.ErrCode(msg.Id, define.MsgCode500))
			return
		}
	default:
		sess.Send(protocol.ErrCode(msg.Id, define.MsgCode400))
		return
	}
	// 推送
	sess.Send(protocol.Success(msg.Id))
}

// 装备升级
func (g *Game) EquipmentUpgrade(sess iface.ISession, msg *codec.Message) {
	var code = define.MsgCode200
	defer func() {
		if code != define.MsgCode200 {
			sess.Send(protocol.ErrCode(msg.Id, code))
		}
	}()
	req := new(request.EquipmentUpgradeReq)
	err := jsoniter.Unmarshal(msg.Data, &req)
	if err != nil {
		log.Error(err)
		code = define.MsgCode400
		return
	}
	player := g.GetPlayer(sess.GetUid())
	if player == nil {
		code = define.MsgCode401
		return
	}
	eqModel, ok := player.Equip[req.Mid]
	if !ok {
		code = define.MsgCode404
		return
	}
	var (
		conf     *excel.EquipagePublicConfig
		exist    bool
		cards    *model.UpgradeCard
		lv       = eqModel.Level
		start    = eqModel.Quality
		gold     = player.Account.Gold
		diamond  = player.Account.Diamond
		taskInfo []*model.TaskProgress
		taskType int
		cardNum  int
	)

	exist = false
	switch eqModel.Type {
	case define.MainWeaponType:
		conf, exist = excel.RoleMainWeaponConf[excel.GetConfigId(eqModel.RelationId, eqModel.Level)]
		taskType = define.TaskMainWeaponLv
	case define.SubWeaponType:
		conf, exist = excel.RoleSubWeaponConf[excel.GetConfigId(eqModel.RelationId, eqModel.Level)]
		taskType = define.TaskSubWeaponLv
	case define.ArmorType:
		conf, exist = excel.RoleArmorDataConf[excel.GetConfigId(eqModel.RelationId, eqModel.Level)]
		taskType = define.TaskArmorLv
	case define.OrnamentType:
		conf, exist = excel.RoleOrnamentsDataConf[excel.GetConfigId(eqModel.RelationId, eqModel.Level)]
		taskType = define.TaskOrnamentLv
	default:
		sess.Send(protocol.ErrCode(msg.Id, define.MsgCode400))
		return
	}

	if !exist {
		log.Error("没获取到装备数据表配置 table ID : %s", excel.GetConfigId(eqModel.RelationId, eqModel.Level))
		code = define.MsgCode404
		return
	}

	if req.Ty == 0 { // 提升等级
		if conf.Level == conf.MaxLevel {
			if conf.StarLevel == define.MaxStarLevel {
				// 等级 星级 皆已经满级
				code = define.MsgCode800
				return
			}
			// 提升星级
			code = define.MsgCode802
			return
		}
		if player.Account.Gold < conf.UpgradeGold {
			code = define.MsgCode604
			return
		}

		gold -= conf.UpgradeGold
		diamond -= conf.UpgradeDiamond
		lv++
		taskInfo = append(taskInfo, &model.TaskProgress{
			Type:      taskType,
			Condition: conf.RelationId,
			Lv:        lv,
		})

	} else {
		// 提升星级
		if conf.Level < conf.MaxLevel { // 先提升等级
			code = define.MsgCode801
			return
		}

		if conf.Level == conf.MaxLevel && conf.StarLevel == define.MaxStarLevel { //大佬 都已经满级
			code = define.MsgCode800
			return
		}

		if player.Account.Gold < conf.StarGold { // 穷
			code = define.MsgCode604
			return
		}

		// 材料是否足够
		packTabId, num := utils.UnmarshalItemsKV(conf.StarDeplete)
		if num > 0 {
			if pack, ok := player.packIndex[packTabId]; !ok {
				code = define.MsgCode404
				return
			} else {
				cards = pack.(*model.UpgradeCard)
				if cards.Num < num {
					code = define.MsgCode606
					return
				}
				cardNum = num
			}
		}
		gold -= conf.StarGold
		diamond -= conf.StarDiamond
		start++
		lv++
	}

	// 更新装备
	updater := bson.M{"level": lv, "quality": start}
	n, err := model.GetEquipmentCollection().UpdateOne(nil, bson.M{"_id": eqModel.Id}, updater)
	if err != nil || n == 0 {
		log.Error(err, "\r\n", n)
		code = define.MsgCode500
		return
	}
	if req.Ty == 1 && cardNum > 0 {
		cards.Num -= cardNum
		// 更新背包
		updater = bson.M{"num": cards.Num}
		n, err = model.GetBackPackColl().UpdateOne(nil, bson.M{"_id": cards.Id}, updater)
		if err != nil || n == 0 {
			log.Error(err, "\r\n", n)
			code = define.MsgCode500
			return
		}
	}
	// 更新账户金币
	if gold != player.Account.Gold || diamond != player.Account.Diamond {
		updater = bson.M{"gold": gold, "diamond": diamond}
		n, err = model.GetUserCollection().UpdateOne(nil, bson.M{"uid": player.Uid}, updater)
		if err != nil || n == 0 {
			log.Error(err, "\r\n", n)
			code = define.MsgCode500
			return
		}
		if conf.UpgradeGold != 0 {
			taskInfo = append(taskInfo, &model.TaskProgress{
				Type:      define.TaskConsumeGold,
				Condition: "",
				Num:       conf.StarGold,
			})
		}
		if conf.UpgradeDiamond != 0 {
			taskInfo = append(taskInfo, &model.TaskProgress{
				Type:      define.TaskConsumeDiamond,
				Condition: "",
				Num:       conf.StarGold,
			})
		}

	}

	// 存入库
	//err = mgoDB.GetMgo().DB().Client().UseSession(ctx, func(sessionContext mongo.SessionContext) error {
	//	sessionContext.StartTransaction()
	//	updater := bson.M{"level": eqModel.Level}
	//	if req.Ty == 1 {
	//		updater["quality"] = eqModel.Quality
	//		n, err := model.GetBackPackColl().UpdateOne(sessionContext, bson.M{"uid": uid, "tabId": itemId}, bson.M{"num": cards.Num})
	//		if err != nil || n == 0 {
	//			sessionContext.AbortTransaction(sessionContext)
	//			return fmt.Errorf("model.GetBackPackColl().UpdateOne : %v ___ n %d", err, n)
	//		}
	//	}
	//	n, err := model.GetEquipmentCollection().UpdateOne(sessionContext, bson.M{"mid": eqModel.Mid}, updater)
	//	if err != nil || n == 0 {
	//		sessionContext.AbortTransaction(sessionContext)
	//		return fmt.Errorf("model.GetEquipmentCollection().UpdateOne : %v ___ n %d", err, n)
	//	}
	//
	//	n, err = model.GetUserCollection().UpdateOne(sessionContext, bson.M{"uid": uid}, bson.M{"gold": gold})
	//	if err != nil || n == 0 {
	//		sessionContext.AbortTransaction(sessionContext)
	//		return fmt.Errorf("model.GetUserCollection().UpdateOne : %v ___ n %d", err, n)
	//	}
	//	return sessionContext.CommitTransaction(sessionContext)
	//})
	player.Account.Gold = gold
	player.Account.Diamond = diamond
	eqModel.Level = lv
	eqModel.Quality = start

	protocol.NoticeUserGoldAndDiamond(sess, player.Account.Gold, player.Account.Diamond)
	if _, ok := player.useEquip[req.Mid]; ok {
		g.flushAttr(player, sess)
	}
	sess.Send(protocol.SuccessData(msg.Id, eqModel))
	if len(taskInfo) > 0 {
		UpdateTask(player.Task,sess,taskInfo)
	}
}
