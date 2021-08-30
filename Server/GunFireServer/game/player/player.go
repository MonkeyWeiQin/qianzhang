package player

import (
	"fmt"
	"github.com/gin-gonic/gin"
)



type UserInfo struct {
	Mobile         string  `json:"mobile"           bson:"mobile"`         // 手机号
	AccountId      string  `json:"accountId"        bson:"accountId"`      // 用户唯一的ID
	DevID          string  `json:"devId"            bson:"devId"`          // 用户登陆时的硬件ID
	UserId         string  `json:"userId"           bson:"userId"`         // 用户唯一的ID
	Status         int     `json:"status"           bson:"status"`         // 账户状态 1正常  0 禁用 -1 临时冻结
	Username       string  `json:"username"         bson:"username"`       // 用户名
	RegisterTime   int64   `json:"registertime"     bson:"registertime"`   // 注册时间
	LoginTime      int64   `json:"logintime"        bson:"logintime"`      // 登陆时间
	Head           string  `json:"head"             bson:"head"`           // 用户头像
	Sicon          string  `json:"sicon"            bson:"sicon"`          // 头像框
	Sex            int     `json:"sex"              bson:"sex"`            // 性别 1 男 2 女
	Birthday       string  `json:"birthday"         bson:"birthday"`       // 生日 "1980-02-21"
	Email          string  `json:"email"            bson:"email"`          // 邮箱
	EmailCode      string  `json:"emailcode"        bson:"emailcode"`      // 邮箱验证码
	Password       string  `json:"password"         bson:"password"`       // 登陆密码
	Level          int     `json:"level"            bson:"level"`          // 账号等级
	Exp            int     `json:"exp"              bson:"exp"`            // 经验值
	Money          float64 `json:"money"            bson:"money"`          // 钻石
	Gold           float64 `json:"gold"             bson:"gold"`           // 金币
	ProhibitedTime int     `json:"prohibitedTime"   bson:"prohibitedTime"` // status 状态为冻结时的解冻日期时间点(unix)
	Strength       int     `json:"strength"         bson:"strength"`       //Physical strength, 体力，
	VIP            int     `json:"vip"              bson:"vip"`            //VIP level,vip等级，
	CurStageNo     int     `json:"stageNo"          bson:"stageNo"`        //The current level,当前关卡，

}


// 不需要存储，临时计算结果设置标志
// 标志信息：
//flag information:
type Flag_Info struct {
	Mall_Flag         int `json:"mallflag"`      //Mall flag, 商城标志，
	Character_Flag    int `json:"characterflag"` //Character flag, 角色标志，
	Hero_Flag         int `json:"heroflag"`      //The hero flag, 英雄标志，
	Camp_Flag         int `json:"campflag"`      //Camp flag, 基地标志，
	Recharge_Flag     int `json:"rechargeflag"`  //Recharge flag, 充值标志，
	Mail_Flag         int `json:"mailflag"`      //Mail flag, 邮件标志，
}


// 不需要存储，临时计算结果设置标志
//活动信息： Activity information:
type Activity_Info struct {
	Season_Flag      int   `json:"seasonflag"`      //Season marks, 赛季标志，
	Skin_Flag        int   `json:"skinflag"`        //The main product is the changing (skin) logo, 主推换装（皮肤）标志，
	Special_1_Flag   int   `json:"special1flag"`    //Special event 1, 特惠活动1，
	Special_2_Flag   int   `json:"special2flag"`    //Special offer 2 特惠活动2

}

// 登陆后首页信息内容：
// 基本信息：体力，金币，钻石，名字，等级，经验，头像，头像框，vip等级，当前关卡，当前装备（枪，副武器，其他外挂），显示宠物/英雄
// 标志信息：商城标志，角色标志，英雄标志，基地标志，充值标志，邮件标志，
// 活动信息：赛季标志，主推换装（皮肤）标志，特惠活动1，特惠活动2
func GetPlayerBaseInfo(c *gin.Context) {
	fmt.Println("GetPlayerBaseInfo")
}
//-----------------------------------天赋系统操作-----------------------------------------------
// 初始化时，可将天赋技能所有的效果都初始化进去，默认等级0，就相当于没有
type PassiveSkillCSVTableInfo struct {
	Passive_ID     string `csv:"passive_id"`       //天赋技能ID
	Passive_Name   string `csv:"passive_name"`     //天赋技能名称
	Pssive_LV      int    `csv:"passive_lv"`       //天赋技能等级
	Passive_Icon   string `csv:"passive_icon"`     //天赋技能图标
	Passive_Effect string `csv:"passive_effect"`   //天赋技能效果
	Passive_Data   int    `csv:"passive_data"`     //天赋技能数据
	Lvup_Iiamond   int64  `csv:"lvup_diamond"`     //直接升级需求钻石
}

type PassiveSkillTableData struct {
	Passive_ID     string `json:"passive_id"      bson:"passive_id"`         //天赋技能ID
	Passive_Name   string `json:"passive_name"    bson:"passive_name"`       //天赋技能名称
	Pssive_LV      int    `json:"passive_lv"      bson:"passive_lv"`         //天赋技能等级
	Passive_Icon   string `json:"passive_icon"    bson:"passive_icon"`       //天赋技能图标
	Passive_Effect string `json:"passive_effect"  bson:"passive_effect"`     //天赋技能效果
	Passive_Data   int    `json:"passive_data"    bson:"passive_data"`       //天赋技能数据
	Lvup_Iiamond   int64  `json:"lvup_diamond"    bson:"lvup_diamond"`       //直接升级需求钻石
}

