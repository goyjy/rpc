package jsonRpc

import (
	"fmt"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"time"
)

var client *rpc.Client

type JsonHandle int

func Server() {
	rpcHandle := new(JsonHandle)
	rpc.Register(rpcHandle)

	listen, _ := net.Listen("tcp", ":8080")
	go func() {
		for {
			conn, _ := listen.Accept()
			go jsonrpc.ServeConn(conn)
		}
	}()
}

type Args struct {
	In string
}

type Reply struct {
	Out string
}

func (r JsonHandle) RpcFunc(args *Args, reply *Reply) error {
	reply.Out = "json rpc reply : " + args.In
	return nil
}

func NewClient() {
	client, _ = jsonrpc.Dial("tcp", "127.0.0.1:8080")
}

func JRpcTest() {
	for i := 0; i < 10; i++ {
		go func() {
			var args = Args{In:"json test"}
			var reply Reply
			begin := time.Now()
			client.Call("JsonHandle.RpcFunc", &args, &reply)
			fmt.Printf("返回结果[%s] 耗时[%v]\n", reply.Out, time.Since(begin))
		}()
	}
}