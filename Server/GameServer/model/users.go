package model

import (
	"battle_rabbit/define"
	"battle_rabbit/excel"
	"battle_rabbit/global"
	"battle_rabbit/service/log"
	"battle_rabbit/service/mgoDB"
	"battle_rabbit/service/redisDB"
	"battle_rabbit/utils"
	"battle_rabbit/utils/xid"
	"context"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"

	"math/rand"
	"strconv"
)

type PlatformType int8     // 第三方平台类型
type StageOtherType string // 关卡类型

const (
	QQPlatform     PlatformType = 1 // "QQ"
	WechatPlatform PlatformType = 2 // "wechat"
	HuaweiPlatform PlatformType = 3 // "Huawei"

	AccountStatusOk        = 0  // 正常账户
	AccountStatusDisable   = -1 // 永久封停
	AccountStatusTemporary = 1  // 暂时冻结

	DefaultUserDiamond = 0
	DefaultUserGold    = 300000
	StrengthIncTime    = 60 // s秒 每60秒增加体力1
)

var (
	userColl = newUserColl()
)

type UserColl struct {
	*mgoDB.DbBase
}

func newUserColl() *UserColl {
	return &UserColl{DbBase: mgoDB.NewDbBase(define.MgoDBNameBattle, define.TableNameUser)}
}
func GetUserCollection() *UserColl { return userColl }

// Redis Key
func UserModelRedisKey(uid int) string {
	return fmt.Sprintf("%s%d", define.RedisUserDataPrefix, uid)
}

/*
角色: 玩家主角色 可以换角色 (不同角色有不同的数据加成)
英雄: 玩家拥有的多个佣兵,进入游戏最多三个佣兵.
*/
type UserModel struct {
	Id            string            `json:"-"                bson:"_id"`
	Uid           int               `json:"uid"              bson:"uid"`          // 用户唯一的ID
	Mobile        string            `json:"mobile"           bson:"mobile"`       // 手机号
	DevID         string            `json:"-"                bson:"devId"`        // 用户登陆时的硬件ID
	Status        int8              `json:"-"                bson:"status"`       // 账户状态 0 正常  -1 永久禁用  unix时间戳表示临时禁用,并且当前时间大于这个时间戳时禁用结束
	Username      string            `json:"username"         bson:"username"`     // 用户名
	RegisterTime  int               `json:"-"                bson:"registerTime"` // 注册时间
	LoginTime     int64             `json:"-"                bson:"loginTime"`    // 登陆时间
	Avatar        string            `json:"avatar"           bson:"avatar"`       // 用户头像
	AvatarBorder  string            `json:"avatarBorder"     bson:"avatarBorder"` // 头像框
	Sex           int8              `json:"sex"              bson:"sex"`          // 性别 1 男 2 女
	Birthday      string            `json:"birthday"         bson:"birthday"`     // 生日 "1980-02-21"
	Email         string            `json:"email"            bson:"email"`        // 邮箱
	Password      string            `json:"-"                bson:"password"`     // 登陆密码
	Level         int               `json:"level"            bson:"level"`        // 账号等级
	Exp           int               `json:"exp"              bson:"exp"`          // 经验值
	Diamond       int               `json:"diamond"          bson:"diamond"`      // 钻石
	Gold          int               `json:"gold"             bson:"gold"`         // 金币
	Strength      int               `json:"strength"         bson:"strength"`     // 体力
	StrengthTime  int64             `json:"-"                bson:"strengthTime"` // 体力
	Vip           int8              `json:"vip"              bson:"vip"`          // VIP 等级
	VipExp        int               `json:"vipExp"           bson:"vipExp"`       // Vip 经验值
	VipEndTime    int               `json:"vipEndTime"       bson:"vipEndTime"`   // Vip 结束时间
	Auth          *Auth             `json:"auth"             bson:"auth"`         // 实名认证
	Platform      *Platform         `json:"-"                bson:"platform"`     // 第三方平台信息
	GuideStep     string            `json:"guideStep"        bson:"guideStep"`    // 新手引导记录
	Chest         map[string]*Chest `json:"-"                bson:"chest"`        // 开宝箱数据
	Sign          *Sign             `json:"sign"             bson:"sign"`         // 用户签到
	SyncSysMailTi int               `json:"-" bson:"syncSysMailTi"`               // 最近同步系统邮件的时间
}

