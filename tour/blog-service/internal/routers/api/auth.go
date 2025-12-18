package api

import (
	"blog-service/global"
	"blog-service/internal/service"
	"blog-service/pkg/app"
	"blog-service/pkg/errcode"

	"github.com/gin-gonic/gin"
)

// @Summary 获取登录认证
// @Tags Auth
// @Description 根据 `app_key` 和 `app_secret` 生成访问 token
// @Accept json
// @Produce json
// @Param body body service.AuthRequest true "请求参数：{app_key, app_secret}"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} errcode.Error
// @Failure 401 {object} errcode.Error
// @Router /auth [post]
func GetAuth(c *gin.Context) {
	param := service.AuthRequest{}
	responses := app.NewResponse(c)
	vaild, errs := app.BindAndValid(c, &param)
	if !vaild {
		global.Logger.Errorf("app.BindAndVaild errs:%v", errs)
		responses.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Error()))
		return
	}
	svc := service.New(c.Request.Context())
	err := svc.CheckAuth(&param)
	if err != nil {
		global.Logger.Errorf("svc.CheckAuth err:%v", err)
		responses.ToErrorResponse(errcode.UnauthorizedAuthNotExist)
		return
	}

	token, err := app.GenerateToken(param.AppKey, param.AppSecret)
	if err != nil {
		global.Logger.Errorf("app.GenerateToken err:%v", err)
		responses.ToErrorResponse(errcode.UnauthorizedTokenGenerate)
		return
	}
	responses.ToResponse(gin.H{
		"token": token,
	})
}
