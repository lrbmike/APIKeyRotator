package handlers

import (
	"net/http"
	"strconv"

	"api-key-rotator/backend/internal/config"
	"api-key-rotator/backend/internal/dto"
	"api-key-rotator/backend/internal/infrastructure/database"
	"api-key-rotator/backend/internal/models"

	"github.com/gin-gonic/gin"
)

// ManagementHandler 管理API处理器 - 使用接口抽象架构
type ManagementHandler struct {
	cfg    *config.Config
	dbRepo database.Repository
}

// NewManagementHandler 创建管理处理器实例
func NewManagementHandler(cfg *config.Config, dbRepo database.Repository) *ManagementHandler {
	return &ManagementHandler{
		cfg:    cfg,
		dbRepo: dbRepo,
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
	var loginReq dto.LoginRequest

	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 从环境变量获取管理员凭据
	adminUsername := h.cfg.AdminUsername
	adminPassword := h.cfg.AdminPassword

	if loginReq.Username == adminUsername && loginReq.Password == adminPassword {
		// 简单的成功响应（实际生产中应该使用JWT）
		c.JSON(http.StatusOK, gin.H{
			"message": "Login successful",
			"success": true,
		})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Invalid credentials",
			"success": false,
		})
	}
}

// CreateConfig 创建代理配置
func (h *ManagementHandler) CreateConfig(c *gin.Context) {
	var req dto.ProxyConfigCreate

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 创建代理配置
	config := &models.ProxyConfig{
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

	if err := h.dbRepo.CreateProxyConfig(config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := dto.ToProxyConfigResponse(config)
	c.JSON(http.StatusCreated, response)
}

// GetAllConfigs 获取所有代理配置
func (h *ManagementHandler) GetAllConfigs(c *gin.Context) {
	configs, err := h.dbRepo.ListProxyConfigs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 转换为响应格式
	var response []dto.ProxyConfigResponse
	for _, config := range configs {
		response = append(response, dto.ToProxyConfigResponse(config))
	}

	c.JSON(http.StatusOK, response)
}

// GetConfigByID 根据ID获取代理配置
func (h *ManagementHandler) GetConfigByID(c *gin.Context) {
	idStr := c.Param("id")
	id64, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	id := uint(id64)

	config, err := h.dbRepo.GetProxyConfigByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Config not found"})
		return
	}

	response := dto.ToProxyConfigResponse(config)
	c.JSON(http.StatusOK, response)
}

// UpdateConfig 更新代理配置
func (h *ManagementHandler) UpdateConfig(c *gin.Context) {
	idStr := c.Param("id")
	id64, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	id := uint(id64)

	var req dto.ProxyConfigCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取现有配置
	config, err := h.dbRepo.GetProxyConfigByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Config not found"})
		return
	}

	// 更新字段
	config.Name = req.Name
	config.Slug = req.Slug
	config.ConfigType = req.ConfigType
	config.APIKeyLocation = req.APIKeyLocation
	config.APIKeyName = req.APIKeyName
	config.IsActive = req.IsActive
	config.Method = req.Method
	config.TargetURL = req.TargetURL
	config.TargetBaseURL = req.TargetBaseURL
	config.APIFormat = req.APIFormat

	if err := h.dbRepo.UpdateProxyConfig(config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := dto.ToProxyConfigResponse(config)
	c.JSON(http.StatusOK, response)
}

// UpdateConfigStatus 更新代理配置状态
func (h *ManagementHandler) UpdateConfigStatus(c *gin.Context) {
	idStr := c.Param("id")
	id64, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	id := uint(id64)

	var req dto.ProxyConfigStatusUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取现有配置
	config, err := h.dbRepo.GetProxyConfigByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Config not found"})
		return
	}

	// 更新状态
	config.IsActive = req.IsActive

	if err := h.dbRepo.UpdateProxyConfig(config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status updated successfully"})
}

// DeleteConfig 删除代理配置
func (h *ManagementHandler) DeleteConfig(c *gin.Context) {
	idStr := c.Param("id")
	id64, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	id := uint(id64)

	if err := h.dbRepo.DeleteProxyConfig(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Config deleted successfully"})
}

// GetKeysForConfig 获取配置的API密钥
func (h *ManagementHandler) GetKeysForConfig(c *gin.Context) {
	idStr := c.Param("id")
	id64, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	id := uint(id64)

	// 检查配置是否存在
	config, err := h.dbRepo.GetProxyConfigByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Config not found"})
		return
	}

	// 返回配置关联的API密钥
	c.JSON(http.StatusOK, config.APIKeys)
}

// CreateAPIKeyForConfig 为配置创建API密钥
func (h *ManagementHandler) CreateAPIKeyForConfig(c *gin.Context) {
	idStr := c.Param("id")
	configID64, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid config ID"})
		return
	}
	configID := uint(configID64)

	var req dto.APIKeyCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查配置是否存在
	_, err = h.dbRepo.GetProxyConfigByID(uint(configID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Config not found"})
		return
	}

	// 创建API密钥
	apiKey := &models.APIKey{
		KeyValue:      req.KeyValue,
		IsActive:      req.IsActive,
		ProxyConfigID: int32(configID64),
	}

	if err := h.dbRepo.CreateAPIKey(apiKey); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, apiKey)
}

// UpdateAPIKeyStatus 更新API密钥状态
func (h *ManagementHandler) UpdateAPIKeyStatus(c *gin.Context) {
	keyIDStr := c.Param("keyID")
	keyID64, err := strconv.ParseInt(keyIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid key ID"})
		return
	}
	keyID := uint(keyID64)

	var req dto.APIKeyStatusUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取现有API密钥
	apiKey, err := h.dbRepo.GetAPIKeyByID(uint(keyID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "API key not found"})
		return
	}

	// 更新状态
	apiKey.IsActive = req.IsActive

	if err := h.dbRepo.UpdateAPIKey(apiKey); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "API key status updated successfully"})
}

// DeleteAPIKey 删除API密钥
func (h *ManagementHandler) DeleteAPIKey(c *gin.Context) {
	keyIDStr := c.Param("keyID")
	keyID64, err := strconv.ParseInt(keyIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid key ID"})
		return
	}
	keyID := uint(keyID64)

	if err := h.dbRepo.DeleteAPIKey(uint(keyID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "API key deleted successfully"})
}