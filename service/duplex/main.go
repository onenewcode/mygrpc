package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	pb "mygrpc/proto/hello"
	"net"
	"sync"
)

var (
	port = flag.Int("port", 50051, "The service port")
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) ProcessOrders(stream pb.Greeter_ProcessOrdersServer) error {
	var (
		waitGroup sync.WaitGroup // 一组 goroutine 的结束
		// 设置通道
		msgCh = make(chan *pb.HelloReply)
	)
	// 计数器加1
	waitGroup.Add(1)
	// 消费队列中的内容
	go func() {
		// 计数器减一
		defer waitGroup.Done()
		for {
			v := <-msgCh
			fmt.Println(v)
			err := stream.Send(v)
			if err != nil {
				fmt.Println("Send error:", err)
				break
			}
		}
	}()
	waitGroup.Add(1)
	// 向队列中添加内容
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
			fmt.Printf("Recved :%v \n", req.GetName())
			msgCh <- &pb.HelloReply{Message: "服务端传输数据"}
		}
		close(msgCh)
	}()
	// 等待 计数器问0 推出
	waitGroup.Wait()

	// 返回nil表示已经完成响应
	return nil
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
