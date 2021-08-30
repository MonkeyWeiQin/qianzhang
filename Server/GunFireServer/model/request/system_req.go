package request

import "com.xv.admin.server/model/db"

// CloseGameForIDRequest 管理员更改游戏开关
type CloseGameForIDRequest struct {
	GameID int `form:"gameid" binding:"required"`
}

// PaginationRequest Pagination 分页参数
type PaginationRequest struct {
	Page     int64 `form:"page"`
	PageSize int64 `form:"page_size"`
}


type UpdateBlindBoxItemRequest struct {
	Type int                `json:"type" bson:"type"`
	Item []*db.BlindBoxItem `json:"item" bson:"item"`
}
type UpdateBlindBoxPriceRequest struct {
	Type            int `json:"type" bson:"type"`
	UnitPrice       int `json:"unitPrice" bson:"unitPrice"`             //单次抽奖价格
	RepeatedlyPrice int `json:"repeatedlyPrice" bson:"repeatedlyPrice"` //10次抽奖价格
}