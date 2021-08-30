package global

import (
	"battle_rabbit/define"
	"battle_rabbit/utils"
)

type Item struct {
	ItemID   string          `json:"itemId" bson:"itemId"`     // 配置表ID 或者 关联ID,"gold","diamond","exp"
	Count    int             `json:"count" bson:"count"`       // 数量
	ItemType define.ItemType `json:"itemType" bson:"itemType"` // 0 经验 1 金币  2钻石 3体力 4卡片 5英雄 6主武器 7副武器 8装备  9饰品 10 英雄装备
}

// 无条件拆解,所以传进来得参数必须是符合拆解格式得,并且不能为 空字符串
func DisassembleToItem(str string, ty define.ItemType) (items []*Item) {
	m := utils.UnmarshalItemsInt(str)
	for s, i := range m {
		items = append(items, &Item{
			ItemID:   s,
			Count:    i,
			ItemType: ty,
		})
	}
	return items
}
