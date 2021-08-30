package global

import (
	"com.xv.admin.server/utils/xlsx"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
)

const (
	HeroDataXlsx          = "/xlsx/HeroData.xlsx"          // 英雄表
	ChestDataXlsx         = "/xlsx/ChestData.xlsx"         // 宝箱列表
	HeroChestDataXlsx     = "/xlsx/HeroChestData.xlsx"     // 英雄宝箱
	MatChestDataXlsx      = "/xlsx/MatChestData.xlsx"      // 材料宝箱
	EquippedChestDataXlsx = "/xlsx/EquippedChestData.xlsx" // 装备宝箱
	WeaponChestDataXlsx   = "/xlsx/WeaponChestData.xlsx"   // 武器宝箱
	WeaponDataXlsx        = "/xlsx/WeaponData.xlsx"        // 武器表,主 副 护甲 饰品 英雄武器,英雄装备
	CardDataXlsx          = "/xlsx/CardData.xlsx"          // 各种升级卡配置表
	GoodsDataDataXlsx     = "/xlsx/GoodsData.xlsx"         //商城商品价格

)

// 导表
func Init() {
	if runtime.GOOS != "windows" {
		return
	}
	src := "E:\\BattleRabbit\\Doc\\XLS\\project\\"                         // 源文件夹
	targetDir := "E:\\BattleRabbit\\Server\\GunFireServer\\config\\xlsx\\" // 目标文件夹
	files := []string{
		//SevenDayXlsx,
		//AdvertisingVideoXlsx,
		//WeaponReformXlsx,
		HeroDataXlsx,
		ChestDataXlsx,
		HeroChestDataXlsx,
		MatChestDataXlsx,
		EquippedChestDataXlsx,
		WeaponChestDataXlsx,
		WeaponDataXlsx,
		CardDataXlsx,
		GoodsDataDataXlsx,
	}
	for _, file := range files {
		file = strings.Split(file, "/")[2]
		err := os.Remove(targetDir + "\\" + file)
		if err != nil {
			log.Fatal(err)
		}
		source, err := os.Open(src + "\\" + file)
		if err != nil {
			log.Fatal(err)
		}
		defer source.Close()

		destination, err := os.Create(targetDir + "\\" + file)
		if err != nil {
			log.Fatal(err)
		}
		defer destination.Close()
		_, err = io.Copy(destination, source)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("导表完成: file ====> %s", file)
	}
}

type readRule struct {
	filename string              // 上面定义的文件名称
	storage  func(v interface{}) // 单行回调处理
	sheet    string              // 单页标题 默认sheet1
	tmp      interface{}         // 临时交换变量  // 必须是地址类型的struct
}

// 定义全局访问的Xls Config 变量(只读,不可修改)
var (
	HeroDataConf          = make(map[string]*HeroDataConfig)       //英雄数据表
	ChestDataConf         = make(map[string]*ChestDataConfig)      // 宝箱列表
	ChestHeroDataConf     = make(map[string]*ChestListConfig)      // 英雄宝箱
	ChestMatDataConf      = make(map[string]*ChestListConfig)      // 材料宝箱
	ChestEquippedDataConf = make(map[string]*ChestListConfig)      // 装备宝箱
	ChestWeaponDataConf   = make(map[string]*ChestListConfig)      // 武器宝箱
	RoleMainWeaponConf    = make(map[string]*EquipagePublicConfig) //角色主武器表
	RoleSubWeaponConf     = make(map[string]*EquipagePublicConfig) //副武器表
	ArmorDataConf         = make(map[string]*EquipagePublicConfig) //护甲装备表
	OrnamentsDataConf     = make(map[string]*EquipagePublicConfig) //护甲装备表
	CardDataConf          = make(map[string]*CardDataConfig)       // 各种升级卡配置表
	GoodsDataConf         = make(map[string]*GoodsDataConfig)      //商城商品价格
)

func InitXlsxConfig(binPath string) {
	Init()
	for _, collection := range GetCollections() {
		err := xlsx.LoadXlsxFile(binPath+collection.filename, collection.tmp, collection.sheet, collection.storage)
		if err != nil {
			log.Fatal(err)
		}
	}
	return
}

// 通过关联ID和等级获得配置表ID
func GetConfigId(relationId string, lv int) string {
	return fmt.Sprintf("%s_%d", relationId, lv)
}

