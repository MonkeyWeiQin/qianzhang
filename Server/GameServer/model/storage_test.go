package model

import (
	"battle_rabbit/config"
	"battle_rabbit/service/mgoDB"
	"github.com/gomodule/redigo/redis"
	"testing"
	"time"
)

var (
	Conn redis.Conn
)

func init() {
	InitMgoDB()
	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			option := []redis.DialOption{
				redis.DialKeepAlive(10 * time.Second),      // Redis服务器的TCP连接的保持活动的时间 10分钟
				redis.DialConnectTimeout(10 * time.Second), // 连接超时时间 10s
				redis.DialReadTimeout(5 * time.Second),     // 读超时 5s
				redis.DialWriteTimeout(5 * time.Second),    // 写超时 5s
				redis.DialDatabase(8),                      // 修改库号
			}

			return redis.Dial("tcp", "192.168.0.108:6379", option...)
		},
	}
	Conn = pool.Get()
}

func InitMgoDB() {
	conf := config.LoadConfigFile("D:\\project\\go\\BattleRabbit\\Server\\GameServer\\bin\\server.json")
	mgoDB.OnInit(conf.MongoConf)
}

func TestBaseModel_SetId(t *testing.T) {

}