package main

import (
	"context"
	"fmt"
	"golang-example/grpc/demo05/proto"
	"io"
	"log"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type Echo struct {
	proto.UnimplementedEchoServer
}

// UnaryEcho 一个普通的UnaryAPI
func (e *Echo) UnaryEcho(ctx context.Context, req *proto.EchoRequest) (*proto.EchoResponse, error) {
	log.Printf("Recved: %v", req.GetMessage())
	resp := &proto.EchoResponse{Message: req.GetMessage()}
	return resp, nil
}

//	ServerStreamingEcho 客户端发送一个请求 服务端以流的形式循环发送多个响应
//
// 1. 获取客户端请求参数
// 2. 处理完成后返回过个响应
// 3. 最后返回nil表示已经完成响应
func (e *Echo) ServerStreamingEcho(req *proto.EchoRequest, stream proto.Echo_ServerStreamingEchoServer) error {
	log.Printf("Recved %v", req.GetMessage())
	// 具体返回多少个response根据业务逻辑调整
	for i := 0; i < 2; i++ {
		// 通过 send 方法不断推送数据
		err := stream.Send(&proto.EchoResponse{Message: req.GetMessage()})
		if err != nil {
			log.Fatalf("Send error:%v", err)
			return err
		}
	}

	// 返回nil表示已经完成响应
	return nil
}

// ClientStreamingEcho 客户端流
//
// 1. for循环中通过stream.Recv()不断接收client传来的数据
// 2. err == io.EOF表示客户端已经发送完毕关闭连接了,此时在等待服务端处理完并返回消息
// 3. stream.SendAndClose() 发送消息并关闭连接(虽然在客户端流里服务器这边并不需要关闭 但是方法还是叫的这个名字，内部也只会调用Send())
func (e *Echo) ClientStreamingEcho(stream proto.Echo_ClientStreamingEchoServer) error {
	// 1.for循环接收客户端发送的消息
	for {
		// 2. 通过 Recv() 不断获取客户端 send()推送的消息
		req, err := stream.Recv() // Recv内部也是调用RecvMsg
		// 3. err == io.EOF表示已经获取全部数据
		if err == io.EOF {
			log.Println("client closed")
			// 4.SendAndClose 返回并关闭连接
			// 在客户端发送完毕后服务端即可返回响应
			return stream.SendAndClose(&proto.EchoResponse{Message: "ok"})
		}
		if err != nil {
			return err
		}
		log.Printf("Recved %v", req.GetMessage())
	}
}

// BidirectionalStreamingEcho 双向流服务端
// 1. 建立连接 获取client
// 2. 通过client调用方法获取stream
// 3. 开两个goroutine（使用 chan 传递数据） 分别用于Recv()和Send()
// 3.1 一直Recv()到err==io.EOF(即客户端关闭stream)
// 3.2 Send()则自己控制什么时候Close 服务端stream没有close 只要跳出循环就算close了。 具体见https://github.com/grpc/grpc-go/issues/444
func (e *Echo) BidirectionalStreamingEcho(stream proto.Echo_BidirectionalStreamingEchoServer) error {
	var (
		waitGroup sync.WaitGroup
		msgCh     = make(chan string)
	)

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()

		for v := range msgCh {
			err := stream.Send(&proto.EchoResponse{Message: v})
			if err != nil {
				fmt.Println("Send error:", err)
				continue
			}
		}
	}()

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		for {
			req, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("recv error:%v", err)
			}
			fmt.Printf("Recved :%v \n", req.GetMessage())
			msgCh <- req.GetMessage()
		}
		close(msgCh)
	}()
	waitGroup.Wait()

	// 返回nil表示已经完成响应
	return nil
}

func main() {
	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(grpc.StreamInterceptor(streamInterceptor))
	// 将服务描述(server)及其具体实现(greeterServer)注册到 gRPC 中去.
	// 内部使用的是一个 map 结构存储，类似 HTTP server。
	proto.RegisterEchoServer(s, &Echo{})
	log.Println("Serving gRPC on 0.0.0.0" + port)
	if err := s.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	m, err := handler(ctx, req)
	end := time.Now()
	// 记录请求参数 耗时 错误信息等数据
	fmt.Printf("RPC: %s,req:%v start time: %s, end time: %s, err: %v\n", info.FullMethod, req, start.Format(time.RFC3339), end.Format(time.RFC3339), err)
	return m, err
}

type wrappedStream struct {
	grpc.ServerStream
}

func newWrappedStream(s grpc.ServerStream) grpc.ServerStream {
	return &wrappedStream{s}
}

func (w *wrappedStream) RecvMsg(m interface{}) error {
	fmt.Printf("Receive a message (Type: %T) at %s\n", m, time.Now().Format(time.RFC3339))
	return w.ServerStream.RecvMsg(m)
}

func (w *wrappedStream) SendMsg(m interface{}) error {
	fmt.Printf("Send a message (Type: %T) at %v\n", m, time.Now().Format(time.RFC3339))
	return w.ServerStream.SendMsg(m)
}

func streamInterceptor(srv interface{}, ss grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	// 包装 grpc.ServerStream 以替换 RecvMsg SendMsg这两个方法。
	err := handler(srv, newWrappedStream(ss))
	if err != nil {
		fmt.Printf("RPC failed with error %v\n", err)
	}
	return err
}
