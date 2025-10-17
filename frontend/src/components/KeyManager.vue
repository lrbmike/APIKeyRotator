<template>
  <el-dialog v-model="visible" :title="t('keyManager.title', { name: configName })" width="60%" @close="handleClose">
    
    <!-- 添加新Key的表单 -->
    <el-form :inline="true" :model="newKeyForm" class="add-key-form">
      <el-form-item :label="t('keyManager.newKey')">
        <el-input v-model="newKeyForm.key_value" :placeholder="t('keyManager.placeholder')" style="width: 400px;"/>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="handleAddNewKey" :loading="addLoading">{{ t('keyManager.add') }}</el-button>
      </el-form-item>
    </el-form>

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
import { getKeysForConfig, addApiKeyToConfig, updateApiKeyStatus, deleteApiKey } from '../api'
import { ElMessage } from 'element-plus'
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

const newKeyForm = reactive({
  key_value: '',
  is_active: true
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

const handleClose = () => {
  emit('update:modelValue', false)
}
</script>

<style scoped>
.add-key-form {
  margin-bottom: 20px;
}
</style>