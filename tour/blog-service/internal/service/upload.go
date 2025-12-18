package service

import (
	"blog-service/global"
	"blog-service/pkg/upload"
	"errors"
	"mime/multipart"
	"os"
)

type FileInfo struct {
	Name      string
	AccessUrl string
}

// 上传文件
func (svc *Service) UploadFile(fileType upload.FileType, file multipart.File, fileHeader *multipart.FileHeader) (*FileInfo, error) {
	fileName := upload.GetFileName(fileHeader.Filename)
	// 检查后缀名是否允许上传
	if !upload.CheckContainExt(fileType, fileName) {
		return nil, errors.New("file suffix is not supported")
	}
	// 检查文件大小是否超出限制
	if upload.CheckMaxSize(fileType, file) {
		return nil, errors.New("exceeds the maximum limit of file size")
	}

	// 获取配置的保存目录
	uploadSavePath := upload.GetSavePath()
	// 检查保存目录是否存在
	if upload.CheckSavePath(uploadSavePath) {
		// 创建保存目录
		if err := upload.CreateSavePath(uploadSavePath, os.ModePerm); err != nil {
			return nil, errors.New("failed to create save directory")
		}
	}

	// 检查文件权限是否足够
	if upload.CheckPermission(uploadSavePath) {
		return nil, errors.New("insufficient file permissions")
	}

	dst := uploadSavePath + "/" + fileName
	// 保存文件
	if err := upload.SaveFile(fileHeader, dst); err != nil {
		return nil, err
	}

	// 构造文件信息
	accessUrl := global.AppSetting.Upload.UploadServerUrl + "/" + fileName
	return &FileInfo{Name: fileName, AccessUrl: accessUrl}, nil
}
