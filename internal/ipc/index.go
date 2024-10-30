package ipc

import (
	"context"
	"fmt"
	"io"
	"runtime"
	"task/internal/global"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/issueye/ipc_grpc/client"
	"github.com/issueye/ipc_grpc/grpc/pb"
	"github.com/issueye/ipc_grpc/vars"
)

func InitIpc() {
	global.Logger.Sugar().Debug("初始化 ipc")
	c, err := client.NewClient()
	if err != nil {
		global.Logger.Sugar().Errorf("初始化 ipc 失败 %s", err.Error())
		panic(err)
	}

	global.IpcClient = c

	vars.AppName = "task"
	vars.CookieKey = "tool_plugin_task"
	vars.CookieValue = "JRJYKeDCR6Q20RGitb4Gk2y9b4mtP69h"
	vars.GoVersion = runtime.Version()
	global.IpcClient.Register()
	go global.IpcClient.Heartbeat()

	// 监听事件
	global.Logger.Sugar().Debug("监听宿主程序事件")
	event, err := global.IpcClient.HostHelper().Event(context.Background(), &pb.ClientRequest{
		CookieKey: vars.CookieKey,
	})

	if err != nil {
		panic(err)
	}

	go func() {
		errCount := 0
		for {

			if errCount > 3 {
				global.Logger.Sugar().Errorf("Event 监听超过3次失败，退出监听")
				break
			}

			msg, err := event.Recv()
			if err == io.EOF {
				global.Logger.Sugar().Infof("Event 退出监听 %s", err.Error())
				break
			}

			if err != nil {
				global.Logger.Sugar().Errorf("Event 监听错误 %s", err.Error())
				errCount++
				continue
			}

			global.Logger.Sugar().Debugf("收到事件 %s", msg.Type.String())
			fmt.Println(msg.Type.String(), msg.Server.Name)
			errCount = 0
		}
	}()

	global.Logger.Sugar().Debug("监听宿主程序指令")
	command, err := global.IpcClient.HostHelper().Command(context.Background(), &pb.ClientRequest{
		CookieKey: vars.CookieKey,
	})

	if err != nil {
		panic(err)
	}

	go func() {
		errCount := 0
		for {

			if errCount > 3 {
				global.Logger.Sugar().Errorf("Command 监听超过3次失败，退出监听")
				break
			}

			msg, err := command.Recv()
			if err == io.EOF {
				global.Logger.Sugar().Infof("Command 退出监听 %s", err.Error())
				break
			}

			if err != nil {
				global.Logger.Sugar().Errorf("Command 监听错误 %s", err.Error())
				errCount++
				continue
			}

			global.Logger.Sugar().Debugf("收到指令 %s", msg.Command)
			if msg.Command == "show" {
				msg := message.NewMessage(watermill.NewUUID(), message.Payload([]byte{}))
				global.PubSub.Publish(global.TOPIC_SHOW_HOME, msg)
			}

			errCount = 0
		}
	}()
}
