package utils

import (
	"battle_rabbit/service/log"
	"battle_rabbit/utils/types"
	"reflect"
	"strings"
	"time"
)

//查找字符是否在数组中
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

// xxxx*5|xxxx*1|xxxx*2
// xxxx : 物品ID
// number : 数量
// * : 数量分隔符
// | : 物品分隔符
func UnmarshalItemsInt(src string, delimiter ...string) (re map[string]int) {
	var (
		limit1 = "|"
		limit2 = "*"
	)
	re = make(map[string]int)
	if len(delimiter) == 1 {
		limit1 = delimiter[0]
	}
	for _, s := range strings.Split(src, limit1) {
		arr := strings.Split(s, limit2)
		if len(arr) == 2 {
			v, err := types.ToInt(arr[1])
			if err != nil {
				log.Error(err)
				return nil
			}
			re[arr[0]] = int(v)
		}
	}
	return re
}

// xxxx*5
func UnmarshalItemsKV(src string, delimiter ...string) (k string, v int) {
	limit := "*"
	if len(delimiter) > 0 {
		limit = delimiter[0]
	}
	arr := strings.Split(src, limit)
	if len(arr) == 2 {
		v, err := types.ToInt(arr[1])
		if err != nil {
			log.Error(err)
			return "", 0
		}
		return arr[0], int(v)
	}
	return
}

// 今天的零晨时间点
func GetMidnightTime() time.Time {
	timeStr := time.Now().Format("2006-01-02")
	ti, _ := time.ParseInLocation("2006-01-02", timeStr, time.Local)
	return ti
}
