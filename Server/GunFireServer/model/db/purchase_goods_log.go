package db

import (
	"com.xv.admin.server/config"
	"com.xv.admin.server/service/mgoDB"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)
type PurchaseGoodsLogModel struct {
	Id           primitive.ObjectID `json:"-"  bson:"-"`
	Mid          string             `json:"mid" bson:"mid"`
	Uid          int                `json:"uid" bson:"uid"`
	GoodsId      string             `json:"goodsId" bson:"goodsId"`     // 商品ID
	GoodsType    int                `json:"goodsType" bson:"goodsType"` // 商品类型
	Time         int                `json:"time" json:"time"`           //中奖时间
	Price        int                `json:"price" bson:"price"`
	Count        int                `json:"count" bson:"count"`
	Presentation int                `json:"presentation" bson:"presentation"`
}

func (m *PurchaseGoodsLogModel) CollectionName() string {
	return config.CollPurchaseGoodsLog
}

func (m *PurchaseGoodsLogModel) GetId() primitive.ObjectID {
	return m.Id
}

func (m *PurchaseGoodsLogModel) SetId(id primitive.ObjectID) {
	m.Id = id
}

func (m *PurchaseGoodsLogModel) GetList(filter bson.D, options *options.FindOptions) (data []*PurchaseGoodsLogModel, err error) {
	finder := mgoDB.NewFinder(m).Options(options).Records(&data)
	if filter != nil {
		finder.Where(filter)
	}
	err = mgoDB.GetMgo().FindMany(context.TODO(), finder)
	return
}

func (m *PurchaseGoodsLogModel) GetTotalNum(filter bson.D) (int64, error) {
	finder := mgoDB.NewCounter(m).Where(filter)
	return mgoDB.GetMgo().CountDocuments(nil, finder)
}
