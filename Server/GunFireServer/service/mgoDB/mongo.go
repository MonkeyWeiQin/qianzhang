package mgoDB

import (
	"com.xv.admin.server/config"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

//type MgoConf struct {
//	User        string
//	Password    string
//	DataSource  []string
//	AuthDB      string
//	DB          string
//	ReplicaSet  string
//	MaxPoolSize uint64
//}

var (
	storge = make(map[string]*MgoStore)
)

type MgoStore struct {
	cli     *mongo.Client
	db      *mongo.Database
	timeOut time.Duration
}

// 创建一个读/写连接
func NewClient(cfg *config.MongoConf) {
	if cfg == nil {
		cfg = config.ENV.Mongo
	}
	cli, err := newClient(cfg, 0)
	if err != nil {
		panic(err)
	}
	storge["mgo"] = NewStore(cli, cfg.DbName)
	storge["gm_mgo"] = NewStore(cli, cfg.LogDbName)
}

func GetMgo() *MgoStore {
	return storge["mgo"]
}
func GetGmMgo() *MgoStore {
	return storge["gm_mgo"]
}

// 创建一个只读连接 (默认数据源是写入的服务器,通过第二个参数指定再其他副本集中读取)
func NewQueryClient(cfg *config.MongoConf) (*mongo.Client, error) {
	return newClient(cfg, readpref.SecondaryPreferredMode)
}

//func newClient2(cfg *config.MongoConf, readPreference string) (*mongo.Client, error) {
//	uri := ""
//	if cfg.Password == "" {
//		uri = fmt.Sprintf("mongodb://%s/%s?",strings.Join(cfg.DataSource, ","), cfg.AuthDB)
//	}else{
//		uri = fmt.Sprintf("mongodb://%s:%s@%s/%s?", cfg.Username, cfg.Password, strings.Join(cfg.DataSource, ","), cfg.AuthDB)
//	}
//	if cfg.ReplicaSet != "" {
//		uri = fmt.Sprintf("%sreplicaSet=%s&", uri, cfg.ReplicaSet)
//	}
//	if readPreference != "" {
//		uri = fmt.Sprintf("%sreadPreference=%s&", uri, readPreference)
//	}
//	uri = strings.TrimRight(strings.TrimRight(uri, "?"), "&")
//	client, err := mongo.NewClient(options.Client().ApplyURI(uri).SetMaxPoolSize(cfg.MaxPoolSize))
//	if err != nil {
//		return nil, err
//	}
//	err = client.Connect(context.Background())
//	if err != nil {
//		return nil, err
//	}
//	ctxTimeout, cancel := context.WithTimeout(context.Background(), 30*time.Second)
//	defer cancel()
//
//	err = client.Ping(ctxTimeout, nil)
//	if err != nil {
//		return nil, err
//	}
//
//	return client, nil
//}

func NewStore(cli *mongo.Client, db string) *MgoStore {
	return &MgoStore{
		cli:     cli,
		db:      cli.Database(db),
		timeOut: time.Second * 5,
	}
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

// 单个查询
func (s *MgoStore) FindOne(ctx context.Context, o *OneFinder) (bool, error) {
	if o == nil || o.col == nil {
		return false, errors.New("oneFinder is invalid")
	}
	result := s.db.Collection(o.col.CollectionName()).FindOne(ctx, o.filter, o.options...)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, result.Err()
	}
	err := result.Decode(o.record)
	if err != nil {
		return false, err
	}
	return true, nil
}


// 批量查询
func (s *MgoStore) FindMany(ctx context.Context, o *Finder) error {
	if o == nil || o.col == nil || o.records == nil {
		return errors.New("finder is invalid")
	}
	cursor, err := s.db.Collection(o.col.CollectionName()).Find(ctx, o.filter, o.options...)
	if err != nil {
		if cursor != nil {
			s.CloseCursor(ctx, cursor)
		}
		return err
	}
	defer s.CloseCursor(ctx, cursor)
	err = cursor.All(ctx, o.records)
	if err != nil {
		return err
	}
	return nil
}

// 单个创建
func (s *MgoStore) InsertOne(ctx context.Context, col ICollection) error {
	if col == nil {
		return errors.New("collection is invalid")
	}
	result, err := s.db.Collection(col.CollectionName()).InsertOne(ctx, col)
	if err != nil {
		return err
	}
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		col.SetId(oid)
	}
	return err
}

// 单个删除
func (s *MgoStore) DeleteOne(ctx context.Context, col ICollection) (int64, error) {
	if col == nil {
		return 0, errors.New("collection is invalid")
	}
	filter := bson.D{
		{"_id", col.GetId()},
	}
	result, err := s.db.Collection(col.CollectionName()).DeleteOne(ctx, filter)
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, err
}

