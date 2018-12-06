package session

import (
	"MagicseaServer/GAServer/service"
	"fmt"
	"MagicseaServer/gameproto/msgs"
	"reflect"
	"MagicseaServer/GAServer/log"
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
	as.RegisterMsg(reflect.TypeOf(&msgs.UserLeave{}), s.OnUserLeave)                   //玩家掉线
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

	ss := &PlayerSession{userInfo: cresult.BaseInfo, gatePid: cresult.TransData.GatePID, key: cresult.TransData.Key}
	s.sessionMgr.AddSession(ss)
	ss.agentPid = cresult.TransData.AgentPID
	ss.gamePlayerPid = cresult.PlayerPID

	//发送消息
	// gsValue := msgs.UserBindServer{msgs.GameServer, cresult.GetPlayerPID()}
	// bsValue := msgs.UserBindServer{msgs.BattleServer, cresult.GetRoomPID()}
	// context.Tell(sender, &msgs.CheckLoginResult{
	// 	Result:      msgs.OK,
	// 	BaseInfo:    ss.userInfo,
	// 	BindServers: []*msgs.UserBindServer{&gsValue, &bsValue}})

	log.Info("SessionService.OnUserCheckLogin ok:", id)
}

//查询玩家信息
func (s *SessionService) GetSessionInfo(context service.Context) {
	fmt.Println("SessionService.GetSessionInfo:", context.Message())
	msg := context.Message().(*msgs.GetSessionInfo)
	ss:= s.sessionMgr.GetSession(msg.Uid)
	if ss != nil {
		context.Tell(context.Sender(),&msgs.GetSessionInfoResult{Result:msgs.OK,UserInfo:ss.userInfo,AgentPID:ss.agentPid})
	}else{
		context.Tell(context.Sender(),&msgs.GetSessionInfoResult{Result:msgs.Fail})
	}

}


func GetServiceValue(key string, values []*msgs.ServiceValue) string {
	for _,v :=range values {
		if v.Key == key {
			return v.Value
		}
	}
	return ""
}

//查询玩家信息 by name
func (s *SessionService) GetSessionInfoByName(context service.Context) {
	fmt.Println("SessionService.GetSessionInfoByName:", context.Message())
	msg := context.Message().(*msgs.GetSessionInfoByName)
	ss := s.sessionMgr.GetSessionByName(msg.Name)
	if ss != nil {
		context.Tell(context.Sender(),&msgs.GetSessionInfoResult{Result:msgs.OK,UserInfo:ss.userInfo,AgentPID:ss.agentPid})
	}else{
		context.Tell(context.Sender(),&msgs.GetSessionInfoResult{Result:msgs.Fail})
	}
}

//离线
func (s *SessionService) OnUserLeave(context service.Context) {
	fmt.Println("SessionService.OnUserLeave:", context.Message())
	msg := context.Message().(*msgs.UserLeave)
	//内存移除
	ss:= s.sessionMgr.RemoveSession(msg.Uid)
	//踢人
	if ss!= nil{
		ss.Kick(msg.Reason,msg.From)
	}
}