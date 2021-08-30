package admin_req

/**
 * 创建新菜单请求数据
 */
type CreateMenuRequest struct {
	PMid      int    `form:"pmid" `                   // 上级目录的ID  1为顶级目录
	Type      int8   `form:"type" binding:"required"` // 1  菜单 2 路由
	Name      string `form:"name" binding:"required"` // 前端vue组件名称
	Path      string `form:"path" binding:"required"` // 前端URL栏 /#/xxx 或 后台接口路由
	Component string `form:"component"`               // 前端vue组件componentMap中定义的组件key,前端通过key加载指定的vue组件
	//Title     string `form:"title" binding:"required"`  // 菜单标题
	Icon string `form:"icon" ` // 图标
}
type UpdateMenuRequest struct {
	Mid       int    `form:"mid" binging:"required"`
	PMid      int    `form:"pmid" `                   // 上级目录的ID  1为顶级目录
	Type      int8   `form:"type" binding:"required"` // 1  菜单 2 路由
	Name      string `form:"name" binding:"required"` // 前端vue组件名称
	Path      string `form:"path" binding:"required"` // 前端URL栏 /#/xxx 或 后台接口路由
	Component string `form:"component"`               // 前端vue组件componentMap中定义的组件key,前端通过key加载指定的vue组件
	//Title     string `form:"title" binding:"required"`  // 菜单标题
	Icon string `form:"icon" ` // 图标
}
type DelMenuByMidRequest struct {
	Mid int `form:"mid"`
}

type GetLevelMenuListRequest struct {
	Type int `form:"type"` //1只获取目录 0获取全部
}
