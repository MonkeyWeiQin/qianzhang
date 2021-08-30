package model

import (
	"battle_rabbit/define"
	"battle_rabbit/global"
	"battle_rabbit/service/mgoDB"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GiftCodeModel struct {
	Id            primitive.ObjectID `json:"-" bson:"-"`
	Mid           int                `json:"mid" bson:"mid"`                     //唯一的ID
	Name          string             `json:"name"         bson:"name"`           //名称
	Count         int                `json:"count" bson:"count"`                 //兑换次数
	Info          string             `json:"info"         bson:"info"`           //礼包说明
	Attachment    []*global.Item      `json:"attachment" bson:"attachment"`       //奖励表
	EffectiveTime int64              `json:"effectiveTime" bson:"effectiveTime"` //有效时间 0表示永久有效
}

var (
	giftBagColl = newGiftBagColl()
)

type GiftBagCollection struct {
	*mgoDB.DbBase
}

func newGiftBagColl() *GiftBagCollection {
	return &GiftBagCollection{mgoDB.NewDbBase(define.MgoDBNameBattle,define.TableNameBroadcast)}
}
func GetGiftBagCollection() *GiftBagCollection {
	return giftBagColl
}

func (m *GiftBagCollection) GetGift(mid int) (data *GiftCodeModel, err error) {
	err = m.FindOne(nil, bson.D{{"mid", mid}}, &data)
	return

}
