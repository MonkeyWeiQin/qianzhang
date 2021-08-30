package excel

import (
	"battle_rabbit/define"
	"battle_rabbit/global"
	"battle_rabbit/service/log"
	"battle_rabbit/utils/xlsx"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
)

// 导表
func Init() {
	if runtime.GOOS != "windows" {
		return
	}
	src := "D:\\project\\go\\BattleRabbit\\Doc\\XLS\\project\\"                      // 源文件夹
	targetDir := "D:\\project\\go\\BattleRabbit\\Server\\GameServer\\bin\\xlsx\\" // 目标文件夹
	files := []string{
		define.SevenDayXlsx,
		//define.AdvertisingVideoXlsx,
		define.GeneralStageXlsx,
		define.BaseMat,
		define.WeaponDataXlsx,
		define.AccountUpgradeXlsx,
		define.HeroDataXlsx,
		define.RolePassiveExpendXlsx,
		define.RolePassiveSkillXlsx,
		define.RoleSkinXlsx,
		define.SkillDataXlsx,
		define.RoleDataXlsx,
		define.CardDataXlsx,
		define.RewardDataXlsx,
		define.ChestDataXlsx,
		define.ShopChestDataXlsx,
		define.GoodsDataDataXlsx,
		define.ChallengeStageXlsx,
		define.ResourceStageXlsx,
		define.TaskDataXlsx,
	}
	for _, file := range files {
		file = strings.Split(file, "/")[2]
		err := os.Remove(targetDir + "\\" + file)
		if err != nil {
			log.Warn(err)
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
		log.Debug("导表完成: file ====> %s", file)
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
	SevenDayConf          = make(map[int]*SevenDayConfig)
	AdvertisingVideoConf  = make(map[int]*AdvertisingVideoConfig)
	GeneralStageConf      = make(map[int]*StageConfig)
	DifficultStageConf    = make(map[int]*StageConfig)
	BaseMatConf           = make(map[string]*BaseMatConfig)
	WeaponReformConf      = make(map[string]*WeaponReformConfig)
	RoleMainWeaponConf    = make(map[string]*EquipagePublicConfig)   //角色主武器表
	RoleSubWeaponConf     = make(map[string]*EquipagePublicConfig)   //副武器表
	RoleArmorDataConf     = make(map[string]*EquipagePublicConfig)   //护甲装备表
	RoleOrnamentsDataConf = make(map[string]*EquipagePublicConfig)   //护甲装备表
	HeroWeaponDataConf    = make(map[string]*EquipagePublicConfig)   //英雄武器表
	HeroEquipmentDataConf = make(map[string]*EquipagePublicConfig)   //英雄装备表
	AccountUpgradeConf    = make(map[int]*AccountUpgradeConfig)      //账户升级表
	HeroDataConf          = make(map[string]*HeroDataConfig)         //英雄数据表
	RolePassiveExpendConf = make(map[int]*RolePassiveExpendConfig)   //天赋升级消耗表
	RolePassiveSkillConf  = make(map[string]*RolePassiveSkillConfig) //天赋属性表 读数据: 关联ID+'_'+当前等级,passive01_1
	RoleSkinDataConf      = make(map[string]*RoleSkinDataConfig)     // 角色皮肤表
	HeroPassiveSkillConf  = make(map[string]*SkillConfig)            // 英雄被动技能表
	HeroFetterSkillConf   = make(map[string]*SkillConfig)            // 英雄羁绊技能表
	HeroBuffSKillConf     = make(map[string]*SkillConfig)            // 英雄Buff技能表
	HeroActiveSkillConf   = make(map[string]*SkillConfig)            // 英雄主动技能表
	RoleDataConf          = make(map[string]*RoleDataConfig)         // 角色数据表
	CardDataConf          = make(map[string]*CardDataConfig)         // 各种升级卡配置表
	RewardDataConf        = make(map[string]*RewardDataConfig)       // 关卡奖励 // 普通关卡 | 困难关卡 | 章节奖励
	ActiveStageRewardConf = make(map[string]*RewardDataConfig)       // 活动关卡奖励表 1:ChallengeReward 2: ResourceReward
	ShopRewardDataConf    = make(map[string]*RewardDataConfig)       // 商店累计购买满次数后必得奖励表
	ChestDataConf         = make(map[string]*ChestDataConfig)        // 宝箱列表
	ChestHeroDataConf     = make(map[string]*ChestListConfig)        // 英雄宝箱
	ChestMatDataConf      = make(map[string]*ChestListConfig)        // 材料宝箱
	ChestEquippedDataConf = make(map[string]*ChestListConfig)        // 装备宝箱
	ChestWeaponDataConf   = make(map[string]*ChestListConfig)        // 武器宝箱
	GoodsDataConf         = make(map[string]*GoodsDataConfig)        // 商城商品价格
	ChallengeStageConf    = make(map[string]*ChallengeStageConfig)   // 挑战关卡配置
	ResourceStageConf     = make(map[string]*ResourceStageConfig)    // 资源关卡关卡配置
	TaskMainDataConf      = make(map[int]*TaskDataConfig)            // 主线任务
	TaskDayDataConf       = make(map[int]*TaskDataConfig)            // 每日任务
	TaskRewardDataConf    = make(map[string]*RewardDataConfig)       // 任务奖励表
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
			filename: define.SevenDayXlsx,
			tmp:      &SevenDayConfig{},
			storage: func(v interface{}) {
				if val, ok := v.(*SevenDayConfig); ok && val.Id != 0 {
					SevenDayConf[val.Id] = val
				}
			},
		}, {
			filename: define.AdvertisingVideoXlsx,
			tmp:      &AdvertisingVideoConfig{},
			storage: func(v interface{}) {
				if val, ok := v.(*AdvertisingVideoConfig); ok && val.Id != 0 {
					AdvertisingVideoConf[val.Id] = val
				}
			},
		}, {
			filename: define.GeneralStageXlsx,
			tmp:      &StageConfig{},
			sheet:    "GeneralStage",
			storage: func(v interface{}) {
				if val, ok := v.(*StageConfig); ok && val.Id != 0 {
					GeneralStageConf[val.Id] = val
				}
			},
		}, {
			filename: define.BaseMat,
			tmp:      &BaseMatConfig{},
			storage: func(v interface{}) {
				if val, ok := v.(*BaseMatConfig); ok && val.Id != "" {
					BaseMatConf[val.Id] = val
				}
			},
		}, {
			filename: define.GeneralStageXlsx,
			tmp:      &StageConfig{},
			sheet:    "DifficultStage",
			storage: func(v interface{}) {
				if val, ok := v.(*StageConfig); ok && val.Id != 0 {
					DifficultStageConf[val.Id] = val
				}
			},
		}, {
			filename: define.WeaponReformXlsx,
			tmp:      &WeaponReformConfig{},
			storage: func(v interface{}) {
				if val, ok := v.(*WeaponReformConfig); ok && val.Id != "" {
					WeaponReformConf[val.Id] = val
				}
			},
		}, {
			filename: define.WeaponDataXlsx,
			tmp:      &EquipagePublicConfig{},
			sheet:    "RoleMainWeaponDataTable",
			storage: func(v interface{}) {
				if val, ok := v.(*EquipagePublicConfig); ok && val.Id != "" {
					RoleMainWeaponConf[val.Id] = val
				}
			},
		}, {
			filename: define.WeaponDataXlsx,
			tmp:      &EquipagePublicConfig{},
			sheet:    "RoleSubWeaponDataTable",
			storage: func(v interface{}) {
				if val, ok := v.(*EquipagePublicConfig); ok && val.Id != "" {
					RoleSubWeaponConf[val.Id] = val
				}
			},
		}, {
			filename: define.WeaponDataXlsx,
			tmp:      &EquipagePublicConfig{},
			sheet:    "ArmorTable",
			storage: func(v interface{}) {
				if val, ok := v.(*EquipagePublicConfig); ok && val.Id != "" {
					RoleArmorDataConf[val.Id] = val
				}
			},
		}, {
			filename: define.WeaponDataXlsx,
			tmp:      &EquipagePublicConfig{},
			sheet:    "OrnamentsDataTable",
			storage: func(v interface{}) {
				if val, ok := v.(*EquipagePublicConfig); ok && val.Id != "" {
					RoleOrnamentsDataConf[val.Id] = val
				}
			},
		}, {
			filename: define.WeaponDataXlsx,
			tmp:      &EquipagePublicConfig{},
			sheet:    "HeroWeaponData",
			storage: func(v interface{}) {
				if val, ok := v.(*EquipagePublicConfig); ok {
					HeroWeaponDataConf[val.Id] = val
				}
			},
		}, {
			filename: define.WeaponDataXlsx,
			tmp:      &EquipagePublicConfig{},
			sheet:    "HeroEquipment",
			storage: func(v interface{}) {
				if val, ok := v.(*EquipagePublicConfig); ok {
					HeroEquipmentDataConf[val.Id] = val
				}
			},
		}, {
			filename: define.AccountUpgradeXlsx,
			tmp:      &AccountUpgradeConfig{},
			storage: func(v interface{}) {
				if val, ok := v.(*AccountUpgradeConfig); ok {
					AccountUpgradeConf[val.Lv] = val
				}
			},
		}, {
			filename: define.HeroDataXlsx,
			tmp:      &HeroDataConfig{},
			sheet:    "HeroDataTable",
			storage: func(v interface{}) {
				if val, ok := v.(*HeroDataConfig); ok {
					HeroDataConf[val.Id] = val
				}
			},
		}, {
			filename: define.RolePassiveExpendXlsx,
			tmp:      &RolePassiveExpendConfig{},
			storage: func(v interface{}) {
				if val, ok := v.(*RolePassiveExpendConfig); ok {
					RolePassiveExpendConf[val.Lv] = val
				}
			},
		}, {
			filename: define.RolePassiveSkillXlsx,
			tmp:      &RolePassiveSkillConfig{},
			storage: func(v interface{}) {
				if val, ok := v.(*RolePassiveSkillConfig); ok {
					RolePassiveSkillConf[val.Id] = val
				}
			},
		}, {
			filename: define.RoleSkinXlsx,
			tmp:      &RoleSkinDataConfig{},
			storage: func(v interface{}) {
				if val, ok := v.(*RoleSkinDataConfig); ok {
					RoleSkinDataConf[val.Id] = val
				}
			},
		}, {
			filename: define.SkillDataXlsx,
			tmp:      &SkillConfig{},
			sheet:    "HeroPassiveSkill", // 被动
			storage: func(v interface{}) {
				if val, ok := v.(*SkillConfig); ok {
					HeroPassiveSkillConf[val.Id] = val
				}
			},
		}, {
			filename: define.SkillDataXlsx,
			tmp:      &SkillConfig{},
			sheet:    "HeroFetterSkill", // 羁绊
			storage: func(v interface{}) {
				if val, ok := v.(*SkillConfig); ok {
					HeroFetterSkillConf[val.Id] = val
				}
			},
		}, {
			filename: define.SkillDataXlsx,
			tmp:      &SkillConfig{},
			sheet:    "HeroBuffSKill", // BUff
			storage: func(v interface{}) {
				if val, ok := v.(*SkillConfig); ok {
					HeroBuffSKillConf[val.Id] = val
				}
			},
		}, {
			filename: define.SkillDataXlsx,
			tmp:      &SkillConfig{},
			sheet:    "HeroActiveSkill", // 英雄主动技能
			storage: func(v interface{}) {
				if val, ok := v.(*SkillConfig); ok {
					HeroActiveSkillConf[val.Id] = val
				}
			},
		}, {
			filename: define.RoleDataXlsx,
			tmp:      &RoleDataConfig{},
			sheet:    "Sheet2",
			storage: func(v interface{}) {
				if val, ok := v.(*RoleDataConfig); ok {
					RoleDataConf[val.Id] = val
				}
			},
		}, {
			filename: define.CardDataXlsx,
			tmp:      &CardDataConfig{},
			storage: func(v interface{}) {
				if val, ok := v.(*CardDataConfig); ok {
					CardDataConf[val.Id] = val
				}
			},
		}, {
			filename: define.RewardDataXlsx,
			tmp:      &RewardDataConfig{},
			sheet:    "FirstBrushStageReward", // 第一次通关奖励
			storage: func(v interface{}) {
				if val, ok := v.(*RewardDataConfig); ok {
					RewardDataConf[val.Id] = val
				}
			},
		}, {
			filename: define.RewardDataXlsx,
			tmp:      &RewardDataConfig{},
			sheet:    "SecondBrushStageReward", // 第二次通关奖励
			storage: func(v interface{}) {
				if val, ok := v.(*RewardDataConfig); ok {
					RewardDataConf[val.Id] = val
				}
			},
		}, {
			filename: define.RewardDataXlsx,
			tmp:      &RewardDataConfig{},
			sheet:    "ChapterReward", // 大章节通关奖励
			storage: func(v interface{}) {
				if val, ok := v.(*RewardDataConfig); ok {
					RewardDataConf[val.Id] = val
				}
			},
		},  {
			filename: define.RewardDataXlsx,
			tmp:      &RewardDataConfig{},
			sheet:    "TaskReward", // 任务完成奖励(主线和每日任务)
			storage: func(v interface{}) {
				if val, ok := v.(*RewardDataConfig); ok {
					TaskRewardDataConf[val.Id] = val
				}
			},
		}, {
			filename: define.RewardDataXlsx,
			tmp:      &RewardDataConfig{},
			sheet:    "ChallengeReward", // 挑战关卡奖励
			storage: func(v interface{}) {
				if val, ok := v.(*RewardDataConfig); ok {
					ActiveStageRewardConf[val.Id] = val
				}
			},
		}, {
			filename: define.RewardDataXlsx,
			tmp:      &RewardDataConfig{},
			sheet:    "EndlessReward", // 无尽关卡奖励(月末排行奖励)
			storage: func(v interface{}) {
				if val, ok := v.(*RewardDataConfig); ok {
					ActiveStageRewardConf[val.Id] = val
				}
			},
		}, {
			filename: define.RewardDataXlsx,
			tmp:      &RewardDataConfig{},
			sheet:    "ShopReward", // 商店累计购买次数达标后的奖励
			storage: func(v interface{}) {
				if val, ok := v.(*RewardDataConfig); ok {
					ShopRewardDataConf[val.Id] = val
				}
			},
		}, {
			filename: define.ChestDataXlsx,
			tmp:      &ChestDataConfig{},
			sheet:    "Sheet1",
			storage: func(v interface{}) {
				if val, ok := v.(*ChestDataConfig); ok {
					ChestDataConf[val.Id] = val
				}
			},
		}, {
			filename: define.ShopChestDataXlsx,
			tmp:      &ChestListConfig{},
			sheet:    "heroChest",
			storage: func(v interface{}) {
				if val, ok := v.(*ChestListConfig); ok {
					if val.Droprate > 0 {
						ChestHeroDataConf[val.Id] = val
					}
				}
			},
		}, {
			filename: define.ShopChestDataXlsx,
			tmp:      &ChestListConfig{},
			sheet:    "matChest",
			storage: func(v interface{}) {
				if val, ok := v.(*ChestListConfig); ok {
					if val.Droprate > 0 {
						ChestMatDataConf[val.Id] = val
					}
				}
			},
		}, {
			filename: define.ShopChestDataXlsx,
			tmp:      &ChestListConfig{},
			sheet:    "equipChest",
			storage: func(v interface{}) {
				if val, ok := v.(*ChestListConfig); ok {
					if val.Droprate > 0 {
						ChestEquippedDataConf[val.Id] = val
					}
				}
			},
		}, {
			filename: define.ShopChestDataXlsx,
			tmp:      &ChestListConfig{},
			sheet:    "weaponChest",
			storage: func(v interface{}) {
				if val, ok := v.(*ChestListConfig); ok {
					if val.Droprate > 0 {
						ChestWeaponDataConf[val.Id] = val
					}
				}
			},
		}, {
			filename: define.GoodsDataDataXlsx,
			tmp:      &GoodsDataConfig{},
			storage: func(v interface{}) {
				if val, ok := v.(*GoodsDataConfig); ok && val.Id != "" {
					GoodsDataConf[val.Id] = val
				}
			},
		}, {
			filename: define.ChallengeStageXlsx,
			tmp:      &ChallengeStageConfig{},
			storage: func(v interface{}) {
				if val, ok := v.(*ChallengeStageConfig); ok && val.Id != "" {
					ChallengeStageConf[val.Id] = val
				}
			},
		}, {
			filename: define.ResourceStageXlsx,
			tmp:      &ResourceStageConfig{},
			storage: func(v interface{}) {
				if val, ok := v.(*ResourceStageConfig); ok && val.Id != "" {
					ResourceStageConf[val.Id] = val
				}
			},
		}, {
			filename: define.TaskDataXlsx,
			tmp:      &TaskDataConfig{},
			sheet:    "MainTask",
			storage: func(v interface{}) {
				if val, ok := v.(*TaskDataConfig); ok && val.Id != 0 {
					TaskMainDataConf[val.Id] = val
				}
			},
		}, {
			filename: define.TaskDataXlsx,
			tmp:      &TaskDataConfig{},
			sheet:    "DayTask",
			storage: func(v interface{}) {
				if val, ok := v.(*TaskDataConfig); ok && val.Id != 0 {
					TaskDayDataConf[val.Id] = val
				}
			},
		},
	}
}

