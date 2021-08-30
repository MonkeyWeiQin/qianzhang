package mgoDB

import (
	"battle_rabbit/config"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"sync"
	"time"
)

type MgoStore struct {
	cli  *mongo.Client
	db   *mongo.Database
	cols map[string]*mongo.Collection
	mtx  sync.RWMutex
}

func NewStore(cli *mongo.Client, db string) *MgoStore {
	return &MgoStore{
		cli:  cli,
		db:   cli.Database(db),
		cols: make(map[string]*mongo.Collection),
	}
}

func (s *MgoStore)Exit(ctx context.Context)  {
	_ = s.cli.Disconnect(ctx)
}


// 创建一个只读连接 (默认数据源是写入的服务器,通过第二个参数指定再其他副本集中读取)
func NewQueryClient(cfg *config.MongoConfig) (*mongo.Client, error) {
	return newClient(cfg, readpref.SecondaryPreferredMode)
}

func NewClient(cfg *config.MongoConfig) (*mongo.Client, error) {
	return newClient(cfg, readpref.PrimaryMode)
}

func (s *MgoStore) GetCol(tableName string) *mongo.Collection {
	s.mtx.RLock()
	if col, ok := s.cols[tableName]; ok {
		s.mtx.RUnlock()
		return col
	}
	s.mtx.RUnlock()

	col := s.db.Collection(tableName)
	s.mtx.Lock()
	s.cols[tableName] = col
	s.mtx.Unlock()
	return col
}
func (s *MgoStore) DB() *mongo.Database {
	return s.db
}

func (s *MgoStore) Cli() *mongo.Client {
	return s.cli
}

func (s *MgoStore) CloseCursor(ctx context.Context, cursor *mongo.Cursor) {
	err := cursor.Close(ctx)
	if err != nil {
		log.Print("CloseCursor", err)
	}
}

// 为集合创建索引
func (s *MgoStore) CreateIndexMany(indexes []Index) error {
	indexModels := make(map[string][]mongo.IndexModel)
	for _, index := range indexes {
		if err := index.Validate(); err != nil {
			return err
		}
		model := mongo.IndexModel{
			Keys: index.Keys,
		}
		opt := options.Index()
		if index.Name != "" {
			opt.SetName(index.Name)
		}
		opt.SetUnique(index.Unique)
		opt.SetBackground(index.Background)

		if index.ExpireAfterSeconds > 0 {
			opt.SetExpireAfterSeconds(index.ExpireAfterSeconds)
		}

		model.Options = opt

		v, ok := indexModels[index.Collection]
		if ok {
			indexModels[index.Collection] = append(v, model)
		} else {
			indexModels[index.Collection] = []mongo.IndexModel{model}
		}
	}

	for collection, index := range indexModels {
		_, err := s.db.Collection(collection).Indexes().CreateMany(context.Background(), index)
		if err != nil {
			return err
		}
	}

	return nil
}

func newClient(conf *config.MongoConfig, model readpref.Mode) (*mongo.Client, error) {
	opt := new(options.ClientOptions).SetHosts(conf.Url)
	opt = opt.SetConnectTimeout(time.Duration(conf.ConnectTimeout) * time.Second)
	opt = opt.SetMaxConnIdleTime(time.Duration(conf.MaxConnIdleTime) * time.Second)
	opt = opt.SetMaxPoolSize(conf.MaxPoolSize)
	opt = opt.SetMinPoolSize(conf.MinPoolSize)

	if conf.Username != "" && conf.Password != "" {
		opt = opt.SetAuth(options.Credential{
			Username:    conf.Username,
			Password:    conf.Password,
			AuthSource:  "admin_req",
			PasswordSet: true,
		})
	}

	if conf.ReplicaSet != "" {
		opt = opt.SetReplicaSet(conf.ReplicaSet)
	}

	if model != 0 {
		pref, _ := readpref.New(model)
		opt = opt.SetReadPreference(pref)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	if client, err := mongo.Connect(ctx, opt); err != nil {
		return nil, err
	}else{
		// Check the connection
		if err = client.Ping(context.TODO(), nil); err != nil {
			return nil, err
		}
		return client, nil
	}
}

