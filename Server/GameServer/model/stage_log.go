package model

import (
	"battle_rabbit/define"
	"battle_rabbit/global"
	"battle_rabbit/service/mgoDB"
)



var (
	stageLogColl = newStageLogColl()
	StageMap     = map[int]string{
		GeneralStage:   "generalStage",
		DifficultStage: "difficultStage",
	}
)

type PlayerStageLogCollection struct{ *mgoDB.DbBase }

func newStageLogColl() *PlayerStageLogCollection {
	return &PlayerStageLogCollection{mgoDB.NewDbBase(define.MgoDBNameBattle, define.TableNameStageLog)}
}
func GetStageLogCollection() *PlayerStageLogCollection { return stageLogColl }

type PlayerStageLogModel struct {
	Uid        int            `json:"uid"         bson:"uid"`        //用户
	StageId    int            `json:"stageId"     bson:"stageId"`    //关卡ID
	StagePass  bool           `json:"stagePass"   bson:"stagePass"`  //0未通过  1通过
	CreateTime int64          `json:"createTime"  bson:"createTime"` //创建时间
	StageType  int            `json:"stageType" bson:"stageType"`    // 关卡类型
	Items      []*global.Item `json:"items" bson:"items"`            // 获得的奖励物品
}

func (m *PlayerStageLogCollection) Insert(playerStageLogModel *PlayerStageLogModel) error {
	err := m.InsertOne(nil, playerStageLogModel)
	return err
}

// TODO 7.0
//func (m *PlayerStageLogCollection) GetList(userId int, day string, stageType int, options *options.FindOptions) (data []PlayerStageLogModel, err error) {
//	data = []PlayerStageLogModel{}
//	finder := mgoDB.NewFinder(m).Options(options).Records(&data)
//	if userId > 0 {
//		finder.Where(bson.D{{"uid", userId}})
//	}
//	if stageType > 0 {
//		finder.Where(bson.D{{"stageType", stageType}})
//	}
//	if day != "" {
//		parse, _ := time.ParseInLocation("2006-01-02 15:04:05", day+" 00:00:00", time.Local)
//		finder.Where(bson.D{{"createTime", bson.M{"$gt": parse.Unix()}}})
//		finder.Where(bson.D{{"createTime", bson.M{"$lt": parse.AddDate(0, 0, 1).Unix() - 1}}})
//	}
//	err = mgoDB.GetMgo().FindMany(context.TODO(), finder)
//	if err != nil {
//		return nil, err
//	}
//	return
//}
// TODO 7.0
//func (m *PlayerStageLogCollection) GetCount(userId int, day string, stageType int) (int64, error) {
//	finder := mgoDB.NewCounter(m)
//	if userId > 0 {
//		finder.Where(bson.D{{"uid", userId}})
//	}
//	if stageType > 0 {
//		finder.Where(bson.D{{"stageType", stageType}})
//	}
//	if day != "" {
//		parse, _ := time.ParseInLocation("2006-01-02 15:04:05", day+" 00:00:00", time.Local)
//		finder.Where(bson.D{{"createTime", bson.M{"$gt": parse.Unix()}}})
//		finder.Where(bson.D{{"createTime", bson.M{"$lt": parse.AddDate(0, 0, 1).Unix() - 1}}})
//	}
//	count, err := mgoDB.GetMgo().CountDocuments(nil, finder)
//	if err != nil {
//		return 0, err
//	}
//	return count, nil
//}
