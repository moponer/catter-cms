<template>
  <div class="scheduled-tasks-container">
    <el-card class="header-card">
      <div class="header-content">
        <h2>定时任务管理</h2>
        <div class="header-actions">
          <el-button type="primary" :icon="Plus" @click="openCreateDialog">
            新建任务
          </el-button>
          <el-button :icon="Refresh" @click="loadTasks">
            刷新
          </el-button>
        </div>
      </div>

      <!-- 调度器状态 -->
      <div class="scheduler-status" v-if="schedulerStatus">
        <el-tag :type="schedulerStatus.is_running ? 'success' : 'danger'" size="large">
          <el-icon><CircleCheck v-if="schedulerStatus.is_running" /><CircleClose v-else /></el-icon>
          调度器状态: {{ schedulerStatus.is_running ? '运行中' : '已停止' }}
        </el-tag>
        <el-button 
          v-if="!schedulerStatus.is_running" 
          type="primary" 
          size="small" 
          @click="restartScheduler"
          :icon="VideoPlay"
        >
          启动调度器
        </el-button>
        <el-button 
          v-else 
          type="warning" 
          size="small" 
          @click="restartScheduler"
          :icon="Refresh"
        >
          重启调度器
        </el-button>
        <span class="status-stats">
          总任务: {{ schedulerStatus.total_tasks }},
          已启用: {{ schedulerStatus.enabled_tasks }},
          运行中: {{ schedulerStatus.running_tasks }},
          空闲: {{ schedulerStatus.idle_tasks }},
          失败: {{ schedulerStatus.failed_tasks }}
        </span>
      </div>
    </el-card>

    <!-- 任务列表 -->
    <el-card class="tasks-card">
      <el-table :data="tasks" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="任务名称" min-width="150" />
        <el-table-column label="任务类型" width="120">
          <template #default="{ row }">
            <el-tag :type="getTaskTypeColor(row.task_type)">
              {{ getTaskTypeName(row.task_type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="cron_expr" label="Cron表达式" width="150">
          <template #default="{ row }">
            <code class="cron-code">{{ row.cron_expr }}</code>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusColor(row.status)">
              {{ getStatusName(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="启用状态" width="80">
          <template #default="{ row }">
            <el-switch
              v-model="row.enabled"
              @change="handleEnableChange(row)"
              :loading="row.switching"
            />
          </template>
        </el-table-column>
        <el-table-column label="执行次数" width="100">
          <template #default="{ row }">
            {{ row.run_count }}
          </template>
        </el-table-column>
        <el-table-column label="下次执行" width="160">
          <template #default="{ row }">
            {{ row.next_run_at ? formatDate(row.next_run_at) : '-' }}
          </template>
        </el-table-column>
        <el-table-column label="最后执行" width="160">
          <template #default="{ row }">
            {{ row.last_run_at ? formatDate(row.last_run_at) : '-' }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="{ row }">
            <el-button
              type="primary"
              size="small"
              @click="handleRunNow(row)"
              :loading="row.running"
              :disabled="!row.enabled"
            >
              立即执行
            </el-button>
            <el-button
              type="warning"
              size="small"
              @click="openEditDialog(row)"
            >
              编辑
            </el-button>
            <el-button
              type="danger"
              size="small"
              @click="handleDelete(row)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 任务编辑/创建对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEditMode ? '编辑定时任务' : '新建定时任务'"
      width="600px"
    >
      <el-form :model="taskForm" :rules="taskRules" ref="taskFormRef" label-width="120px">
        <el-form-item label="任务名称" prop="name">
          <el-input v-model="taskForm.name" placeholder="请输入任务名称" />
        </el-form-item>

        <el-form-item label="任务类型" prop="task_type">
          <el-select v-model="taskForm.task_type" placeholder="请选择任务类型" style="width: 100%">
            <el-option label="定时发布" value="publish" />
            <el-option label="数据备份" value="backup" />
            <el-option label="系统清理" value="cleanup" />
            <el-option label="缓存清理" value="cache_clear" />
            <el-option label="自定义命令" value="custom" />
          </el-select>
        </el-form-item>

        <el-form-item label="Cron表达式" prop="cron_expr">
          <el-input
            v-model="taskForm.cron_expr"
            placeholder="例如: 0 2 * * * (每天凌晨2点)"
          >
            <template #append>
              <el-button @click="showCronHelp = true">帮助</el-button>
            </template>
          </el-input>
          <div class="cron-tips">
            <small>格式: 分 时 日 月 周</small>
          </div>
        </el-form-item>

        <el-form-item
          v-if="taskForm.task_type === 'custom'"
          label="自定义命令"
          prop="command"
        >
          <el-input
            v-model="taskForm.command"
            type="textarea"
            :rows="3"
            placeholder="请输入自定义命令"
          />
        </el-form-item>

        <el-form-item label="启用任务">
          <el-switch v-model="taskForm.enabled" />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSaveTask" :loading="saving">
          保存
        </el-button>
      </template>
    </el-dialog>

    <!-- Cron帮助对话框 -->
    <el-dialog v-model="showCronHelp" title="Cron表达式帮助" width="700px">
      <div class="cron-help">
        <h3>Cron表达式格式</h3>
        <p>Cron表达式由5个字段组成,分别代表:分 时 日 月 周</p>
        
        <h3>字段说明</h3>
        <el-table :data="cronFields" stripe>
          <el-table-column prop="field" label="字段" width="100" />
          <el-table-column prop="range" label="取值范围" width="150" />
          <el-table-column prop="desc" label="说明" />
        </el-table>

        <h3>常用示例</h3>
        <el-table :data="cronExamples" stripe>
          <el-table-column prop="expr" label="表达式" width="200" />
          <el-table-column prop="desc" label="说明" />
        </el-table>

        <h3>特殊符号</h3>
        <ul>
          <li><code>*</code> - 匹配所有值</li>
          <li><code>,</code> - 分隔多个值</li>
          <li><code>-</code> - 指定范围</li>
          <li><code>/</code> - 指定步长</li>
        </ul>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Refresh, CircleCheck, CircleClose, VideoPlay } from '@element-plus/icons-vue'
import axios from 'axios'

// 数据
const tasks = ref([])
const loading = ref(false)
const dialogVisible = ref(false)
const isEditMode = ref(false)
const showCronHelp = ref(false)
const saving = ref(false)
const schedulerStatus = ref(null)

// 表单
const taskFormRef = ref(null)
const taskForm = reactive({
  id: null,
  name: '',
  task_type: '',
  cron_expr: '',
  command: '',
  enabled: true
})

// 表单验证规则
const taskRules = {
  name: [{ required: true, message: '请输入任务名称', trigger: 'blur' }],
  task_type: [{ required: true, message: '请选择任务类型', trigger: 'change' }],
  cron_expr: [
    { required: true, message: '请输入Cron表达式', trigger: 'blur' },
    { pattern: /^(\*|([0-9]|1[0-9]|2[0-9]|3[0-9]|4[0-9]|5[0-9])|\*\/([0-9]+)) (\*|([0-9]|1[0-9]|2[0-3])|\*\/([0-9]+)) (\*|([1-9]|1[0-9]|2[0-9]|3[0-1])|\*\/([0-9]+)) (\*|([1-9]|1[0-2])|\*\/([0-9]+)) (\*|([0-6])|\*\/([0-9]+))$/, message: 'Cron表达式格式不正确', trigger: 'blur' }
  ],
  command: [{ required: true, message: '请输入自定义命令', trigger: 'blur' }]
}

// Cron字段说明
const cronFields = [
  { field: '分', range: '0-59', desc: '分钟' },
  { field: '时', range: '0-23', desc: '小时' },
  { field: '日', range: '1-31', desc: '日期' },
  { field: '月', range: '1-12', desc: '月份' },
  { field: '周', range: '0-6', desc: '星期(0=周日)' }
]

// Cron示例
const cronExamples = [
  { expr: '0 * * * *', desc: '每小时执行一次' },
  { expr: '0 2 * * *', desc: '每天凌晨2点执行' },
  { expr: '0 */6 * * *', desc: '每6小时执行一次' },
  { expr: '0 0 * * 0', desc: '每周日凌晨执行' },
  { expr: '0 0 1 * *', desc: '每月1号执行' },
  { expr: '*/5 * * * *', desc: '每5分钟执行一次' },
  { expr: '0 0,12 * * *', desc: '每天0点和12点执行' }
]

// 加载任务列表
const loadTasks = async () => {
  loading.value = true
  try {
    const response = await axios.get('/api/admin/scheduled-tasks')
    if (response.data.success) {
      tasks.value = response.data.data || []
    }
  } catch (error) {
    ElMessage.error('加载任务列表失败')
  } finally {
    loading.value = false
  }
}

// 加载调度器状态
const loadSchedulerStatus = async () => {
  try {
    const response = await axios.get('/api/admin/scheduler/status')
    if (response.data.success) {
      schedulerStatus.value = response.data.data
    }
  } catch (error) {
    // 错误已在 request 拦截器中处理
  }
}

// 打开创建对话框
const openCreateDialog = () => {
  isEditMode.value = false
  Object.assign(taskForm, {
    id: null,
    name: '',
    task_type: '',
    cron_expr: '',
    command: '',
    enabled: true
  })
  dialogVisible.value = true
}

// 打开编辑对话框
const openEditDialog = (task) => {
  isEditMode.value = true
  Object.assign(taskForm, task)
  dialogVisible.value = true
}

// 保存任务
const handleSaveTask = async () => {
  if (!taskFormRef.value) return
  
  await taskFormRef.value.validate(async (valid) => {
    if (!valid) return
    
    saving.value = true
    try {
      if (isEditMode.value) {
        await axios.put(`/api/admin/scheduled-tasks/${taskForm.id}`, taskForm)
        ElMessage.success('任务更新成功')
      } else {
        await axios.post('/api/admin/scheduled-tasks', taskForm)
        ElMessage.success('任务创建成功')
      }
      dialogVisible.value = false
      await loadTasks()
    } catch (error) {
      ElMessage.error(isEditMode.value ? '任务更新失败' : '任务创建失败')
    } finally {
      saving.value = false
    }
  })
}

// 删除任务
const handleDelete = async (task) => {
  try {
    await ElMessageBox.confirm('确定要删除这个任务吗?', '确认删除', {
      type: 'warning'
    })
    
    await axios.delete(`/api/admin/scheduled-tasks/${task.id}`)
    ElMessage.success('任务删除成功')
    await loadTasks()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('任务删除失败')
    }
  }
}

// 立即执行任务
const handleRunNow = async (task) => {
  try {
    await ElMessageBox.confirm('确定要立即执行这个任务吗?', '确认执行', {
      type: 'warning'
    })
    
    task.running = true
    await axios.post(`/api/admin/scheduled-tasks/${task.id}/run-now`)
    ElMessage.success('任务已开始执行')
    await loadTasks()
    await loadSchedulerStatus()
  } catch (error) {
    task.running = false
    if (error !== 'cancel') {
      ElMessage.error('任务执行失败')
    }
  }
}

// 启用/禁用任务
const handleEnableChange = async (task) => {
  task.switching = true
  try {
    if (task.enabled) {
      await axios.post(`/api/admin/scheduled-tasks/${task.id}/enable`)
      ElMessage.success('任务已启用')
    } else {
      await axios.post(`/api/admin/scheduled-tasks/${task.id}/disable`)
      ElMessage.success('任务已禁用')
    }
    await loadSchedulerStatus()
  } catch (error) {
    task.enabled = !task.enabled
    ElMessage.error(task.enabled ? '任务禁用失败' : '任务启用失败')
  } finally {
    task.switching = false
  }
}

// 重启调度器
const restartScheduler = async () => {
  try {
    await ElMessageBox.confirm('确定要重启调度器吗?', '确认重启', {
      type: 'warning'
    })
    
    await axios.post('/api/admin/scheduler/restart')
    ElMessage.success('调度器重启成功')
    await loadSchedulerStatus()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('调度器重启失败')
    }
  }
}

// 辅助函数
const getTaskTypeName = (type) => {
  const names = {
    publish: '定时发布',
    backup: '数据备份',
    cleanup: '系统清理',
    cache_clear: '缓存清理',
    custom: '自定义命令'
  }
  return names[type] || type
}

const getTaskTypeColor = (type) => {
  const colors = {
    publish: 'primary',
    backup: 'warning',
    cleanup: 'info',
    cache_clear: 'success',
    custom: 'danger'
  }
  return colors[type] || 'info'
}

const getStatusName = (status) => {
  const names = {
    idle: '空闲',
    running: '运行中',
    failed: '失败'
  }
  return names[status] || status
}

const getStatusColor = (status) => {
  const colors = {
    idle: 'success',
    running: 'primary',
    failed: 'danger'
  }
  return colors[status] || 'info'
}

const formatDate = (dateString) => {
  if (!dateString) return '-'
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN')
}

// 生命周期
onMounted(() => {
  loadTasks()
  loadSchedulerStatus()
})
</script>

<style scoped>
.scheduled-tasks-container {
  padding: 20px;
}

.header-card {
  margin-bottom: 20px;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.header-content h2 {
  margin: 0;
  font-size: 24px;
  color: #303133;
}

.header-actions {
  display: flex;
  gap: 10px;
}

.scheduler-status {
  display: flex;
  align-items: center;
  gap: 15px;
  padding: 15px;
  background: #f5f7fa;
  border-radius: 4px;
}

.status-stats {
  font-size: 14px;
  color: #606266;
}

.tasks-card {
  margin-bottom: 20px;
}

.cron-code {
  background: #f5f7fa;
  padding: 2px 6px;
  border-radius: 3px;
  font-size: 12px;
  color: #606266;
}

.cron-tips {
  margin-top: 5px;
  color: #909399;
}

.cron-help {
  padding: 10px 0;
}

.cron-help h3 {
  margin: 20px 0 10px;
  font-size: 16px;
  color: #303133;
}

.cron-help p {
  margin: 5px 0;
  color: #606266;
}

.cron-help ul {
  margin: 10px 0;
  padding-left: 20px;
}

.cron-help li {
  margin: 5px 0;
  color: #606266;
}

.cron-help code {
  background: #f5f7fa;
  padding: 2px 6px;
  border-radius: 3px;
  font-size: 13px;
  color: #e74c3c;
}
</style>
