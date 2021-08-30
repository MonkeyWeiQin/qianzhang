package lobby

import (
	"com.xv.admin.server/model/db"
	"com.xv.admin.server/service/redisDB"
	"com.xv.admin.server/utils"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/json-iterator/go"
	"time"
)

//广告视频奖励的表格

type AdvertisingVideo_CSV struct {
	ID       int     `csv:"id"`       // 类型编号
	Name     string  `csv:"name"`     // 内部使用，不显示
	AvType   int     `csv:"av_type"`  // 奖励种类：随机 倍数 次数 固定值 必须看
	Multiple float64 `csv:"multiple"` // 倍数
	MinGold  float64 `csv:"min_gold"` // 最小数
	MaxGold  float64 `csv:"max_gold"` // 最大数
	Times    int     `csv:"times"`    // 次数
}

type AdvertisingVideo_Data struct {
	ID       int     `json:"id"       bson:"id"`       // 类型编号
	Name     string  `json:"name"     bson:"name"`     // 内部使用，不显示
	AvType   int     `json:"av_type"  bson:"av_type"`  // 奖励种类：随机 倍数 次数 固定值 必须看
	Multiple float64 `json:"multiple" bson:"multiple"` // 倍数
	MinGold  float64 `json:"min_gold" bson:"min_gold"` // 最小数
	MaxGold  float64 `json:"max_gold" bson:"max_gold"` // 最大数
	Times    int     `json:"times"    bson:"times"`    // 次数
}

// User_AdvertisingVideo_Info 记录玩家的观看视频的相关信息
type User_AdvertisingVideo_Info struct {
	//ID        int         `json:"id"       bson:"id"`0       // 类型编号
	AV_ID     int   `json:"av_id"    bson:"av_id"` // 奖励种类：随机 倍数 次数 固定值 必须看
	Times     int   `json:"times"    bson:"times"` // 已经观看次数
	TimeStamp int64 `json:"stamp"    bson:"stamp"` // 时间戳
}

// 记录玩家的观看视频的相关信息
type User_AdvertisingVideo_Data struct {
	Today         string `json:"today"    bosn:"today"` // 日期
	BrowseABSInfo map[int]User_AdvertisingVideo_Info
}
var (
	AdvertisingVideo_DataMap = make(map[int]*AdvertisingVideo_Data)
)

const (
	ABS_TYPE_RANDOM      int = 1 // 随机
	ABS_TYPE_MULTIPLE    int = 2 // 倍数
	ABS_TYPE_AWARDS      int = 3 // 奖励次数
	ABS_TYPE_FIXED_VALUE int = 4 // 固定值
	ABS_TYPE_MUST_SEE    int = 5 // 必须看
	ABS_TYPE_DAILY_TIMES int = 6 // 每天观看次数
)

func init() {
	InitPublicAdvertisingVideo()
}

func InitPublicAdvertisingVideo() {
	//filename := "系统路径" + "/conf/AdvertisingVideo.csv" // 在线奖励数据
	var tmpAc []*AdvertisingVideo_CSV
	//common.LoadCsvConfig(filename, &tmpAc)
	for _, v := range tmpAc {
		data := new(AdvertisingVideo_Data)
		data.ID = v.ID
		data.Name = v.Name
		data.AvType = v.AvType
		data.Multiple = v.Multiple
		data.MinGold = v.MinGold
		data.MaxGold = v.MaxGold
		data.Times = v.Times
		AdvertisingVideo_DataMap[v.ID] = data
	}

	fmt.Println("InitPublicAdvertisingVideo ----- 加载广告点数据完成")
}
func GetABSVideoInfoByRedis(userId string) (User_AdvertisingVideo_Data, error) {
	userStr, err := redis.Bytes(redisDB.Client.Hget(db.TBALE_NAME_ADVERTISING_VIDEO, userId))
	if err != nil {
		fmt.Println("GetABSVideoInfoByRedis redis.Bytes err is : ", err.Error())
	}
	user := new(User_AdvertisingVideo_Data) // make(map[string]interface{})
	user.BrowseABSInfo = make(map[int]User_AdvertisingVideo_Info)

	err = jsoniter.Unmarshal([]byte(userStr), user)
	if err != nil {
		fmt.Println("GetABSVideoInfoByRedis Unmarshal err is : ", err.Error())
	}
	return *user, nil
}

