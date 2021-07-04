package g

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"time"
)

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

func Client() {
	// grpc.WithInsecure() 禁用传输安全性
	conn, _ := grpc.Dial("127.0.0.1:8050", grpc.WithInsecure())
	defer conn.Close()
	// 创建一个client
	client := NewUserServiceClient(conn)
	// 调用的是UserServiceClient中的方法

	resp := &Resp{}
	begin := time.Now()
	for i := 0; i < 100; i++ {
		resp, _ = client.Create(context.Background(), &User{Name: "ljq", Age: "25"})
	}

	fmt.Printf("返回结果[%s] 耗时[%v]\n", resp.Message, time.Since(begin)/100)
}