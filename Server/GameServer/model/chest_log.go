package model

import (
	"battle_rabbit/define"
	"battle_rabbit/global"
	"battle_rabbit/service/mgoDB"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type ChestLog struct {
	Id        primitive.ObjectID `json:"-"  bson:"-"`
	Mid       string             `json:"mid" bson:"mid"`
	Uid       int                `json:"uid" bson:"uid"`
	ChestId   string             `json:"chestId" bson:"chestId"`     // 宝箱ID
	Price     int                `json:"price" bson:"price"`         //价格
	Count     int                `json:"count" bson:"count"`         //次数
	ChestItem []*global.Item       `json:"chestItem" bson:"chestItem"` //中奖内容
	Time      int                `json:"time" json:"time"`           //中奖时间
}

var (
	chestLogColl = newChestLogColl()
)

type ChestLogCollection struct {
	*mgoDB.DbBase
}

func newChestLogColl() *ChestLogCollection {
	return &ChestLogCollection{mgoDB.NewDbBase(define.MgoDBNameBattle,define.TableNameBroadcast)}
}
func GetChestLogCollection() *ChestLogCollection {
	return chestLogColl
}

func (e *ChestLogCollection) Create(uid int, ChestId string, Price int, Count int, Attachment []*global.Item) error {
	return e.InsertOne(nil, bson.D{
		{"mid", primitive.NewObjectID().Hex()},
		{"uid", uid},
		{"chestId", ChestId},
		{"price", Price},
		{"count", Count},
		{"chestItem", Attachment},
		{"time", time.Now().Unix()},
	})
}
