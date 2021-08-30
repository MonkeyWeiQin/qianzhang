package v1

import (
	"com.xv.admin.server/config"
	"com.xv.admin.server/middleware"
	"com.xv.admin.server/model/db"
	"com.xv.admin.server/model/request/admin_req"
	"com.xv.admin.server/model/response"
	"com.xv.admin.server/utils"
	"crypto/md5"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// AdminLogin login
func AdminLogin(c *gin.Context) {
	var req admin_req.LoginRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		utils.FailWithMessage(c, "参数异常")
		return
	}
	// 加密Md5
	req.Password = createPassword(req.Password)
	var AdminModel = &db.AdminModel{}

	Admin, LoginByUserErr := AdminModel.LoginByUser(req.Username, req.Password)
	if LoginByUserErr != nil || Admin == nil {
		utils.FailWithMessage(c, "登录失败")
		return
	}
	if Admin.Status != 0 {
		utils.FailWithMessage(c, "账号无效")
	} else {
		createToken(c, Admin)
	}
}

func createPassword(password string) string {
	data := []byte(password + db.MD5Secret)
	has := md5.Sum(data)
	return fmt.Sprintf("%x", has) //将[]byte转成16进制
}

// 登录以后签发jwt
func createToken(c *gin.Context, user *db.AdminModel) {
	j := &middleware.JWT{
		SigningKey: []byte(config.ENV.JwtSecret), // 唯一签名
	}
	clams := jwt.StandardClaims{
		NotBefore: time.Now().Unix() - 10,               // 签名生效时间
		ExpiresAt: time.Now().Add(time.Hour * 3).Unix(), // 过期时间3小时
		Issuer:    "qmPlus",                             // 签名的发行者
		Id:        user.Username,
	}
	token, err := j.CreateToken(clams)
	if err != nil {
		utils.FailWithMessage(c, "获取token失败")
		return
	}
	updateLoginTime, UpdateErr := user.Update(bson.D{{"username", user.Username}}, bson.D{{"last_login_time", time.Now().Unix()}})
	if UpdateErr != nil {
		utils.FailWithMessage(c, "更新登录时间失败："+UpdateErr.Error())
		return
	}
	if updateLoginTime < 1 {
		utils.FailWithMessage(c, "更新登录时间失败")
		return
	}
	utils.OkWithData(c, &response.LoginResponse{
		User:      user,
		Token:     token,
		ExpiresAt: clams.ExpiresAt,
	})
}

func GetAdminInfo(c *gin.Context) {
	utils.OkWithData(c, response.AdminInfoResponse{
		Roles:  []string{"admin"},
		Avatar: "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif",
		Name:   "超级管理员",
	})
}

func CreateAdmin(c *gin.Context) {
	req := new(admin_req.CreateRequest)
	err := c.ShouldBindJSON(req)
	if err != nil {
		utils.FailWithMessage(c, "参数错误"+err.Error())
		return
	}
	var admin = &db.AdminModel{}
	num, GetTotalNumErr := admin.GetTotalNum(bson.D{{"username", req.Username}})
	if GetTotalNumErr != nil {
		utils.FailWithMessage(c, "创建管理员失败"+GetTotalNumErr.Error())
		return
	}
	if num >= 1 {
		utils.FailWithMessage(c, "创建管理员失败:账号重复")
		return
	}

	admin.Username = req.Username
	admin.Avatar = "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif"
	admin.Password = createPassword(req.Password)
	if CreateErr := admin.Create(); CreateErr != nil {
		utils.FailWithMessage(c, "创建管理员失败"+CreateErr.Error())
		return
	}
	utils.Ok(c)
	return
}

func GetAdminList(c *gin.Context) {
	req := new(admin_req.GetAdminListRequest)
	err := c.ShouldBind(req)
	if err != nil {
		utils.FailWithMessage(c, "参数错误"+err.Error())
		return
	}

	if req.Limit == 0 {
		req.Limit = 10
	}
	if req.Page == 0 {
		req.Page = 1
	}
	model := &db.AdminModel{}
	var filter bson.D

	skip := (req.Page - 1) * req.Limit
	opt := &options.FindOptions{
		Limit: &req.Limit,
		Skip:  &skip,
		Sort:  bson.D{{"uid", -1}},
	}
	datas, err := model.GetList(filter, opt)
	if err != nil {
		utils.FailWithMessage(c, "参数错误"+err.Error())
		return
	}
	count, err := model.GetTotalNum(filter)
	utils.OkWithData(c, map[string]interface{}{
		"list":  datas,
		"page":  req.Page,
		"limit": req.Limit,
		"total": count,
	})
}

func UpdatePassword(c *gin.Context) {
	req := new(admin_req.CreateRequest)
	err := c.ShouldBindJSON(req)
	if err != nil {
		utils.FailWithMessage(c, "参数错误"+err.Error())
		return
	}
	model := &db.AdminModel{}
	if len(req.Password) < 5 {
		utils.FailWithMessage(c, "修改失败，密码最小长度为5位")
		return
	}
	res, UpdateErr := model.Update(bson.D{{"username", req.Username}}, bson.D{{"password", req.Password}})
	if UpdateErr != nil {
		utils.FailWithMessage(c, "修改失败:"+UpdateErr.Error())
		return
	}
	if res > 0 {
		utils.Ok(c)
		return
	}
	utils.FailWithMessage(c, "更新失败")
	return
}

func UpdateRole(c *gin.Context) {
	req := new(admin_req.UpdateRole)
	err := c.ShouldBindJSON(req)
	if err != nil {
		utils.FailWithMessage(c, "参数错误"+err.Error())
		return
	}

	model := &db.AdminModel{}
	menuModel := &db.MenuModel{}
	_, GetMenuNameByMidErr := menuModel.GetMenuNameByMids(req.Role)
	if GetMenuNameByMidErr != nil {
		utils.FailWithMessage(c, "修改失败:"+GetMenuNameByMidErr.Error())
		return
	}
	if req.Role[0] ==  1 {
		req.Role = append(req.Role[:0], req.Role[1:]...)
	}
	res, UpdateErr := model.Update(bson.D{{"username", req.Username}}, bson.D{{"role", req.Role}})
	if UpdateErr != nil {
		utils.FailWithMessage(c, "修改失败:"+UpdateErr.Error())
		return
	}
	if res > 0 {
		utils.Ok(c)
		return
	}else{
		utils.FailWithMessage(c, "更新失败:数据未发生改变" )
		return
	}
}
