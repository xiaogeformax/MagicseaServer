package game

import (
	"MagicseaServer/GAServer/service"
	"github.com/AsynkronIT/protoactor-go/actor"
	"MagicseaServer/gameproto/msgs"
	"MagicseaServer/GAServer/log"
	"reflect"
	"time"
	"MagicseaServer/Server/cluster"
	"github.com/magicsea/ganet/util"
	"fmt"
	"github.com/go-redis/redis"
	"MagicseaServer/Server/db"
)

type GameService struct {
	service.ServiceData
	isReg bool
}


type PlayerInitEnd struct {
	Result    msgs.GAErrorCode
	BaseInfo  *msgs.UserBaseInfo
	TransData *msgs.CreatePlayer
	Sender    *actor.PID
	RoomPID   *actor.PID
}

//Service 获取服务对象
func Service() service.IService {
	return new(GameService)
}

func Type() string {
	return "game"
}

//以下为接口函数
func (s *GameService) OnReceive(context service.Context) {
	log.Debug("game.OnReceive:", context.Message())

}
func (s *GameService) OnInit() {

}

func (s *GameService) OnStart(as *service.ActorService) {

	as.RegisterMsg(reflect.TypeOf(&msgs.ServerCheckLogin{}), s.OnUserCheckLogin) //二次验证
	as.RegisterMsg(reflect.TypeOf(&msgs.Tick{}), s.OnTick) //定时任务
	as.RegisterMsg(reflect.TypeOf(&msgs.AddServiceRep{}), s.OnRegOK)          //注册完成
	as.RegisterMsg(reflect.TypeOf(&actor.Terminated{}), s.OnDisconnectCenter) //被动断开服务器

}

func (s *GameService) OnRun() {
	//注册到center
	s.RegToCenter()
	//cluster.RegServerWork(&s.ServiceData, nil)
	//定时任务
	util.StartLoopTask(time.Second*5, func() {
		s.Pid.Tell(&msgs.Tick{}) //转主线程执行
	})
}

//注册到中心服务器
func (s *GameService) RegToCenter() {
	//注册到center
	r := cluster.GetServicePID("center")
	msg := msgs.AddService{
		ServiceName: s.Name,
		ServiceType: s.TypeName,
		Pid:         s.GetPID(),
		Values:      nil}
	r.GetActorPID().Request(&msg, s.Pid)
	log.Info("game RegToCenter !!!")
}

//注册成功
func (s *GameService) OnRegOK(context service.Context) {
	s.isReg = true
	log.Info("game reg ok!!!")
	context.Watch(context.Sender())
}


//从中心断开
func (s *GameService) OnDisconnectCenter(context service.Context) {
	s.isReg = false
	log.Info("game OnDisconnectCenter !!!")
}

func (s *GameService) OnTick(context service.Context) {
	if !s.isReg {
		s.RegToCenter()
		return
	}
	load := len(context.Children())
	cluster.UpdateServiceLoad(s.Name, uint32(load), msgs.ServiceStateFree)
}


//玩家验证
func (s *GameService) OnUserCheckLogin(context service.Context) {
	log.Info("GameService.OnUserCheckLogin:", context.Message())
	msg := context.Message().(*msgs.ServerCheckLogin)
	//验证token
	tokenkey := fmt.Sprintf("UserToken:%v_%v", msg.Uid, msg.Key)
	_, err := GetRedisGame().Get(tokenkey).Result()
	if err != nil {
		log.Error("OnUserCheckLogin,no found player,token=%v", tokenkey)
		context.Tell(context.Sender(), &msgs.CheckLoginResult{Result: msgs.KeyError})
		return
	}

	//创建角色
	NewPlayer(msg.Uid, msg.AgentPID, &msgs.CreatePlayer{Uid: msg.Uid, AgentPID: msg.AgentPID, Sender: context.Sender(),
		GatePID: nil}, context) //异步载入完成才发回
	log.Info("GameService.OnUserCheckLogin pre:", msg.Uid)

}

func GetRedisGame() *redis.Client {
	return db.GetRedisGame()
}