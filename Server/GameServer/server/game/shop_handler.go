package game

import (
	"battle_rabbit/codec"
	"battle_rabbit/define"
	"battle_rabbit/excel"
	"battle_rabbit/global"
	"battle_rabbit/iface"
	"battle_rabbit/model"
	"battle_rabbit/protocol"
	"battle_rabbit/service/log"
	"battle_rabbit/service/redisDB"
	"github.com/gomodule/redigo/redis"
	jsoniter "github.com/json-iterator/go"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type BuyGoods struct {
	Id string `json:"id"`
}

type ShopList struct {
	Goods  map[string]*excel.GoodsDataConfig `json:"goods"`
	Chests map[string]*excel.ChestDataConfig `json:"chests"`
}

func (g *Game) LoadShopListToRedis() {
	ok, err := redis.Bool(redisDB.Client.Exists(define.RedisGoodsListKey))
	if err != nil {
		log.Fatal("加载商店数据失败", err)
	}
	if !ok {
		shopList := &ShopList{
			Goods:  excel.GoodsDataConf,
			Chests: excel.ChestDataConf,
		}
		list, _ := jsoniter.Marshal(shopList)
		_, err = redisDB.Client.SET(define.RedisGoodsListKey, string(list))
		if err != nil {
			log.Fatal("缓存商店数据到redis失败:", err)
		}
	}
}

func (g *Game) canShopList() *ShopList {
	i, ok := g.hotData.Get(define.RedisGoodsListKey)
	if ok {
		return i.(*ShopList)
	}
	by, err := redis.Bytes(redisDB.Client.GET(define.RedisGoodsListKey))
	if err != nil {
		log.Error(err)
		return nil
	}
	var shopList *ShopList
	err = jsoniter.Unmarshal(by, &shopList)
	if err != nil {
		log.Error(err)
		return nil
	}
	g.hotData.Set(define.RedisGoodsListKey, shopList, time.Minute*5)
	return shopList
}

// 获取商店的列表数据
func (g *Game) GetShopList(sess iface.ISession, msg *codec.Message) {
	sess.Send(protocol.SuccessData(msg.Id, g.canShopList()))
}

//PurchaseGoods 购买商品
func (g *Game) PurchaseGoods(sess iface.ISession, msg *codec.Message) {
	var code = define.MsgCode200
	defer func() {
		if code != define.MsgCode200 {
			sess.Send(protocol.ErrCode(msg.Id, code))
		}
	}()
	req := new(BuyGoods)
	err := jsoniter.Unmarshal(msg.GetData(), &req)
	if err != nil {
		log.Error("解析请求出错 :", err)
		code = define.MsgCode404
		return
	}
	shop := g.canShopList()
	if shop == nil {
		log.Error("获取商品信息失败 :")
		code = define.MsgCode500
		return
	}

	player := g.GetPlayer(sess.GetUid())
	if player == nil {
		code = define.MsgCode401
		return
	}

	goods, ok := shop.Goods[req.Id]
	if !ok {
		log.Error("获取商品信息失败 :")
		code = define.MsgCode500
		return
	}

	now := int(time.Now().Unix())
	count := goods.Count
	price := goods.Price
	presentation := 0
	if now >= goods.ActiveStart && now <= goods.ActiveEnd {
		count += goods.Discount
		presentation = goods.Discount
		if goods.ActivePrice != 0 {
			price = goods.ActivePrice
		}
	}

	switch define.ItemType(goods.Ty) {
	case define.ItemGoldType, define.ItemStrengthType:
		if player.Account.Diamond < goods.Price {
			log.Error("玩家钻石不足 :", err)
			code = define.MsgCode605
			return
		}
		player.Account.Diamond -= price

		r, err := model.GetUserCollection().UpdateOne(nil, bson.M{"uid": player.Uid}, bson.M{"diamond": player.Account.Diamond})
		if err != nil || r == 0 {
			log.Error("购买商品,扣除玩家钻石失败!: ", err, r)
			// 扣除钻石
			code = define.MsgCode500
			return
		}
	case define.ItemDiamondType:
		//todo 第三方支付 (完善的时候单写一个函数)

	}

	attachment := []*global.Item{{ItemType: define.ItemType(goods.Ty), Count: count, ItemID: req.Id}}
	if resp, flush, err := player.AddAttachments(attachment); err != nil {
		log.Error("添加物品失败 :", err)
		code = define.MsgCode500
		return
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
		if model.GetPurchaseGoodsLogCollection().Create(player.Uid, goods.Id, presentation) != nil {
			log.Error("创建记录失败 :")
			code = define.MsgCode500
			return
		}
		sess.Send(protocol.SuccessData(msg.Id, resp))
		return
	}
}

//GetFirstPurchase 获取首次购买列表
func (g *Game) GetFirstPurchase(sess iface.ISession, msg *codec.Message) {
	player := g.GetPlayer(sess.GetUid())
	if player == nil {
		sess.Send(protocol.ErrCode(msg.Id, define.MsgCode401))
		return
	}
	res := make(map[string]bool, len(excel.GoodsDataConf))
	sess.Send(protocol.SuccessData(msg.Id, res))
	return
}
