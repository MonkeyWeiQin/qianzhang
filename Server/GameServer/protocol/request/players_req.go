package request

// DevId http 登录或者注册
type PlayerDevIdLoginReq struct {
	DevId string `form:"devId" binding:"required"`
}

// 手机账户登录或者注册
type PlayerPasswordLoginReq struct {
	Mobile   string `form:"mobile" binding:"required"`
	Password string `form:"password" binding:"required"`
}

// 手机验证码登录或者注册
type PlayerSmsCodeLoginReq struct {
	Mobile string `form:"mobile" binding:"required"`
	Code   string `form:"code" binding:"required"`
}

// TCP 绑定连接认证请求
type BindSessionReq struct {
	Token string `json:"token"` // 用户原始token
	Url   string `json:"url"`   // 校验地址,也就是客户端登录服务器的地址,分布式中保持校验和登录地址一直
}

type UpgradeRoleTalentLvReq struct {
	Ty int8 `form:"ty" binding:"required"` // 升级类型
}

type ChangeRoleSkinReq struct {
	SkinId string `json:"skinId"` // 皮肤的ID
}

// ty  1    1    1  max = 7
//    role main sub
type FlushUserAttributePushData struct {
	Ty int8 `json:"ty"` // 推送类别 二进制表示三个推送选择
}

type UpgradeRoleLevelByLvCardReq struct {
	Mid string `json:"mid"`
}
