package mgoDB

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DbBase struct {
	DbName       string
	CollName     string
}

func NewDbBase(DbName string, collName string) *DbBase {
	return &DbBase{
		DbName:       DbName,
		CollName: collName,
	}
}

func (s *DbBase) CollectionName() string {
	return s.CollName
}

func (s *DbBase) FindOne(ctx context.Context, filter interface{}, recode interface{}, opts ...*options.FindOneOptions) error {
	return GetMgo(s.DbName).GetCol(s.CollName).FindOne(ctx, filter, opts...).Decode(recode)
}

// 批量查询
func (s *DbBase) FindAll(ctx context.Context,  filter interface{}, recode interface{}, opts ...*options.FindOptions) error {
	cursor, err := GetMgo(s.DbName).GetCol(s.CollName).Find(ctx, filter, opts...)
	if err != nil {
		if cursor != nil {
			_ = cursor.Close(ctx)
		}
		return err
	}
	defer cursor.Close(ctx)
	return cursor.All(ctx, recode)
}


// 单个创建
func (s *DbBase) InsertOne(ctx context.Context,  record interface{}) error {
	_, err := GetMgo(s.DbName).GetCol(s.CollName).InsertOne(ctx, record)
	return err
}

// 单个删除
func (s *DbBase) DeleteOne(ctx context.Context,  filter interface{}, opts ...*options.DeleteOptions) (int64, error) {
	result, err := GetMgo(s.DbName).GetCol(s.CollName).DeleteOne(ctx, filter, opts...)
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, err
}

// 更新一条
func (s *DbBase) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (int64, error) {
	result, err := GetMgo(s.DbName).GetCol(s.CollName).UpdateOne(ctx, filter, bson.M{"$set": update}, opts...)
	if err != nil {
		return 0, err
	}
	return result.ModifiedCount, nil
}

// 数据自增长
func (s *DbBase) IncOne(ctx context.Context, filter interface{}, inc interface{}, opt ...*options.UpdateOptions) (int64, error) {
	update := bson.M{"$inc": inc}
	result, err := GetMgo(s.DbName).GetCol(s.CollName).UpdateOne(ctx, filter, update, opt...)
	if err != nil {
		return 0, err
	}
	return result.ModifiedCount, nil
}

// 批量创建
func (s *DbBase) InsertMany(ctx context.Context, cols []interface{}) (int, error) {
	result, err := GetMgo(s.DbName).GetCol(s.CollName).InsertMany(ctx, cols)
	if err != nil {
		return 0, err
	}
	return len(result.InsertedIDs), nil
}

// 批量删除
func (s *DbBase) DeleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (int64, error) {
	result, err := GetMgo(s.DbName).GetCol(s.CollName).DeleteMany(ctx, filter, opts...)
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, err
}

// 批量更新
func (s *DbBase) UpdateMany(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (int64, error) {
	update = bson.M{"$set": update}
	result, err := GetMgo(s.DbName).GetCol(s.CollName).UpdateMany(ctx, filter, update, opts...)
	if err != nil {
		return 0, err
	}
	return result.ModifiedCount, nil
}

// 聚合查询
func (s *DbBase) Aggregate(ctx context.Context, pipeline bson.A, records interface{}, opts ...*options.AggregateOptions) error {
	cursor, err := GetMgo(s.DbName).GetCol(s.CollName).Aggregate(ctx, pipeline, opts...)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)
	return cursor.All(ctx, records)
}

// 对符合条件的集合文档数量进行统计
func (s *DbBase) CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	return GetMgo(s.DbName).GetCol(s.CollName).CountDocuments(ctx, filter, opts...)
}

//(快速计数) 使用集合元数据返回集合中文档数的估计值
func (s *DbBase) CountEstimateDocuments(ctx context.Context, opts ...*options.EstimatedDocumentCountOptions) (int64, error) {
	return GetMgo(s.DbName).GetCol(s.CollName).EstimatedDocumentCount(ctx, opts...)
}

// 批量更新
func (s *DbBase) PushOne(ctx context.Context, filter interface{}, push interface{}, opts ...*options.UpdateOptions) (int64, error) {
	push = bson.M{"$push": push}
	result, err := GetMgo(s.DbName).GetCol(s.CollName).UpdateOne(ctx, filter, push, opts...)
	if err != nil {
		return 0, err
	}
	return result.MatchedCount, nil
}
