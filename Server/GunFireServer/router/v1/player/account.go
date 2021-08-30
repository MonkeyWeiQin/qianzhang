package player

import (
	"github.com/gin-gonic/gin"
)

func InitAccountRouter(gin *gin.Engine) {
	//router := gin.Group("v1")
	//{
	//	router.POST("/SendRegisterSms", player.SendRegisterSms)       //发送注册短信验证码
	//	router.POST("/VailRegisterSms", player.VailRegisterSms)       //验证注册短信验证码
	//	router.POST("/UserRegister", player.UserRegister)             //用户注册
	//	router.GET("/GetPlatformList", player.GetPlatformList)        //获取第三方登录平台
	//	router.POST("/OtherPlatformLogin", player.OtherPlatformLogin) //第三方登录
	//	router.POST("/Login", player.Login)                           //登录
	//
	//	router.POST("/notice", public.CreateNotice)   //创建公告
	//	router.GET("/notice", public.GetNotice)       //获取公告
	//	router.GET("/AllNotice", public.GetAllNotice) //获取所有公告
	//	router.PUT("/notice", public.UpdateNotice)    //更新公告
	//
	//	router.Any("/CheckToken", player.CheckToken) //token验证
	//}
	//
	//userRouter := gin.Group("/v1/User").Use(middleware.JwtAuth())
	//{
	//	userRouter.PUT("/", player.UpdateUser)                          //修改用户信息
	//	userRouter.GET("/", player.GetUser)                             //获取用户信息
	//	userRouter.POST("/BindMobile", player.BindMobile)               //绑定手机号
	//	userRouter.POST("/SendBindMobileSms", player.SendBindMobileSms) //发送绑定手机验证码
	//
	//	userRouter.POST("/mail", public.CreatePlayerMail) //创建邮件
	//	userRouter.GET("/mail", public.GetPlayerMailList) //获取邮件列表
	//	userRouter.PUT("/mail", public.SetPlayerMailRead) //设置邮件已读
	//	userRouter.DELETE("/mail", public.DelPlayerMail)  //设置邮件删除
	//
	//}
}
