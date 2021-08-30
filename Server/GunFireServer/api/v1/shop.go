package v1

import (
	"com.xv.admin.server/global"
	"com.xv.admin.server/model/db"
	"com.xv.admin.server/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetGoodsPurchaseListRequest 请求参数
type GetGoodsPurchaseListRequest struct {
	GoodsId string `form:"goodsId"`
	Uid     int    `form:"uid"`
	Page    int64  `form:"page"`
	Limit   int64  `form:"limit"`
}
// GetChestListRequest 请求参数
type GetChestListRequest struct {
	ChestId string `form:"chestId"`
	Uid     int    `form:"uid"`
	Page    int64  `form:"page"`
	Limit   int64  `form:"limit"`
}

var (
	chestBoxList = map[string]map[string]*global.ChestListConfig{
		"treasurebox001": global.ChestHeroDataConf,
		"treasurebox002": global.ChestMatDataConf,
		"treasurebox003": global.ChestEquippedDataConf,
		"treasurebox004": global.ChestWeaponDataConf,
	}
)

func GetGoodsList(c *gin.Context) {
	utils.OkWithData(c, global.GoodsDataConf)
	return
}

func GetGoodsPurchaseList(c *gin.Context) {
	req := new(GetGoodsPurchaseListRequest)
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
	var filter bson.D
	if req.Uid > 0 {
		filter = append(filter, bson.E{Key: "uid", Value: req.Uid})
	}
	if req.GoodsId != "" {
		filter = append(filter, bson.E{Key: "goodsId", Value: req.GoodsId})
	}
	model := new(db.PurchaseGoodsLogModel)

	skip := (req.Page - 1) * req.Limit
	opt := &options.FindOptions{
		Limit: &req.Limit,
		Skip:  &skip,
		Sort:  bson.D{{"time", -1}},
	}
	datas, err := model.GetList(filter, opt)
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


func GetChestList(c *gin.Context) {
	utils.OkWithData(c, global.ChestDataConf)
	return
}


func GetChestContent(c *gin.Context) {
	utils.OkWithData(c, chestBoxList)
	return
}



func GetChestPurchaseList(c *gin.Context) {
	req := new(GetChestListRequest)
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
	var filter bson.D
	if req.Uid > 0 {
		filter = append(filter, bson.E{Key: "uid", Value: req.Uid})
	}
	if req.ChestId != "" {
		filter = append(filter, bson.E{Key: "chestId", Value: req.ChestId})
	}
	model := new(db.ChestLogModel)

	skip := (req.Page - 1) * req.Limit
	opt := &options.FindOptions{
		Limit: &req.Limit,
		Skip:  &skip,
		Sort:  bson.D{{"time", -1}},
	}
	datas, err := model.GetList(filter, opt)
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
