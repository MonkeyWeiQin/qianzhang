package public

import (
	"com.xv.admin.server/model/db"
	"com.xv.admin.server/utils"
	"github.com/gin-gonic/gin"
)

// AdvertisingVideoCsv 广告视频奖励的表格
type AdvertisingVideoCsv struct {
	AvId     int     `csv:"ad_id"`    // 广告编号
	Name     string  `csv:"name"`     // 内部使用，不显示
	AvType   int     `csv:"av_type"`  // 奖励种类：随机 倍数 次数 固定值 必须看
	Multiple float64 `csv:"multiple"` // 倍数
	MinGold  float64 `csv:"min_gold"` // 最小数
	MaxGold  float64 `csv:"max_gold"` // 最大数
	Times    int     `csv:"times"`    // 次数
}

type CreateAdvertisingVideoRequest struct {
	AvId     int     `json:"ad_id"    binding:"required"` // 广告编号
	Name     string  `json:"name"     binding:"required"` // 内部使用，不显示
	AvType   int     `json:"av_type"  binding:"required"` // 奖励种类：随机 倍数 次数 固定值 必须看
	Multiple float64 `json:"multiple" binding:"required"` // 倍数
	MinGold  float64 `json:"min_gold" binding:"required"` // 最小数
	MaxGold  float64 `json:"max_gold" binding:"required"` // 最大数
	Times    int     `json:"times"    binding:"required"` // 次数
}

type UserAdsVideoLogRequest struct {
	AvId     int     `json:"av_id"      binding:"required"`   // 广告编号
	TaskType int     `json:"task_type"  binding:"required"`   // 任务类型，日常 1 成就 2 登陆 3
	TaskId   int     `json:"task_id"    binding:"required"`   // 对应的任务id号
}
type UsersBrowseAVTypeResponse struct {
	//State int       `json:"state"      bson:"state"` // 状态：成功 1 ， 失败 0
	AvType  int     `json:"av_type"    bson:"av_type"` // 奖励种类：随机 倍数 次数 固定值 必须看
	Times int       `json:"times"      bson:"times"`   // 剩余次数
	Gold  float64   `json:"gold"       bson:"gold"`    // 获得金币数
}

func GetAdsVideoConfig(c *gin.Context) {
	AdsVideoConfig, GetConfigErr := new(db.AdvertisingVideoConfigModel).GetConfig()
	if GetConfigErr != nil {
		panic(GetConfigErr.Error())
	}
	utils.OkWithData(c, AdsVideoConfig)
}

func CreateAdsVideoConfig(c *gin.Context) {
	CreateAdvertisingVideoRequest := new(CreateAdvertisingVideoRequest)
	err := c.ShouldBindJSON(CreateAdvertisingVideoRequest)
	if err != nil {
		panic(10001)
	}
	adsVideoConfigModel := new(db.AdvertisingVideoConfigModel)
	res, FindConfigByAvIdOrNameErr := adsVideoConfigModel.FindConfigByAvIdOrName(CreateAdvertisingVideoRequest.AvId, CreateAdvertisingVideoRequest.Name)
	if FindConfigByAvIdOrNameErr != nil {
		panic(FindConfigByAvIdOrNameErr.Error())
	}
	if res {
		panic("Id或者name重复")
	}
	config := make([]interface{}, 0)
	config = append(config, CreateAdvertisingVideoRequest)
	SetConfigErr := adsVideoConfigModel.CreateConfig(config)
	if SetConfigErr != nil {
		panic("操作失败:" + SetConfigErr.Error())
	}
	utils.OkWithData(c, CreateAdvertisingVideoRequest)
}

//func UserAdsVideoLog(c *gin.Context) {
//	UserAdsVideoLogRequest := new(UserAdsVideoLogRequest)
//
//	err := c.ShouldBindJSON(UserAdsVideoLogRequest)
//	if err != nil {
//		panic(err.Error())
//	}
//	userId, _ := c.Get("userId")
//	AdvertisingVideoConfigModel:=new(db.AdvertisingVideoConfigModel)
//	_, FindConfigByAvIdOrNameErr := AdvertisingVideoConfigModel.FindConfigByAvIdOrName(UserAdsVideoLogRequest.AvId , "")
//	if FindConfigByAvIdOrNameErr != nil {
//		panic("广告数据未找到")
//	}
//	UserAdvertisingVideoDataModel := new(db.UserAdvertisingVideoDataModel)
//
//	//初始化当天的数据
//	if CheckUserDataErr := UserAdvertisingVideoDataModel.CheckUserData(userId.(int64), ""); CheckUserDataErr != nil {
//		panic(CheckUserDataErr.Error())
//	}
//
//	gold , times , HandleErr := adsvideo.Handle(AdvertisingVideoConfigModel , userId.(int64))
//	if HandleErr != nil {
//		panic(HandleErr.Error())
//	}
//
//	BrowseABSInfo := db.UserAdvertisingVideoInfo{
//		Gold: gold,
//		AvId: AdvertisingVideoConfigModel.AvId,
//		TimeStamp: time.Now().Unix(),
//		AvType:AdvertisingVideoConfigModel.AvType,
//	}
//	if AddUserAdsVideoLogErr := UserAdvertisingVideoDataModel.AddUserAdsVideoLog(BrowseABSInfo); AddUserAdsVideoLogErr != nil {
//		panic(AddUserAdsVideoLogErr.Error())
//	}
//
//	utils.OkWithData(c, &UsersBrowseAVTypeResponse{
//		AvType: AdvertisingVideoConfigModel.AvType,
//		Times:  times,
//		Gold:   gold,
//	})
//}



func GetUserAdsVideoConfig(c *gin.Context) {
	//userId, _ := c.Get("userId")
	//UserAdvertisingVideoDataModel := new(db.UserAdvertisingVideoDataModel)
	//UserAdvertisingVideoDataModel.

}