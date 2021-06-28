package httpRpc

import (
	"fmt"
	"net/http"
	"net/rpc"
	"time"
)

type HttpHandle int

func Server() {
	rpcHandle := new(HttpHandle)
	rpc.Register(rpcHandle)
	rpc.HandleHTTP()
	go http.ListenAndServe(":8060", nil)
}

type Args struct {
	In string
}

type Reply struct {
	Out string
}

func (r HttpHandle) RpcFunc(args *Args, reply *Reply) error {
	reply.Out = "http rpc reply : " + args.In
	return nil
}

func Client() {
	var args = Args{In:"http test"}
	var reply Reply

	var begin, end time.Time
	begin = time.Now()
	client, _ := rpc.DialHTTP("tcp", "127.0.0.1:8060")
	client.Call("HttpHandle.RpcFunc", &args, &reply)
	end = time.Now()
	fmt.Printf("返回结果[%s] 耗时[%v]\n", reply.Out, end.Sub(begin))
}
