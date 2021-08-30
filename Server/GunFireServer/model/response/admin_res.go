package response

import (
	"com.xv.admin.server/model/db"
)

// 登录响应数据结构
type LoginResponse struct {
	User      *db.AdminModel `json:"user"`
	Token     string         `json:"token"`
	ExpiresAt int64          `json:"expires_at"`
}

// 当前登录的管理员信息和权限角色数据响应数据结构
type AdminInfoResponse struct {
	Roles  []string `json:"roles"`
	Avatar string   `json:"avatar"`
	Name   string   `json:"name"`
}
