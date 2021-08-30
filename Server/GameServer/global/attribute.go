package global

import (
	jsoniter "github.com/json-iterator/go"
)

// 属性和技能等类型定义,全局使用
type Attribute struct {
	LifeA         int     `xlsx:"life_a"          bson:"lifeA"          json:"lifeA"`         //生命值
	LifeC         float32 `xlsx:"life_c"          bson:"lifeC"          json:"lifeC"`         //生命加成
	DefenseA      int     `xlsx:"defense_a"       bson:"defenseA"       json:"defenseA"`      //防御值
	DefenseC      float32 `xlsx:"defense_c"       bson:"defenseC"       json:"defenseC"`      //防御力加成
	AttackA       int     `xlsx:"attack_a"        bson:"attackA"        json:"attackA"`       //攻击值
	AttackC       float32 `xlsx:"attack_c"        bson:"attackC"        json:"attackC"`       //攻击加成
	AttackSpeedA  int     `xlsx:"attack_speed_a"  bson:"attackSpeedA"   json:"attackSpeedA"`  //攻速
	AttackSpeedC  float32 `xlsx:"attack_speed_c"  bson:"attackSpeedC"   json:"attackSpeedC"`  //攻速加成
	MoveSpeedA    int     `xlsx:"move_speed_a"    bson:"moveSpeedA"     json:"moveSpeedA"`    //移速
	MoveSpeedC    float32 `xlsx:"move_speed_c"    bson:"moveSpeedC"     json:"moveSpeedC"`    //移速加成
	DodgeB        int     `xlsx:"dodge_b"         bson:"dodgeB"         json:"dodgeB"`        //闪避几率
	DodgeC        int     `xlsx:"dodge_c"         bson:"dodgeC"         json:"dodgeC"`        //闪避几率加成
	CriticalB     int     `xlsx:"critical_b"      bson:"criticalB"      json:"criticalB"`     //暴击几率
	CriticalC     int     `xlsx:"critical_c"      bson:"criticalC"      json:"criticalC"`     //暴击几率加成
	BuffTimeA     float32 `xlsx:"buff_time_a"     bson:"buffTimeA"      json:"buffTimeA"`     //buff时长
	BuffTimeC     float32 `xlsx:"buff_time_c"     bson:"buffTimeC"      json:"buffTimeC"`     //buff时长加成
	BossHurtA     float32 `xlsx:"boss_hurt_a"     bson:"bossHurtA"      json:"bossHurtA"`     //boss伤害
	BossHurtC     float32 `xlsx:"boss_hurt_c"     bson:"bossHurtC"      json:"bossHurtC"`     //boss伤害加成
	CriticalHurtA float32 `xlsx:"critical_hurt_a" bson:"criticalHurtA"  json:"criticalHurtA"` //暴击伤害
	CriticalHurtC float32 `xlsx:"critical_hurt_c" bson:"criticalHurtC"  json:"criticalHurtC"` //暴击伤害加成
	GoldAddC      float32 `xlsx:"gold_add_c"      bson:"goldAddC"       json:"goldAddC"`      //金币加成
}

func (t *Attribute) RedisScan(src interface{}) error { return jsoniter.Unmarshal(src.([]byte), t) }
func (t *Attribute) RedisArg() interface{}           { d, _ := jsoniter.Marshal(t); return d }
func (t Attribute) Add(attrs ...*Attribute) *Attribute {
	for _, attr := range attrs {
		if attr == nil {
			continue
		}
		t.LifeA += attr.LifeA
		t.LifeC += attr.LifeC
		t.DefenseA += attr.DefenseA
		t.DefenseC += attr.DefenseC
		t.AttackA += attr.AttackA
		t.AttackC += attr.AttackC
		t.AttackSpeedA += attr.AttackSpeedA
		t.AttackSpeedC += attr.AttackSpeedC
		t.MoveSpeedA += attr.MoveSpeedA
		t.MoveSpeedC += attr.MoveSpeedC
		t.DodgeB += attr.DodgeB
		t.DodgeC += attr.DodgeC
		t.CriticalB += attr.CriticalB
		t.CriticalC += attr.CriticalC
		t.BuffTimeA += attr.BuffTimeA
		t.BuffTimeC += attr.BuffTimeC
		t.BossHurtA += attr.BossHurtA
		t.BossHurtC += attr.BossHurtC
		t.CriticalHurtA += attr.CriticalHurtA
		t.CriticalHurtC += attr.CriticalHurtC
		t.GoldAddC += attr.GoldAddC
	}
	return &t
}
func (attr *Attribute) AttributeConst() *Attribute {
	attr.LifeA = (attr.LifeA * int((attr.LifeC+1)*100)) / 100
	attr.DefenseA = (attr.DefenseA * int((attr.DefenseC+1)*100)) / 100
	attr.AttackA = (attr.AttackA * int((attr.AttackC+1)*100)) / 100
	attr.AttackSpeedA = (attr.AttackSpeedA * int((attr.AttackSpeedC+1)*100)) / 100
	attr.MoveSpeedA = (attr.MoveSpeedA * int((attr.MoveSpeedC+1)*100)) / 100
	attr.DodgeB += attr.DodgeC
	attr.CriticalB += attr.CriticalC
	attr.BuffTimeA += attr.BuffTimeC // TODO 不确定是不是相加
	attr.BossHurtA += attr.BossHurtC // TODO 不确定是不是相加
	attr.CriticalHurtA += attr.CriticalHurtC
	return attr
}
