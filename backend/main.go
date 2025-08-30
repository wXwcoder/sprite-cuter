package main

import (
	"SpriteCuter/router"
	"SpriteCuter/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// 设置Gin为发布模式
	gin.SetMode(gin.ReleaseMode)

	// 初始化路由
	Router := router.SetupRouter()

	// 启动临时文件清理任务
	go utils.CleanupTempFiles(30 * time.Minute)

	// 记录服务器启动信息
	utils.InfoLogger.Println("服务器启动在端口 8080")

	// 启动服务器
	Router.Run(":8080")
}
