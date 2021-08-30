package response

// 注册留存率响应
type RegisterRetentionRateResponse struct {
	Time        int64 `json:"time"`
	RegisterNum int    `json:"register_num"`
	One         int    `json:"one"`
	Two         int    `json:"two"`
	Three       int    `json:"three"`
	Four        int    `json:"four"`
	Five        int    `json:"five"`
	Six         int    `json:"six"`
	Seven       int    `json:"seven"`
	Fifteen     int    `json:"fifteen"`
	Month       int    `json:"month"`
}

// UserLoginResponse  登录返回参数
type UserLoginResponse struct {
	Token 		string 	`json:"token"`
	Username 	string 	`json:"username"`
	UserId      int     `json:"uid"`
}