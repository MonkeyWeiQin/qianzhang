package model

import (
	"battle_rabbit/define"
	"battle_rabbit/service/mgoDB"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type GiftCodeLogModel struct {
	Id       primitive.ObjectID `json:"-" bson:"-"`
	Mid      int                `json:"mid" bson:"mid"`            // 兑换码ID
	Code     string             `json:"code" bson:"code"`          // 唯一的ID
	Uid      int                `json:"uid" bson:"uid"`            // 使用者ID
	UsedTime int                `json:"usedTime" bson:"used_time"` // 使用时间
}

var (
	giftBagLogColl = newGiftBagLogColl()
)

type GiftBagLogCollection struct {
	*mgoDB.DbBase
}

func newGiftBagLogColl() *GiftBagLogCollection {
	return &GiftBagLogCollection{mgoDB.NewDbBase(define.MgoDBNameBattle,define.TableNameBroadcast)}
}

func GetGiftBagLogCollection() *GiftBagLogCollection {
	return giftBagLogColl
}

func (m *GiftBagLogCollection) GetGift(code string) (data *GiftCodeLogModel, err error) {
	err = m.FindOne(nil, bson.D{{"code", code}}, &data)
	return
}

// todo 更改修改
func (m *GiftBagLogCollection) CheckUsedByMid(mid int, Uid int) (bool, error) {
	data := &GiftCodeLogModel{}
	err := m.FindOne(nil, bson.D{{"mid", mid}, {"uid", Uid}}, data)
	if err != nil  {
		if err == mongo.ErrNilDocument {
			return false,nil
		}
		return false, err
	}
	return  true,nil
}
func (m *GiftBagLogCollection) SetCodeUsed(code string, uid int) (bool, error) {
	res, err := m.UpdateOne(nil, bson.D{{"code", code}}, bson.D{{"uid", uid}, {"usedTime", int(time.Now().Unix())}})
	if err != nil {
		return false, err
	}
	if res < 0 {
		return false, nil
	}
	return true, nil
}
