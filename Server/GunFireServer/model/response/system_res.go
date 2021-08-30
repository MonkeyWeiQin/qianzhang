package response


// 当前登录的管理员信息和权限角色数据响应数据结构
type CloseGameInfoResponse struct {
	State int      `json:"state"`
	Error string   `json:"error"`
}

type GetMenuListResponse struct {
	Value int `json:"value"`
	Label string `json:"label"`
	Children []GetMenuListResponse `json:"children"`
}





