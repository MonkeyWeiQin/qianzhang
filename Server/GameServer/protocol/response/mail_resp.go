package response

import (
	"battle_rabbit/define"
	"battle_rabbit/global"
)

type GetPlayerMailListResponse struct {
	MailType     int              `json:"type"` //邮件类型 1：系统邮件  2其他邮件
	MailId       string           `json:"id"`   //邮件ID
	Title        string           `json:"title"`
	Content      string           `json:"content"`      //内容
	Attachment   []*global.Item `json:"attachment"`   //附件
	State        int              `json:"state"`        // 邮件状态：未读 0 已读 1 删除 2
	ReceiveState bool             `json:"receiveState"` // 领取状态：false未领取  true 已领取
	SendTime     int              `json:"sendTime"`     //发送时间
}

type ReplaceColl struct {
	NewLv int  	// 最新的等级
	OldLv int   // 被替换的装备的等级
	RelationId string // 装备关联ID
	ItemType define.ItemType // 装备类型
}

type AttachmentsResp struct {
	Attachments []*global.Item
	Replace map[string]*ReplaceColl
}