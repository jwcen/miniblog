// Copyright 2024 jayvee <jvvcen@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"

	apiv1 "github.com/jwcen/miniblog/pkg/api/apiserver/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

var (
	// 定义命令行参数
	addr  = flag.String("addr", "localhost:6666", "The grpc server address to connect to.") // gRPC 服务的地址
	limit = flag.Int64("limit", 10, "Limit to list users.")                                 // 限制列出用户的数量
)

func main() {
	flag.Parse()

	conn, err := grpc.NewClient(*addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(unaryClientInterceptor()),
	)
	if err != nil {
		log.Fatalf("Failed to connect to grpc server: %v", err)
	}
	defer conn.Close()

	client := apiv1.NewMiniBlogClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// 创建一个 Metadata 用于传递请求元数据
	md := metadata.Pairs("custom-header", "value123")
	ctx = metadata.NewOutgoingContext(ctx, md)

	// 调用 MiniBlog 的 Healthz 方法，检查服务健康状况
	var header metadata.MD                                      // 用于存储返回的 Header 元数据
	resp, err := client.Healthz(ctx, nil, grpc.Header(&header)) // 发起 gRPC 请求
	if err != nil {
		log.Fatalf("Failed to call healthz: %v", err) // 如果调用失败，记录错误并退出程序
	}

	for key, val := range header {
		fmt.Printf("Response Header (key: %s, value: %s)\n", key, val)
	}

	// 将返回的响应数据转换为 JSON 格式
	jsonData, _ := json.Marshal(resp) // 使用 json.Marshal 将响应对象序列化为 JSON 格式
	fmt.Println(string(jsonData))     // 输出 JSON 数据到终端
}

func unaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		// 打印请求方法和请求参数
		log.Printf("[UnaryClientInterceptor] Invoking method: %s", method)

		// 添加自定义元数据
		md := metadata.Pairs("interceptor-header", "interceptor-value")
		ctx = metadata.NewOutgoingContext(ctx, md)

		// 调用实际的 gRPC 方法
		err := invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {
			log.Printf("[UnaryClientInterceptor] Method: %s, Error: %v", method, err)
			return err
		}

		return nil
	}
}
