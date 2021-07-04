package g

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
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

func GRpcRequest() {
	for i := 0; i < 10; i++ {
		go func() {
			begin := time.Now()
			resp, err := client.Create(context.Background(), &User{Name: "ljq", Age: "25"})
			if err != nil {
				fmt.Println("请求失败 ", err)
			} else {
				fmt.Printf("返回结果[%s] 耗时[%v]\n",
					resp.Message, time.Since(begin))
			}
		}()
	}


}