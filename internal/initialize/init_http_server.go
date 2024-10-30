package initialize

import (
	"context"
	"fmt"
	"net/http"
	"task/internal/global"
	"task/internal/router"
	"task/internal/utils"

	"github.com/gin-gonic/gin"
)

func InitHttpServer(_ context.Context) {
	// gin引擎对象
	global.Router = gin.New()
	router.InitRouter(global.Router)

	global.Logger.Sugar().Debug("http 服务获取端口中...")

	port, err := utils.GetFreePort()
	if err != nil {
		panic(err)
	}

	global.Logger.Sugar().Debugf("http 服务获取端口成功: %d", port)
	httpSrv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: global.Router,
	}

	fmt.Println("Http Server Start At Port:", port)
	httpSrv.ListenAndServe()
}
