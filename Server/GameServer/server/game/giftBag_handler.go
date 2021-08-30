package game

import (
	"battle_rabbit/codec"
	"battle_rabbit/define"
	"battle_rabbit/iface"
	"battle_rabbit/model"
	"battle_rabbit/protocol"
	"battle_rabbit/protocol/response"
	"battle_rabbit/service/log"
	jsoniter "github.com/json-iterator/go"
	"time"
)

type UseGiftBagRequest struct {
	Code string `json:"code"`
}

func (g *Game) UseGiftBag(sess iface.ISession, msg *codec.Message) {
	var code = define.MsgCode200
	defer func() {
		if code != define.MsgCode200 {
			sess.Send(protocol.ErrCode(msg.Id, code))
		}
	}()

	req := new(UseGiftBagRequest)
	err := jsoniter.Unmarshal(msg.GetData(), &req)
	if err != nil {
		log.Error("解析请求出错 :", err)
		code = define.MsgCode400
		return
	}
	userId := sess.GetUid()

	gift, GetGiftErr := model.GetGiftBagLogCollection().GetGift(req.Code)
	if GetGiftErr != nil || gift == nil {
		log.Error("获取兑换码失败 :", GetGiftErr)
		code = define.MsgCode1204
		return
	}
	if gift.Uid != 0 {
		log.Error("兑换码已被使用 :")
		code = define.MsgCode1201
		return
	}

	used, err := model.GetGiftBagLogCollection().CheckUsedByMid(gift.Mid, userId)
	if err != nil {
		log.Error("获取兑换信息失败 :", err)
		code = define.MsgCode1204
		return
	}

	if used {
		code = define.MsgCode1203
		return
	}

	getGift, err := model.GetGiftBagCollection().GetGift(gift.Mid)
	if err != nil || getGift == nil {
		log.Error("获取兑换信息失败 :", GetGiftErr)
		code = define.MsgCode1204
		return
	}

	if getGift.EffectiveTime < time.Now().Unix() {
		log.Error("兑换码已过期")
		code = define.MsgCode1202
		return
	}
	player := g.GetPlayer(sess.GetUid())
	if player == nil {
		code = define.MsgCode401
		return
	}
	var resp *response.AttachmentsResp
	if resp, code = saveItem(g, player, sess, getGift.Attachment); code != define.MsgCode200 {
		code = define.MsgCode1204
		return
	} else {
		_, SetCodeUsedErr := model.GetGiftBagLogCollection().SetCodeUsed(req.Code, userId)
		if SetCodeUsedErr != nil {
			log.Error("设置兑换码已使用失败 :", SetCodeUsedErr)
			code = define.MsgCode1204
			return
		}
		sess.Send(protocol.SuccessData(msg.Id,resp))
	}
}
