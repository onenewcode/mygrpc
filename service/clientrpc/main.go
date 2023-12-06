package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	pb "mygrpc/proto/hello"
	"net"
)

var (
	port = flag.Int("port", 50051, "The service port")
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
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
func (s *server) UpdateOrders(stream pb.Greeter_UpdateOrdersServer) error {

	for {
		log.Println("开始接受客户端的流")
		// Recv 对客户端发来的请求接收
		order, err := stream.Recv()
		if err == io.EOF {
			// 流结束，关闭并发送响应给客户端
			return stream.Send(&pb.HelloReply{Message: "接受客户流结束"})
		}
		if err != nil {
			return err
		}
		// 更新数据
		log.Printf("Order ID : %s - %s", order.GetName(), "Updated")
	}
}
func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// 开启rpc
	s := grpc.NewServer()
	// 注册服务
	pb.RegisterGreeterServer(s, &server{})
	log.Printf("service listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
