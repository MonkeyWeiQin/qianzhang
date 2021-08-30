package lobby

import (
	"com.xv.admin.server/model/db"
	"com.xv.admin.server/service/mgoDB"
	"com.xv.admin.server/service/redisDB"
	"context"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/json-iterator/go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sort"
	"strconv"
	"time"
)

// 使用redis的有序集合进行排序

//服务器启动时，初始化活动列表
func InitPublicRankingList() {
	fmt.Println("InitPublicRankingList ----------------------------- 初始化排行榜数据")
	//初始化最大列表数量
	coll := mgoDB.GetMgo().DB().Collection(db.TABLE_NAME_PUB_RANKING)
	//加载激活状态下的活动
	cursor,err :=  coll.Find(context.Background(),bson.M{})
	if err != nil {
		//log.Error("Lobby Run error , Find mongo  pub_ranking failed ==> :",err)
		return
	}

	for cursor.Next(nil)  {

		g := new(PublicRankingList)
		err = cursor.Decode(g)
		if err != nil {
			return
		}

		data_day := []Ranking{}
		data_week := []Ranking{}
		data_month := []Ranking{}
		j := 0
		_count := len(g.DayRankingList)
		for (j < _count){
			v := new(Ranking)
			v.UserName = g.DayRankingList[j].UserName
			v.UserId = g.DayRankingList[j].UserId
			v.RankingNo = g.DayRankingList[j].RankingNo
			v.Head = g.DayRankingList[j].Head
			v.Bonus = g.DayRankingList[j].Bonus
			data_day = append(data_day, *v)
			j++
		}
		j = 0
		for (j < _count){
			v := new(Ranking)
			v.UserName = g.DayRankingList[j].UserName
			v.UserId = g.DayRankingList[j].UserId
			v.RankingNo = g.WeekRankingList[j].RankingNo
			v.Head = g.WeekRankingList[j].Head
			v.Bonus = g.WeekRankingList[j].Bonus
			data_week = append(data_week, *v)
			j++
		}
		j = 0
		for (j < _count){
			v := new(Ranking)
			v.UserName = g.DayRankingList[j].UserName
			v.UserId = g.DayRankingList[j].UserId
			v.RankingNo = g.MonthRankingList[j].RankingNo
			v.Head = g.MonthRankingList[j].Head
			v.Bonus = g.MonthRankingList[j].Bonus
			data_month = append(data_month, *v)
			j++
		}
		_data := PublicRankingList{
			ID: primitive.NewObjectID(),
			DayRankingList:data_day,//		[]Ranking			`json:"day" 		bson:"day"`		//日排行榜
			WeekRankingList:data_day,//		[]Ranking			`json:"week" 		bson:"week"`	//周排行榜
			MonthRankingList:data_day,//		[]Ranking			`json:"month" 		bson:"month"`	//月排行榜
			RefreshDayTime:g.RefreshDayTime,//			int64					`json:"refreshday" 		bson:"refreshday"`
			RefreshWeekTime:g.RefreshWeekTime,//		int64					`json:"refreshweek" 		bson:"refreshweek"`
			RefreshMonthTime:g.RefreshMonthTime,//		int64					`json:"refreshmonth" 		bson:"refreshmonth"`
		}
		conn := redisDB.Client.Get()
		dataStr ,_ := jsoniter.Marshal(_data)
		_, err := conn.Do("HSETNX", db.TABLE_NAME_PUB_RANKING, db.TABLE_NAME_PUB_RANKING, dataStr)
		if err != nil{
			fmt.Println("InitPublicRankingList error is ", err.Error())
		}
		conn.Close()
	}

	//添加一个定时每天刷新的驱动时间
}

//战力排行榜
//关闭界面的时候再向上请求战力升级信息，包含英雄总战力

//解锁排行榜

