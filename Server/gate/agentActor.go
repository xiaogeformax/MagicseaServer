package gate

import (
	gfw "MagicseaServer/GAServer/gateframework"
	"github.com/AsynkronIT/protoactor-go/actor"
	"MagicseaServer/gameproto/msgs"
	"MagicseaServer/GAServer/log"
	"github.com/golang/protobuf/proto"
	"errors"
	"MagicseaServer/gameproto"
	"time"
)
type AgentActor struct {
	key string
	verified bool
	bindAgent gfw.Agent
	pid *actor.PID
	parentPid *actor.PID
	baseInfo *msgs.UserBaseInfo
	bindServers []*msgs.UserBindServer
	wantDead bool
}

func NewAgentActor(ag gfw.Agent, parentPid *actor.PID) *AgentActor {
	//创建actor
	ab := &AgentActor{verified: false, bindAgent: ag}
	pid := actor.Spawn(actor.FromInstance(ab))


	ab.pid = pid
	ab.parentPid = parentPid
	log.Println("new agent actor:", pid, "  parent:", parentPid)
	return ab
}

//外部调用tell
func (ab *AgentActor)Tell(msg proto.Message){

}

//收到后端消息
func (ab *AgentActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *msgs.Kick:
		ab.OnStop()
		//todo:not safe
		ab.bindAgent.SetDead() //被动死亡，防止二次关闭
		ab.bindAgent.Close()   //关闭连接
	case *msgs.ClientDisconnect:
		//上报 todo
		if ab.baseInfo != nil {
			//ss := cluster.GetServicePID("session")
			//ss.Tell(&msgs.UserLeave{Uid: ab.baseInfo.Uid, From: msgs.ST_GateServer, Reason: "client disconnect"})
		}
		ab.OnStop()
		context.Self().Stop()

	case *msgs.ReceviceClientMsg:
		//收到客户端消息
		ab.ReceviceClientMsg(msg.Rawdata)
	}
}

//收到前端消息
func (ab *AgentActor) ReceviceClientMsg(data []byte) error {
	pack := new(NetPack)
	if !pack.Read(data) {
		log.Error("AgentActor recv too short:", data)
		return errors.New("AgentActor recv too short")
	}
	//心跳包
	channel := msgs.ChannelType(pack.channel)
	if channel == msgs.Heartbeat {
		ab.SendClientPack(pack)
		return nil
	}
	//认证
	if !ab.verified {
		return ab.CheckLogin(pack)
	}

	//转发
	return ab.forward(pack)
}

//验证消息
func (ab *AgentActor) CheckLogin(pack *NetPack) error {
	log.Info("checklogin....")
	msg := gameproto.PlatformUser{}
	err := proto.Unmarshal(pack.rawData, &msg)
	if err != nil {
		log.Error("CheckLogin fail:%v,msgid:%d", err, pack.msgID)
		return err
	}

	//todo cluster
	//pretime := time.Now()
	//smsg := &msgs.ServerCheckLogin{Uid: uint64(msg.PlatformUid), Key: msg.Key, AgentPID: ab.pid}
	//result, err := cluster.GetServicePID("session").Ask(smsg)

	/*if err == nil {
		checkResult := result.(*msgs.CheckLoginResult)
		if checkResult.Result == msgs.OK {
			//登录成功
			usetime := time.Now().Sub(pretime)
			log.Info("CheckLogin success:%v,time:%v", checkResult, usetime.Seconds())
			ab.baseInfo = checkResult.BaseInfo
			ab.bindServers = checkResult.BindServers
			ab.verified = true
			ab.parentPid.Tell(&msgs.AddAgentToParent{Uid: checkResult.BaseInfo.Uid, Sender: ab.pid})
		} else {
			log.Println("###CheckLogin fail:", checkResult)
		}

		ret := &gameproto.LoginReturn{ErrCode: int32(checkResult.Result), ServerTime: int32(time.Now().Unix())}
		ab.SendClient(msgs.Login, byte(gameproto.S2C_LOGIN_END), ret)

	} else {
		log.Error("CheckLogin error :" + err.Error())
	}*/

	return nil
}

func (ab *AgentActor) SendClientPack(pack *NetPack) {
	data := pack.Write()
	ab.bindAgent.WriteMsg(data)
}


func (ab *AgentActor) OnStop() {
	if ab.verified && ab.baseInfo != nil {
		ab.parentPid.Tell(&msgs.RemoveAgentFromParent{Uid: ab.baseInfo.Uid})
	}
}