package model

import (
)

var (
	skinColl = new(SkinCollection)
)

type SkinCollection struct{}
/*
func GetSkinCollection() *SkinCollection         { return skinColl }
func (m *SkinCollection) CollectionName() string { return TableNameUser }

// SkinModel 玩家皮肤
type SkinModel struct {
	Id     primitive.ObjectID `bson:"_id"`
	Uid    int                `json:"uid"         bson:"uid"`   // 用户的ID
	SkinId string             `json:"skinId"     bson:"skinId"` // 当前皮肤的ID
	Use    bool               `json:"use"         bson:"use"`   // 是否已装备
}

func (m *SkinCollection) GetSkinList(filter bson.D, opt *options.FindOptions) (data []*SkinModel, err error) {
	finder := mgoDB.NewFinder(m).Where(filter).Options(opt).Records(&data)
	err = mgoDB.GetMgo().FindMany(context.TODO(), finder)
	return
}

func (m *SkinCollection) Create(skin *SkinModel) error {
	skin.Id = primitive.NewObjectID()
	err := mgoDB.GetMgo().InsertOne(nil, m, skin)
	return err
}

func (m *SkinCollection) ChangeSkin(uid int, id string) error {
	updater := mgoDB.NewUpdater(m).Where(bson.D{{"uid", uid}}).Update(bson.D{{"use", false}})
	_, err := mgoDB.GetMgo().UpdateMany(nil, updater)
	if err != nil {
		return err
	}
	_, err = mgoDB.GetMgo().UpdateOne(nil, mgoDB.NewUpdater(m).Where(bson.D{{"skin_id", id}}).Update(bson.D{{"use", true}}))
	return err
}
*/