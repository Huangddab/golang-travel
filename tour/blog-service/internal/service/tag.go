package service

import (
	"blog-service/internal/model"
	"blog-service/pkg/app"
	"blog-service/pkg/errcode"
	"fmt"
)

// 接口中定义的的增删改查和统计行为进行了 Request 结构体编写
type CountTagRequest struct {
	Name  string `json:"name" form:"name" binding:"max=100"`
	State uint8  `json:"state" form:"state,default=1" binding:"oneof=0 1"`
}

type TagListRequest struct {
	Name  string `json:"name" form:"name" binding:"max=100"`
	State uint8  `json:"state" form:"state,default=1" binding:"oneof=0 1"`
}

type CreateTagRequest struct {
	Name      string `json:"name" form:"name" binding:"required,min=2,max=100"`
	CreatedBy string `json:"created_by" form:"created_by" binding:"required,min=2,max=100"`
	State     uint8  `json:"state" form:"state,default=1" binding:"oneof=0 1"`
}

type UpdateTagRequest struct {
	ID         uint32 `json:"id" form:"id" binding:"required,gte=1"`
	Name       string `json:"name" form:"name" binding:"max=100"`
	State      uint8  `json:"state" form:"state" binding:"oneof=0 1"`
	ModifiedBy string `json:"modified_by" form:"modified_by" binding:"required,min=2,max=100"`
}

type DeleteTagRequest struct {
	ID uint32 `json:"id" form:"id" binding:"required,gte=1"`
}

func (svc *Service) CountTag(param *CountTagRequest) (int64, error) {
	return svc.dao.CountTag(param.Name, param.State)
}

func (svc *Service) GetTagList(param *TagListRequest, pager *app.Pager) ([]*model.Tag, error) {
	return svc.dao.GetTagList(param.Name, param.State, pager.Page, pager.PageSize)
}

func (svc *Service) CreateTag(param *CreateTagRequest) error {
	exists, err := svc.dao.TagExists(param.Name)
	if err != nil {
		return err
	}
	if exists {
		fmt.Println("true====")
		return errcode.ErrorTagAlreadyExists
	}
	return svc.dao.CreateTag(param.Name, param.State, param.CreatedBy)
}

func (svc *Service) UpdateTag(param *UpdateTagRequest) error {
	return svc.dao.UpdateTag(param.ID, param.Name, param.State, param.ModifiedBy)
}

func (svc *Service) DeleteTag(param *DeleteTagRequest) error {
	return svc.dao.DeleteTag(param.ID)
}