func InitAdvertisingVideoDataForUser(userId string) string {
	if userId == "" {
		return "10104"
	}

	userData, err := GetABSVideoInfoByRedis(userId)
	if err != nil {
		fmt.Println("InitAdvertisingVideoDataForUser error :", err.Error())
	}
	userData.Today = time.Now().Format("2006-01-02")

	userData.BrowseABSInfo = make(map[int]User_AdvertisingVideo_Info)

	for k, vv := range AdvertisingVideo_DataMap {
		if vv.AvType == ABS_TYPE_AWARDS || vv.AvType == ABS_TYPE_DAILY_TIMES || vv.AvType == ABS_TYPE_RANDOM {
			abs_data := new(User_AdvertisingVideo_Info)
			abs_data.AV_ID = k
			abs_data.Times = vv.Times
			userData.BrowseABSInfo[k] = *abs_data
		}
	}

	M, _ := jsoniter.Marshal(userData)
	_, err = redisDB.Client.Hset(db.TBALE_NAME_ADVERTISING_VIDEO, userId, M)
	if err != nil {
		fmt.Println(" InitAdvertisingVideoDataForUser error is : ", err.Error())
		return err.Error()
	}

	return "0"
}

// 获得广告点相关配置数据
func GetAdvertisingVideoData(userId string, msg map[string]interface{}) (map[int]*AdvertisingVideo_Data, string) {
	return AdvertisingVideo_DataMap, "0"
}

//初始化 有序集合排行榜  把每个人都加进去
//func InitAdvertisingVideoDataForUser(userid string) string {
//
//	err := InitAdvertisingVideoDataForUser(userid)
//	if err != "0" {
//		return err
//	}
//	return "0"
//}
// 获得玩家的视频浏览记录
func getUserAdvertisingVideoInfo(userId string, msg map[string]interface{}) string {

	userData, err := GetABSVideoInfoByRedis(userId)
	if err != nil {
		return "0"
	}
	s_today := time.Now().Format("2006-01-02")
	if userData.Today != s_today || userData.Today == "" {
		// 说明是新的一天
		userData.Today = s_today
		userData.BrowseABSInfo = make(map[int]User_AdvertisingVideo_Info)

		for k, vv := range AdvertisingVideo_DataMap {
			if vv.AvType == ABS_TYPE_AWARDS || vv.AvType == ABS_TYPE_DAILY_TIMES || vv.AvType == ABS_TYPE_RANDOM {
				abs_data := new(User_AdvertisingVideo_Info)
				abs_data.AV_ID = k
				abs_data.Times = vv.Times
				userData.BrowseABSInfo[k] = *abs_data
			}
		}
	}
	M, _ := jsoniter.Marshal(userData)
	_, err = redisDB.Client.Hset(db.TBALE_NAME_ADVERTISING_VIDEO, userId, M)
	if err != nil {
		fmt.Println(" getUserAdvertisingVideoInfo error is : ", err.Error())
		return err.Error()
	}
	return "0"
}

type UsersBrowseAVTypeRequest struct {
	AV_ID     int    `json:"av_id"      bson:"av_id"`     // 奖励种类：随机 倍数 次数 固定值 必须看
	Task_Type int    `json:"task_type"  bson:"task_type"` // 任务类型，日常 1 成就 2 登陆 3
	Task_ID   string `json:"task_id"    bson:"task_id"`   // 对应的任务id号
}

type UsersBrowseAVTypeResult struct {
	State int     `json:"state"    bson:"state"` // 状态：成功 1 ， 失败 0
	AV_ID int     `json:"av_id"    bson:"av_id"` // 奖励种类：随机 倍数 次数 固定值 必须看
	Times int     `json:"times"    bson:"times"` // 剩余次数
	Gold  float64 `json:"gold"     bson:"gold"`  // 获得金币数
}

