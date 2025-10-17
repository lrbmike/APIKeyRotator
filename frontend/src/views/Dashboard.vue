<template>
  <div>
    <div class="header-bar">
      <h1>代理服务配置</h1>
      <el-button type="primary" @click="openCreateDialog">创建新服务</el-button>
    </div>

    <el-table :data="configs" v-loading="loading" style="width: 100%">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="name" label="服务名称" />
      <el-table-column prop="slug" label="服务标识 (Slug)" />
      <el-table-column prop="config_type" label="类型">
        <template #default="scope">
          <el-tag :type="scope.row.config_type === 'LLM' ? 'success' : 'primary'">
            {{ scope.row.config_type }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="目标地址" width="400">
         <template #default="scope">
          <span class="target-url">{{ scope.row.config_type === 'LLM' ? scope.row.target_base_url : scope.row.target_url }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="is_active" label="状态">
         <template #default="scope">
           <el-popconfirm
             :title="`确定要'${scope.row.is_active ? '禁用' : '启用'}'这个服务吗?`"
             width="220"
             @confirm="handleStatusChange(scope.row)"
           >
             <template #reference>
                <!-- 
                  我们在这里加一个div来阻止switch的默认点击行为,
                  因为点击事件将由Popconfirm来触发。
                -->
               <div @click.stop.prevent>
                 <el-switch :model-value="scope.row.is_active" />
               </div>
             </template>
           </el-popconfirm>
         </template>
       </el-table-column>
       <el-table-column label="操作" width="280">
         <template #default="scope">
          <div class="action-buttons-horizontal">
            <el-button size="small" type="primary" @click="handleCopy(scope.row)">复制地址</el-button>
            <el-button size="small" @click="handleEdit(scope.row)">编辑</el-button>
            <el-button size="small" type="primary" @click="openKeyManager(scope.row)">管理Key</el-button>
           </div>
         </template>
       </el-table-column>
    </el-table>

    <!-- 创建/编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="50%" @close="resetForm">
      <el-form ref="configFormRef" :model="configForm" :rules="formRules" label-width="120px">
        <el-form-item label="服务类型" prop="config_type">
          <el-radio-group v-model="configForm.config_type" :disabled="isEditMode">
            <el-radio-button label="GENERIC">通用API</el-radio-button>
            <el-radio-button label="LLM">LLM API</el-radio-button>
          </el-radio-group>
        </el-form-item>
        
        <el-form-item label="服务名称" prop="name">
          <el-input v-model="configForm.name" placeholder="例如：官方OpenAI生产环境" />
        </el-form-item>
        <el-form-item label="服务标识 (Slug)" prop="slug">
          <el-input v-model="configForm.slug" placeholder="例如：openai-prod (全局唯一, 只能是小写字母和-)" />
        </el-form-item>
        
        <!-- 通用API专属字段 -->
        <div v-if="configForm.config_type === 'GENERIC'">
          <el-form-item label="请求方法" prop="method">
            <el-select v-model="configForm.method" placeholder="选择请求方法">
              <el-option label="GET" value="GET" />
              <el-option label="POST" value="POST" />
              <el-option label="PUT" value="PUT" />
              <el-option label="DELETE" value="DELETE" />
            </el-select>
          </el-form-item>
          <el-form-item label="目标URL" prop="target_url">
            <el-input v-model="configForm.target_url" placeholder="https://api.weatherstack.com/current" />
          </el-form-item>
        </div>

        <!-- LLM API专属字段 -->
        <div v-if="configForm.config_type === 'LLM'">
          <el-form-item label="API 格式" prop="api_format">
            <el-select v-model="configForm.api_format" placeholder="选择API格式">
              <el-option label="OpenAI Compatible" value="openai_compatible" />
              <el-option label="Gemini Native" value="gemini_native" />
            </el-select>
          </el-form-item>

          <el-form-item label="目标Base URL" prop="target_base_url">
            <el-input v-model="configForm.target_base_url" placeholder="https://api.openai.com/v1" />
          </el-form-item>
        </div>

        <el-form-item label="密钥位置" prop="api_key_location">
           <el-select v-model="configForm.api_key_location" placeholder="选择密钥注入位置">
              <el-option label="Header" value="header" />
              <el-option label="Query" value="query" />
            </el-select>
        </el-form-item>
        <el-form-item label="密钥名称" prop="api_key_name">
          <el-input v-model="configForm.api_key_name" placeholder="例如：Authorization 或 access_key" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitForm" :loading="formLoading">
            确认
          </el-button>
        </span>
      </template>
    </el-dialog>

    <KeyManager 
      v-if="keyManagerVisible"
      v-model="keyManagerVisible" 
      :config-id="selectedConfig.id"
      :config-name="selectedConfig.name"
    />

  </div>
</template>

<script setup>
import { ref, onMounted, reactive, computed } from 'vue'
import { 
  getAllConfigs, 
  createProxyConfig,
  updateProxyConfig,
  getAppConfig,
  updateConfigStatus
} from '../api' 
import { ElMessage } from 'element-plus'
import KeyManager from '../components/KeyManager.vue'

const configs = ref([])
const loading = ref(true)

const dialogMode = ref('create') // 'create' or 'edit'
const editingConfigId = ref(null)

const isEditMode = computed(() => dialogMode.value === 'edit')
const dialogTitle = computed(() => isEditMode.value ? '编辑代理服务' : '创建代理服务')

// 对话框相关
const dialogVisible = ref(false)
const formLoading = ref(false)
const configFormRef = ref(null)

// 表单数据
const initialFormState = {
  config_type: 'GENERIC',
  name: '',
  slug: '',
  method: 'GET',
  api_format: 'openai_compatible',
  target_url: '',
  target_base_url: '',
  api_key_location: 'header',
  api_key_name: ''
}
const configForm = reactive({ ...initialFormState })

// 表单校验规则
const formRules = computed(() => ({
  config_type: [{ required: true, message: '请选择服务类型' }],
  name: [{ required: true, message: '请输入服务名称', trigger: 'blur' }],
  slug: [{ required: true, message: '请输入服务标识', trigger: 'blur' }],
  method: [{ required: configForm.config_type === 'GENERIC', message: '请选择请求方法' }],
  target_url: [{ required: configForm.config_type === 'GENERIC', message: '请输入目标URL', trigger: 'blur' }],
  target_base_url: [{ required: configForm.config_type === 'LLM', message: '请输入目标Base URL', trigger: 'blur' }],
  api_key_location: [{ required: true, message: '请选择密钥位置' }],
  api_key_name: [{ required: true, message: '请输入密钥名称', trigger: 'blur' }],
}))

const keyManagerVisible = ref(false)
const selectedConfig = ref(null)

const openKeyManager = (row) => {
  selectedConfig.value = row
  keyManagerVisible.value = true
}

// 用于存储从后端获取的公共配置
const appConfig = ref({
  proxy_public_base_url: ''
})


// --- 方法 ---

const fetchConfigs = async () => {
  try {
    loading.value = true
    const response = await getAllConfigs()
    configs.value = response.data
  } catch (error) {
    ElMessage.error('加载配置列表失败')
  } finally {
    loading.value = false
  }
}

const openCreateDialog = () => {
  resetForm(); // 确保是空表单
  dialogVisible.value = true;
}

const resetForm = () => {
  dialogMode.value = 'create';
  editingConfigId.value = null;
  Object.assign(configForm, initialFormState);
  if (configFormRef.value) {
    configFormRef.value.clearValidate();
  }
}

const submitForm = async () => {
  await configFormRef.value.validate(async (valid) => {
    if (valid) {
      formLoading.value = true;
      try {
        const payload = {
          name: configForm.name,
          slug: configForm.slug,
          config_type: configForm.config_type,
          api_key_location: configForm.api_key_location,
          api_key_name: configForm.api_key_name,
          method: configForm.config_type === 'GENERIC' ? configForm.method : null,
          target_url: configForm.config_type === 'GENERIC' ? configForm.target_url : null,
          api_format: configForm.config_type === 'LLM' ? configForm.api_format : null,
          target_base_url: configForm.config_type === 'LLM' ? configForm.target_base_url : null,
        };

        if (isEditMode.value) {
          // --- 编辑逻辑 ---
          await updateProxyConfig(editingConfigId.value, payload);
          ElMessage.success('更新成功！');
        } else {
          // --- 创建逻辑 ---
          await createProxyConfig(payload);
          ElMessage.success('创建成功！');
        }

        dialogVisible.value = false;
        fetchConfigs(); // 重新加载列表
      } catch (error) {
        ElMessage.error(`${isEditMode.value ? '更新' : '创建'}失败，请检查输入或查看控制台`);
        console.error(error);
      } finally {
        formLoading.value = false;
      }
    }
  });
}

const handleEdit = (row) => {
  dialogMode.value = 'edit';
  editingConfigId.value = row.id;
  
  // 将行数据填充到表单中
  // 使用 Object.assign 确保只填充表单中存在的字段
  Object.assign(configForm, row);

  dialogVisible.value = true;
}

const fetchAppConfig = async () => {
  try {
    const response = await getAppConfig();
    appConfig.value = response.data;
  } catch (error) {
    ElMessage.error('获取应用配置失败');
    console.error(error);
  }
}

// 处理状态改变的方法
const handleStatusChange = async (row) => {
  const originalStatus = row.is_active;
  const newStatus = !originalStatus;

  try {
    // 调用API更新后端状态
    await updateConfigStatus(row.id, newStatus);
    
    // API调用成功后，才更新前端UI
    row.is_active = newStatus;
    
    ElMessage.success('状态更新成功！');
  } catch (error) {
    ElMessage.error('状态更新失败');
    console.error(error);
    // 如果API调用失败，前端UI不会改变，保持原状
  }
}

const handleCopy = (row) => {
  if (!appConfig.value.proxy_public_base_url) {
    ElMessage.error('代理域名未配置，无法复制');
    return;
  }

  let proxyUrl = '';
  const baseUrl = appConfig.value.proxy_public_base_url.replace(/\/$/, ''); // 移除末尾的斜杠

  if (row.config_type === 'GENERIC') {
    // 示例: http://localhost:8000/proxy/weather
    proxyUrl = `${baseUrl}/proxy/${row.slug}`;
  } else if (row.config_type === 'LLM') {
    // 示例: http://localhost:8000/llm/openai-prod
    // 注意：对于LLM，我们只复制base_url，因为具体的action由SDK调用时决定
    proxyUrl = `${baseUrl}/llm/${row.slug}`;
  }

  // 检查 Clipboard API 是否可用
  if (navigator.clipboard && window.isSecureContext) {
    // 使用 Clipboard API 进行复制，这是现代浏览器的标准做法
    navigator.clipboard.writeText(proxyUrl).then(() => {
      ElMessage.success({
        message: `已复制: ${proxyUrl}`,
        duration: 2000 // 消息显示2秒
      });
    }).catch(err => {
      ElMessage.error('复制失败，请检查浏览器权限或手动复制');
      console.error('Could not copy text: ', err);
    });
  } else {
    // 降级方案：使用传统的 document.execCommand('copy')
    try {
      const textArea = document.createElement("textarea");
      textArea.value = proxyUrl;
      
      // 避免滚动到底部
      textArea.style.top = "0";
      textArea.style.left = "0";
      textArea.style.position = "fixed";
      textArea.style.opacity = "0";
      
      document.body.appendChild(textArea);
      textArea.focus();
      textArea.select();
      
      const successful = document.execCommand('copy');
      document.body.removeChild(textArea);
      
      if (successful) {
        ElMessage.success({
          message: `已复制: ${proxyUrl}`,
          duration: 2000 // 消息显示2秒
        });
      } else {
        ElMessage.error('复制失败，请手动复制');
      }
    } catch (err) {
      ElMessage.error('复制失败，请手动复制');
      console.error('Could not copy text: ', err);
    }
  }
}

onMounted(() => {
  fetchConfigs()
  fetchAppConfig() 
})
</script>

<style scoped>
.header-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}
.target-url {
  word-break: break-all;
}
.action-buttons-horizontal {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}
.action-buttons-horizontal .el-button {
  margin-left: 0 !important;
}
</style>