package players

import "com.xv.admin.server/model/db"

// GetGiftListRequest 请求参数
type GetGiftListRequest struct {
	Code  string `form:"code"`
	Mid   int    `form:"mid"`
	Page  int64  `form:"page"`
	Limit int64  `form:"limit"`
}

type CreateGiftRequest struct {
	Name          string           `json:"name"         bson:"name"`           //名称
	Count         int              `json:"count" bson:"count"`                 //兑换次数
	Info          string           `json:"info"         bson:"info"`           //礼包说明
	Code          string           `json:"code" bson:"code"`                   //兑换码
	Attachment    []*db.Attachment `json:"attachment" bson:"attachment"`       //奖励表
	EffectiveTime int64            `json:"effectiveTime" bson:"effectiveTime"` //有效时间 0表示永久有效
}

// GiftDownloadRequest 请求参数
type GiftDownloadRequest struct {
	Mid int `form:"mid"`
}

// GetGiftLogListRequest 请求参数
type GetGiftLogListRequest struct {
	Uid   int    `form:"uid"`
	Code  string `form:"code"`
	Mid   int    `form:"mid"`
	Used  int8   `form:"used"`
	Page  int64  `form:"page"`
	Limit int64  `form:"limit"`
}
