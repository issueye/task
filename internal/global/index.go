package global

import (
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/gin-gonic/gin"
	"github.com/issueye/ipc_grpc/client"
	"go.uber.org/zap"
)

var (
	AppName     string       // 应用程序名称
	Version     string = "1" // 版本号
	GitHash     string       // git commit hash
	GitBranch   string       // git branch
	BuildTime   string       // 构建时间
	CookieKey   = "PLUGIN_TASK"
	CookieValue = "vz9mr6vevv50zyd5pgnkw6vtkhvhzzm0"

	Router    *gin.Engine
	IpcClient *client.Client
	PubSub    *gochannel.GoChannel
	Logger    *zap.Logger
	LogClose  func()
)

const (
	TOPIC_SHOW_HOME = "TOPIC_SHOW_HOME"
)
