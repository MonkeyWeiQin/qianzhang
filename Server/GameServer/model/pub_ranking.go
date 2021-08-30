package model

import "go.mongodb.org/mongo-driver/bson/primitive"

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

