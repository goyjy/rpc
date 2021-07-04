package main

import (
	"fmt"
	g "grpc"
	jsonRpc "json"
	"os"
	"strconv"
	tcpRpc "tcp"
	"time"
)

func main() {
	args := os.Args

	caps, _ := strconv.Atoi(args[1])

	fmt.Printf("caps [%v]\n", caps)

	g.Server()
	g.NewClient()

	jsonRpc.Server()
	jsonRpc.NewClient()

	tcpRpc.Server()
	tcpRpc.NewClient()

	time.Sleep(time.Second * 3)

	//go g.GRpcRequest(caps)
	//go jsonRpc.JRpcTest(caps)
	go tcpRpc.TRpcTest(caps)

	select {

	}
}