//赛季排行榜
//关卡排行榜
//--------------------------------------------功能函数-----------------------------------------------------
//判断时间是当年的第几周
func DayByDate() int {
	return int(time.Now().Day())
}
//获得今天是周几
func DayInWeek() int {
	_day := int(time.Now().Weekday())
	return _day
}

//判断时间是当年的第几周
func WeekByDate() int {
	t :=time.Now()
	yearDay := t.YearDay()
	yearFirstDay := t.AddDate(0, 0, -yearDay+1)
	firstDayInWeek := int(yearFirstDay.Weekday())

	//今年第一周有几天
	firstWeekDays := 1
	if firstDayInWeek != 0 {
		firstWeekDays = 7 - firstDayInWeek + 1
	}
	var week int
	if yearDay <= firstWeekDays {
		week =  1
	} else {
		week = (yearDay-firstWeekDays)/7 + 2
	}
	return week
}
//判断时间是当年的第几月
func MonthByDate() int {
	monthint := time.Now().Month()
	return int(monthint)
}

//------------------------------------功能需求-----表格：pub_ranking------------------------------------------//
//从redis公共数据库中，获得排行榜内容，后台管理可以查看该数据


type MobileRankingRequest struct {
	RankingType   	int `json:"type"`  //获得排行榜类型
	RankingPages 	int `json:"pag"`	//获得第几页
}
type RankingInfo struct {
	UserName 			string				`json:"username"`		//玩家名字
	RankingNo			int					`json:"nomber"`		//排名名次
	Icon 				string				`json:"icon"`			//头像
	Bonus 				float64			    `json:"bonus"`			//获奖总额
}
type RankingResult struct {
	RankingType 		int 				`json:"type"`
	List  				[]RankingInfo   	`json:"list"`
}

//---------------------------------------有序集合排行榜-------------------------------------------------------
//初始化 有序集合排行榜  把每个人都加进去
func InitUserRankingData(userid string) ( string){

	if userid == "" {
		return "10104"
	}

	ok, err := redis.Bool(redisDB.Client.Hexists(db.TABLE_NAME_PRV_RANKING, userid))
	if err != nil{
		fmt.Println(" InitUserRankingData  TABLE_NAME_PRV_RANKING 1 ", err.Error())
	}
	if ok {
		fmt.Println(" InitUserRankingData  TABLE_NAME_PRV_RANKING 2 ", err.Error())
		return "10103"
	}else {
		fmt.Println(" InitUserRankingData  TABLE_NAME_PRV_RANKING 3 ")
		//初始化 TABLE_NAME_PRV_RANKING 排行榜的有序集合

		//TABLE_NAME_PRV_RANKING
		//初始化排行榜的有序集合 		//TABLE_NAME_PRV_RANKING
		_, err2 := redisDB.Client.Do("ZADD", db.TABLE_NAME_PRV_RANKING,  0 ,userid)
		if err2 != nil {
			return err2.Error()
		}
	}
	fmt.Println(" InitUserRankingData end ")
	return "0"
}
type PublicRankingList struct {
	ID              	primitive.ObjectID  `json:"id" bson:"_id"`
	DayRankingList		[]Ranking			`json:"day" bson:"day"`		//日排行榜
	WeekRankingList		[]Ranking			`json:"week" bson:"week"`	//周排行榜
	MonthRankingList	[]Ranking			`json:"month" bson:"month"`	//月排行榜
	RefreshDayTime		int				`json:"refreshday" bson:"refreshday"`
	RefreshWeekTime		int				`json:"refreshweek" bson:"refreshweek"`
	RefreshMonthTime	int				`json:"refreshmonth" bson:"refreshmonth"`
}

type Ranking struct {
	UserName 			string				`json:"username" bson:"username"`		//玩家ID
	UserId				string 				`json:"userId" bson:"userId"`
	RankingNo			int					`json:"rankingno" bson:"rankingno"`		//排名名次
	Head 				string				`json:"head" bson:"head"`			//头像
	Bonus 				float64				`json:"bonus" bson:"bonus"`			//获奖总额
}

