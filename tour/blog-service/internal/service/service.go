package service

import (
	"blog-service/global"
	"blog-service/internal/dao"
	"context"
)

// 业务层
type Service struct {
	ctx context.Context
	dao *dao.Dao
}

func New(ctx context.Context) Service {
	svc := Service{ctx: ctx}
	svc.dao = dao.New(global.DBEngine)
	return svc
}
