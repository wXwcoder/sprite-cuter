package utils

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

// CleanupTempFiles 定期清理临时文件
func CleanupTempFiles(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			cleanDir("./uploads/")
			cleanDir("./export/")
		}
	}
}

// cleanDir 清理指定目录下的文件
func cleanDir(dirPath string) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		log.Printf("读取目录 %s 失败: %v", dirPath, err)
		return
	}

	for _, file := range files {
		filePath := filepath.Join(dirPath, file.Name())
		info, err := file.Info()
		if err != nil {
			log.Printf("获取文件信息 %s 失败: %v", filePath, err)
			continue
		}

		// 删除超过1小时的文件
		if time.Since(info.ModTime()) > time.Hour {
			if err := os.RemoveAll(filePath); err != nil {
				log.Printf("删除文件 %s 失败: %v", filePath, err)
			} else {
				log.Printf("已删除文件: %s", filePath)
			}
		}
	}
}