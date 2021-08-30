package define

const (
	MgoDBNameBattle = "battle"


	RedisGoodsListKey  string = "goods_list"

	RedisUserIndexMap string = "user_id_index" // userId的索引
	// 维护  User Id 自增和唯一性 其他 待定
	RedisIncrementKeyMap  string = "increment_key"      // 主Key
	RedisUserIncrementKey string = "user_increment_key" // 查询自增的user Id 使用的二级Key

	RedisUserDataPrefix              string = "user_data_"               // 后面会拼接一个用户的UserID
	REDIS_USER_INFO                  string = "user_info"                // 用户基本信息json string
	RedisUserShopDataPrefix          string = "user_shop_data_"          // 后面会拼接一个用户的UserID
	RedisUserStrengthGoodsDataPrefix string = "user_strength_shop_data_" // 后面会拼接一个用户的UserID



	TABLE_NAME_PRV_RANKING string = "prv_ranking"
	TABLE_NAME_PUB_RANKING string = "pub_ranking"

	TableNameUser             = "users"
	TableNameRole             = "role"
	TableNameStage            = "stage"
	TableNameStageLog         = "stage_log"
	TableNameHero             = "hero"
	TableNameMail             = "mail"
	TableNameBackPack         = "backpack"
	TableNameSystemMail       = "system_mail"
	TableNameGiftCode         = "gift_code"
	TableNameGiftCodeLog      = "gift_code_log"
	TableNameNotice           = "notice"
	TableNameChestLog         = "chest_log"
	TableNamePurchaseGoodsLog = "purchase_goods_log"
	PurchaseRecode            = "purchase_recode"
	TableNameEquipment        = "equipment"
	TableNameBroadcast        = "broadcast"

	TableNameTask = "task"
)
