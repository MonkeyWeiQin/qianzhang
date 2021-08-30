package task

import (
	"com.xv.admin.server/model/db"
	"com.xv.admin.server/service/redisDB"
	"com.xv.admin.server/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/json-iterator/go"
	"time"
)


const (
	//TASK_TYPE_DEFAULT    int=1//=完成对局X次
	//TASK_TYPE_WIN        int=2//=对局获胜X次
	//TASK_TYPE_ALLGOLD    int=3//=累积赢得X金币
	//TASK_TYPE_ONEGOLD    int=4//=单次对局赢得X金币
	TASK_TYPE_LEVELUP    int=5//=等级or赛季等级提升X次
)
const (
	TASK_COUNT int = 5
)
var (
	TaskTypeCount int
	UserTaskDataMap = make(map[int]TaskTableDataMap) //UserTaskData
	TaskDataConfArray []*TaskDataConf //tasktabledata
)

// 玩家存储数据
type TaskData struct {
	TaskID       int64      `json:"taskid"     bosn:"taskid"`   //编号
	SubID        int        `json:"subid"      bosn:"subid"`    //副编号
	//TaskType     int64      `json:"tasktype" bosn:"tasktype"` //任务类型
	CurCount     float64    `json:"curcount"   bosn:"curcount"` //任务当前计数
	Status       int        `json:"status"     bosn:"status"`   //未完成 0， 完成未领取 1， 完成已领取 2
	//RewardGold   float64    `json:"rewardgold" bosn:"rewardgold"`//奖励金币
	//RewardExp    int64      `json:"rewardexp" bosn:"rewardexp"`//奖励经验
}

type UserTaskData struct {
	UserTaskMap 		 map[int64]TaskData    `json:"usertaskmap" bosn:"usertaskmap"`
	RefreshTaskTime      int64                 `json:"refreshtasktime" bosn:"refreshtasktime"`
}

//中间公共数据状态
type TaskTableDataMap struct {
	SubTaskCount  int
	SubTaskList   map[int64]TaskDataMap
}
//中间状态专用
type TaskDataMap struct {
	TaskID       int64      `json:"taskid"`     //编号
	SubID        int        `json:"subid"`      //副编号
	TaskType     int        `json:"tasktype"`   //任务类型
	GameID       int        `json:"gameid"`     //对应游戏
	Need         float64    `json:"need"`       //任务当前计数
	RewardGold   float64    `json:"rewardgold"` //奖励金币
	RewardExp    int32      `json:"rewardexp"`  //奖励经验
}

// TaskTableData 任务数据配置表
type TaskDataConf struct {
	TaskID        int64     `csv:"id"`
	SubID         int       `csv:"sub_id"`
	Des           string    `csv:"description"`
	ENDes         string    `csv:"english"`
	TType         int       `csv:"type"`
	Need          float64   `csv:"need"`
	Gold          float64   `csv:"gold"`
	Exp           int64     `csv:"exp"`
}

// 任务更新消息
type TaskUpdateInfo struct {
	TaskID     int64    `json:"taskid"`
	CurCount   float64  `json:"curcount"`
	IsOver     int      `json:"isover"`
}

type TempTaskList struct {
	TaskID int64
	SubID  int
}
// 领取任务奖励请求消息
type ReciveTaskID struct {
	TaskID string `json:"taskid"`
}

func init()  {
	InitTaskData()
}



// 获得主线任务列表
func GetPlayerMainTaskList(c *gin.Context) {

}
// 获得当天的日常任务列表
func GetPlayerDayTaskList(c *gin.Context) {

}

// 获得剧情任务列表
func GetPlayerPlotTaskList(c *gin.Context) {

}
//获得玩家任务列表
func getUserTaskList(userId string) (string) {

	if userId == "" {
		return "0"
	}
	taskdata, err := GetUserTaskByRedis(userId)
	if err != nil {
		return "0"
	}
	fmt.Println("usertaskmap = ", taskdata)
	return "1"//taskdata["usertaskmap"], "0")
}


func GetUserTaskByRedis(userId string) (map[string]interface{}, error) {
	userStr ,err := redis.Bytes(redisDB.Client.Hget( userId, db.REDIS_USER_TASK))
	if err != nil {
		return nil,err
	}
	user := make(map[string]interface{})
	err = jsoniter.Unmarshal(userStr,&user)
	return user,err
}

