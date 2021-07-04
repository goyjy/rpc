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
	g.GRpcRequest()

	time.Sleep(time.Second)

	jsonRpc.Server()
	jsonRpc.NewClient()
	jsonRpc.JRpcTest()

	time.Sleep(time.Second)

	tcpRpc.Server()
	tcpRpc.NewClient()
	tcpRpc.TRpcTest()

	time.Sleep(time.Second)
}



