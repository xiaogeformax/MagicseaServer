package db

import (
	"github.com/go-redis/redis"
	"MagicseaServer/Server/config"
	"time"
	"MagicseaServer/GAServer/log"
)

var redisMgr *RedisMgr

//redis db用途
type RedisDBUse int

const (
	RedisDBUseGame       RedisDBUse = iota //游戏角色相关数据 0
	RedisDBUseBattleLoad                   //战场负载数据 1
	RedisDBUseBattleInfo                   //战斗相关数据 2
	RedisDBFriend                          //好友相关数据 3
	RedisDBGuild                           //公会相关数据 4
	RedisDBConfig        = 10              //一些及时配置
	RedisDBUseMax
)



type RedisMgr struct {
	clients map[RedisDBUse]*redis.Client
}


func (mgr *RedisMgr) OnInit() bool {
      if config.GetAppConf().Redis != nil{
      	//todo
		  //for _, v := range config.GetAppConf().Redis.DBs {
		  	//if !mgr.NewRe
		  //}
	  }
}

func NewRedisMgr() *RedisMgr {
	r:= new (RedisMgr)
	r.clients = make(map[RedisDBUse]*redis.Client)
	redisMgr = r
	return r
}

func (mgr *RedisMgr) NewRedisClient(dbIndex int, addr string, poolsize int, passsword string) bool {
	redclient := redis.NewClient(&redis.Options{
		Addr:         addr, //":6379",
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     poolsize, //10,
		Password:     passsword,
		PoolTimeout:  30 * time.Second,
		DB:           dbIndex,
	})
	_, err := redclient.Set("test____", 1, time.Second*10).Result()

	if err!=nil {
		log.Error("%v", err)
		return false
	}
	mgr.clients[RedisDBUse(dbIndex)] = redclient
	return true
}