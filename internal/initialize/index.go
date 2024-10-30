package initialize

import (
	"context"
	_ "embed"
	"task/internal/ipc"
)

//go:embed banner.txt
var banner string

func Init() {
	// 启用mqtt
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// 初始化日志
	InitLog()
	// 初始化事件
	InitEvetnBus()
	// 初始化 ipc客户端
	ipc.InitIpc()
	// 初始化HTTP 服务
	go InitHttpServer(ctx)
}
