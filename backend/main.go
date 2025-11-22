package main

import (
	"log"
	"os"

	"api-key-rotator/backend/internal/config"
	"api-key-rotator/backend/internal/logger"
	"api-key-rotator/backend/internal/router"

	"github.com/joho/godotenv"
)

func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// 初始化日志
	logger.Setup()

	// 加载配置
	cfg := config.Load()

	// 打印当前配置信息
	log.Printf("Database Type: %s", cfg.DBType)
	log.Printf("Cache Type: %s", cfg.CacheType)

	// 创建基础设施工厂
	factory := config.NewInfrastructureFactory(cfg)

	// 初始化数据库仓库
	dbRepo, err := factory.CreateDatabaseRepository()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// 初始化缓存接口
	cacheInterface, err := factory.CreateCacheInterface()
	if err != nil {
		log.Fatal("Failed to initialize cache:", err)
	}

	// 创建数据库表
	resetTables := os.Getenv("RESET_DB_TABLES")
	if resetTables == "true" {
		log.Println("Resetting database tables...")
		if err := dbRepo.Reset(); err != nil {
			log.Fatal("Failed to reset database tables:", err)
		}
		log.Println("Database tables reset successfully")
	} else {
		if err := dbRepo.Migrate(); err != nil {
			log.Fatal("Failed to migrate database:", err)
		}
	}

	// 初始化路由
	r := router.Setup(cfg, dbRepo, cacheInterface)

	log.Println("Backend services initialized successfully")
	log.Printf("Database: tables migrated successfully")
	log.Printf("Cache: %s interface initialized", cfg.CacheType)

	// 启动Web服务器
	log.Printf("Starting server on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}