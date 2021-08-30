package db

const (
	REDIS_USER_INDEX_MAP string = "user_id_index" // userId的索引
	// 维护  User Id 自增和唯一性 其他 待定
	REDIS_INCREMENT_KEY_MAP  string = "increment_key"      // 主Key
	REDIS_USER_INCREMENT_KEY string = "user_increment_key" // 查询自增的user Id 使用的二级Key

	REDIS_USER_DATA_PREFIX string = "user_data_" // 后面会拼接一个用户的UserID
	REDIS_USER_INFO        string = "user_info"  // 用户基本信息json string

	TABLE_NAME_PRV_RANKING string = "prv_ranking"
	TABLE_NAME_PUB_RANKING string = "pub_ranking"

	TBALE_NAME_ADVERTISING_VIDEO string = "pub_video"        //广告点视频观看
	TABLE_NAME_ONLINE_TABLE      string = "pub_online_table" // 在线奖励
	TABLE_NAME_SEVEN_DAY         string = "pub_seven_day"    // 七天奖励
	REDIS_USER_TASK              string = "task"             // 任务表
	REDIS_ADMIN_ROLE             string = "admin_role"       // 后台账号权限

)
