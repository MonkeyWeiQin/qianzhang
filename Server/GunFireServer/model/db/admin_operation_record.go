package db

import (
	"com.xv.admin.server/config"
	"com.xv.admin.server/service/mgoDB"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type AdminOperationRecordModel struct {
	Id            primitive.ObjectID `json:"-"        bson:"-"`
	AdminId       int8               `json:"admin_id"  bson:"admin_id"`
	Uid           int               `json:"uid"  bson:"uid"`
	BusinessType  int8               `json:"business_type"  bson:"business_type"`   // 业务大类  1 用户
	BusinessStyle int8               `json:"business_style"  bson:"business_style"` // 业务小类  1 修改体力
	Description   string             `json:"description"  bson:"description"`       // 描述
	CreateTime    int64              `json:"create_time"  bson:"create_time"`       // 创建时间

}

func (c *AdminOperationRecordModel) CollectionName() string {
	return config.CollAdminOperationRecordConfig
}

func (c *AdminOperationRecordModel) GetId() primitive.ObjectID {
	return c.Id
}

func (c *AdminOperationRecordModel) SetId(id primitive.ObjectID) {
	c.Id = id
}

func (c *AdminOperationRecordModel) Create() error {
	c.CreateTime = time.Now().Unix()
	err := mgoDB.GetMgo().InsertOne(nil, c)
	return err
}
