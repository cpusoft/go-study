package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	proto "grpc2/proto"
	"io/ioutil"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

/*
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
*/
func main() {
	cert, err := tls.LoadX509KeyPair("../cert/clienttlscrt.cer", "../cert/clienttlskey.pem")
	if err != nil {
		panic(err)
	}
	// 将根证书加入证书池
	certPool := x509.NewCertPool()
	bs, err := ioutil.ReadFile("../cert/catlsroot.cer")
	if err != nil {
		panic(err)
	}

	if !certPool.AppendCertsFromPEM(bs) {
		panic("fail to append test ca")
	}

	// 新建凭证
	// ServerName 需要与服务器证书内的通用名称一致
	transportCreds := credentials.NewTLS(&tls.Config{
		//	ServerName:   "server.razeen.me",
		Certificates: []tls.Certificate{cert},
		RootCAs:      certPool,
	})

	dialOpt := grpc.WithTransportCredentials(transportCreds)

	conn, err := grpc.Dial("rpki-rtr-server.zdns.cn:8080", dialOpt)
	if err != nil {
		log.Fatalf("Dial failed:%v", err)
	}
	defer conn.Close()

	client := proto.NewGreeterClient(conn)
	resp1, err := client.SayHello(context.Background(), &proto.HelloRequest{
		Name: "Hello Server 1 !!",
	})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Resp1:", resp1)

	resp2, err := client.SayHello(context.Background(), &proto.HelloRequest{
		Name: "Hello Server 2 !!",
	})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Resp2:", resp2)
}
