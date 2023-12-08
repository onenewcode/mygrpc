package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	pb "mygrpc/proto/hello"
	"net"
	"time"
)

var (
	port = flag.Int("port", 50051, "The service port")
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SearchOrders(req *pb.HelloRequest, stream pb.Greeter_SearchOrdersServer) error {
	log.Printf("Recved %v", req.GetName())
	// 具体返回多少个response根据业务逻辑调整
	for i := 0; i < 10; i++ {
		// 通过 send 方法不断推送数据
		err := stream.Send(&pb.HelloReply{})
		if err != nil {
			log.Fatalf("Send error:%v", err)
			return err
		}
	}
	return nil
}

type wrappedStream struct {
	// 包装器流
	grpc.ServerStream
}

// 接受信息拦截器
func (w *wrappedStream) RecvMsg(m interface{}) error {
	log.Printf("====== [Server Stream Interceptor Wrapper] Receive a message (Type: %T) at %s", m, time.Now().Format(time.RFC3339))
	return w.ServerStream.RecvMsg(m)
}

// 发送消息拦截器
func (w *wrappedStream) SendMsg(m interface{}) error {
	log.Printf("====== [Server Stream Interceptor Wrapper] Send a message (Type: %T) at %v", m, time.Now().Format(time.RFC3339))
	return w.ServerStream.SendMsg(m)
}

func newWrappedStream(s grpc.ServerStream) grpc.ServerStream {
	return &wrappedStream{s}
}

func orderServerStreamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	// 前置处理
	log.Println("====== [Server Stream Interceptor] ", info.FullMethod)

	// 包装器流调用 流RPC
	err := handler(srv, newWrappedStream(ss))
	if err != nil {
		log.Printf("RPC failed with error %v", err)
	}
	return err
}
func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// 开启rpc
	s := grpc.NewServer(grpc.StreamInterceptor(orderServerStreamInterceptor))
	// 注册服务
	pb.RegisterGreeterServer(s, &server{})
	log.Printf("service listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
