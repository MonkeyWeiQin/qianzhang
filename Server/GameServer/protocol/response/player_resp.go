package response

import "battle_rabbit/model"

// http 校验用户Token HTTP响应数据
type ValidateUserResp struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"` // uid
	Msg  string      `json:"msg"`
}

type WebUserLoginByDevIdResp struct {
	Token   string `json:"token"`
	TcpHost string `json:"tcpHost"`
}

type RoleInfoResp struct {
	Role       *model.RoleModel      `json:"role"`       // 玩家角色数据
	WeaponMain *model.EquipmentModel `json:"weaponMain"` // 主武器数据
	WeaponSub  *model.EquipmentModel `json:"weaponSub"`  // 武器数据
	Armor      *model.EquipmentModel `json:"armor"`      // 护甲装备数据
	Ornaments  *model.EquipmentModel `json:"ornaments"`  // 饰品装备数据
}


type GetRoleTalentLvResp struct {
	Lv      int8 `json:"lv"`      // 升级后的等级
	Ty      int8 `json:"ty"`      // 升级的天赋类型
	TotalLv int  `json:"totalLv"` // 升级后的总等级
}

type FlushStrengthResp struct {
	Strength int `json:"strength"` //体力值
	Time     int `json:"time"`     // 最后一次刷新的时间戳
}

type UpgradeRoleStarResp struct {
	RLv        int    `json:"rlv"`        // 角色等级
	RExp       int    `json:"rExp"`       // 角色经验值
	Quality    int8   `json:"quality"`    // 当前星级
	RelationId string `json:"relationId"` // 数据表关联ID
}

// 角色/账号 升级.加经验推送
type UpgradeLevelResp struct {
	Lv     int  `json:"lv"`
	Exp    int  `json:"exp"`
	IsRole bool `json:"isRole"`
}
