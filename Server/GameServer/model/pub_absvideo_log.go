package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	AbsTypeRandom     = 1
	AbsTypeMultiple   = 2
	AbsTypeAwards     = 3
	AbsTypeFixedValue = 4
	AbsTypeMustSee    = 5
)

// BrowseABSInfo 记录玩家的观看视频的相关信息
type BrowseABSInfo struct {
	AbsId     int     `json:"abs_id"    bson:"abs_id"`   // 广告ID
	TimeStamp int64   `json:"stamp"     bson:"stamp"`    // 时间戳
	Gold      float64 `json:"gold"      bson:"gold"`     // 获得的金币数
	AbsType   int     `json:"abs_type"  bson:"abs_type"` // 奖励种类：随机 倍数 次数 固定值 必须看
	Times     int     `json:"times"    bson:"times"`     // 已经观看次数
}

// UserAbsVideoDataModel 记录玩家的观看视频的相关信息
type UserAbsVideoDataModel struct {
	Id            primitive.ObjectID `json:"-"                bson:"-"`
	Day           string             `json:"day"              bson:"day"` // 日期 2020-10-10
	UserId        int64              `json:"user_id"          bson:"user_id"`
	BrowseABSInfo []BrowseABSInfo    `json:"browse_abs_info"  bson:"browse_abs_info"`
}

var (
	//absVideoColl = new(AbsVideoCollection)
)
//
//type AbsVideoCollection struct{}
//
//func (m *AbsVideoCollection) CollectionName() string {
//	return TableNameUser
//}
//func GetAbsVideoCollection() *AbsVideoCollection {
//	return absVideoColl
//}
//
//func (UserAdsData *UserAbsVideoDataModel) SetId(id primitive.ObjectID) {
//	UserAdsData.Id = id
//}
//
//// InitUserAdsVideoLogByDay 初始化当天数据
//func (m *AbsVideoCollection) InitUserAdsVideoLogByDay(absVideo *UserAbsVideoDataModel) error {
//	err := mgoDB.GetMgo().InsertOne(nil, m,absVideo)
//	return err
//}
//
//// AddUserAdsVideoLog 增加浏览记录
//func (m *AbsVideoCollection) AddUserAdsVideoLog(uid int , day string,BrowseABSInfo BrowseABSInfo) error {
//	updater := mgoDB.NewUpdater(m).Where(bson.D{{"user_id", uid}, {"day", day}}).Push(bson.D{{"browse_abs_info", BrowseABSInfo}})
//	_, err := mgoDB.GetMgo().PushOne(nil, updater)
//	return err
//}
//
//func (m *AbsVideoCollection) GetList(filter bson.D, opt *options.FindOptions) (data []map[string]interface{}, err error) {
//	data = []map[string]interface{}{}
//	finder := mgoDB.NewFinder(m).Where(filter).Options(opt).Records(&data)
//	err = mgoDB.GetMgo().FindMany(context.TODO(), finder)
//	return
//}
//
//func (m *AbsVideoCollection) CheckUserData(userId int64, day string) error {
//	if day == "" {
//		day = time.Now().Format("2006-01-02")
//	}
//	filter := bson.D{{"user_id", userId}, {"day", day}}
//	finder := mgoDB.NewOneFinder(m).Where(filter)
//	res, err := mgoDB.GetMgo().FindOne(context.TODO(), finder)
//	if err != nil {
//		return err
//	}
//	if !res {
//		absVideo :=new(UserAbsVideoDataModel)
//		absVideo.UserId = userId
//		absVideo.Day = day
//		absVideo.BrowseABSInfo = []BrowseABSInfo{}
//		InsertUserAdsVideoLogErr := m.InitUserAdsVideoLogByDay(absVideo)
//		if InsertUserAdsVideoLogErr != nil {
//			return InsertUserAdsVideoLogErr
//		}
//	}
//	return err
//}
//
//// GetTotalNum 获取数量
//func (m *AbsVideoCollection) GetTotalNum(userId int64, AvType int, day string) (int64, error) {
//	filter := bson.D{{"user_id", userId}}
//	if AvType > 0 {
//		filter = append(filter, bson.E{Key: "browse_abs_info.av_type", Value: AvType})
//	}
//	if day == "" {
//		filter = append(filter, bson.E{Key: "day", Value: day})
//	}
//	finder := mgoDB.NewCounter(m).Where(filter)
//	return mgoDB.GetMgo().CountDocuments(nil, finder)
//}
//
//// GetAbsVideoDataList 获取数据列表
//func (m *AbsVideoCollection) GetAbsVideoDataList(userId int64, day string, AbsType int) ([]map[string]interface{}, error) {
//	where := bson.D{}
//	if userId != 0 {
//		where = append(where, bson.E{Key: "user_id", Value: userId})
//	}
//	if AbsType != 0 {
//		where = append(where, bson.E{Key: "browse_abs_info.abs_type", Value: AbsType})
//	}
//	if day != "" {
//		where = append(where, bson.E{Key: "day", Value: day})
//	} else {
//		where = append(where, bson.E{Key: "day", Value: time.Now().Format("2006-01-02")})
//	}
//	option := new(options.FindOptions)
//	return m.GetList(where, option)
//}
