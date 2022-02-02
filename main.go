package main

import (
	ser "./grpc/mrpc"
	protos "./grpc/protos"
	"fmt"
	"github.com/golang/protobuf/proto"
	"time"
)

//mrpc测试
func main() {
	ser.InitRedis()

	server := ser.NewServer()

	helloServer := ser.NewHelloServer()

	protos.RegisterHelloServiceServer(server, helloServer)

	//
	ser.InitSubManager(server)

	ser.Subscribe("c1") //nolint:errcheck

	message := &protos.String{
		Value: "rpc test",
	}
	marshal, err := proto.Marshal(message)
	if err != nil {
		panic(err)
	}

	//after := time.After(1 * time.Second)
	//sendSum := 0

	//go func() {
	//	duration := time.NewTicker(200* time.Millisecond)
	//
	//
	//	for  {
	//		select{
	//			case <-duration.C:
	//				ser.Publish("c1","proto.HelloService/HelloTest", marshal) //nolint:errcheck
	//				fmt.Println("send c1")
	//				sendSum++
	//			case <-after:
	//				fmt.Println("send sum :",sendSum)
	//				return
	//		}
	//	}
	//}()

	time.Sleep(1 * time.Second)

	ser.Publish("c1", "proto.HelloService/HelloTest", marshal) //nolint:errcheck
	fmt.Println("send c1")

	ser.Publish("c1", "proto.HelloService/HelloTest", marshal)
	fmt.Println("send c1")

	time.Sleep(10 * time.Minute)
}
