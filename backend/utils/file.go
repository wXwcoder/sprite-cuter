package utils

import (
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// GenerateUniqueFilename 生成唯一文件名
func GenerateUniqueFilename(originalName string) string {
	// 获取文件扩展名
	ext := filepath.Ext(originalName)
	
	// 生成随机字符串
	randomStr := generateRandomString(8)
	
	// 获取当前时间戳
	timestamp := time.Now().Unix()
	
	// 组合生成唯一文件名
	filename := strings.TrimSuffix(originalName, ext)
	return fmt.Sprintf("%s_%s_%d%s", filename, randomStr, timestamp, ext)
}

// generateRandomString 生成指定长度的随机字符串
func generateRandomString(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		// 如果随机数生成失败，使用时间戳作为后备
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[bytes[i]%byte(len(charset))]
	}
	return string(b)
}

// DirExists 检查目录是否存在
func DirExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// CreateDir 创建目录
func CreateDir(path string) error {
	return os.MkdirAll(path, 0755)
}