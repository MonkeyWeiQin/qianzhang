package lobby

import (
	"com.xv.admin.server/config"
	"com.xv.admin.server/model/db"
	"com.xv.admin.server/service/mgoDB"
	"com.xv.admin.server/service/redisDB"
	"context"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/json-iterator/go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)
var G_SevenDayData  map[int]float64

//--------------------------------------------------------登陆奖励活动-----------------------------------------------------------------
type DayInfoData struct {
	DayInfo   string                    `json:"dayinfo" bson:"dayinfo"`
	Status    int                       `json:"status" bson:"status"`
}

type SevenDayLoginData struct {
	Id            primitive.ObjectID     `json:"-" bson:"-"`
	UserId	      string	             `json:"userId"  bson:"userId"`
	Today         int64                  `json:"today" bson:"today"`
	Day_Statue    map[int]DayInfoData    `json:"day_status"  bson:"day_status"`
}

func (user *SevenDayLoginData) CollectionName() string {
	return config.CollUsers
}

func (user *SevenDayLoginData) GetId() primitive.ObjectID {
	return user.Id
}

func (user *SevenDayLoginData) SetId(id primitive.ObjectID) {
	user.Id = id
}

func (user *SevenDayLoginData)IsExists(userId string)(bool, string){

	return false, ""
}

func (user *SevenDayLoginData)GetData(userId string)(b bool, err error){
	filter := bson.D{{"userId", userId}}
	finder := mgoDB.NewOneFinder(user).Where(filter)
	return mgoDB.GetMgo().FindOne(context.TODO(), finder)
}

func (user *SevenDayLoginData)SetData(userId string, userdata *SevenDayLoginData)(string){

	return ""
}


func Init7DayLogin(userid string) string {

	if userid == "" {
		return "0"
	}

	ok, err := redis.Bool(redisDB.Client.Hexists(db.TABLE_NAME_SEVEN_DAY, userid))
	if err != nil {
		fmt.Println(" Init7DayLogin 1 ", err.Error())
	}
	if ok {
		fmt.Println(" Init7DayLogin 2 ", err.Error())
		return "10103"
	} else {
		fmt.Println(" Init7DayLogin 3 ")
		//初始化 七天活动表
		data := new(SevenDayLoginData)
		data.UserId = userid
		data.Today = time.Now().Unix()
		_data := make(map[int]DayInfoData)
		_info := new(DayInfoData)
		for i := 1; i <= 7; i++ {
			_data[i] = *_info
		}
		_info.DayInfo = time.Now().Format("2006-01-02") //记录第一天的
		_info.Status = 0

		_data[1] = *_info
		data.Day_Statue = _data

		userInfo, _ := jsoniter.Marshal(data)
		_, err1 := redisDB.Client.Hset(db.TABLE_NAME_SEVEN_DAY, userid, userInfo)
		if err1 != nil {
			fmt.Println(" Init7DayLogin 4 ", err1.Error())
			return err1.Error()
		}
	}
	return "0"
}

//登陆的时候，需要刷新7天登陆数据，并给予玩家奖励
func UpdataSevenDayLoginData(userId string) (string) {

	fmt.Println("UpdataSevenDayLoginData", userId)
	userStr, err := redis.Bytes(redisDB.Client.Hget(db.TABLE_NAME_SEVEN_DAY, userId))
	if err != nil {
		fmt.Println(" UpdataSevenDayLoginData 1 ", err.Error())
		return "0"
	}
	_today := time.Now().Format("2006-01-02")
	user := new(SevenDayLoginData) // make(map[string]interface{})

	err = jsoniter.Unmarshal([]byte(userStr), user)
	if err != nil {
		return "0"
	}
	user.Today = time.Now().Unix()
	day_index := 0
	is_revert := false
	for i := 1; i <= 7; i++ {
		if i == 1 {
			if user.Day_Statue[1].Status == 0 {
				//第一天处理
				day_index = 1
				_dayinfo := new(DayInfoData)
				_dayinfo.DayInfo = _today
				_dayinfo.Status = 1
				user.Day_Statue[1] = *_dayinfo
				break
			}
			continue
		}
		_date := user.Day_Statue[i-1] //取昨天的

		if _date.DayInfo == "" { //如果昨天的日期为空，则说明断了，需要重置
			is_revert = true //重置
			break
		}

		//if common.IsYesterday(_date.DayInfo) && { //取出来的日期是昨天 和昨天进行比较，是否相同
		if _date.Status == 1 {
			if IsYesterday(_date.DayInfo) {
				//领取过奖励，且是昨天
				day_index = i
				_dayinfo := new(DayInfoData)
				_dayinfo.DayInfo = _today
				_dayinfo.Status = 1
				user.Day_Statue[i] = *_dayinfo
				break
			}

			if i == 7 {
				//说明，七天登陆已经完成
				//今天可以领取和第7天相同的奖励
				day_index = i
				_dayinfo := new(DayInfoData)
				_dayinfo.DayInfo = _today
				_dayinfo.Status = 1
				user.Day_Statue[i] = *_dayinfo
			}
		} else {
			is_revert = true
			break
		}
	}

	sevenday_reward := 0.0
	if day_index == 0 && is_revert == true {
		//重置操作
		fmt.Println("is_revert == true ")
		_dayinfo := new(DayInfoData)
		for i := 1; i < 7; i++ {
			user.Day_Statue[i] = *_dayinfo
		}
		_dayinfo.DayInfo = _today
		_dayinfo.Status = 1
		user.Day_Statue[1] = *_dayinfo
		day_index = 1
	}
	//获得今天的奖励金币数量
	sevenday_reward = G_SevenDayData[day_index]
	fmt.Println("serverday_reward = ", sevenday_reward)
	// 发金币给玩家
	//err = model.UpdateUserWalletAndPush(session, sevenday_reward, 0)
	//if err != nil {
	//	log.Error("Lobby UpdataSevenDayLoginData() 1 : %s", err.Error())
	//	return self.App.ProtocolMarshal("", nil, "6101")
	//}

	fmt.Println("day_index = ", day_index)
	userInfo, _ := jsoniter.Marshal(user)
	_, err1 := redisDB.Client.Hset(db.TABLE_NAME_SEVEN_DAY, userId, userInfo)
	if err1 != nil {
		fmt.Println(" UpdataSevenDayLoginData 4 ", err1.Error())
		return "0"
	}
	// 成就 [登录获得金币] start
	//model.PutAchievementToQueue(&model.AchievementQueueData{
	//	Sess:     session,
	//	Actions:  []model.Achievement_Type{model.Achievement_Type_2},
	//	GameType: model.SYSTEM_TYPE,
	//	WinMoney: int(sevenday_reward),
	//})
	// end

	return "0"
}

