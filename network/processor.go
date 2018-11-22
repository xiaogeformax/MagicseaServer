package network

type Processor interface {
	// must goroutine safe
	Route(msg interface{}, userData interface{}) error
	// must goroutine safe
	Unmarshal(data []byte) (interface{}, error) //解析protobuf
	// must goroutine safe
	Marshal(msg interface{}) ([][]byte, error)
}
