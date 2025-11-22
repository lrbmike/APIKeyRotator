package sqlite

import (
	"api-key-rotator/backend/internal/infrastructure/database"
)

// Manager SQLite数据库管理器
type Manager struct {
	databasePath string
	repo         database.Repository
}

// NewSQLiteManager 创建SQLite管理器实例
func NewSQLiteManager(databasePath string) *Manager {
	return &Manager{
		databasePath: databasePath,
	}
}

// Initialize 初始化SQLite数据库
func (m *Manager) Initialize() (database.Repository, error) {
	// 使用传入的数据库路径，如果为空则使用默认路径
	databasePath := m.databasePath
	if databasePath == "" {
		databasePath = "./api_key_rotator.db"
	}

	repo, err := NewSQLiteRepository(databasePath)
	if err != nil {
		return nil, err
	}

	m.repo = repo
	return repo, nil
}

// Close 关闭数据库连接
func (m *Manager) Close() error {
	// SQLite的GORM连接通常不需要显式关闭
	return nil
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