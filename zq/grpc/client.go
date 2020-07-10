package main

import (
	"fmt"

	pb "go_lession/gRPC_test/my_rpc_proto"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	address    = "localhost:18881"
	clientName = "GreenHat"
)

type server struct{}

func main() {

	//得到 gRPC 链接客户端句柄
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		fmt.Println("did not connetc : ", err)
		return
	}
	defer conn.Close()

	//将 proto 里面的服务句柄 和 gRPC句柄绑定
	c := pb.NewHelloServerClient(conn)

	//远程调用 SayHello接口
	r1, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: clientName})
	if err != nil {
		fmt.Println("cloud not get Hello server ..", err)
		return
	}

	fmt.Println("HelloServer resp: ", r1.Message)

	//远程调用 GetHelloMsg接口
	r2, err := c.GetHelloMsg(context.Background(), &pb.HelloRequest{Name: clientName})
	if err != nil {
		fmt.Println("cloud not get hello msg ..", err)
		return
	}

	fmt.Println("HelloServer resp: ", r2.Msg)
}
