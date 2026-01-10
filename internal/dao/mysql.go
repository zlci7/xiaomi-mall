package dao

import (
	"fmt"
	"log"
	"time"
	"xiaomi-mall/config"
	"xiaomi-mall/internal/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitMySQL() {
	dsn := config.AppConfig.Database.MySQL

	// 根据环境决定日志级别
	var logLevel logger.LogLevel
	if config.AppConfig.Server.Mode == "release" {
		logLevel = logger.Error
	} else {
		logLevel = logger.Info
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})

	if err != nil {
		log.Fatalf("❌ 连接数据库失败: %v", err)
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("❌ 获取数据库实例失败: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// ⬅️ 必须赋值给全局变量！
	DB = db
	fmt.Println("✅ MySQL 连接成功！")

	// 执行数据库迁移
	if err := model.Migrate(DB); err != nil {
		log.Fatalf("❌ 数据库迁移失败: %v", err)
	}
	fmt.Println("✅ 数据库迁移完成！")
}
