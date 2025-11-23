import axios from 'axios'
import { ElMessage } from 'element-plus'

// 创建一个通用的API客户端
const apiClient = axios.create({
  baseURL: '/admin',
  headers: {
    'Content-Type': 'application/json',
  }
});

// 添加请求拦截器
apiClient.interceptors.request.use(config => {
  const token = localStorage.getItem('authToken');
  if (token) {
    // 未来使用JWT时，可以在这里添加 Authorization 头
    // config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
}, error => {
  return Promise.reject(error);
});

// 添加响应拦截器，用于统一处理错误
apiClient.interceptors.response.use(
  response => response,
  error => {
    // 在这里可以处理类似401未授权、403禁止访问等全局错误
    if (error.response) {
      // 例如，如果服务器返回了错误消息，则显示它
      ElMessage.error(error.response.data.message || '请求失败');
    } else {
      // 处理网络错误等
      ElMessage.error('网络错误或服务器无响应');
    }
    return Promise.reject(error);
  }
);


// --- 认证 API ---
export const login = (username, password) => {
  const payload = {
    username: username,
    password: password
  };
  // login请求不需要token，但为了统一也使用apiClient
  return apiClient.post('/login', payload);
}


// --- 配置管理 API ---

export const getAllConfigs = () => {
  return apiClient.get('/proxy-configs')
}

// 创建通用API配置
export const createProxyConfig = (data) => {
  return apiClient.post('/proxy-configs', data)
}

// 创建LLM API配置
export const createLlmConfig = (data) => {
  return apiClient.post('/llm-configs', data)
}

// 获取指定配置的Key列表
export const getKeysForConfig = (configId) => {
  return apiClient.get(`/proxy-configs/${configId}/keys`);
}

// 为配置添加Key
export const addApiKeyToConfig = (configId, keyData) => {
  return apiClient.post(`/proxy-configs/${configId}/keys`, keyData);
}

// 更新Key的状态
export const updateApiKeyStatus = (keyId, isActive) => {
  return apiClient.patch(`/keys/${keyId}`, { is_active: isActive });
}

// 删除一个Key
export const deleteApiKey = (keyId) => {
  return apiClient.delete(`/keys/${keyId}`);
}

// 获取应用公共配置
export const getAppConfig = () => {
  return apiClient.get('/app-config');
}

// 更新通用API配置
export const updateProxyConfig = (id, data) => {
  return apiClient.put(`/proxy-configs/${id}`, data);
}

// 更新LLM API配置
export const updateLlmConfig = (id, data) => {
  return apiClient.put(`/llm-configs/${id}`, data);
}

// 更新配置的状态
export const updateConfigStatus = (id, isActive) => {
  return apiClient.put(`/proxy-configs/${id}/status`, { is_active: isActive });
}