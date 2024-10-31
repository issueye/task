package initialize

import (
	"path/filepath"
	"task/internal/global"
	"task/internal/logger"
)

func InitLog() {
	path := filepath.Join(global.RuntimePath, "logs", "task_log.log")
	l, close, err := logger.NewZap(path, logger.LOM_RELEASE)
	if err != nil {
		panic(err)
	}

	global.Logger = l
	global.LogClose = close
}