type RedisRanking struct {
	UserID				int 				`json:"userId" bson:"userId"`
	Bonus 				float64				`json:"bonus" bson:"bonus"`			//获奖总额
}

//将前100名玩家的数据刷新到 总排行榜中
func RefreshRankingDataToRedis()(result string, error string) {
	scoreMap, err := redis.StringMap(redisDB.Client.Do("ZREVRANGE", db.TABLE_NAME_PRV_RANKING, 0, 99, "withscores"))
	if err != nil {
		return "", err.Error()
	}

	_result_list := []Ranking{}
	_no := 1
	for name := range scoreMap {
		//fmt.Println(name, scoreMap[name])

		userid := name
		if userid == "" {
			continue
		}

		_bons, _ := strconv.ParseFloat(scoreMap[name], 64)

		_user_info := new(Ranking)
		//user_info, err2 := model.GetNewUserInfoByRedis(userid)
		//if err2 != nil {
		//	fmt.Println("userid is not find :", userid, " err:", err2.Error())
		//	continue
		//}
		//_user_info.UserId = userid
		//_user_info.Head = user_info.Head
		//_user_info.UserName = user_info.UserName
		_user_info.Bonus = _bons
		//_user_info.RankingNo = _no

		_result_list = append(_result_list, *_user_info)

		_no += 1
	}
	//按照下注的多少，从大到小排列
	sort.Slice(_result_list, func(i, j int) bool {
		return _result_list[i].Bonus > _result_list[j].Bonus
	})

	for i:=0; i< len(_result_list);i++{
		_result_list[i].RankingNo = i+1
	}

	if len(_result_list) > 0{
		//从redis中找出表格
		data_v := new(PublicRankingList)
		//将byte流解析为json格式
		if data_v == nil{
			return "","0"
		}
		data_v.DayRankingList = _result_list
		data_v.RefreshDayTime = time.Now().Day()

		this_week := WeekByDate()
		if data_v.RefreshWeekTime != this_week {
			data_v.RefreshWeekTime = this_week
			data_v.WeekRankingList = _result_list
		}

		this_mouth := MonthByDate()
		if data_v.RefreshMonthTime != this_mouth {
			data_v.RefreshMonthTime = this_mouth
			data_v.MonthRankingList = _result_list
		}

		dataStr ,_ := jsoniter.Marshal(data_v)
		_, err2 := redisDB.Client.Do("HSET", db.TABLE_NAME_PUB_RANKING, db.TABLE_NAME_PUB_RANKING, dataStr)
		if err2 != nil{
			//log.Error("RefreshRankingDataToRedis error is ", err2.Error())
		}
	}

	return "", "0"
}

// 刷新玩家排行
// matchReward 玩家本次赢得的money
func SetRankingDataByUser(userid string, matchReward float64)(err error) {

	_, err = redisDB.Client.Do("ZINCRBY",db.TABLE_NAME_PRV_RANKING, matchReward, userid )
	if err != nil{
		return  err
	}

	//_rs.HighestReward += matchReward
	_index, err := redis.Int(redisDB.Client.Do("ZREVRANK",db.TABLE_NAME_PRV_RANKING, userid))
	if err != nil {
		return  err
	}
	fmt.Println("_index = ", _index)
	//更新自己的数据
	/*_user_stistics, err := redis.Bytes(redisDB.Client.Hget(REDIS_USER_DATA_PREFIX + userid, REDIS_USER_STATISTICS))
	if err != nil{
		return err
	}
	_rs := new(UserStatistics)
	err = jsoniter.Unmarshal(_user_stistics, _rs)

	if _rs.HighestRanking > _index {
		_rs.HighestRanking = _index
	}
	M,_:= jsoniter.Marshal(_rs)
	_, err = global.Client.Hset(REDIS_USER_DATA_PREFIX + userid, REDIS_USER_STATISTICS, M)
	if err != nil{
		return err
	}
*/
	return nil
}