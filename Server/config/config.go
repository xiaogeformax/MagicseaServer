package config

import (
	"MagicseaServer/GAServer/config"

	"io/ioutil"
	"encoding/json"

	"strconv"
	"MagicseaServer/GAServer/log"
)

type Config struct {

	Base         config.ServiceConfig `json:"config"`
	Redis          *RedisConf           `json:"redis"`
	DB           map[string]string    `json:"db"`
	DesignConfig map[string]string    `json:"design"`
	Ver          string               `json:"ver"`
}

type RedisConf struct {
	Addr     string `json:"addr"`
	Password string `json:"password"`
	PoolSize int    `json:"poolsize"`
	DBs      []int  `json:"dbs"`
}

var appConfig *Config

func LoadConfig(confPath string)(*Config,error){
	if data,err := ioutil.ReadFile(confPath);err!=nil{
		return nil,err
	}else{
		var conf = & Config{}
		err:= json.Unmarshal(data, conf)
		appConfig = conf
		return conf,err
	}

}

func (conf *Config) GetDBConfigInt(key string) (int, bool) {
	if v,ok := conf.DB[key];ok {
		i,e := strconv.Atoi(v)
		if e!= nil {
			log.Fatal("GetDBConfigInt err:%v", key)
		}
		return i,true
	}
	return 0,false
}

func GetAppConf() *Config {
	return appConfig
}

func SetConfig(c *Config) {
	appConfig = c
}
