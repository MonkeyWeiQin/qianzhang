package game

import (
	"battle_rabbit/codec"
	"battle_rabbit/define"
	"battle_rabbit/global"
	"battle_rabbit/iface"
	"battle_rabbit/model"
	"battle_rabbit/protocol"
	"battle_rabbit/protocol/request"
	"battle_rabbit/protocol/response"
	"battle_rabbit/service/log"
	jsoniter "github.com/json-iterator/go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math"
)

// CreatePlayerMail 增加邮件
func (g *Game) CreatePlayerMail(sess iface.ISession, msg *codec.Message) {
	model.GetMailCollection().SendNewUserMail(0000)
	return
}

// GetPlayerMailList 获得玩家的邮件列表
func (g *Game) GetPlayerMailList(sess iface.ISession, msg *codec.Message) {
	req := new(request.GetPlayerMailListRequest)
	err := jsoniter.Unmarshal(msg.Data, &req)
	if err != nil {
		log.Error("解析请求出错 :", err)
		sess.Send(protocol.ErrCode(msg.Id, define.MsgCode400))
		return
	}

	userId := sess.GetUid()
	if req.Page < 1 {
		req.Page = 1
	}
	var (
		limit = 10
		skip  = (req.Page - 1) * limit
		resp  []*response.GetPlayerMailListResponse
	)

	num, err := model.GetMailCollection().GetTotalNum(userId, nil)
	if err != nil {
		log.Error("获取邮件失败 :", err)
		sess.Send(protocol.ErrCode(msg.Id, define.MsgCode400))
		return
	}

	findOptions := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)).SetSort(bson.M{"state": 1})
	list, err := model.GetMailCollection().GetPlayerMailList(nil, userId, findOptions)
	if err != nil {
		log.Error("获取邮件失败 :", err)
		sess.Send(protocol.ErrCode(msg.Id, define.MsgCode400))
		return
	}

	for _, item := range list {
		resp = append(resp, &response.GetPlayerMailListResponse{
			MailType:     2,
			MailId:       item.Mid,
			Title:        item.Title,
			Content:      item.Content,
			Attachment:   item.Attachment,
			State:        item.State,
			ReceiveState: item.ReceiveState,
			SendTime:     item.SendTime,
		})
	}
	r := protocol.SuccessData(msg.Id, map[string]interface{}{
		"list":      resp,
		"page":      req.Page,
		"pageTotal": math.Ceil(float64(num) / float64(limit)),
	})
	sess.Send(r)
}

// SetPlayerMailRead 设置邮件已读
func (g *Game) SetPlayerMailRead(sess iface.ISession, msg *codec.Message) {
	ReceiveMailGiftRequest := new(request.ReceiveMailGiftRequest)
	err := jsoniter.Unmarshal(msg.GetData(), &ReceiveMailGiftRequest)
	if err != nil {
		log.Error("解析请求出错 :", err)
		sess.Send(protocol.ErrCode(msg.Id, define.MsgCode400))
		return
	}
	err = model.GetMailCollection().SetMailState([]string{ReceiveMailGiftRequest.MailId}, sess.GetUid(), model.STATE_HAVE_READ)
	if err != nil {
		log.Error(err)
		sess.Send(protocol.Err(msg.Id))
		return
	}
	sess.Send(protocol.Success(msg.Id))
}

// DelPlayerMail 删除邮件
func (g *Game) DelPlayerMail(sess iface.ISession, msg *codec.Message) {
	ReceiveMailGiftRequest := new(request.ReceiveMailGiftRequest)
	err := jsoniter.Unmarshal(msg.GetData(), &ReceiveMailGiftRequest)
	if err != nil {
		log.Error("解析请求出错 :", err)
		sess.Send(protocol.ErrCode(msg.Id, define.MsgCode400))
		return
	}

	err = model.GetMailCollection().SetMailState([]string{ReceiveMailGiftRequest.MailId}, sess.GetUid(), model.STATE_DELETE)
	if err != nil {
		log.Error(err)
		sess.Send(protocol.ErrCode(msg.Id, define.MsgCode500))
		return
	}
	sess.Send(protocol.Success(msg.Id))
}

