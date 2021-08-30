package db

import (
	"com.xv.admin.server/config"
	"com.xv.admin.server/service/mgoDB"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type GiftCodeModel struct {
	Id            primitive.ObjectID `json:"-" bson:"-"`
	Mid           int                `json:"mid" bson:"mid"`                     //唯一的ID
	Name          string             `json:"name"         bson:"name"`           //名称
	Count         int                `json:"count" bson:"count"`                 //兑换次数
	Info          string             `json:"info"         bson:"info"`           //礼包说明
	Attachment    []*Attachment      `json:"attachment" bson:"attachment"`       //奖励表
	EffectiveTime int64              `json:"effectiveTime" bson:"effectiveTime"` //有效时间 0表示永久有效
}

func (gift *GiftCodeModel) CollectionName() string {
	return config.CollGiftCode
}

func (gift *GiftCodeModel) GetId() primitive.ObjectID {
	return gift.Id
}

func (gift *GiftCodeModel) SetId(id primitive.ObjectID) {
	gift.Id = id
}

// Insert 写入一条数据
func (gift *GiftCodeModel) Insert() error {
	if gift.Mid == 0 {
		opt := options.FindOne().SetSort(bson.D{{"mid", -1}})
		Tmp, err := new(GiftCodeModel).SelectOne(nil, opt)
		if err != nil {
			return err
		}
		if Tmp == nil {
			gift.Mid = 2
		} else {
			gift.Mid = Tmp.Mid + 1
		}
	}
	gift.Id = primitive.NewObjectID()
	err := mgoDB.GetMgo().InsertOne(nil, gift)
	return err
}
func (gift *GiftCodeModel) SelectOne(filter bson.D, opt *options.FindOneOptions) (data *GiftCodeModel, err error) {
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

func (gift *GiftCodeModel) GetList(filter bson.D, opt *options.FindOptions) (data []map[string]interface{}, err error) {
	data = []map[string]interface{}{}
	finder := mgoDB.NewFinder(gift).Where(filter).Options(opt).Records(&data)
	err = mgoDB.GetMgo().FindMany(context.TODO(), finder)
	return
}
func (gift *GiftCodeModel) GetTotalNum(filter bson.D) (int64, error) {
	finder := mgoDB.NewCounter(gift).Where(filter)
	return mgoDB.GetMgo().CountDocuments(nil, finder)
}
