package game

import (
	"battle_rabbit/excel"
	"battle_rabbit/global"
	"battle_rabbit/iface"
	"battle_rabbit/model"
	"battle_rabbit/protocol"
	"battle_rabbit/service/log"
	"battle_rabbit/utils"
	"battle_rabbit/utils/xid"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"math/rand"
	"time"
)

/*
初始化为task配置文件为方便查找的索引map
1 初始化任务系统的主线任务集合
2 有等级限制的任务集合
*/

var (
	InitTaskDefaultMap = make(map[int]*model.Subtask)          // 玩家初始化任务系统时默认的初始任务
	LevelTaskIndexMap  = make(map[int][]*excel.TaskDataConfig) // 等级限制开放的任务集合,这里只保留等级对应的任务的第一个任务,后面的任务通过下一个任务ID关联
	DayTaskLvIndexMap  = make(map[int][]*model.Subtask)        // 每日任务根据等级分类的集合
)

/*
初始化索引
启动每日任务刷新计时器
启动主线任务定时存储落地
(
	内存 ==[sync]==> redis
	redis ==[async]==> mongo
)
*/
func (g *Game) InitTaskIndex() {
	for _, conf := range excel.TaskMainDataConf {
		if conf.IsInitTask {
			subTask := &model.Subtask{
				TaskId:    conf.Id,
				Type:      conf.Type,
				LimitTime: conf.LimitTime,
				Condition: conf.Condition,
				Target:    conf.Count,
				RewardId:  conf.Reward,
			}
			InitTaskDefaultMap[conf.Id] = subTask
			continue
		}

		if conf.IsFirstTask && conf.Level > 0 {
			LevelTaskIndexMap[conf.Id] = append(LevelTaskIndexMap[conf.Id], conf)
		}
	}

	for _, conf := range excel.TaskDayDataConf {
		subTask := &model.Subtask{
			TaskId:    conf.Id,
			Type:      conf.Type,
			LimitTime: conf.LimitTime,
			Condition: conf.Condition,
			Target:    conf.Count,
			RewardId:  conf.Reward,
		}
		DayTaskLvIndexMap[conf.Id] = append(DayTaskLvIndexMap[conf.Id], subTask)
	}
}

type TaskCompletePushMsg struct {
	TaskId int `json:"taskId"`
	Type   int `json:"type"`
	State  int `json:"state"`
	Num    int `json:"num"`
}

func UpdateTask(task *model.TaskModel, sess iface.ISession, info []*model.TaskProgress) {

	return
	var (
		overTaskList []*TaskCompletePushMsg
		mainUpdate   = false
		dayUpdate    = false
	)
	for _, progress := range info {
		// 主线任务
		for _, subtask := range task.MainTask {
			if subtask.Type == progress.Type &&
				subtask.Condition == progress.Condition &&
				subtask.State == model.TaskStateBegin {

				if progress.Num != 0 {
					subtask.Progress += progress.Num
				}else{ // lv
					subtask.Progress = progress.Lv
				}

				if subtask.Progress >= subtask.Target { // 任务完成
					subtask.State = model.TaskStateComplete
					overTaskList = append(overTaskList, &TaskCompletePushMsg{
						TaskId: subtask.TaskId,
						Type:   subtask.Type,
						State:  subtask.State,
						Num:    subtask.Target,
					})
				}
				mainUpdate = true
				continue
			}
		}
		// 每日任务
		for _, subtask := range task.DayTask {
			if subtask.Type == progress.Type &&
				subtask.Condition == progress.Condition &&
				subtask.State == model.TaskStateBegin {

				subtask.Progress += progress.Num
				if subtask.Progress >= subtask.Target { // 任务完成
					subtask.State = model.TaskStateComplete
					overTaskList = append(overTaskList, &TaskCompletePushMsg{
						TaskId: subtask.TaskId,
						Type:   subtask.Type,
						State:  subtask.State,
						Num:    subtask.Target,
					})
				}
				dayUpdate = true
				continue
			}
		}
	}
	if len(overTaskList) > 0 {
		err := task.Save(mainUpdate, dayUpdate)
		if err != nil {
			log.Error(err)
			return
		}
		protocol.NoticeUserTaskComplete(sess, overTaskList)
	}
}

// 刷新每日任务
// 任务完成没有领取奖励的需要通过邮件发放给玩家
func flushDayTask(player *Player) {
	var (
		tm    = player.Task
		items []*global.Item
	)
	if tm.DayTaskTime == utils.GetDayBeginTimeStamp(utils.GetNowTimeStamp()) {
		return
	}

	for _, subtask := range tm.DayTask {
		// 完成了,没有领取,需要将奖励通过邮件发送到用户邮箱
		if subtask.State == model.TaskStateComplete {
			// 获得奖励列表
			rewardConf := excel.TaskRewardDataConf[subtask.RewardId]
			items = append(items, rewardToItems(rewardConf)...)
		}
	}
	if len(items) > 0 {
		// 创建邮件
		mail := &model.MailModel{
			Mid:          xid.New().String(),
			UserId:       player.Uid,
			Name:         "每日任务奖励",
			Sender:       "",
			SenderID:     0,
			SendTime:     int(time.Now().Unix()),
			Validity:     0,
			Title:        "每日任务奖励",
			Content:      fmt.Sprintf("请领取您在%s的每日任务奖励!", utils.TransformationTime(tm.DayTaskTime, "Y-m-d")),
			State:        0,
			ReceiveState: false,
			Attachment:   items,
		}
		err := model.GetMailCollection().InsertOne(nil, mail)
		if err != nil {
			log.Error("每日任务奖励发放失败!!! ", err)
		}
	}

	// 刷新新的任务给玩家
	tasks := DayTaskLvIndexMap[player.Role.RLv]

	var newTask = make(map[int]*model.Subtask)
	if len(tasks) > 0 {
		for _, task := range tasks {
			if rand.Intn(100) > 50 && len(newTask) < 5 {
				t := *task
				if t.LimitTime != 0 {
					t.LimitTime = (t.LimitTime * 60) + int(utils.GetNowTimeStamp())
				}
				newTask[t.TaskId] = &t
			}
		}
		if len(newTask) == 0 {
			t := *tasks[0]
			if t.LimitTime != 0 {
				t.LimitTime = (t.LimitTime * 60) + int(utils.GetNowTimeStamp())
			}
			newTask[t.TaskId] = &t
		}
	}
	tm.DayTask = newTask
	tm.DayTaskTime = utils.GetDayBeginTimeStamp(utils.GetNowTimeStamp())
	_, err := model.GetTaskCollection().UpdateOne(nil, bson.M{"uid": tm.Uid}, bson.M{"dayTask": tm.DayTask, "dayTaskTime": tm.DayTaskTime})
	if err != nil {
		log.Error("更新每日任务失败!!")
		log.Error(err)
	}
}
