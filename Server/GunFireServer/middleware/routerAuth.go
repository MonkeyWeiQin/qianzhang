package middleware

import (
	"com.xv.admin.server/model/db"
	"com.xv.admin.server/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)


func whitelistRoute() []string {
	return []string{
		"/admin/login",
		"/admin/info",
		"/system/downLoadGiftList",
		"/system/getLevelMenuList",
		"/system/createMenuData",
		"/admin/list",
		"/admin/updateRole",
		"/system/getAttachmentList",
		"/system/getAttachmentLabel",
	}
}

func RouterAuth() gin.HandlerFunc{
	return func(c *gin.Context) {
		currentRoute := c.Request.URL.Path
		if utils.InArray(currentRoute,whitelistRoute()) {
			c.Next()
			return
		}
		username, _ := c.Get("username")
		GetRoleMid,GetMenuMidErr := new(db.MenuModel).GetMenuMid(bson.D{{"path",currentRoute}})
		if GetMenuMidErr != nil {
			c.Abort()
			utils.Send(c, utils.UNAUTHORIZED, nil, GetMenuMidErr.Error())
			return
		}
		if GetRoleMid <= 0 {
			c.Abort()
			utils.Send(c, utils.UNAUTHORIZED, nil, "接口未找到")
			return
		}
		AdminRole,GetRoleErr := new(db.AdminModel).GetRole(bson.D{{"username",username}})
		if GetRoleErr != nil {
			c.Abort()
			utils.Send(c, utils.UNAUTHORIZED, nil, GetRoleErr.Error())
			return
		}
		if !utils.InArray(GetRoleMid,AdminRole) {
			c.Abort()
			utils.Send(c, utils.UNAUTHORIZED, nil, "权限不足")
			return
		}

		c.Next()
	}
}