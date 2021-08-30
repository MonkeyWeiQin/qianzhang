package model

/* =============== 英雄装备集合 =======================*/

type HeroEquipment struct {
	Mid    string `json:"mid"    bson:"_id"`
	Level  int    `json:"level"  bson:"level"`
	StarLv int    `json:"starLv" bson:"starLv"`
	
}
