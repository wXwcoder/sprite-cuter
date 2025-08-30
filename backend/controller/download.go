package controller

import (
	"SpriteCuter/utils"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// DownloadFile 提供文件下载
func DownloadFile(c *gin.Context) {
	filename := c.Param("filename")
	filePath := filepath.Join("./export/", filename)

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		utils.ErrorLogger.Printf("文件不存在: %s", filename)
		c.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
		return
	}

	// 提供文件下载
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/octet-stream")
	c.File(filePath)
}