// 载入游戏任务公共数据表
func InitTaskData(){
	//filename:= common.AppBinPath+"/conf/Task.csv" // 普通赛规则
	//common.LoadCsvConfig(filename,&TaskDataConfArray)
	//self.UserTaskData = make(map[int]&taskTableDataMap)
	for _, v := range TaskDataConfArray {

		data := new(TaskDataMap)
		data.TaskID     = int64(v.TaskID)
		data.SubID      = int(v.SubID)   //int     `json:"subid"` //副编号
		data.TaskType   = int(v.TType) //int64      `json:"tasktype"`//任务类型
		data.Need       = v.Need         //float64    `json:"need"`//任务当前计数
		data.RewardGold = v.Gold         //float64    `json:"rewardgold"`//奖励金币
		data.RewardExp  = int32(v.Exp)   //int64      `json:"rewardexp"`//奖励经验
		//fmt.Println("data  = " , *data)
		_sub_id := int(v.SubID)

		sub_data := UserTaskDataMap[_sub_id]
		if sub_data.SubTaskList == nil{
			M := make(map[int64]TaskDataMap)
			sub_data.SubTaskList = M
		}
		sub_data.SubTaskList[data.TaskID] = *data
		sub_data.SubTaskCount = len(sub_data.SubTaskList)
		UserTaskDataMap[_sub_id] = sub_data
	}
	TaskTypeCount = len(UserTaskDataMap)
	//fmt.Println("UserTaskData = ", UserTaskDataMap)
	fmt.Println("-------------------InitTaskData  加载玩家任务数据表完成")
}

//刷新 玩家任务数据
// isWin 是否是赢家, // gold 本局赢的金币数量
func UpdateUserTaskInfo(userId string) (err error) {

	return nil
}

func UserTaskLevelUp(userId string) error {
	user_taskmap , err:= GetUserTaskByRedis(userId)
	if err != nil {
		return err
	}
	_task_map := user_taskmap["usertaskmap"].(map[string]interface{})

	if _task_map == nil{
		return nil//errors.New(" UpdateUserTaskInfo usertaskmap is nil , userid = " + userId)
	}
	_update_task_list := []TaskUpdateInfo{}

	for k, v := range _task_map {
		_info := v.(map[string]interface{})
		sub_id := int(_info["subid"].(float64))
		task_id := int64(_info["taskid"].(float64))
		status := int(_info["status"].(float64))
		curcount := _info["curcount"].(float64)
		fmt.Println("task_id = ", task_id, " sub_id = ", sub_id, " status = ", status, " curcount = ", curcount)
		if status > 0 { //说明这个任务已经完成
			continue
		}

		_is_change := false
		g_task_data := UserTaskDataMap[sub_id].SubTaskList[task_id]

		if g_task_data.TaskType == TASK_TYPE_LEVELUP {
			//    int=5//=等级or赛季等级提升X次
			if curcount < g_task_data.Need {
				curcount += 1
				_is_change = true
			}
		}

		_is_status := 0
		if g_task_data.Need <= curcount {
			_is_status = 1
		}
		_info["curcount"] = curcount
		_info["status"] = _is_status
		_task_map[k] = _info

		if _is_change == true{
			_send_info := new(TaskUpdateInfo)
			_send_info.TaskID = task_id
			_send_info.CurCount = curcount
			_send_info.IsOver = _is_status

			_update_task_list = append(_update_task_list, *_send_info)
		}
	}
	user_taskmap["usertaskmap"] = _task_map
	fmt.Println("UserTaskLevelUp    user_taskmap = ", user_taskmap)
	// 存储

	// 加金币

	return nil
}

//通知客户端玩家任务数据刷新
func userTaskUpPushToC(userId string, n []TaskUpdateInfo) error {
	return nil
}

// 初始化玩家任务数据表
func InitUserTaskData(userid string) string {
	user_taskdata, err := GetUserTaskByRedis(userid)
	if err != nil {
		return ""
	}
	_taskmap := make(map[int64]TaskData)
	_list := RandTask()
	for _, v := range _list {
		taskinfo := new(TaskData)
		taskinfo.TaskID = v.TaskID
		taskinfo.SubID = v.SubID
		taskinfo.Status = 0
		taskinfo.CurCount = float64(0)
		//_index_s := strconv.Itoa(int(v.TaskID))
		_taskmap[int64(v.TaskID)] = *taskinfo
	}
	user_taskdata["usertaskmap"] = _taskmap
	user_taskdata["refreshtasktime"] = time.Now().Unix()

	fmt.Println("user_taskdata = ", user_taskdata)
	// 存储操作
	// ...
	return "0"
}

