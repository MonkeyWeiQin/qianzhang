package model

import (
	"battle_rabbit/define"
	"battle_rabbit/global"
	"battle_rabbit/service/mgoDB"
	"battle_rabbit/utils/xid"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	DefaultRoleTableConfigId = "role_01_1" // 账号角色初始时使用的配置表中的数据ID
)

var (
	rollColl = newRollColl()
)

type roleColl struct {
	*mgoDB.DbBase
}

func newRollColl() *roleColl {
	return &roleColl{mgoDB.NewDbBase(define.MgoDBNameBattle, define.TableNameRole)}
}
func RoleCollection() *roleColl {
	return rollColl
}

// 角色天赋等级数据
// 都记录的是数据表中的等级值
type TalentLv struct {
	TotalLv  int  `json:"totalLv"   bson:"totalLv"`  //天赋总等级
	Hp       int8 `json:"hp"        bson:"hp"`       //生命加成
	Attack   int8 `json:"attack"    bson:"attack"`   //攻击加成
	AttSpeed int8 `json:"attSpeed"  bson:"attSpeed"` //攻速加成
	Def      int8 `json:"def"       bson:"def"`      //防御加成
	Violence int8 `json:"violence"  bson:"violence"` //暴击加成
	Gold     int8 `json:"gold"      bson:"gold"`     //金币加成
	Boss     int8 `json:"boss"      bson:"boss"`     //boss伤害加成
	Move     int8 `json:"move"      bson:"move"`     //移动速度加成
	Dodge    int8 `json:"dodge"     bson:"dodge"`    //闪避加成
	Buff     int8 `json:"buff"      bson:"buff"`     //buff时长加成
}

// 主角属性
type RoleModel struct {
	Id         string            `json:"id"         bson:"_id"`
	Uid        int               `json:"uid"        bson:"uid"`
	RLv        int               `json:"rlv"        bson:"rlv"`        // 角色等级
	RExp       int               `json:"rExp"       bson:"rExp"`       // 角色经验值
	Quality    int               `json:"quality"    bson:"quality"`    // 星级 初始为1星 最高可升至6星
	RelationId string            `json:"relationId" bson:"relationId"` // 数据表关联ID
	SkinId     string            `json:"skinId"     bson:"skinId"`     // 角色皮肤ID
	TalentLv   *TalentLv         `json:"talentLv"   bson:"talentLv"`   // 天赋等级
	Attribute  *global.Attribute `json:"attribute"  bson:"attribute"`  // 角色属性
	Skins      []string          `json:"skins"      bson:"skins"`      // 角色皮肤集合
}

// 创建默认的角色
func NewDefaultRole(uid int) *RoleModel {

	return &RoleModel{
		Id:         xid.New().String(),
		Uid:        uid,
		RLv:        1,
		Quality:    1,
		RelationId: "role_01",
		SkinId:     "skin_1",
		TalentLv:   new(TalentLv),
		Attribute:  new(global.Attribute),
		Skins:      []string{"skin_1"},
	}
}

func (rc *roleColl) LoadRoleFromDB(uid int) (*RoleModel, error) {
	role := new(RoleModel)
	err := mgoDB.GetMgoSecondary(define.MgoDBNameBattle).GetCol(rc.CollectionName()).FindOne(nil, bson.M{"uid": uid}).Decode(role)
	if err != nil {
		return nil, err
	}
	return role, nil
}
