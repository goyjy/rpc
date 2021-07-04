package jsonRpc

import (
	"fmt"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"sync"
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

func JRpcTest(caps int) {
	wait := sync.WaitGroup{}
	sum := make([]time.Duration, caps)
	for i := 0; i < caps; i++ {
		wait.Add(1)
		go func(in int) {
			var args = Args{In:"json test"}
			var reply Reply
			begin := time.Now()
			client.Call("JsonHandle.RpcFunc", &args, &reply)
			//fmt.Printf("返回结果[%s] 耗时[%v]\n", reply.Out, time.Since(begin))
			sum[in] = time.Since(begin)
			wait.Done()
		}(i)
	}
	wait.Wait()
	var all time.Duration
	for _, t := range sum {
		all += t
	}
	fmt.Printf("jRpc耗时 [%v]\n", all/time.Duration(caps))
}