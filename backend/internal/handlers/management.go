package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"api-key-rotator/backend/internal/config"
	"api-key-rotator/backend/internal/dto"
	"api-key-rotator/backend/internal/logger"
	"api-key-rotator/backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ManagementHandler 管理API处理器
type ManagementHandler struct {
	cfg *config.Config
	db  *gorm.DB
}

// NewManagementHandler 创建管理处理器实例
func NewManagementHandler(cfg *config.Config, db *gorm.DB) *ManagementHandler {
	return &ManagementHandler{
		cfg: cfg,
		db:  db,
	}
}

// GetAppConfig 获取应用配置
func (h *ManagementHandler) GetAppConfig(c *gin.Context) {
	response := dto.AppConfigResponse{
		ProxyPublicBaseURL: h.cfg.ProxyPublicBaseURL,
	}
	c.JSON(http.StatusOK, response)
}

// Login 处理登录
func (h *ManagementHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Username == h.cfg.AdminUser && req.Password == h.cfg.AdminPassword {
		response := dto.LoginResponse{
			AccessToken: "a_simple_mock_token",
			TokenType:   "bearer",
		}
		c.JSON(http.StatusOK, response)
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"detail": "Incorrect username or password"})
}

// CreateConfig 创建统一的代理配置
func (h *ManagementHandler) CreateConfig(c *gin.Context) {
	var req dto.ProxyConfigCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if strings.Contains(req.Slug, "/") || strings.Contains(req.Slug, " ") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "slug cannot contain slashes or spaces"})
		return
	}

	proxyConfig := models.ProxyConfig{
		Name:           req.Name,
		Slug:           req.Slug,
		ConfigType:     req.ConfigType,
		APIKeyLocation: req.APIKeyLocation,
		APIKeyName:     req.APIKeyName,
		IsActive:       req.IsActive,
		Method:         req.Method,
		TargetURL:      req.TargetURL,
		TargetBaseURL:  req.TargetBaseURL,
		APIFormat:      req.APIFormat,
	}

	if err := h.db.Create(&proxyConfig).Error; err != nil {
		logger.Errorf("Failed to create proxy config: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create configuration"})
		return
	}

	response := dto.ToProxyConfigResponse(&proxyConfig)
	c.JSON(http.StatusCreated, response)
}

// GetAllConfigs 获取所有代理配置
func (h *ManagementHandler) GetAllConfigs(c *gin.Context) {
	skip, _ := strconv.Atoi(c.DefaultQuery("skip", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))

	var proxyConfigs []models.ProxyConfig
	if err := h.db.Preload("APIKeys").Offset(skip).Limit(limit).Find(&proxyConfigs).Error; err != nil {
		logger.Errorf("Failed to get proxy configs: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get configurations"})
		return
	}

	var responses []dto.ProxyConfigResponse
	for _, pc := range proxyConfigs {
		responses = append(responses, dto.ToProxyConfigResponse(&pc))
	}

	c.JSON(http.StatusOK, responses)
}

// GetConfigByID 获取单个代理配置
func (h *ManagementHandler) GetConfigByID(c *gin.Context) {
	id, err := h.parseID(c)
	if err != nil {
		return
	}

	var proxyConfig models.ProxyConfig
	if err := h.db.Preload("APIKeys").First(&proxyConfig, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"detail": "ProxyConfig not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get configuration"})
		return
	}

	response := dto.ToProxyConfigResponse(&proxyConfig)
	c.JSON(http.StatusOK, response)
}