// ReceiveMailGifts 一键领取邮件附件
func (g *Game) ReceiveMailGifts(sess iface.ISession, msg *codec.Message) {
	userId := sess.GetUid()
	listAttachment, err := model.GetMailCollection().GetPlayerMailValidityAttachment(nil, userId)
	if err != nil {
		log.Error("获取邮件附件失败 :", err)
		sess.Send(protocol.ErrCode(msg.Id, define.MsgCode1101))
		return
	}

	var attachments []*global.Item
	var mailIds []string

	for _, attachment := range listAttachment {
		mailIds = append(mailIds, attachment.Mid)
		attachments = append(attachments, attachment.Attachment...)
	}

	if len(attachments) == 0 {
		log.Error("获取邮件附件失败 : 没有可领取的附件!!")
		sess.Send(protocol.ErrCode(msg.Id, define.MsgCode1102))
		return
	}

	player := g.GetPlayer(sess.GetUid())
	if player == nil { // 重新登录
		sess.Send(protocol.ErrCode(msg.Id, define.MsgCode401))
		return
	}

	resp, code := saveItem(g, player, sess, attachments)
	if code != define.MsgCode200 {
		sess.Send(protocol.ErrCode(msg.Id, code))
		return
	}

	err = model.GetMailCollection().SetMailReceive(mailIds, userId)
	if err != nil {
		log.Error("领取失败 : 设置邮件状态失败")
		sess.Send(protocol.ErrCode(msg.Id, define.MsgCode1101))
		return
	}
	sess.Send(protocol.SuccessData(msg.Id, resp.Attachments))
}

// ReceiveMailGift  单个领取邮件附件
func (g *Game) ReceiveMailGift(sess iface.ISession, msg *codec.Message) {
	req := new(request.ReceiveMailGiftRequest)
	err := jsoniter.Unmarshal(msg.GetData(), &req)
	if err != nil {
		log.Error("解析请求出错 :", err)
		sess.Send(protocol.ErrCode(msg.Id, define.MsgCode400))
		return
	}

	player := g.GetPlayer(sess.GetUid())
	if player == nil { // 重新登录
		sess.Send(protocol.ErrCode(msg.Id, define.MsgCode401))
		return
	}

	var Attachment model.AttachmentItem
	mail, err := model.GetMailCollection().GetMailById(nil, player.Uid, req.MailId)
	if err != nil {
		log.Error("获取邮件附件失败 :", err)
		sess.Send(protocol.ErrCode(msg.Id, define.MsgCode1101))
		return
	}
	if mail.ReceiveState != false {
		sess.Send(protocol.ErrCode(msg.Id, define.MsgCode1103))
		return
	}
	Attachment.Attachment = mail.Attachment
	Attachment.Mid = mail.Mid

	if len(Attachment.Attachment) <= 0 {
		sess.Send(protocol.ErrCode(msg.Id, define.MsgCode1102))
		return
	}

	resp, code := saveItem(g, player,sess, mail.Attachment)
	if code != define.MsgCode200 {
		sess.Send(protocol.ErrCode(msg.Id, define.MsgCode1101))
		return
	}

	if err := model.GetMailCollection().SetMailReceive([]string{Attachment.Mid}, player.Uid); err != nil {
		log.Error("领取失败:设置邮件状态失败", err)
		sess.Send(protocol.ErrCode(msg.Id, define.MsgCode1101))
		return
	}
	sess.Send(protocol.SuccessData(msg.Id,resp.Attachments))
}

// DelPlayerMails 一键删除
func (g *Game) DelPlayerMails(sess iface.ISession, msg *codec.Message) {
	userId := sess.GetUid()
	_, err := model.GetMailCollection().SetAllMailDel(userId)
	if err != nil {
		log.Error(err.Error())
		sess.Send(protocol.Err(msg.Id))
		return
	}
	sess.Send(protocol.Success(msg.Id))
}
