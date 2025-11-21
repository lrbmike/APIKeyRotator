package database

import (
	"api-key-rotator/backend/internal/config"
	"api-key-rotator/backend/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Initialize 初始化数据库连接
func Initialize(cfg *config.Config) (*gorm.DB, error) {
	// 使用 SQLite，数据库文件路径从配置中读取
	db, err := gorm.Open(sqlite.Open(cfg.DatabasePath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Migrate 执行数据库迁移
func Migrate(db *gorm.DB) error {
	// 创建表结构
	err := db.AutoMigrate(
		&models.ProxyConfig{},
		&models.APIKey{},
	)
	if err != nil {
		return err
	}

	return nil
}

// ResetTables 重置数据库表（删除并重新创建）
func ResetTables(db *gorm.DB) error {
	// 按依赖顺序删除表
	tables := []interface{}{
		&models.APIKey{},
		&models.ProxyConfig{},
	}

	for _, table := range tables {
		if err := db.Migrator().DropTable(table); err != nil {
			// 忽略表不存在的错误
			continue
		}
	}

	// 重新创建表
	return Migrate(db)
}
