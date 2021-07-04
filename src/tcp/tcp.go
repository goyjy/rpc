package tcpRpc

import (
	"fmt"
	"net"
	"net/rpc"
	"sync"
	"sync/atomic"
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

func TRpcTest() {
	wait := sync.WaitGroup{}
	var sum int32
	for i := 0; i < 100; i++ {
		wait.Add(1)
		go func() {
			var args = Args{In:"tcp test"}
			var reply Reply
			begin := time.Now()
			client.Call("RpcHandle.RpcFunc", &args, &reply)
			//fmt.Printf("返回结果[%s] 耗时[%v]\n", reply.Out, time.Since(begin))
			atomic.AddInt32(&sum, int32(time.Since(begin)))
			wait.Done()
		}()
	}
	wait.Wait()
	sum = sum / 100
	fmt.Println("tRpc耗时 ", time.Duration(sum))
}
