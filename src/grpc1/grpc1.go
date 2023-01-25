package main

import (
	"context"
	"log"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
)

type SayHelloServer struct{}

func (s *SayHelloServer) SayHelloWorld(ctx context.Context, in *HelloWorldRequest) (res *pb.HelloWorldResponse, err error) {
	log.Printf("Client Greeting:%s", in.Greeting)
	log.Printf("Client Info:%v", in.Infos)

	var an *any.Any
	if in.Infos["hello"] == "world" {
		an, err = ptypes.MarshalAny(&HelloWorld{Msg: "Good Request"})
	} else {
		an, err = ptypes.MarshalAny(&Error{Msg: []string{"Bad Request", "Wrong Info Msg"}})
	}

	if err != nil {
		return
	}
	return &HelloWorldResponse{
		Reply:   "Hello World !!",
		Details: []*any.Any{an},
	}, nil
}
