package v1

import (
	"battle_rabbit/define"
	"battle_rabbit/excel"
	"battle_rabbit/model"
	"battle_rabbit/protocol/request"
	"battle_rabbit/protocol/response"
	"battle_rabbit/server/login/middleware"
	"battle_rabbit/server/login/proto"
	"battle_rabbit/service/log"
	"battle_rabbit/service/mgoDB"
	"battle_rabbit/utils/xid"
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"strconv"
	"strings"
	"time"
)

const (
	GateAddr = "127.0.0.1:8854"
)

// 设备账号注册或者登录
func LoginByDevId(c *gin.Context) {
	req := new(request.PlayerDevIdLoginReq)
	err := c.ShouldBind(req)
	if err != nil {
		log.Error(err)
		proto.ResponseFailWithCode(c, define.MsgCode400)
		return
	}
	// 注册或者登录
	user, err := model.GetUserCollection().UserLoginByDevId(req.DevId)
	if err != nil && err != mongo.ErrNoDocuments {
		log.Error(err)
		proto.ResponseFailWithCode(c, define.MsgCode402)
		return
	}

	if user == nil {
		// 创建用户
		user, err = CreateUsersByDevId(c, req.DevId)
		if err != nil {
			log.Error(err)
			proto.ResponseFailWithCode(c, define.MsgCode402)
			return
		}
		time.Sleep(time.Second * 2)
	} else {
		if user.Status != model.AccountStatusOk {
			if user.Status == model.AccountStatusDisable {
				proto.ResponseFailWithCode(c, define.MsgCode700)
			} else if user.Status == model.AccountStatusDisable {
				proto.ResponseFailWithData(c, define.MsgCode701, user.Status)
			}
			return
		}
		//var m = make(map[string]interface{})
		//b,err := model.GetStageColl().FindOne(c,bson.M{"uid":uid},&m)
		//if err != nil {
		//	return
		//}
		//if !b {
		//	err = model.GetStageColl().InsertOne(c,model.NewStageModel(user.Uid))
		//}

	}

	//err = userDataPushCache(user)
	//if err != nil {
	//	log.Error(err)
	//	proto.ResponseFailWithCode(c, define.MsgCode402)
	//	return
	//}

	token, err := createToken(user.Uid)
	if err != nil {
		proto.ResponseFailWithCode(c, define.MsgCode402)
		return
	}
	syncSystemMail(c, user.Uid, user.LoginTime)

	proto.ResponseOkWithData(c, &response.WebUserLoginByDevIdResp{Token: token, TcpHost: GateAddr})
}

// CreateToken 创建token
func createToken(userId int) (token string, err error) {
	clams := jwt.StandardClaims{
		NotBefore: time.Now().Unix() - 10,                // 签名生效时间
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // 过期时间24小时
		Issuer:    "qmPlus",                              // 签名的发行者
		Id:        strconv.Itoa(userId),
	}
	token, err = middleware.JwtObj.CreateToken(clams)
	return
}

// 通过devId创建用户
func CreateUsersByDevId(ctx context.Context, devId string) (*model.UserModel, error) {
	user := model.GetUserCollection().NewDefeatUserObject()

	if strings.Contains(devId,"abab1212") {
		/* ------------ debug 默认全部主 副 武器.全英雄------------------------------ */
		DebugCreate(user.Uid)
		// 各种升级卡,升星卡,强化卡,装备卡,经验卡 各100张
		DebugCreateUpgradeCard(user.Uid)
		/*--------- DEBUG END---------------------------------------------------------------------------*/
		// 创建玩家关卡数据
		stage := model.NewStageModel(user.Uid)
		stage.StageMain.GeneralStage = 100010
		// 创建玩家关卡数据
		model.GetStageColl().InsertOne(ctx, stage)
		// 角色数据
		role := model.NewDefaultRole(user.Uid)
		role.Skins = []string{"skin_1","skin_2","skin_3"}
		_, err := mgoDB.GetMgo(define.MgoDBNameBattle).GetCol(model.RoleCollection().CollectionName()).InsertOne(nil, role)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		log.Debug("角色创建完成!! : %d 个", role.Id)
	} else {

		//创建默认的主武器
		wmConfig := excel.RoleMainWeaponConf[model.DefeatWeaponMainConfigId]
		wm := model.NewEquipmentModel(user.Uid, define.MainWeaponType)
		wm.Use = true
		wm.RelationId = wmConfig.RelationId
		wm.Index = wmConfig.Index
		err := model.GetEquipmentCollection().InsertOne(ctx, wm)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		// 创建玩家关卡数据
		model.GetStageColl().InsertOne(ctx, model.NewStageModel(user.Uid))
		// 角色数据
		role := model.NewDefaultRole(user.Uid)
		_, err = mgoDB.GetMgo(define.MgoDBNameBattle).GetCol(model.RoleCollection().CollectionName()).InsertOne(nil, role)
		if err != nil {
			log.Error(err)
			return nil, err
		}
	}
	user.DevID = devId
	user.Username = "user_" + devId[:4]
	err := model.GetUserCollection().InsertOne(ctx, user)

	return user, err
}

