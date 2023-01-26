package main

import (
	"context"
	"fmt"
	proto "grpc4/proto"

	server "grpc4/server"
)

type HelloServer struct {
	proto.UnimplementedGreeterServer
}

func (s *HelloServer) SayHello(ctx context.Context, in *proto.HelloRequest) (*proto.HelloReply, error) {
	fmt.Println("receive from client:" + in.Name)
	return &proto.HelloReply{Message: "server reply: hi, " + in.Name}, nil
}

func StartRpcTcpServer() {
	listener, grpcServer, _ := server.InitGrpcTcpServer(":8080")
	// 注册服务
	proto.RegisterGreeterServer(grpcServer, &HelloServer{})
	server.StartGrpcTcpServer(listener, grpcServer)
}
func StartRpcTlsServer() {
	listener, grpcServer, _ := server.InitGrpcTlsServer(":8080")
	proto.RegisterGreeterServer(grpcServer, &HelloServer{}) // &proto.SayHelloServer{})
	server.StartGrpcTlsServer(listener, grpcServer)
}
func main() {
	//StartRpcTcpServer()
	StartRpcTlsServer()
}
