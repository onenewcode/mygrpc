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

func orderUnaryClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	// 前置处理逻辑
	log.Println("Method : " + method)

	// 调用invoker 执行远程方法
	err := invoker(ctx, method, req, reply, cc, opts...)

	// 后置处理逻辑
	log.Println(reply)
	return err
}

func main() {
	flag.Parse()
	// 与服务建立连接.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(orderUnaryClientInterceptor)) //添加拦截器
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	// 创建指定服务的客户端
	c := pb.NewGreeterClient(conn)

	// 连接服务器并打印出其响应。
	ctx, cancel := context.WithTimeout(context.Background(), time.Second) // 设置超时时间为一秒
	defer cancel()
	// 调用指定方法
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: *name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())
}
