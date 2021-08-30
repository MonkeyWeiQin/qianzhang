package request

import "battle_rabbit/define"

// ChangeArmorRequest 切换装备/取消装备
type ChangeEquipmentRequest struct {
	EquipmentType define.EquipageType `json:"equipmentType"`
	EquipmentId   string              `json:"equipmentId"`
}

// GetArmorListRequest 获取装备列表
type GetEquipmentListRequest struct {
	Page          int                 `json:"page"`
	EquipmentType define.EquipageType `json:"equipmentType"`
}

type EquipmentUpgradeReq struct {
	Ty  define.EquipageType `json:"ty"` // 0 升级 1 升星
	Mid string              `json:"mid"`
}
