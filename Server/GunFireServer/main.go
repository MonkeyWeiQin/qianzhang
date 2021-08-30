package main

import (
	"com.xv.admin.server/config"
	"com.xv.admin.server/core"
	"com.xv.admin.server/model/db"
	"com.xv.admin.server/service/redisDB"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
)


func init() {
	// 检查数据库
	// 1 索引集合是否存在

}


func main() {
	// 加载配置文件
	config.InitConfig()
	// 创建server
	gin.SetMode(gin.DebugMode)
	r := gin.Default()

	// 启动内核
	core.Start(r)

	increment_key, err := redis.Bool(redisDB.Client.Exists(db.REDIS_INCREMENT_KEY_MAP))
	if err != nil || !increment_key {
		panic(fmt.Sprintf("err:%v, increment_key:%v", err, increment_key))
	}

	// 启动端口监听
	r.Run(fmt.Sprintf("%s:%s", config.ENV.Server.Host, config.ENV.Server.Port))


}
