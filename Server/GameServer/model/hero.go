package model

import (
	"battle_rabbit/define"
	"battle_rabbit/service/mgoDB"
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	DefeatHeroRelationId = "hero_01_1" // 关联ID
)

var (
	heroColl = newHeroCollection()
)

type HeroCollection struct {
	*mgoDB.DbBase
}

func newHeroCollection() *HeroCollection {
	return &HeroCollection{mgoDB.NewDbBase(define.MgoDBNameBattle, define.TableNameHero)}
}

func GetHeroCollection() *HeroCollection { return heroColl }

// HeroModel 英雄的属性结构
type HeroModel struct {
	Id          string   `json:"mid"              bson:"_id"       redis:"id"`
	Uid         int      `json:"-"               bson:"uid"`         //英雄所属用户的ID
	RelationId  string   `json:"relationId"      bson:"relationId"`  //数据表关联ID
	Use         bool     `json:"use"             bson:"use"`         //出战状态 true 为出战，false 不出战
	Rarity      int8     `json:"rarity"          bson:"rarity"`      //稀有度  1:N 2:R 3:SR 4:SSR
	Lv          int      `json:"level"           bson:"level"`       //等级
	StarLv      int      `json:"starLv"          bson:"starLv"`      //星级
	TotalPower  int      `json:"totalPower"      bson:"totalPower"`  //总战力
	Index       int      `json:"index"           bson:"index"`       // 排序索引,前端用
	ASId        string   `json:"asId"            bson:"asId"`        // 英雄主动技能ID ,技能可以升级
	PsIds       []string `json:"psIds"           bson:"psIds"`       // 英雄被动技能ID集合[1,2,3,4] ,技能可以升级
	BuffKill    string   `json:"buffKill"        bson:"buffKill"`    // buff 技能关联ID
	FetterSkill string   `json:"fetterSkill"     bson:"fetterSkill"` // 羁绊技能
	Strengthen  int8     `json:"strengthen"      bson:"strengthen"`  // 强化次数
	EquipmentId string   `json:"equipmentId"     bson:"equipmentId"` // 英雄装备ID

}

type MidResult struct {
	Mid string `json:"mid"          bson:"mid"` //英雄ID
}

func (m *HeroCollection) Upgrade(ctx context.Context, uid, gold, diamond int, heroModel *HeroModel, packageId string, packNum int) (err error) {
	db := mgoDB.GetMgo(define.MgoDBNameBattle)
	if packageId != "" {
		_, err = m.UpdateOne(ctx, bson.M{"_id": packageId}, bson.M{"num": packNum})
		if err != nil {
			return err
		}
	}
	if gold != -1 || diamond != -1 {
		_, err = userColl.UpdateOne(ctx, bson.M{"uid": uid}, bson.M{"gold": gold, "diamond": diamond})
		if err != nil {
			return err
		}
	}
	_, err = db.GetCol(m.CollName).UpdateByID(ctx, heroModel.Id, bson.M{"$set": heroModel})
	return err
}

func (m *HeroCollection) LoadHeroFromDB(uid int) (map[string]*HeroModel, error) {
	cursor, err := mgoDB.GetMgoSecondary(define.MgoDBNameBattle).GetCol(m.CollName).Find(nil, bson.M{"uid": uid})
	if err != nil {
		return nil, err
	}

	var heroList []*HeroModel
	if err = cursor.All(nil, &heroList); err != nil {
		return nil, err
	}
	heroMap := make(map[string]*HeroModel)
	for _, model := range heroList {
		heroMap[model.Id] = model
	}
	return heroMap, nil
}
