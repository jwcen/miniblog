package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	// "google.golang.org/protobuf/types/known/emptypb"

	apiv1 "github.com/jwcen/miniblog/pkg/api/apiserver/v1"
)

var (
	// 定义命令行参数
	addr  = flag.String("addr", "localhost:6666", "The grpc server address to connect to.") // gRPC 服务的地址
	limit = flag.Int64("limit", 10, "Limit to list users.")                                 // 限制列出用户的数量
)

func main() {
	flag.Parse()

	log.Printf("正在连接到 gRPC 服务器: %s", *addr)

	// 使用 grpc.Dial 建立客户端与 gRPC 服务端的连接
	conn, err := grpc.Dial(*addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(), // 添加阻塞选项，确保连接建立
	)
	if err != nil {
		log.Fatalf("连接 gRPC 服务器失败: %v", err)
	}
	defer conn.Close()

	log.Printf("成功连接到 gRPC 服务器")

	// 使用连接创建一个 MiniBlog 的 gRPC 客户端实例
	client := apiv1.NewMiniBlogClient(conn)

	// 创建一个带超时的上下文，增加超时时间到 10 秒
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Printf("正在调用 Healthz 方法...")

	// 调用 MiniBlog 的 Healthz 方法，检查服务健康状况
	resp, err := client.Healthz(ctx, nil)
	if err != nil {
		log.Fatalf("调用 Healthz 失败: %v", err)
	}

	log.Printf("成功获取 Healthz 响应")

	jsonData, _ := json.MarshalIndent(resp, "", "  ") // 使用 json.MarshalIndent 美化输出
	fmt.Println(string(jsonData))
}