type SevenDayConfig struct {
	Id       int    `xlsx:"id"`
	ItemType int    `xlsx:"award_type"`
	ItemID   string `xlsx:"award_id"`
	Count    int    `xlsx:"award_count"`
}

type AdvertisingVideoConfig struct {
	Id          int     `xlsx:"id"`           // ID
	AbsName     string  `xlsx:"abs_name"`     // 内部使用，不显示
	AbsType     int     `xlsx:"abs_type"`     // 奖励种类：随机 倍数 次数 固定值 必须看
	AbsMultiple float64 `xlsx:"abs_multiple"` // 倍数
	AbsMinGold  float64 `xlsx:"abs_min_gold"` // 最小数
	AbsMaxGold  float64 `xlsx:"abs_max_gold"` // 最大数
	AbsTimes    int     `xlsx:"abs_times"`    // 次数
}
type StageConfig struct {
	Id            int    `xlsx:"id"` // 关卡ID
	MapId         string `xlsx:"mapid"`
	ChapterReward string `xlsx:"chapterreward"` // 章节奖励
	NextLevelId   int    `xlsx:"nextlevelid"`   // 下一个关卡的ID
	LevelType     int    `xlsx:"leveltype"`     // 关卡类型
	EnergyConsume int    `xlsx:"energyconsume"` // 体力消耗
	SucceedReward string `xlsx:"succeedreward"` // 通关奖励id
	RepeatReward  string `xlsx:"repeatreward"`  // 重复玩奖励id
	FailureReward string `xlsx:"failurereward"` // 通关失败奖励id
	Experience    int    `xlsx:"experience"`    // 经验值
}

