package response

import (
	"battle_rabbit/global"
	"battle_rabbit/model"
)

type JoinGameRoleResp struct {
	RLv        int               `json:"rlv"        bson:"rlv"`        // 角色等级
	RExp       int               `json:"rExp"       bson:"rExp"`       // 角色经验值
	Quality    int               `json:"quality"    bson:"quality"`    // 星级 初始为1星 最高可升至6星
	RelationId string            `json:"relationId" bson:"relationId"` // 数据表关联ID
	SkinId     string            `json:"skinId"     bson:"skinId"`     // 角色皮肤ID
	Attribute  *global.Attribute `json:"attribute"`
}

type JoinGameWeaponMainResp struct {
	Level      int               `json:"level"       bson:"level"`      // 角色等级
	RelationId string            `json:"relationId"  bson:"relationId"` // 数据表关联ID
	Attribute  *global.Attribute `json:"attribute"`
}

type JoinGameWeaponSubResp struct {
	Level      int               `json:"level"       bson:"level"`      // 角色等级
	RelationId string            `json:"relationId"  bson:"relationId"` // 数据表关联ID
	Attribute  *global.Attribute `json:"attribute"`
}

type JoinGameHeroResp struct {
	Mid           string            `json:"mid"           bson:"mid"`           //英雄ID
	RelationId    string            `json:"relationId"    bson:"relationId"`    //数据表关联ID
	Lv            int               `json:"level"         bson:"level"`         //等级
	StarLv        int               `json:"starLv"        bson:"starLv"`        //星级
	Attribute     *global.Attribute `json:"attribute"     bson:"attribute"`     // 属性
	ActiveSkillLv int               `json:"activeSkillLv" bson:"activeSkillLv"` // 主动技能等级
	PassSkillIds  []string          `json:"passSkillIds" bson:"passSkillIds"`   // 被动技能ID集合
}

type PlayerJoinGameResp struct {
	Role       *JoinGameRoleResp       `json:"role"`       // 玩家角色数据
	WeaponMain *JoinGameWeaponMainResp `json:"weaponMain"` // 主武器数据
	WeaponSub  *JoinGameWeaponSubResp  `json:"weaponSub"`  // 武器数据
	HeroList   []*JoinGameHeroResp     `json:"heroList"`   // 英雄集合
	Armor      *JoinGameWeaponMainResp `json:"Armor"`      // 护甲装备数据
	Ornaments  *JoinGameWeaponMainResp `json:"ornaments"`  // 饰品装备数据
}

// 关卡结算需要返回的数据
type SettlementStageResp struct {
	StageId     interface{}  `json:"stageId"`   //关卡ID
	StageType   int  `json:"stageType"` //关卡类型(枚举值)
	StagePass   bool `json:"stagePass"` //0未通过  1通过
	Attachments []*global.Item
	Replace     map[string]*ReplaceColl
}

// ============================================================================

type GetUserStageRecordResp struct {
	StageMain model.StageMain `json:"stageMain" bson:"stageMain"`
}
