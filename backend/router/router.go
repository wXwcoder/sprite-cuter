package router

import (
	"github.com/gin-gonic/gin"
	"SpriteCuter/controller"
)

func SetupRouter() *gin.Engine {
	// 设置为发布模式
	gin.SetMode(gin.ReleaseMode)
	
	// 创建gin引擎
	r := gin.Default()

	// 配置CORS
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")
		c.Header("Access-Control-Allow-Credentials", "true")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	// 设置路由
	v1 := r.Group("/api/v1")
	{
		// 图片上传接口
		v1.POST("/upload", controller.UploadImage)
		
		// 图片切割接口
		v1.POST("/process", controller.ProcessImage)
		
		// 文件下载接口
		v1.GET("/download/:filename", controller.DownloadFile)
	}

	return r
}