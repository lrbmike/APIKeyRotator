package router

import (
	"fmt"
	
	"api-key-rotator/backend/internal/config"
	"api-key-rotator/backend/internal/handlers"
	"api-key-rotator/backend/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// Setup 设置路由
func Setup(cfg *config.Config, db *gorm.DB, redisClient *redis.Client) *gin.Engine {
	// 设置Gin模式为调试模式以便看到更多日志
	gin.SetMode(gin.DebugMode)

	r := gin.New()

	// 添加中间件
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format("02/Jan/2006:15:04:05 -0700"),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())

	// 添加全局请求调试中间件
	r.Use(func(c *gin.Context) {
		fmt.Printf("Global Request: %s %s\n", c.Request.Method, c.Request.URL.Path)
		c.Next()
	})

	// 添加404处理器
	r.NoRoute(func(c *gin.Context) {
		fmt.Printf("404 Not Found: %s %s\n", c.Request.Method, c.Request.URL.Path)
		c.JSON(404, gin.H{"error": "Route not found"})
	})

	// 创建处理器实例
	managementHandler := handlers.NewManagementHandler(cfg, db)
	proxyHandler := handlers.NewProxyHandler(cfg, db, redisClient)
	llmProxyHandler := handlers.NewLLMProxyHandler(cfg, db, redisClient)

	// 根路径
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to the API Key Rotator",
		})
	})

	// 管理API路由组
	adminAPI := r.Group("/admin")
	{
		// 添加路由级别的调试中间件
		adminAPI.Use(func(c *gin.Context) {
			fmt.Printf("Admin API Request: %s %s\n", c.Request.Method, c.Request.URL.Path)
			c.Next()
		})

		// 应用配置和认证
		adminAPI.GET("/app-config", managementHandler.GetAppConfig)
		adminAPI.POST("/login", managementHandler.Login)

		// 代理配置管理
		adminAPI.POST("/proxy-configs", managementHandler.CreateConfig)
		adminAPI.GET("/proxy-configs", managementHandler.GetAllConfigs)
		adminAPI.GET("/proxy-configs/:id", managementHandler.GetConfigByID)
		adminAPI.PUT("/proxy-configs/:id", managementHandler.UpdateConfig)
		adminAPI.DELETE("/proxy-configs/:id", managementHandler.DeleteConfig)

		// API密钥管理
		adminAPI.GET("/proxy-configs/:id/keys", managementHandler.GetKeysForConfig)
		adminAPI.POST("/proxy-configs/:id/keys", managementHandler.CreateAPIKeyForConfig)
		adminAPI.PATCH("/keys/:keyID", managementHandler.UpdateAPIKeyStatus)
		adminAPI.DELETE("/keys/:keyID", managementHandler.DeleteAPIKey)
	}

	// 通用代理路由组
	proxyGroup := r.Group("/proxy")
	{
		proxyGroup.Any("/*slug", proxyHandler.HandleGenericProxy)
	}

	// LLM代理路由组
	llmGroup := r.Group("/llm")
	{
		llmGroup.Any("/:slug/*action", llmProxyHandler.HandleLLMProxy)
	}

	return r
}