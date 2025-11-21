package config

import (
	"os"
	"strconv"
	"strings"
)

// Config 应用配置结构
type Config struct {
	// 数据库配置
	DatabasePath string // SQLite 数据库文件路径

	// 服务器配置
	Port string

	// JWT配置
	JWTSecret string

	// 管理员配置
	AdminUsername string
	AdminPassword string
	AdminUser     string // 别名，兼容性

	// 代理配置
	ProxyTimeout       int
	GlobalProxyKeys    string // 逗号分隔的多个密钥，也支持单个密钥
	ProxyPublicBaseURL string

	// 日志配置
	LogLevel string
}

// GetGlobalProxyKeys 获取所有有效的代理密钥列表
// 支持单个密钥或逗号分隔的多个密钥
func (c *Config) GetGlobalProxyKeys() []string {
	// 分割逗号分隔的密钥，并去除空白字符
	keys := strings.Split(c.GlobalProxyKeys, ",")
	var result []string
	for _, key := range keys {
		trimmed := strings.TrimSpace(key)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// Load 加载配置
func Load() *Config {
	// 加载敏感或必需的配置，如果没有设置则会 panic
	jwtSecret := getRequiredEnv("JWT_SECRET")
	adminPassword := getRequiredEnv("ADMIN_PASSWORD")
	globalProxyKeys := getRequiredEnv("GLOBAL_PROXY_KEYS")

	// 加载可选配置，如果未设置则使用默认值
	adminUsername := getEnv("ADMIN_USERNAME", "admin")

	config := &Config{
		// 必填项
		JWTSecret:       jwtSecret,
		AdminPassword:   adminPassword,
		GlobalProxyKeys: globalProxyKeys,

		// 可选项
		DatabasePath:       getEnv("DATABASE_PATH", "./data/api_key_rotator.db"),
		Port:               getEnv("BACKEND_PORT", "8000"),
		AdminUsername:      adminUsername,
		AdminUser:          adminUsername, // 别名，兼容性
		ProxyTimeout:       getEnvAsInt("PROXY_TIMEOUT", 30),
		ProxyPublicBaseURL: getEnv("PROXY_PUBLIC_BASE_URL", "http://localhost:8000"),
		LogLevel:           getEnv("LOG_LEVEL", "info"),
	}

	return config
}

// getRequiredEnv 获取一个必需的环境变量，如果不存在则 panic
func getRequiredEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic("FATAL: Required environment variable not set: " + key)
	}
	return value
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt 获取环境变量并转换为整数，如果不存在或转换失败则返回默认值
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
