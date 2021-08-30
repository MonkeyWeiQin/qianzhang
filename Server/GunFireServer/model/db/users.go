package db

import (
	"com.xv.admin.server/service/mgoDB"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	//"time"
)

type PlatformType int8 // 第三方平台类型

type UserModel struct {
	Id           primitive.ObjectID                  `json:"-"                bson:"-"         redis:"-"`
	Uid          int                                 `json:"uid"              bson:"uid"`           // 用户唯一的ID
	Mobile       string                              `json:"mobile"           bson:"mobile"`        // 手机号
	DevID        string                              `json:"dev_id"           bson:"dev_id"`        // 用户登陆时的硬件ID
	Status       int8                                `json:"status"           bson:"status"`        // 账户状态 0 正常  -1 永久禁用  unix时间戳表示临时禁用,并且当前时间大于这个时间戳时禁用结束
	Username     string                              `json:"username"         bson:"username"`      // 用户名
	RegisterTime int                                 `json:"-"                bson:"register_time"` // 注册时间
	LoginTime    int                                 `json:"-"                bson:"login_time"`    // 登陆时间
	Avatar       string                              `json:"avatar"           bson:"avatar"`        // 用户头像
	AvatarBorder string                              `json:"avatar_border"    bson:"avatar_border"` // 头像框
	Sex          int8                                `json:"sex"              bson:"sex"`           // 性别 1 男 2 女
	Birthday     string                              `json:"birthday"         bson:"birthday"`      // 生日 "1980-02-21"
	Email        string                              `json:"email"            bson:"email"`         // 邮箱
	Password     string                              `json:"-"                bson:"password"`      // 登陆密码
	Level        int8                                `json:"level"            bson:"level"`         // 账号等级
	Exp          int                                 `json:"exp"              bson:"exp"`           // 经验值
	Diamond      int                                 `json:"diamond"          bson:"diamond"`       // 钻石
	Gold         int                                 `json:"gold"             bson:"gold"`          // 金币
	Strength     int8                                `json:"strength"         bson:"strength"`      // 体力
	Vip          int8                                `json:"vip"              bson:"vip"`           // VIP 等级
	VipExp       int                                 `json:"vip_exp"          bson:"vip_exp"`       // Vip 经验值
	VipEndTime   int                                 `json:"vip_end_time"     bson:"vip_end_time"`  // Vip 结束时间
	//HeroList     model.SliceString                   `json:"hero_list" bson:"hero_list"`            // 使用的英雄列表
	//TalentLv     *TalentLv                           `json:"talent_lv" bson:"talent_lv"`            // 天赋等级
	EquipageInfo *EquipageInfo                       `json:"equipage_info" bson:"equipage_info"`    // 装备数据
	*Stage       `json:"stage"      bson:"stage"`    // 当前关卡
	*Auth        `json:"auth"       bson:"auth"`     // 实名认证
	*Platform    `json:"-"          bson:"platform"` // 第三方平台信息
	*Role        `json:"role"       bson:"role"`     // 玩家角色数据
}
// Stage 当前关卡
type Stage struct {
	GeneralStage   string `json:"general_stage"   bson:"general_stage"`   //普通主线任务
	DifficultStage string `json:"difficult_stage" bson:"difficult_stage"` //困难主线任务
	StageType      int8   `json:"stage_type"      bson:"stage_type"`      //当前选择的关卡模式  0:简单 1:困难
}

// Auth 实名认证信息
type Auth struct {
	Name     string `json:"name"      bson:"name"`      // 真实姓名
	IdCard   string `json:"id_card"   bson:"id_card"`   // 证件编号
	AuthType string `json:"auth_type" bson:"auth_type"` // 证件类型
}

// Platform 第三方平台关联信息
type Platform struct {
	Type   PlatformType `json:"type"    bson:"type"`    // 平台类型
	OpenId string       `json:"open_id" bson:"open_id"` // openId
}

