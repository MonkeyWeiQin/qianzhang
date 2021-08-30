package router

import (
	"com.xv.admin.server/api/v1"
	"com.xv.admin.server/middleware"
	"github.com/gin-gonic/gin"
)

func InitSystemRouter(gin *gin.Engine) {
	// 简单的路由组: v1
	router := gin.Group("/system").
		Use(middleware.JwtAuth()).Use(middleware.RouterAuth())
	{
		router.GET("/getMenuList", v1.GetMenuList)           //获取菜单列表
		router.POST("/createMenuData", v1.CreateMenuData)    //创建菜单
		router.POST("/updateMenuData", v1.UpdateMenuData)    //更新菜单
		router.GET("/getLevelMenuList", v1.GetLevelMenuList) //获取菜单列表--子父级架构
		router.POST("/deleteMenuData", v1.DeleteMenuData)    //删除菜单

		router.GET("/getNotice", v1.GetNoticeList)
		router.POST("/addNotice", v1.CreateNotice)
		router.POST("/updateNotice", v1.UpdateNotice)
		//router.POST("/delNotice", v1.DelPublicNotice)

		router.POST("/createSystemMail", v1.CreateSystemMail) //新增系统邮件
		router.GET("/getSystemMail", v1.GetSystemMail)        //获取系统邮件列表
		router.POST("/updateSystemMail", v1.UpdateSystemMail) //修改系统邮件
		router.GET("/downLoadGiftList", v1.DownLoadGiftCode)  //导出兑换码

		router.POST("/createGift", v1.CreateGift)        //新增兑换码
		router.GET("/getGiftList", v1.GetGiftList)       //获取兑换码
		router.GET("/getGiftLogList", v1.GetGiftLogList) //获取兑换码兑换列表

		router.GET("/getBroadcast", v1.GetBroadcastList)    //获取广播
		router.POST("/createBroadcast", v1.CreateBroadcast) //创建广播
		router.POST("/updateBroadcast", v1.UpdateBroadcast) //更新广播

		router.GET("/getGoodsList", v1.GetGoodsList)                 //获取商品列表
		router.GET("/getGoodsPurchaseList", v1.GetGoodsPurchaseList) //获取购买记录

		router.GET("/getChestList", v1.GetChestList)                 //获取宝箱列表
		router.GET("/getChestPurchaseList", v1.GetChestPurchaseList) //获取购买宝箱记录
		router.GET("/getChestContent", v1.GetChestContent)           //获取宝箱内容
	}

	attachmentRouter := gin.Group("/system")
	{
		attachmentRouter.GET("/getAttachmentLabel", v1.GetAttachmentLabel) //获取奖品配置
		attachmentRouter.GET("/getAttachmentList", v1.GetAttachmentList)   //获取奖品配置
	}

}
