package db

import (
	"com.xv.admin.server/config"
	"com.xv.admin.server/service/mgoDB"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// UserAdvertisingVideoInfo 记录玩家的观看视频的相关信息
type UserAdvertisingVideoInfo struct {
	AvId      int     `json:"av_id"    bson:"av_id"`   // 广告ID
	TimeStamp int64   `json:"stamp"    bson:"stamp"`   // 时间戳
	Gold      float64 `json:"gold"     bson:"gold"`    // 获得的金币数
	AvType    int     `json:"av_type"  bson:"av_type"` // 奖励种类：随机 倍数 次数 固定值 必须看

	//Times     int   `json:"times"    bson:"times"`     // 已经观看次数
}

// UserAdvertisingVideoDataModel 记录玩家的观看视频的相关信息
type UserAdvertisingVideoDataModel struct {
	Id            primitive.ObjectID         `json:"-"                bson:"-"`
	Today         string                     `json:"day"              bson:"day"` // 日期 2020-10-10
	UserId        int64                      `json:"user_id"          bson:"user_id"`
	BrowseABSInfo []UserAdvertisingVideoInfo `json:"browse_abs_info"  bson:"browse_abs_info"`
}

func (UserAdsData *UserAdvertisingVideoDataModel) CollectionName() string {
	return config.CollUserAdvertisingVideoLog
}

func (UserAdsData *UserAdvertisingVideoDataModel) GetId() primitive.ObjectID {
	return UserAdsData.Id
}

func (UserAdsData *UserAdvertisingVideoDataModel) SetId(id primitive.ObjectID) {
	UserAdsData.Id = id
}

// InitUserAdsVideoLogByDay 初始化当天数据
func (UserAdsData *UserAdvertisingVideoDataModel) InitUserAdsVideoLogByDay() error {
	err := mgoDB.GetMgo().InsertOne(nil, UserAdsData)
	return err
}

// AddUserAdsVideoLog 增加浏览记录
func (UserAdsData *UserAdvertisingVideoDataModel) AddUserAdsVideoLog(UserAdvertisingVideoInfo UserAdvertisingVideoInfo) error {
	updater := mgoDB.NewUpdater(UserAdsData).Where(bson.D{{"user_id", UserAdsData.UserId}, {"day", UserAdsData.Today}}).Push(bson.D{{"browse_abs_info", UserAdvertisingVideoInfo}})
	_, err := mgoDB.GetMgo().PushOne(nil, updater)
	return err
}

//func (UserAdsData *UserAdvertisingVideoDataModel) GetUsersList(filter bson.D, opt *options.FindOptions) (data []map[string]interface{}, err error) {
//	//fmt.Println(user.CollectionName())
//	data = []map[string]interface{}{}
//	finder := mgoDB.NewFinder(user).Where(filter).Options(opt).Records(&data)
//	err = mgoDB.GetMgo().FindMany(context.TODO(), finder)
//	return
//}

func (UserAdsData *UserAdvertisingVideoDataModel) CheckUserData(userId int64, day string) error {
	if day == "" {
		day = time.Now().Format("2006-01-02")
	}
	filter := bson.D{{"user_id", userId}, {"day", day}}
	finder := mgoDB.NewOneFinder(UserAdsData).Where(filter)
	res, err := mgoDB.GetMgo().FindOne(context.TODO(), finder)
	if err != nil {
		return err
	}
	if !res {
		UserAdsData.UserId = userId
		UserAdsData.Today = day
		UserAdsData.BrowseABSInfo = []UserAdvertisingVideoInfo{}
		InsertUserAdsVideoLogErr := UserAdsData.InitUserAdsVideoLogByDay()
		if InsertUserAdsVideoLogErr != nil {
			return InsertUserAdsVideoLogErr
		}
	}
	return err
}

// GetTotalNum 获取数量
func (UserAdsData *UserAdvertisingVideoDataModel) GetTotalNum(userId int64, AvType int, day string) (int64, error) {
	filter := bson.D{{"user_id", userId}}
	if AvType > 0 {
		filter = append(filter, bson.E{Key: "browse_abs_info.av_type", Value: AvType})
	}
	if day == "" {
		filter = append(filter, bson.E{Key: "day", Value: day})
	}
	finder := mgoDB.NewCounter(UserAdsData).Where(filter)
	return mgoDB.GetMgo().CountDocuments(nil, finder)
}