func GetCollections() []*readRule {
	// 读取集合
	return []*readRule{
		{
			filename: HeroDataXlsx,
			tmp:      &HeroDataConfig{},
			sheet:    "HeroDataTable",
			storage: func(v interface{}) {
				if val, ok := v.(*HeroDataConfig); ok {
					HeroDataConf[val.RelationId] = val
				}
			},
		}, {
			filename: ChestDataXlsx,
			tmp:      &ChestDataConfig{},
			sheet:    "Sheet1",
			storage: func(v interface{}) {
				if val, ok := v.(*ChestDataConfig); ok {
					ChestDataConf[val.Id] = val
				}
			},
		}, {
			filename: HeroChestDataXlsx,
			tmp:      &ChestListConfig{},
			sheet:    "Sheet1",
			storage: func(v interface{}) {
				if val, ok := v.(*ChestListConfig); ok {
					ChestHeroDataConf[val.Id] = val
				}
			},
		}, {
			filename: MatChestDataXlsx,
			tmp:      &ChestListConfig{},
			sheet:    "Sheet1",
			storage: func(v interface{}) {
				if val, ok := v.(*ChestListConfig); ok {
					ChestMatDataConf[val.Id] = val
				}
			},
		}, {
			filename: EquippedChestDataXlsx,
			tmp:      &ChestListConfig{},
			sheet:    "Sheet1",
			storage: func(v interface{}) {
				if val, ok := v.(*ChestListConfig); ok {
					ChestEquippedDataConf[val.Id] = val
				}
			},
		}, {
			filename: WeaponChestDataXlsx,
			tmp:      &ChestListConfig{},
			storage: func(v interface{}) {
				if val, ok := v.(*ChestListConfig); ok {
					ChestWeaponDataConf[val.Id] = val
				}
			},
		}, {
			filename: WeaponDataXlsx,
			tmp:      &EquipagePublicConfig{},
			sheet:    "RoleMainWeaponDataTable",
			storage: func(v interface{}) {
				if val, ok := v.(*EquipagePublicConfig); ok && val.Id != "" {
					RoleMainWeaponConf[val.RelationId] = val
				}
			},
		}, {
			filename: WeaponDataXlsx,
			tmp:      &EquipagePublicConfig{},
			sheet:    "RoleSubWeaponDataTable",
			storage: func(v interface{}) {
				if val, ok := v.(*EquipagePublicConfig); ok && val.Id != "" {
					RoleSubWeaponConf[val.RelationId] = val
				}
			},
		}, {
			filename: WeaponDataXlsx,
			tmp:      &EquipagePublicConfig{},
			sheet:    "ArmorTable",
			storage: func(v interface{}) {
				if val, ok := v.(*EquipagePublicConfig); ok && val.Id != "" {
					ArmorDataConf[val.RelationId] = val
				}
			},
		}, {
			filename: WeaponDataXlsx,
			tmp:      &EquipagePublicConfig{},
			sheet:    "OrnamentsDataTable",
			storage: func(v interface{}) {
				if val, ok := v.(*EquipagePublicConfig); ok && val.Id != "" {
					OrnamentsDataConf[val.RelationId] = val
				}
			},
		}, {
			filename: CardDataXlsx,
			tmp:      &CardDataConfig{},
			storage: func(v interface{}) {
				if val, ok := v.(*CardDataConfig); ok {
					CardDataConf[val.Id] = val
				}
			},
		}, {
			filename: GoodsDataDataXlsx,
			tmp:      &GoodsDataConfig{},
			storage: func(v interface{}) {
				if val, ok := v.(*GoodsDataConfig); ok {
					if val.Id != "" {
						GoodsDataConf[val.Id] = val
					}
				}
			},
		},
	}
}

type CardDataConfig struct {
	Id   string `xlsx:"id"`
	Name string `xlsx:"name"`
	Ty   int8   `xlsx:"type"` // 1 //主角经验卡  2 //主角升星卡  3 //装备升品卡  4 //副武器改造卡  5 //英雄升级卡  6 //英雄升星卡  7 //英雄强化卡  8 //英雄技能卡  9 //英雄装备卡
	Rare int    `xlsx:"rare"` // 稀有度 1:N 2:R 3:SR 4:SSR (升星卡 专属)
	Exp  int    `xlsx:"exp"`  // 经验卡携带的经验值
}
type ChestDataConfig struct {
	Id                  string `xlsx:"id"`
	OneRequirementCount int    `xlsx:"onerequirementcount"` //单次价格
	TenRequirementCount int    `xlsx:"tenrequirementcount"` //多次价格
	TimeLimit           int    `xlsx:"timelimit"`           //免费分钟数
}

