package db

import (
	"com.xv.admin.server/config"
	"com.xv.admin.server/service/mgoDB"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const MD5Secret = "Game"

type AdminModel struct {
	Id            primitive.ObjectID `json:"-" bson:"-"`
	Username      string `json:"username" bson:"username"`
	Password      string `json:"-" bson:"password"`
	Status        int8   `json:"status" bson:"status"`
	Avatar        string `json:"avatar" bson:"avatar"`
	LastLoginTime int64  `json:"last_login_time" bson:"last_login_time"`
	Role          []int   `json:"role" bson:"role"`
}

func (admin *AdminModel) CollectionName() string {
	return config.CollGmAdmin
}

func (admin *AdminModel) GetId() primitive.ObjectID {
	return admin.Id
}

func (admin *AdminModel) SetId(id primitive.ObjectID) {
	admin.Id = id
}

func (admin *AdminModel) LoginByUser(username, pass string) (*AdminModel, error)  {
	result := new(AdminModel)
	filter := bson.D{
		{"username", username},
		{"password", pass},
	}
	finder := mgoDB.NewOneFinder(admin).Where(filter).Record(result)
	res , err := mgoDB.GetMgo().FindOne(context.TODO(), finder)

	if err != nil {
		return nil, err
	}
	if res {
		return result, nil
	}
	return nil, nil
}

func (admin *AdminModel) Create() error {
	err := mgoDB.GetMgo().InsertOne(nil, admin)
	return err
}

func (admin *AdminModel) GetList(filter bson.D, opt *options.FindOptions) (data []map[string]interface{}, err error) {
	data = []map[string]interface{}{}
	finder := mgoDB.NewFinder(admin).Where(filter).Options(opt).Records(&data)
	err = mgoDB.GetMgo().FindMany(context.TODO(), finder)
	return
}
func (admin *AdminModel) GetTotalNum(filter bson.D) (int64, error) {
	finder := mgoDB.NewCounter(admin).Where(filter)
	return mgoDB.GetMgo().CountDocuments(nil, finder)
}


func (admin *AdminModel) Update(filter bson.D, update bson.D) (int64, error) {
	updater := mgoDB.NewUpdater(admin).Where(filter).Update(update)
	res, err := mgoDB.GetMgo().UpdateOne(nil, updater)
	return res,err
}

func (admin *AdminModel) GetRole(d bson.D) ([]int, error) {
	result := make(map[string][]int)
	finder := mgoDB.NewOneFinder(admin).Where(d).Options(options.FindOne().SetProjection(bson.M{"role": 1, "_id": 0})).Record(result)
	res, err := mgoDB.GetMgo().FindOne(context.TODO(), finder)
	if err != nil {
		return nil, err
	}
	fmt.Println(result)
	if res {
		return result["role"], nil
	}
	return nil, nil
}