// 更新一条
func (s *MgoStore) UpdateOne(ctx context.Context, o *Updater) (int64, error) {
	if o == nil || o.col == nil || len(o.filter) == 0 || len(o.update) == 0 {
		return 0, errors.New("updater is invalid")
	}
	update := bson.D{
		{"$set", o.update},
	}
	result, err := s.db.Collection(o.col.CollectionName()).UpdateOne(ctx, o.filter, update, o.options...)
	if err != nil {
		return 0, err
	}
	return result.ModifiedCount, nil
}

// 批量创建
func (s *MgoStore) InsertMany(ctx context.Context, cols []interface{}) error {
	if len(cols) == 0 {
		return errors.New("cols is invalid")
	}

	var name string
	if v, ok := cols[0].(ICollection); ok {
		name = v.CollectionName()
	} else {
		return errors.New("cols not implement collection interface")
	}

	result, err := s.db.Collection(name).InsertMany(ctx, cols)
	if err != nil {
		return err
	}

	for i, _ := range result.InsertedIDs {
		if id, ok := result.InsertedIDs[i].(primitive.ObjectID); ok {
			if v, ok := cols[i].(ICollection); ok {
				v.SetId(id)
			}
		}
	}

	return err
}

// 批量删除
func (s *MgoStore) DeleteMany(ctx context.Context, o *Deleter) (int64, error) {
	if o == nil || o.col == nil {
		return 0, errors.New("deleter is invalid")
	}
	result, err := s.db.Collection(o.col.CollectionName()).DeleteMany(ctx, o.filter, o.options...)
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, err
}

// 批量更新
func (s *MgoStore) UpdateMany(ctx context.Context, o *Updater) (int64, error) {
	if o == nil || o.col == nil || len(o.update) == 0 {
		return 0, errors.New("updater is invalid")
	}
	update := bson.D{
		{"$set", o.update},
	}
	result, err := s.db.Collection(o.col.CollectionName()).UpdateMany(ctx, o.filter, update, o.options...)
	if err != nil {
		return 0, err
	}
	return result.ModifiedCount, nil
}

// 聚合查询
func (s *MgoStore) Aggregate(ctx context.Context, o *Aggregator) error {
	if o == nil || o.col == nil || len(o.pipeline) == 0 || o.records == nil {
		return errors.New("aggregator is invalid")
	}
	cursor, err := s.db.Collection(o.col.CollectionName()).Aggregate(ctx, o.pipeline, o.options...)
	if err != nil {
		return err
	}
	defer s.CloseCursor(ctx, cursor)
	err = cursor.All(ctx, o.records)
	if err != nil {
		return err
	}
	return nil
}

// 对符合条件的集合文档数量进行统计
func (s *MgoStore) CountDocuments(ctx context.Context, o *Counter) (int64, error) {
	if o == nil || o.col == nil {
		return 0, errors.New("counter is invalid")
	}
	cnt, err := s.db.Collection(o.col.CollectionName()).CountDocuments(ctx, o.filter)
	if err != nil {
		return 0, err
	}
	return cnt, nil
}

//(快速计数) 使用集合元数据返回集合中文档数的估计值
func (s *MgoStore) CountEstimateDocuments(ctx context.Context, o *EstimateCounter) (int64, error) {
	if o == nil || o.col == nil {
		return 0, errors.New("estimateCounter is invalid")
	}
	cnt, err := s.db.Collection(o.col.CollectionName()).EstimatedDocumentCount(ctx, o.options...)
	if err != nil {
		return 0, err
	}
	return cnt, nil
}

// 增加一条子集
func (s *MgoStore) PushOne(ctx context.Context, o *Updater) (int64, error) {
	if o == nil || o.col == nil || len(o.filter) == 0 || len(o.push) == 0 {
		return 0, errors.New("push is invalid")
	}
	update := bson.D{
		{"$push", o.push},
	}
	result, err := s.db.Collection(o.col.CollectionName()).UpdateOne(ctx, o.filter, update, o.options...)
	if err != nil {
		return 0, err
	}
	return result.ModifiedCount, nil
}

// 自增
func (s *MgoStore) IncOne(ctx context.Context, o *Updater) (int64, error) {
	if o == nil || o.col == nil || len(o.filter) == 0 || len(o.inc) == 0 {
		return 0, errors.New("IncOne is invalid")
	}
	inc := bson.D{
		{"$inc", o.inc},
	}
	result, err := s.db.Collection(o.col.CollectionName()).UpdateOne(ctx, o.filter, inc, o.options...)
	if err != nil {
		return 0, err
	}
	return result.ModifiedCount, nil
}
