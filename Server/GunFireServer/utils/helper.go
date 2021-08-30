package utils

import (
	"math/rand"
	"reflect"
	"time"
)

// 检查字符串是否在数组内
func InArray(obj interface{}, target interface{}) bool {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true
		}
	}

	return false
}

func RandInt64(min, max int64) int64 {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	rand.Seed(time.Now().UnixNano())
	return rand.Int63n(max-min) + min
}
func RandInt(min, max int) int {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}
