package define

// 消息错误码,对应execl状态表的错ID
// 格式 MsgCodexxxx => int

const (
	MsgCode200 = 200 // success
	MsgCode400 = 400 // 参数错误
	MsgCode401 = 401 // 无权限操作
	MsgCode402 = 402 // 登录失败
	MsgCode404 = 404 // 数据未找到
	MsgCode500 = 500 // 加载失败

	MsgCode601 = 601 //暂无可领取的附件
	MsgCode602 = 602 //体力不足
	MsgCode603 = 603 //次数不足
	MsgCode604 = 604 //金币不足
	MsgCode605 = 605 //钻石不足
	MsgCode606 = 606 //数量不足

	MsgCode700 = 700 // 账户永久封停
	MsgCode701 = 701 // 账户已被冻结
	MsgCode702 = 702 // 账号不存在
	MsgCode703 = 703 // 角色已经创建

	MsgCode800 = 800 // 已达到最高等级
	MsgCode801 = 801 // 请先提升等级
	MsgCode802 = 802 // 请先提升品质

	MsgCode1001 = 1001 //已领取过英雄
	MsgCode1002 = 1002 //领取英雄失败
	MsgCode1003 = 1003 //英雄数量上限
	MsgCode1004 = 1004 //已经出战中
	MsgCode1005 = 1005 //未找到英雄数据

	MsgCode1101 = 1101 //附件领取失败
	MsgCode1102 = 1102 //暂无可领取的附件
	MsgCode1103 = 1103 //附件已领取

	MsgCode1201 = 1201 //兑换码已被使用
	MsgCode1202 = 1202 //兑换码已过期
	MsgCode1203 = 1203 //同一类型的兑换码只能使用一次
	MsgCode1204 = 1204 //兑换码失败

	MsgCode1301 = 1301 //还未到免费抽奖时间
	MsgCode1302 = 1302 //购买次数已达上线

	MsgCode1303 = 1303 //当天已经签到

	MsgCode1401 = 1401 // 关卡次数已用完
)
