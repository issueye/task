package initialize

import (
	"task/internal/global"
	"task/internal/logger"
)

func InitLog() {
	l, close, err := logger.NewZap("task_log.log", logger.LOM_RELEASE)
	if err != nil {
		panic(err)
	}

	global.Logger = l
	global.LogClose = close
}
