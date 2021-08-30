package db

import (
	"com.xv.admin.server/config"
	"com.xv.admin.server/service/mgoDB"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// BroadcastModel 公告信息数据
type BroadcastModel struct {
	Id          primitive.ObjectID `json:"-"  bson:"-"`
	Mid         int                `json:"mid" bson:"mid"`                 //广播ID
	Content     string             `json:"content" bson:"content"`         //公告的字符内容
	Color       string             `json:"color" bson:"color"`             //展示颜色
	StartTime   int                `json:"startTime" bson:"startTime"`     //开始时间
	EndTime     int                `json:"endTime" bson:"endTime"`         //结束时间
	SpacingTime int64              `json:"spacingTime" bson:"spacingTime"` //展示时间  单位：秒
}

func (broadcast *BroadcastModel) CollectionName() string {
	return config.CollBroadcast
}

func (broadcast *BroadcastModel) GetId() primitive.ObjectID {
	return broadcast.Id
}

func (broadcast *BroadcastModel) SetId(id primitive.ObjectID) {
	broadcast.Id = id
}

// InsertBroadcast 新增数据
func (broadcast *BroadcastModel) InsertBroadcast() error {
	if broadcast.Mid == 0 {
		opt := options.FindOne().SetSort(bson.D{{"mid", -1}})
		Tmp, err := new(BroadcastModel).SelectOne(nil, opt)
		if err != nil {
			return err
		}
		if Tmp == nil {
			broadcast.Mid = 2
		} else {
			broadcast.Mid = Tmp.Mid + 1
		}
	}
	broadcast.Id = primitive.NewObjectID()
	err := mgoDB.GetMgo().InsertOne(nil, broadcast)
	return err
}

func (broadcast *BroadcastModel) SelectOne(filter bson.D, opt *options.FindOneOptions) (data *NoticeMsgModel, err error) {
	oneFilter := mgoDB.NewOneFinder(broadcast).Where(filter).Options(opt).Record(&data)
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
func (broadcast *BroadcastModel) GetNotice(filter bson.D, options *options.FindOptions) (data []map[string]interface{}, err error) {
	data = []map[string]interface{}{}
	finder := mgoDB.NewFinder(broadcast).Options(options).Records(&data)
	if filter != nil {
		finder.Where(filter)
	}
	err = mgoDB.GetMgo().FindMany(nil, finder)
	return
}

// UpdateNotice 更新公告
func (broadcast *BroadcastModel) UpdateBroadcast() error {
	updater := mgoDB.NewUpdater(broadcast).Where(bson.D{{"mid", broadcast.Mid}}).Update(bson.D{
		{"content", broadcast.Content},
		{"startTime", broadcast.StartTime},
		{"endTime", broadcast.EndTime},
		{"spacingTime", broadcast.SpacingTime},
		{"color", broadcast.Color},
	})
	_, err := mgoDB.GetMgo().UpdateOne(nil, updater)
	return err
}

// GetTotalNum 按统计条数
func (broadcast *BroadcastModel) GetTotalNum(filter bson.D) (int64, error) {
	finder := mgoDB.NewCounter(broadcast).Where(filter)
	return mgoDB.GetMgo().CountDocuments(nil, finder)
}
