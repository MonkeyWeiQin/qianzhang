package mgoDB

import (
	"com.xv.admin.server/config"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

//type MongoUtils struct {
//	Conn *mongo.Client
//	Db   *mongo.Database
//	Ctx  context.Context
//}

func newClient(conf *config.MongoConf, model readpref.Mode) (*mongo.Client, error) {
	var mgOptions = new(options.ClientOptions)
	mgOptions = mgOptions.SetHosts(conf.Host)
	mgOptions = mgOptions.SetConnectTimeout(time.Duration(conf.ConnectTimeout) * time.Second)
	mgOptions = mgOptions.SetMaxConnIdleTime(time.Duration(conf.MaxConnIdleTime) * time.Second)
	mgOptions = mgOptions.SetMaxPoolSize(conf.MaxPoolSize)
	mgOptions = mgOptions.SetMinPoolSize(conf.MinPoolSize)

	if conf.Username != "" && conf.Password != "" {
		mgOptions = mgOptions.SetAuth(options.Credential{
			Username:    conf.Username,
			Password:    conf.Password,
			AuthSource:  "admin_req",
			PasswordSet: true,
		})
	}

	if conf.ReplicaSet != "" {
		mgOptions = mgOptions.SetReplicaSet(conf.ReplicaSet)
	}

	if model != 0 {
		pref, _ := readpref.New(model)
		mgOptions = mgOptions.SetReadPreference(pref)
	}

	return connectionMgo(mgOptions)
	//
	//Mgo = new(MongoUtils)
	//Mgo.Conn = connectionMgo(mgOptions)
	//Mgo.Db = Mgo.Conn.Database(conf.DbName)
	//
	//GmMgo = new(MongoUtils)
	//GmMgo.Conn = connectionMgo(mgOptions)
	//GmMgo.Db = Mgo.Conn.Database(conf.LogDbName)

}

func connectionMgo(op *options.ClientOptions) (*mongo.Client, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, op)

	if err != nil {
		return nil, err
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		return nil, err
	}
	return client, nil
}

//
//// [切换数据库] 返会一个新的连接,并且这个连接已连接到指定名称的数据库
//func (mgo *MongoUtils) SetDb(dbname string) (*MongoUtils, error) {
//	if mgo.Conn == nil {
//		return nil, errors.New("连接未初始化!!")
//	}
//	var newMgo = new(MongoUtils)
//	newMgo.Conn = mgo.Conn
//	newMgo.Db = newMgo.Conn.Database(dbname)
//	return newMgo, nil
//}
//
//// 获取一个集合的句柄
//func (mgo *MongoUtils) GetCollByName(n string) *mongo.Collection {
//	return mgo.Db.Collection(n)
//}
