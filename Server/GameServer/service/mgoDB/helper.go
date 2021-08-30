package mgoDB

import "context"

type IStore interface {
	// Object
	FindOne(ctx context.Context, storeType int, key interface{}, x interface{}) error
	FindAll(ctx context.Context, storeType int, keyName string, keyValue interface{}) (map[string]interface{}, error)
	UpdateOne(ctx context.Context, storeType int, k interface{}, x interface{}, immediately ...bool) error
	UpdateFields(ctx context.Context, storeType int, k interface{}, fields map[string]interface{}, immediately ...bool) error
	DeleteOne(ctx context.Context, storeType int, k interface{}, immediately ...bool) error
	DeleteFields(ctx context.Context, storeType int, k interface{}, fields []string, immediately ...bool) error

	// deprecated
	PushArray(ctx context.Context, storeType int, k interface{}, arrayName string, x interface{}) error
	PullArray(ctx context.Context, storeType int, k interface{}, arrayName string, xKey interface{}) error
	UpdateArray(ctx context.Context, storeType int, k interface{}, arrayName string, xKey interface{}, fields map[string]interface{}) error
	SaveHashObjectFields(storeType int, k interface{}, field interface{}, x interface{}, fields map[string]interface{}) error
	SaveHashObject(storeType int, k interface{}, field interface{}, x interface{}) error
	DeleteObjectFields(storeType int, k interface{}, x interface{}, fields []string) error
	DeleteHashObject(storeType int, k interface{}, field interface{}) error
	DeleteHashObjectFields(storeType int, k interface{}, field interface{}, x interface{}, fields []string) error
}

