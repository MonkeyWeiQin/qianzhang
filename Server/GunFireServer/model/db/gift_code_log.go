package db

import (
	"com.xv.admin.server/config"
	"com.xv.admin.server/service/mgoDB"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type GiftCodeLogModel struct {
	Id       primitive.ObjectID `json:"-" bson:"-"`
	Mid      int                `json:"mid" bson:"mid"`           // 兑换码ID
	Code     string             `json:"code" bson:"code"`         // 唯一的ID
	Uid      int                `json:"uid" bson:"uid"`           // 使用者ID
	UsedTime int                `json:"usedTime" bson:"usedTime"` // 使用时间
}

func (gift *GiftCodeLogModel) CollectionName() string {
	return config.CollGiftCodeLog
}

func (gift *GiftCodeLogModel) GetId() primitive.ObjectID {
	return gift.Id
}

func (gift *GiftCodeLogModel) SetId(id primitive.ObjectID) {
	gift.Id = id
}

// Insert 写入一条数据
func (gift *GiftCodeLogModel) Insert(data []interface{}) error {
	err := mgoDB.GetMgo().InsertMany(nil, data)
	return err
}
func (gift *GiftCodeLogModel) SelectOneMenu(filter bson.D, opt *options.FindOneOptions) (data *GiftCodeModel, err error) {
	oneFilter := mgoDB.NewOneFinder(gift).Where(filter).Options(opt).Record(&data)
	res, err := mgoDB.GetMgo().FindOne(nil, oneFilter)
	if err != nil {
		return nil, err
	}
	if !res {
		return nil, nil
	}
	return data, err
}

func (gift *GiftCodeLogModel) GetList(filter bson.D, opt *options.FindOptions) (data []map[string]interface{}, err error) {
	data = []map[string]interface{}{}
	finder := mgoDB.NewFinder(gift).Where(filter).Options(opt).Records(&data)
	err = mgoDB.GetMgo().FindMany(context.TODO(), finder)
	return
}
func (gift *GiftCodeLogModel) GetTotalNum(filter bson.D) (int64, error) {
	finder := mgoDB.NewCounter(gift).Where(filter)
	return mgoDB.GetMgo().CountDocuments(nil, finder)
}
