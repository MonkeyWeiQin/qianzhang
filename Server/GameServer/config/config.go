package config

import (
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"os"
)

type (
	//RpcConf struct {
	//	RpcAddr      string   `json:"RpcAddr"`
	//	RegisterAddr []string `json:"RegisterAddr"`
	//	BasePath     string   `json:"BasePath"`
	//}

	LoginConfig struct {
		Id        string `json:"Id"`
		Host      string `json:"Host"`
		JwtSecret string `json:"JwtSecret"`
		Resource  string `json:"Resource"`
		Mode      string `json:"Mode"`
	}

	GameConf struct {
		Id      string `json:"Id"`
		TCPAddr    string `json:"TCPAddr"`
		MaxMsgLen  int    `json:"MaxMsgLen"`
		MaxConnNum int    `json:"MaxConnNum"`
		SecretKey  string `json:"SecretKey"`
	}

	MongoConfig struct {
		Key             string   `json:"Key"` // 必须与常量定义中的Key定义,以选中要操作的数据库
		Url             []string `json:"Url"`
		ReplicaSet      string   `json:"ReplicaSet"`
		DbName          string   `json:"DbName"`
		Username        string   `json:"Username"`
		Password        string   `json:"Password"`
		MaxPoolSize     uint64   `json:"MaxPoolSize"`
		MinPoolSize     uint64   `json:"MinPoolSize"`
		ConnectTimeout  uint     `json:"ConnectTimeout"`
		MaxConnIdleTime uint     `json:"MaxConnIdleTime"`
	}
	RedisConfig struct {
		Host        string `json:"Host"`
		Password    string `json:"Password"`
		MaxIdle     int    `json:"MaxIdle"`
		MaxActive   int    `json:"MaxActive"`
		IdleTimeout int    `json:"IdleTimeout"`
		DataBase    int    `json:"DataBase"`
	}

	// 只定义基础结构公共结构,根据业务在各个模块中转化为需要的的数据结构
	Config struct {
		LoginConf map[string]*LoginConfig           `json:"Login"`
		GameConf  map[string]*GameConf              `json:"Game"`
		MongoConf []*MongoConfig                    `json:"Mongo"`
		RedisConf *RedisConfig                      `json:"Redis"`
		Log       map[string]interface{}            `json:"Log"`
		Settings  map[string]map[string]interface{} `json:"Settings"`
		Dev       map[string]string                 `json:"Dev"`
		Debug     bool                              `json:"Debug"`
	}
)

func LoadConfigFile(path string) *Config {
	f, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	conf := new(Config)
	err = jsoniter.Unmarshal(data, conf)
	if err != nil {
		panic(err)
	}
	return conf
}
