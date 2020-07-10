package main

import (
	"log"
	pb "xhx/microserver/grpc/proto/hello" // 引入proto包

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	// "google.golang.org/grpc/grpclog"
)

const (
	// Address gRPC服务地址
	Address = "www.zq.com:50052"
)

func main() {
	// 连接
	conn, err := grpc.Dial(Address, grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	// 初始化客户端
	c := pb.NewHelloClient(conn)

	// 调用方法
	req := &pb.HelloRequest{Name: "gRPC"}
	res, err := c.SayHello(context.Background(), req)

	if err != nil {
		log.Fatalln(err)
	}

	log.Println(res.Message)
}
