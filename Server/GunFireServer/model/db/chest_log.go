package db

import (
	"com.xv.admin.server/config"
	"com.xv.admin.server/service/mgoDB"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ChestLogModel struct {
	Id        primitive.ObjectID `json:"-"  bson:"-"`
	Mid       string             `json:"mid" bson:"mid"`
	Uid       int                `json:"uid" bson:"uid"`
	ChestId   string             `json:"chestId" bson:"chestId"`     // 宝箱ID
	Price     int                `json:"price" bson:"price"`         //价格
	Count     int                `json:"count" bson:"count"`         //次数
	ChestItem []*Attachment      `json:"chestItem" bson:"chestItem"` //中奖内容
	Time      int                `json:"time" json:"time"`           //中奖时间
}

func (m *ChestLogModel) CollectionName() string {
	return config.CollChestLog
}

func (m *ChestLogModel) GetId() primitive.ObjectID {
	return m.Id
}

func (m *ChestLogModel) SetId(id primitive.ObjectID) {
	m.Id = id
}
func (m *ChestLogModel) GetList(filter bson.D, options *options.FindOptions) (data []*ChestLogModel, err error) {
	finder := mgoDB.NewFinder(m).Options(options).Records(&data)
	if filter != nil {
		finder.Where(filter)
	}
	err = mgoDB.GetMgo().FindMany(context.TODO(), finder)
	return
}

func (m *ChestLogModel) GetTotalNum(filter bson.D) (int64, error) {
	finder := mgoDB.NewCounter(m).Where(filter)
	return mgoDB.GetMgo().CountDocuments(nil, finder)
}
