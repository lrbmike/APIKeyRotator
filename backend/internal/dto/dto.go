package dto

import "api-key-rotator/backend/internal/models"

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

// AppConfigResponse 应用配置响应
type AppConfigResponse struct {
	ProxyPublicBaseURL string `json:"proxy_public_base_url"`
}

// ProxyConfigCreate 创建或更新代理配置的统一请求
type ProxyConfigCreate struct {
	Name           string  `json:"name" binding:"required"`
	Slug           string  `json:"slug" binding:"required"`
	ConfigType     string  `json:"config_type" binding:"required"` // "generic" or "llm"
	APIKeyLocation *string `json:"api_key_location,omitempty"`
	APIKeyName     *string `json:"api_key_name,omitempty"`
	IsActive       bool    `json:"is_active"`
	Method         *string `json:"method,omitempty"`
	TargetURL      *string `json:"target_url,omitempty"`
	TargetBaseURL  *string `json:"target_base_url,omitempty"`
	APIFormat      *string `json:"api_format,omitempty"`
	OutputFormat   *string `json:"output_format,omitempty"`
}

// ProxyConfigStatusUpdate 更新代理配置状态的请求
type ProxyConfigStatusUpdate struct {
	IsActive bool `json:"is_active"`
}

// APIKeyCreate 创建API密钥请求
type APIKeyCreate struct {
	KeyValue string `json:"key_value" binding:"required"`
	IsActive bool   `json:"is_active"`
}

// APIKeyStatusUpdate 更新API密钥状态请求
type APIKeyStatusUpdate struct {
	IsActive bool `json:"is_active"`
}

// ProxyConfigResponse 代理配置的统一响应
type ProxyConfigResponse struct {
	ID             int32           `json:"id"`
	Name           string          `json:"name"`
	Slug           string          `json:"slug"`
	ConfigType     string          `json:"config_type"`
	APIKeyLocation *string         `json:"api_key_location,omitempty"`
	APIKeyName     *string         `json:"api_key_name,omitempty"`
	IsActive       bool            `json:"is_active"`
	APIKeys        []models.APIKey `json:"api_keys"`
	Method         *string         `json:"method,omitempty"`
	TargetURL      *string         `json:"target_url,omitempty"`
	TargetBaseURL  *string         `json:"target_base_url,omitempty"`
	APIFormat      *string         `json:"api_format,omitempty"`
	OutputFormat   *string         `json:"output_format,omitempty"`
}

// ToProxyConfigResponse 将模型转换为响应DTO
func ToProxyConfigResponse(proxyConfig *models.ProxyConfig) ProxyConfigResponse {
	resp := ProxyConfigResponse{
		ID:            proxyConfig.ID,
		Name:          proxyConfig.Name,
		Slug:          proxyConfig.Slug,
		ConfigType:    proxyConfig.ConfigType,
		IsActive:      proxyConfig.IsActive,
		APIKeys:       proxyConfig.APIKeys,
		Method:        proxyConfig.Method,
		TargetURL:     proxyConfig.TargetURL,
		TargetBaseURL: proxyConfig.TargetBaseURL,
		APIFormat:     proxyConfig.APIFormat,
		OutputFormat:  proxyConfig.OutputFormat,
	}

	if proxyConfig.APIKeyLocation != nil {
		resp.APIKeyLocation = proxyConfig.APIKeyLocation
	}
	if proxyConfig.APIKeyName != nil {
		resp.APIKeyName = proxyConfig.APIKeyName
	}

	return resp
}

// BatchAPIKeyCreate 批量创建API密钥请求
type BatchAPIKeyCreate struct {
	Keys []string `json:"keys" binding:"required"`
}

// BatchAPIKeyImportResponse 批量导入API密钥响应
type BatchAPIKeyImportResponse struct {
	SuccessCount int      `json:"success_count"`
	FailedCount  int      `json:"failed_count"`
	SkippedCount int      `json:"skipped_count"`
	FailedKeys   []string `json:"failed_keys,omitempty"`
}

// ClearAllAPIKeysResponse 清除所有API密钥响应
type ClearAllAPIKeysResponse struct {
	DeletedCount int `json:"deleted_count"`
}