// CheckToken 检查Token是否合法
func CheckToken(c *gin.Context) {
	uid, ok := c.Get("userId")
	if uid, _ok := uid.(int); ok && _ok {
		proto.ResponseOkWithData(c, uid)
		return
	}
	proto.ResponseFailWithCode(c, define.MsgCode400)
}

func DebugCreate(uid int) (r1 []*model.EquipmentModel, r3 *model.HeroModel) {
	// 创建主武器
	m := map[string]*excel.EquipagePublicConfig{}
	for _, config := range excel.RoleMainWeaponConf {
		if _, ok := m[config.RelationId]; !ok {
			m[config.RelationId] = config
		}
	}
	var datas = []interface{}{}
	var i = 0
	use := false
	for s, conf := range m {
		if conf.RelationId == "weapon_01" {
			use = true
		}
		datas = append(datas,
			&model.EquipmentModel{
				Id:         xid.New().String(),
				Uid:        uid,
				Type:       define.MainWeaponType,
				RelationId: s,
				Level:      1,
				Quality:    1,
				Index:      conf.Index,
				Use:        use,
			})
		use = false
	}
	r1 = append(r1, datas[0].(*model.EquipmentModel))

	i, err := model.GetEquipmentCollection().InsertMany(nil, datas)
	if err != nil {
		log.Error(err)
		return
	}
	log.Debug("生成主武器: %d 个", i)

	m = map[string]*excel.EquipagePublicConfig{}
	datas = []interface{}{}
	for _, config := range excel.RoleSubWeaponConf {
		if _, ok := m[config.RelationId]; !ok {
			m[config.RelationId] = config
		}
	}
	i = 0
	use = false
	for s, conf := range m {
		if conf.RelationId == "sub_01" {
			use = true
		}
		datas = append(datas,
			&model.EquipmentModel{
				Id:         xid.New().String(),
				Uid:        uid,
				Type:       define.SubWeaponType,
				RelationId: s,
				Level:      1,
				Index:      conf.Index,
				Quality:    1,
				Use:        use,
			})
		use = false
	}
	r1 = append(r1, datas[0].(*model.EquipmentModel))
	i, err = model.GetEquipmentCollection().InsertMany(nil, datas)
	if err != nil {
		log.Error(err)
		return
	}
	log.Debug("生成副武器: %d 个", i)

	// 护甲
	datas = []interface{}{}
	use = false
	m = map[string]*excel.EquipagePublicConfig{}
	datas = []interface{}{}
	for _, config := range excel.RoleArmorDataConf {
		if _, ok := m[config.RelationId]; !ok {
			m[config.RelationId] = config
		}
	}

	for _, armor := range m {
		if armor.RelationId == "armor_01" {
			use = true
		}
		datas = append(datas,
			&model.EquipmentModel{
				Id:         xid.New().String(),
				Uid:        uid,
				Type:       define.ArmorType,
				RelationId: armor.RelationId,
				Level:      1,
				Index:      armor.Index,
				Quality:    1,
				Use:        use,
			})
		use = false
	}
	r1 = append(r1, datas[0].(*model.EquipmentModel))
	i, err = model.GetEquipmentCollection().InsertMany(nil, datas)
	if err != nil {
		log.Error(err)
		return
	}
	log.Debug("生成护甲: %d 个", i)
	//饰品
	use = false
	m = map[string]*excel.EquipagePublicConfig{}
	datas = []interface{}{}
	for _, config := range excel.RoleOrnamentsDataConf {
		if _, ok := m[config.RelationId]; !ok {
			m[config.RelationId] = config
		}
	}
	for _, ornaments := range m {
		if ornaments.RelationId == "jewel_01" {
			use = true
		}
		datas = append(datas,
			&model.EquipmentModel{
				Id:         xid.New().String(),
				Uid:        uid,
				Type:       define.OrnamentType,
				RelationId: ornaments.RelationId,
				Level:      1,
				Index:      ornaments.Index,
				Quality:    1,
				Use:        use,
			})
		use = false
	}
	r1 = append(r1, datas[0].(*model.EquipmentModel))
	i, err = model.GetEquipmentCollection().InsertMany(nil, datas)
	if err != nil {
		log.Error(err)
		return
	}
	log.Debug("生成饰品: %d 个", i)

	// 英雄
	heroMap := map[string]*excel.HeroDataConfig{}
	datas = []interface{}{}
	for _, config := range excel.HeroDataConf {
		if _, ok := heroMap[config.RelationId]; !ok {
			heroMap[config.RelationId] = config
		}
	}
	for s, conf := range heroMap {
		if conf.RelationId == "battery_01" || conf.RelationId == "battery_02" || conf.RelationId == "battery_03" {
			continue
		}
		var passSkillIds []string
		RIds := []string{conf.PSId1, conf.PSId2, conf.PSId3, conf.PSId4}
		for _, id := range RIds {
			if psConf, ok := excel.HeroPassiveSkillConf[excel.GetConfigId(id, 1)]; ok {
				passSkillIds = append(passSkillIds, psConf.Id)
			}
		}
		datas = append(datas,
			&model.HeroModel{
				Id:         xid.New().String(),
				Uid:        uid,
				RelationId: s,
				Use:        false,
				Rarity:     conf.Rarity,
				Lv:         1,
				StarLv:     1,
				Index:      conf.Index,
				TotalPower: 0, // TODO 还需添加默认总战力,目前先为0
				ASId:       excel.HeroActiveSkillConf[excel.GetConfigId(conf.ASId, 1)].Id,
				PsIds:      passSkillIds,
				//WeaponId:   excel.HeroWeaponDataConf[conf.WeaponId].Id,
			})
	}
	r3 = datas[0].(*model.HeroModel)
	i, err = model.GetHeroCollection().InsertMany(nil, datas)
	if err != nil {
		log.Error(err)
		return
	}
	log.Debug("生成英雄: %d 个", i)

	return
}

