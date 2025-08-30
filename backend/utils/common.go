package utils

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// FileExists 检查文件是否存在
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// GetBaseName 获取文件基础名（不包含扩展名）
func GetBaseName(filename string) string {
	return strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
}

// CreateZipFromDir 将目录打包成ZIP文件
func CreateZipFromDir(sourceDir, zipPath string) error {
	// 创建ZIP文件
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	// 创建ZIP写入器
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// 遍历目录中的所有文件
	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 获取相对路径
		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}

		// 跳过根目录
		if relPath == "." {
			return nil
		}

		// 创建ZIP文件头
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = relPath

		// 如果是目录，添加斜杠
		if info.IsDir() {
			header.Name += "/"
		} else {
			// 设置压缩方法
			header.Method = zip.Deflate
		}

		// 创建ZIP文件写入器
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		// 如果是文件，写入内容
		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(writer, file)
			if err != nil {
				return err
			}
		}

		return nil
	})
}
