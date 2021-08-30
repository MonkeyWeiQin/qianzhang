package game

import (
	"battle_rabbit/codec"
	"battle_rabbit/define"
	"battle_rabbit/excel"
	"battle_rabbit/global"
	"battle_rabbit/iface"
	"battle_rabbit/model"
	"battle_rabbit/protocol"
	"battle_rabbit/protocol/response"
	"battle_rabbit/service/log"
	jsoniter "github.com/json-iterator/go"
	"go.mongodb.org/mongo-driver/bson"
	"math/rand"
	"time"
)

var (
	chestBoxList = map[string]map[string]*excel.ChestListConfig{
		"treasurebox001": excel.ChestHeroDataConf,
		"treasurebox002": excel.ChestMatDataConf,
		"treasurebox003": excel.ChestEquippedDataConf,
		"treasurebox004": excel.ChestWeaponDataConf,
	}
)

// OpenChest 开宝箱
type OpenChest struct {
	Id string `json:"id"`
}

type ChestFreeTimeResponse struct {
	FreeTime int `json:"freeTime"`
}

func (g *Game) GetFreeTime(sess iface.ISession, msg *codec.Message) {
	player := g.GetPlayer(sess.GetUid())
	if player == nil {
		// 重新登录
		sess.Send(protocol.ErrCode(msg.Id, define.MsgCode401))
		return
	}

	userChest := player.Account.Chest
	chestResponse := map[string]*ChestFreeTimeResponse{}
	for _, config := range excel.ChestDataConf {
		freeTime := 0
		if chest, ok := userChest[config.Id]; ok {
			if (chest.LastTime + config.TimeLimit*60) > int(time.Now().Unix()) {
				freeTime = (chest.LastTime + config.TimeLimit*60) - int(time.Now().Unix())
			}
		}
		chestResponse[config.Id] = &ChestFreeTimeResponse{
			FreeTime: freeTime,
		}
	}
	sess.Send(protocol.SuccessData(msg.Id,chestResponse))
}

// FreeOpenChest 免费开一次宝箱
func (g *Game) FreeOpenChest(sess iface.ISession, msg *codec.Message) {

	var code = define.MsgCode200
	defer func() {
		if code != define.MsgCode200 {
			sess.Send(protocol.ErrCode(msg.Id,code))
		}
	}()

	req := new(OpenChest)
	err := jsoniter.Unmarshal(msg.GetData(), &req)
	if err != nil {
		log.Error("解析请求出错 :", err)
		code = define.MsgCode400
		return
	}

	player := g.GetPlayer(sess.GetUid())
	if player == nil {
		// 重新登录
		code = define.MsgCode401
		return
	}

	if code = checkFree(player, req.Id); code != define.MsgCode200 {
		log.Error("开宝箱失败")
		return
	}

	var chestItems []*global.Item
	chestItems, code = getChestItem(1, req.Id)
	if code != define.MsgCode200 {
		log.Error("获取奖品失败")
		return
	}
	b, _ := jsoniter.Marshal(chestItems)
	log.Debug("chestItems ==== ::: ", string(b))

	var resp *response.AttachmentsResp
	if resp, code = saveItem(g, player, sess, chestItems); code != define.MsgCode200 {
		log.Error("保存奖品失败")

		return
	} else {
		player.Account.FreeOpenChestFlushTime(req.Id)
		if _, err := model.GetUserCollection().SetFreeOpenChestTime(player.Account, req.Id); err != nil {
			log.Error(err)
			code = define.MsgCode500
			return
		}
		saveChestLog(player.Uid, req.Id, 0, 1, chestItems)
		sess.Send(protocol.SuccessData(msg.Id,resp))
		// 开宝箱,任务统计
		UpdateTask(player.Task,sess,[]*model.TaskProgress{{
			Type:      define.TaskOpenChest,
			Num:       1,
		}})
	}
}

