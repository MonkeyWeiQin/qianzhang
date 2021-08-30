package model

import (
	"battle_rabbit/define"
	"battle_rabbit/service/log"
	"battle_rabbit/service/mgoDB"
	"battle_rabbit/utils/xid"
	jsoniter "github.com/json-iterator/go"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	GeneralStage   = 1 //普通关卡
	DifficultStage = 2 //困难关卡

)

var stageColl = newStageColl()

type StageCollection struct {
	*mgoDB.DbBase
}

func newStageColl() *StageCollection {
	return &StageCollection{mgoDB.NewDbBase(define.MgoDBNameBattle, define.TableNameStage)}
}

func GetStageColl() *StageCollection {
	return stageColl
}

/*
普通关卡
困难关卡
活动关卡 // 时间 次数
*/
//
//func InitStageConfigToRedis()  {
//	s,err := redis.String(redisDB.Client.GET(define.RedisActiveStageConf))
//	if err != nil {
//		if err == redis.ErrNil {
//
//		}
//	}
//}
//
//func newActiveStageConfig()  {
//	m := map[string]interface{}{
//		"ChallengeStage":1,
//		"DefenseStage":1,
//
//	}
//}
func GetActiveStageConf() map[string]interface{} {
	m := map[string]interface{}{
		"challengeStage": true,                                                                       // 挑战
		"defenseStage":   true,                                                                       // 营地防守
		"goldStage":      true,                                                                       // 金币
		"towerStage":     true,                                                                       // 塔防
		"resourceStage":  []string{"task001", "task002", "task003", "task004", "task005", "task006"}, // 资源
		"endlessStage":   true,                                                                       // 无尽
	}
	return m
}

// Stage 关卡
type StageModel struct {
	Id             string        `json:"id"  bson:"_id"`
	Uid            int           `json:"uid" bson:"uid"`
	StageMain      *StageMain    `json:"stageMain" bson:"stageMain"`
	ChallengeStage *ActiveStage  `json:"challengeStage" bson:"challengeStage"` // 挑战
	DefenseStage   *ActiveStage  `json:"defenseStage"   bson:"defenseStage"`   // 营地防守
	GoldStage      *ActiveStage  `json:"goldStage"      bson:"goldStage"`      // 金币
	ResourceStage  *ActiveStage  `json:"resourceStage"  bson:"resourceStage"`  // 资源
	EndlessStage   *EndlessStage `json:"endlessStage"   bson:"endlessStage"`   // 无尽
	TowerStage     *ActiveStage  `json:"towerStage" bson:"towerStage"`         // 塔防
}

// 主线关卡(简单和困难)
type StageMain struct {
	GeneralStage    int               `json:"generalStage"    bson:"generalStage"`   // 普通主线任务
	DifficultStage  int               `json:"difficultStage"  bson:"difficultStage"` // 困难主线任务
	StageType       int8              `json:"stageType"       bson:"stageType"`      // 当前选择的关卡模式  0:简单 1:困难
	RewardGeneral   map[string]string `json:"rewardGeneral"   bson:"rewardGeneral"`
	RewardDifficult map[string]string `json:"rewardDifficult" bson:"rewardDifficult"`
}

// 只限制挑战次数的关卡通用
/*
金币关卡：
	1.每天可以进行2次。
	2.获得击杀怪物掉落的金币，怪物只掉落金币。
	3.游戏结算通知服务端获得的金币数量。

挑战关卡	（可以携带英雄协同战斗）
	1.怪物总波次100
	2.每挑战胜利一波，提供一次奖励  （TODO 通知服务端按表给与奖励）
	3.每天可以有一次免费挑战的机会，可消耗钻石购买挑战次数（TODO 通知服务端消耗钻石购买次数）
营地防守关卡:
	1.每天最多一次
	2.完成防守获得奖励
资源关卡: (可以携带英雄协同战斗)
	1.根据一个星期来计算  周一开经验关卡 周二升星关卡 周三升品关卡 周四升级关卡 周五技能卡关卡 周六装备卡关卡 周日 同时开启上述六种关卡。
	2.关卡分五种难度，玩家自己选择难度挑战。
	3.玩家每天可以进行两次资源关卡的战斗，次数到达上限之后不允许进入。
	4.每一个关卡分成三个小关卡，每一个小关卡怪物击杀完成之后都会提供奖励，如果在那一波战斗中失败，则不获得奖励。
*/

