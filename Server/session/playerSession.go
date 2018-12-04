package session

import (
	"MagicseaServer/gameproto/msgs"
	"github.com/AsynkronIT/protoactor-go/actor"
)

type PlayerSession struct {
	userInfo      *msgs.UserBaseInfo
	gatePid       *actor.PID  //gate 服地址
	agentPid      *actor.PID  //agent 对象地址
	gamePlayerPid *actor.PID  //player 对象地址

	key string //动态生成密码
}

// 踢下线
func (p *PlayerSession) Kick(reason string,ignServer msgs.ServerType){
	msg := &msgs.Kick{Uid:p.userInfo.Uid,Reason:reason}
	if ignServer != msgs.ST_GateServer && p.agentPid != nil{
		p.agentPid.Tell(msg)
	}
	if ignServer != msgs.ST_GateServer && p.gamePlayerPid!=nil{
		p.gamePlayerPid.Tell(msg)
	}
}