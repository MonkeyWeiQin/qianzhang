package define

type (
	AttributeType int8
	PackageType   int8
	TalentType    int8
	RarityLevel   string
	EquipageType  int8
	ItemType      int8 //物品类型
	//ItemType     int8 //宝箱奖品类型
)

// 其他
//武器类型
const (
	MainWeaponType EquipageType = 1
	SubWeaponType  EquipageType = 2
	ArmorType      EquipageType = 3
	OrnamentType   EquipageType = 4
	HeroEquipType  EquipageType = 5

	DefaultWeaponMainId = "weapon_01_1"
	DefaultWeaponSubId  = "sub_01_1"
	DefaultArmorId      = "weapon_01_1"
	DefaultOrnamentId   = "weapon_01_1"

	// 天赋类型
	TalentHpType       TalentType = 1  // 生命加成类型
	TalentAttackType   TalentType = 2  // 攻击加成类型
	TalentAttSpeedType TalentType = 3  // 攻速加成类型
	TalentDefType      TalentType = 4  // 防御加成类型
	TalentViolenceType TalentType = 5  // 暴击加成类型
	TalentGoldType     TalentType = 6  // 金币加成类型
	TalentBossType     TalentType = 7  // boss伤害加成类型
	TalentMoveType     TalentType = 8  // 移动速度加成类型
	TalentDodgeType    TalentType = 9  // 闪避加成类型
	TalentBuffType     TalentType = 10 // buff时长加成类型

	TalentAttackRelationId   = "passive01" // 攻击加成关联ID
	TalentAttSpeedRelationId = "passive02" // 攻速加成关联ID
	TalentViolenceRelationId = "passive03" // 暴击加成关联ID
	TalentMoveRelationId     = "passive04" // 移动速度加成关联ID
	TalentBuffRelationId     = "passive05" // buff时长加成关联ID
	TalentDodgeRelationId    = "passive06" // 闪避加成关联ID
	TalentBossRelationId     = "passive07" // boss伤害加成关联ID
	TalentGoldRelationId     = "passive08" // 金币加成关联ID
	TalentHpRelationId       = "passive09" // 生命加成关联ID
	TalentDefRelationId      = "passive10" // 防御加成关联ID
	TalentHpAllRelationId    = "passive11" // 总等级生命加成关联ID

	//CardRoleExpType        PackageType = 1  //主角经验卡
	//CardRoleStarType       PackageType = 2  //主角升星卡
	//CardReformerType       PackageType = 3  //副武器改造卡
	//CardHeroUpgradeType    PackageType = 4  //英雄升级卡
	//CardHeroStarType       PackageType = 6  //英雄升星卡
	//CardHeroStrengthenType PackageType = 7  //英雄强化卡
	//CardHeroSkillType      PackageType = 8  //英雄技能卡
	//CardHeroArmorType      PackageType = 9  //英雄装备卡
	//CardWeaponMainType     PackageType = 10 //主武器升品卡
	//CardWeaponSubType      PackageType = 11 //副武器升品卡
	//CardArmorType          PackageType = 12 //护甲升品卡
	//CardOrnamentType       PackageType = 13 //饰品升品卡
	//
	//CrystalGreen  PackageType = 20 // 绿色晶体
	//CrystalBlue   PackageType = 21 // 蓝色晶体
	//CrystalGolden PackageType = 22 // 金色晶体
	//CrystalViolet PackageType = 23 // 紫色晶体
	//CrystalOrange PackageType = 24 // 橙色晶体

	RarityLevelN   RarityLevel = "N"
	RarityLevelR   RarityLevel = "R"
	RarityLevelSR  RarityLevel = "SR"
	RarityLevelSRR RarityLevel = "SRR"

	MaxStarLevel     = 6 // 最大星级
	MaxStrengthenNum = 9 // 最多强化次数

	ItemGoldType       ItemType = 0  //金币
	ItemDiamondType    ItemType = 1  //钻石
	ItemStrengthType   ItemType = 2  //体力值
	ItemExpType        ItemType = 3  //经验值
	ItemHeroType       ItemType = 4  //英雄 对应HeroData
	ItemCardType       ItemType = 5  //卡片 对应CardData
	ItemSubWeaponType  ItemType = 6  //副武器 对应weaponData.RoleSubWeaponDataTable
	ItemMainWeaponType ItemType = 7  //主武器 对应weaponData.RoleMainWeaponDataTable
	ItemArmorType      ItemType = 8  //护甲 对应weaponData.ArmorTable
	ItemOrnamentsType  ItemType = 9  //饰品 对应weaponData.OrnamentsDataTable
	ItemHeroEquipType  ItemType = 10 //英雄饰品 对应 HeroEquipType

	GeneralStageType         = 1 // 普通关卡
	DifficultStageType       = 2 // 困难关卡
	ActiveStageTypeChallenge = 3 // 挑战关卡
	ActiveStageTypeDefense   = 4 // 防守关卡
	ActiveStageTypeGold      = 5 // 金币关卡
	ActiveStageTypeResource  = 6 // 资源关卡
	ActiveStageTypeEndless   = 7 // 无尽模式
	ActiveStageTypeTower     = 8 // 塔防关卡

	TaskGotGold         = 1  // 金币获得
	TaskConsumeGold     = 2  // 金币消耗
	TaskGotDiamond      = 3  // 钻石获得
	TaskConsumeDiamond  = 4  // 钻石消耗
	TaskStageOver       = 5  // 通过指定关卡
	TaskRoleLevel       = 6  // 角色达到指定等级
	TaskAccountLevel    = 7  // 账号达到指定的等级
	TaskContinuousLogin = 8  // 连续登录指定天数
	TaskDesStrength     = 9  // 消耗指定数量的体力
	TaskOpenChest       = 10 // 开宝箱达到指定次数
	TaskGotHero         = 11 // 获得某个英雄
	TaskHeroNum         = 12 // 拥有指定数量的英雄
	TaskHeroLevel       = 13 // 某个英雄达到指定等级
	TaskGotMainWeapon   = 14 // 获得指定的主武器
	TaskMainWeaponNum   = 15 // 拥有指定数量的主武器
	TaskMainWeaponLv    = 16 // 某个主武达到指定等级
	TaskGotSubWeapon    = 17 // 获得指定的副武
	TaskSubWeaponNum    = 18 // 拥有指定数量的副武
	TaskSubWeaponLv     = 19 // 某个副武达到指定等级
	TaskGotArmor        = 20 // 获得指定的护甲
	TaskArmorNum        = 21 // 拥有指定数量的护甲
	TaskArmorLv         = 22 // 某个护甲达到指定等级
	TaskGotOrnament     = 23 // 获得指定的饰品
	TaskOrnamentNum     = 24 // 拥有指定数量的饰品
	TaskOrnamentLv      = 25 // 某个饰品达到指定等级
	TaskGotTalent       = 26 // 解锁某项角色天赋
	TaskTalentLv        = 27 // 某项角色天赋提升到某个等级
	TaskKillBoss        = 28 // 击杀某个boss
	TaskKillBossNum     = 29 // 击杀boss类型的怪物的总数量
	TaskKillEnemy       = 30 // 击杀所有类型的敌人到达指定数量

)
