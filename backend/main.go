package main

import (
	"log"
	"os"

	"api-key-rotator/backend/internal/cache"
	"api-key-rotator/backend/internal/config"
	"api-key-rotator/backend/internal/database"
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

	// 初始化数据库
	db, err := database.Initialize(cfg)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// 初始化缓存
	cacheClient, err := cache.Initialize(cfg)
	if err != nil {
		log.Fatal("Failed to initialize cache:", err)
	}

	// 创建数据库表
	resetTables := os.Getenv("RESET_DB_TABLES")
	if resetTables == "true" {
		log.Println("Resetting database tables...")
		if err := database.ResetTables(db); err != nil {
			log.Fatal("Failed to reset database tables:", err)
		}
		log.Println("Database tables reset successfully")
	} else {
		if err := database.Migrate(db); err != nil {
			log.Fatal("Failed to migrate database:", err)
		}
	}

	// 初始化路由
	r := router.Setup(cfg, db, cacheClient)

	// 启动服务器
	port := os.Getenv("BACKEND_PORT")
	if port == "" {
		port = "8000"
	}

	log.Printf("Starting server on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}