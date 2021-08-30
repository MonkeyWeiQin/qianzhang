package db

import (
	"com.xv.admin.server/config"
	"com.xv.admin.server/service/mgoDB"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MailItemInfo struct {
	ItemID string `json:"itemid" bson:"itemid"`
	Count  int    `json:"count" bson:"count"`
}

type PlayerMailModel struct {
	Id           primitive.ObjectID `json:"_id"        bson:"_id"`
	Mid          string             `json:"mid"  bson:"mid"`
	UserId       int                `json:"uid"  bson:"uid"`
	Name         string             `json:"name"     bson:"name"`                // 邮件名称
	Sender       string             `json:"sender"   bson:"sender"`              // 发送者名称
	SenderID     int                `json:"senderId" bson:"senderId"`            // 发送者ID   可以退回邮件
	SendTime     int                `json:"sendTime" bson:"sendTime"`            // 发送时间
	Validity     int64              `json:"validity" bson:"validity"`            // 邮件有效期
	Title        string             `json:"title"    bson:"title"`               // 邮件标题
	Content      string             `json:"content"   bson:"content"`            // 邮件文字内容
	State        int                `json:"state"    bson:"state"`               // 邮件状态：未读 0 已读 1 删除 2
	ReceiveState bool               `json:"receiveState"    bson:"receiveState"` // 领取状态：false未领取  true 已领取
	Attachment   []*Attachment      `json:"attachment"    bson:"attachment"`     // 附件
}

const (
	State_UNREAD = iota
	State_HAVE_READ
	State_delete
)

func (mail *PlayerMailModel) CollectionName() string {
	return config.CollMail
}

func (mail *PlayerMailModel) GetId() primitive.ObjectID {
	return mail.Id
}

func (mail *PlayerMailModel) SetId(id primitive.ObjectID) {
	mail.Id = id
}

// GetPlayerMailList 获取邮箱列表
func (mail *PlayerMailModel) GetPlayerMailList(filter bson.D, options *options.FindOptions) (data []*PlayerMailModel, err error) {
	finder := mgoDB.NewFinder(mail).Where(filter).Options(options).Records(&data)
	err = mgoDB.GetMgo().FindMany(context.TODO(), finder)
	if err != nil {
		return nil, err
	}
	return
}

// InsertPlayerMail 新增邮件
func (mail *PlayerMailModel) InsertPlayerMail() error {
	err := mgoDB.GetMgo().InsertOne(nil, mail)
	return err
}

// UpdatePlayerMail 更新邮件
func (mail *PlayerMailModel) UpdatePlayerMail(filter bson.D, update bson.D, opt *options.UpdateOptions) error {
	updater := mgoDB.NewUpdater(mail).Where(filter).Update(update).Options(opt)
	_, err := mgoDB.GetMgo().UpdateOne(nil, updater)
	return err
}

func (mail *PlayerMailModel) GetTotalNum(filter bson.D) (int64, error) {
	finder := mgoDB.NewCounter(mail).Where(filter)
	return mgoDB.GetMgo().CountDocuments(nil, finder)
}

// SetMailRead 设置邮件已读
func (mail *PlayerMailModel) SetMailRead(id interface{}, userId int, opt *options.UpdateOptions) error {
	return mail.UpdatePlayerMail(bson.D{{"mid", id}, {"uid", userId}}, bson.D{{"state", State_HAVE_READ}}, opt)
}

// SetMailDel 设置邮件删除
func (mail *PlayerMailModel) SetMailDel(id string, userId int, opt *options.UpdateOptions) error {
	return mail.UpdatePlayerMail(bson.D{{"mid", id}, {"uid", userId}}, bson.D{{"state", State_delete}}, opt)
}

// SendNewPlayerMail 发送新注册用户邮件
func (mail *PlayerMailModel) SendNewPlayerMail(userId int) {
	mail.UserId = userId
	mail.Name = "欢迎注册"
	mail.Sender = "系统"
	mail.SenderID = 00000
	mail.SendTime = int(time.Now().Unix())
	mail.Validity = 642000
	mail.Title = "欢迎注册"
	mail.Content = "欢迎注册欢迎注册欢迎注册欢迎注册欢迎注册欢迎注册欢迎注册欢迎注册"
	mail.State = State_UNREAD
	mail.InsertPlayerMail()
}
