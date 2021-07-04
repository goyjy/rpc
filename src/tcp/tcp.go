package tcpRpc

import (
	"fmt"
	"net"
	"net/rpc"
	"sync"
	"time"
)

var client *rpc.Client

type RpcHandle int

func Server() {
	rpcHandle := new(RpcHandle)

	listen, _ := net.Listen("tcp", ":8070")
	go func() {
		for {
			conn, _ := listen.Accept()
			serve := rpc.NewServer()
			serve.Register(rpcHandle)
			go serve.ServeConn(conn)
		}
	}()
}

type Args struct {
	In string
}

type Reply struct {
	Out string
}

func (r RpcHandle) RpcFunc(args *Args, reply *Reply) error {
	reply.Out = "tcp rpc reply : " + args.In
	return nil
}

func NewClient() {
	client, _ = rpc.Dial("tcp", "127.0.0.1:8070")
}

func TRpcTest(caps int) {
	wait := sync.WaitGroup{}
	sum := make([]time.Duration, caps)
	for i := 0; i < caps; i++ {
		wait.Add(1)
		go func(in int) {
			var args = Args{In:"tcp test"}
			var reply Reply
			begin := time.Now()
			client.Call("RpcHandle.RpcFunc", &args, &reply)
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
	fmt.Printf("tRpc耗时 [%v]\n", all/time.Duration(caps))
}
