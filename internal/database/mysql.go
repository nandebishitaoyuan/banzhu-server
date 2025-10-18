package database

import (
	"fmt"
	"httpServerTest/internal/config"
	"httpServerTest/internal/model"
	"reflect"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

// Init 初始化数据库
func Init(cfg *config.Config) error {
	db, err := gorm.Open(mysql.Open(cfg.Database.DSN), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
			TablePrefix: "t_",
		},
	})

	if err != nil {
		return err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(1 * time.Hour)

	DB = db

	return nil
}

// Paginate 分页工具
func Paginate[T any](db *gorm.DB, param *model.PageParam) (*model.PageResult[T], error) {
	var total int64
	offset := (param.PageIndex - 1) * param.PageSize
	var data []T
	db = db.Model(new(T)).Count(&total).Limit(param.PageSize).Offset(offset)

	if param.Order != "" {

		// 判断排序是否包含 ASC 或 DESC
		if !strings.Contains(param.Order, "ASC") && !strings.Contains(param.Order, "DESC") {
			return nil, fmt.Errorf("排序值无效，必须包含 ASC 或 DESC")
		}

		db = db.Order(param.Order)
	}

	if err := db.Find(&data).Error; err != nil {
		return nil, err
	}

	return &model.PageResult[T]{
		Total:     total,
		PageIndex: param.PageIndex,
		PageSize:  param.PageSize,
		Data:      data,
	}, nil
}

// ApplyKeywordSearch 对model中标注了like字段生成Sql模糊查询
func ApplyKeywordSearch[T any](db *gorm.DB, keyword string) *gorm.DB {
	if keyword == "" {
		return db
	}

	var entity T
	typ := reflect.TypeOf(entity)

	// 反射遍历字段
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.Tag.Get("like") == "true" {
			colName := DB.NamingStrategy.ColumnName("", field.Name)
			db = db.Or(fmt.Sprintf("%s LIKE ?", colName), "%"+keyword+"%")
		}
	}
	return db
}