//系统邮件
type SystemMail struct {
	State        int  `json:"state"`        // 邮件状态：已读 1 删除 2
	ReceiveState bool `json:"receiveState"` // 领取状态：false未领取  true 已领取
}

//开宝箱次数统计
type Chest struct {
	Count    int `json:"count" bson:"count"`       //开箱次数
	LastTime int `json:"lastTime" bson:"lastTime"` //上次免费时间
}

//用户签到
type Sign struct {
	Count    int `json:"count"    bson:"count"`    //签到次数
	LastTime int `json:"lastTime" bson:"lastTime"` //上次签到时间
}

// Auth 实名认证信息
type Auth struct {
	Name     string `json:"name"   bson:"name"`     // 真实姓名
	IdCard   string `json:"-"      bson:"idCard"`   // 证件编号
	AuthType string `json:"-"      bson:"authType"` // 证件类型
}

// Platform 第三方平台关联信息
type Platform struct {
	Type   PlatformType `json:"type"    bson:"type"`   // 平台类型
	OpenId string       `json:"openId"  bson:"openId"` // openId
}

// 获取一个默认的user对象
func (c *UserColl) NewDefeatUserObject() *UserModel {
	return &UserModel{
		Id:           xid.New().String(),
		Uid:          int(createUserId()),
		RegisterTime: int(time.Now().Unix()),
		LoginTime:    time.Now().Unix(),
		Level:        1,
		Sex:          1,
		Diamond:      DefaultUserDiamond,
		Gold:         DefaultUserGold,
		Strength:     100,
		Chest:        map[string]*Chest{},
		Auth:         new(Auth),
		Platform:     new(Platform),
		GuideStep:    "1",
		Sign:         new(Sign),
	}
}

