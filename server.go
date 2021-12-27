package main

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"time"

	// 导入grpc包
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

const (
	// Address 监听地址
	Address string = ":8000"
	// Network 网络通信协议
	Network             string = "tcp"
	HEALTHCHECK_SERVICE        = "grpc.health.v1.Health"
)

// StreamService 定义我们的服务
type StreamService struct{}

// Conversations 实现Conversations方法
func (s *StreamService) Conversations(srv Stream_ConversationsServer) error {
	n := 1
	for {
		req, err := srv.Recv()
		if err == io.EOF {
			fmt.Println(err)
			return nil
		}
		if err != nil {
			fmt.Println(err)
			return err
		}
		err = srv.Send(&StreamResponse{
			Answer: "from stream server answer: the " + strconv.Itoa(n) + " question is " + req.Question,
		})
		if err != nil {
			fmt.Println(err)
			return err
		}
		n++
		fmt.Printf("from stream client question: %s\n", req.Question)
		time.Sleep(5 * time.Second)
	}
}

func main() {
	// 监听本地端口
	listener, err := net.Listen(Network, Address)
	if err != nil {
		fmt.Printf("net.Listen err: %v\n", err)
	}
	fmt.Println(Address + " net.Listing...")
	// 新建gRPC服务器实例
	grpcServer := grpc.NewServer()

	healthserver := health.NewServer()
	healthserver.SetServingStatus(HEALTHCHECK_SERVICE, healthpb.HealthCheckResponse_SERVING)
	healthpb.RegisterHealthServer(grpcServer, healthserver)

	// 在gRPC服务器注册我们的服务
	RegisterStreamServer(grpcServer, &StreamService{})

	go func() {
		//用服务器 Serve() 方法以及我们的端口信息区实现阻塞等待，直到进程被杀死或者 Stop() 被调用
		err = grpcServer.Serve(listener)
		if err != nil {
			fmt.Printf("grpcServer.Serve err: %v\n", err)
		}
	}()
	time.Sleep(time.Minute)

	grpcServer.GracefulStop()
	fmt.Println("GracefulStop")
	time.Sleep(time.Minute)
}
