package upload

import (
	"blog-service/global"
	"blog-service/pkg/util"
	"io"
	"mime/multipart"
	"os"
	"path"
	"strings"
)

type FileType int

const TypeImage FileType = iota + 1

// 获取文件名称，先是通过获取文件后缀并筛出原始文件名进行 MD5 加密，最后返回经过加密处理后的文件名。
func GetFileName(name string) string {
	ext := GetFileExt(name)
	fileName := strings.TrimSuffix(name, ext)
	fileName = util.EncodeMD5(fileName)
	return fileName + ext
}

// 获取文件名后缀，主要通过path.Ext方法进行循环查找“.”符号，最后通过切片索引返回文件后缀名称
func GetFileExt(name string) string {
	return path.Ext(name)
}

// 获取文件存储路径
func GetSavePath() string {
	return global.AppSetting.Upload.UploadSavePath
}

// 检查文件相关方法

// 检查保存目录是否存在
func CheckSavePath(dst string) bool {
	_, err := os.Stat(dst)
	// Stat 返回描述文件的 FileInfo 结构。
	return os.IsNotExist(err)
	// IsNotExist 返回一个布尔值，指示其参数是否已知
}

// 检查文件权限是否足够
func CheckPermission(dst string) bool {
	_, err := os.Stat(dst)
	return os.IsPermission(err)
	// 指示其参数是否已知报告许可被拒绝
}

// 检查文件后缀是否包含在约定的后缀配置项中
func CheckContainExt(t FileType, name string) bool {
	ext := GetFileExt(name) // 后缀名称
	ext = strings.ToUpper(ext)
	switch t {
	case TypeImage:
		for _, allowExt := range global.AppSetting.Upload.UploadImageAllowExts {
			if strings.ToUpper(allowExt) == ext {
				return true
			}
		}
	}
	return false
}

// 检查文件大小是否超过最大大小限制
func CheckMaxSize(t FileType, f multipart.File) bool {
	// content, _ := ioutil.ReadAll(f)
	content, _ := io.ReadAll(f)
	size := len(content)
	switch t {
	case TypeImage:
		// config `UploadImageMaxSize` is specified in MB in `configs/config.yaml`.
		// Convert to bytes when comparing with the actual file size (bytes).
		maxSizeBytes := global.AppSetting.Upload.UploadImageMaxSize * 1024 * 1024
		if size > maxSizeBytes {
			return true
		}
	}

	return false
}

// 文件写入/创建的相关操作

// 创建在上传文件时所使用的保存目录
func CreateSavePath(dst string, perm os.FileMode) error {
	err := os.MkdirAll(dst, perm)
	if err != nil {
		return err
	}
	return nil
}

// 保存所上传的文件
func SaveFile(file *multipart.FileHeader, dst string) error {
	// 通过 file.Open 方法打开源地址的文件
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// os.Create 方法创建目标地址的文件
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	// os.Create 方法创建目标地址的文件
	_, err = io.Copy(out, src)
	return err
}
