package main

import (
	"fmt"
	"net"
	
	pb "xhx/microserver/grpc/proto/hello"
	
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
)

const = (
	Address = ":50052"
)

type helloService struct{}

var HelloService = helloService{}

func (h helloService) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	resp := new(pb.HelloResponse)
	resp.Message = fmt.Sprintf("Hello %s.", in.Name)
	
	return resp, nil
}

func main() {
	listen, err := net.Listen("tcp", Address)
	if err != nil {
		grpclog.Fatalf("Failed to listen: %v", err)
	}
	
	var opts []grpc.ServerOption
	
	// TLS认证
	creds, err := credentials.NewServerTLSFromFile("../../keys/server.pem", "../../keys/server.key")
	if err != nil {
		grpclog.Fatalf("Failed to generate credentials %v", err)
	}
	
	opts = append(opts, grpc.Creds(creds))
	
	// 注册interceptor
	opts = append(opts, grpc.UnaryInterceptor(interceptor))
	
	// 实例化grpc server
	s := grpc.NewServer(opts...)
	
	// 注册HelloService
	pb.RegisterHelloServer(s, HelloService)
	
	grpclog.Println(fmt.Sprintf("Listen on %s with TLS + Token + Interceptor", Address))
	
	s.Serve(listen)
}

// auth 验证token
func auth(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return grpc.Errorf(codes.Unauthenticated, "无token认证信息")
	}
	
	var (
		appid string
		appkey string
	)
	if val, ok := md["appid"]; ok {
		appid = val[0]
	}
	if val, ok := md["appkey"]; ok {
		appkey = val[0]
	}
	if appid != "101010" || appkey != "i am key" {
		return grpc.Errorf(codes.Unauthenticated, "Token认证信息无效：appid=%s, appkey=%s", appid, appkey)
	}
	
	return nil
}

func interceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler)(interface{}, error) {
	err := auth(ctx)
	if err != nil {
		return nil, err
	}
	// 继续处理请求
	return handler(ctx, req)
}