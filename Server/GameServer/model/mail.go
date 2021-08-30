package model

import (
	"battle_rabbit/define"
	"battle_rabbit/global"
	"battle_rabbit/service/mgoDB"
	"battle_rabbit/utils/xid"
	"context"
	"fmt"
	bson "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var (
	mailColl = newMailColl()
)

type MailCollection struct {
	*mgoDB.DbBase
}

func newMailColl() *MailCollection {
	return &MailCollection{mgoDB.NewDbBase(define.MgoDBNameBattle, define.TableNameBroadcast)}
}
func GetMailCollection() *MailCollection {
	return mailColl
}

/*
全局邮件: 由后台管理员发送,每个全局邮件在系统邮件中只存在一份.每个玩家登录后先读取全局邮件,如有系统邮件,则将系统邮件转为个人邮件并保存
个人邮件: 玩家每次从个人邮件中读邮件列表并做标记
*/
type MailModel struct {
	Mid          string         `json:"mid"           bson:"_id"`
	UserId       int            `json:"uid"           bson:"uid"`          // -1 为全局邮件邮件
	Name         string         `json:"name"          bson:"name"`         // 邮件名称
	Sender       string         `json:"sender"        bson:"sender"`       // 发送者名称
	SenderID     int            `json:"senderId"      bson:"senderId"`     // 发送者ID   可以退回邮件 (-1 为系统管理发送)
	SendTime     int            `json:"sendTime"      bson:"sendTime"`     // 发送时间
	Validity     int64          `json:"validity"      bson:"validity"`     // 邮件有效期
	Title        string         `json:"title"         bson:"title"`        // 邮件标题
	Content      string         `json:"content"       bson:"content"`      // 邮件文字内容
	State        int            `json:"state"         bson:"state"`        // 邮件状态：未读 0 已读 1 删除 2
	ReceiveState bool           `json:"receiveState"  bson:"receiveState"` // 领取状态：false未领取  true 已领取
	Attachment   []*global.Item `json:"attachment"    bson:"attachment"`   // 附件
}

type AttachmentItem struct {
	Mid        string         `json:"mid" bson:"mid"`
	Attachment []*global.Item `json:"attachment" bson:"attachment"` //附件
}

const (
	STATE_UNREAD    int = 0
	STATE_HAVE_READ int = 1
	STATE_DELETE    int = 2
)

// GetPlayerMailList 获邮件列表
func (m *MailCollection) GetPlayerMailList(ctx context.Context, userId int, options ...*options.FindOptions) (data []*MailModel, err error) {
	err = m.FindAll(ctx, bson.D{
		{"uid", userId},
		{"state", bson.M{"$ne": STATE_DELETE}},
		{"validity", bson.M{"$gt": time.Now().Unix()}}}, &data, options...)
	return
}

// 获取用户有效期内，未删除未领取的附件
func (m *MailCollection) GetPlayerMailValidityAttachment(ctx context.Context, userId int) (data []*AttachmentItem, err error) {
	err = m.FindAll(ctx, bson.D{
		{"uid", userId},
		{"receiveState", false},
		{"state", bson.M{"$ne": STATE_DELETE}},
		{"validity", bson.M{"$gt": time.Now().Unix()}},
	}, &data, options.Find().SetProjection(bson.M{"_id": 0, "mid": 1, "attachment": 1}))
	return
}



// UpdatePlayerMail 更新邮件
func (m *MailCollection) UpdatePlayerMail(filter interface{}, update interface{}, opt *options.UpdateOptions) (int64, error) {
	return m.UpdateMany(context.TODO(), filter, update, opt)
	//updater := mgoDB.NewUpdater(m).Where(filter).Update(update).Options(opt)
	//res, err := mgoDB.GetMgo().UpdateMany(nil, updater)
	//return res, err
}

// GetTotalNum 按统计条数
func (m *MailCollection) GetTotalNum(uid int, opt *options.FindOptions) (int, error) {
	num, err := mailColl.CountDocuments(nil, bson.M{"uid": uid})
	return int(num), err
}

