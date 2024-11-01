package initialize

import (
	_ "embed"
	"task/internal/bdb"
	"task/internal/ipc"
	"task/internal/task"
)

//go:embed banner.txt
var banner string

func Init() {
	InitRuntime()
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()
	// 初始化日志
	InitLog()
	// 初始化事件
	InitEventBus()
	// 初始化 task
	task.InitTask()
	// 初始化 ipc客户端
	ipc.InitIpc()
	// 初始化数据库
	bdb.InitBdb()
	// 初始化脚本引擎
	InitCodeEngine()
}
