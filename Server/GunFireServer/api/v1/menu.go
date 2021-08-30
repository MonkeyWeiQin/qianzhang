package v1

import (
	"com.xv.admin.server/model/db"
	"com.xv.admin.server/model/request/admin_req"
	"com.xv.admin.server/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

/*
 * 获取菜单列表
 */
func GetMenuList(c *gin.Context) {
	var menu = new(db.MenuModel)
	menus, err := menu.SelectAllMenus(nil, nil)
	if err != nil {
		utils.FailWithMessage(c, "获取失败!")
		return
	}
	utils.OkWithData(c, menus)
	return
}

func GetLevelMenuList(c *gin.Context) {
	GetMenuListResponse := db.LevelMenu{PMid: 1, MId: 1, Name: "系统", Type: 1, Children: []*db.LevelMenu{}}
	GetLevelMenuListRequest := new(admin_req.GetLevelMenuListRequest)
	var menu = new(db.MenuModel)
	//mid , err := strconv.Atoi(GetMenuListResponse.MId)
	menus, err := menu.GetLevelMenuList(GetMenuListResponse.MId, GetLevelMenuListRequest.Type)
	if err != nil {
		fmt.Println(err)
		utils.FailWithMessage(c, "获取失败!")
		return
	}
	GetMenuListResponse.Children = menus

	utils.OkWithData(c, []db.LevelMenu{GetMenuListResponse})
	return
}

func CreateMenuData(c *gin.Context) {
	var menuReq admin_req.CreateMenuRequest

	err := c.ShouldBind(&menuReq)
	if err != nil {
		fmt.Println(err)
		utils.FailWithMessage(c, "参数异常")
		return
	}

	menu := &db.MenuModel{
		PMid:      menuReq.PMid,
		Type:      menuReq.Type,
		Name:      menuReq.Name,
		Path:      menuReq.Path,
		Component: menuReq.Component,
		//Title:     menuReq.Title,
		Icon: menuReq.Icon,
	}
	err = menu.InsertMenu()
	if err != nil {
		log.Print(err)
		utils.FailWithMessage(c, "添加失败!")
		return
	}
	utils.Ok(c)
}

func UpdateMenuData(c *gin.Context) {
	var menuReq admin_req.UpdateMenuRequest

	err := c.ShouldBind(&menuReq)
	if err != nil {
		fmt.Println(err)
		utils.FailWithMessage(c, "参数异常")
		return
	}

	menu := &db.MenuModel{
		PMid:      menuReq.PMid,
		Type:      menuReq.Type,
		Name:      menuReq.Name,
		Path:      menuReq.Path,
		Component: menuReq.Component,
		Icon:      menuReq.Icon,
		MId:       menuReq.Mid,
	}
	err = menu.UpdateMenuByMid()
	if err != nil {
		utils.FailWithMessage(c, "更新失败!")
		log.Println(err)
		return
	}
	utils.Ok(c)
}

func DeleteMenuData(c *gin.Context) {
	var req admin_req.DelMenuByMidRequest
	err := c.ShouldBind(&req)
	if err != nil {
		fmt.Println(err)
		utils.FailWithMessage(c, "参数异常")
		return
	}
	m := &db.MenuModel{}
	err = m.DelMenuById(req.Mid)
	if err != nil {
		log.Println(err.Error())
		utils.Fail(c)
		return
	}
	utils.Ok(c)
}
