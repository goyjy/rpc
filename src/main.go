package main

import (
	g "grpc"
	jsonRpc "json"
	tcpRpc "tcp"
	"time"
)

func main() {
	tcpRpc.Server()
	jsonRpc.Server()
	g.Server()

	time.Sleep(time.Second)

	tcpRpc.Client()
	time.Sleep(time.Second)

	jsonRpc.Client()
	time.Sleep(time.Second)

	g.Client()
}