// UpdateConfig 更新统一的代理配置
func (h *ManagementHandler) UpdateConfig(c *gin.Context) {
	id, err := h.parseID(c)
	if err != nil {
		return
	}

	var req dto.ProxyConfigCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var proxyConfig models.ProxyConfig
	if err := h.db.First(&proxyConfig, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"detail": "ProxyConfig not found"})
		return
	}

	// 更新字段
	proxyConfig.Name = req.Name
	proxyConfig.Slug = req.Slug
	proxyConfig.APIKeyLocation = req.APIKeyLocation
	proxyConfig.APIKeyName = req.APIKeyName
	proxyConfig.IsActive = req.IsActive
	proxyConfig.Method = req.Method
	proxyConfig.TargetURL = req.TargetURL
	proxyConfig.TargetBaseURL = req.TargetBaseURL
	proxyConfig.APIFormat = req.APIFormat
    // ConfigType 不应被更新

	if err := h.db.Save(&proxyConfig).Error; err != nil {
		logger.Errorf("Failed to update proxy config: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update configuration"})
		return
	}

	// 重新加载APIKeys以获得完整响应
	h.db.Preload("APIKeys").First(&proxyConfig, id)

	response := dto.ToProxyConfigResponse(&proxyConfig)
	c.JSON(http.StatusOK, response)
}

// DeleteConfig 删除代理配置
func (h *ManagementHandler) DeleteConfig(c *gin.Context) {
	id, err := h.parseID(c)
	if err != nil {
		return
	}

	// GORM的级联删除约束会自动删除关联的APIKeys
	if err := h.db.Delete(&models.ProxyConfig{}, id).Error; err != nil {
		logger.Errorf("Failed to delete config: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete configuration"})
		return
	}

	c.Status(http.StatusNoContent)
}


// GetKeysForConfig 获取指定配置的所有API密钥
func (h *ManagementHandler) GetKeysForConfig(c *gin.Context) {
	id, err := h.parseID(c)
	if err != nil {
		return
	}

	var config models.ProxyConfig
	if err := h.db.Preload("APIKeys").First(&config, id).Error; err != nil {
		logger.Errorf("Failed to find ProxyConfig with ID %d: %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{"detail": "ProxyConfig not found"})
		return
	}

	c.JSON(http.StatusOK, config.APIKeys)
}

// CreateAPIKeyForConfig 为指定配置创建API密钥
func (h *ManagementHandler) CreateAPIKeyForConfig(c *gin.Context) {
	id, err := h.parseID(c)
	if err != nil {
		return
	}

	var req dto.APIKeyCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var config models.ProxyConfig
	if err := h.db.First(&config, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"detail": "ProxyConfig not found"})
		return
	}

	apiKey := &models.APIKey{
		KeyValue:      req.KeyValue,
		IsActive:      req.IsActive,
		ProxyConfigID: id,
	}

	if err := h.db.Create(apiKey).Error; err != nil {
		logger.Errorf("Failed to create API key: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create API key"})
		return
	}

	c.JSON(http.StatusCreated, apiKey)
}

// UpdateAPIKeyStatus 更新API密钥状态
func (h *ManagementHandler) UpdateAPIKeyStatus(c *gin.Context) {
	keyID, err := h.parseID(c)
	if err != nil {
		return
	}

	var req dto.APIKeyStatusUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var apiKey models.APIKey
	if err := h.db.First(&apiKey, keyID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"detail": "ApiKey not found"})
		return
	}

	apiKey.IsActive = req.IsActive
	if err := h.db.Save(&apiKey).Error; err != nil {
		logger.Errorf("Failed to update API key status: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update API key"})
		return
	}

	c.JSON(http.StatusOK, apiKey)
}

// DeleteAPIKey 删除API密钥
func (h *ManagementHandler) DeleteAPIKey(c *gin.Context) {
	keyID, err := h_parseKeyID(c)
	if err != nil {
		return
	}

	var apiKey models.APIKey
	if err := h.db.First(&apiKey, keyID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"detail": "ApiKey not found"})
		return
	}

	if err := h.db.Delete(&apiKey).Error; err != nil {
		logger.Errorf("Failed to delete API key: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete API key"})
		return
	}

	c.Status(http.StatusNoContent)
}

// parseID 是一个辅助函数，用于从URL参数解析ID
func (h *ManagementHandler) parseID(c *gin.Context) (int32, error) {
	id64, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return 0, err
	}
	return int32(id64), nil
}

// h_parseKeyID 是一个辅助函数，用于从URL参数解析keyID
func h_parseKeyID(c *gin.Context) (int32, error) {
	id64, err := strconv.ParseInt(c.Param("keyID"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Key ID"})
		return 0, err
	}
	return int32(id64), nil
}