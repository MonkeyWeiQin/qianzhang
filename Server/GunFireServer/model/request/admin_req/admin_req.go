package admin_req

// LoginRequest 管理员登录请求参数
type LoginRequest struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

// CreateRequest 管理员登录请求参数
type CreateRequest struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

// GetAdminListRequest 管理员登录请求参数
type GetAdminListRequest struct {
	Uid   string `form:"uid"`
	Page  int64  `form:"page"`
	Limit int64  `form:"limit"`
}

// UpdateRole 修改权限
type UpdateRole struct {
	Username string `form:"username" binding:"required"`
	Role     []int  `form:"role" binding:"required"`
}