// 营地物品掉落表
type BaseMatConfig struct {
	Id          string `xlsx:"id"`          //营地材料编号
	Name        string `xlsx:"name"`        //材料名称
	DropTime    int    `xlsx:"droptime"`    //掉落有效时间/秒
	Probability int    `xlsx:"probability"` //掉落机率
	SalePrice   int    `xlsx:"sale_price"`  //出售价格
}

type WeaponReformConfig struct {
	Id              string `xlsx:"id"`               //武器ID
	RelationId      string `xlsx:"relation_id"`      //关联ID
	ReformMaterials int    `xlsx:"reform_materials"` //改造材料
	ReformFragments int    `xlsx:"reform_fragments"` //武器碎片
	ReformDiamond   int    `xlsx:"reform_diamond"`   //钻石
}

//装备公共数据结构,包括以下表使用
//角色主武器表,
//副武器表,
//装备部件表,
//饰品部件表,
//英雄武器表
//角色皮肤表
type EquipagePublicConfig struct {
	Id               string            `xlsx:"id"`                 // 主武器ID
	RelationId       string            `xlsx:"relation_id"`        // 数据表关联ID
	Index            int               `xlsx:"index"`              // 排序索引
	Level            int               `xlsx:"level"`              // 等级
	MaxLevel         int               `xlsx:"maxlevel"`           // 最高等级
	StarLevel        int               `xlsx:"starlevel"`          // 星级
	UpgradeMaterials string            `xlsx:"upgrade_materials"`  // 升级材料.后期若材料多种格式就使用: 材料ID:需求数量|材料ID:需求数量|材料ID:需求数量
	UpgradeGold      int               `xlsx:"upgrade_gold"`       // 升级消耗的金币
	UpgradeDiamond   int               `xlsx:"upgrade_diamond"`    // 升级消耗的钻石
	SplitUpDeplete   string            `xlsx:"split_up_deplete"`   // 拆分为升级材料
	StarDeplete      string            `xlsx:"star_deplete"`       // 升星消耗
	StarGold         int               `xlsx:"star_gold"`          // 升星金币
	StarDiamond      int               `xlsx:"star_diamond"`       // 升星钻石
	SplitStarDeplete string            `xlsx:"split_star_deplete"` // 拆分为星级材料
	Attribute        *global.Attribute `xlsx:"-"`                  // 属性列表
}

