package protocol

import (
	"battle_rabbit/define"
	"battle_rabbit/global"
	"battle_rabbit/iface"
)

type LevelUpgradeNoticeData struct {
	UpAcc  bool
	UpRole bool
	Lv     int `json:"lv"`   // 账号等级
	Exp    int `json:"exp"`  // 账号经验
	RLv    int `json:"rLv"`  // 角色等级
	RExp   int `json:"rExp"` // 角色经验
}

// 角色升级变更通知
func NoticeRoleLevelUpgrade(session iface.ISession, data *LevelUpgradeNoticeData) {
	session.Send(SuccessData(define.PushMsgId1101, data))
}

type NoticeUserGoldAndDiamondData struct {
	Gold    int `json:"gold"`
	Diamond int `json:"diamond"`
}

// 金币和钻石变化通知
func NoticeUserGoldAndDiamond(session iface.ISession, gold, Diamond int) {
	session.Send(SuccessData(define.PushMsgId1020, &NoticeUserGoldAndDiamondData{
		Gold:    gold,
		Diamond: Diamond,
	}))
}

// 角色属性更变通知
func NoticeRoleAttribute(session iface.ISession, attribute *global.Attribute) {
	session.Send(SuccessData(define.PushMsgId1021, map[string]*global.Attribute{"attribute": attribute}))
}

type NoticeWeaponAttributeData struct {
	Ty        define.EquipageType `json:"ty"`
	Attribute *global.Attribute   `json:"attribute"`
}

// 武器属性更变通知
func NoticeWeaponAttribute(session iface.ISession, ty define.EquipageType, attribute *global.Attribute) {
	session.Send(SuccessData(define.PushMsgId1022, NoticeWeaponAttributeData{
		Ty:        ty,
		Attribute: attribute,
	}))
}

type NoticeUserStrengthChangeData struct {
	Strength int `json:"strength"`
	Time     int64 `json:"time"`
}

func NoticeUserStrengthChange(session iface.ISession, s int, t int64) {
	session.Send(SuccessData(define.PushMsgId1023, &NoticeUserStrengthChangeData{Strength: s, Time: t}))
}

// 任务完成通知
func NoticeUserTaskComplete(sess iface.ISession, data interface{}) {
	sess.Send(SuccessData(define.PushMsgId2001, data))
}