type ActiveStage struct {
	Time int `json:"time"` // 时间(当前时间与这个时间相差一天,这个次数重置)
	Num  int `json:"num"`  // 次数(达到次数后不可进入)
}

/*
无尽斗技场关卡	（不可携带英雄协同战斗）
	1.只能主角单人进行
	2.每波怪消灭成功之后恢复20%血量，等待30秒进行下一场战斗
	3.怪物的实力会随着波数的增加而提高
	4.击杀的怪物数量累计月底参与击杀数排行
	5.每天1次斗技场挑战
	6.记录上一次失败的怪物波数，接着打
	7.每个月底结算一次
	8。每月1号数据重置 从第一波开始战斗
*/

type EndlessStage struct {
	Time        int    `json:"time"       bson:"time"`         // 时间(当前时间与这个时间相差一天,这个次数重置)
	Num         int    `json:"num"        bson:"num"`          // 次数(达到次数后不可进入)
	Monster     int    `json:"monster"    bson:"monster"`      // 本月怪物累计
	LastStageId string `json:"lastStageId" bson:"lastStageId"` // 上一次无尽模式通关的关卡ID
}

func NewStageModel(uid int) *StageModel {
	stage := &StageModel{}
	stage.Id = xid.New().String()
	stage.StageMain = &StageMain{
		GeneralStage:    100001,
		DifficultStage:  100001,
		StageType:       define.GeneralStageType,
		RewardGeneral:   make(map[string]string),
		RewardDifficult: make(map[string]string),
	}
	stage.ChallengeStage = &ActiveStage{}
	stage.DefenseStage = &ActiveStage{}
	stage.GoldStage = &ActiveStage{}
	stage.ResourceStage = &ActiveStage{}
	stage.EndlessStage = &EndlessStage{
		LastStageId: "arena0001",
	}
	stage.Uid = uid
	return stage
}

func (s *StageModel) CanActiveStageUseNumber(uid, ty int) bool {
	update := bson.M{}
	switch ty {
	case define.ActiveStageTypeChallenge:
		if s.ChallengeStage.Num > 0 {
			s.ChallengeStage.Num--
			update = bson.M{"challengeStage.num": s.ChallengeStage.Num, "challengeStage.time": s.ChallengeStage.Time}
		}
	case define.ActiveStageTypeDefense:
		if s.DefenseStage.Num > 0 {
			s.DefenseStage.Num--
			update = bson.M{"defenseStage.num": s.DefenseStage.Num, "defenseStage.time": s.DefenseStage.Time}
		}
	case define.ActiveStageTypeGold:
		if s.GoldStage.Num > 0 {
			s.GoldStage.Num--
			update = bson.M{"goldStage.num": s.GoldStage.Num, "goldStage.time": s.GoldStage.Time}
		}
	case define.ActiveStageTypeResource:
		if s.ResourceStage.Num > 0 {
			s.ResourceStage.Num--
			update = bson.M{"resourceStage.num": s.ResourceStage.Num, "resourceStage.time": s.ResourceStage.Time}
		}
	case define.ActiveStageTypeTower:
		if s.TowerStage.Num > 0 {
			s.TowerStage.Num--
			update = bson.M{"towerStage.num": s.TowerStage.Num, "towerStage.time": s.TowerStage.Time}
		}
	case define.ActiveStageTypeEndless: // 无尽
		if s.EndlessStage.Num > 0 {
			s.EndlessStage.Num--
			update = bson.M{"endlessStage.num": s.EndlessStage.Num, "endlessStage.time": s.EndlessStage.Time}
		}
	default:
		log.Debug("没有这个关卡类型!!")
		return false
	}
	if len(update) == 0 {
		b, _ := jsoniter.Marshal(s)
		log.Debug("次数不够了: %s", string(b))
		return false
	}
	_, err := GetStageColl().UpdateOne(nil, bson.M{"uid": uid}, update)
	if err != nil {
		log.Error("查询数据出错!! : ", err)
		return false
	}
	return true
}

func (s *StageCollection) LoadStageFromDB(uid int) (*StageModel, error) {
	stage := new(StageModel)
	err := mgoDB.GetMgoSecondary(define.MgoDBNameBattle).GetCol(s.CollName).FindOne(nil, bson.M{"uid": uid}).Decode(stage)
	return stage, err
}
