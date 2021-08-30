package define

//type uint32 uint32

// Client -> Server
const (
	HeartbeatMsgId uint32 = 100001 // 心跳

	SevenDaySignInMsgId uint32 = 101030 //签到
	GetNoticeMsgId      uint32 = 101050 //获取公告
	GetBroadCastMsgId   uint32 = 101060 //获取跑马灯

	GetPlayerMailListMsgId uint32 = 102010 //获取玩家邮件
	SetPlayerMailReadMsgId uint32 = 102020 //设置邮件已读
	DelPlayerMailsMsgId    uint32 = 102030 //批量删除邮件
	ReceiveMailGiftsMsgId  uint32 = 102040 //批量领取邮件附件
	ReceiveMailGiftMsgId   uint32 = 102050 //单个领取邮件附件
	CreatePlayerMailMsgId  uint32 = 102060 //增加邮件

	CreatePlayerStageMsgId  uint32 = 104010 //更新用户最新关卡
	StartPlayStageMsgId     uint32 = 104030 //玩家开始进入普通关卡
	ActiveStageConfMsgId    uint32 = 104040 // 获取已经开启的活动关卡
	GetUserStageRecordMsgId uint32 = 104050 // 玩家关卡记录
	ReceiveChapterMsgId     uint32 = 104060 // 领取章节奖励

	CheckTokenMsgId          uint32 = 106010 //通过token绑定认证并绑定这个连接
	CheckTokenTestMsgId      uint32 = 106011 //测试 通过token绑定认证并绑定这个连接
	GetPlayerInfoMsgId       uint32 = 106020 //获取玩家信息
	GetRoleInfoMsgId         uint32 = 106030 //获取主角信息
	GetRoleTalentLvMsgId     uint32 = 106040 //获取天赋技能等级
	UpgradeRoleTalentLvMsgId uint32 = 106050 //升级天赋技能等级
	ChangeRoleSkinMsgId      uint32 = 106060 //玩家更改角色皮肤
	UpdateGuideStepMsgId     uint32 = 106070 //跟新玩家新手引导步数

	//UpgradeRoleLevelByLvCardMsgId uint32 = 106080 //角色/账号 加经验或者升级

	GetEquipmentListMsgId uint32 = 107010 //获取玩家所有装备列表
	ChangeEquipmentMsgId  uint32 = 107020 //切换装备
	CancelEquipmentMsgId  uint32 = 107030 //取消装备
	EquipmentUpgradeMsgId uint32 = 107040 //装备升级/升星

	GetPlayerSkinListMsgId uint32 = 108010 //获取玩家皮肤列表
	ChangePlayerSkinMsgId  uint32 = 108020 //更换皮肤
	UpgradeRoleStarMsgId   uint32 = 108030 //角色升星

	GetHeroListMsgId       uint32 = 110020 //获取英雄列表
	HeroGoToWarMsgId       uint32 = 110030 //英雄出战
	HeroCancelGoToWarMsgId uint32 = 110040 //英雄取消出战
	HeroUpgradeMsgId       uint32 = 110050 //英雄升级/升星/主动/被动/装备提升等级操作

	UseGiftBagMsgId uint32 = 120010 //使用兑换码

	FreeOpenChestMsgId     uint32 = 130010 //免费开宝箱
	OneTimesOpenChestMsgId uint32 = 130020 //付费开一次
	TenTimesOpenChestMsgId uint32 = 130030 //付费开十次
	GetFreeTimeMsgId       uint32 = 130040 //获取免费开宝箱时间

	GetFirstPurchaseMsgId uint32 = 130110 //商品是否第一次购买列表
	PurchaseGoodsMsgId    uint32 = 130120 //购买商品
	GetShopListMsgId      uint32 = 130130 //获取商品列表

	GetBackpackDataMsgId uint32 = 140010 // 获取背包数据
	UseExpCardMsgId      uint32 = 140020 // 使用角色经验卡

)

// server -> client  推送
// 推送消息ID 1000 < id < 10000 范围
const (
	PushMsgId1001 = 1001 // 强制下线
	PushMsgId1020 = 1020 // 金币钻石更新
	PushMsgId1021 = 1021 // 角色属性通知
	PushMsgId1022 = 1022 // 主/副武器属性通知
	PushMsgId1023 = 1023 // 体力推送

	PushMsgId1101 = 1101 // 角色升级
	PushMsgId2001 = 2001 // 任务完成提醒
)