// OneTimesOpenChest 付费开一次
func (g *Game) OneTimesOpenChest(sess iface.ISession, msg *codec.Message) {
	var code = define.MsgCode200
	defer func() {
		if code != define.MsgCode200 {
			sess.Send(protocol.ErrCode(msg.Id,code))
		}
	}()

	req := new(OpenChest)
	err := jsoniter.Unmarshal(msg.GetData(), &req)
	if err != nil {
		log.Error("解析请求出错 :", err)
		code = define.MsgCode400
		return
	}
	player := g.GetPlayer(sess.GetUid())
	if player == nil {
		// 重新登录
		sess.Send(protocol.ErrCode(msg.Id, define.MsgCode401))
		return
	}

	ChestInfo := getChestData(req.Id)
	if ChestInfo == nil {
		log.Error("宝箱不存在 :")
		code = define.MsgCode400
		return
	}

	if player.Account.Diamond < ChestInfo.OpenOne {
		log.Warn("钻石不足!!")
		code = define.MsgCode605
		return
	}

	var chestItems []*global.Item
	chestItems, code = getChestItem(1, req.Id)
	if code != define.MsgCode200 {
		log.Error("获取奖品失败")
		return
	}
	player.Account.Diamond -= ChestInfo.OpenOne
	_, err = model.GetUserCollection().UpdateOne(nil, bson.M{"uid": player.Uid}, bson.M{"diamond": player.Account.Diamond})
	if err != nil {
		log.Error("扣除宝石失败::", err)
		code =  define.MsgCode500
		return
	}

	if resp, c := saveItem(g, player, sess, chestItems); c != define.MsgCode200 {
		log.Error("保存奖品失败")
		code = c
		return
	} else {
		saveChestLog(player.Uid, req.Id, ChestInfo.OpenOne, 1, chestItems)
		protocol.NoticeUserGoldAndDiamond(sess, player.Account.Gold, player.Account.Diamond)
		sess.Send(protocol.SuccessData(msg.Id,resp))
		// 开宝箱,任务统计
		UpdateTask(player.Task,sess,[]*model.TaskProgress{{
			Type:      define.TaskOpenChest,
			Num:       1,
		},{
			Type:      define.TaskConsumeDiamond,
			Num:       ChestInfo.OpenOne,
		}})
	}
}

// TenTimesOpenChest 开十次宝箱
func (g *Game) TenTimesOpenChest(sess iface.ISession, msg *codec.Message) {
	var code = define.MsgCode200
	defer func() {
		if code != define.MsgCode200 {
			sess.Send(protocol.ErrCode(msg.Id,code))
		}
	}()

	req := new(OpenChest)
	err := jsoniter.Unmarshal(msg.GetData(), &req)
	if err != nil {
		log.Error("解析请求出错 :", err)
		code = define.MsgCode400
		return
	}
	player := g.GetPlayer(sess.GetUid())
	if player == nil {
		// 重新登录
		code = define.MsgCode401
		return
	}

	ChestInfo := getChestData(req.Id)
	if ChestInfo == nil {
		log.Error("宝箱不存在 :")
		code = define.MsgCode400
		return
	}

	if player.Account.Diamond < ChestInfo.OpenTen {
		log.Warn("宝石不足!!")
		code = define.MsgCode605
		return
	}

	chestItems, code := getChestItem(10, req.Id)
	if code != define.MsgCode200 {
		log.Error("获取奖品失败")
		return
	}
	player.Account.Diamond -= ChestInfo.OpenTen
	_, err = model.GetUserCollection().UpdateOne(nil, bson.M{"uid": player.Uid}, bson.M{"diamond": player.Account.Diamond})
	if err != nil {
		log.Error("扣除宝石失败::", err)
		code =  define.MsgCode500
		return
	}

	if resp, c := saveItem(g, player, sess, chestItems); c != define.MsgCode200 {
		log.Error("保存奖品失败")
		code = c
		return
	} else {
		saveChestLog(player.Uid, req.Id, ChestInfo.OpenOne, 10, chestItems)
		protocol.NoticeUserGoldAndDiamond(sess, player.Account.Gold, player.Account.Diamond)
		sess.Send(protocol.SuccessData(msg.Id,resp))
		// 开宝箱,任务统计
		UpdateTask(player.Task,sess,[]*model.TaskProgress{{
			Type:      define.TaskOpenChest,
			Num:       10,
		},{
			Type:      define.TaskConsumeDiamond,
			Num:       ChestInfo.OpenOne,
		}})
	}
}

