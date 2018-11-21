package gateframework

import (
	"time"
	"github.com/magicsea/ganet/network"
	"github.com/AsynkronIT/protoactor-go/actor"
	"log"
)

type IGateService interface {
	GetAgentActor(Agent) (*actor.PID, error)
}
type Gate struct {
	MaxConnNum      int
	PendingWriteNum int
	MaxMsgLen       uint32
	Processor       network.Processor
	//AgentChanRPC    *chanrpc.Server

	// websocket
	WSAddr      string
	HTTPTimeout time.Duration
	CertFile    string
	KeyFile     string

	// tcp
	TCPAddr      string
	LenMsgLen    int
	LittleEndian bool

	//实例
	wsServer  *network.WSServer
	tcpServer *network.TCPServer
}

func (gate *Gate) Run(gs IGateService) {

	//todo websock


	var tcpServer *network.TCPServer

	if gate.TCPAddr != "" {
		tcpServer = new(network.TCPServer)
		tcpServer.Addr = gate.TCPAddr
		tcpServer.MaxConnNum = gate.MaxConnNum
		tcpServer.PendingWriteNum = gate.PendingWriteNum
		tcpServer.LenMsgLen = gate.LenMsgLen
		tcpServer.MaxMsgLen = gate.MaxMsgLen
		tcpServer.LittleEndian = gate.LittleEndian

		tcpServer.NewAgent = func(conn *network.TCPConn) network.Agent {
			a := &GFAgent{conn: conn, gate: gate, netType: TCP}
			//ab := NewAgentActor(a, pid)
			//gs.Pid.Tell(new(messages.NewChild)) //请求一个actor
			//a.agentActor = <-gs.actorchan
			//a.agentActor.bindAgent = a
			ac, err := gs.GetAgentActor(a)
			if err != nil {
				//todo:应该不会发生吧
				log.Println("NewAgent fail:%v", err.Error())
			}
			a.agentActor = ac
			return a
		}
	}
	//todo websock

	if tcpServer != nil {
		tcpServer.Start()
	}
	gate.tcpServer = tcpServer
}

func (gate *Gate) OnDestroy() {
	//todo websock
	if gate.tcpServer != nil {
		gate.tcpServer.Close()
	}
}
