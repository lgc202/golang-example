package main

import (
	"context"
	"fmt"
	"golang-example/grpc/demo05/proto"
	"io"
	"log"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	address = "localhost:50051"
)

func main() {
	// 1.建立连接 获取client
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStreamInterceptor(streamInterceptor))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	client := proto.NewEchoClient(conn)

	// 2.执行各个Stream的对应方法
	// unary(client)
	// serverStream(client)
	// clientStream(client)
	bidirectionalStream(client)
}

func unary(client proto.EchoClient) {
	resp, err := client.UnaryEcho(context.Background(), &proto.EchoRequest{Message: "hello world"})
	if err != nil {
		log.Printf("send error:%v\n", err)
		return
	}

	fmt.Printf("Recved:%v \n", resp.GetMessage())
}

// 1. 建立连接 获取client
// 2. 通过 client 获取stream
// 3. for循环中通过stream.Recv()依次获取服务端推送的消息
// 4. err==io.EOF则表示服务端关闭stream了
func serverStream(client proto.EchoClient) {
	// 2.调用获取stream
	stream, err := client.ServerStreamingEcho(context.Background(), &proto.EchoRequest{Message: "Hello World"})
	if err != nil {
		log.Fatalf("could not echo: %v", err)
	}

	// 3. for循环获取服务端推送的消息
	for {
		// 通过 Recv() 不断获取服务端send()推送的消息
		resp, err := stream.Recv()
		// 4. err==io.EOF则表示服务端关闭stream了 退出
		if err == io.EOF {
			log.Println("server closed")
			break
		}

		if err != nil {
			log.Printf("Recv error:%v", err)
			continue
		}

		log.Printf("Recv data:%v", resp.GetMessage())
	}
}

// clientStream 客户端流
// 1. 建立连接并获取client
// 2. 获取 stream 并通过 Send 方法不断推送数据到服务端
// 3. 发送完成后通过stream.CloseAndRecv() 关闭steam并接收服务端返回结果
func clientStream(client proto.EchoClient) {
	// 2.获取 stream 并通过 Send 方法不断推送数据到服务端
	stream, err := client.ClientStreamingEcho(context.Background())
	if err != nil {
		log.Fatalf("ClientStreamingEcho() error: %v", err)
	}
	for i := int64(0); i < 2; i++ {
		err := stream.Send(&proto.EchoRequest{Message: "hello world"})
		if err != nil {
			log.Printf("send error: %v", err)
			continue
		}
	}

	// 3. 发送完成后通过stream.CloseAndRecv() 关闭steam并接收服务端返回结果
	// (服务端则根据err==io.EOF来判断client是否关闭stream)
	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("CloseAndRecv() error: %v", err)
	}
	log.Printf("message: %v", resp.GetMessage())
}

// bidirectionalStream 双向流
//  1. 建立连接 获取client
//  2. 通过client获取stream
//  3. 开两个goroutine 分别用于Recv()和Send()
//     3.1 一直Recv()到err==io.EOF(即服务端关闭stream)
//     3.2 Send()则由自己控制
//  4. 发送完毕调用 stream.CloseSend()关闭stream 必须调用关闭 否则Server会一直尝试接收数据 一直报错...
func bidirectionalStream(client proto.EchoClient) {
	var wg sync.WaitGroup
	// 2. 调用方法获取stream
	stream, err := client.BidirectionalStreamingEcho(context.Background())
	if err != nil {
		panic(err)
	}
	// 3.开两个goroutine 分别用于Recv()和Send()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			req, err := stream.Recv()
			if err == io.EOF {
				fmt.Println("Server Closed")
				break
			}
			if err != nil {
				continue
			}
			fmt.Printf("Recv Data:%v \n", req.GetMessage())
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		for i := 0; i < 2; i++ {
			err := stream.Send(&proto.EchoRequest{Message: "hello world"})
			if err != nil {
				log.Printf("send error:%v\n", err)
			}
			time.Sleep(time.Second)
		}
		// 4. 发送完毕关闭stream
		err := stream.CloseSend()
		if err != nil {
			log.Printf("Send error:%v\n", err)
			return
		}
	}()
	wg.Wait()
}

// unaryInterceptor 一个简单的 unary interceptor 示例。
func unaryInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	// pre-processing
	start := time.Now()
	err := invoker(ctx, method, req, reply, cc, opts...) // invoking RPC method
	// post-processing
	end := time.Now()
	fmt.Printf("RPC: %s, req:%v start time: %s, end time: %s, err: %v\n", method, req, start.Format(time.RFC3339), end.Format(time.RFC3339), err)
	return err
}

// wrappedStream  用于包装 grpc.ClientStream 结构体并拦截其对应的方法。
type wrappedStream struct {
	grpc.ClientStream
}

func newWrappedStream(s grpc.ClientStream) grpc.ClientStream {
	return &wrappedStream{s}
}

func (w *wrappedStream) RecvMsg(m interface{}) error {
	fmt.Printf("Receive a message (Type: %T) at %v\n", m, time.Now().Format(time.RFC3339))
	return w.ClientStream.RecvMsg(m)
}

func (w *wrappedStream) SendMsg(m interface{}) error {
	fmt.Printf("Send a message (Type: %T) at %v\n", m, time.Now().Format(time.RFC3339))
	return w.ClientStream.SendMsg(m)
}

// streamInterceptor 一个简单的 stream interceptor 示例。
func streamInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	s, err := streamer(ctx, desc, cc, method, opts...)
	if err != nil {
		return nil, err
	}
    // 返回的是自定义的封装过的 stream
	return newWrappedStream(s), nil
}
