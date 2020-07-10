// https://aceld.gitbooks.io/go/grpc/kuai-su-kai-59cb-zhun-bei.html

package main

import (
	"fmt"
	"net"

	pb "go_lession/gRPC_test/my_rpc_proto"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":18881"
)

type server struct{}

//实现RPC SayHello 接口
func (this *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "hello" + in.Name}, nil
}

//实现RPC GetHelloMsg 接口
func (this *server) GetHelloMsg(ctx context.Context, in *pb.HelloRequest) (*pb.HelloMessage, error) {
	return &pb.HelloMessage{Msg: "this is from server HAHA!"}, nil
}

func main() {

	listen, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("failed to listen : ", err)
		return
	}

	//得到一个gRPC 服务句柄
	srv := grpc.NewServer()

	//将 gRPC服务句柄 和 我们的server结构体绑定
	pb.RegisterHelloServerServer(srv, &server{})

	//注册服务
	reflection.Register(srv)

	//启动监听gRPC服务
	if err := srv.Serve(listen); err != nil {
		fmt.Println("failed to serve, ", err)
		return
	}

}