func DebugCreateUpgradeCard(uid int) {
	var packList []interface{}
	for tid, cardConf := range excel.CardDataConf {
		//ty := 1
		//if cardConf.Ty == 10 {
		//	ty = 2
		//}
		cardModel := &model.UpgradeCard{
			Id:    xid.New().String(),
			Uid:   uid,
			TabId: tid,
			Ty:    define.PackageType(cardConf.Ty),
			//Val:    cardConf.Exp,
			//Rarity: cardConf.Rare,
			Num: 10000,
		}
		packList = append(packList, cardModel)
	}
	_, err := model.GetBackPackColl().InsertMany(context.TODO(), packList)
	if err != nil {
		log.Error(err)
	}
	return
}

// 将系统邮件同步到个人账户中
func syncSystemMail(ctx context.Context, uid int, lastSyncTi int64) {
	filter := bson.M{"uid": -1, "sendTime": bson.M{"$gt": lastSyncTi}, "validity": bson.M{"$gt": time.Now().Unix()}}
	var recode []*model.MailModel
	err := model.GetMailCollection().FindAll(ctx, filter, &recode)
	if err != nil {
		log.Error(err)
		return
	}

	if len(recode) > 0 {
		var cols []interface{}
		for _, mailModel := range recode {
			mailModel.UserId = uid
			cols = append(cols, mailModel)
		}
		n, err := model.GetMailCollection().InsertMany(ctx, cols)
		if err != nil {
			log.Error(err)
		}
		log.Debug("同步系统邮件到个人账户中! 数量:%d", n)
	}
}
