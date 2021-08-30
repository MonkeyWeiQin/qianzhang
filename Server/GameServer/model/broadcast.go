package model

import (
	"battle_rabbit/define"
	"battle_rabbit/service/mgoDB"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// BroadcastModel 公告信息数据
type BroadcastModel struct {
	Id          primitive.ObjectID `json:"-"  bson:"-"`
	Mid         int                `json:"mid" bson:"mid"`                 //广播ID
	Content     string             `json:"content" bson:"content"`         //公告的字符内容
	Color       string             `json:"color" bson:"color"`             //展示颜色 1红色  2蓝色  3黄色
	StartTime   int                `json:"startTime" bson:"startTime"`     //开始时间
	EndTime     int                `json:"endTime" bson:"endTime"`         //结束时间
	SpacingTime int64              `json:"spacingTime" bson:"spacingTime"` //展示时间  单位：秒
}

type GetBroadcastListResult struct {
	Mid         int    `json:"mid"  bson:"mid"`                //公告ID
	Content     string `json:"content" bson:"content"`         //内容
	Color       string `json:"color" bson:"color"`             //展示颜色 1红色  2蓝色  3黄色
	SpacingTime int64  `json:"spacingTime" bson:"spacingTime"` //展示时间  单位：秒
}

var (
	broadcastColl = newBroadcastCollection()
)

type BroadcastCollection struct {
	*mgoDB.DbBase
}

func newBroadcastCollection() *BroadcastCollection {
	return &BroadcastCollection{mgoDB.NewDbBase(define.MgoDBNameBattle,define.TableNameBroadcast)}
}

func GetBroadcastCollection() *BroadcastCollection {
	return broadcastColl
}

func (m *BroadcastCollection) GetBroadcast(ctx context.Context) (data []*GetBroadcastListResult, err error) {
	err = m.FindAll(ctx, bson.D{{"$and", []bson.M{{"startTime": bson.M{"$lte": time.Now().Unix()}}, {"endTime": bson.M{"$gte": time.Now().Unix()}}}}}, &data, options.Find().SetProjection(bson.M{"_id": 0, "mid": 1, "content": 1, "color": 1, "spacingTime": 1}))
	return
}
