package request

import "battle_rabbit/global"

//  关卡结算
type SettlementStageReq struct {
	StageId   interface{}    `json:"stageId"`   // 关卡ID
	StageType int            `json:"stageType"` // 关卡类型(枚举值)
	StagePass bool           `json:"stagePass"` // 0未通过  1通过
	Items     []*global.Item `json:"items"`     // 本关卡获得得物品
	KillPeople int  	      `json:"killPeople"` // 无尽关卡杀死敌人的数量
}

// DecPlayerStrengthRequest 扣除玩家体力
type DecPlayerStrengthRequest struct {
	StageType     interface{}    `json:"stageType"` // 关卡类型
	StageId       string `json:"stageId"`       // 关卡ID
}

type StartPlayStageReq struct {
	StageType int         `json:"stageType"` // 关卡类型
	StageId   interface{} `json:"stageId"`   // 关卡ID
}

type PlayerJoinStageByTypeReq struct {
	Ty int `json:"ty"`
}

type ReceiveChapterRewardReq struct {
	Type  int   `json:"type"`
	MapId string `json:"mapId"`
}