// 账号升级表
type AccountUpgradeConfig struct {
	Lv       int `xlsx:"id"`        // 玩家等级
	Exp      int `xlsx:"exp"`       // 升级所需经验值
	TotalExp int `xlsx:"total_exp"` // 升级所需经验值
	MaxLv    int `xlsx:"max_lv"`    // 总等级
}

// 英雄数据表
type HeroDataConfig struct {
	Id                string `xlsx:"id"`                 // 英雄ID
	RelationId        string `xlsx:"relation_id"`        // 关联ID
	WeaponId          string `xlsx:"weapon_id"`          // 武器ID
	Index             int    `xlsx:"index"`              // 排序序列
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
	SplitUpDeplete    string `xlsx:"split_up_deplete"`   // 拆分为升级材料
	SplitStarDeplete  string `xlsx:"split_star_deplete"` // 拆分为星级材料
}

// 天赋升级消耗表
type RolePassiveExpendConfig struct {
	Lv     int  `xlsx:"id"`        // 玩家等级
	Gold   int  `xlsx:"lvup_gold"` // 升级所需金币
	MaxLv  int8 `xlsx:"max_lv"`    // 单技能最高等级
	LiftLv int8 `xlsx:"lift_lv"`   // 生命值等级
	UpNum  int  `xlsx:"up_num"`    // 总升级次数

}

