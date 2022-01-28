package main

import (
	ser "./grpc/mrpc"
	"time"
)

//mrpc测试
func main() {
	ser.InitRedis()

	server := ser.NewServer()

	//
	ser.InitSubManager(server)

	ser.Subscribe("c1")

	ser.Publish("c1", "testwfw")

	time.Sleep(1 * time.Minute)
}
