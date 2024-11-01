package initialize

import (
	"path/filepath"
	"task/internal/code_engine"
	"task/internal/global"
)

func InitCodeEngine() {
	logPath := filepath.Join(global.RuntimePath, "logs")
	global.CodeEngine = code_engine.NewCore(
		code_engine.OptionLog(logPath, global.Logger.Named("code_engine")),
	)
	scriptPath := filepath.Join(global.RuntimePath, "scripts")
	global.CodeEngine.SetGlobalPath(scriptPath)
}
