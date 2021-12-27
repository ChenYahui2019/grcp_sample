package main

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"golang.org/x/net/context"
	// 导入grpc包
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	_ "google.golang.org/grpc/health"
	"google.golang.org/grpc/status"
)

const (
	Address             string = ":8000"
	HEALTHCHECK_SERVICE        = "grpc.health.v1.Health"
)

var client StreamClient

// conversations 调用服务端的Conversations方法
func conversations() {
	clientDeadline := time.Now().Add(time.Duration(3 * time.Second))
	ctx, cancel := context.WithDeadline(context.Background(), clientDeadline)
	defer cancel()

	//调用服务端的Conversations方法，获取流
	stream, err := client.Conversations(ctx)
	//stream, err := client.Conversations(context.Background())
	if err != nil {
		fmt.Printf("get conversations stream err: %v\n", err)
	}
	for n := 0; n < 5; n++ {
		err := stream.Send(&StreamRequest{Question: "stream client rpc " + strconv.Itoa(n)})
		if err != nil {
			fmt.Printf("stream request err: %v\n", err)
		}
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		stat, ok := status.FromError(err)
		if ok {
			//判断是否为调用超时
			if stat.Code() == codes.DeadlineExceeded {
				fmt.Println("timeout!")
				err = stream.CloseSend()
				if err != nil {
					fmt.Printf("Conversations close stream err: %v\n", err)
				}
				return
			}
		}
		if err != nil {
			fmt.Printf("Conversations get stream err: %v\n", err)
		}
		// 打印返回值
		fmt.Println(res.Answer)
	}
	//最后关闭流
	err = stream.CloseSend()
	if err != nil {
		fmt.Printf("Conversations close stream err: %v\n", err)
	}
}

func main() {
	// 连接服务器
	conn, err := grpc.Dial(Address, grpc.WithInsecure(), grpc.WithDefaultServiceConfig(
		fmt.Sprintf(`{"HealthCheckConfig": {"ServiceName": "%s"}}`, HEALTHCHECK_SERVICE)))
	//conn, err := grpc.Dial(Address, grpc.WithInsecure())
	if err != nil {
		fmt.Printf("net.Connect err: %v\n", err)
	}
	defer conn.Close()

	// 建立gRPC连接
	client = NewStreamClient(conn)
	conversations()
}
