package db

import (
	"com.xv.admin.server/config"
	"com.xv.admin.server/service/mgoDB"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type STYLE int

const (
	NoticeTypeLogin = 1 // 1 为永久公告
	NoticeTypeGame  = 2 // 2 为时限公告
)

type NoticeMsgModel struct {
	Id        primitive.ObjectID `json:"-"         bson:"-"`
	Mid       int                `json:"mid"  bson:"mid"`             // 公告ID
	Title     string             `json:"title"   bson:"title"`        // 公告名称
	Content   string             `json:"content" bson:"content"`      // 公告的字符内容
	Type      int                `json:"type"  bson:"type"`           // 公告类型  //1登录公告  2/游戏内公告
	StartTime int                `json:"startTime"  bson:"startTime"` // 开始时间
	EndTime   int                `json:"endTime"  bson:"endTime"`     // 结束时间
}

func (notice *NoticeMsgModel) CollectionName() string {
	return config.CollNotice
}

func (notice *NoticeMsgModel) GetId() primitive.ObjectID {
	return notice.Id
}

func (notice *NoticeMsgModel) SetId(id primitive.ObjectID) {
	notice.Id = id
}

// InsertNotice 新增数据
func (notice *NoticeMsgModel) InsertNotice() error {
	if notice.Mid == 0 {
		opt := options.FindOne().SetSort(bson.D{{"mid", -1}})
		Tmp, err := new(NoticeMsgModel).SelectOne(nil, opt)
		if err != nil {
			return err
		}
		if Tmp == nil {
			notice.Mid = 2
		} else {
			notice.Mid = Tmp.Mid + 1
		}
	}
	notice.Id = primitive.NewObjectID()
	err := mgoDB.GetMgo().InsertOne(nil, notice)
	return err
}

func (notice *NoticeMsgModel) SelectOne(filter bson.D, opt *options.FindOneOptions) (data *NoticeMsgModel, err error) {
	oneFilter := mgoDB.NewOneFinder(notice).Where(filter).Options(opt).Record(&data)
	res, err := mgoDB.GetMgo().FindOne(nil, oneFilter)
	if err != nil {
		return nil, err
	}
	if !res {
		return nil, nil
	}
	return data, err
}

// GetNotice 获取所有公告
func (notice *NoticeMsgModel) GetNotice(filter bson.D, options *options.FindOptions) (data []map[string]interface{}, err error) {
	data = []map[string]interface{}{}
	finder := mgoDB.NewFinder(notice).Options(options).Records(&data)
	if filter != nil {
		finder.Where(filter)
	}
	err = mgoDB.GetMgo().FindMany(context.TODO(), finder)
	return
}

// CheckAdvIdDoesItExist 检查advID是否重复
func (notice *NoticeMsgModel) CheckAdvIdDoesItExist(advId int, ruleOutId int) (bool, error) {
	finder := mgoDB.NewOneFinder(notice).Where(bson.D{{"advid", advId}})
	if ruleOutId > 0 {
		finder.Where(bson.D{{"$no in", ruleOutId}})
	}
	return mgoDB.GetMgo().FindOne(context.TODO(), finder)
}

// GetNoticeById 检查advID是否重复
func (notice *NoticeMsgModel) GetNoticeById(id int) bool {
	finder := mgoDB.NewOneFinder(notice).Where(bson.D{{"_id", id}})
	one, err := mgoDB.GetMgo().FindOne(context.TODO(), finder)
	if err != nil {
		return false
	}
	return one
}

// UpdateNotice 更新公告
func (notice *NoticeMsgModel) UpdateNotice() error {
	updater := mgoDB.NewUpdater(notice).Where(bson.D{{"mid", notice.Mid}}).Update(bson.D{
		{"title", notice.Title},
		{"content", notice.Content},
		{"type", notice.Type},
		{"startTime", notice.StartTime},
		{"endTime", notice.EndTime},
	})
	_, err := mgoDB.GetMgo().UpdateOne(nil, updater)
	return err
}

// GetTotalNum 按统计条数
func (notice *NoticeMsgModel) GetTotalNum(filter bson.D) (int64, error) {
	finder := mgoDB.NewCounter(notice).Where(filter)
	return mgoDB.GetMgo().CountDocuments(nil, finder)
}