type EquipageInfo struct { // 当前使用的装备数据
	MainWeaponId      string `json:"main_weapon_id"      bson:"main_weapon_id"`      // 主武器
	SecondaryWeaponId string `json:"secondary_weapon_id" bson:"secondary_weapon_id"` // 副武器
	EquipmentId       string `json:"equipment_id"        bson:"equipment_id"`        // 部件
	OrnamentsId       string `json:"ornaments_id"        bson:"ornaments_id"`        // 饰品ID
}

// 主角属性
type Role struct {
	RLv        int               `json:"rlv"         bson:"rlv"`         // 角色等级
	RExp       int               `json:"r_exp"       bson:"r_exp"`       // 角色经验值
	Quality    int8              `json:"quality"     bson:"quality"`     // 星级 初始为1星 最高可升至6星
	RelationId string            `json:"relation_id" bson:"relation_id"` // 数据表关联ID
	SkinId     string            `json:"skin_id"     bson:"skin_id"`     // 角色皮肤ID
	//Attribute  *global.Attribute `json:"attribute"   bson:"attribute"`   // 角色属性
}

// 都记录的是数据表中的等级值
type TalentLv struct {
	TotalLv  int  `json:"total_lv"  bson:"total_lv"`  //天赋总等级
	Hp       int8 `json:"hp"        bson:"hp"`        //生命加成
	Attack   int8 `json:"attack"    bson:"attack"`    //攻击加成
	AttSpeed int8 `json:"att_speed" bson:"att_speed"` //攻速加成
	Def      int8 `json:"def"       bson:"def"`       //防御加成
	Violence int8 `json:"violence"  bson:"violence"`  //暴击加成
	Gold     int8 `json:"gold"      bson:"gold"`      //金币加成
	Boss     int8 `json:"boss"      bson:"boss"`      //boss伤害加成
	Move     int8 `json:"move"      bson:"move"`      //移动速度加成
	Dodge    int8 `json:"dodge"     bson:"dodge"`     //闪避加成
	Buff     int8 `json:"buff"      bson:"buff"`      //buff时长加成
}
func DefaultUserModelRedisKey(uid int) string { return fmt.Sprintf("%s_%d", REDIS_USER_DATA_PREFIX, uid) }


func (user *UserModel) CollectionName() string {
	//suffix := fmt.Sprintf("_%d_%d", time.Now().Add(time.Hour*-6).Month(), time.Now().Add(time.Hour*-6).Day())
	//return config.CollUsers + suffix
	return "users"
}

func (user *UserModel) GetId() primitive.ObjectID {
	return user.Id
}

func (user *UserModel) SetId(id primitive.ObjectID) {
	user.Id = id
}


func (user *UserModel) GetUsersList(filter bson.D, opt *options.FindOptions) (data []map[string]interface{}, err error) {
	data = []map[string]interface{}{}
	finder := mgoDB.NewFinder(user).Where(filter).Options(opt).Records(&data)
	err = mgoDB.GetMgo().FindMany(context.TODO(), finder)
	return
}

func (user *UserModel) GetTotalNum(filter bson.D) (int64, error) {
	finder := mgoDB.NewCounter(user).Where(filter)
	return mgoDB.GetMgo().CountDocuments(nil, finder)
}
func (user *UserModel) UpdateUser(filter bson.D, update bson.D, opt *options.UpdateOptions) error {
	updater := mgoDB.NewUpdater(user).Where(filter).Update(update).Options(opt)
	_, err := mgoDB.GetMgo().UpdateOne(nil, updater)
	return err
}
// GetUser 获取用户
func (user *UserModel) GetUser(filter bson.D) (data *UserModel, err error) {
	finder := mgoDB.NewOneFinder(user).Where(filter).Record(&data)
	_ , err = mgoDB.GetMgo().FindOne(context.TODO(), finder)
	return
}

func (user *UserModel) Inc(filter bson.D, inc bson.D) error {
	updater := mgoDB.NewUpdater(user).Where(filter).Inc(inc)
	_, err := mgoDB.GetMgo().IncOne(nil, updater)
	return err
}
