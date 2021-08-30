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

var (
	noticeColl = newNoticeCollection()
)

type NoticeCollection struct {
	*mgoDB.DbBase
}

func newNoticeCollection() *NoticeCollection {
	return &NoticeCollection{mgoDB.NewDbBase(define.MgoDBNameBattle, define.TableNameBroadcast)}
}

func GetNoticeCollection() *NoticeCollection {
	return noticeColl
}

type NoticeModel struct {
	Id        primitive.ObjectID `json:"-"         bson:"-"`
	Mid       int                `json:"mid"       bson:"mid"`       // 公告ID
	Title     string             `json:"title"     bson:"title"`     // 公告名称
	Content   string             `json:"content"   bson:"content"`   // 公告的字符内容
	Type      int                `json:"type"      bson:"type"`      // 公告类型  //1登录公告  2/游戏内公告
	StartTime int                `json:"startTime" bson:"startTime"` // 开始时间
	EndTime   int                `json:"endTime"   bson:"endTime"`   // 结束时间
}

type GetNoticeListResult struct {
	Mid     int    `json:"mid"  bson:"mid"`        //公告ID
	Title   string `json:"title"  bson:"title"`    //标题
	Content string `json:"content" bson:"content"` //内容
}

func (m *NoticeCollection) GetNotice(ctx context.Context) (data []*GetNoticeListResult, err error) {
	err = m.FindAll(ctx, bson.D{{"$and", []bson.M{{"startTime": bson.M{"$lte": time.Now().Unix()}}, {"endTime": bson.M{"$gte": time.Now().Unix()}}}}}, &data, options.Find().SetProjection(bson.M{"_id": 0, "mid": 1, "title": 1, "content": 1}))
	return
}
