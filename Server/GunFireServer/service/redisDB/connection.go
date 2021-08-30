package redisDB

import (
	"com.xv.admin.server/config"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"time"
)

type RedisHelper struct {
	*redis.Pool
}

var Client = new(RedisHelper)

func ConnectRedis() {
	// 建立连接池
	Pool := &redis.Pool{
		MaxIdle:     config.ENV.Redis.MaxIdle,
		MaxActive:   config.ENV.Redis.MaxActive,
		IdleTimeout: time.Duration(config.ENV.Redis.IdleTimeout) * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			option := []redis.DialOption{
				redis.DialKeepAlive(10 * time.Second),      // Redis服务器的TCP连接的保持活动的时间 10分钟
				redis.DialConnectTimeout(10 * time.Second), // 连接超时时间 10s
				redis.DialReadTimeout(5 * time.Second),     // 读超时 5s
				redis.DialWriteTimeout(5 * time.Second),    // 写超时 5s
				redis.DialDatabase(config.ENV.Redis.DataBase), // 修改库号
			}
			if config.ENV.Redis.Password != "" {
				option = append(option, redis.DialPassword(config.ENV.Redis.Password)) // 设置密码进行连接
			}

			return redis.Dial("tcp",
				fmt.Sprintf("%s:%s", config.ENV.Redis.Host, config.ENV.Redis.Port),
				option...,
			)
		},
	}
	Client.Pool = Pool
}

func (this *RedisHelper) Do(cmd string, args ...interface{}) (interface{}, error) {
	pool := this.Get()
	defer pool.Close()

	return pool.Do(cmd, args...)
}
