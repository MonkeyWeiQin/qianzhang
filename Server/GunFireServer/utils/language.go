package utils

type En struct{}
type Zh struct{}
type Temp struct {
}

func (*En) message(temp interface{}) string {
	mess := map[interface{}]string{
		00000: "successful operation",
		10001: "parameter error",
		10002: "operation failed",
		10004: "The registration verification code is: [%d], the expiration time is 5 minutes",
		10005: "The phone is being bound, the verification code is:【%d】，the expiration time is 5 minutes",
		10006: "The binding phone is being modified, the verification code is: [%d], the expiration time is 5 minutes",
		10007: "Binding failed: the user information was not queried",
		10008: "Binding failed: you have already bound your phone",
		10009: "Registration failed",
		10010: "Failed to obtain token",
		10011: "Sending failed, the phone number has been registered",
		10012: "Failed to send",
		10013: "Binding failed: the user information was not find",
		10014: "Verification failed",
		10015: "The phone did not send the verification code or the verification code has expired",
		10016: "Verification code error",
		10017: "Binding failed: the phone number has been registered",
		10018: "Binding failed",
		10019: "Failed to obtain user information",
		10020: "Please enter account or password",
		10021: "Incorrect username or password",
		10022: "Login failed",
		10023: "The phone number is not registered",
		10024: "Update failed",
		10025: "Update completed",
		10026: "failed to delete",
		10027: "successfully deleted",
		10028: "The announcement ID entered already exists",
		10029: "Data does not exist",
	}
	return mess[temp]
}
func (*Zh) message(temp interface{}) string {
	mess := map[interface{}]string{
		00000: "操作成功",
		10001: "参数错误",
		10002: "操作失败",
		10004: "注册验证码为:【%d】，过期时间为5分钟",
		10005: "正在绑定手机,验证码为:【%d】，过期时间为5分钟",
		10006: "正在修改绑定手机,验证码为:【%d】，过期时间为5分钟",
		10007: "绑定失败：未查询到该用户信息",
		10008: "绑定失败：你已经绑定过手机",
		10009: "注册失败",
		10010: "获取token失败",
		10011: "发送失败，该手机号已被注册",
		10012: "发送失败",
		10013: "绑定失败：未查询到该用户信息",
		10014: "验证失败",
		10015: "该手机未发送验证码或验证码已失效",
		10016: "验证码错误",
		10017: "绑定失败：该手机号已被注册",
		10018: "绑定失败",
		10019: "获取用户信息失败",
		10020: "请输入账号或密码",
		10021: "账号或密码错误",
		10022: "登录失败",
		10023: "该手机号没有注册",
		10024: "更新失败",
		10025: "更新成功",
		10026: "删除失败",
		10027: "删除成功",
		10028: "输入的公告ID已存在",
		10029: "数据不存在",
	}
	return mess[temp]
}

type LanguageType struct {
	strategy LanguageTypeStrategy
}

func language(strategy LanguageTypeStrategy) *LanguageType {
	return &LanguageType{
		strategy: strategy,
	}
}

type LanguageTypeStrategy interface {
	message(interface{}) string
}