func RandTask() []TempTaskList {
	//从9中不同类型中领取 5个任务
	nums := make([]int, 0)
	for len(nums) < TASK_COUNT {
		exist := false
		num := utils.RandInt(1, TaskTypeCount)
		for _, v := range nums {
			if v == num {
				exist = true
				break
			}
		}
		if !exist {
			nums = append(nums, num)
		}
	}
	fmt.Println("nums = ", nums)
	taskid_list := []TempTaskList{}
	for _, v := range nums {
		sub_task_count := UserTaskDataMap[v].SubTaskCount
		_task := new(TempTaskList)
		_task.SubID = v
		if sub_task_count > 1 {
			_index := utils.RandInt(0, sub_task_count-1)
			i := 0
			for k, _ := range UserTaskDataMap[v].SubTaskList {
				if i == _index {
					_task.TaskID = k
				}
				i++
			}

		} else {
			if sub_task_count == 1 {
				for k, _ := range UserTaskDataMap[v].SubTaskList {
					_task.TaskID = k
				}
			}
		}
		taskid_list = append(taskid_list, *_task)
	}

	return taskid_list
}



func IsRefresh(old_time int64) bool {
	_is_refresh := false
	tt := time.Now()
	zero_tm := time.Date(tt.Year(), tt.Month(), tt.Day(), 3, 0, 0, 0, tt.Location()).Unix()

	if old_time < zero_tm {
		_is_refresh = true
	}
	return _is_refresh
}

//重置当天玩家任务列表
func RefreshUserTaskList(userid string) string {

	fmt.Println("RefreshUserTaskList userid = ", userid)

	user_taskmap, err := GetUserTaskByRedis(userid)
	if err != nil {
		return ""
	}
	fmt.Println("RefreshUserTaskList user_taskmap = ", user_taskmap)
	_task_map := user_taskmap["usertaskmap"].(map[string]interface{})
	_task_refresh_time := int64(user_taskmap["refreshtasktime"].(float64))
	fmt.Println("_task_refresh_time = ", _task_refresh_time)
	if _task_map == nil {
		return " UpdateUserTaskInfo usertaskmap is nil , userid = " + userid
	}
	if IsRefresh(_task_refresh_time) {
		_taskmap := make(map[int64]TaskData)
		_list := RandTask()
		for _, v := range _list {
			taskinfo := new(TaskData)
			taskinfo.TaskID = v.TaskID
			taskinfo.SubID = v.SubID
			taskinfo.Status = 0
			taskinfo.CurCount = float64(0)
			//_index_s := strconv.Itoa(int(v.TaskID))
			_taskmap[int64(v.TaskID)] = *taskinfo
		}
		user_taskmap["usertaskmap"] = _taskmap
		user_taskmap["refreshtasktime"] = time.Now().Unix()

		fmt.Println("user_taskdata = ", user_taskmap)
		// 存储操作
		// ......
	}
	return "0"
}


// 领取任务奖励
func reciveUserTask(userId string, taskId string) (string) {

	_taskid := taskId


	fmt.Println("reciveUserTask  userid = ", userId, " req.TaskID = ", _taskid)

	user_taskmap, err := GetUserTaskByRedis(userId)
	if err != nil {
		return err.Error()
	}
	_task_map := user_taskmap["usertaskmap"].(map[string]interface{})

	if _task_map == nil {
		return " UpdateUserTaskInfo usertaskmap is nil , userid = "+ userId
	}
	fmt.Println(" _taskid = ", _taskid)
	_task := _task_map[_taskid].(map[string]interface{})

	exp := int32(0)
	gold := float64(0)
	if _task != nil {
		sub_id := int(_task["subid"].(float64))
		task_id := int64(_task["taskid"].(float64))
		status := int(_task["status"].(float64))
		curcount := _task["curcount"].(float64)
		fmt.Println("task_id = ", task_id, " sub_id = ", sub_id, " status = ", status, " curcount = ", curcount)
		if status == 1 {
			status = 2
			exp = UserTaskDataMap[sub_id].SubTaskList[task_id].RewardExp
			gold = UserTaskDataMap[sub_id].SubTaskList[task_id].RewardGold
			fmt.Println("exp = ", exp, " gold = ", gold)
		}
		if status == 2 {
			//升级操作
			//LevelUp(session, common.TASK_OVER, int(exp))
			//self.AddGold(session, gold)
		    //增加金币操作

			_task["status"] = status
			_task_map[_taskid] = _task

			user_taskmap["usertaskmap"] = _task_map
			// 存储操作
			// ...
			// 成就 [完成任务次数] start
			// 完成任务,任务获得金币 两个事件
			//model.PutAchievementToQueue(&model.AchievementQueueData{
			//	Sess: session,
			//	Actions: []model.Achievement_Type{
			//		model.Achievement_Type_9,
			//		model.Achievement_Type_10,
			//	},
			//	GameType: model.SYSTEM_TYPE,
			//	WinMoney: int(gold),
			//})
			// end

			_send_info := new(TaskUpdateInfo)
			_send_info.TaskID = task_id
			_send_info.CurCount = curcount
			_send_info.IsOver = status
			return  "" //self.App.ProtocolMarshal("", _send_info, "0")
		}
	}

	return "1"
}

