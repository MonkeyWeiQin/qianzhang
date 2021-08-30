package user_req

type GetUsersListRequest struct {
	Uid               string `form:"uid"`
	Page              int64  `form:"page"`
	Limit             int64  `form:"limit"`
	Username          string `form:"username"`
	RegisterEndTime   int64  `form:"register_end_time"`
	RegisterStartTime int64  `form:"register_start_time"`
}

// ModifyNumberRequest 修改钻石/体力/金币
type ModifyNumberRequest struct {
	Uid    int    `form:"uid"`
	Type   string `form:"type"`
	Number int64  `form:"number"`
}

// ModifyStatusRequest 修改当前状态
type ModifyStatusRequest struct {
	Uid    int   `form:"uid"`
	Status int   `form:"status"`
	Time   int64 `form:"time"`
}
