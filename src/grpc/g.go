package g

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"sync"
	"time"
)

var client UserServiceClient

type server struct {

}

func (s *server) Create(ctx context.Context, in *User) (*Resp, error) {
	return &Resp{Message: "gRpc test " + in.Name}, nil
}

func Server() {
	go func() {
		listener, _ := net.Listen("tcp", ":8050")
		newServer := grpc.NewServer()
		RegisterUserServiceServer(newServer, &server{})
		newServer.Serve(listener)
	}()
}

func NewClient() {
	// grpc.WithInsecure() 禁用传输安全性
	conn, _ := grpc.Dial("127.0.0.1:8050", grpc.WithInsecure())
	// 创建一个client
	client = NewUserServiceClient(conn)
	// 调用的是UserServiceClient中的方法
}

func GRpcRequest(caps int) {
	wait := sync.WaitGroup{}
	sum := make([]time.Duration, caps)
	for i := 0; i < caps; i++ {
		wait.Add(1)
		go func(in int) {
			begin := time.Now()
			_, _ = client.Create(context.Background(), &User{Name: "ljq", Age: "25"})
			//fmt.Printf("返回结果[%s] 耗时[%v]\n", resp.Message, time.Since(begin))
			sum[in] = time.Since(begin)
			wait.Done()
		}(i)
	}
	wait.Wait()
	var all time.Duration
	for _, t := range sum {
		all += t
	}
	fmt.Printf("gRpc耗时 [%v]\n", all/time.Duration(caps))
}