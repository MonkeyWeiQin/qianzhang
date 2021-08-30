package user

import (
	"com.xv.admin.server/model/db"
	"com.xv.admin.server/model/request/user_req"
	"com.xv.admin.server/service/redisDB"
	"com.xv.admin.server/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
	"strings"
)

func GetUsersList(c *gin.Context) {
	req := new(user_req.GetUsersListRequest)
	err := c.ShouldBind(req)
	if err != nil {
		fmt.Println(err)
		utils.FailWithMessage(c, "参数错误"+err.Error())
		return
	}

	if req.Limit == 0 {
		req.Limit = 10
	}
	if req.Page == 0 {
		req.Page = 1
	}
	model := &db.UserModel{}
	var filter bson.D
	if req.Uid != "" {
		var uidArr []int
		for _, item := range strings.Split(req.Uid, ",") {
			if uid, AtoiErr := strconv.Atoi(item); AtoiErr == nil {
				uidArr = append(uidArr, uid)
			}
		}
		filter = append(filter, bson.E{Key: "uid", Value: bson.M{"$in": uidArr}})
	}
	if req.Username != "" {
		filter = append(filter, bson.E{Key: "username", Value: req.Username})
	}
	if req.RegisterStartTime > 0 {
		filter = append(filter, bson.E{Key: "registerTime", Value: bson.M{"$gte": req.RegisterStartTime}})
	}
	if req.RegisterEndTime > 0 {
		filter = append(filter, bson.E{Key: "registerTime", Value: bson.M{"$lte": req.RegisterStartTime}})
	}
	skip := (req.Page - 1) * req.Limit
	opt := &options.FindOptions{
		Limit: &req.Limit,
		Skip:  &skip,
		Sort:  bson.D{{"uid", -1}},
	}
	datas, err := model.GetUsersList(filter, opt)
	if err != nil {
		utils.FailWithMessage(c, "参数错误"+err.Error())
		return
	}
	count, err := model.GetTotalNum(filter)
	utils.OkWithData(c, map[string]interface{}{
		"list":  datas,
		"page":  req.Page,
		"limit": req.Limit,
		"total": count,
	})
}
func ModifyDiamond(c *gin.Context) {
	req := new(user_req.ModifyNumberRequest)
	err := c.ShouldBindJSON(req)
	if err != nil {
		utils.FailWithMessage(c, "参数错误"+err.Error())
		return
	}
	model := &db.UserModel{}
	number := 0
	if req.Type == "add" {
		number = int(req.Number)
	} else {
		number = int(-req.Number)
	}
	if err := model.Inc(bson.D{{"uid", req.Uid}}, bson.D{{"diamond", number}}); err != nil {
		utils.FailWithMessage(c, "更新失败")
		return
	}

	res, err := model.GetUser(bson.D{{"uid", req.Uid}})
	if err != nil {
		utils.FailWithMessage(c, "更新失败")
		return
	}
	_, RedisErr := redisDB.Client.Hmset(db.DefaultUserModelRedisKey(req.Uid), "Diamond", res.Diamond)
	if RedisErr != nil {
		fmt.Println("刷新缓存失败：", RedisErr)
		utils.FailWithMessage(c, "更新失败")
		return
	}
	AdminOperationRecordModel := &db.AdminOperationRecordModel{}
	AdminOperationRecordModel.BusinessType = 1
	AdminOperationRecordModel.BusinessStyle = 1
	AdminOperationRecordModel.Uid = req.Uid
	AdminOperationRecordModel.Description = fmt.Sprintf("修改用户钻石数据:  %d", number)
	_ = AdminOperationRecordModel.Create()

	utils.OkWithMessage(c, "更新成功")
}
func ModifyGold(c *gin.Context) {
	req := new(user_req.ModifyNumberRequest)
	err := c.ShouldBindJSON(req)
	if err != nil {
		utils.FailWithMessage(c, "参数错误"+err.Error())
		return
	}
	model := &db.UserModel{}
	number := 0
	if req.Type == "add" {
		number = int(req.Number)
	} else {
		number = int(-req.Number)
	}
	if err := model.Inc(bson.D{{"uid", req.Uid}}, bson.D{{"gold", number}}); err != nil {
		utils.FailWithMessage(c, "更新失败")
		return
	}
	res, err := model.GetUser(bson.D{{"uid", req.Uid}})
	if err != nil {
		utils.FailWithMessage(c, "更新失败")
		return
	}
	_, RedisErr := redisDB.Client.Hmset(db.DefaultUserModelRedisKey(req.Uid), "Gold", res.Gold)
	if RedisErr != nil {
		fmt.Println("刷新缓存失败：", RedisErr)
		utils.FailWithMessage(c, "更新失败")
		return
	}
	AdminOperationRecordModel := &db.AdminOperationRecordModel{}
	AdminOperationRecordModel.BusinessType = 1
	AdminOperationRecordModel.BusinessStyle = 2
	AdminOperationRecordModel.Uid = req.Uid
	AdminOperationRecordModel.Description = fmt.Sprintf("修改用户金币数据:  %d", number)
	_ = AdminOperationRecordModel.Create()

	utils.OkWithMessage(c, "更新成功")
}
func ModifyStrength(c *gin.Context) {
	req := new(user_req.ModifyNumberRequest)
	err := c.ShouldBindJSON(req)
	if err != nil {
		utils.FailWithMessage(c, "参数错误"+err.Error())
		return
	}
	model := &db.UserModel{}
	number := 0
	if req.Type == "add" {
		number = int(req.Number)
	} else {
		number = int(-req.Number)
	}
	if err := model.Inc(bson.D{{"uid", req.Uid}}, bson.D{{"strength", number}}); err != nil {
		utils.FailWithMessage(c, "更新失败")
		return
	}
	res, err := model.GetUser(bson.D{{"uid", req.Uid}})
	if err != nil {
		utils.FailWithMessage(c, "更新失败")
		return
	}
	_, RedisErr := redisDB.Client.Hmset(db.DefaultUserModelRedisKey(req.Uid), "Strength", res.Strength)
	if RedisErr != nil {
		fmt.Println("刷新缓存失败：", RedisErr)
		utils.FailWithMessage(c, "更新失败")
		return
	}
	AdminOperationRecordModel := &db.AdminOperationRecordModel{}
	AdminOperationRecordModel.BusinessType = 1
	AdminOperationRecordModel.BusinessStyle = 3
	AdminOperationRecordModel.Uid = req.Uid
	AdminOperationRecordModel.Description = fmt.Sprintf("修改用户体力数据:  %d", number)
	_ = AdminOperationRecordModel.Create()

	utils.OkWithMessage(c, "更新成功")
}
func ModifyStatus(c *gin.Context) {
	req := new(user_req.ModifyStatusRequest)
	err := c.ShouldBindJSON(req)
	if err != nil {
		utils.FailWithMessage(c, "参数错误"+err.Error())
		return
	}
	model := &db.UserModel{}
	update := bson.D{}

	if req.Status != 0 {
		update = append(update, bson.E{Key: "status", Value: req.Time})
	} else {
		update = bson.D{{"status", 0}}
	}

	if err := model.UpdateUser(bson.D{{"uid", req.Uid}}, update, nil); err != nil {
		utils.FailWithMessage(c, "更新失败")
		return
	}

	AdminOperationRecordModel := &db.AdminOperationRecordModel{}
	AdminOperationRecordModel.BusinessType = 1
	AdminOperationRecordModel.BusinessStyle = 4
	AdminOperationRecordModel.Uid = req.Uid
	AdminOperationRecordModel.Description = fmt.Sprintf("修改用户状态数据:  %d", req.Status)
	_ = AdminOperationRecordModel.Create()

	utils.OkWithMessage(c, "更新成功")
}
