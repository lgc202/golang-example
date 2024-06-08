package main

import (
	"context"
	"golang-example/grpc/demo03/proto"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	port = ":50051"
)

// greeterServer 定义一个结构体用于实现 .proto文件中定义的方法
// 新版本 gRPC 要求必须嵌入 pb.UnimplementedGreeterServer 结构体
type greeterServer struct {
	proto.UnimplementedGreeterServer
}

// SayHello 简单实现一下.proto文件中定义的 SayHello 方法
func (g *greeterServer) SayHello(ctx context.Context, in *proto.HelloRequest) (*proto.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &proto.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func main() {
	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// 指定使用服务端证书创建一个 TLS credentials。
	creds, err := credentials.NewServerTLSFromFile("../x509/server.crt", "../x509/server.key")
	if err != nil {
		log.Fatalf("failed to create credentials: %v", err)
	}

	s := grpc.NewServer(grpc.Creds(creds))
	// 将服务描述(server)及其具体实现(greeterServer)注册到 gRPC 中去.
	// 内部使用的是一个 map 结构存储，类似 HTTP server。
	proto.RegisterGreeterServer(s, &greeterServer{})
	log.Println("Serving gRPC on 0.0.0.0" + port)
	if err := s.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
