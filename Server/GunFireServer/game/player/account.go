package player

import (
	"com.xv.admin.server/model/db"
	"com.xv.admin.server/service/redisDB"
	"com.xv.admin.server/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/json-iterator/go"
	"log"
	"math/rand"
	"strconv"
	"time"
)

//新建角色的游客登陆及普通登陆
type AppLoginInRequest struct {
	DevID    string `form:"devid" binding:"required"`
	Password string `form:"password" binding:"required"`
	UserId   string `form:"userid" binding:"required"`
}

// 玩家登陆
func PlayerLogin(c *gin.Context) {
	fmt.Println("PlayerLogin")
	//db.REDIS_USER_INDEX_MAP
	req := new(AppLoginInRequest)
	err := c.ShouldBind(req)
	if err != nil {
		utils.FailWithMessage(c, "参数错误")
		return
	}
	userId := req.UserId
	devId := req.DevID
	passWord := req.Password
	fmt.Println("UserId = ", userId, " DevId = ", devId, " PassWord = ", passWord)

	ok, _ := redis.Bool(redisDB.Client.Hexists(db.REDIS_USER_INDEX_MAP, devId))
	if err != nil {
		log.Println("试图从Redis获取用户ID 失败！！：", err.Error())
		return
	}
	if ok{

	}else{
		newUserRegister(1, devId, passWord)
	}
}
func newUserRegister(action int, devid string, password string ){

	switch action{
	case 1 :
		break
	case 2 :
		break
	default:
		break
	}
	userId, err1 := redis.String(redisDB.Client.Hget(db.REDIS_INCREMENT_KEY_MAP, db.REDIS_USER_INCREMENT_KEY))
	if err1 != nil {
		return
	}

	userInfo := new(UserInfo)
	userInfo.Mobile         = "" // 手机号
	userInfo.AccountId      = "" // 用户唯一的ID
	userInfo.DevID          = devid // 用户登陆时的硬件ID
	userInfo.UserId         = userId // 用户唯一的ID
	userInfo.Status         = 1 // 账户状态 1正常  0 禁用 -1 临时冻结
	userInfo.Username       = "User_" + devid[:4] + "***" + devid[len(devid)-5:] // 用户名
	userInfo.RegisterTime   = time.Now().Unix() // 注册时间
	userInfo.LoginTime      = time.Now().Unix() // 登陆时间
	userInfo.Head           = "hot_res/Lobby/sprite/headIcon/image_touxiang_0" + strconv.Itoa(int(rand.Intn(9))) // 用户头像
	userInfo.Sicon          = "" // 头像框
	userInfo.Sex            = 2 // 性别 1 男 2 女
	userInfo.Birthday       = "" // 生日 "1980-02-21"
	userInfo.Email          = "" // 邮箱
	userInfo.EmailCode      = "" // 邮箱验证码
	userInfo.Password       = password // 登陆密码
	userInfo.Level          = 1 // 账号等级
	userInfo.Exp            = 0 // 经验值
	userInfo.Money          = 0 // 钻石
	userInfo.Gold           = 3000 // 金币
	userInfo.ProhibitedTime = 0 // status 状态为冻结时的解冻日期时间点(unix)
	userInfo.Strength       = 20 //Physical strength, 体力，
	userInfo.VIP            = 0 //VIP level,vip等级，
	userInfo.CurStageNo     = 1 //The current level,当前关卡，

	M, _ := jsoniter.Marshal(userInfo)
	_, _ = redisDB.Client.Hset(db.REDIS_USER_DATA_PREFIX+userId, db.REDIS_USER_INFO, M)

	_, err3 := redisDB.Client.Hset(db.REDIS_USER_INDEX_MAP, devid, userId)
	if err3 != nil {
		log.Println(err3.Error())
		return
	}

	// 更新user唯一索引
	_, err2 := redisDB.Client.Hincrby(db.REDIS_INCREMENT_KEY_MAP, db.REDIS_USER_INCREMENT_KEY, 1)
	if err2 != nil {
		fmt.Println("  更新user唯一索引 error = ", err2.Error())
	}

}
// 玩家注册
func PlayerRegister(c *gin.Context) {
	fmt.Println("PlayerRegister")
}