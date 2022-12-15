package models

import (
	"Yiban3/FrontWeb/helper"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"time"
)

// 在其它model的实体类中可直接调用
var db *gorm.DB

func InitDB(connStr string) (*gorm.DB, *sql.DB) {
	var err error
	db, err = gorm.Open(mysql.Open(connStr), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true}, // 使用单数表名
	})
	if err != nil {
		helper.LogError(err.Error())
		panic("系统错误, 连接数据库失败")
	}
	sqlDB, err := db.DB()
	if err != nil {
		helper.LogError(err.Error())
		panic("系统错误, 数据库对象创建失败")
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, sqlDB
}

// Paginate 分页封装
func Paginate(page int, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

// PageQuery 分页查询参数结构体
type PageQuery struct {
	Page int
	Size int
}

// Option 定义配置选项函数（关键）
type Option func(*PageQuery)

func SetPage(page int) Option {
	// 返回一个Option类型的函数（闭包）：接受ExampleClient类型指针参数并修改之
	return func(this *PageQuery) {
		this.Page = page
	}
}
func SetSize(size int) Option {
	// 返回一个Option类型的函数（闭包）：接受ExampleClient类型指针参数并修改之
	return func(this *PageQuery) {
		this.Size = size
	}
}

// NewPageQuery 应用函数选项配置（样例）
func NewPageQuery(opts ...Option) PageQuery {
	// 初始化默认值
	defaultClient := PageQuery{
		Page: 1,
		Size: 10,
	}

	// 依次调用opts函数列表中的函数，为结构体成员赋值
	for _, o := range opts {
		o(&defaultClient)
	}

	return defaultClient
}
