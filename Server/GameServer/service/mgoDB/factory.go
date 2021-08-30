/*******************************************
		暂时未使用,主要用在多副本集服务器
********************************************/
package mgoDB

import (
	"battle_rabbit/config"
	"battle_rabbit/service/log"
	"context"
)

var (
	mgo *MgoFactory
)

type MgoStoreSet struct {
	secondary *MgoStore
	primary   *MgoStore
}

type MgoFactory struct {
	stores map[string]*MgoStoreSet
}

// 创建一个读/写连接
func OnInit(configs []*config.MongoConfig) {
	mgo = &MgoFactory{
		stores: make(map[string]*MgoStoreSet),
	}
	for _, conf := range configs {
		mgo.stores[conf.Key] = &MgoStoreSet{
			secondary: createSecondary(conf),
			primary:   createPrimary(conf),
		}
	}
}

// 获取store
// key：数据库关键字,配置文件中的Key
// secondary：读模式是否为secondary模式
func GetStore(key string, secondary bool) *MgoStore {
	if secondary {
		return mgo.stores[key].secondary
	}
	return mgo.stores[key].primary
}

func  GetMgo(dbName string) *MgoStore {
	return GetStore(dbName, false)
}

func GetMgoSecondary(dbName string) *MgoStore {
	return GetStore(dbName, true)
}

func (f *MgoFactory) Exit(ctx context.Context) {
	for k, set := range f.stores {
		if set.secondary != nil {
			set.secondary.Exit(ctx)
		}
		if set.primary != nil {
			set.primary.Exit(ctx)
		}
		delete(f.stores, k)
	}
}

// 添加secondary模式的store
// conf：数据库配置文件
func createSecondary(conf *config.MongoConfig) *MgoStore {
	cli, err := NewQueryClient(conf)
	if err != nil {
		log.Painc(err)
	}
	dbStore := NewStore(cli, conf.DbName)
	return dbStore
}

// 添加primary模式的store
// conf：数据库配置文件
func createPrimary(conf *config.MongoConfig) *MgoStore {
	cli, err := NewClient(conf)
	if err != nil {
		log.Painc(err)
	}
	dbStore := NewStore(cli, conf.DbName)
	return dbStore
}
