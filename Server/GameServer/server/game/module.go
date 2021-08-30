package game

import (
	"battle_rabbit/config"
	"battle_rabbit/define"
	"battle_rabbit/iface"
	"battle_rabbit/network"
	"battle_rabbit/service/cache"
	"battle_rabbit/service/log"
	"fmt"
	"sync"
	"time"
)

var (
	expirationTime  = time.Minute * 10
	cleanupInterval = time.Minute * 10
)

func GetGameModel() *Game {
	return new(Game)
}

type Game struct {
	App     iface.IApp
	NodeId  string
	gate    *network.Gate
	conf    *config.GameConf
	players *cache.Cache
	mx      sync.RWMutex
	hotData *cache.Cache // 热数据,1:商店list 2 ...
}

func (g *Game) GetModuleType() string {
	return "Game"
}

func (g *Game) OnInit(app iface.IApp, nodeId string) {
	conf, ok := app.GetConfig().GameConf[nodeId]
	if !ok {
		log.Fatal("game config is nil !!! ")
	}

	g.gate = network.NewGate(conf.TCPAddr)
	g.App = app
	g.NodeId = nodeId
	g.players = cache.New(expirationTime, cleanupInterval)
	g.hotData = cache.New(expirationTime, time.Minute*5)

	// 加载商店列表数据到redis
	g.LoadShopListToRedis()
	g.RegisterRouter()
}

func (g *Game) OnRun() {
	g.gate.Start()
}

func (g *Game) OnStop() {
	log.Debug("Game Stop ....")
	g.gate.OnStop()
}

func (g *Game) RegisterRouter() {
	// 心跳
	g.gate.Register(define.HeartbeatMsgId, g.Heartbeat)

	// 账户
	g.gate.Register(define.CheckTokenMsgId, g.CheckToken)
	g.gate.Register(define.CheckTokenTestMsgId, g.CheckTokenTest)
	g.gate.Register(define.GetPlayerInfoMsgId, g.GetPlayerInfo)
	g.gate.Register(define.UpdateGuideStepMsgId, g.UpdateGuideStep)

	// player
	g.gate.Register(define.GetRoleInfoMsgId, g.GetRoleInfo)
	g.gate.Register(define.GetRoleTalentLvMsgId, g.GetRoleTalentLv)
	g.gate.Register(define.UpgradeRoleTalentLvMsgId, g.UpgradeRoleTalentLv)
	g.gate.Register(define.ChangeRoleSkinMsgId, g.ChangeSkin)
	g.gate.Register(define.UpgradeRoleStarMsgId, g.UpgradeRoleStar)
	//routers.Register(routers.NewRouterInfo(define.UpgradeRoleLevelByLvCardMsgId, true, "Game.UpgradeRoleLevelByLvCard"), g.UpgradeRoleLevelByLvCard)

	// hero
	g.gate.Register(define.GetHeroListMsgId, g.GetHeroList)
	g.gate.Register(define.HeroGoToWarMsgId, g.HeroGoToWar)
	g.gate.Register(define.HeroCancelGoToWarMsgId, g.HeroCancelGoToWar)
	g.gate.Register(define.HeroUpgradeMsgId, g.HeroUpgrade)

	// 装备
	g.gate.Register(define.GetEquipmentListMsgId, g.GetEquipmentList)
	g.gate.Register(define.ChangeEquipmentMsgId, g.ChangeEquipment)
	g.gate.Register(define.CancelEquipmentMsgId, g.CancelEquipment)
	g.gate.Register(define.EquipmentUpgradeMsgId, g.EquipmentUpgrade)

	//背包
	g.gate.Register(define.GetBackpackDataMsgId, g.GetBackpackData)
	g.gate.Register(define.UseExpCardMsgId, g.UseExpCard)

	// stage
	g.gate.Register(define.ActiveStageConfMsgId, g.GetActiveStageConf)
	g.gate.Register(define.GetUserStageRecordMsgId, g.GetUserStageRecord)
	g.gate.Register(define.StartPlayStageMsgId, g.StartPlayStage)
	g.gate.Register(define.CreatePlayerStageMsgId, g.SettlementStage)
	g.gate.Register(define.ReceiveChapterMsgId, g.ReceiveChapterReward) // 章节通关奖励

	// shop
	g.gate.Register(define.GetShopListMsgId, g.GetShopList)
	g.gate.Register(define.GetFirstPurchaseMsgId, g.GetFirstPurchase)
	g.gate.Register(define.PurchaseGoodsMsgId, g.PurchaseGoods)

	// notice
	g.gate.Register(define.GetNoticeMsgId, g.GetNotice)
	g.gate.Register(define.GetBroadCastMsgId, g.GetBroadCast)

	// 邮件
	g.gate.Register(define.GetPlayerMailListMsgId, g.GetPlayerMailList)
	g.gate.Register(define.SetPlayerMailReadMsgId, g.SetPlayerMailRead)
	g.gate.Register(define.DelPlayerMailsMsgId, g.DelPlayerMails)
	g.gate.Register(define.ReceiveMailGiftsMsgId, g.ReceiveMailGifts)
	g.gate.Register(define.ReceiveMailGiftMsgId, g.ReceiveMailGift)
	g.gate.Register(define.CreatePlayerMailMsgId, g.CreatePlayerMail)

	// 宝箱
	g.gate.Register(define.FreeOpenChestMsgId, g.FreeOpenChest)
	g.gate.Register(define.OneTimesOpenChestMsgId, g.OneTimesOpenChest)
	g.gate.Register(define.TenTimesOpenChestMsgId, g.TenTimesOpenChest)
	g.gate.Register(define.GetFreeTimeMsgId, g.GetFreeTime)

	// 签到
	g.gate.Register(define.SevenDaySignInMsgId, g.SevenDaySignIn)
}

func (g *Game) GetPlayer(uid int) *Player {
	i, ok := g.players.GetWithFlushExpiration(fmt.Sprintf("user_%d", uid), expirationTime)
	if ok {
		return i.(*Player)
	}
	p := NewPlayer(uid, g)
	err := p.LoadDataFromDB()
	if err != nil {
		return nil
	}
	g.AddPlayer(uid, p)
	return p
}
func (g *Game) AddPlayer(uid int, p *Player) {
	g.players.Set(fmt.Sprintf("user_%d", uid), p, expirationTime)
}
