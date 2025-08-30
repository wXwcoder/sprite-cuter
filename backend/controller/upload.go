package controller

import (
	"SpriteCuter/utils"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// UploadImage 处理图片上传请求
func UploadImage(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "上传文件失败: " + err.Error()})
		return
	}

	// 生成唯一文件名
	uniqueFilename := utils.GenerateUniqueFilename(file.Filename)
	filePath := filepath.Join("./uploads/", uniqueFilename)

	// 保存文件
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		utils.ErrorLogger.Printf("文件保存失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文件保存失败: " + err.Error()})
		return
	}

	// 返回成功响应
	utils.InfoLogger.Printf("文件上传成功: %s", uniqueFilename)
	c.JSON(http.StatusOK, gin.H{
		"message":  "文件上传成功",
		"filename": uniqueFilename,
	})
}
