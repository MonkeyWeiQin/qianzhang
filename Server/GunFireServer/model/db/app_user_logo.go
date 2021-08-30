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

/**
app_user_log 表操作文件
*/

type BusinessType int8

const (
	BusinessTypeLogin    BusinessType = 1 //登录
	BusinessTypeRegister BusinessType = 2 //注册


	BusinessTypeSevenDay BusinessType = 4 //领取七天活动

)

type AppUserLogModel struct {
	Id           primitive.ObjectID `json:"-" bson:"-"`
	UserId       int              `json:"user_id" bson:"user_id"`
	Time         int64              `json:"time"   bson:"time"`
	Ip           string             `json:"ip"     bson:"ip"`
	DevId        string             `json:"devid"  bson:"devid"`
	BusinessType BusinessType       `json:"business_type" bson:"business_type"`
	Description  string       		`json:"description" bson:"description"`
}

func (AppUserLog *AppUserLogModel) CollectionName() string {
	return config.CollAppUserLog
}

func (AppUserLog *AppUserLogModel) GetId() primitive.ObjectID {
	return AppUserLog.Id
}

func (AppUserLog *AppUserLogModel) SetId(id primitive.ObjectID) {
	AppUserLog.Id = id
}

// GetAppUserLogList 获取日志列表
func (AppUserLog *AppUserLogModel) GetAppUserLogList(userId int64, start int64, end int64, BusinessType BusinessType, options *options.FindOptions) (data []AppUserLogModel, err error) {
	data = []AppUserLogModel{}
	finder := mgoDB.NewFinder(AppUserLog).Options(options).Records(&data)

	if userId > 0 {
		finder.Where(bson.D{{"user_id", userId}})
	}
	if BusinessType > 0 {
		finder.Where(bson.D{{"business_type", BusinessType}})
	}
	if start > 0 {
		finder.Where(bson.D{{"time", bson.M{"$gt": start}}})
	}
	if end > 0 {
		finder.Where(bson.D{{"time", bson.M{"$lt": end}}})
	}

	err = mgoDB.GetMgo().FindMany(context.TODO(), finder)
	if err != nil {
		return nil, err
	}
	return
}

// CreateAppUserLog 创建日志
func (AppUserLog *AppUserLogModel) CreateAppUserLog(userId int, BusinessType BusinessType, Ip string, DevId string , Description string) {
	AppUserLog.UserId = userId
	AppUserLog.BusinessType = BusinessType
	AppUserLog.Ip = Ip
	AppUserLog.DevId = DevId
	AppUserLog.Time = time.Now().Unix()
	AppUserLog.Description = Description

	err := mgoDB.GetMgo().InsertOne(nil, AppUserLog)
	if err != nil {
		return
	}
}

// GetTotalNum 按统计条数
func (AppUserLog *AppUserLogModel) GetTotalNum(userId int64, start int64, end int64, BusinessType BusinessType) (int64, error) {
	finder := mgoDB.NewCounter(AppUserLog)
	if userId > 0 {
		finder = finder.Where(bson.D{{"user_id", userId}})
	}
	if BusinessType > 0 {
		finder = finder.Where(bson.D{{"business_type", BusinessType}})
	}
	if start > 0 {
		finder = finder.Where(bson.D{{"time", bson.M{"$gt": start}}})
	}
	if end > 0 {
		finder = finder.Where(bson.D{{"time", bson.M{"$lt": end}}})
	}
	return mgoDB.GetMgo().CountDocuments(nil, finder)
}
