package utils

import (
	"com.xv.admin.server/service/redisDB"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"math/rand"
	"strings"
	"time"
)

const REGISTER_SMS = "user_register_"           // 用户注册验证码
const BIND_MOBILE_SMS = "user_bind_mobile_"     // 用户绑定手机验证码
const UPDATE_MOBILE_SMS = "user_update_mobile_" // 用户修改手机号验证码

// SendRegisterMsg 发送注册短信
func SendRegisterMsg(mobile string) bool {
	code := SendMsg(GetMessage(10004), mobile)
	redisDB.Client.Setex(REGISTER_SMS+mobile, code, 5*60)
	return true
}

// SendBindMobileMsg 发送绑定手机短信
func SendBindMobileMsg(mobile string) bool {
	code := SendMsg(GetMessage(10005), mobile)
	redisDB.Client.Setex(BIND_MOBILE_SMS+mobile, code, 5*60)
	return true
}

// SendUpdateMobileMsg 发送修改手机短信
func SendUpdateMobileMsg(mobile string) bool {
	code := SendMsg(GetMessage(10006), mobile)
	redisDB.Client.Setex(UPDATE_MOBILE_SMS+mobile, code, 5*60)
	return true
}

// VailSms 验证注册短信验证码是否正确
func VailSms(mobile string, code string, temp string) {
	if !InArray(temp, GetAllSms()) {
		panic(10014)
	}
	sendCode, _ := redis.String(redisDB.Client.GET(temp + mobile))
	if sendCode == "" {
		panic(10015)
	}
	if sendCode != code {
		panic(10016)
	}
}

// GetAllSms 获取所有的短信模板标识
func GetAllSms() []string {
	return []string{
		REGISTER_SMS,
		BIND_MOBILE_SMS,
		UPDATE_MOBILE_SMS,
	}
}

// GetCode 获取验证码
func GetCode() string {
	return "1234" //todo 检查是否是调试模式
	return fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000))
}

// SendMsg 发送短信
func SendMsg(text string, mobile string) string {
	code := GetCode()
	fmt.Println(strings.Replace(text, "%d", code, -1))
	return code
}
