package initialize

import (
	"task/internal/global"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
)

func InitEvetnBus() {
	global.Logger.Sugar().Debug("初始化事件总线")
	pubSub := gochannel.NewGoChannel(
		gochannel.Config{},
		watermill.NewStdLogger(false, false),
	)

	global.PubSub = pubSub
}
