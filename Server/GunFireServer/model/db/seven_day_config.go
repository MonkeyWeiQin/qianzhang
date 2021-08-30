package db

import (
	"com.xv.admin.server/config"
	"com.xv.admin.server/service/mgoDB"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AwardType int //奖品类型

const (
	AwardTypeGold       int = 1  //金币
	AwardTypeDiamond	int = 2  //钻石
	AwardTypeProps      int = 3  //物品
)

// SevenDayConfigModel  7天签到活动配置
type SevenDayConfigModel struct {
	Id        	primitive.ObjectID `json:"-"           bson:"-"`
	Times      	int                `json:"times"       bson:"times"`             		// 关卡数
	AwardType 	AwardType          `json:"award_type"  bson:"award_type"`   		// 奖励类型
	AwardId   	int                `json:"award_td"    bson:"award_td"`      		// 奖励物品ID
	AwardCount  int64              `json:"award_count" bson:"award_count"`     // 奖励数量
}

func (sevenDayConfig *SevenDayConfigModel) CollectionName() string {
	return config.CollSevenDayConfig
}

func (sevenDayConfig *SevenDayConfigModel) GetId() primitive.ObjectID {
	return sevenDayConfig.Id
}

func (sevenDayConfig *SevenDayConfigModel) SetId(id primitive.ObjectID) {
	sevenDayConfig.Id = id
}

// SetConfig 设置规则
func (sevenDayConfig *SevenDayConfigModel) SetConfig(cols []interface{}) error {
	_, DelConfigErr := sevenDayConfig.DelConfig()
	if DelConfigErr != nil {
		return DelConfigErr
	}
	err := mgoDB.GetMgo().InsertMany(nil, cols)
	return err
}

func (sevenDayConfig *SevenDayConfigModel) DelConfig() (int64, error) {
	filter := mgoDB.NewDeleter(sevenDayConfig)
	return mgoDB.GetMgo().DeleteMany(nil, filter)
}

func (sevenDayConfig *SevenDayConfigModel) GetConfig(options *options.FindOptions) (data []SevenDayConfigModel, err error) {
	data = []SevenDayConfigModel{}
	finder := mgoDB.NewFinder(sevenDayConfig).Options(options).Records(&data)
	err = mgoDB.GetMgo().FindMany(context.TODO(), finder)
	if err != nil {
		return nil, err
	}
	return
}

// GetConfigByTimes 获取单个配置
func (sevenDayConfig *SevenDayConfigModel) GetConfigByTimes(Times int) (res bool , err error) {
	filter := bson.D{{"times",Times}}
	finder := mgoDB.NewOneFinder(sevenDayConfig).Where(filter)
	res, err = mgoDB.GetMgo().FindOne(context.TODO(), finder)
	return
}

