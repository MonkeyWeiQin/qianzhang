package request

type GetPlayerMailListRequest struct {
	Page     int `form:"page" `
}

type SetPlayerMailStateRequest struct {
	Id string `form:"id"  binding:"required"`
}

type ReceiveMailGiftRequest struct {
	MailType int    `json:"type"` //邮件类型 1：系统邮件  2其他邮件
	MailId   string `json:"id"`   //邮件ID
}
