package service

import (
	"testing"
	"github.com/AsynkronIT/protoactor-go/actor"
	"time"
	"fmt"
)

type hello struct {
	Who string
}

func Test1(t *testing.T) {
	fmt.Println("service_test Example pass")
	props := actor.FromInstance(&BaseServer{})
	pid := actor.Spawn(props)
	pid.Tell(&hello{Who: "Roger"})
	time.Sleep(1)
	fmt.Println("service_test Example pass")
	pid.GracefulStop()
}
/*func Example(t *testing.T) {
	fmt.Println("service_test Example pass")
	props := actor.FromInstance(&BaseServer{})
	pid := actor.Spawn(props)
	pid.Tell(&hello{Who: "Roger"})
	time.Sleep(1)
	fmt.Println("service_test Example pass")
	pid.GracefulStop()
}
*/