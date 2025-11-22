package database

import (
	"api-key-rotator/backend/internal/config"
	"api-key-rotator/backend/internal/models"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Initialize 初始化数据库连接（工厂模式）
func Initialize(cfg *config.Config) (*gorm.DB, error) {
	switch cfg.DBType {
	case "sqlite":
		return initializeSQLite(cfg)
	case "mysql":
		return initializeMySQL(cfg)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", cfg.DBType)
	}
}

// initializeSQLite 初始化SQLite数据库
func initializeSQLite(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(cfg.DatabasePath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SQLite database at %s: %w", cfg.DatabasePath, err)
	}

	// 对于SQLite，确保数据目录存在
	// 这在Docker环境中尤其重要

	return db, nil
}

// initializeMySQL 初始化MySQL数据库
func initializeMySQL(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(cfg.DatabaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL database: %w", err)
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