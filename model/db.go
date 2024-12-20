package model

import (
	"fmt"
	"Go_blog/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"time"
	"gorm.io/gorm/logger"
	"os"
)

var db *gorm.DB
var err error

func InitDb() {
	dns :=  fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
			utils.DbUser,
			utils.DbPassword,
			utils.DbHost,
			utils.DbPort,
			utils.DbName,
	)
	db, err = gorm.Open(mysql.Open(dns), &gorm.Config{
		// gorm日志模式:silent
		Logger: logger.Default.LogMode(logger.Silent),
		// 禁用外键约束
		DisableForeignKeyConstraintWhenMigrating: true,
		// 禁用默认事务
		SkipDefaultTransaction: true,
		// 命名策略
		NamingStrategy: schema.NamingStrategy{
			// 使用单数表名，启用该选项，此时，`User` 的表名应该是 `user`
			SingularTable: true,
		},
	})
	if err != nil {
		fmt.Println("连接数据库失败，请检查参数：", err)
		os.Exit(1)
	}

	// 迁移数据表，在没有数据表结构变化时，不要轻易执行
	//_ = db.AutoMigrate(&User{}, &Article{}, &Category{}, Profile{}, Comment{})
	_=db.AutoMigrate(&User{}, &Article{}, &Category{})

	sqlDB, _ := db.DB()

	// 设置连接池
	sqlDB.SetMaxIdleConns(10)

	// 设置最大连接数
	sqlDB.SetMaxOpenConns(100)

	// 设置连接的最大可复用时间
	sqlDB.SetConnMaxLifetime(10 * time.Second)
}

// GetDB 返回数据库实例
func GetDB() *gorm.DB {
	return db
}