package main

import (
	"context"
	"fmt"
	proto "grpc2/proto"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	proto.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *proto.HelloRequest) (*proto.HelloReply, error) {
	return &proto.HelloReply{Message: "hello " + in.Name}, nil
}
func main() {
	// 监听本地端口
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Printf("监听端口失败: %s", err)
		return
	}
	// 创建gRPC服务器
	s := grpc.NewServer()
	// 注册服务
	proto.RegisterGreeterServer(s, &server{})
	reflection.Register(s)
	err = s.Serve(lis)
	if err != nil {
		fmt.Printf("开启服务失败: %s", err)
		return
	}
}
