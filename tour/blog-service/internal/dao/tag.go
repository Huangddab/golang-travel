package dao

import (
	"blog-service/internal/model"
)

// 用于处理标签模块的 dao 操作
func (d *Dao) CountTag(name string, state uint8) (int64, error) {
	tag := model.Tag{Name: name, State: state}
	return tag.Count(d.engine)
}


