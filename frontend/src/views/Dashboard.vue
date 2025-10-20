<template>
  <div>
    <div class="header-bar">
      <h1>{{ $t('dashboard.title') }}</h1>
      <el-button type="primary" @click="openCreateDialog">{{ $t('dashboard.createButton') }}</el-button>
    </div>

    <el-table :data="configs" v-loading="loading" style="width: 100%">
      <el-table-column prop="id" :label="$t('dashboard.table.id')" width="80" />
      <el-table-column prop="name" :label="$t('dashboard.table.name')" />
      <el-table-column prop="slug" :label="$t('dashboard.table.slug')" />
      <el-table-column prop="config_type" :label="$t('dashboard.table.type')">
        <template #default="scope">
          <el-tag :type="scope.row.config_type === 'LLM' ? 'success' : 'primary'">
            {{ scope.row.config_type }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="$t('dashboard.table.targetUrl')" width="400">
         <template #default="scope">
          <span class="target-url">{{ scope.row.config_type === 'LLM' ? scope.row.target_base_url : scope.row.target_url }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="is_active" :label="$t('dashboard.table.status')">
         <template #default="scope">
           <el-popconfirm
             :title="scope.row.is_active ? $t('dashboard.disableConfirm') : $t('dashboard.enableConfirm')"
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
       <el-table-column :label="$t('dashboard.table.actions')" width="280">
         <template #default="scope">
          <div class="action-buttons-horizontal">
            <el-button size="small" type="primary" @click="handleCopy(scope.row)">{{ $t('dashboard.actions.copy') }}</el-button>
            <el-button size="small" @click="handleEdit(scope.row)">{{ $t('dashboard.actions.edit') }}</el-button>
            <el-button size="small" type="primary" @click="openKeyManager(scope.row)">{{ $t('dashboard.actions.manageKeys') }}</el-button>
           </div>
         </template>
       </el-table-column>
    </el-table>

    <!-- 创建/编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="50%" @close="resetForm">
      <el-form ref="configFormRef" :model="configForm" :rules="formRules" label-position="top" label-width="auto">
        <el-form-item :label="$t('dashboard.form.type')" prop="config_type">
          <el-radio-group v-model="configForm.config_type" :disabled="isEditMode">
            <el-radio-button label="GENERIC">{{ $t('dashboard.form.genericApi') }}</el-radio-button>
            <el-radio-button label="LLM">{{ $t('dashboard.form.llmApi') }}</el-radio-button>
          </el-radio-group>
        </el-form-item>
        
        <el-form-item :label="$t('dashboard.form.name')" prop="name">
          <el-input v-model="configForm.name" :placeholder="$t('dashboard.form.namePlaceholder')" />
        </el-form-item>
        <el-form-item :label="$t('dashboard.form.slug')" prop="slug">
          <el-input v-model="configForm.slug" :placeholder="$t('dashboard.form.slugPlaceholder')" />
        </el-form-item>
        
        <!-- 通用API专属字段 -->
        <div v-if="configForm.config_type === 'GENERIC'">
          <el-form-item :label="$t('dashboard.form.method')" prop="method">
            <el-select v-model="configForm.method" :placeholder="$t('dashboard.form.methodPlaceholder')">
              <el-option label="GET" value="GET" />
              <el-option label="POST" value="POST" />
              <el-option label="PUT" value="PUT" />
              <el-option label="DELETE" value="DELETE" />
            </el-select>
          </el-form-item>
          <el-form-item :label="$t('dashboard.form.targetUrl')" prop="target_url">
            <el-input v-model="configForm.target_url" :placeholder="$t('dashboard.form.targetUrlPlaceholder')" />
          </el-form-item>
        </div>

        <!-- LLM API专属字段 -->
        <div v-if="configForm.config_type === 'LLM'">
          <el-form-item :label="$t('dashboard.form.apiFormat')" prop="api_format">
            <el-select v-model="configForm.api_format" :placeholder="$t('dashboard.form.apiFormatPlaceholder')">
              <el-option :label="$t('dashboard.form.openaiCompatible')" value="openai_compatible" />
              <el-option :label="$t('dashboard.form.geminiNative')" value="gemini_native" />
              <el-option :label="$t('dashboard.form.anthropicNative')" value="anthropic_native" />
            </el-select>
          </el-form-item>

          <el-form-item :label="$t('dashboard.form.targetBaseUrl')" prop="target_base_url">
            <el-input v-model="configForm.target_base_url" :placeholder="$t('dashboard.form.targetBaseUrlPlaceholder')" />
          </el-form-item>
        </div>

        <el-form-item :label="$t('dashboard.form.keyLocation')" prop="api_key_location">
           <el-select v-model="configForm.api_key_location" :placeholder="$t('dashboard.form.keyLocationPlaceholder')">
              <el-option label="Header" value="header" />
              <el-option label="Query" value="query" />
            </el-select>
        </el-form-item>
        <el-form-item :label="$t('dashboard.form.keyName')" prop="api_key_name">
          <el-input v-model="configForm.api_key_name" :placeholder="$t('dashboard.form.keyNamePlaceholder')" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">{{ $t('dashboard.form.cancel') }}</el-button>
          <el-button type="primary" @click="submitForm" :loading="formLoading">
            {{ $t('dashboard.form.confirm') }}
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
import { useI18n } from 'vue-i18n'
import {
  getAllConfigs,
  createProxyConfig,
  updateProxyConfig,
  getAppConfig,
  updateConfigStatus
} from '../api'
import { ElMessage } from 'element-plus'
import KeyManager from '../components/KeyManager.vue'

const { t } = useI18n()

const configs = ref([])
const loading = ref(true)

const dialogMode = ref('create') // 'create' or 'edit'
const editingConfigId = ref(null)

const isEditMode = computed(() => dialogMode.value === 'edit')
const dialogTitle = computed(() => isEditMode.value ? t('dashboard.editTitle') : t('dashboard.createTitle'))

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
  config_type: [{ required: true, message: t('dashboard.form.validation.type') }],
  name: [{ required: true, message: t('dashboard.form.validation.name'), trigger: 'blur' }],
  slug: [{ required: true, message: t('dashboard.form.validation.slug'), trigger: 'blur' }],
  method: [{ required: configForm.config_type === 'GENERIC', message: t('dashboard.form.validation.method') }],
  target_url: [{ required: configForm.config_type === 'GENERIC', message: t('dashboard.form.validation.targetUrl'), trigger: 'blur' }],
  target_base_url: [{ required: configForm.config_type === 'LLM', message: t('dashboard.form.validation.targetBaseUrl'), trigger: 'blur' }],
  api_key_location: [{ required: true, message: t('dashboard.form.validation.keyLocation') }],
  api_key_name: [{ required: true, message: t('dashboard.form.validation.keyName'), trigger: 'blur' }],
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
    ElMessage.error(t('dashboard.messages.loadFailed'))
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
          ElMessage.success(t('dashboard.messages.updateSuccess'));
        } else {
          // --- 创建逻辑 ---
          await createProxyConfig(payload);
          ElMessage.success(t('dashboard.messages.createSuccess'));
        }

        dialogVisible.value = false;
        fetchConfigs(); // 重新加载列表
      } catch (error) {
        ElMessage.error(isEditMode.value ? t('dashboard.messages.updateFailed') : t('dashboard.messages.createFailed'));
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
    ElMessage.error(t('dashboard.messages.appConfigFailed'));
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
    
    ElMessage.success(t('dashboard.messages.statusUpdateSuccess'));
  } catch (error) {
    ElMessage.error(t('dashboard.messages.statusUpdateFailed'));
    console.error(error);
    // 如果API调用失败，前端UI不会改变，保持原状
  }
}

const handleCopy = (row) => {
  if (!appConfig.value.proxy_public_base_url) {
    ElMessage.error(t('dashboard.messages.copyBaseUrlMissing'));
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
        message: t('dashboard.messages.copySuccess', { url: proxyUrl }),
        duration: 2000 // 消息显示2秒
      });
    }).catch(err => {
      ElMessage.error(t('dashboard.messages.copyFailed'));
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
          message: t('dashboard.messages.copySuccess', { url: proxyUrl }),
          duration: 2000 // 消息显示2秒
        });
      } else {
        ElMessage.error(t('dashboard.messages.copyFailed'));
      }
    } catch (err) {
      ElMessage.error(t('dashboard.messages.copyFailed'));
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