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
	"time"
)

// 获取角色属性以及4种装备
func (g *Game) GetRoleInfo(sess iface.ISession, msg *codec.Message) {
	var code = define.MsgCode200
	defer func() {
		if code != define.MsgCode200 {
			sess.Send(protocol.ErrCode(msg.Id, code))
		}
	}()
	resp := new(response.RoleInfoResp)
	player := g.GetPlayer(sess.GetUid())
	if player == nil {
		code = define.MsgCode401
		return
	}

	role := *player.Role
	resp.Role = &role

	for _, equipment := range player.useEquip {
		if equipment == nil {
			continue
		}
		switch equipment.Type {
		case define.ArmorType: //护甲的属性
			resp.Armor = equipment
		case define.OrnamentType: //饰品的属性
			resp.Ornaments = equipment
		case define.MainWeaponType:
			resp.WeaponMain = equipment
		case define.SubWeaponType:
			resp.WeaponSub = equipment
		}
	}
	if player.roleAttr == nil {
		err := player.TotalAttribute()
		if err != nil {
			log.Error(err)
			code = define.MsgCode500
			return
		}
	}
	resp.Role.Attribute = player.roleAttr
	sess.Send(protocol.SuccessData(msg.Id, resp))
}

// 获取角色天赋等级
func (g *Game) GetRoleTalentLv(sess iface.ISession, msg *codec.Message) {
	player := g.GetPlayer(sess.GetUid())
	if player == nil {
		sess.Send(protocol.ErrCode(msg.Id, define.MsgCode401))
		return
	}
	sess.Send(protocol.SuccessData(msg.Id, player.Role.TalentLv))
}

// 角色天赋升级
func (g *Game) UpgradeRoleTalentLv(sess iface.ISession, msg *codec.Message) {
	var code = define.MsgCode200
	defer func() {
		if code != define.MsgCode200 {
			sess.Send(protocol.ErrCode(msg.Id, code))
		}
	}()

	req := new(request.UpgradeRoleTalentLvReq)
	err := jsoniter.Unmarshal(msg.GetData(), req)
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

	tl := player.Role.TalentLv
	conf := excel.RolePassiveExpendConf[tl.TotalLv+1]
	if player.Account.Gold < conf.Gold {
		code = define.MsgCode604
		return
	}
	if conf.Gold == 0 {
		code = define.MsgCode800
		return
	}

	//升级
	lv, ty := tl.UpgradeRoleTalentOne(conf)
	// 扣金币
	player.Account.Gold -= conf.Gold

	// mongo 更新
	_, err = model.GetUserCollection().UpdateOne(nil, bson.M{"uid": sess.GetUid()}, bson.M{"gold": player.Account.Gold})
	if err != nil {
		log.Error(err)
		code = define.MsgCode500
		return
	}
	_, err = model.RoleCollection().UpdateOne(nil, bson.M{"uid": sess.GetUid()}, bson.M{"talentLv": tl})
	if err != nil {
		log.Error(err)
		code = define.MsgCode500
		return
	}
	resp := &response.GetRoleTalentLvResp{
		Lv:      lv,
		Ty:      ty,
		TotalLv: tl.TotalLv,
	}
	sess.Send(protocol.SuccessData(msg.Id, resp))
	// 推送金币变化
	protocol.NoticeUserGoldAndDiamond(sess, player.Account.Gold, player.Account.Diamond)
	// 刷新属性并推送
	g.flushAttr(player, sess)

}

// 更换角色皮肤
func (g *Game) ChangeSkin(sess iface.ISession, msg *codec.Message) {
	var req *request.ChangeRoleSkinReq
	err := jsoniter.Unmarshal(msg.GetData(), &req)
	if err != nil {
		sess.Send(protocol.ErrCode(msg.Id, define.MsgCode400))
		return
	}

	player := g.GetPlayer(sess.GetUid())
	if player == nil {
		sess.Send(protocol.ErrCode(msg.Id, define.MsgCode401))
		return
	}
	change := false
	for _, skinId := range player.Role.Skins {
		if skinId == req.SkinId {
			player.Role.SkinId = skinId
			change = true
			break
		}
	}

	if !change {
		sess.Send(protocol.ErrCode(msg.Id, define.MsgCode404))
		return
	}
	_, err = model.RoleCollection().UpdateOne(nil, bson.M{"uid": player.Uid}, bson.M{"skinId": req.SkinId})
	if err != nil {
		sess.Send(protocol.Err(msg.Id))
		return
	}
	sess.Send(protocol.Success(msg.Id))
	g.flushAttr(player, sess)
}

