package activity

import (
	"com.xv.admin.server/model/db"
	"com.xv.admin.server/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
	"time"
)

//--------------------------------------------------------登陆奖励活动-----------------------------------------------------------------

type SetSevenDayRuleRequest struct {
	Data []*SevenDayConfig `json:"data" binding:"required"` //
}
type SevenDayConfig struct {
	Times      	int   	`json:"times"       binding:"required"`        // 次数
	AwardType 	int   	`json:"awardType"   binding:"required"`   //  奖励类型 （金币，体力，道具）
	AwardId   	int 	`json:"awardId"`  	//  道具Id
	AwardCount  int64 	`json:"awardCount"  binding:"required"`  //  数量
}

type SevenDayListResponse struct {
	Times      	int  	`json:"times"`      //次数
	AwardType 	int  	`json:"awardType"` // 奖励类型 （金币，体力，道具）
	AwardId   	int  	`json:"awardId" `  // 道具Id
	AwardCount  int64	`json:"amount" `   // 数量
	State     	bool 	`json:"state" `    // 领取状态  true 已经领取  false未领取
}

type SignInRequest struct {
	Times int `json:"times"` // 次数
}

// SetSevenDayRule 设置活动规则
func SetSevenDayRule(c *gin.Context) {
	SetSevenDayRuleRequest := new(SetSevenDayRuleRequest)
	err := c.ShouldBindJSON(SetSevenDayRuleRequest)
	if err != nil {
		panic(10001)
	}
	if len(SetSevenDayRuleRequest.Data) < 7 {
		panic("请配置完整")
	}

	SevenDayConfig := make([]interface{}, 0)

	for index, item := range SetSevenDayRuleRequest.Data {
		if index != item.Times-1 {
			panic("第" + strconv.Itoa(item.Times) + "项的配置与其他项重复")
		}
		if !(item.AwardType == db.AwardTypeGold || item.AwardType == db.AwardTypeDiamond || item.AwardType == db.AwardTypeProps) {
			panic("第" + strconv.Itoa(item.Times) + "项的奖励类型配置有误")
		}
		if db.AwardTypeProps == item.AwardType {
			if item.AwardId  < 0 {
				//todo  检查商品ID是否存在
				panic("第" + strconv.Itoa(item.Times) + "项的道具ID:" + strconv.Itoa(item.Times) + "不存在")
			}
		}
		SevenDayConfig = append(SevenDayConfig, &db.SevenDayConfigModel{
			Times:      	item.Times,
			AwardType: 		db.AwardType(item.AwardType),
			AwardId:		item.AwardId,
			AwardCount:   	item.AwardCount,
		})
	}
	SevenDayConfigModel := new(db.SevenDayConfigModel)
	SetConfigErr := SevenDayConfigModel.SetConfig(SevenDayConfig)
	if SetConfigErr != nil {
		panic("操作失败:" + SetConfigErr.Error())
	}
	utils.Ok(c)
}

// GetSevenDay 获取活动列表
func GetSevenDay(c *gin.Context) {
	userId, _ := c.Get("userId")
	activity := make([]SevenDayListResponse, 0)

	AppUserLog 	:= new(db.AppUserLogModel)
	Count, GetTotalNumErr := AppUserLog.GetTotalNum(userId.(int64), 0, 0, db.BusinessTypeSevenDay)
	if GetTotalNumErr != nil {
		panic(GetTotalNumErr.Error())
	}

	SevenDayConfigModel := new(db.SevenDayConfigModel)
	option 				:= new(options.FindOptions).SetSkip((Count/7) * 7 ).SetLimit(7)
	config, GetConfigErr := SevenDayConfigModel.GetConfig(option)
	if GetConfigErr != nil {
		panic(GetConfigErr.Error())
	}

	for _, item := range config {
		amount,_:= GetPrize(item.AwardCount , item.AwardType)
		state   := false
		if int(Count) >= item.Times {
			state = true
		}
		activity = append(activity, SevenDayListResponse{
			AwardType: int(item.AwardType),
			AwardCount:    amount,
			Times:      item.Times,
			State:     	state,
		})
	}
	utils.OkWithData(c, activity)
}

// SevenDaySignIn 签到
func SevenDaySignIn(c *gin.Context) {
	userId, _ := c.Get("userId")
	SignInRequest := new(SignInRequest)
	err := c.ShouldBindJSON(SignInRequest)
	if err != nil {
		panic(10001)
	}
	currentTime := time.Now()
	startTime 	:= time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 00, 00, 00, 0, currentTime.Location())
	endTime 	:= time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 23, 59, 59, 0, currentTime.Location())
	AppUserLog 	:= new(db.AppUserLogModel)

	ToDayNum, GetToDayTotalNumErr := AppUserLog.GetTotalNum(userId.(int64), startTime.Unix(), endTime.Unix(), db.BusinessTypeSevenDay)
	if GetToDayTotalNumErr != nil {
		panic(GetToDayTotalNumErr.Error())
	}
	if ToDayNum > 0 {
		//panic("今天已经领取过了")
	}

	num, GetTotalNumErr := AppUserLog.GetTotalNum(userId.(int64), 0, 0, db.BusinessTypeSevenDay)

	if GetTotalNumErr != nil {
		panic(GetTotalNumErr.Error())
	}
	SevenDayConfigModel := new(db.SevenDayConfigModel)
	res, GetConfigErr := SevenDayConfigModel.GetConfigByTimes(SignInRequest.Times)
	if GetConfigErr != nil {
		panic(GetConfigErr.Error())
	}
	if !res {
		panic("获取奖品列表失败")
	}
	if int(num+1) < SignInRequest.Times {
		panic("该奖品已被领取过了")
	}
	if int(num+1) > SignInRequest.Times {
		panic("请先领取上一个奖品")
	}

	//_, message := GetPrize(SevenDayConfigModel.AwardCount, SevenDayConfigModel.AwardType)
	//AppUserLog.CreateAppUserLog(userId.(int64), db.BusinessTypeSevenDay, c.ClientIP(), "", "领取7天签到活动奖励:获得奖品为"+message)
	//utils.OkWithMessage(c, "领取成功")
}


func GetPrize(amount int64, awardType db.AwardType) (int64, string) {
	message := ""
	if awardType == db.AwardType(db.AwardTypeGold) {
		message = strconv.Itoa(int(amount)) + "金币"
	}
	if awardType == db.AwardType(db.AwardTypeDiamond) {
		message =  strconv.Itoa(int(amount)) + "钻石"

	}
	if awardType == db.AwardType(db.AwardTypeProps) {
		message =  "道具AAA"
	}
	return amount,message
}