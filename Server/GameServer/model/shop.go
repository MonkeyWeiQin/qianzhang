package model

type ShopModel struct {
	Goods  []*ShopGoods
	Chests []*ShopChest
}

// 商品
type ShopGoods struct {
	Id          string `json:"id"           bson:"id"`          //商品ID
	Ty          string `json:"type"         bson:"des"`         //商品
	Price       int    `json:"price"        bson:"price"`       //价格
	ActivePrice int    `json:"activePrice"  bson:"activePrice"` //活动价格
	Count       int    `json:"count"        bson:"count"`       //数量
	Icon        string `json:"icon"         bson:"icon"`        //资源icon
	Discount    string `json:"discount"     bson:"discount"`    //活动赠送数量
	ActiveStart int    `json:"activeStart"  bson:"activeStart"` //活动开始时间
	ActiveEnd   int    `json:"activeEnd"  bson:"activeEnd"`     //活动结束时间
}

// 宝箱
type ShopChest struct {
	ID          string `json:"id"`          //宝箱id
	Ty          int    `json:"type"`        //宝箱描述
	Des         string `json:"des"`         //类型
	OpenOne     int    `json:"openOne"`     //开宝箱一次消耗
	AcOpenOne   int    `json:"acOpenOne"`   //活动开宝箱一次消耗
	OpenTen     int    `json:"openTen"`     //开宝箱十次消耗
	AcOpenTen   int    `json:"acOpenTen"`   //活动开宝箱十次消耗
	TimeLimit   int    `json:"timeLimit"`   //免费开箱倒计时（分）
	LimitCount  int    `json:"limitCount"`  //连续开箱必得上限
	RewardId    string `json:"rewardId"`    //连抽必中奖励id
	ActiveStart int    `json:"activeStart"` //活动开始时间
	ActiveEnd   int    `json:"activeEnd"`   //活动结束时间
}
