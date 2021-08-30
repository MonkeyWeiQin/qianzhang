package db

import (
	"com.xv.admin.server/config"
	"com.xv.admin.server/service/mgoDB"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	AdvertisingVideoDataMap = make(map[int]*AdvertisingVideoConfigModel)
)

type AdvertisingVideoConfigModel struct {
	Id       primitive.ObjectID `json:"-"        bson:"-"`
	AvId     int                `json:"av_id"    bson:"av_id"`
	Name     string             `json:"name"     bson:"name"`     // 内部使用，不显示
	AvType   int                `json:"av_type"  bson:"av_type"`  // 奖励种类：随机 倍数 次数 固定值 必须看
	Multiple float64            `json:"multiple" bson:"multiple"` // 倍数
	MinGold  float64            `json:"min_gold" bson:"min_gold"` // 最小数
	MaxGold  float64            `json:"max_gold" bson:"max_gold"` // 最大数
	Times    int                `json:"times"    bson:"times"`    // 次数
}

func (AdvertisingVideoConfig *AdvertisingVideoConfigModel) CollectionName() string {
	return config.CollAdvertisingVideoConfig
}

func (AdvertisingVideoConfig *AdvertisingVideoConfigModel) GetId() primitive.ObjectID {
	return AdvertisingVideoConfig.Id
}

func (AdvertisingVideoConfig *AdvertisingVideoConfigModel) SetId(id primitive.ObjectID) {
	AdvertisingVideoConfig.Id = id
}

func (AdvertisingVideoConfig *AdvertisingVideoConfigModel) GetConfig() (data []map[string]interface{}, err error) {
	data = []map[string]interface{}{}
	finder := mgoDB.NewFinder(AdvertisingVideoConfig).Records(&data)
	err = mgoDB.GetMgo().FindMany(context.TODO(), finder)
	if err != nil {
		return nil, err
	}
	return
}

func (AdvertisingVideoConfig *AdvertisingVideoConfigModel) FindConfig(AvId int, Name string) (res bool , err error) {
	filter := bson.D{}
	if AvId > 0 {
		filter = append(filter , bson.E{Key: "av_id",Value: AvId})
	}
	if Name != "" {
		filter = append(filter , bson.E{Key: "name",Value: Name})
	}
	finder := mgoDB.NewOneFinder(AdvertisingVideoConfig).Where(filter)
	res, err = mgoDB.GetMgo().FindOne(context.TODO(), finder)
	return
}

// CreateConfig 新增配置
func (AdvertisingVideoConfig *AdvertisingVideoConfigModel) CreateConfig(cols []interface{}) error {
	err := mgoDB.GetMgo().InsertMany(nil, cols)
	return err
}

// DelConfig 删除配置
func (AdvertisingVideoConfig *AdvertisingVideoConfigModel) DelConfig(res int64, err error) {
	finder := mgoDB.NewDeleter(AdvertisingVideoConfig)
	res, err = mgoDB.GetMgo().DeleteMany(context.TODO(), finder)
	return
}

func (AdvertisingVideoConfig *AdvertisingVideoConfigModel) FindConfigByAvIdOrName(AvId int, Name string) (res bool , err error) {
	filter := bson.D{{"$or",[]bson.M{{"av_id":AvId},{"name":Name}}}}
	finder := mgoDB.NewOneFinder(AdvertisingVideoConfig).Where(filter)
	res, err = mgoDB.GetMgo().FindOne(context.TODO(), finder)
	return
}
