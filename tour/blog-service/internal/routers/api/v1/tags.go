package v1

import (
	"blog-service/global"
	"blog-service/internal/service"
	"blog-service/pkg/app"
	"blog-service/pkg/convert"
	"blog-service/pkg/errcode"

	"github.com/gin-gonic/gin"
)

type Tag struct {
}

// 构造函数
// 生产和返回一个初始化好的结构体实例
func NewTag() Tag {
	return Tag{}
}

// 这里使用的是t Tag类型的值接收者 意味着每次调用方法时
// 都会得到Tag类型的一个副本 而不是对原始结构体的引用
// 适用于不需要修改结构体内部状态的方法

// 更常见的是使用指针接收者 以避免复制开销 并允许方法修改结构体的状态

// @Summary 获取多个标签
// @Tags Tags
// @Security ApiKeyAuth
// @Produce json
// @Param name query string false "标签名称" maxlength(100)
// @Param state query int false "状态" Enums(0, 1) default(1)
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} model.TagSwagger "成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /api/v1/tags [get]
func (t Tag) List(c *gin.Context) {
	// 定义请求参数结构体
	param := service.TagListRequest{}
	// 创建响应实例
	response := app.NewResponse(c)
	// 参数绑定和校验
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Logger.Errorf("app.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	// 创建服务实例
	svc := service.New(c.Request.Context())
	// 设置分页参数
	pager := app.Pager{Page: app.GetPage(c), PageSize: app.GetPageSize(c)}
	// 统计标签总数
	totalRows, err := svc.CountTag(&service.CountTagRequest{Name: param.Name, State: param.State})
	if err != nil {
		global.Logger.Errorf("svc.CountTag err:%v", err)
		response.ToErrorResponse(errcode.ErrorCountTagFail)
		return
	}
	// 获取标签列表
	tags, err := svc.GetTagList(&param, &pager)
	if err != nil {
		global.Logger.Errorf("svc.GetTagList err:%v", err)
		response.ToErrorResponse(errcode.ErrorGetTagListFail)
		return
	}
	response.ToResponseList(tags, int(totalRows))
}

//
// func (t Tag) List(c *gin.Context) {
// 	// app.NewResponse(c).ToErrorResponse(errcode.ServerError)
// 	// {"code":10000000,"msg":"服务内部错误"}* Connection #0 to host localhost:8088 left intact

// 	param := struct {
// 		Name  string `form:"name" binding:"max=100"`
// 		State string `form:"state,default=1" binding:"oneof=0 1"`
// 	}{}
// 	response := app.NewResponse(c)
// 	valid, errs := app.BindAndValid(c, &param)
// 	if !valid {
// 		global.Logger.Errorf("app.BindAndValid errs: %v", errs)
// 		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
// 		return
// 	}
// 	response.ToResponse(gin.H{})
// 	return
// }

// @Summary 新增标签
// @Tags Tags
// @Security ApiKeyAuth
// @Produce  json
// @Param data body service.CreateTagRequest true "创建标签"
// @Success 200 {object} model.TagSwagger "成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /api/v1/tags [post]
func (t Tag) Create(c *gin.Context) {
	param := service.CreateTagRequest{}
	response := app.NewResponse(c)
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Logger.Errorf("app.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	svc := service.New(c.Request.Context())
	err := svc.CreateTag(&param)
	if err != nil {
		global.Logger.Errorf("svc.CreateTag err:%v", err)
		response.ToErrorResponse(errcode.ErrorCreateTagFail)
		return
	}

	response.ToResponse(gin.H{})
	return
	// c.JSON(http.StatusOK, gin.H{"message": "创建标签"})
}

// @Summary 更新标签
// @Tags Tags
// @Security ApiKeyAuth
// @Produce  json
// @Param id path int true "标签 ID"
// @Param data body service.UpdateTagRequest true "更新标签"
// @Success 200 {array} model.TagSwagger "成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /api/v1/tags/{id} [put]
func (t Tag) Update(c *gin.Context) {
	param := service.UpdateTagRequest{ID: convert.StrTo(c.Param("id")).MustUint32()}
	response := app.NewResponse(c)
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Logger.Errorf("app.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	svc := service.New(c.Request.Context())
	err := svc.UpdateTag(&param)
	if err != nil {
		global.Logger.Errorf("svc.UpdateTag err:%v", err)
		response.ToErrorResponse(errcode.ErrorUpdateTagFail)
		return
	}
	response.ToResponse(gin.H{})
	return

	// c.JSON(http.StatusOK, gin.H{"message": "更新标签"})
}

// @Summary 删除标签
// @Tags Tags
// @Security ApiKeyAuth
// @Produce  json
// @Param id path int true "标签 ID"
// @Success 200 {string} string "成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /api/v1/tags/{id} [delete]
func (t Tag) Delete(c *gin.Context) {
	param := service.DeleteTagRequest{ID: convert.StrTo(c.Param("id")).MustUint32()}
	response := app.NewResponse(c)
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Logger.Errorf("app.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	svc := service.New(c.Request.Context())
	err := svc.DeleteTag(&param)
	if err != nil {
		global.Logger.Errorf("svc.DeleteTag err:%v", err)
		response.ToErrorResponse(errcode.ErrorDeleteTagFail)
		return
	}
	response.ToResponse(gin.H{})
	return

	// c.JSON(http.StatusOK, gin.H{"message": "删除标签"})
}
