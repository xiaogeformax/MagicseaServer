//服务类型，可作为独立进程或线程使用
package service

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"reflect"
	"fmt"
	"MagicseaServer/GAServer/log"
	"github.com/AsynkronIT/protoactor-go/remote"
)

type Context actor.Context

type IService interface {
	IServiceData
	OnReceive(context Context)
	OnInit()
	OnStart(as *ActorService)
	//正式运行(服务线程)
	OnRun()

	OnDestory()
}

type ServiceRun struct {
}

type MessageFunc func(context Context)

//服务的代理
type ActorService struct {
	serviceIns IService
	rounter    map[reflect.Type]MessageFunc
}

func (s *ActorService) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *actor.Started:
		fmt.Println("Started, initialize actor here")
	case *actor.Stopping:
		fmt.Println("Stopping, actor is about shut down")
	case *actor.Stopped:
		fmt.Println("Stopped, actor and its children are stopped")
	case *actor.Restarting:
		fmt.Println("Restarting, actor is about restart")
	case *ServiceRun:
		fmt.Println("ServiceRun ", s.serviceIns.GetName())
		s.serviceIns.OnRun()
	default:
		log.Debug("recv defalult:", msg)
		s.serviceIns.OnReceive(context.(Context))
		fun := s.rounter[reflect.TypeOf(msg)]
		if fun != nil {
		 	fun(context.(Context))
		}
	}
}

func (s *ActorService) RegisterMsg(t reflect.Type, f MessageFunc) {
	s.rounter[t] = f
}

func StartService(s IService) {
	ac := &ActorService{s, make(map[reflect.Type]MessageFunc)}
	props := actor.FromProducer(func() actor.Actor { return ac }) //.WithSupervisor(supervisor)
	if s.GetAddress() != "" {
		remote.Start(s.GetAddress())
	}
	pid, err := actor.SpawnNamed(props, s.GetName())
	if err == nil {
		s.SetPID(pid)
		s.OnStart(ac)
	} else {
		log.Error("#############actor.SpawnNamed error:%v", err)
	}

}
func DestoryService(s *ActorService) {
	s.serviceIns.OnDestory()
}
