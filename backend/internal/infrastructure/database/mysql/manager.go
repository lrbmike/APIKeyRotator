package mysql

import (
	"api-key-rotator/backend/internal/infrastructure/database"
)

// Manager MySQL数据库管理器
type Manager struct {
	repo database.Repository
}

// NewMySQLManager 创建MySQL管理器实例
func NewMySQLManager() *Manager {
	return &Manager{}
}

// Initialize 初始化MySQL数据库
func (m *Manager) Initialize() (database.Repository, error) {
	// 这里应该从配置中获取数据库连接URL
	// 暂时使用默认的连接字符串
	databaseURL := "user:password@tcp(localhost:3306)/database?charset=utf8mb4&parseTime=True&loc=Local"

	repo, err := NewMySQLRepository(databaseURL)
	if err != nil {
		return nil, err
	}

	m.repo = repo
	return repo, nil
}

// Close 关闭数据库连接
func (m *Manager) Close() error {
	if m.repo == nil {
		return nil
	}

	sqlDB, err := m.repo.GetDB().DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

// Ping 测试数据库连接
func (m *Manager) Ping() error {
	if m.repo == nil {
		return nil
	}

	sqlDB, err := m.repo.GetDB().DB()
	if err != nil {
		return err
	}

	return sqlDB.Ping()
}