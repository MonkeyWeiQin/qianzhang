package request

type HeroGoToWarRequest struct {
	HeroId string `json:"heroId"`
}

type ObtainHeroRequest struct {
	RelationId string `json:"relationId" ` //数据表关联ID
}

// 英雄升级/升星/主动/被动/装备提升等级操作
type HeroUpgradeReq struct {
	Ty  int8   `json:"ty"`  // 升级类型
	Rid string `json:"rid"` // 数据表关联ID
}
