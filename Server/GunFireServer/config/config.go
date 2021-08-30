package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

var ENV *Conf

type Server struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Root     string `json:"root"`
	Resource string `json:"resource"`
	LogDir   string `json:"log_dir"`
	Language string `json:"language"`
}

type MongoConf struct {
	Host            []string `json:"host"`
	ReplicaSet      string   `json:"replica_set"`
	DbName          string   `json:"db_name"`
	LogDbName       string   `json:"log_db_name"`
	Username        string   `json:"username"`
	Password        string   `json:"password"`
	MaxPoolSize     uint64   `json:"max_pool_size"`
	MinPoolSize     uint64   `json:"min_pool_size"`
	ConnectTimeout  uint     `json:"connect_timeout"`
	MaxConnIdleTime uint     `json:"max_conn_idle_time"`
}

type RedisConf struct {
	Host        string `json:"host"`
	Port        string `json:"port"`
	Password    string `json:"password"`
	MaxIdle     int    `json:"max_idle"`
	MaxActive   int    `json:"max_active"`
	IdleTimeout int    `json:"idle_timeout"`
	DataBase    int    `json:"data_base"`
}

type Conf struct {
	AppPath     string
	JwtSecret   string     `json:"jwt_secret"`
	LobbyServer string     `json:"lobby_server"`
	Server      Server     `json:"server"`
	Mongo       *MongoConf `json:"mongo"`
	Redis       *RedisConf `json:"redis"`
}

func InitConfig() {
	c := new(Conf)
	c.AppPath = getAppPath()
	f, err := os.Open(c.AppPath + "/config.json")
	defer f.Close()
	if err != nil {
		panic(err)
	}
	r := io.Reader(f)
	if err = json.NewDecoder(r).Decode(c); err != nil {
		panic(err)
	}

	ENV = c
}

func getAppPath() string {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(dir)

	return dir

}
