package gate

import (
	gfw "MagicseaServer/GAServer/gateframework"
	"github.com/AsynkronIT/protoactor-go/actor"
	"MagicseaServer/gameproto/msgs"
	"MagicseaServer/GAServer/log"
	"github.com/golang/protobuf/proto"
	gp "github.com/magicsea/ganet/proto"
	"errors"
	"MagicseaServer/gameproto"

	"MagicseaServer/Server/cluster"
	"time"
	"MagicseaServer/GAServer/config"
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
		if ab.baseInfo != nil {
			ss := cluster.GetServicePID("session")
			ss.Tell(&msgs.UserLeave{Uid: ab.baseInfo.Uid, From: msgs.ST_GateServer, Reason: "client disconnect"})
		}
		ab.OnStop()
		context.Self().Stop()

	case *msgs.ReceviceClientMsg:
		//收到客户端消息
		ab.ReceviceClientMsg(msg.Rawdata)
	}
}

func (ab *AgentActor) GetNetPack() NetPack {
	return GetNetPackByConf()
}

func (ab *AgentActor) GetChannelServer(channel int) *actor.PID {
	c := msgs.ChannelType(int(channel) / 100 * 100) //简单对应
	//log.Info("GetChannelServer,%v,%v", channel, c)
	if ab.bindServers == nil {
		return nil
	}
	for _, v := range ab.bindServers {
		//log.Info("try GetChannelServer,%v", v.Channel)
		if v.Channel == c {
			return v.GetPid()
		}
	}
	return nil
}


//收到前端消息
func (ab *AgentActor) ReceviceClientMsg(data []byte) error {
	pack := ab.GetNetPack()

	msgId, rawdata, err := pack.Unmarshal(data)
	if msgId != "b_move" {
		log.Info("recv:%v", msgId)
	}

	if err != nil {
		log.Error("pack.Unmarshal error:%v,%v", data, err)
		return errors.New("AgentActor recv too short")
	}


	//心跳包
	channel := pack.GetChannelType(msgId)
	if channel ==ChannelHeartbeat {
		ab.SendClientPack(msgId, rawdata)
		return nil
	}
	//认证
	if !ab.verified {
		return ab.CheckLogin(msgId, rawdata)
	}

	//转发
	return ab.forward(msgId, rawdata, channel)
}

//验证消息
func (ab *AgentActor) CheckLogin(msgId interface{}, rawdata []byte) error {
	log.Info("checklogin....")
	msg := gameproto.PlatformUser{}
	err := proto.Unmarshal(rawdata, &msg)
	if err != nil {
		log.Error("CheckLogin fail:%v,msgid:%v", err, msgId)
		return err
	}

	//todo cluster
	pretime := time.Now()
	smsg := &msgs.ServerCheckLogin{Uid: uint64(msg.PlatformUid), Key: msg.Key, AgentPID: ab.pid}
	result, err := cluster.GetServicePID("session").Ask(smsg)

	if err == nil {
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
		ab.SendClient(msgId,ret)

	} else {
		log.Error("CheckLogin error :" + err.Error())
	}

	return nil
}
//发送消息到客户端
func (ab *AgentActor) SendClient(msgId interface{},  msg proto.Message ) {
	mdata, err := gp.Marshal(msg)
	if err != nil {
		log.Error("SendClient marshal error:%v", err)
		return
	}
	//log.Info("sendclient:msg%v,data:%d=>%v", pack.msgID, len(pack.rawData), pack.rawData)
	ab.SendClientPack(msgId, mdata)

}

func (ab *AgentActor) SendClientPack(msgId interface{}, rawdata []byte) {
	var pack = ab.GetNetPack()
	data, err := pack.Marshal(msgId, rawdata)
	if err != nil {
		log.Error("SendClientPack marshal error:id=%v,%s", msgId, err.Error())
		return
	}
	ab.bindAgent.WriteMsg(data)
}

//转发
func (ab *AgentActor)  forward(msgId interface{}, rawdata []byte, channel ChannelType) error {
	//channel := pack.channel
	//msgid := pack.msgID
	//test gate
	//if channel == byte(msgs.Shop) {
	//	ab.SendClient(msgs.Shop, byte(msgs.S2C_ShopBuy), &msgs.S2C_ShopBuyMsg{ItemId: 1, Result: msgs.OK})
	//	return nil
	//}

	if msgId != "b_move" {
		log.Info("=========>forward msg:%v", msgId)
	}

	pid := ab.GetChannelServer(int(channel))
	if pid == nil {
		log.Error("forward server nil:%+v,c=%v,m=%v", pid, channel, msgId)
		return nil
	}
	if config.IsJsonProto() {
		frame := &msgs.FrameMsgJson{MsgId: msgId.(string), RawData: rawdata, Uid: ab.baseInfo.Uid}
		pid.Request(frame, ab.pid)
	} else {
		frame := &msgs.FrameMsg{MsgId: uint32(msgId.(byte)), RawData: rawdata, Uid: ab.baseInfo.Uid}
		pid.Request(frame, ab.pid)
	}

	//frame := &msgs.FrameMsg{channel, uint32(msgid), pack.rawData}
	//pid.Tell(frame)
	//r, e := pid.RequestFuture(frame, time.Second*3).Result()
	//if e != nil {
	//	log.Error("forward error:id=%v, err=%v", ab.baseInfo.Uid, e)
	//}

	//rep := r.(*msgs.FrameMsgRep)
	//repMsg := &gameproto.S2C_ConfirmInfo{MsgHead: int32(msgid), Code: int32(rep.ErrCode)}
	//ab.SendClient(msgs.GameServer, byte(gameproto.S2C_CONFIRM), repMsg)
	return nil
}

func (ab *AgentActor) OnStop() {
	if ab.verified && ab.baseInfo != nil {
		ab.parentPid.Tell(&msgs.RemoveAgentFromParent{Uid: ab.baseInfo.Uid})
	}
}