type ChestListConfig struct {
	Id       string `xlsx:"id"`       //当前宝箱奖励的物品顺序id
	MatId    string `xlsx:"matid"`    //材料id
	Droprate int    `xlsx:"droprate"` //掉落几率
	MinCount int    `xlsx:"mincount"` //最小数
	MaxCount int    `xlsx:"maxcount"` //最大数
	Type     int    `xlsx:"type"`     //类型
	Des      string `xlsx:"des"`      //描述
}

// 英雄数据表
type HeroDataConfig struct {
	Id                string `xlsx:"id"`                 // 英雄ID
	Name              string `xlsx:"name"`               // 英雄ID
	RelationId        string `xlsx:"relation_id"`        // 关联ID
	WeaponId          string `xlsx:"weapon_id"`          // 英雄武器ID
	Rarity            int8   `xlsx:"rarity"`             // 稀有度
	MaxLevel          int    `xlsx:"max_level"`          // 最大等级
	Level             int    `xlsx:"level"`              // 等级
	StarLevel         int    `xlsx:"star_level"`         // 星级
	Exp               int    `xlsx:"experience"`         // 经验值
	ASId              string `xlsx:"active_skill_id"`    // 主动技能[关联]ID
	PSId1             string `xlsx:"passive_skill_id_1"` // 被动技能1[关联]ID
	PSId2             string `xlsx:"passive_skill_id_2"` // 被动技能2[关联]ID
	PSId3             string `xlsx:"passive_skill_id_3"` // 被动技能3[关联]ID
	PSId4             string `xlsx:"passive_skill_id_4"` // 被动技能4[关联]ID
	BSId              string `xlsx:"buff_id"`            // buff技能ID
	FSId              string `xlsx:"fetter_skill_id"`    // 羁绊技能ID
	Attack            int    `xlsx:"attack"`             // 攻击力
	Critical          int    `xlsx:"critical"`           // 暴击
	MoveSpeed         int    `xlsx:"move_speed"`         // 移速
	StarDeplete       string `xlsx:"star_deplete"`       // 升星消耗
	StarGold          int    `xlsx:"star_gold"`          // 升星金币
	StarDiamond       int    `xlsx:"star_diamond"`       // 升星钻石
	UpDeplete         string `xlsx:"up_deplete"`         // 升级消耗
	UpGold            int    `xlsx:"up_gold"`            // 升级金币
	UpDiamond         int    `xlsx:"up_diamond"`         // 升级钻石
	StrengthenDep     string `xlsx:"strengthen_deplete"` // 强化消耗
	StrengthenGold    int    `xlsx:"strengthen_gold"`    // 强化金币
	StrengthenDiamond int    `xlsx:"strengthen_diamond"` // 强化钻石
}

type EquipagePublicConfig struct {
	Id               string `xlsx:"id"`                // 主武器ID
	Name             string `xlsx:"name"`              // 数据表关联ID
	RelationId       string `xlsx:"relation_id"`       // 数据表关联ID
	Index            int    `xlsx:"index"`             // 排序索引
	Level            int    `xlsx:"level"`             // 等级
	MaxLevel         int    `xlsx:"maxlevel"`          // 最高等级
	StarLevel        int    `xlsx:"starlevel"`         // 星级
	HurtAttribute    int    `xlsx:"hurt_attribute"`    // 属性(不知道在那儿用)
	UpgradeMaterials string `xlsx:"upgrade_materials"` // 升级材料.后期若材料多种格式就使用: 材料ID:需求数量|材料ID:需求数量|材料ID:需求数量
	UpgradeGold      int    `xlsx:"upgrade_gold"`      // 升级消耗的金币
	UpgradeDiamond   int    `xlsx:"upgrade_diamond"`   // 升级消耗的钻石
}
type GoodsDataConfig struct {
	Id           string `xlsx:"id"`           //ID
	Type         int    `xlsx:"type"`         //商品类型
	Price        int    `xlsx:"price"`        //价格
	Count        int    `xlsx:"count"`        //获得数量
	Presentation int    `xlsx:"presentation"` //赠送数量
}
