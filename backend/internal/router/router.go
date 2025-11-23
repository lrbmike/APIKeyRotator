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
	r.Static("/static", "./static")

	// 添加SPA支持 - 对于非API路径，返回index.html
	r.NoRoute(func(c *gin.Context) {
		// 检查是否是API请求
		if c.Request.URL.Path[0:1] == "/" &&
		   (c.Request.URL.Path == "/admin" ||
		    c.Request.URL.Path[:5] == "/admin" ||
		    c.Request.URL.Path[:4] == "/api/" ||
		    c.Request.URL.Path == "/api") {
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

	// 根路径 - 返回前端应用
	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
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
		adminAPI.PUT("/proxy-configs/:id/status", managementHandler.UpdateConfigStatus)
		adminAPI.DELETE("/proxy-configs/:id", managementHandler.DeleteConfig)

		// API密钥管理
		adminAPI.GET("/proxy-configs/:id/keys", managementHandler.GetKeysForConfig)
		adminAPI.POST("/proxy-configs/:id/keys", managementHandler.CreateAPIKeyForConfig)
		adminAPI.PATCH("/keys/:keyID", managementHandler.UpdateAPIKeyStatus)
		adminAPI.DELETE("/keys/:keyID", managementHandler.DeleteAPIKey)
	}

	// 通用代理路由组 - 暂时禁用
	// TODO: 重新实现代理处理器以支持新的接口抽象架构
	// proxyGroup := r.Group("/proxy")
	// {
	//     proxyGroup.Any("/*slug", proxyHandler.HandleGenericProxy)
	// }

	// LLM代理路由组 - 暂时禁用
	// TODO: 重新实现LLM代理处理器以支持新的接口抽象架构
	// llmGroup := r.Group("/llm")
	// {
	//     llmGroup.Any("/:slug/*action", llmProxyHandler.HandleLLMProxy)
	// }

	return r
}