// GetUserByDevId 通过设备号获取用户信息
func (c *UserColl) UserLoginByDevId(devId string) (*UserModel, error) {
	var user *UserModel
	err := c.FindOne(nil, bson.D{{"devId", devId}}, &user, options.FindOne().SetProjection(bson.M{"_id": 0, "uid": 1, "status": 1, "loginTime": 1}))
	if err != nil {
		if err == mongo.ErrNilDocument {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

// GetUserByDevId 通过设备号获取用户信息
func (c *UserColl) GetUserByDevId(devId string) (*UserModel, error) {
	var userModel *UserModel
	err := c.FindOne(nil, bson.D{{"devId", devId}}, &userModel)
	if err != nil {
		if err == mongo.ErrNilDocument {
			return nil, nil
		}
		return nil, err
	}
	return userModel, nil
}

// TODO 优化
func createUserId() int64 {
	_, err := redisDB.Client.Hincrby(define.RedisIncrementKeyMap, define.RedisUserIncrementKey, 1)
	if err != nil {
		return 0
	}
	userId, _ := redis.String(redisDB.Client.Hget(define.RedisIncrementKeyMap, define.RedisUserIncrementKey))
	id, _ := strconv.ParseInt(userId, 10, 64)
	return id
}

func (c *UserColl) Inc(filter bson.D, inc bson.D) error {
	_, err := c.IncOne(context.TODO(), filter, inc)
	return err
}

func (c *UserColl) IncGold(userId int, Gold int) error {
	return c.Inc(bson.D{{"uid", userId}}, bson.D{{"gold", Gold}})
}
func (c *UserColl) DecGold(userId int, Gold int) error {
	return c.Inc(bson.D{{"uid", userId}}, bson.D{{"gold", -Gold}})
}
func (c *UserColl) IncDiamond(userId int, Diamond int) error {
	return c.Inc(bson.D{{"uid", userId}}, bson.D{{"diamond", Diamond}})
}
func (c *UserColl) DecDiamond(userId int, Diamond int) error {
	return c.Inc(bson.D{{"uid", userId}}, bson.D{{"diamond", -Diamond}})
}
func (c *UserColl) IncStrength(userId int, Strength int) error {
	return c.Inc(bson.D{{"uid", userId}}, bson.D{{"strength", Strength}})
}
func (c *UserColl) DecStrength(userId int, Strength int) error {
	return c.Inc(bson.D{{"uid", userId}}, bson.D{{"strength", -Strength}})
}

// SetChestInfo 设置上一次免费开宝箱时间
func (c *UserColl) SetFreeOpenChestTime(user *UserModel, chestId string) (bool, error) {
	if user.Chest != nil {
		if chest, ok := user.Chest[chestId]; ok {
			chest.LastTime = int(time.Now().Unix())

		} else {
			user.Chest[chestId] = &Chest{
				LastTime: int(time.Now().Unix()),
			}
		}
	} else {
		user.Chest = map[string]*Chest{chestId: &Chest{
			LastTime: int(time.Now().Unix()),
		}}
	}
	if res, err := c.UpdateOne(nil, bson.D{{"uid", user.Uid}}, bson.D{{"chest." + chestId, user.Chest[chestId]}}, options.Update().SetUpsert(true)); err != nil {
		log.Error(err)
		return false, err
	} else {
		if res > 0 {
			return true, nil
		}
		return false, errors.New("设置免费开宝箱时间失败")
	}
}

// SetSignInfo 更新玩家签到
func (c *UserColl) SetSignInfo(uid int, count int) (bool, error) {
	if res, err := c.UpdateOne(nil, bson.D{{"uid", uid}}, bson.D{{"sign.count", count}, {"sign.lastTime", time.Now().Unix()}}); err != nil {
		return false, err
	} else {
		if res > 0 {
			return true, nil
		}
		return false, errors.New("更新玩家签到失败")
	}
}

// version 07
// ---------------------------------------------------------------------

func (c *UserColl) LoadUser(ctx context.Context, uid int) (*UserModel, error) {
	var user *UserModel
	err := mgoDB.GetMgoSecondary(define.MgoDBNameBattle).GetCol(c.CollName).FindOne(ctx, bson.M{"uid": uid}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return user, err

}

//刷新体力
// ok 扣除体力是否成功
// n < 0  扣除体力 | n == 0 根据当前时间刷新体力 |  n > 0 增加体力(若当前时间计算后的体力为100, 那么这个体力值还能在100上累加)
func (c *UserModel) FlushStrength(n int) (ok bool) {
	var (
		update     = false
		now        = utils.GetNowTimeStamp()
		intervalTi = now - c.StrengthTime
	)

	if intervalTi/StrengthIncTime > 0 {
		if c.Strength < 100 {
			ad := intervalTi / StrengthIncTime
			c.Strength += int(ad)
			if c.Strength > 100 {
				c.Strength = 100
			}
			c.StrengthTime += ad * StrengthIncTime
			update = true
		}else{
			c.StrengthTime = now
		}
	}

	if n != 0 && c.Strength+n >= 0 {
		ok = true
		update = true
		c.Strength += n
	}

	if update {
		userColl.UpdateOne(nil, bson.M{"uid": c.Uid}, bson.M{"strength": c.Strength, "strengthTime": c.StrengthTime})
	}
	return
}

func (tl *TalentLv) AddAttribute(rLv int, attr global.Attribute) *global.Attribute {
	//TotalLv  int  `json:"totalLv"   bson:"totalLv"`  //天赋总等级
	//Hp       int8 `json:"hp"        bson:"hp"`       //生命加成
	//Attack   int8 `json:"attack"    bson:"attack"`   //攻击加成
	//AttSpeed int8 `json:"attSpeed"  bson:"attSpeed"` //攻速加成
	//Def      int8 `json:"def"       bson:"def"`      //防御加成
	//Violence int8 `json:"violence"  bson:"violence"` //暴击加成
	//Gold     int8 `json:"gold"      bson:"gold"`     //金币加成
	//Boss     int8 `json:"boss"      bson:"boss"`     //boss伤害加成
	//Move     int8 `json:"move"      bson:"move"`     //移动速度加成
	//Dodge    int8 `json:"dodge"     bson:"dodge"`    //闪避加成
	//Buff     int8 `json:"buff"      bson:"buff"`     //buff时长加成

	if tl.Hp != 0 {
		attr.LifeC += excel.RolePassiveSkillConf[excel.GetConfigId(define.TalentHpRelationId, int(tl.Hp))].LifeC
	}
	if tl.Attack != 0 {
		attr.AttackC += excel.RolePassiveSkillConf[excel.GetConfigId(define.TalentAttackRelationId, int(tl.Attack))].AttackC
	}
	if tl.AttSpeed != 0 {
		attr.AttackSpeedC += excel.RolePassiveSkillConf[excel.GetConfigId(define.TalentAttSpeedRelationId, int(tl.AttSpeed))].AttackSpeedC
	}
	if tl.Def != 0 {
		attr.DefenseC += excel.RolePassiveSkillConf[excel.GetConfigId(define.TalentDefRelationId, int(tl.Def))].DefenseC
	}
	if tl.Violence != 0 {
		attr.CriticalC += excel.RolePassiveSkillConf[excel.GetConfigId(define.TalentViolenceRelationId, int(tl.Violence))].CriticalC
	}
	if tl.Gold != 0 {
		attr.GoldAddC += excel.RolePassiveSkillConf[excel.GetConfigId(define.TalentGoldRelationId, int(tl.Gold))].GoldAddC
	}
	if tl.Boss != 0 {
		attr.BossHurtC += excel.RolePassiveSkillConf[excel.GetConfigId(define.TalentBossRelationId, int(tl.Boss))].BossHurtC
	}
	if tl.Move != 0 {
		attr.MoveSpeedC += excel.RolePassiveSkillConf[excel.GetConfigId(define.TalentMoveRelationId, int(tl.Move))].MoveSpeedC
	}
	if tl.Dodge != 0 {
		attr.DodgeC += excel.RolePassiveSkillConf[excel.GetConfigId(define.TalentDodgeRelationId, int(tl.Dodge))].DodgeC
	}
	if tl.Buff != 0 {
		attr.BuffTimeC += excel.RolePassiveSkillConf[excel.GetConfigId(define.TalentBuffRelationId, int(tl.Buff))].BuffTimeC
	}

	expendConf := excel.RolePassiveExpendConf[rLv]
	if expendConf != nil && expendConf.LiftLv > 0 {
		attr.LifeA += excel.RolePassiveSkillConf[excel.GetConfigId(define.TalentHpAllRelationId, int(expendConf.LiftLv))].LifeA
	}
	return &attr
}

// TODO 某个等级满级后的跳过处理
func (t *TalentLv) UpgradeRoleTalentOne(conf *excel.RolePassiveExpendConfig) (lv, ty int8) {
	//	var tmp = map[int]int{1: 1, 2: 2, 3: 3, 4: 4, 5: 5, 6: 6, 7: 7, 8: 8, 9: 9, 10: 10}
	//SelectTalent:
	//	// 10个天赋技能,随机选择一个
	ty = int8(rand.Intn(10) + 1)
	//	selectTy ,ok := tmp[k]
	//
	//	if ! ok {
	//		goto SelectTalent
	//	}

	switch define.TalentType(ty) {
	case define.TalentHpType: // 生命加成类型
		t.Hp++
		lv = t.Hp
	case define.TalentAttackType: // 攻击加成类型
		t.Attack++
		lv = t.Attack
	case define.TalentAttSpeedType: // 攻速加成类型
		t.AttSpeed++
		lv = t.AttSpeed
	case define.TalentDefType: // 防御加成类型
		t.Def++
		lv = t.Def
	case define.TalentViolenceType: // 暴击加成类型
		t.Violence++
		lv = t.Violence
	case define.TalentGoldType: // 金币加成类型
		t.Gold++
		lv = t.Gold
	case define.TalentBossType: // boss伤害加成类型
		t.Boss++
		lv = t.Boss
	case define.TalentMoveType: // 移动速度加成类型
		t.Move++
		lv = t.Move
	case define.TalentDodgeType: // 闪避加成类型
		t.Dodge++
		lv = t.Dodge
	case define.TalentBuffType: // buff时长加成类型
		t.Buff++
		lv = t.Buff
	}
	t.TotalLv++
	return lv, ty
}

func (user *UserModel) FreeOpenChestFlushTime(chestId string) {

}
