package public

import (
	"com.xv.admin.server/model/db"
	"com.xv.admin.server/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type GetPlayerMailListRequest struct {
	Page  int64 `form:"page" `
	Limit int64 `form:"limit"`
	Uid   int   `form:"uid"`
}

type SetPlayerMailStateRequest struct {
	Id  string `json:"mid"  binding:"required"`
	Uid int    `json:"uid"  binding:"required"`
}

// CreatePlayerMail 增加邮件
func CreatePlayerMail(c *gin.Context) {
	new(db.PlayerMailModel).SendNewPlayerMail(00000)
}

// GetPlayerMailList 获得玩家的邮件列表
func GetPlayerMailList(c *gin.Context) {
	req := new(GetPlayerMailListRequest)
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

	if req.Uid > 0 {
		filter = append(filter, bson.E{Key: "uid", Value: req.Uid})
	}

	playerMail := new(db.PlayerMailModel)

	skip := (req.Page - 1) * req.Limit
	opt := &options.FindOptions{
		Limit: &req.Limit,
		Skip:  &skip,
		Sort:  bson.D{{"sendTime", -1}},
	}
	mail, GetPlayerMailListErr := playerMail.GetPlayerMailList(filter, opt)
	count, err := playerMail.GetTotalNum(filter)
	fmt.Println(GetPlayerMailListErr)
	utils.OkWithData(c, map[string]interface{}{
		"list":  mail,
		"page":  req.Page,
		"limit": req.Limit,
		"total": count,
	})
}

// SetPlayerMailRead 设置邮件的阅读
func SetPlayerMailRead(c *gin.Context) {
	SetPlayerMailStateRequest := new(SetPlayerMailStateRequest)
	err := c.ShouldBind(SetPlayerMailStateRequest)
	if err != nil {
		panic(10001)
	}

	id := SetPlayerMailStateRequest.Id
	userId:=SetPlayerMailStateRequest.Uid

	playerMail := new(db.PlayerMailModel)

	if err := playerMail.SetMailRead(id, userId, nil); err != nil {
		panic(10024)
	}
	utils.OkWithMessage(c, utils.GetMessage(10025))
}

// DelPlayerMail 删除邮件
func DelPlayerMail(c *gin.Context) {
	SetPlayerMailStateRequest := new(SetPlayerMailStateRequest)
	err := c.ShouldBind(SetPlayerMailStateRequest)
	if err != nil {
		panic(10001)
	}
	id := SetPlayerMailStateRequest.Id
	userId:=SetPlayerMailStateRequest.Uid
	playerMail := new(db.PlayerMailModel)
	if err := playerMail.SetMailDel(id, userId, nil); err != nil {
		panic(10026)
	}
	utils.OkWithMessage(c, utils.GetMessage(10027))
}
