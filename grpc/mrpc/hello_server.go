package mrpc

import (
	protos "../protos"
	"context"
	"fmt"
)

type HelloServer struct {
	protos.UnimplementedHelloServiceServer
}

var HelloServerInstacne *HelloServer

func NewHelloServer() *HelloServer {
	HelloServerInstacne = &HelloServer{}

	return HelloServerInstacne
}

func (s *HelloServer) HelloTest(ctx context.Context, in *protos.String) (*protos.String, error) {
	fmt.Println(in)

	return &protos.String{
		Value: "123",
	}, nil
}
