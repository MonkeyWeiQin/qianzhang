package redisDB

import (
	"battle_rabbit/config"
	"github.com/gomodule/redigo/redis"
	"time"
)



type RedisHelper struct {
	*redis.Pool
}

var Client = new(RedisHelper)

func OnInit(conf *config.RedisConfig) {
	// 建立连接池
	Pool := &redis.Pool{
		MaxIdle:     conf.MaxIdle,
		MaxActive:   conf.MaxActive,
		IdleTimeout: time.Duration(conf.IdleTimeout) * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			option := []redis.DialOption{
				redis.DialKeepAlive(10 * time.Second),      // Redis服务器的TCP连接的保持活动的时间 10分钟
				redis.DialConnectTimeout(10 * time.Second), // 连接超时时间 10s
				redis.DialReadTimeout(5 * time.Second),     // 读超时 5s
				redis.DialWriteTimeout(5 * time.Second),    // 写超时 5s
				redis.DialDatabase(conf.DataBase), // 修改库号
			}
			if conf.Password != "" {
				option = append(option, redis.DialPassword(conf.Password)) // 设置密码进行连接
			}

			return redis.Dial("tcp", conf.Host, option...)
		},
	}
	Client.Pool = Pool
}

func (this *RedisHelper) Do(cmd string, args ...interface{}) (interface{}, error) {
	pool := this.Get()
	defer pool.Close()

	return pool.Do(cmd, args...)
}
