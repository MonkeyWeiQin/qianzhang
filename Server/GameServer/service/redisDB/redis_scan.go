package redisDB

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

// redis Hash 表 HMGET 命令,返回多个参数,映射到结构体
// args redis.Args
func HMGEToStructByArgs(args redis.Args, data []interface{}, re ...interface{}) error {
	if len(re) == 0 {
		return fmt.Errorf("not struct!! ")
	}
	for i, val := range args {
		if i%2 == 1 {
			if i == 1 {
				data = append([]interface{}{[]byte(val.(string))}, data...)
			} else {
				data[i-1] = []byte(val.(string))
			}
		}
	}
	data = data[:len(data)-1]
	if len(re) == 1 {
		err := redis.ScanStruct(data, re[0])
		if err != nil {
			return err
		}
	} else {
		for i := 0; i < len(re); i++ {
			err := redis.ScanStruct(data, re[i])
			if err != nil {
				return err
			}
		}
	}
	return nil
}
