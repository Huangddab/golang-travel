package api

import (
	"blog-service/global"
	"blog-service/internal/service"
	"blog-service/pkg/app"
	"blog-service/pkg/convert"
	"blog-service/pkg/errcode"
	"blog-service/pkg/upload"

	"github.com/gin-gonic/gin"
)

type Upload struct{}

func NewUpload() Upload { return Upload{} }

// @Summary 文件上传
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "上传的文件"
// @Param type formData int true "文件类型 (例如: 1=jpg,2=jpeg,3=png,4=gif)"
// @Success 200 {object} service.FileInfo "成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /upload/file [post]
func (u Upload) UploadFile(c *gin.Context) {
	reponse := app.NewResponse(c)
	// c.Request.FormFile 读取入参 file 字段的上传文件信息
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		reponse.ToErrorResponse(errcode.InvalidParams.WithDetails(err.Error()))
		return
	}

	fileType := convert.StrTo(c.PostForm("type")).MustInt()
	if fileHeader == nil || fileType <= 0 {
		reponse.ToErrorResponse(errcode.InvalidParams)
		return
	}

	svc := service.New(c.Request.Context())
	fileInfo, err := svc.UploadFile(upload.FileType(fileType), file, fileHeader)
	if err != nil {
		global.Logger.Errorf("svc.UploadFile err: %v", err)
		reponse.ToErrorResponse(errcode.ErrorUploadFileFail.WithDetails(err.Error()))
		return
	}
	reponse.ToResponse(gin.H{
		"file_access_url": fileInfo.AccessUrl,
	})
}
