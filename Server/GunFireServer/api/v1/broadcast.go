package v1

import (
	"com.xv.admin.server/model/db"
	"com.xv.admin.server/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UpdateBroadcastRequest 游戏系统消息
type UpdateBroadcastRequest struct {
	Mid         int    `json:"mid"`         //广播ID 手动操作时，避免重复添加
	SpacingTime int64  `json:"spacingTime"` // 显示时间  单位：秒
	Content     string `json:"content"`     // 广播的字符内容
	StartTime   int    `json:"startTime"`   // 开始时间
	EndTime     int    `json:"endTime"`     // 结束时间
	Color       string `json:"color"`       // 展示颜色
}

type CreateBroadcastRequest struct {
	SpacingTime int64  `json:"spacingTime"` // 显示时间  单位：秒
	Content     string `json:"content"`     // 广播的字符内容
	StartTime   int    `json:"startTime"`   // 开始时间
	EndTime     int    `json:"endTime"`     // 结束时间
	Color       string `json:"color"`       // 展示颜色
}

type GetBroadcastListRequest struct {
	Page     int64  `form:"page" `
	Limit    int64  `form:"limit"`
	Mid      int    `form:"mid"`
	Position string `form:"position"`
}

// GetBroadcastList 获取所有广播
func GetBroadcastList(c *gin.Context) {
	req := new(GetBroadcastListRequest)
	err := c.ShouldBind(req)
	if err != nil {
		panic(10001)
	}
	if req.Limit == 0 {
		req.Limit = 10
	}
	if req.Page == 0 {
		req.Page = 1
	}
	var filter bson.D

	if req.Mid > 0 {
		filter = append(filter, bson.E{Key: "mid", Value: req.Mid})
	}

	BroadcastModel := new(db.BroadcastModel)

	skip := (req.Page - 1) * req.Limit
	opt := &options.FindOptions{
		Limit: &req.Limit,
		Skip:  &skip,
		Sort:  bson.D{{"mid", -1}},
	}
	Broadcast, _ := BroadcastModel.GetNotice(filter, opt)
	count, err := BroadcastModel.GetTotalNum(filter)

	utils.OkWithData(c, map[string]interface{}{
		"list":  Broadcast,
		"page":  req.Page,
		"limit": req.Limit,
		"total": count,
	})
}

// CreateBroadcast 创建广播
func CreateBroadcast(c *gin.Context) {
	BroadcastRequest := new(CreateBroadcastRequest)
	err := c.ShouldBindJSON(BroadcastRequest)
	if err != nil {
		panic(utils.GetMessage(10001) + ":" + err.Error())
	}
	broadcast := new(db.BroadcastModel)

	if BroadcastRequest.StartTime >= BroadcastRequest.EndTime {
		panic(10001)
	}

	broadcast.EndTime = BroadcastRequest.EndTime
	broadcast.StartTime = BroadcastRequest.StartTime
	broadcast.SpacingTime = BroadcastRequest.SpacingTime
	broadcast.Content = BroadcastRequest.Content
	broadcast.Color = BroadcastRequest.Color

	if err := broadcast.InsertBroadcast(); err != nil {
		utils.Fail(c)
		return
	}
	utils.Ok(c)
}

// UpdateBroadcast 更新广播
func UpdateBroadcast(c *gin.Context) {
	BroadcastRequest := new(UpdateBroadcastRequest)
	err := c.ShouldBindJSON(BroadcastRequest)
	if err != nil {
		panic(10001)
	}
	broadcast := new(db.BroadcastModel)

	broadcast.EndTime = BroadcastRequest.EndTime
	broadcast.StartTime = BroadcastRequest.StartTime
	broadcast.SpacingTime = BroadcastRequest.SpacingTime
	broadcast.Content = BroadcastRequest.Content
	broadcast.Color = BroadcastRequest.Color
	broadcast.Mid = BroadcastRequest.Mid

	if err := broadcast.UpdateBroadcast(); err != nil {
		utils.FailWithMessage(c, utils.GetMessage(10024))
		return
	}
	utils.OkWithMessage(c, utils.GetMessage(10025))
}
