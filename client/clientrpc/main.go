package main

import (
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "mygrpc/proto/hello" // 引入编译生成的包
)

const (
	defaultName = "world"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
)

func main() {
	flag.Parse()
	// 与服务建立连接.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)
	// 创建指定服务的客户端
	c := pb.NewGreeterClient(conn)

	// 连接服务器并打印出其响应。
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// 调用指定方法
	updateStream, err := c.UpdateOrders(ctx)
	if err != nil {
		log.Fatalf("%v.UpdateOrders(_) = _, %v", c, err)
	}
	// Updating order 1
	if err := updateStream.Send(&pb.HelloRequest{Name: "1"}); err != nil {
		log.Fatalf("%v.Send(%v) = %v", updateStream, &pb.HelloRequest{Name: "1"}, err)
	}
	// Updating order 2
	if err := updateStream.Send(&pb.HelloRequest{Name: "2"}); err != nil {
		log.Fatalf("%v.Send(%v) = %v", updateStream, &pb.HelloRequest{Name: "2"}, err)
	}
	// 发送关闭信号并接收服务端响应
	err = updateStream.CloseSend()
	if err != nil {
		log.Fatalf("%v.CloseAndRecv() got error %v, want %v", updateStream, err, nil)
	}
	log.Printf("客户端流传输结束")
}
