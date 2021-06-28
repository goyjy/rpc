package main

import (
	jsonRpc "json"
	tcpRpc "tcp"
)

func main() {
	tcpRpc.Server()
	tcpRpc.Client()

	jsonRpc.Server()
	jsonRpc.Client()
}

