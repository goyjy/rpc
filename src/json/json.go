package jsonRpc

import (
	"fmt"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"time"
)

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

func Client() {
	var args = Args{In:"json test"}
	var reply Reply
	begin := time.Now()
	client, _ := jsonrpc.Dial("tcp", "127.0.0.1:8080")
	for i := 0; i < 100; i++ {
		client.Call("JsonHandle.RpcFunc", &args, &reply)
	}
	fmt.Printf("返回结果[%s] 耗时[%v]\n", reply.Out, time.Since(begin)/100)
}