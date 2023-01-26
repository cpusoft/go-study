package server

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

// port:":8080"
func InitGrpcTcpServer(port string) (net.Listener, *grpc.Server, error) {
	// 监听本地端口
	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("InitGrpcTcpServer(): Listen fail,", err)
		return nil, nil, err
	}
	// 创建gRPC服务器
	grpcServer := grpc.NewServer()
	return listener, grpcServer, nil
}

func StartGrpcTcpServer(listener net.Listener, grpcServer *grpc.Server) error {
	reflection.Register(grpcServer)
	fmt.Println("StartGrpcTcpServer():start")
	err := grpcServer.Serve(listener)
	if err != nil {
		fmt.Println("StartGrpcTcpServer(): Serve fail,", err)
		return err
	}
	return nil
}

func InitGrpcTlsServer(port string) (net.Listener, *grpc.Server, error) {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("InitGrpcTlsServer(): Listen fail,", err)
		return nil, nil, err
	}

	// 加载证书和密钥 （同时能验证证书与私钥是否匹配）
	cert, err := tls.LoadX509KeyPair("../cert/servertlscrt.cer", "../cert/servertlskey.pem")
	if err != nil {
		fmt.Println("InitGrpcTlsServer(): LoadX509KeyPair fail,", err)
		return nil, nil, err
	}

	// 将根证书加入证书池
	// 测试证书的根如果不加入可信池，那么测试证书将视为不可惜，无法通过验证。
	certPool := x509.NewCertPool()
	rootBuf, err := ioutil.ReadFile("../cert/catlsroot.cer")
	if err != nil {
		fmt.Println("InitGrpcTlsServer(): ReadFile fail,", err)
		return nil, nil, err
	}

	if !certPool.AppendCertsFromPEM(rootBuf) {
		fmt.Println("InitGrpcTlsServer(): AppendCertsFromPEM fail,", err)
		return nil, nil, errors.New("append certs fail")
	}

	tlsConf := &tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{cert},
		ClientCAs:    certPool,
	}

	serverOpt := grpc.Creds(credentials.NewTLS(tlsConf))
	grpcServer := grpc.NewServer(serverOpt)
	return listener, grpcServer, nil
}

func StartGrpcTlsServer(listener net.Listener, grpcServer *grpc.Server) error {
	err := grpcServer.Serve(listener)
	fmt.Println("StartGrpcTlsServer():start")
	if err != nil {
		fmt.Println("StartGrpcTlsServer(): Listen fail,", err)
		return err
	}
	return nil
}
