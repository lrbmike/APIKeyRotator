package models

import (
	"time"
)

// ProxyConfig 统一的代理配置模型
type ProxyConfig struct {
	ID             int32     `json:"id" gorm:"primaryKey"`
	Name           string    `json:"name" gorm:"uniqueIndex;size:100;not null"`
	Slug           string    `json:"slug" gorm:"uniqueIndex;size:100;not null"`
	ConfigType     string    `json:"config_type" gorm:"size:50;not null;index"` // "generic" or "llm"
	APIKeyLocation *string   `json:"api_key_location,omitempty" gorm:"size:50"`
	APIKeyName     *string   `json:"api_key_name,omitempty" gorm:"size:100"`
	IsActive       bool      `json:"is_active" gorm:"default:true"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Fields from GenericAPIConfig (nullable)
	Method    *string `json:"method,omitempty" gorm:"size:10"`
	TargetURL *string `json:"target_url,omitempty" gorm:"size:255"`

	// Fields from LLMAPIConfig (nullable)
	TargetBaseURL *string `json:"target_base_url,omitempty" gorm:"size:255"`
	APIFormat     *string `json:"api_format,omitempty" gorm:"size:50;default:openai_compatible"`

	// 关系: 一个配置可以有多个 API Key
	APIKeys []APIKey `json:"api_keys" gorm:"foreignKey:ProxyConfigID;constraint:OnDelete:CASCADE"`
}

// APIKey API密钥模型
type APIKey struct {
	ID            int32        `json:"id" gorm:"primaryKey"`
	KeyValue      string       `json:"key_value" gorm:"size:255;not null"`
	IsActive      bool         `json:"is_active" gorm:"default:true"`
	ProxyConfigID int32        `json:"proxy_config_id"`
	ProxyConfig   *ProxyConfig `json:"-" gorm:"foreignKey:ProxyConfigID"`
}

// TableName 设置ProxyConfig表名
func (ProxyConfig) TableName() string {
	return "proxy_configs"
}

// TableName 设置APIKey表名
func (APIKey) TableName() string {
	return "api_keys"
}