// 刷新体力值
func (g *Game) FlushStrength(sess iface.ISession, msg *codec.Message) {
	player := g.GetPlayer(sess.GetUid())
	if player == nil {
		sess.Send(protocol.ErrCode(msg.Id, define.MsgCode401))
		return
	}
	ok := player.Account.FlushStrength(0)
	if !ok {
		sess.Send(protocol.ErrCode(msg.Id, define.MsgCode500))
		return
	}
	resp := protocol.SuccessData(msg.Id, &response.FlushStrengthResp{
		Strength: player.Account.Strength,
		Time:     int(time.Now().Unix()),
	})
	sess.Send(resp)
}

// 角色升星
func (g *Game) UpgradeRoleStar(sess iface.ISession, msg *codec.Message) {
	var code = define.MsgCode200
	defer func() {
		if code != define.MsgCode200 {
			sess.Send(protocol.ErrCode(msg.Id, code))
		}
	}()

	player := g.GetPlayer(sess.GetUid())
	if player == nil {
		code = define.MsgCode401
		return
	}
	var (
		gold    = player.Account.Gold
		diamond = player.Account.Diamond
		cards   *model.UpgradeCard
		role    = player.Role
	)

	conf, ok := excel.RoleDataConf[excel.GetConfigId(role.RelationId, role.RLv)]
	if !ok {
		code = define.MsgCode500
		return
	}
	if conf.NextId == "0" { // 满级
		code = define.MsgCode800
		return
	}

	if conf.StarGold != 0 {
		if conf.StarGold > gold {
			code = define.MsgCode604
			return
		} else {
			gold -= conf.StarGold
		}
	}
	if conf.StarDiamond != 0 {
		if conf.StarDiamond > diamond {
			code = define.MsgCode605
			return
		} else {
			diamond -= conf.UpDiamond
		}
	}

	nextConf := excel.RoleDataConf[conf.NextId]

	if nextConf.StarLevel != conf.StarLevel+1 {

		/* ----------debug todo 默认跳转到下一个升星等级 ------------- */
		//if nextConf == nil || conf.NextId == "0" {
		//	return reply.ErrCode(args.Msg.GetMsgId(), define.MsgCode800)
		//}
		//for nextConf.StarMaterials == "0" {
		//	nextConf = excel.RoleDataConf[nextConf.NextId]
		//	if nextConf == nil || conf.NextId == "0" {
		//		return reply.ErrCode(args.Msg.GetMsgId(), define.MsgCode800)
		//	}
		//}
		//_, err = redisDB.Client.Hmset(model.DefaultUserModelRedisKey(uid), "RLv", nextConf.Level, "Quality", nextConf.StarLevel)
		//if err != nil {
		//	log.Error(err)
		//	return reply.ErrCode(args.Msg.GetMsgId(), define.MsgCode500)
		//}
		//
		//n, err := model.GetUserCollection().UpdateOne(ctx, bson.M{"uid": uid}, bson.M{"rLv": nextConf.Level, "quality": nextConf.StarLevel})
		//if err != nil || n == 0 {
		//	log.Error(err)
		//	return reply.ErrCode(args.Msg.GetMsgId(), define.MsgCode500)
		//}
		//conf = nextConf
		//nextConf = excel.RoleDataConf[conf.NextId]
		//log.Debug("升星测试:: 提升从 %d 级提升到 %d 级 进行升星测试! ", rlv, nextConf.Level)

		/*---------------debug end --------------------------------------------------------*/

		// 等级不够,不能升星
		code = define.MsgCode801
		return
	}

	packTabId, num := utils.UnmarshalItemsKV(conf.StarMaterials)
	if packTabId == "" || num == 0 { // 升星卡是必须物品
		code = define.MsgCode404
		return
	}

	for _, pack := range player.Package {
		if pack.GetType() == model.PackageCardType {
			cards = pack.(*model.UpgradeCard)
			if cards.TabId == packTabId {
				break
			}
			cards = nil
		}
	}

	if cards == nil {
		log.Error("没有升星卡: id: %s", packTabId)
		code = define.MsgCode404
		return
	}
	if cards.Num < num {
		code = define.MsgCode606
		return
	} else {
		cards.Num -= num
	}
	// 更新星级,等级,玩家属性
	// 升星操作,玩家星级+1 等级+1 扣除相应升星卡,金币,钻石等消耗
	role.RLv = nextConf.Level
	role.RExp = 0
	role.Quality = nextConf.StarLevel

	player.Account.Gold = gold

	update := bson.M{"rlv": role.RLv, "quality": role.Quality, "rExp": 0}
	//if conf.StarDiamond != 0 {
	//	update["diamond"] = diamond
	//}

	_, err := model.RoleCollection().UpdateOne(nil, bson.M{"uid": player.Uid}, update)
	if err != nil {
		log.Error(err)
		code = define.MsgCode500
		return
	}

	_, err = model.GetUserCollection().UpdateOne(nil, bson.M{"uid": player.Uid}, bson.M{"gold": player.Account.Gold})
	if err != nil {
		log.Error(err)
		code = define.MsgCode500
		return
	}

	_, err = model.GetBackPackColl().UpdateOne(nil, bson.M{"_id": cards.Id}, bson.M{"num": cards.Num})
	if err != nil {
		log.Error(err)
		code = define.MsgCode500
		return
	}
	//err = mgoDB.GetMgo().Cli().UseSession(ctx, func(sessCtx mongo.SessionContext) error {
	//	if err := sessCtx.StartTransaction(); err != nil {
	//		return err
	//	}
	//
	//	defer sessCtx.EndSession(sessCtx)
	//	_, err := model.GetBackPackColl().IncOne(sessCtx, bson.M{"uid": uid, "ty": define.CardRoleStarType}, bson.M{"num": -num})
	//	if err != nil {
	//		sessCtx.AbortTransaction(sessCtx)
	//		return err
	//	}
	//
	//	_, err = model.GetUserCollection().UpdateOne(sessCtx, bson.M{"uid": uid}, update)
	//	if err != nil {
	//		sessCtx.AbortTransaction(sessCtx)
	//		return err
	//	}
	//	return sessCtx.CommitTransaction(sessCtx)
	//})
	//
	//if err != nil {
	//	log.Error(err)
	//	return reply.ErrCode(args.Msg.GetMsgId(), define.MsgCode500)
	//}

	// 升级成功
	// 刷新战力并推送 TODO

	// 推送升级后的等级和星级,以及消耗
	resp := &response.UpgradeRoleStarResp{
		RLv:        nextConf.Level,
		Quality:    int8(nextConf.StarLevel),
		RelationId: role.RelationId,
	}
	sess.Send(protocol.SuccessData(msg.Id, resp))
	// 刷新属性并推送
	g.flushAttr(player, sess)
	// 推送金币变化
	protocol.NoticeUserGoldAndDiamond(sess, player.Account.Gold, player.Account.Gold)
}

// 刷新账号角色,装备的属性
func (g *Game) FlushUserAttributePush(sess iface.ISession, msg *codec.Message) {
	var req *request.FlushUserAttributePushData
	err := jsoniter.Unmarshal(msg.Data, &req)
	if err != nil {
		sess.Send(protocol.ErrCode(msg.Id, define.MsgCode500))
		return
	}
	player := g.GetPlayer(sess.GetUid())
	if player == nil {
		sess.Send(protocol.ErrCode(msg.Id, define.MsgCode401))
		return
	}
	g.flushAttr(player, sess)
	return
}

func (g *Game) flushAttr(player *Player, sess iface.ISession) {
	 err := player.TotalAttribute()
	if err != nil {
		log.Error(err)
		return
	}
	if player.subAttr != nil {
		protocol.NoticeWeaponAttribute(sess, define.SubWeaponType, player.subAttr)
	}
	//if req.Ty>>1&1 == 1 {
	protocol.NoticeWeaponAttribute(sess, define.MainWeaponType, player.mainAttr)
	//}
	//if req.Ty>>2&1 == 1 {
	protocol.NoticeRoleAttribute(sess, player.roleAttr)
}
