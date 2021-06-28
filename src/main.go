package main

import (
	httpRpc "http"
	jsonRpc "json"
	tcpRpc "tcp"
)

func main() {
	httpRpc.Server()
	httpRpc.Client()

	tcpRpc.Server()
	tcpRpc.Client()

	jsonRpc.Server()
	jsonRpc.Client()
}

