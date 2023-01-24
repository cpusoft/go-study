package main

import (
	"context"
	proto "grpc2/proto"
	"log"

	"google.golang.org/grpc"
)

func main() {
	// 创建一个 gRPC channel 和服务器交互
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Dial failed:%v", err)
	}
	defer conn.Close()

	// 创建客户端
	client := proto.NewGreeterClient(conn)

	// 直接调用
	resp1, err := client.SayHello(context.Background(), &proto.HelloRequest{
		Name: "Hello Server 1 !!",
	})

	log.Printf("Resp1:%+v", resp1)

	resp2, err := client.SayHello(context.Background(), &proto.HelloRequest{
		Name: "Hello Server 2 !!",
	})

	log.Printf("Resp2:%+v", resp2)
}