// 返回7天登陆活动的标准数据
func GetSevenDayData(userId string, msg map[string]interface{}) (string) {
	return "0"
}

// 返回 自己的七天登录情况
func GetSelfSevenDayData(userId string, msg map[string]interface{}) (string) {

	fmt.Println("getSelfSevenDayData  userid = ", userId)
	userStr, err := redis.Bytes(redisDB.Client.Hget(db.TABLE_NAME_SEVEN_DAY, userId))
	if err != nil {
		fmt.Println(" getSelfSevenDayData 1 ", err.Error())
		return "0"
	}
	//_today := time.Now().Format("2006-01-02")
	user := new(SevenDayLoginData) // make(map[string]interface{})
	err = jsoniter.Unmarshal([]byte(userStr), user)
	if err != nil {
		return "0"
	}
	user.Today = time.Now().Unix()
	is_revert := true
	//超过七天
	is_over7day := false
	for i := 1; i <= 7; i++ {
		_dayinfo := user.Day_Statue[i].DayInfo
		if IsYesterday(_dayinfo) {
			is_revert = false
		}
		if IsToday(_dayinfo) {
			is_revert = false
		}

		if i == 7 {
			_yesterdayinfo := user.Day_Statue[i-1].DayInfo
			if _dayinfo == "" && _yesterdayinfo != "" {
				//这是第一轮 7天
				//fmt.Println("1 ")
				break
			}

			if _dayinfo != "" && _yesterdayinfo != "" {
				if IsToday(_dayinfo) {
					//第七天是今天  无论它签了没签
					is_revert = false
					_stauts := user.Day_Statue[i].Status
					if _stauts == 1 {
						is_over7day = false
						//fmt.Println("2.0 ")
						break
					} else {
						if _stauts == 0 {
							is_over7day = true
							//fmt.Println("2.1 ")
							break
						}
					}
					//fmt.Println("2 .2")

					break
				} else {
					if IsYesterday(_dayinfo) {
						//第七天是昨天
						//并且 要签过的
						_stauts := user.Day_Statue[i].Status
						if _stauts == 0 {
							is_revert = true
							//fmt.Println("3 ")
							break
						} else {
							if _stauts == 1 {
								is_revert = false
								is_over7day = true
								//fmt.Println("4 ")
							}
						}
					}
				}
			}
		}

	}

	_dayinfo := new(DayInfoData)
	_dayinfo.DayInfo = ""
	_dayinfo.Status = 0

	fmt.Println(" is_revert = ", is_revert, " is_over7day = ", is_over7day)
	if is_revert == true {

		for j := 1; j <= 7; j++ {
			user.Day_Statue[j] = *_dayinfo
		}
	} else {
		if is_over7day == true {
			_dayinfo.DayInfo = time.Now().Format("2006-01-02")
			_dayinfo.Status = 0

			user.Day_Statue[7] = *_dayinfo
		}
	}
	userInfo, _ := jsoniter.Marshal(user)
	_, err1 := redisDB.Client.Hset(db.TABLE_NAME_SEVEN_DAY, userId, userInfo)
	if err1 != nil {
		fmt.Println(" getSelfSevenDayData 4 ", err1.Error())
		return "0"
	}
	fmt.Println(" user.Day_Statue = ", user)
	return "0"
}

// 领取当前的奖励
func ReceiveTodaySevenDayReward(userId string, msg map[string]interface{}) (string) {
	return "0"
}


/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func IsToday(day_str string) bool {
	_today := time.Now().Format("2006-01-02")
	//fmt.Println("today = ", _today)
	if _today == day_str {
		return true
	} else {
		return false
	}
	return false
}
func IsYesterday(t string) bool {
	_yesterday := time.Now().Add(-24 * time.Hour).Format("2006-01-02")

	if _yesterday == t {
		return true
	}
	return false
}