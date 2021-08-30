package model

import (
	"battle_rabbit/define"
	"battle_rabbit/service/mgoDB"
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

// 背包
// 存放一下数据:
// 1 各种升级卡
// 2 宝石
// 3 英雄碎片
// 其他补充

const (
	PackageCardType    = 1
	PackageCrystalType = 2
)

var (
	backpackColl = newBackPackColl()
)

type BackPackColl struct {
	*mgoDB.DbBase
}

func newBackPackColl() *BackPackColl {
	return &BackPackColl{mgoDB.NewDbBase(define.MgoDBNameBattle,define.TableNameBackPack)}
}

func GetBackPackColl() *BackPackColl { return backpackColl }

type IBackPack interface {
	GetType() int
}

// 提升卡
// Uid tid 组成复合唯一索引
type UpgradeCard struct {
	Id    string             `json:"id" bson:"_id"`
	Uid   int                `json:"uid"    bson:"uid"`
	TabId string             `json:"tabId"  bson:"tabId"` // 关联ID
	Ty    define.PackageType `json:"ty"     bson:"ty"`    // 类型
	//Rarity int                `json:"rarity" bson:"rarity"` // 稀有度  1:N 2:R 3:SR 4:SSR (升星卡 专属)
	Num int `json:"num"    bson:"num"` // 数量
}

func (card *UpgradeCard) GetType() int {
	return 1
}

//func (c *BackPackColl) updateOne(ctx context.Context, filter interface{}, update interface{}, options ...*options.UpdateOptions) (int64, error) {
//	return c.UpdateOne(ctx, filter, update, options...)
//}
//func (c *BackPackColl) updateMany(ctx context.Context, filter interface{}, update interface{}, options ...*options.UpdateOptions) (int64, error) {
//	return c.UpdateMany(ctx, filter, update, options...)
//}

func (c *BackPackColl) GetCardByTabId(ctx context.Context, uid int, tabId string) (result *UpgradeCard, err error) {
	err = c.FindOne(ctx, bson.M{"uid": uid, "tabId": tabId}, &result)
	return
}

func (c *BackPackColl) GetCardByMid(ctx context.Context, uid int, mid string) (result *UpgradeCard, err error) {
	err = c.FindOne(ctx, bson.M{"uid": uid, "mid": mid}, &result)
	return
}

// 是个于用户ID和类型组成符合缩影的数据(英雄升星卡不能使用此方法,要区分稀有度!)
func (c *BackPackColl) GetCardByType(ctx context.Context, uid int, ty define.PackageType) (result *UpgradeCard, err error) {
	err = c.FindOne(ctx, bson.M{"uid": uid, "ty": ty}, &result)
	return
}

// AddCardToBackpack 添加卡片到背包
func (c *BackPackColl) AddCardToBackpack(uid int, id string, count int) (err error) {
	if result, err := c.GetCardByTabId(nil, uid, id); err != nil {
		return err
	} else {
		if result != nil {
			one, err := c.IncOne(nil, bson.D{{"uid", uid}, {"tabId", id}}, bson.D{{"num", count}})
			if err != nil || one <= 0 {
				return err
			}
			return err
		} else {
			err := c.InsertOne(nil, &UpgradeCard{
				Uid:   uid,
				TabId: id,
				//Ty:     define.PackageType(excel.CardDataConf[id].Ty),
				//Val:    excel.CardDataConf[id].Exp,
				//Rarity: excel.CardDataConf[id].Rare,
				Num: count,
			})
			if err != nil {
				return err
			}
		}
	}
	return
}

//
////CreateEquipment 添加武器
//func (c *BackPackColl) CreateUpgradeCard(uid int, relationId string) error {
//	return c.InsertOne(nil, &UpgradeCard{
//		Uid:        uid,
//		RelationId: relationId,
//		Level:      1,
//		Quality:    1,
//		Type:       ArmorType(1),
//		Use:        false,
//	})
//}

func (c *BackPackColl) LoadBackpackFromDB(uid int) (map[string]IBackPack, error) {
	cursor, err := mgoDB.GetMgoSecondary(define.MgoDBNameBattle).GetCol(c.CollName).Find(nil, bson.M{"uid": uid})
	if err != nil {
		return nil, err
	}
	var list []*UpgradeCard
	if err = cursor.All(nil,&list); err != nil {
		return nil, err
	}
	pack := make(map[string]IBackPack)
	for _, card := range list {
		pack[card.Id] = card

		// 分类型处理数据
		//if ty, ok := data["ty"].(int32); !ok {
		//	b, _ := jsoniter.Marshal(data)
		//	return nil, fmt.Errorf("背包中的数据有误!无法识别类型 data: %s", string(b))
		//} else {
			//switch int(ty) {
			//case PackageCardType:
			//	card := new(UpgradeCard)
			//	if err = mapstructure.Decode(data, card); err != nil {
			//		return nil, err
			//	}
			//	pack[card.Id] = card
			//case PackageCrystalType:
			//	crystal := new(UpgradeCard)
			//	if err = mapstructure.Decode(data, crystal); err != nil {
			//		return nil, err
			//	}
			//	pack[crystal.Id] = crystal
			//
			//default:
			//	break
			//}
		//}
	}

	return pack, nil

}