// 主角天赋技能表
type RolePassiveSkillConfig struct {
	Id               string  `xlsx:"id"`                //编号ID
	Relation         string  `xlsx:"relation_id"`       //关联ID
	StarLevel        int     `xlsx:"star_level"`        //星级
	Level            int     `xlsx:"level"`             //等级
	NextId           string  `xlsx:"next_id"`           //下一级id
	UpgradeGold      int     `xlsx:"upgrade_gold"`      //升级金币
	UpgradeDiamond   int     `xlsx:"upgrade_diamond"`   //升级钻石
	UpgradeMaterials string  `xlsx:"upgrade_materials"` //升级材料
	AttackA          string  `xlsx:"attack_a"`          //攻击力
	AttackC          float32 `xlsx:"attack_c"`          //攻击加成
	AttackSpeedC     float32 `xlsx:"attack_speed_c"`    //攻速加成
	DefenseC         float32 `xlsx:"defense_c"`         //防御加成
	CriticalC        int     `xlsx:"critical_c"`        //暴击几率加成
	GoldAddC         float32 `xlsx:"gold_add_c"`        //金币加成
	LifeC            float32 `xlsx:"life_c"`            //生命加成
	LifeA            int     `xlsx:"life_a"`            //生命加成
	BossHurtC        float32 `xlsx:"boss_hurt_c"`       //boss伤害加成
	MoveSpeedC       float32 `xlsx:"move_speed_c"`      //移动速度加成
	DodgeC           int     `xlsx:"dodge_c"`           //闪避几率加成
	BuffTimeC        float32 `xlsx:"buff_time_c"`       //buff时长加成
}

