package main

import (
	g "grpc"
	jsonRpc "json"
	tcpRpc "tcp"
	"time"
)

func main() {
	g.Server()
	g.NewClient()

	jsonRpc.Server()
	jsonRpc.NewClient()

	tcpRpc.Server()
	tcpRpc.NewClient()

	time.Sleep(time.Second * 3)

	go g.GRpcRequest()
	go jsonRpc.JRpcTest()
	go tcpRpc.TRpcTest()

	time.Sleep(time.Second * 3)
}



