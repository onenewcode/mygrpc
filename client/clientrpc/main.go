package main

import (
	"context"
	"flag"
	"io"
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
	defer conn.Close()
	// 创建指定服务的客户端
	c := pb.NewGreeterClient(conn)

	// 连接服务器并打印出其响应。
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// 调用指定方法
	searchStream, _ := c.SearchOrders(ctx, &pb.HelloRequest{Name: "开始服务端rpc流测试"})
	for {
		// 客户端 Recv 方法接收服务端发送的流
		searchOrder, err := searchStream.Recv()
		if err == io.EOF {
			log.Print("EOF")
			break
		}
		if err == nil {
			log.Print("Search Result : ", searchOrder)
		}
	}
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: *name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())
}