// 技能表通用数据结构 目前包括以下表使用
// 英雄被动技能表
// 英雄羁绊技能表
// 英雄buff技能表
type SkillConfig struct {
	Id               string            `xlsx:"id"`                // ID
	RelationId       string            `xlsx:"relation_id"`       // 数据表关联ID
	Level            int               `xlsx:"level"`             // 等级
	NextId           string            `xlsx:"next_id"`           // 下一级ID
	StarLevel        int               `xlsx:"star_level"`        // 星级
	HurtAttribute    int               `xlsx:"hurt_attribute"`    // 属性(不知道在那儿用)
	UpgradeMaterials string            `xlsx:"upgrade_materials"` // 升级材料.后期若材料多种格式就使用: 材料ID:需求数量|材料ID:需求数量|材料ID:需求数量
	UpgradeGold      int               `xlsx:"upgrade_gold"`      // 升级消耗的金币
	UpgradeDiamond   int               `xlsx:"upgrade_diamond"`   // 升级消耗的钻石
	TargetType       int               `xlsx:"target_type"`       // 被动技能接收对象(加成的目标对象) 0:自己(谁拥有这个技能就给谁) 1：角色 2：英雄	3：敌人(不处理)
	RealHurt         int               `xlsx:"real_hurt"`         // 伤害效果 (冗余数据,服务器暂不需要)
	RecoveryLift     int               `xlsx:"recovery_lift"`     // 回复效果 (冗余数据,服务器暂不需要)
	Attribute        *global.Attribute `xlsx:"-"`                 // 属性列表
}

