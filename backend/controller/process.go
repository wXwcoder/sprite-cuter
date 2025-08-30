package controller

import (
	"SpriteCuter/core"
	"SpriteCuter/utils"
	"image"
	"image/png"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// ProcessImage 处理图片切割请求
func ProcessImage(c *gin.Context) {
	var req struct {
		Filename string `json:"filename" binding:"required"`
	}

	// 绑定请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorLogger.Printf("请求参数错误: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	// 检查文件是否存在
	uploadPath := filepath.Join("./uploads/", req.Filename)
	if !utils.FileExists(uploadPath) {
		utils.ErrorLogger.Printf("文件不存在: %s", req.Filename)
		c.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
		return
	}

	// 读取PNG文件
	img, err := readPNG(uploadPath)
	if err != nil {
		utils.ErrorLogger.Printf("读取图片时发生错误: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "读取图片时发生错误: " + err.Error()})
		return
	}

	// 获取输出目录名
	outDir := utils.GetBaseName(req.Filename)
	exportPath := filepath.Join("./export/", outDir)

	// 创建输出目录
	if err := utils.CreateDir(exportPath); err != nil {
		utils.ErrorLogger.Printf("创建输出目录失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建输出目录失败: " + err.Error()})
		return
	}

	// 调用核心逻辑进行图片切割
	spritesArray := core.GetSprites(img)

	// 生成CSS文件
	css := core.GetCSS(spritesArray, req.Filename)
	//utils.ErrorLogger.Printf("CSS内容: %v", css)
	cssPath := filepath.Join(exportPath, outDir+".css")
	if err := os.WriteFile(cssPath, []byte(css), 0644); err != nil {
		utils.ErrorLogger.Printf("保存CSS文件失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存CSS文件失败: " + err.Error()})
		return
	}

	// 生成JSON文件
	json := core.GetJson(spritesArray, req.Filename)
	//utils.ErrorLogger.Printf("JSON内容: %v", json)
	jsonPath := filepath.Join(exportPath, outDir+".json")
	if err := os.WriteFile(jsonPath, []byte(json), 0644); err != nil {
		utils.ErrorLogger.Printf("保存JSON文件失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存JSON文件失败: " + err.Error()})
		return
	}

	// 切割并保存精灵图
	for i, rect := range spritesArray {
		if err := core.SaveSprite(img, rect, outDir, i); err != nil {
			// 记录错误但继续处理其他精灵
			utils.ErrorLogger.Printf("保存精灵 %d 时出错: %v", i, err)
		}
	}

	// 打包成ZIP文件
	zipFilename := outDir + ".zip"
	zipPath := filepath.Join("./export/", zipFilename)
	if err := utils.CreateZipFromDir(exportPath, zipPath); err != nil {
		utils.ErrorLogger.Printf("打包ZIP文件失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "打包ZIP文件失败: " + err.Error()})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"message":      "图片切割成功",
		"download_url": "/api/v1/download/" + zipFilename,
	})
}

// readPNG 读取PNG文件
func readPNG(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return png.Decode(file)
}
