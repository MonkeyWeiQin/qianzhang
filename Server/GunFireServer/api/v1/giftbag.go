package v1

import (
	"com.xv.admin.server/model/db"
	"com.xv.admin.server/model/request/players"
	"com.xv.admin.server/utils"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func GetGiftList(c *gin.Context) {
	req := new(players.GetGiftListRequest)
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
	model := &db.GiftCodeModel{}
	var filter bson.D

	skip := (req.Page - 1) * req.Limit
	opt := &options.FindOptions{
		Limit: &req.Limit,
		Skip:  &skip,
		Sort:  bson.D{{"uid", -1}},
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

func CreateGift(c *gin.Context) {
	req := new(players.CreateGiftRequest)
	err := c.ShouldBindJSON(req)
	if err != nil {
		utils.FailWithMessage(c, "参数错误"+err.Error())
		return
	}
	var giftCode = &db.GiftCodeModel{}
	num, GetTotalNumErr := giftCode.GetTotalNum(bson.D{{"name", req.Name}})
	if GetTotalNumErr != nil {
		utils.FailWithMessage(c, "创建礼品码失败"+GetTotalNumErr.Error())
		return
	}
	if num >= 1 {
		utils.FailWithMessage(c, "创建礼品码失败:名称重复")
		return
	}
	giftCode.Name = req.Name
	giftCode.Count = req.Count
	giftCode.Info = req.Info
	giftCode.Attachment = req.Attachment
	giftCode.EffectiveTime = req.EffectiveTime
	if CreateErr := giftCode.Insert(); CreateErr != nil {
		utils.FailWithMessage(c, "创建礼品码失败"+CreateErr.Error())
		return
	}
	count := 10000
	for giftCode.Count > 0 {
		if giftCode.Count < count {
			count = giftCode.Count
		} else {
			giftCode.Count = giftCode.Count - count
		}
		go func() {
			err := generateCode(count, giftCode.Mid, int64(giftCode.Count))
			if err != nil {
				utils.FailWithMessage(c, "创建礼品码失败:生成兑换码失败："+err.Error())
				return
			}
		}()
	}
	utils.Ok(c)
	return
}

func generateCode(number int, mit int,seed int64) error {
	rand.Seed(time.Now().UnixNano() + seed)
	bytes := make([]byte, 8)
	code := make([]interface{}, number)
	for i := 0; i < number; i++ {
		for i := 0; i < 8; i++ {
			b := rand.Intn(26) + 97
			bytes[i] = byte(b)
		}
		code[i] = &db.GiftCodeLogModel{
			Mid:  mit,
			Code: string(bytes),
		}
	}
	err := new(db.GiftCodeLogModel).Insert(code)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func DownLoadGiftCode(c *gin.Context) {
	req := new(players.GiftDownloadRequest)
	err := c.ShouldBind(req)
	if err != nil {
		utils.FailWithMessage(c, "参数错误"+err.Error())
		return
	}
	var giftCode = &db.GiftCodeModel{}
	GiftModel, SelectOneErr := giftCode.SelectOne(bson.D{{"mid",req.Mid}}, options.FindOne())

	if SelectOneErr != nil {
		utils.FailWithMessage(c, "导出失败"+SelectOneErr.Error())
		return
	}

	f := excelize.NewFile()
	// Set value of a cell.
	list, err := new(db.GiftCodeLogModel).GetList(bson.D{{"mid",req.Mid}}, nil)
	if err != nil {
		return
	}
	//writer, err := f.NewStreamWriter("Sheet1")
	if err != nil {
		return
	}
	for index, item := range list {
		err := f.SetCellValue("Sheet1", "A"+strconv.Itoa(index+1), item["code"])
		if err != nil {
			utils.FailWithMessage(c, "导出失败"+err.Error())
			return
		}
	}
	// Set active sheet of the workbook.
	fileName := GiftModel.Name + "_" +strconv.Itoa(int(time.Now().Unix())) + ".csv"
	// Save spreadsheet by the given path.\
	//

	if err := f.SaveAs(fileName); err != nil {
		fmt.Println(err)
	}

	//t := "wwwwwww,qweqw,wqeqw,eqwe, www,ee,rrr,ttt"

	fileContentDisposition := "attachment;filename=" + fileName
	c.Header("Content-Type", "application/octet-stream") // 这里是压缩文件类型 .zip
	c.Header("Content-Disposition", fileContentDisposition)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("response-type", "blob")
	data, _ := f.WriteToBuffer()
	c.Data(http.StatusOK, "application/vnd.ms-excel", data.Bytes())
}

func GetGiftLogList(c *gin.Context){
	req := new(players.GetGiftLogListRequest)
	err := c.ShouldBind(req)
	if err != nil {
		fmt.Println(err)
		utils.FailWithMessage(c, "参数错误"+err.Error())
		return
	}
	var filter bson.D
	if req.Limit == 0 {
		req.Limit = 10
	}
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Used == 1 {
		filter =append(filter,bson.E{Key: "uid",Value: bson.E{Key: "$eq",Value: 0}})
	}
	if req.Used == 2 {
		filter =append(filter,bson.E{Key: "uid",Value: 0})
	}
	if req.Uid !=0  {
		filter =append(filter,bson.E{Key: "uid",Value: req.Uid })
	}
	if req.Mid != 0  {
		filter =append(filter,bson.E{Key: "mid",Value: req.Mid })
	}
	if req.Code != ""  {
		filter =append(filter,bson.E{Key: "code",Value: req.Code })
	}

	model := &db.GiftCodeLogModel{}

	skip := (req.Page - 1) * req.Limit
	opt := &options.FindOptions{
		Limit: &req.Limit,
		Skip:  &skip,
		Sort:  bson.D{{"uid", -1}},
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