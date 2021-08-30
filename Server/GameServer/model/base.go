package model

//基地系统不开发

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	BreadStuff      = "basemat01" //面包
	WaterStuff      = "basemat02" //水
	MedicalKitStuff = "basemat03" //医疗包
	WoodStuff       = "basemat04" //木
	RockStuff       = "basemat05" //石
	SteelStuff      = "basemat06" //铁
)

//var Stuff = []int{RockStuff, SteelStuff, WoodStuff}

// BaseModel  基地
type BaseModel struct {
	Id            primitive.ObjectID `json:"-"                bson:"-"`
	UserId        int                `json:"user_id" bson:"user_id"`                //用户ID
	IsUnlock      bool               `json:"is_unlock" bson:"is_unlock"`            //是否解锁
	Grade         int                `json:"grade" bson:"grade"`                    //基地等级
	People        int                `json:"people" bson:"people"`                  //人数
	Depot         BaseDepot          `json:"depot" bson:"depot"`                    //仓库
	Institute     BaseInstitute      `json:"institute" bson:"institute"`            //研究所
	BulletinBoard BaseBulletinBoard  `json:"bulletin_board"  bson:"bulletin_board"` //公告栏
	Guard         BaseGuard          `json:"guard"  bson:"guard"`                   //守卫
	Task          BaseTask           `json:"task" bson:"task"`                      //任务
	MedalOfHonor  BaseMedalOfHonor   `json:"medal_of_honor" bson:"medal_of_honor"`  //荣誉勋章
	Buff          BaseBuff           `json:"buff" bson:"buff"`                      //基地buff
}

// BaseDepot 仓库
type BaseDepot struct {
	State                int `json:"state" bson:"state"`                                     //仓库状态    0未解锁  1正常 2损坏
	BreadStuffCount      int `json:"bread_stuff_count" bson:"bread_stuff_count"`             //面包个数
	WaterStuffCount      int `json:"water_stuff_count" bson:"water_stuff_count"`             //水个数
	MedicalKitStuffCount int `json:"medical_kit_stuff_count" bson:"medical_kit_stuff_count"` //医疗包个数
	WoodStuffCount       int `json:"wood_stuff_count" bson:"wood_stuff_count"`               //木个数
	RockStuffCount       int `json:"rock_stuff_count" bson:"rock_stuff_count"`               //石个数
	SteelStuffCount      int `json:"steel_stuff_count" bson:"steel_stuff_count"`             //铁个数
}

// BaseInstitute 研究所
type BaseInstitute struct {
	State           int `json:"state" bson:"state"`                           //状态    0未解锁  1正常 2损坏
	GoodsASalePrice int `json:"goods_a_sale_price" bson:"goods_a_sale_price"` //商品A出售价格  todo 特殊商品
	//todo 改造副武器
}

// BaseBulletinBoard 公告栏
type BaseBulletinBoard struct {
	State              int `json:"state" bson:"state"`                               //状态    0未解锁  1正常 2损坏
	NextCollectionTime int `json:"next_collection_time" bson:"next_collection_time"` //下一次领取时间
}

// BaseGuard 守卫
type BaseGuard struct {
	State                  int `json:"state" bson:"state"`                                         //状态    0未解锁  1正常 2损坏
	LateGoldCollectionTime int `json:"late_gold_collection_time" bson:"late_gold_collection_time"` //最后一次领取时间
}

// BaseTask 任务
type BaseTask struct {
	State        int `json:"state" bson:"state"`                   //状态    0未解锁  1正常 2损坏
	MainTaskId   int `json:"main_task_id" bson:"main_task_id"`     //主线任务ID
	EveryDayTask int `json:"every_day_task" bson:"every_day_task"` //每日任务 todo
}

// BaseMedalOfHonor 荣誉勋章
type BaseMedalOfHonor struct {
	State int `json:"state" bson:"state"` //状态    0未解锁  1正常 2损坏
	Count int `json:"count" bson:"count"` //数量
}

// BaseBuff 基地buff
type BaseBuff struct {
	Buff1 struct {
		State int `json:"state" bson:"state"` //状态    0未解锁  1正常 2损坏  //todo
	}
}

func (base *BaseModel) CollectionName() string {
	return "base"
}

func (base *BaseModel) GetId() primitive.ObjectID {
	return base.Id
}

func (base *BaseModel) SetId(id primitive.ObjectID) {
	base.Id = id
}

func (base *BaseModel) UpdateBase(filter bson.D, update bson.D, opt *options.UpdateOptions) error {
	//updater := mgoDB.NewUpdater(base).Where(filter).Update(update).Options(opt)
	//_, err := mgoDB.GetMgo().UpdateOne(nil, updater)
	//return err
	return nil
}

// IncPeople 增加基地人口
func (base *BaseModel) IncPeople(people int) error {
	base.People = base.People + people
	update := bson.D{{"people", base.People}}
	if base.Grade <= 3 {
		if grade := base.checkBaseGrade(); base.Grade != grade {
			base.Grade = grade
			update = append(update, bson.E{Key: "grade", Value: base.Grade})
		}
	}
	if UpdateBaseErr := base.UpdateBase(bson.D{{"user_id", base.UserId}}, update, nil); UpdateBaseErr != nil {
		return UpdateBaseErr
	}
	return nil
}

func (base *BaseModel) checkBaseGrade() int {
	grade := 1
	switch {
	case base.People > 0 && base.People <= 49:
		grade = 1
	case base.People > 49 && base.People <= 199:
		grade = 2
	case base.People > 199 && base.People <= 399:
		grade = 3
	case base.People > 399:
		grade = 4
	}
	return grade
}

//// Insert 新增数据
//func (base *BaseModel) Insert() error {
//	base.Id = primitive.NewObjectID()
//	err := mgoDB.GetMgo().InsertOne(nil, base)
//	return err
//}
//
//func (base *BaseModel) GetBase(userId int) (err error) {
//	filter := bson.D{{"user_id", userId}}
//	finder := mgoDB.NewOneFinder(base).Where(filter)
//	res, err := mgoDB.GetMgo().FindOne(context.TODO(), finder)
//	if !res { //未找到就新增
//		base.UserId = userId
//		return base.Insert()
//	}
//	return
//}
