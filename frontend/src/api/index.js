import axios from 'axios'
import { ElMessage } from 'element-plus'

// 创建一个通用的、携带认证信息的客户端，用于登录后的所有请求
const authApiClient = () => {
  const token = localStorage.getItem('authToken')
  const client = axios.create({
    baseURL: '/api/admin',
    headers: {
      'Content-Type': 'application/json',
      // 如果未来使用JWT，可以在这里添加 Authorization 头
      // 'Authorization': `Bearer ${token}`
    }
  })
  return client
}


// --- 认证 API ---
export const login = (username, password) => {
  const payload = {
    username: username,
    password: password
  };
  // 当axios的post请求第二个参数是对象时，它会自动设置Content-Type为application/json
  return axios.post('/api/admin/login', payload);
}


// --- 配置管理 API ---

export const getAllConfigs = () => {
  return authApiClient().get('/proxy-configs')
}

// 创建通用API配置
export const createProxyConfig = (data) => {
  return authApiClient().post('/proxy-configs', data)
}

// 创建LLM API配置
export const createLlmConfig = (data) => {
  return authApiClient().post('/llm-configs', data)
}

// 获取指定配置的Key列表
export const getKeysForConfig = (configId) => {
  return authApiClient().get(`/proxy-configs/${configId}/keys`);
}

// 为配置添加Key
export const addApiKeyToConfig = (configId, keyData) => {
  return authApiClient().post(`/proxy-configs/${configId}/keys`, keyData);
}

// 更新Key的状态
export const updateApiKeyStatus = (keyId, isActive) => {
  // PATCH请求的body需要是对象
  return authApiClient().patch(`/keys/${keyId}`, { is_active: isActive });
}

// 删除一个Key
export const deleteApiKey = (keyId) => {
  return authApiClient().delete(`/keys/${keyId}`);
}

// 获取应用公共配置
export const getAppConfig = () => {
  return authApiClient().get('/app-config');
}

// 更新通用API配置
export const updateProxyConfig = (id, data) => {
  return authApiClient().put(`/proxy-configs/${id}`, data);
}

// 更新LLM API配置
export const updateLlmConfig = (id, data) => {
  return authApiClient().put(`/llm-configs/${id}`, data);
}

// 更新配置的状态
export const updateConfigStatus = (id, isActive) => {
  return authApiClient().patch(`/proxy-configs/${id}/status`, { is_active: isActive });
}