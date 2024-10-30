package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	{
		v1.GET("/ping", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{"code": 200, "msg": "pong"})
		})
	}

	r.NoRoute(func(ctx *gin.Context) {
		fmt.Printf("404- %s\n", ctx.Request.URL.Path)
		ctx.JSON(404, gin.H{"code": 404, "msg": "------ 404"})
	})
}
