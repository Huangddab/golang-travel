package model

import (
	"blog-service/global"
	"blog-service/pkg/setting"
	"fmt"
	"reflect"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Model struct {
	ID         uint32 `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	CreatedOn  int64  `json:"created_on"`
	ModifiedOn int64  `json:"modified_on"`
	DeletedOn  int64  `json:"deleted_on"`
	IsDel      int8   `json:"is_del"`
}

func setupCallbacks(db *gorm.DB) {
	// 注册创建回调
	db.Callback().Create().Before("gorm:create").Register("update_time_stamp", updateTimeStampForCreateCallback)

	// 注册更新回调
	db.Callback().Update().Before("gorm:update").Register("update_time_stamp", updateTimeStampForUpdateCallback)

	// 注册删除回调
	db.Callback().Delete().Before("gorm:delete").Register("soft_delete", deleteCallback)
}

func NewDBEngine(databaseSetting *setting.DatabaseSettingS) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=Local",
		databaseSetting.UserName,
		databaseSetting.Password,
		databaseSetting.Host,
		databaseSetting.DBName,
		databaseSetting.Charset,
		databaseSetting.ParseTime,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(databaseSetting.MaxIdleConns)
	sqlDB.SetMaxOpenConns(databaseSetting.MaxOpenConns)

	if global.ServerSetting.RunMode == "debug" {
		db.Logger = logger.Default.LogMode(logger.Info)
	}

	setupCallbacks(db)
	return db, nil
}

func updateTimeStampForCreateCallback(db *gorm.DB) {
	if db.Error != nil {
		return
	}
	// 获取当前模型的值
	// reflect
	// ValueOf用来获取输入参数接口中的数据的值，如果接口为空则返回0
	// TypeOf用来动态获取输入参数接口中的值的类型，如果接口为空则返回nil
	modelValue := reflect.ValueOf(db.Statement.Model)
	if modelValue.Kind() == reflect.Ptr {
		modelValue = modelValue.Elem()
	}
	now := time.Now()

	// 设置CreatedOn
	if createdOnField := modelValue.FieldByName("CreatesOn"); createdOnField.IsValid() {
		if createdOnField.Int() == 0 { // 检查是否为0
			db.Statement.SetColumn("CreatedOn", now)
		}
	}
	// 设置ModifiedOn
	if modifiedOnField := modelValue.FieldByName("ModifiedOn"); modifiedOnField.IsValid() {
		if modifiedOnField.Int() == 0 { // 检查是否为0
			db.Statement.SetColumn("ModifiedOn", now)
		}
	}
}

// 更新回调
func updateTimeStampForUpdateCallback(db *gorm.DB) {
	if db.Error != nil {
		return
	}

	// 检查是否已经手动设置了更新列
	if _, ok := db.Statement.Clauses["SET"]; !ok {
		db.Statement.SetColumn("ModifiedOn", time.Now().Unix())
	}
}

// 删除回调
func deleteCallback(db *gorm.DB) {
	if db.Error != nil {
		return
	}

	modelValue := reflect.ValueOf(db.Statement.Model)
	if modelValue.Kind() == reflect.Ptr {
		modelValue = modelValue.Elem()
	}

	hasDeletedOn := modelValue.FieldByName("DeletedOn").IsValid()
	hasIsDel := modelValue.FieldByName("IsDel").IsValid()

	if !db.Statement.Unscoped && hasDeletedOn && hasIsDel {
		now := time.Now().Unix()

		// 新版GORM的正确写法：使用主键执行软删除的 UPDATE
		tableName := db.Statement.Table

		// 尝试从模型中读取主键 ID 值（大多数模型使用 uint32 类型的 ID 字段）
		var id interface{}
		if idField := modelValue.FieldByName("ID"); idField.IsValid() {
			id = idField.Interface()
		}

		if id != nil && id != 0 {
			// 使用 fmt.Sprintf 拼接表名并通过参数传入值，避免 GORM 对占位符的错误替换
			sql := fmt.Sprintf("UPDATE %s SET deleted_on = ?, is_del = ? WHERE id = ?", tableName)
			db.Exec(sql, now, 1, id)
		}
	}
	// 否则让GORM执行默认的DELETE

}
