package router

import (
	"fmt"

	"api-key-rotator/backend/internal/config"
	"api-key-rotator/backend/internal/handlers"
	"api-key-rotator/backend/internal/middleware"
	"api-key-rotator/backend/internal/infrastructure/database"
	"api-key-rotator/backend/internal/infrastructure/cache"

	"github.com/gin-gonic/gin"
)

// Setup 设置路由
func Setup(cfg *config.Config, dbRepo database.Repository, cacheInterface cache.CacheInterface) *gin.Engine {
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

	// 静态文件服务 - 为前端提供静态资源
	r.StaticFile("/", "./static/index.html")
	r.Static("/assets", "./static/assets")

	// 添加SPA支持 - 对于非API路径，返回index.html
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		// 检查是否是API请求
		if len(path) >= 6 && path[:6] == "/admin" ||
		   len(path) >= 7 && path[:7] == "/proxy" ||
		   len(path) >= 5 && path[:5] == "/llm" {
			// API路径返回404
			fmt.Printf("404 Not Found: %s %s\n", c.Request.Method, c.Request.URL.Path)
			c.JSON(404, gin.H{"error": "Route not found"})
		} else {
			// 非API路径，返回前端index.html
			c.File("./static/index.html")
		}
	})

	// 创建处理器实例，使用完整版本
	managementHandler := handlers.NewManagementHandler(cfg, dbRepo)
	proxyHandler := handlers.NewProxyHandler(cfg, dbRepo.GetDB(), cacheInterface)
	llmProxyHandler := handlers.NewLLMProxyHandler(cfg, dbRepo.GetDB(), cacheInterface)

	
	// 管理API路由组 - 后台管理接口
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
		adminAPI.PUT("/proxy-configs/:id/status", managementHandler.UpdateConfigStatus)
		adminAPI.DELETE("/proxy-configs/:id", managementHandler.DeleteConfig)

		// API密钥管理
		adminAPI.GET("/proxy-configs/:id/keys", managementHandler.GetKeysForConfig)
		adminAPI.POST("/proxy-configs/:id/keys", managementHandler.CreateAPIKeyForConfig)
		adminAPI.POST("/proxy-configs/:id/keys/batch", managementHandler.BatchCreateAPIKeys)
		adminAPI.DELETE("/proxy-configs/:id/keys", managementHandler.ClearAllAPIKeys)
		adminAPI.PATCH("/keys/:keyID", managementHandler.UpdateAPIKeyStatus)
		adminAPI.DELETE("/keys/:keyID", managementHandler.DeleteAPIKey)
	}

	// 通用代理路由组 - 公开API接口
	proxyGroup := r.Group("/proxy")
	proxyGroup.Any("/*slug", proxyHandler.HandleGenericProxy)

	// LLM代理路由组 - 公开API接口
	llmGroup := r.Group("/llm")
	llmGroup.Any("/:slug/*action", llmProxyHandler.HandleLLMProxy)

	return r
}