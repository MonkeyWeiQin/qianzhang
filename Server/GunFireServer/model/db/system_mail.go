package db

import (
	"com.xv.admin.server/config"
	"com.xv.admin.server/service/mgoDB"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	Attachment  []*Attachment      `json:"attachment" bson:"attachment"`   //附件
}

type Attachment struct {
	ItemID   string `json:"itemId" bson:"itemId"`
	Count    int    `json:"count" bson:"count"`       //数量
	ItemType int    `json:"itemType" bson:"itemType"` //1 金币  2钻石  3
}

const (
	NewPlayerWelcomeMailType = 1
)

func (mail *SystemMailModel) CollectionName() string {
	return config.CollSystemMail
}

func (mail *SystemMailModel) GetId() primitive.ObjectID {
	return mail.Id
}

func (mail *SystemMailModel) SetId(id primitive.ObjectID) {
	mail.Id = id
}

// Insert 写入一条数据
func (mail *SystemMailModel) Insert() error {
	mail.Mid = primitive.NewObjectID().Hex()
	mail.Id = primitive.NewObjectID()
	err := mgoDB.GetMgo().InsertOne(nil, mail)
	return err
}
func (mail *SystemMailModel) SelectOneMenu(filter bson.D, opt *options.FindOneOptions) (data *SystemMailModel, err error) {
	oneFilter := mgoDB.NewOneFinder(mail).Where(filter).Options(opt).Record(&data)
	res, err := mgoDB.GetMgo().FindOne(nil, oneFilter)
	if err != nil {
		return nil, err
	}
	if !res {
		return nil, nil
	}

	return data, err
}
// GetPlayerMailList 获取邮箱列表
func (mail *SystemMailModel) GetPlayerMailList(filter bson.D, options *options.FindOptions) (data []*SystemMailModel, err error) {
	finder := mgoDB.NewFinder(mail).Where(filter).Options(options).Records(&data)
	err = mgoDB.GetMgo().FindMany(context.TODO(), finder)
	if err != nil {
		return nil, err
	}
	return
}
func (mail *SystemMailModel) GetTotalNum(filter bson.D) (int64, error) {
	finder := mgoDB.NewCounter(mail).Where(filter)
	return mgoDB.GetMgo().CountDocuments(nil, finder)
}


// UpdatePlayerMail 更新邮件
func (mail *SystemMailModel) UpdatePlayerMail(filter bson.D, update bson.D, opt *options.UpdateOptions) error {
	updater := mgoDB.NewUpdater(mail).Where(filter).Update(update).Options(opt)
	_, err := mgoDB.GetMgo().UpdateOne(nil, updater)
	return err
}