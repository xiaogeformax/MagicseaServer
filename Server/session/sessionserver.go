package session

import (
	"MagicseaServer/GAServer/service"
	"fmt"
	"MagicseaServer/gameproto/msgs"
	"reflect"
)

type SessionService struct{
	service.ServiceData

	sessionMgr *SessionManager
}

// Service 获取服务对象
func Service() service.IService{
	return new (SessionService)
}

func Type() string {
	return "session"
}

//接口函数
func (s *SessionService) OnReceive(context service.Context) {
	fmt.Println("session.OnReceive:", context.Message())
}
//接口函数
func (s *SessionService) OnInit() {
	//todo
	//s.sessionMgr = NewSessionManager()
}
//接口函数
func (s *SessionService) OnStart(as *service.ActorService) {
	as.RegisterMsg(reflect.TypeOf(&msgs.CreatePlayerResult{}), s.OnUserCheckLoginGsBack) //二次验证gs返回
	as.RegisterMsg(reflect.TypeOf(&msgs.GetSessionInfo{}), s.GetSessionInfo)
	as.RegisterMsg(reflect.TypeOf(&msgs.GetSessionInfoByName{}), s.GetSessionInfoByName) //查询玩家信息通过名字
	as.RegisterMsg(reflect.TypeOf(&msgs.UserLeave{}), s.OnUserLeave)
}

func (s *SessionService) OnUserCheckLoginGsBack(context service.Context) {
	fmt.Println("SessionService.OnUserCheckLogin:", context.Message())
	cresult := context.Message().(*msgs.CreatePlayerResult)

	//踢掉老玩家
	id:= cresult.BaseInfo.Uid
	oldSession := s.sessionMgr.GetSession(id)
	if oldSession != nil {
		oldSession.Kick("try kick same", msgs.ST_NONE)
		s.sessionMgr.RemoveSession(id)
	}
}