func UsersBrowseABSSuccessfully(userId string, msg map[string]interface{}) string {

	var req = new(UsersBrowseAVTypeRequest)
	//error := common.ValidateParam(&msg, req)
	//if error != "" {
	//	log.Error("reciveUserSeasonReward 1 _____: ", error)
	//	return self.App.ProtocolMarshal("", nil, "20206")
	//}
	av_id := req.AV_ID
	//userid := session.GetUserID()
	fmt.Println("UsersBrowseABSSuccessfully  av_id = ", av_id, " userid = ", userId)
	userData, err := GetABSVideoInfoByRedis(userId)
	if err != nil {
		//log.Error(" UsersBrowseABSSuccessfully err = ", err.Error())
		return "0"
	}

	data := new(UsersBrowseAVTypeResult)
	data.AV_ID = av_id
	gold := float64(0)
	av_type := AdvertisingVideo_DataMap[av_id].AvType
	switch av_type {
	case ABS_TYPE_RANDOM: //1	随机

		//cur_times := 0
		//cur_TimeStamp := 0
		s_today := time.Now().Format("2006-01-02")
		if userData.Today != s_today {
			// 说明不是同一天,说明是新的一天，容错处理，
			userData.Today = s_today
			if times_info, ok := userData.BrowseABSInfo[av_id]; ok {
				if v, ok := AdvertisingVideo_DataMap[av_id]; ok {
					times_info.Times = v.Times
					times_info.TimeStamp = 0
				}
				userData.BrowseABSInfo[av_id] = times_info
			}
		}

		if v, ok := AdvertisingVideo_DataMap[av_id]; ok {
			if times_info, ok := userData.BrowseABSInfo[av_id]; ok {
				if times_info.Times <= v.Times && times_info.Times > 0 { //还有剩余次数
					min_gold := v.MinGold
					max_gold := v.MaxGold
					gold = float64(utils.RandInt(int(min_gold), int(max_gold)))

					data.State = 1
					data.Times = times_info.Times - 1
					times_info.Times -= 1
					times_info.TimeStamp = time.Now().Unix()
					userData.BrowseABSInfo[av_id] = times_info
				} else { // 没有剩余次数
					data.State = 0
					data.Times = 0
				}
			}
		}
		fmt.Println("av_type = ", av_type, " gold = ", gold)
		break
	case ABS_TYPE_MULTIPLE: //2	倍数
		//日常 1 成就 2 登陆 3
		if req.Task_Type == 1 {
			_taskid := req.Task_ID
			//UserTaskDataMap
			if _taskid != "" {
				//user_taskmap, err := GetUserTaskByRedis(userId)
				//if err != nil {
				//	return err.Error()
				//}
				//_task_map := user_taskmap["usertaskmap"].(map[string]interface{})
				//if _task_map == nil {
				//	return "UpdateUserTaskInfo usertaskmap is nil , userid = "+ userId
				//}
				//fmt.Println(" _taskid = ", _taskid)
				//_task := _task_map[_taskid].(map[string]interface{})
				//
				//if _task != nil {
				//	sub_id := int(_task["subid"].(float64))
				//	task_id := int64(_task["taskid"].(float64))
				//	status := int(_task["status"].(float64))
				//	//curcount := _task["curcount"].(float64)
				//	fmt.Println("task_id = ", task_id, " sub_id = ", sub_id, " status = ", status)
				//	if status == 1 {
				//		if v, ok := AdvertisingVideo_DataMap[av_id];ok{
				//			gold = UserTaskDataMap[sub_id].SubTaskList[task_id].RewardGold *  (v.Multiple - 1)
				//			data.State = 1
				//		}
				//	}
				//}
			}
		}

		if req.Task_Type == 2 {
			//AchievementCsvDataMap
			//taskid, _ := strconv.Atoi(req.Task_ID)
			//if v, ok := model.AchievementCsvDataMap[taskid];ok{
			//	if v.A_Reward_Type == 1{
			//		_gold, _ := strconv.ParseFloat(v.A_Reward_Data, 64)
			//		if vv, ok := model.AdvertisingVideo_DataMap[av_id];ok{
			//			gold = _gold *  (vv.Multiple - 1)
			//			data.State = 1
			//		}
			//	}
			//}
		}

		if req.Task_Type == 3 {
			//
			//day_index, _ := strconv.Atoi(req.Task_ID)
			//if v, ok := self.sevendaydata[day_index];ok{
			//	if vv, ok := model.AdvertisingVideo_DataMap[av_id];ok{
			//		gold = v * (vv.Multiple - 1)
			//		data.State = 1
			//	}
			//}
		}

		break
	case ABS_TYPE_AWARDS: //3	奖励次数
		//s_today := time.Now().Format("2006-01-02")
		//userinfo, err1 := model.GetLockyTableInfoByRedis(userid)
		//if err1 != nil {
		//	fmt.Println()
		//}
		//if userData.Today != s_today {
		//	// 说明不是同一天,说明是新的一天，容错处理，
		//	userData.Today = s_today
		//	if times_info, ok := userData.BrowseABSInfo[av_id];ok{
		//		data.State = 1
		//		if v, ok := model.AdvertisingVideo_DataMap[av_id];ok{
		//			data.Times = v.Times - 1
		//			times_info.Times = data.Times
		//			userinfo.LockyCount += 1
		//		}
		//		userData.BrowseABSInfo[av_id] = times_info
		//	}
		//}else{
		//	if v, ok := model.AdvertisingVideo_DataMap[av_id];ok{
		//		if times_info,ok := userData.BrowseABSInfo[av_id];ok{
		//			if times_info.Times <= v.Times && times_info.Times > 0{ //还有剩余次数
		//				data.State = 1
		//				data.Times = times_info.Times - 1
		//				times_info.Times -= 1
		//				userData.BrowseABSInfo[av_id] = times_info
		//				userinfo.LockyCount += 1
		//			}else{// 没有剩余次数
		//				data.State = 0
		//				data.Times = 0
		//				//userinfo.LockyCount = int64(0)
		//			}
		//		}
		//	}
		//}
		//_userInfo, _ := jsoniter.Marshal(userinfo)
		//_, err2 := global.Client.Hset(model.TABLE_NAME_LOCKY_TABLE, userid, _userInfo)
		//if err2 != nil {
		//	log.Error(err2.Error())
		//	fmt.Println(" runLockyTableBegin 4 ", err2.Error())
		//	return self.App.ProtocolMarshal("", nil, "0")
		//}
		break
	case ABS_TYPE_FIXED_VALUE: //4	固定值
		if v, ok := AdvertisingVideo_DataMap[av_id]; ok {
			gold = v.MinGold
			data.State = 1
			fmt.Println("av_type = ", av_type, " gold = ", gold)
		}
		break
	case ABS_TYPE_MUST_SEE: //5	必须看
		data.State = 1
		break
	case ABS_TYPE_DAILY_TIMES: //6	每天观看次数
		s_today := time.Now().Format("2006-01-02")
		if userData.Today != s_today {
			// 说明不是同一天,说明是新的一天
			userData.Today = s_today
			if times_info, ok := userData.BrowseABSInfo[av_id]; ok {
				times_info.Times = 1
				userData.BrowseABSInfo[av_id] = times_info
				data.State = 1
				if v, ok := AdvertisingVideo_DataMap[av_id]; ok {
					data.Times = v.Times - 1
					gold = v.MinGold
				}
			}
		} else {
			if v, ok := AdvertisingVideo_DataMap[av_id]; ok {
				if times_info, ok := userData.BrowseABSInfo[av_id]; ok {
					if times_info.Times < v.Times { //还有剩余次数
						gold = v.MinGold
						fmt.Println("av_id = ", av_id, " gold = ", gold)
						data.State = 1
						data.Times = v.Times - times_info.Times - 1
						times_info.Times += 1
						userData.BrowseABSInfo[av_id] = times_info
					} else { // 没有剩余次数
						data.State = 0
						data.Times = 0
					}
				}
			}
		}
		break
	default:
		break
	}
	data.Gold = gold
	M, _ := jsoniter.Marshal(userData)
	_, err = redisDB.Client.Hset(db.TBALE_NAME_ADVERTISING_VIDEO, userId, M)
	if err != nil {
		fmt.Println(" UsersBrowseABSSuccessfully error is : ", err.Error())
		return err.Error()
	}

	//e := model.UpdateUserWalletAndPush(session, gold, 0)
	//if e != nil {
	//	log.Error("Lobby buySkinByGold() 1 : %s", e.Error())
	//	return self.App.ProtocolMarshal("", nil, "6101")
	//}

	return "0"
}
