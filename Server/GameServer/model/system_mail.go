package model

import (
	"battle_rabbit/define"
	"battle_rabbit/global"
	"battle_rabbit/service/mgoDB"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// SystemMailModel  系统邮件功能
type SystemMailModel struct {
	Id          primitive.ObjectID `json:"-"           bson:"-"`
	Mid         string             `json:"mid" bson:"mid"`
	SendTime    int                `json:"sendTime" bson:"sendTime"`       //发送时间
	InvalidTime int                `json:"invalidTime" bson:"invalidTime"` //失效时间
	CronTime    int                `json:"cronTime" bson:"cronTime"`       //定时发送
	Uid         []int              `json:"uid" bson:"uid"`                 //收件人IDs
	ISAllUser   bool               `json:"isAllUser" bson:"isAllUser"`     //是否是全服发送
	Title       string             `json:"title" bson:"title"`             //标题
	Content     string             `json:"content" bson:"content"`         //内容
	Attachment  []*global.Item      `json:"attachment" bson:"attachment"`   //附件
}


var (
	SystemMailColl = newSystemMailCollection()
)

type SystemMailCollection struct {
	*mgoDB.DbBase
}

func newSystemMailCollection() *SystemMailCollection {
	return &SystemMailCollection{mgoDB.NewDbBase(define.MgoDBNameBattle,define.TableNameSystemMail)}
}

func GetSystemMailCollection() *SystemMailCollection {
	return SystemMailColl
}

// GetUserMail 获取用户邮件列表
// 1 失效时间大于当前时间
// 2 用户ID 为当前用户
// 3
func (e *SystemMailCollection) GetUserMail(ctx context.Context, uid int, options ...*options.FindOptions) (data []*SystemMailModel, err error) {
	err = e.FindAll(ctx, bson.D{
		{"$and", []bson.M{
			{"invalidTime": bson.M{"$gte": time.Now().Unix()}},
			{"$or": []bson.M{ {"uid": bson.M{"$in": []int{uid}}}, {"isAllUser": true}}}},
		}}, &data, options...)
	return
}

// GetUserMailIds 获取用户邮件 ids列表
func (e *SystemMailCollection) GetUserMailIds(ctx context.Context, uid int) (data []string, err error) {
	var res []map[string]string
	err = e.FindAll(ctx, bson.D{{"$and", []bson.M{{"invalidTime": bson.M{"$gte": time.Now().Unix()}}, {"$or": []bson.M{{"uid": bson.M{"$in": []int{uid}}}, {"isAllUser": true}}}}}}, &res, options.Find().SetProjection(bson.M{"_id": 0, "mid": 1}))
	for _, item := range res {
		data = append(data, item["mid"])
	}
	return
}

// GetMailAttachment 获取邮件附件列表
func (e *SystemMailCollection) GetMailAttachment(ctx context.Context, uid int) (data []*global.Item, err error) {
	err = e.FindAll(ctx, bson.D{{"$and", []bson.M{{"invalidTime": bson.M{"$gte": time.Now().Unix()}}, {"$or": []bson.M{{"uid": bson.M{"$in": []int{uid}}}, {"isAllUser": true}}}}}}, &data, options.Find().SetProjection(bson.M{"_id": 0, "mid": 1, "attachment": 1}))
	return
}

func (e *SystemMailCollection) GetMailById(ctx context.Context, uid int, mid string) (data *global.Item, err error) {
	err = e.FindOne(ctx, bson.D{{"$and", []bson.M{{"invalidTime": bson.M{"$gte": time.Now().Unix()}}, {"$or": []bson.M{{"uid": bson.M{"$in": []int{uid}}}, {"isAllUser": true}}}, {"mid": mid}}}}, &data, options.FindOne().SetProjection(bson.M{"_id": 0, "mid": 1, "attachment": 1}))
	return
}

// GetMailAttachmentNotInIds 获取邮件附件列表
func (e *SystemMailCollection) GetMailAttachmentNotInIds(ctx context.Context, uid int, ids []string) (data []*global.Item, err error) {
	err = e.FindAll(ctx, bson.D{{"$and", []bson.M{{"invalidTime": bson.M{"$gte": time.Now().Unix()}}, {"$or": []bson.M{{"uid": bson.M{"$in": []int{uid}}}, {"isAllUser": true}}}, {"mid": bson.M{"$nin": ids}}}}}, &data, options.Find().SetProjection(bson.M{"_id": 0, "mid": 1, "attachment": 1}))
	return
}
