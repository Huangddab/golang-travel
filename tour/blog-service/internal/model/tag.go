package model

import (
	"blog-service/pkg/app"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type TagSwagger struct {
	List  []*Tag
	Pager *app.Pager
}

type Tag struct {
	*Model        // 执行运行DB操作的模型实例
	Name   string `json:"name"`  // 标签名称
	State  uint8  `json:"state"` // 状态 (0-禁用 1-启用)
}

func (t *Tag) TableName() string {
	return "blog_tag" // 指定数据库表名
}

// 针对标签模块的模型操作进行封装

// 统计数量
func (t Tag) Count(db *gorm.DB) (int64, error) {
	var count int64
	// 如果名字不为空 则按名称精准匹配
	if t.Name != "" {
		db = db.Where("name = ?", t.Name)
	}
	db = db.Where("state = ?", t.State)
	// is_del软删除：只统计未删除的记录
	if err := db.Model(&t).Where("is_del = ?", 0).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// 获取列表
func (t Tag) List(db *gorm.DB, pageOffset, pageSize int) ([]*Tag, error) {
	var tags []*Tag
	var err error
	if pageOffset >= 0 && pageSize > 0 {
		db = db.Offset(pageOffset).Limit(pageSize)
	}
	if t.Name != "" {
		db = db.Where("name = ?", t.Name)
	}

	db = db.Where("state = ?", t.State)
	if err = db.Where("is_del = ?", 0).Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

// 创建标签
func (t Tag) Create(db *gorm.DB) error {
	return db.Create(&t).Error
}

// 更新标签
func (t Tag) Update(db *gorm.DB, value interface{}) error {
	if db.Model(&Tag{}).Where("id = ? AND is_del = ?", t.ID, 0).Updates(value).Error != nil {
		return db.Error
	}

	return nil
}

// 删除标签
func (t Tag) Delete(db *gorm.DB) error {
	return db.Where("id = ? AND is_del = ?", t.Model.ID, 0).Delete(&t).Error
}

// 查重
func (t Tag) Exists(db *gorm.DB) (bool, error) {
	err := db.Model(&Tag{}).Where("name = ? AND is_del = 0", t.Name).First(&t).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("---1----")
			return false, nil
		}
		fmt.Println("----2---")

		return false, err
	}
	fmt.Println("----3---")
	return true, nil
}
