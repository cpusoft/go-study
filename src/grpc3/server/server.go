package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	proto "grpc2/proto"
	"io/ioutil"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type server struct {
	proto.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *proto.HelloRequest) (*proto.HelloReply, error) {
	fmt.Println("receive from client:" + in.Name)
	return &proto.HelloReply{Message: "server reply: hi, " + in.Name}, nil
}

/*
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
*/
func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	// 加载证书和密钥 （同时能验证证书与私钥是否匹配）
	cert, err := tls.LoadX509KeyPair("../cert/servertlscrt.cer", "../cert/servertlskey.pem")
	if err != nil {
		panic(err)
	}

	// 将根证书加入证书池
	// 测试证书的根如果不加入可信池，那么测试证书将视为不可惜，无法通过验证。
	certPool := x509.NewCertPool()
	rootBuf, err := ioutil.ReadFile("../cert/catlsroot.cer")
	if err != nil {
		panic(err)
	}

	if !certPool.AppendCertsFromPEM(rootBuf) {
		panic("fail to append test ca")
	}

	tlsConf := &tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{cert},
		ClientCAs:    certPool,
	}

	serverOpt := grpc.Creds(credentials.NewTLS(tlsConf))
	grpcServer := grpc.NewServer(serverOpt)

	proto.RegisterGreeterServer(grpcServer, &server{}) // &proto.SayHelloServer{})

	log.Println("Server Start...")
	grpcServer.Serve(lis)
}
