package lobby

import "go.mongodb.org/mongo-driver/bson/primitive"

//------------------------------------功能需求-----表格：pub_notice------------------------------------------//
//从redis/mongo公共数据库中，获得公告表内容，后台管理也就必须向后台redis/mongo输入公告内容，并由撤销功能
//公告必须有有效期
//公告：等级 内容 有效期
//通过 mgodb模块 先建立 数据结构  然后通过  对应数据结构的后台方法 建立数据库表格
//在通过所在模块的module模块进行初始化操作

//---------------------------------------公共数据处理的基本原则--------------------------------------------//
// 后台  通过 HTtp/接口 访问 redis/mongo 增删改
// 运行时，直接从redis/mongo中获得数据，
// 后台 做定时任务对过期数据及切换进行处理，写入redis/mongo

//游戏公告
//游戏公告
//游戏系统消息
type NoticeMsgData struct {
	ID       primitive.ObjectID 	`json:"id"  bson:"_id"`
	AdvID	 int	`json:"advid"     bson:"advid"`			//公告ID 手动操作时，避免重复添加
	Name     string  `json:"name"     bosn:"name"`        // 公告名称
	Notice   string  `json:"notice"   bosn:"notice"`      // 公告的字符内容
	NType    int     `json:"ntype"    bosn:"ntype"`	      // 公告类型
	STime    int64   `json:"stime"    bosn:"stime"`	      // 开始时间
	ETime    int64   `json:"etime"    bosn:"etime"`	      // 结束时间
	LinkURL  string  `json:"link"     bosn:"link"`        // 连接地址
}

//游戏系统消息
type NoticeMsgRequest struct {
	AdvID    int     `json:"advid"`     // 公告ID 手动操作时，避免重复添加
    Name     string  `json:"name"`         // 公告名称
	Notice   string  `json:"notice"`       // 公告的字符内容
    NType    int     `json:"ntype"`	       // 公告类型
	STime    int64   `json:"stime"`	       // 开始时间
    ETime    int64   `json:"etime"`	       // 结束时间
	LinkURL  string  `json:"link"`         // 连接地址
}

const (
	TYPE_PERMANENT  int = 1 //1 为永久公告
	TYPE_TIMELIMIT 	int = 2 // 2 为时限公告
)