//获取宝箱信息
func getChestData(id string) *excel.ChestDataConfig {
	return excel.ChestDataConf[id]
}

//开奖
func getChestItem(count int, id string) ([]*global.Item, int) {
	var chestItemMap = make(map[string]*global.Item)
	var chestItem []*global.Item
	var chestList map[string]*excel.ChestListConfig
	switch id {
	case "treasurebox001":
		chestList = excel.ChestMatDataConf
	case "treasurebox002":
		chestList = excel.ChestWeaponDataConf
	case "treasurebox003":
		chestList = excel.ChestEquippedDataConf
	case "treasurebox004":
		chestList = excel.ChestHeroDataConf
	default:
		return nil, define.MsgCode404
	}
	//chestList := getChestListData(id)
	if chestList == nil || len(chestList) <= 0 {
		log.Error("宝箱数据有误:")
		return nil, define.MsgCode500
	}
	for i := 0; i < count; i++ {
		for _, item := range chestList {
			if checkIsWinALottery(item.Droprate) {
				if _, ok := chestItemMap[item.MatId]; !ok {
					atta := &global.Item{
						Count:    getCount(item.MinCount, item.MaxCount),
						ItemID:   item.MatId,
						ItemType: define.ItemType(item.Type),
					}
					chestItemMap[item.MatId] = atta
					chestItem = append(chestItem, atta)
				} else {
					chestItemMap[item.MatId].Count += getCount(item.MinCount, item.MaxCount)
				}
			}
		}
	}
	return chestItem, define.MsgCode200
}

//计算随机获取个数
func getCount(min int, max int) int {
	if min == max {
		return max
	}
	if min > max { // 配置表不应该出现 最小值 > 最大值的情况
		return 1
	}

	return rand.Intn(max-min+1) + min
}

//检测是否中奖
func checkIsWinALottery(droprate int) bool {
	if droprate >= 100000 {
		return true
	}
	if droprate <= 0 {
		return false
	}
	if droprate >= rand.Intn(100000) {
		return true
	}
	return false
}

//检测是否可以免费抽宝箱
func checkFree(player *Player, id string) int {
	ChestInfo := getChestData(id)
	if ChestInfo == nil {
		log.Error("宝箱不存在 :")
		return define.MsgCode404
	}

	if item, ok := player.Account.Chest[id]; !ok {
		return define.MsgCode200
	} else if int64(item.LastTime+(ChestInfo.TimeLimit*60)) < time.Now().Unix() {
		return define.MsgCode200
	}
	return define.MsgCode1301
}

//保存奖品
func saveItem(g *Game, player *Player, sess iface.ISession, attachment []*global.Item) (*response.AttachmentsResp, int) {
	if resp, flush, err := player.AddAttachments(attachment); err != nil {
		log.Error("添加物品失败 :", err)
		return nil, define.MsgCode500
	} else {
		if flush[0] {
			protocol.NoticeUserGoldAndDiamond(sess, player.Account.Gold, player.Account.Diamond)
		}
		if flush[1] {
			// 推送体力
			protocol.NoticeUserStrengthChange(sess, player.Account.Strength, player.Account.StrengthTime)
		}

		if flush[2] {
			g.flushAttr(player, sess)
		}

		if flush[3] {
			// TODO 刷新角色数据
		}
		if flush[4] {
			push := &protocol.LevelUpgradeNoticeData{Lv: player.Account.Level, Exp: player.Account.Exp, RLv: player.Role.RLv, RExp: player.Role.RExp}
			protocol.NoticeRoleLevelUpgrade(sess, push)
		}
		return resp, define.MsgCode200
	}
}

//保存抽奖记录
func saveChestLog(uid int, ChestId string, Price int, Count int, Attachment []*global.Item) {
	err := model.GetChestLogCollection().Create(uid, ChestId, Price, Count, Attachment)
	if err != nil {
		log.Error(err)
		return
	}
}