type PassiveSkillLevelUpCSVInfo struct {
	Total_LV    int      `csv:"total_lv"`    //天赋总等级
	Lvup_Gold   int64    `csv:"lvup_gold"`   //升级天赋消耗金币
}

type PassiveSkillLevelUpTableData struct {
	Total_LV    int      `json:"total_lv"      bson:"total_lv"`    //天赋总等级
	Lvup_Gold   int64    `json:"lvup_gold"     bson:"lvup_gold"`   //升级天赋消耗金币
}

// 金币升级为随机选中，满级时不在可选择列表中；钻石升级是选中指定天赋进行升级

// 天赋技能的存储数据结构
type PassiveSkillInfo struct {
	SkillID   string  `json:"skillid"    bson:"skillid"`
	Level     int     `json:"level"      bson:"level"`    // 当前等级
}

// 天赋技能的存储数据结构
type UserPassiveSkillData struct {
	UserID           string                       `json:"userid"       bson:"userid"`
	TalentSkillList  map[string]PassiveSkillInfo  `json:"tskilllist"    bson:"tskilllist"`
}
// 玩家的天赋系统数据
func GetPlayerPassiveInfo(c *gin.Context) {
	fmt.Println("GetPlayerTalentInfo")
}
// 天赋技能升级的请求参数数据结构
type UpdatePlayerPassiveLevelRequest struct {
	TType   int         `json:"ttype"`  // 0 为金币 1 为钻石类型
	SkillID string 		`json:"skillid"` // 为钻石类型时，此参数有效
}
// 升级玩家的天赋
func UpdatePlayerPassiveLevel(c *gin.Context) {
	// 区别升级的类型，是金币升级还是钻石升级，客户端要有选择过程，服务器确认后直接返回，客户端延迟显示结果，中间播放选择过程
	//
	fmt.Println("UpdatePlayerTalentLevel")
}
//-----------------------------------武器系统操作-----------------------------------------------
// 玩家武器系统数据
func GetPlayerArmsInfo(c *gin.Context) {
	fmt.Println("GetPlayerArmsInfo")
}
// 设置玩家主武器
func SetPlayerMainWeapon(c *gin.Context) {
	fmt.Println("SetPlayerMainWeapon")
}

// 设置玩家副武器
func SetPlayerAuxWeapon(c *gin.Context) {
	fmt.Println("SetPlayerAuxWeapon")
}


//-----------------------------------关卡系统操作-----------------------------------------------
//关卡存储数据结构
type StageBaseInfo struct {
	Stage_ID      string      //
	Stage_Type    int         // 1 主线关卡，2 挑战关卡，3 防守关卡，4金币关卡四种。
	Stage_State   int         // 0 未通关 1 通关
	Reward_Item   string      // 通关额外奖励，服务器调整
	Recive_State  int         // 领取状态 0 未领取  1 领取
}




// 设置通过主线关卡的索引
func SetPlayerMainStageNob(c *gin.Context){

}

// 获得玩家的关卡列表 需要记录每个关卡过关的成绩等信息
// 需要根据关卡类型来
//战车金币关卡
//坦克材料关卡
//生存模式关卡
//塔防模式关卡

func GetPlayerActivityStageList(c *gin.Context){

}
//-----------------------------------体力相关操作-----------------------------------------------
// 消耗玩家体力，同时开始计时进行体力恢复，就是把消耗体力的时间计下来
func PlayerConsumeStrength(c *gin.Context){

}

// 获得玩家体力，自动恢复体力和获得体力
func PlayerPickUpStrength(c *gin.Context){

}
//-----------------------------------货币相关操作-----------------------------------------------
//金币系统
func PlayerConsumeGold(c *gin.Context){

}

//钻石系统Diamonds   必然是购买物品
func PlayerConsumeDiamonds(c *gin.Context){

}

//-----------------------------------玩家可修改属性相关操作-----------------------------------------------
// 更换皮肤
func PlayerChangeSkin(c *gin.Context){

}
// 更换头像
func PlayerChangeHead(c *gin.Context){

}
// 更换名字
func PlayerChangeName(c *gin.Context){

}

//复活Reincarnate
//冲撞车Ram
//
//步枪/来复枪Rifle
//
//突击步枪Assault Rifle
//
//机枪Machine Gun/Pillbox
//
//冲锋枪Sub-Machine Gun
//
//狙击步枪Sniper Rifle
//
//火焰喷射器Flame Thrower
//
//离子Ion
//
//等离子Plasma
//
//激光Laser
//
//榴弹/手榴弹Grenade
//
//炸弹Bomb
//
//照明弹Flare
//
//发射器Launcher
//
//火箭rocket
//
//六管机枪Minigun/Vulcano
//
//散弹枪Shotgun
//
//加农炮Cannon
//
//榴弹炮Howitzer
//
//闪光弹Flash Grenade
//
//穿甲弹AP
//
//高爆弹HE
//
//防空炮AA-Guns
//
//高斯步枪Gauss
//
//电池Battery
//装甲Armor
//
//医疗箱Med Kit
//
//弹药Ammo
//
//弹夹Clip
//
//散弹枪子弹/炮弹Shell
//
//手电筒Flash Light
//雷达Radar
//背包Backpack