package lobby

// 礼包数据库表名 "giftcode"
// 使用工具，批量生成，形成表格，然后导入数据库，给运营同步
type GiftBagCSVTableInfo struct {
	ID       string        `csv:"id"`           //编号
	Name     string        `csv:"name"`         //名称
	Icon     string        `csv:"icon"`         //图标
	Info     string        `csv:"info"`         //礼包说明
	Gold     int64         `csv:"gold"`         //金币
	Diamond  int64         `csv:"diamond"`      //钻石
	Item_1   string        `csv:"item_1"`       //物品1
	Number_1 int           `csv:"number_1"`     //物品1数量
	Item_2   string        `csv:"item_2"`       //物品2
	Number_2 int           `csv:"number_2"`     //物品2数量
	Item_3   string        `csv:"item_3"`       //物品3
	Number_3 int           `csv:"number_3"`     //物品3数量
}

type GiftBagTableData struct {
	ID       string        `json:"id"           bson:"id"`           //编号
	Name     string        `json:"name"         bson:"name"`         //名称
	Icon     string        `json:"icon"         bson:"icon"`         //图标
	Info     string        `json:"info"         bson:"info"`         //礼包说明
	Gold     int64         `json:"gold"         bson:"gold"`         //金币
	Diamond  int64         `json:"diamond"      bson:"diamond"`      //钻石
	Item_1   string        `json:"item_1"       bson:"item_1"`       //物品1
	Number_1 int           `json:"number_1"     bson:"number_1"`     //物品1数量
	Item_2   string        `json:"item_2"       bson:"item_2"`       //物品2
	Number_2 int           `json:"number_2"     bson:"number_2"`     //物品2数量
	Item_3   string        `json:"item_3"       bson:"item_3"`       //物品3
	Number_3 int           `json:"number_3"     bson:"number_3"`     //物品3数量
}

// 礼包码数据库结构
type GiftCode struct {
	Id 			int	            `bson:"_id" json:"id"`
	GiftID      string          `bson:"giftid" json:"giftid"`
	Status      int             `bson:"status" json:"status"` // 0 未使用，1 已经使用
}

// 客户端请求消息
type GiftCodeRequest struct {
	Giftid string  `json:"giftid"` // 礼包码
}
type GiftListInfo struct {
	ItemID string  `json:"itemid"`
	Count  int     `json:"count"`
}
// 回执客户端消息
type GiftCodeResult struct {
	GiftID     string              `json:"giftid"`  // 礼包码
	Status     int                 `json:"status"`  // 1 已经使用 2 正常使用 3 不存在
	Gold       float64             `json:"gold"`
	Money      float64             `json:"money"`
	ItemList   []GiftListInfo      `json:"itemlist"`
}
// 检查是否存在礼包码，不存在，或者已经使用，如果存在直接发放对应奖励，并通知客户端
func ReciveGiftCode(userId string, giftid string){

}

// 发放物品，并通知客户端