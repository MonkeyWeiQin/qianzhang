package model

import (
	"battle_rabbit/define"
	"battle_rabbit/excel"
	"battle_rabbit/service/mgoDB"
	"battle_rabbit/service/redisDB"
	"battle_rabbit/utils/xid"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"strconv"
	"time"
)

var (
	purchaseGoodsLogColl = newPurchaseGoodsLogColl()
)

type PurchaseGoodsLogCollection struct {
	*mgoDB.DbBase
}

func newPurchaseGoodsLogColl() *PurchaseGoodsLogCollection {
	return &PurchaseGoodsLogCollection{mgoDB.NewDbBase(define.MgoDBNameBattle,define.TableNameBroadcast)}
}
func GetPurchaseGoodsLogCollection() *PurchaseGoodsLogCollection {
	return purchaseGoodsLogColl
}
func DefaultUserShopModelRedisKey(uid int) string {
	return fmt.Sprintf("%s_%d", define.RedisUserShopDataPrefix, uid)
}
func DefaultUserStrengthGoodsModelRedisKey(uid int) string {
	return fmt.Sprintf("%s_%d", define.RedisUserStrengthGoodsDataPrefix, uid)
}

type PurchaseGoodsLogModel struct {
	Id          string             `json:"mid" bson:"_id"`
	Uid          int                `json:"uid" bson:"uid"`
	GoodsId      string             `json:"goodsId" bson:"goodsId"`     // 商品ID
	GoodsType    define.ItemType   `json:"goodsType" bson:"goodsType"` // 商品类型
	Time         int                `json:"time" json:"time"`           // 中奖时间
	Price        int                `json:"price" bson:"price"`
	Count        int                `json:"count" bson:"count"`
	Presentation int                `json:"presentation" bson:"presentation"`
}

//CheckIsFirstPurchase 检查是否是第一次购买
func (m *PurchaseGoodsLogCollection) CheckIsFirstPurchase(GoodsId string, uid int) (bool, error) {
	ss, err := GetUserShopRedisFieldsByHMGET(uid, GoodsId)
	if err != nil {
		return false, err
	}
	var PurchaseCount int
	_, err = redis.Scan(ss, &PurchaseCount)
	if err != nil {
		return false,err
	}
	if PurchaseCount > 0 {
		return false, nil
	}
	return true, nil
}

func (m *PurchaseGoodsLogCollection) Create(uid int, GoodsId string, Presentation int) error {
	return m.InsertOne(nil, &PurchaseGoodsLogModel{
		Id:          xid.New().String(),
		Uid:          uid,
		Time:         int(time.Now().Unix()),
		GoodsId:      GoodsId,
		GoodsType:    define.ItemType(excel.GoodsDataConf[GoodsId].Ty),
		Price:        excel.GoodsDataConf[GoodsId].Price,
		Count:        excel.GoodsDataConf[GoodsId].Count,
		Presentation: Presentation,
	})
}

func (m *PurchaseGoodsLogCollection) setPurchaseRedis(uid int, GoodsId string) error {
	ss, err := GetUserShopRedisFieldsByHMGET(uid, GoodsId)
	if err != nil {
		return err
	}
	var PurchaseCount int
	_, err = redis.Scan(ss, &PurchaseCount)
	if err != nil {
		return err
	}
	if _, err := SetUserShopRedisFieldByHMSET(uid, map[string]int{GoodsId: PurchaseCount + 1}); err != nil {
		return err
	}
	return nil
}
func (m *PurchaseGoodsLogCollection) setStrengthGoodsRedis(uid int) error {
	timeKey := strconv.Itoa(time.Now().Year()) + "-" + time.Now().Month().String() + "-" + strconv.Itoa(time.Now().Day())
	StrengthGoods, err := GetUserStrengthGoodsFieldsByHMGET(uid, timeKey)
	if err != nil {
		return err
	}
	_, err = SetUserStrengthGoodsFieldsByHMGET(uid, map[string]int{timeKey: StrengthGoods + 1})
	if err != nil {
		return err
	}
	return nil
}
func (m *PurchaseGoodsLogCollection) CheckStrengthGoods(uid int,vipLv int) (res bool, err error) {
	//val, err := GetUserRedisFieldsByHMGET(uid, "Vip")
	timeKey := strconv.Itoa(time.Now().Year()) + "-" + time.Now().Month().String() + "-" + strconv.Itoa(time.Now().Day())

	//if err != nil {
	//	return
	//}
	//var Vip int8
	//_, err = redis.Scan(val, &Vip)
	//if err != nil {
	//	return
	//}
	StrengthGoods, err := GetUserStrengthGoodsFieldsByHMGET(uid, timeKey)
	if err != nil {
		return
	}
	count := 1
	switch vipLv {
	case 1:
		count = 2
	case 2:
		count = 3
	case 3:
		count = 4
	case 4:
		count = 5
	case 5:
		count = 6
	case 6:
		count = 7
	case 7:
		count = 8
	case 8:
		count = 9
	case 9:
		count = 10
	}
	if count <= StrengthGoods {
		return false, nil
	}
	return true, nil
}

func SetUserShopRedisFieldByHMSET(uid int, m interface{}) (interface{}, error) {
	return redisDB.Client.Hmset(redis.Args{}.Add(DefaultUserShopModelRedisKey(uid)).AddFlat(m)...)
}
func GetUserShopRedisFieldsByHMGET(uid int, keys ...string) (vals []interface{}, err error) {
	vals, err = redis.Values(redisDB.Client.Hmget(redis.Args{}.Add(DefaultUserShopModelRedisKey(uid)).AddFlat(keys)...))
	return
}
func GetUserStrengthGoodsFieldsByHMGET(uid int, time string) (vals int, err error) {
	ss, err := redis.Values(redisDB.Client.Hmget(redis.Args{}.Add(DefaultUserStrengthGoodsModelRedisKey(uid)).AddFlat(time)...))
	if err != nil {
		return 0, err
	}
	var StrengthGoods int
	_, err = redis.Scan(ss, &StrengthGoods)
	return StrengthGoods, err
}
func SetUserStrengthGoodsFieldsByHMGET(uid int, m interface{}) (interface{}, error) {
	return redisDB.Client.Hmset(redis.Args{}.Add(DefaultUserStrengthGoodsModelRedisKey(uid)).AddFlat(m)...)
}