// 角色配置表
type RoleDataConfig struct {
	Id            string            `xlsx:"id"`             // ID
	RelationId    string            `xlsx:"relation_id"`    // 数据表关联ID
	Level         int               `xlsx:"level"`          // 等级
	StarLevel     int               `xlsx:"star_level"`     // 星级
	Exp           int               `xlsx:"exp"`            // 升级经验消耗
	NextId        string            `xlsx:"next_id"`        // 下一级ID
	StarGold      int               `xlsx:"star_gold"`      // 升星金币消耗
	StarDiamond   int               `xlsx:"star_diamond"`   // 升星钻石消耗
	StarMaterials string            `xlsx:"star_materials"` // 升星材料消耗 (升星卡级数量)
	UpGold        int               `xlsx:"up_gold"`        // 升级金币消耗
	UpDiamond     int               `xlsx:"up_diamond"`     // 升级钻石消耗
	UpMaterials   int               `xlsx:"up_materials"`   // 升级材料消耗消耗
	Attribute     *global.Attribute `xlsx:"-"`              // 属性列表
}

type CardDataConfig struct {
	Id   string `xlsx:"id"`
	Ty   int8   `xlsx:"type"` // 1 //主角经验卡  2 //主角升星卡  3 //装备升品卡  4 //副武器改造卡  5 //英雄升级卡  6 //英雄升星卡  7 //英雄强化卡  8 //英雄技能卡  9 //英雄装备卡
	Rare int    `xlsx:"rare"` // 稀有度 1:N 2:R 3:SR 4:SSR (升星卡 专属)
	Exp  int    `xlsx:"exp"`  // 经验卡携带的经验值
}

