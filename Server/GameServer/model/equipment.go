package model

import (
	"battle_rabbit/define"
	"battle_rabbit/service/mgoDB"
	"battle_rabbit/utils/xid"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	DefeatWeaponMainConfigId = "weapon_14_1" // 关联ID
	DefeatWeaponSubConfigId  = "sub_01"
)

var (
	equipmentColl = newEquipmentColl()
)

type EquipmentCollection struct {
	*mgoDB.DbBase
}

func newEquipmentColl() *EquipmentCollection {
	return &EquipmentCollection{mgoDB.NewDbBase(define.MgoDBNameBattle, define.TableNameEquipment)}
}
func GetEquipmentCollection() *EquipmentCollection { return equipmentColl }

// 装备数据结构
type EquipmentModel struct {
	Id         string              `json:"mid"        bson:"_id"`
	Type       define.EquipageType `json:"type"       bson:"type"`        // 类型
	Uid        int                 `json:"uid"        bson:"uid"`         // 用户的ID
	RelationId string              `json:"relationId" bson:"relationId"`  // 数据表关联ID
	Index      int                 `json:"index"      bson:"index"`       // 排序索引
	Level      int                 `json:"level"      bson:"level"`       // 当前等级
	Quality    int                 `json:"quality"    bson:"quality"`     // 装备品质星级
	GemList    []string            `json:"gemList"    bson:"gemList"`     // 镶嵌宝石的IDs(主武器特有)
	Use        bool                `json:"use"        bson:"use"`         // 是否已装备
	Strengthen bool                `json:"strengthen"  bson:"strengthen"` // 是否强化(副武器特有)
}

func NewEquipmentModel(uid int, ty define.EquipageType) *EquipmentModel {
	return &EquipmentModel{
		Id:      xid.New().String(),
		Type:    ty,
		Uid:     uid,
		Level:   1,
		Quality: 1,
	}
}

func (e *EquipmentCollection) ChangeEquipment(oldId, newId string) error {
	var update = []mongo.WriteModel{
		mongo.NewUpdateOneModel().SetFilter(bson.M{"_id": newId}).SetUpdate(bson.M{"$set": bson.M{"use": true}}),
	}

	if oldId != "" {
		update = append(update, mongo.NewUpdateOneModel().SetFilter(bson.M{"_id": oldId}).SetUpdate(bson.M{"$set": bson.M{"use": false}}) )
	}

	_, err := mgoDB.GetMgo(define.MgoDBNameBattle).GetCol(e.CollName).BulkWrite(nil, update)
	if err != nil {
		return err
	}
	return nil
}

func (e *EquipmentCollection) CancelEquipment(id string) error {
	_, err := e.UpdateOne(nil, bson.M{"mid": id}, bson.D{{"use", false}})
	return err
}

//CreateEquipment 添加武器
func (m *EquipmentCollection) CreateEquipment(uid int, relationId string, ty define.EquipageType, index int) *EquipmentModel {
	return &EquipmentModel{
		Id:         xid.New().String(),
		Uid:        uid,
		Type:       ty,
		RelationId: relationId,
		Level:      1,
		Quality:    1,
		Index:      index,
		GemList:    nil,
		Use:        false,
	}

}

func (m *EquipmentCollection) LoadEquipmentFromDB(uid int) (map[string]*EquipmentModel, error) {
	cursor, err := mgoDB.GetMgo(define.MgoDBNameBattle).GetCol(m.CollectionName()).Find(nil, bson.M{"uid": uid})
	if err != nil {
		return nil, err
	}
	var equipList []*EquipmentModel
	err = cursor.All(nil, &equipList)
	if err != nil {
		return nil, err
	}
	equip := make(map[string]*EquipmentModel)
	for _, model := range equipList {
		equip[model.Id] = model
	}

	if len(equip) == 0 {
		return nil, fmt.Errorf("装备为空! uid:%d", uid)
	}
	return equip, nil
}
