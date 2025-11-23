<template>
  <el-dialog v-model="visible" :title="t('keyManager.title', { name: configName })" width="60%" @close="handleClose">
    
    <!-- 添加新Key的表单 -->
    <el-form :inline="true" :model="newKeyForm" class="add-key-form">
      <el-form-item :label="t('keyManager.newKey')">
        <el-input v-model="newKeyForm.key_value" :placeholder="t('keyManager.placeholder')" style="width: 400px;"/>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="handleAddNewKey" :loading="addLoading">{{ t('keyManager.add') }}</el-button>
        <el-button type="success" @click="showBatchImport = true">{{ t('keyManager.batchImport') }}</el-button>
        <el-button type="info" @click="handleExportKeys" :disabled="keys.length === 0">
          {{ t('keyManager.export') }}
        </el-button>
        <el-button type="danger" @click="handleClearAllKeys" :loading="clearAllLoading" :disabled="keys.length === 0">
          {{ t('keyManager.clearAll') }}
        </el-button>
      </el-form-item>
    </el-form>

    <!-- 批量导入对话框 -->
    <el-dialog v-model="showBatchImport" :title="t('keyManager.batchImportTitle')" width="50%" append-to-body>
      <el-form :model="batchImportForm" label-width="100px">
        <el-form-item :label="t('keyManager.batchImportLabel')">
          <el-input
            v-model="batchImportForm.keys"
            type="textarea"
            :rows="10"
            :placeholder="t('keyManager.batchImportPlaceholder')"
          />
          <div class="batch-import-tip">
            <el-text type="info" size="small">{{ t('keyManager.batchImportTip') }}</el-text>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showBatchImport = false">{{ t('common.cancel') }}</el-button>
          <el-button type="primary" @click="handleBatchImport" :loading="batchImportLoading">
            {{ t('keyManager.batchImportConfirm') }}
          </el-button>
        </span>
      </template>
    </el-dialog>

    <!-- Key列表 -->
    <el-table :data="keys" v-loading="loading">
      <el-table-column prop="id" :label="t('keyManager.table.id')" width="80" />
      <el-table-column :label="t('keyManager.table.key')">
        <template #default="scope">
          <span>{{ maskApiKey(scope.row.key_value) }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="t('keyManager.table.status')" width="120">
        <template #default="scope">
          <el-switch
            :model-value="scope.row.is_active"
            @change="handleStatusChange(scope.row)"
          />
        </template>
      </el-table-column>
      <el-table-column :label="t('keyManager.table.actions')" width="120">
        <template #default="scope">
          <el-popconfirm :title="t('keyManager.deleteConfirm')" @confirm="handleDeleteKey(scope.row.id)">
            <template #reference>
              <el-button size="small" type="danger">{{ t('keyManager.delete') }}</el-button>
            </template>
          </el-popconfirm>
        </template>
      </el-table-column>
    </el-table>

  </el-dialog>
</template>

<script setup>
import { ref, reactive, watch } from 'vue'
import { getKeysForConfig, addApiKeyToConfig, updateApiKeyStatus, deleteApiKey, batchImportApiKeys, clearAllApiKeys } from '../api'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const props = defineProps({
  modelValue: Boolean, // 控制对话框显示
  configId: Number,
  configName: String
})

const emit = defineEmits(['update:modelValue'])

const visible = ref(props.modelValue)
const keys = ref([])
const loading = ref(false)
const addLoading = ref(false)
const batchImportLoading = ref(false)
const clearAllLoading = ref(false)
const showBatchImport = ref(false)

const newKeyForm = reactive({
  key_value: '',
  is_active: true
})

const batchImportForm = reactive({
  keys: ''
})

// 脱敏函数
const maskApiKey = (key) => {
  if (!key || key.length < 10) return '*****'
  return `${key.substring(0, 6)}*****${key.substring(key.length - 4)}`
}

const fetchKeys = async () => {
  if (!props.configId || props.configId <= 0) {
    keys.value = []; // 清空旧数据
    return;
  }
  loading.value = true
  try {
    const response = await getKeysForConfig(props.configId)
    keys.value = response.data
  } catch (error) {
    ElMessage.error(t('keyManager.messages.loadFailed'))
  } finally {
    loading.value = false
  }
}

const handleAddNewKey = async () => {
  if (!newKeyForm.key_value.trim()) {
    ElMessage.warning(t('keyManager.messages.keyRequired'));
    return;
  }
  addLoading.value = true
  try {
    await addApiKeyToConfig(props.configId, newKeyForm);
    ElMessage.success(t('keyManager.messages.addSuccess'));
    newKeyForm.key_value = ''; // 清空输入框
    await fetchKeys(); // 重新加载列表
  } catch (error) {
    ElMessage.error(t('keyManager.messages.addFailed'));
  } finally {
    addLoading.value = false
  }
}

const handleStatusChange = async (row) => {
  // 乐观更新UI
  const originalStatus = row.is_active
  row.is_active = !originalStatus
  try {
    await updateApiKeyStatus(row.id, row.is_active)
    ElMessage.success(t('keyManager.messages.statusUpdateSuccess'))
  } catch (error) {
    // 回滚UI
    row.is_active = originalStatus
    ElMessage.error(t('keyManager.messages.statusUpdateFailed'))
  }
}

const handleDeleteKey = async (keyId) => {
  try {
    await deleteApiKey(keyId);
    ElMessage.success(t('keyManager.messages.deleteSuccess'));
    await fetchKeys(); // 重新加载列表
  } catch (error) {
    ElMessage.error(t('keyManager.messages.deleteFailed'));
  }
}

// 监听v-model的变化来控制对话框
watch(() => props.modelValue, (newValue) => {
  visible.value = newValue
})

watch(() => props.configId, (newId) => {
  if (newId && props.modelValue) { // 只有当对话框可见且ID有效时才加载
    fetchKeys()
  }
}, { immediate: true })

const handleBatchImport = async () => {
  if (!batchImportForm.keys.trim()) {
    ElMessage.warning(t('keyManager.messages.keysRequired'));
    return;
  }

  // 解析输入的keys，一行一个
  const keyList = batchImportForm.keys
    .split('\n')
    .map(key => key.trim())
    .filter(key => key.length > 0);

  if (keyList.length === 0) {
    ElMessage.warning(t('keyManager.messages.noValidKeys'));
    return;
  }

  batchImportLoading.value = true;
  try {
    const response = await batchImportApiKeys(props.configId, { keys: keyList });
    const { success_count, failed_count, failed_keys } = response.data;

    if (failed_count === 0) {
      ElMessage.success(t('keyManager.messages.batchImportSuccess', { count: success_count }));
    } else if (success_count === 0) {
      ElMessage.error(t('keyManager.messages.batchImportFailed', { count: failed_count }));
    } else {
      ElMessage.warning(t('keyManager.messages.batchImportPartial', {
        success: success_count,
        failed: failed_count
      }));
    }

    // 清空表单并关闭对话框
    batchImportForm.keys = '';
    showBatchImport.value = false;

    // 重新加载key列表
    await fetchKeys();
  } catch (error) {
    ElMessage.error(t('keyManager.messages.batchImportError'));
  } finally {
    batchImportLoading.value = false;
  }
}

// 处理清除所有Keys
const handleClearAllKeys = async () => {
  if (keys.value.length === 0) {
    ElMessage.warning(t('keyManager.messages.noKeysToClear'));
    return;
  }

  try {
    await ElMessageBox.confirm(
      t('keyManager.clearAllConfirm', { count: keys.value.length }),
      t('keyManager.clearAllTitle'),
      {
        confirmButtonText: t('keyManager.clearAllConfirmBtn'),
        cancelButtonText: t('common.cancel'),
        type: 'warning',
        confirmButtonClass: 'el-button--danger'
      }
    );

    await confirmClearAllKeys();
  } catch (action) {
    if (action !== 'cancel') {
      console.error('Clear all keys dialog error:', action);
    }
  }
};

// 确认清除所有Keys
const confirmClearAllKeys = async () => {
  clearAllLoading.value = true;
  try {
    const response = await clearAllApiKeys(props.configId);
    const { deleted_count } = response.data;

    ElMessage.success(t('keyManager.messages.clearAllSuccess', { count: deleted_count }));
    await fetchKeys(); // 重新加载列表
  } catch (error) {
    ElMessage.error(t('keyManager.messages.clearAllFailed'));
  } finally {
    clearAllLoading.value = false;
  }
};

// 处理导出Keys
const handleExportKeys = () => {
  if (keys.value.length === 0) {
    ElMessage.warning(t('keyManager.messages.noKeysToExport'));
    return;
  }

  // 过滤出激活状态的keys，并按行拼接
  const activeKeys = keys.value
    .filter(key => key.is_active)
    .map(key => key.key_value)
    .join('\n');

  if (!activeKeys.trim()) {
    ElMessage.warning(t('keyManager.messages.noActiveKeysToExport'));
    return;
  }

  // 创建Blob对象
  const blob = new Blob([activeKeys], { type: 'text/plain;charset=utf-8' });

  // 创建下载链接
  const link = document.createElement('a');
  const url = URL.createObjectURL(blob);
  link.setAttribute('href', url);

  // 设置文件名，包含配置名和当前时间戳
  const timestamp = new Date().toISOString().slice(0, 19).replace(/[:-]/g, '');
  const filename = `${props.configName}_keys_${timestamp}.txt`;
  link.setAttribute('download', filename);

  // 触发下载
  link.style.visibility = 'hidden';
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);

  // 清理URL对象
  URL.revokeObjectURL(url);

  ElMessage.success(t('keyManager.messages.exportSuccess', { count: keys.value.filter(key => key.is_active).length }));
};

const handleClose = () => {
  emit('update:modelValue', false)
}
</script>

<style scoped>
.add-key-form {
  margin-bottom: 20px;
}

.batch-import-tip {
  margin-top: 8px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style>