// SetMailReceive 设置邮件已领取
func (m *MailCollection) SetMailReceive(idArr []string, userId int) error {

	n, err := m.UpdateMany(nil, bson.M{"mid": bson.M{"$in": idArr}, "uid": userId}, bson.M{"state": STATE_HAVE_READ, "receiveState": true}, nil)
	if err != nil || n == 0 {
		return fmt.Errorf("设置邮件为状态失败:: %v , n = %d ", err, n)
	}
	return nil
}

// SetMailDel 设置邮件删除
func (m *MailCollection) SetMailDel(idArr []string, userId int) (bool, error) {
	playerMail, err := m.UpdatePlayerMail(bson.D{{"mid", bson.M{"$in": idArr}}, {"uid", userId}}, bson.D{{"state", STATE_DELETE}}, nil)
	if err != nil {
		return false, err
	}
	if playerMail <= 0 {
		return false, nil
	}
	return true, nil
}

// 设置邮件状态 已读, 已删
func (m *MailCollection) SetMailState(idArr []string, userId, state int) error {
	if len(idArr) == 1 {
		n, err := m.UpdateOne(nil, bson.M{"mid": idArr[0]}, bson.M{"state": state})
		if err != nil || n == 0 {
			return fmt.Errorf("设置邮件状态失败:: %v , n = %d ", err, n)
		}
		return nil

	} else {
		n, err := m.UpdateMany(nil, bson.M{"mid": bson.M{"$in": idArr}, "uid": userId}, bson.M{"state": state}, nil)
		if err != nil || n == 0 {
			return fmt.Errorf("设置邮件为状态失败:: %v , n = %d ", err, n)
		}
		return nil
	}
}

// SendNewUserMail 发送新注册用户邮件
func (m *MailCollection) SendNewUserMail(userId int) {
	mail := new(MailModel)
	mail.UserId = userId
	mail.Mid = xid.New().String()
	mail.Name = "欢迎注册"
	mail.Sender = "系统"
	mail.SenderID = 00000
	mail.SendTime = int(time.Now().Unix())
	mail.Validity = 642000
	mail.Title = "欢迎注册"
	mail.Content = "欢迎注册欢迎注册欢迎注册欢迎注册欢迎注册欢迎注册欢迎注册欢迎注册"
	mail.State = STATE_UNREAD
	mail.ReceiveState = false
	mail.Attachment = append(mail.Attachment, &global.Item{
		ItemID:   "0",
		Count:    12323,
		ItemType: 1,
	})
	mail.Attachment = append(mail.Attachment, &global.Item{
		ItemID:   "0",
		Count:    123232,
		ItemType: 2,
	})
	err := m.InsertOne(nil,mail)
	if err != nil {
		return
	}
}

func (m *MailCollection) SetAllMailRead(userId int) (bool, error) {
	playerMail, err := m.UpdatePlayerMail(bson.D{{"uid", userId}}, bson.D{{"state", STATE_HAVE_READ}}, nil)
	if err != nil {
		return false, err
	}
	if playerMail <= 0 {
		return false, nil
	}
	return true, nil
}
func (m *MailCollection) SetAllMailDel(userId int) (bool, error) {
	playerMail, err := m.UpdatePlayerMail(bson.D{
		{"uid", userId},
		{"$or", []*bson.M{
			{"receiveState": true},
			{"attachment.0": bson.M{"$exists": false}}}}}, bson.D{{"state", STATE_DELETE}}, nil)
	if err != nil {
		return false, err
	}
	if playerMail <= 0 {
		return false, nil
	}
	return true, nil
}

// GetMailById 获邮件列表
func (m *MailCollection) GetMailById(ctx context.Context, userId int, mid string) (data *MailModel, err error) {
	err = m.FindOne(ctx, bson.D{{"uid", userId}, {"mid", mid}}, &data)
	if err != nil {
		return nil, err
	}
	return
}

// GetMailAttachmentById 获邮件列表
func (m *MailCollection) GetMailAttachmentById(ctx context.Context, userId int, mid string) (data *global.Item, err error) {
	err = m.FindOne(ctx, bson.D{{"uid", userId}, {"mid", mid}}, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
