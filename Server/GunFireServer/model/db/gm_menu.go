package db

import (
	"com.xv.admin.server/config"
	"com.xv.admin.server/service/mgoDB"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/***********************************************
		菜单管理
 ***********************************************/
type MenuModel struct {
	Id        primitive.ObjectID `json:"-"bson:"-"`
	MId       int                `json:"mid" bson:"mid"`             // 唯一的ID
	PMid      int                `json:"pmid" bson:"pmid"`           // 上级目录的ID  1为顶级目录( TMD0在前端有回显问题)
	Type      int8               `json:"type" bson:"type"`           // 1  菜单 2 接口
	Path      string             `json:"path" bson:"path"`           // 前端URL栏 /#/xxx
	Component string             `json:"component" bson:"component"` // 前端vue组件componentMap中定义的组件key,前端通过key加载指定的vue组件
	Name      string             `json:"name" bson:"name"`           // 设定路由的名字，一定要填写不然使用<keep-alive>时会出现各种问题,最好与加载的组件名称(name)一样,并保持唯一
	//Title     string             `json:"title" bson:"title"`         // 设置该路由在侧边栏和面包屑中展示的名字
	Icon   string `json:"icon" bson:"icon"`     // 设置该路由的图标，支持 svg-class，也支持 el-icon-x element-ui 的 icon
	Status bool   `json:"status" bson:"status"` // 是否禁用或者是否隐藏
}

type LevelMenu struct {
	MId      int          `json:"mid"`              // 唯一的ID
	PMid     int          `json:"pmid"`             // 上级目录的ID  1为顶级目录( TMD0在前端有回显问题)
	Name     string       `json:"name" bson:"name"` // 设定路由的名字，一定要填写不然使用<keep-alive>时会出现各种问题,最好与加载的组件名称(name)一样,并保持唯一
	Type     int8         `json:"type" bson:"type"` // 1  菜单 2 接口
	Children []*LevelMenu `json:"children"`
}

func (menu *MenuModel) CollectionName() string {
	return config.CollGmMenu
}

func (menu *MenuModel) GetId() primitive.ObjectID {
	return menu.Id
}

func (menu *MenuModel) SetId(id primitive.ObjectID) {
	menu.Id = id
}

// 写入一条数据
func (menu *MenuModel) InsertMenu() error {
	if menu.MId == 0 {
		opt := options.FindOne().SetSort(bson.D{{"mid", -1}})
		menuTmp, err := new(MenuModel).SelectOneMenu(nil, opt)
		if err != nil {
			return err
		}
		if menuTmp == nil {
			menu.MId = 2
		} else {
			menu.MId = menuTmp.MId + 1
		}
	}
	menu.Id = primitive.NewObjectID()
	err := mgoDB.GetMgo().InsertOne(nil, menu)
	return err
}

func (menu *MenuModel) SelectMenuByMId(mid int) error {
	filter := bson.D{
		{"mid", mid},
	}
	finder := mgoDB.NewOneFinder(menu).Where(filter)
	_, err := mgoDB.GetMgo().FindOne(nil, finder)
	return err
}

func (menu *MenuModel) SelectAllMenus(filter bson.D, options *options.FindOptions) (data []*MenuModel, err error) {
	finder := mgoDB.NewFinder(menu).Where(filter).Options(options).Records(&data)
	err = mgoDB.GetMgo().FindMany(nil, finder)
	return
}

func (menu *MenuModel) DelMenuById(mid int) error {
	filter := mgoDB.NewDeleter(menu).Where(bson.D{
		{"mid", mid},
	})
	_, err := mgoDB.GetMgo().DeleteMany(nil, filter)
	return err
}

func (menu *MenuModel) SelectOneMenu(filter bson.D, opt *options.FindOneOptions) (data *MenuModel, err error) {
	oneFilter := mgoDB.NewOneFinder(menu).Where(filter).Options(opt).Record(&data)
	res, err := mgoDB.GetMgo().FindOne(nil, oneFilter)
	if err != nil {
		return nil, err
	}
	if !res {
		return nil, nil
	}

	return data, err
}

func (menu *MenuModel) UpdateMenuByMid() error {
	filter := bson.D{
		{"mid", menu.MId},
	}
	//data,err := bson.Marshal(menu)
	//if err != nil {
	//	return err
	//}
	_, err := mgoDB.GetMgo().DB().Collection(menu.CollectionName()).UpdateOne(nil, filter, bson.D{{"$set", menu}})
	return err
}

func (menu *MenuModel) GetLevelMenuList(pmid int, menuType int) ([]*LevelMenu, error) {
	var LevelMenuData []*LevelMenu
	if pmid == 0 {
		return nil, nil
	}
	where := bson.D{{"pmid", pmid}}
	if menuType > 0 {
		where = append(where, bson.E{Key: "type", Value: menuType})
	}
	Children, err := menu.SelectAllMenus(where, nil)
	if err != nil {
		return nil, err
	}
	for _, item := range Children {
		list, err := menu.GetLevelMenuList(item.MId, menuType)
		if err != nil {
			return nil, err
		}
		LevelMenuData = append(LevelMenuData, &LevelMenu{
			MId:      item.MId,
			Type:     item.Type,
			Name:     item.Name,
			PMid:     item.PMid,
			Children: list,
		})
	}
	return LevelMenuData, nil
}

func (menu *MenuModel) GetMenuNameByMids(mid []int) ([]map[string]string, error) {
	var result []map[string]string
	finder := mgoDB.NewFinder(menu).Where(bson.D{{"mid", bson.M{"$in": mid}}}).Options(options.Find().SetProjection(bson.M{"path": 1, "_id": 0})).Records(&result)
	err := mgoDB.GetMgo().FindMany(nil, finder)
	return result, err
}

func (menu *MenuModel) GetMenuMid(d bson.D) (int, error) {
	result := make(map[string]int)
	finder := mgoDB.NewOneFinder(menu).Where(d).Options(options.FindOne().SetProjection(bson.M{"mid": 1, "_id": 0})).Record(result)
	res, err := mgoDB.GetMgo().FindOne(context.TODO(), finder)
	if err != nil {
		return 0, err
	}
	if res {
		return result["mid"], nil
	}
	return 0, nil
}
