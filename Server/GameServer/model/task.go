package model

import (
	"battle_rabbit/define"
	"battle_rabbit/service/mgoDB"
	"battle_rabbit/utils"
	"battle_rabbit/utils/xid"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	TaskStateBegin    = 0
	TaskStateComplete = 1
	TaskStateOver     = 2

	TaskCacheKeyPrefix = "user_task_"
)

// task配置表转换为使用便捷的map

/*
1 主线任务 | 每日任务 两类
2 主线任务中某一个任务完成后会获得这个子任务的奖励并刷新一个新的子任务
3 每日任务与主线任务不冲突,相同的任务分开统计,每日任务每天清零,若达标的每日任务没有领取奖励,则奖励通过邮件发送

*/

var (
	taskColl = &TaskColl{
		DbBase: &mgoDB.DbBase{
			DbName:   define.MgoDBNameBattle,
			CollName: define.TableNameTask,
		},
	}
)

//TaskColl
type TaskColl struct {
	*mgoDB.DbBase
}

func GetTaskCollection() *TaskColl { return taskColl }

// 玩家存储数据
type (
	TaskModel struct {
		Id          string           `json:"id"           bson:"id"` // 唯一ID
		Uid         int              `json:"uid"          bson:"uid"`
		OverTask    map[int]int      `json:"overTask"     bson:"overTask"`     // 已经完成的任务的ID集合
		MainTask    map[int]*Subtask `json:"mainTaskMap"  bson:"mainTaskMap"`  // 主线任务集合
		DayTask     map[int]*Subtask `json:"everydayTask" bson:"everydayTask"` // 每日任务集合
		DayTaskTime int64            `json:"dayTaskTime"  bson:"dayTaskTime"`  // 每日任务的当天零点日期
	}

	Subtask struct {
		TaskId    int    `json:"taskId"    bson:"taskId"`
		Type      int    `json:"type"      bson:"type"`      // 任务类型
		LimitTime int    `json:"limitTime" bson:"limitTime"` // 时间限制,在生成任务时,配置中不为0,则这个值等于生成任务的时间节点(unix)加上配置文件中的时间,大于这个时间没完成任务则失效
		Progress  int    `json:"progress"  bson:"progress"`  // 任务进度
		Condition string `json:"condition" bson:"condition"` // 附加条件参数 // 若不为空(""),根据类型判别是否是这个条件,比如装备的ID,英雄ID,怪物Boss,天赋技能等
		Target    int    `json:"target"    bson:"target"`    // 任务目标数量,达成任务的条件数量
		State     int    `json:"state"     bson:"state"`     // 任务状态,0 进行中, 1 已完成,2 奖励已领取,3,等待刷新,新的任务,(某些条件没达到,还不能切换到下一个任务,如: 角色等级没达到),若没有写一个任务,则直接结束当前任务
		RewardId  string `json:"rewardId"  bson:"rewardId"`  // 任务奖励
	}

	TaskProgress struct {
		Type      int
		Condition string
		Num       int
		Lv        int
	}
)

func (tm *TaskModel) NewTaskModel(uid int) *TaskModel {
	return &TaskModel{
		Id:          xid.New().String(),
		Uid:         uid,
		MainTask:    make(map[int]*Subtask), // 任务ID => 详细的任务数据
		DayTask:     make(map[int]*Subtask),
		DayTaskTime: utils.GetMidnightTime().Unix(),
	}
}

func (tm *TaskModel) CreateSubTask(isMain bool, taskId, ty, progress, target int) {
	subTask := &Subtask{
		Type:     ty,
		Progress: progress,
		Target:   target,
	}
	if isMain {
		tm.MainTask[ty] = subTask
	} else {
		tm.DayTask[ty] = subTask
	}
}

// 初始化主线和每日任务的任务详细内容
func (tm *TaskModel) InitMainTask() {

}

func (tm *TaskModel) FlushDayTask() {

}

func (*TaskColl)LoadTaskFromDB(uid int) (*TaskModel, error) {
	var tm *TaskModel
	err := taskColl.FindOne(nil, bson.M{"uid": uid}, &tm)
	if err != nil {
		return nil, err
	}
	return tm, err
}

func (tm *TaskModel) Save(main, day bool) error {
	if main == false && day == false {
		return nil
	}

	var updater = bson.M{}
	if main {
		updater["mainTask"] = tm.MainTask
	}

	if day {
		updater["dayTask"] = tm.DayTask

	}

	_, err := taskColl.UpdateOne(nil, bson.M{"uid": tm.Uid}, updater)
	return err
}