type RewardDataConfig struct {
	Id            string `xlsx:"id"`
	Gold          int    `xlsx:"gold"`            //金币
	Diamond       int    `xlsx:"diamond"`         //钻石
	PhysicalPower int    `xlsx:"physicalpower"`   //体力
	Equipment     string `xlsx:"equipmat"`        //装备材料
	Exp           int    `xlsx:"exp"`             //经验
	Hero          string `xlsx:"heroid"`          // 英雄
	Card          string `xlsx:"cardids"`         // 卡
	SubWeapon     string `xlsx:"shieldid"`        // 副武器
	MainWeapon    string `xlsx:"mianweaponid"`    // 主武器
	Armor         string `xlsx:"mianarmorid"`     // 护甲
	Ornament      string `xlsx:"mianornamentsid"` // 饰品
	HeroEquip     string `xlsx:"heroornamentsid"` // 英雄装备
}

type ChestDataConfig struct {
	Id          string `xlsx:"id"`           // 宝箱id
	Ty          int    `xlsx:"type"`         // 类型
	Des         string `xlsx:"des"`          // 宝箱描述
	OpenOne     int    `xlsx:"open_one"`     // 开宝箱一次消耗
	AcOpenOne   int    `xlsx:"ac_open_one"`  // 活动开宝箱一次消耗
	OpenTen     int    `xlsx:"open_ten"`     // 开宝箱十次消耗
	AcOpenTen   int    `xlsx:"ac_open_ten"`  // 活动开宝箱十次消耗
	TimeLimit   int    `xlsx:"timelimit"`    // 免费开箱倒计时（分）
	LimitCount  int    `xlsx:"limitcount"`   // 连续开箱必得上限
	RewardId    string `xlsx:"rewardid"`     // 连抽必中奖励id
	ActiveStart int    `xlsx:"active_start"` // 活动开始时间
	ActiveEnd   int    `xlsx:"active_end"`   // 活动结束时间
}

type ChestListConfig struct {
	Id       string `xlsx:"id"`       //当前宝箱奖励的物品顺序id
	MatId    string `xlsx:"matid"`    //材料id
	Droprate int    `xlsx:"droprate"` //掉落几率
	MinCount int    `xlsx:"mincount"` //最小数
	MaxCount int    `xlsx:"maxcount"` //最大数
	Type     int    `xlsx:"type"`     //类型
}

type GoodsDataConfig struct {
	Id          string `xlsx:"id"           bson:"id"`          //商品ID
	Ty          int    `xlsx:"type"         bson:"des"`         //商品
	Price       int    `xlsx:"price"        bson:"price"`       //价格
	ActivePrice int    `xlsx:"active_price" bson:"activePrice"` //活动价格
	Count       int    `xlsx:"count"        bson:"count"`       //数量
	Icon        string `xlsx:"icon"         bson:"icon"`        //资源icon
	Discount    int    `xlsx:"discount"     bson:"discount"`    //活动赠送数量
	ActiveStart int    `xlsx:"active_start" bson:"activeStart"` //活动开始时间
	ActiveEnd   int    `xlsx:"active_end"   bson:"activeEnd"`   //活动结束时间
}

type ChallengeStageConfig struct {
	Id       string `xlsx:"id"`
	RewardId string `xlsx:"rewardid"`
}

type ResourceStageConfig struct {
	Id        string `xlsx:"id"`
	RewardIds string `xlsx:"rewardid"`
}

type RoleSkinDataConfig struct {
	Id         string            `xlsx:"id"`
	RelationId string            `xlsx:"relation_id"` // 数据表关联ID
	Attribute  *global.Attribute `xlsx:"-"`           // 属性列表
}

/*
主线任务
每日任务
*/
type TaskDataConfig struct {
	Id          int    `xlsx:"id"`
	LimitTime   int    `xlsx:"limit_time"`      //任务限制时间
	Level       int    `xlsx:"level"`           //等级限制
	IsInitTask  bool   `xlsx:"is_init_task"`    // 是否是初始任务
	IsFirstTask bool   `xlsx:"is_first_task"`   // 是否是某系列的第一个任务
	NextId      int    `xlsx:"next_task"`       //下一个任务ID
	Type        int    `xlsx:"type"`            //任务类型
	Condition   string `xlsx:"condition_data"`  //终结任务条件参数
	Count       int    `xlsx:"condition_count"` //终结任务条件数量
	Reward      string `xlsx:"reward_id"`       //奖励物品ID
}
