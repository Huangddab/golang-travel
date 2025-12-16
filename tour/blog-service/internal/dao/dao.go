package dao

import (
	"gorm.io/gorm"
)


// dao层是干嘛？数据访问对象层 专门负责与数据库打交道的中间层
type Dao struct {
	engine *gorm.DB
}

func New(engine *gorm.DB) *Dao {
	return &Dao{engine: engine}
}
