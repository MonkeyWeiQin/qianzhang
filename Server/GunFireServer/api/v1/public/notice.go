package public

import (
	"com.xv.admin.server/model/db"
	"com.xv.admin.server/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UpdateNoticeMsgRequest 游戏系统消息
type UpdateNoticeMsgRequest struct {
	Mid       int    `form:"mid"  binding:"required"`       // 公告ID 手动操作时，避免重复添加
	Title     string `form:"title"   binding:"required"`    // 公告名称
	Content   string `form:"content" binding:"required"`    // 公告的字符内容
	Type      int    `form:"type"  binding:"required"`      // 公告类型  //1登录公告  2/游戏内公告
	StartTime int    `form:"startTime"  binding:"required"` // 开始时间
	EndTime   int    `form:"endTime"  binding:"required"`   // 结束时间
}

type CreateNoticeMsgRequest struct {
	Title     string `form:"title"   binding:"required"`    // 公告名称
	Content   string `form:"content" binding:"required"`    // 公告的字符内容
	Type      int    `form:"type"  binding:"required"`      // 公告类型  //1登录公告  2/游戏内公告
	StartTime int    `form:"startTime"  binding:"required"` // 开始时间
	EndTime   int    `form:"endTime"  binding:"required"`   // 结束时间
}

type GetNoticeListRequest struct {
	Page  int64 `form:"page" `
	Limit int64 `form:"limit"`
	Mid   int   `form:"mid"`
	Type  int   `form:"type"`
}

// GetAllNotice 获取所有公告
func GetAllNotice(c *gin.Context) {
	req := new(GetNoticeListRequest)
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
	if req.Type > 0 {
		filter = append(filter, bson.E{Key: "type", Value: req.Type})
	}
	fmt.Println(req.Type)
	NoticeMail := new(db.NoticeMsgModel)

	skip := (req.Page - 1) * req.Limit
	opt := &options.FindOptions{
		Limit: &req.Limit,
		Skip:  &skip,
		Sort:  bson.D{{"mid", -1}},
	}
	mail, _ := NoticeMail.GetNotice(filter, opt)
	count, err := NoticeMail.GetTotalNum(filter)

	utils.OkWithData(c, map[string]interface{}{
		"list":  mail,
		"page":  req.Page,
		"limit": req.Limit,
		"total": count,
	})
}

// CreateNotice 创建公告
func CreateNotice(c *gin.Context) {
	NoticeMsgRequest := new(CreateNoticeMsgRequest)
	err := c.ShouldBind(NoticeMsgRequest)
	if err != nil {
		panic(utils.GetMessage(10001) + ":" + err.Error())
	}
	notice := new(db.NoticeMsgModel)

	if NoticeMsgRequest.StartTime >= NoticeMsgRequest.EndTime {
		panic(10001)
	}

	notice.EndTime = NoticeMsgRequest.EndTime
	notice.StartTime = NoticeMsgRequest.StartTime
	notice.Title = NoticeMsgRequest.Title
	notice.Content = NoticeMsgRequest.Content
	notice.Type = NoticeMsgRequest.Type

	if err := notice.InsertNotice(); err != nil {
		utils.Fail(c)
	}
	utils.Ok(c)
}

// UpdateNotice 更新公告
func UpdateNotice(c *gin.Context) {
	NoticeMsgRequest := new(UpdateNoticeMsgRequest)
	err := c.ShouldBind(NoticeMsgRequest)
	if err != nil {
		panic(10001)
	}
	notice := new(db.NoticeMsgModel)

	notice.EndTime = NoticeMsgRequest.EndTime
	notice.StartTime = NoticeMsgRequest.StartTime
	notice.Title = NoticeMsgRequest.Title
	notice.Content = NoticeMsgRequest.Content
	notice.Type = NoticeMsgRequest.Type

	if err := notice.UpdateNotice(); err != nil {
		utils.FailWithMessage(c, utils.GetMessage(10024))
	}
	utils.OkWithMessage(c, utils.GetMessage(10025))
}
