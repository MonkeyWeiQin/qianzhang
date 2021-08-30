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

type GetSystemMailRequest struct {
	Page  int64 `form:"page" `
	Limit int64 `form:"limit"`
}
type CreateSystemMailRequest struct {
	SendTime    int64            `form:"sendTime"`    //发送时间
	InvalidTime int64            `form:"invalidTime"` //失效时间
	CronTime    int64            `form:"cronTime"`    //定时发送
	Uid         []int            `form:"uid"`         //收件人IDs
	ISAllUser   bool             `form:"isAllUser"`   //是否是全服发送
	Title       string           `form:"title"`       //标题
	Content     string           `form:"content"`     //内容
	Attachment  []*db.Attachment `form:"attachment"`  //附件
}
type UpdateSystemMailRequest struct {
	Mid         string           `form:"mid"`         //发送时间
	SendTime    int64            `form:"sendTime"`    //发送时间
	InvalidTime int64            `form:"invalidTime"` //失效时间
	CronTime    int64            `form:"cronTime"`    //定时发送
	Uid         []int            `form:"uid"`         //收件人IDs
	ISAllUser   bool             `form:"isAllUser"`   //是否是全服发送
	Title       string           `form:"title"`       //标题
	Content     string           `form:"content"`     //内容
	Attachment  []*db.Attachment `form:"attachment"`  //附件
}

type AttachmentLabelResponse struct {
	Value ChestType `json:"value"`
	Label string `json:"label"`
}
type GetAttachmentListResponse struct {
	Value string `json:"value"`
	Label string `json:"label"`
}




type ChestType int8 //宝箱奖品类型
const (
	ChestGoldType          ChestType = 1 //金币
	ChestDiamondType       ChestType = 2 //钻石
	ChestStrengthType      ChestType = 3 //体力值
	ChestCardType          ChestType = 4 //卡片 对应CardData
	ChestHeroType          ChestType = 5 //英雄 对应HeroData
	ChestMainWeaponType    ChestType = 6 //主武器 对应weaponData.RoleMainWeaponDataTable
	ChestSubWeaponDataType ChestType = 7 //副武器 对应weaponData.RoleSubWeaponDataTable
	ChestArmorType         ChestType = 8 //装备 对应weaponData.ArmorTable
	ChestOrnamentsType     ChestType = 9 //饰品 对应weaponData.OrnamentsDataTable
)

var ChestLabel = map[ChestType]string{
	ChestGoldType:          "金币",
	ChestDiamondType:       "钻石",
	ChestStrengthType:      "体力值",
	ChestCardType:          "卡片",
	ChestHeroType:          "英雄",
	ChestMainWeaponType:    "主武器",
	ChestSubWeaponDataType: "副武器",
	ChestArmorType:         "装备",
	ChestOrnamentsType:     "饰品",
}

func GetSystemMail(c *gin.Context) {
	req := new(GetSystemMailRequest)
	err := c.ShouldBind(req)
	if err != nil {
		panic(10001)
	}
	var filter bson.D

	SystemMail := new(db.SystemMailModel)

	skip := (req.Page - 1) * req.Limit
	opt := &options.FindOptions{
		Limit: &req.Limit,
		Skip:  &skip,
		Sort:  bson.D{{"mid", -1}},
	}
	mail, _ := SystemMail.GetPlayerMailList(filter, opt)
	count, err := SystemMail.GetTotalNum(filter)

	utils.OkWithData(c, map[string]interface{}{
		"list":  mail,
		"page":  req.Page,
		"limit": req.Limit,
		"total": count,
	})
}

func CreateSystemMail(c *gin.Context) {
	req := new(CreateSystemMailRequest)
	err := c.ShouldBind(req)
	if err != nil {
		utils.FailWithMessage(c, err.Error())
		return
	}
	SystemMail := new(db.SystemMailModel)

	SystemMail.Uid = req.Uid
	SystemMail.SendTime = int(req.SendTime)
	SystemMail.InvalidTime = int(req.InvalidTime)
	SystemMail.CronTime = int(req.CronTime)
	SystemMail.ISAllUser = req.ISAllUser
	SystemMail.Title = req.Title
	SystemMail.Content = req.Content
	SystemMail.Attachment = req.Attachment
	InsertErr := SystemMail.Insert()
	if InsertErr != nil {
		utils.FailWithMessage(c, InsertErr.Error())
		return
	}
	utils.Ok(c)
}
func UpdateSystemMail(c *gin.Context) {
	req := new(UpdateSystemMailRequest)
	err := c.ShouldBind(req)
	if err != nil {
		utils.FailWithMessage(c, err.Error())
		return
	}
	SystemMail := new(db.SystemMailModel)

	update := bson.D{{"sendTime", req.SendTime},
		{"invalidTime", req.InvalidTime},
		{"cronTime", req.CronTime},
		{"uid", req.Uid},
		{"isAllUser", req.ISAllUser},
		{"title", req.Title},
		{"content", req.Content},
		{"attachment", req.Attachment}}
	UpdateErr := SystemMail.UpdatePlayerMail(bson.D{{"mid", req.Mid}}, update, nil)
	if UpdateErr != nil {
		utils.FailWithMessage(c, UpdateErr.Error())
		return
	}
	utils.Ok(c)
}

func GetAttachmentLabel(c *gin.Context) {
	var AttachmentLabel []AttachmentLabelResponse
	for index, chestLabel := range ChestLabel {
		AttachmentLabel = append(AttachmentLabel, AttachmentLabelResponse{
			Value: index,
			Label: chestLabel,
		})
	}
	utils.OkWithData(c, AttachmentLabel)
}

func GetAttachmentList(c *gin.Context) {
	fmt.Println(global.ChestWeaponDataConf)
	fmt.Println(global.CardDataConf)

	var AttachmentLabel = make(map[ChestType][]*GetAttachmentListResponse)
	for index, _ := range ChestLabel {
		if index == ChestGoldType || index == ChestDiamondType || index == ChestStrengthType {
			continue
		}
		switch index {
		case ChestCardType:
			for _, item := range global.CardDataConf {
				AttachmentLabel[index] = append(AttachmentLabel[index], &GetAttachmentListResponse{
					Value: item.Id,
					Label: item.Name,
				})
			}
		case ChestHeroType:
			for _, item := range global.HeroDataConf {
				AttachmentLabel[index] = append(AttachmentLabel[index], &GetAttachmentListResponse{
					Value: item.RelationId,
					Label: item.Name,
				})
			}
		case ChestMainWeaponType:
			for _, item := range global.RoleMainWeaponConf {
				AttachmentLabel[index] = append(AttachmentLabel[index], &GetAttachmentListResponse{
					Value: item.RelationId,
					Label: item.Name,
				})
			}
		case ChestSubWeaponDataType:
			for _, item := range global.RoleSubWeaponConf {
				AttachmentLabel[index] = append(AttachmentLabel[index], &GetAttachmentListResponse{
					Value: item.RelationId,
					Label: item.Name,
				})
			}
		case ChestArmorType:
			for _, item := range global.ArmorDataConf {
				AttachmentLabel[index] = append(AttachmentLabel[index], &GetAttachmentListResponse{
					Value: item.RelationId,
					Label: item.Name,
				})
			}
		case ChestOrnamentsType:
			for _, item := range global.OrnamentsDataConf {
				AttachmentLabel[index] = append(AttachmentLabel[index], &GetAttachmentListResponse{
					Value: item.RelationId,
					Label: item.Name,
				})
			}
		}
	}
	utils.OkWithData(c, AttachmentLabel)
}
