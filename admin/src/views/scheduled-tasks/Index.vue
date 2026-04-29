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
        <el-tag :type="schedulerStatus.is_running ? 'success' : 'danger'" size="large" class="status-tag">
          <span class="status-tag-content">
            <el-icon><CircleCheck v-if="schedulerStatus.is_running" /><CircleClose v-else /></el-icon>
            <span>调度器状态: {{ schedulerStatus.is_running ? '运行中' : '已停止' }}</span>
          </span>
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
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column prop="name" label="任务名称" min-width="140" />
        <el-table-column label="任务类型" width="110">
          <template #default="{ row }">
            <el-tag :type="getTaskTypeColor(row.task_type)" size="small">
              {{ getTaskTypeName(row.task_type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="cron_expr" label="Cron" width="130">
          <template #default="{ row }">
            <code class="cron-code">{{ row.cron_expr }}</code>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="getStatusColor(row.status)" size="small">
              {{ getStatusName(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="启用" width="70">
          <template #default="{ row }">
            <el-switch
              v-model="row.enabled"
              @change="handleEnableChange(row)"
              :loading="row.switching"
              size="small"
            />
          </template>
        </el-table-column>
        <el-table-column label="执行" width="60">
          <template #default="{ row }">
            <el-button link type="primary" @click="openExecutionLog(row)" :disabled="row.run_count === 0">
              {{ row.run_count }}
            </el-button>
          </template>
        </el-table-column>
        <el-table-column label="下次执行" width="150">
          <template #default="{ row }">
            <span class="next-run">{{ row.next_run_at ? formatDate(row.next_run_at) : '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="160" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleRunNow(row)" :loading="row.running" :disabled="!row.enabled">
              执行
            </el-button>
            <el-button type="info" link size="small" @click="openExecutionLog(row)" :disabled="row.run_count === 0">
              日志
            </el-button>
            <el-dropdown trigger="click" @command="(cmd) => handleAction(cmd, row)">
              <el-button type="primary" link size="small">更多<el-icon class="dropdown-arrow"><ArrowDown /></el-icon></el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="edit">编辑</el-dropdown-item>
                  <el-dropdown-item command="delete" style="color: var(--el-color-danger)">删除</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 任务编辑/创建对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEditMode ? '编辑定时任务' : '新建定时任务'"
      width="700px"
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

        <!-- 内置任务快捷模板（仅自定义命令类型显示） -->
        <el-form-item 
          v-if="taskForm.task_type === 'custom'"
          label="快捷模板"
        >
          <div class="template-cards">
            <div v-for="(group, groupName) in groupedTemplates" :key="groupName" class="template-group">
              <div class="template-group-title">{{ groupName }}</div>
              <div class="template-group-items">
                <div
                  v-for="tpl in group"
                  :key="tpl.command"
                  class="template-card"
                  :class="{ active: taskForm.command === tpl.command }"
                  @click="selectTemplate(tpl)"
                >
                  <div class="template-card-name">{{ tpl.name }}</div>
                  <div class="template-card-desc">{{ tpl.desc }}</div>
                </div>
              </div>
            </div>
          </div>
        </el-form-item>

        <el-form-item label="执行频率" prop="cron_expr">
          <!-- 人性化频率选择 -->
          <div class="cron-quick-select">
            <div class="freq-group">
              <span class="freq-label">快速选择：</span>
              <el-radio-group v-model="selectedFreq" @change="onFreqChange">
                <el-radio-button value="every_minute">每分钟</el-radio-button>
                <el-radio-button value="every_5_min">每5分钟</el-radio-button>
                <el-radio-button value="every_30_min">每30分钟</el-radio-button>
                <el-radio-button value="hourly">每小时</el-radio-button>
                <el-radio-button value="daily">每天</el-radio-button>
                <el-radio-button value="weekly">每周</el-radio-button>
                <el-radio-button value="monthly">每月</el-radio-button>
              </el-radio-group>
            </div>

            <!-- 每天/每周/每月 的具体时间设置 -->
            <div v-if="showTimeSetting" class="time-setting">
              <span class="setting-label">{{ timeSettingLabel }}：</span>
              <el-time-picker
                v-if="selectedFreq === 'daily'"
                v-model="dailyTime"
                format="HH:mm"
                placeholder="选择时间"
                @change="onTimeChange"
                style="width: 150px;"
              />
              <template v-else-if="selectedFreq === 'weekly'">
                <el-select v-model="weeklyDay" @change="onTimeChange" style="width: 120px; margin-right: 8px;">
                  <el-option label="周一" value="1" />
                  <el-option label="周二" value="2" />
                  <el-option label="周三" value="3" />
                  <el-option label="周四" value="4" />
                  <el-option label="周五" value="5" />
                  <el-option label="周六" value="6" />
                  <el-option label="周日" value="0" />
                </el-select>
                <el-time-picker
                  v-model="weeklyTime"
                  format="HH:mm"
                  placeholder="选择时间"
                  @change="onTimeChange"
                  style="width: 150px;"
                />
              </template>
              <template v-else-if="selectedFreq === 'monthly'">
                <el-select v-model="monthlyDay" @change="onTimeChange" style="width: 100px; margin-right: 8px;">
                  <el-option v-for="d in 28" :key="d" :label="d + '号'" :value="d" />
                </el-select>
                <el-time-picker
                  v-model="monthlyTime"
                  format="HH:mm"
                  placeholder="选择时间"
                  @change="onTimeChange"
                  style="width: 150px;"
                />
              </template>
            </div>
          </div>

          <!-- Cron 表达式显示（只读） -->
          <div class="cron-preview">
            <el-input v-model="taskForm.cron_expr" placeholder="Cron 表达式" :readonly="true">
              <template #prepend>Cron</template>
            </el-input>
            <el-button @click="showCronHelp = true" style="margin-left: 8px;">帮助</el-button>
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
          <div class="command-help">
            <small>支持系统命令、SQLite 命令、以及 forfiles 等文件管理命令。示例：sqlite3 data/catter.db "SELECT COUNT(*) FROM posts;"</small>
          </div>
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
        <p>Cron表达式由5个字段组成，分别代表：分 时 日 月 周</p>
        
        <h3>字段说明</h3>
        <el-table :data="cronFields" stripe size="small">
          <el-table-column prop="field" label="字段" width="80" />
          <el-table-column prop="range" label="取值范围" width="120" />
          <el-table-column prop="desc" label="说明" />
        </el-table>

        <h3>常用示例</h3>
        <el-table :data="cronExamples" stripe size="small">
          <el-table-column prop="expr" label="表达式" width="160" />
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

    <!-- 执行日志对话框 -->
    <el-dialog v-model="showExecutionLog" title="执行日志" width="800px">
      <div v-if="currentTask" class="log-header">
        <strong>{{ currentTask.name }}</strong>
        <span class="log-meta">执行次数: {{ currentTask.run_count }}</span>
      </div>
      <el-empty v-if="executions.length === 0" description="暂无执行记录" />
      <div v-else class="execution-list">
        <div v-for="exec in executions" :key="exec.id" class="execution-item">
          <div class="execution-header">
            <span class="execution-time">{{ formatDate(exec.started_at) }}</span>
            <el-tag :type="exec.status === 'success' ? 'success' : 'danger'" size="small">
              {{ exec.status === 'success' ? '成功' : '失败' }}
            </el-tag>
            <span class="execution-duration">{{ exec.duration ? exec.duration + 's' : '-' }}</span>
          </div>
          <pre class="execution-output">{{ exec.message || '(无输出)' }}</pre>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Refresh, CircleCheck, CircleClose, VideoPlay, ArrowDown } from '@element-plus/icons-vue'
import request from '@/utils/request'

// 数据
const tasks = ref([])
const loading = ref(false)
const dialogVisible = ref(false)
const isEditMode = ref(false)
const showCronHelp = ref(false)
const showExecutionLog = ref(false)
const saving = ref(false)
const schedulerStatus = ref(null)
const currentTask = ref(null)
const executions = ref([])

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

// 内置任务模板（CMS 实用场景）
const builtinTemplates = [
  // 数据备份类
  {
    name: '数据库备份',
    command: 'sqlite3 data/catter.db ".backup backup/backup_$(date +%%Y%%m%%d).db"',
    desc: '每日备份 SQLite 数据库',
    cron_default: '0 3 * * *',
    icon: 'database'
  },
  {
    name: '清理旧备份',
    command: 'forfiles /p backup /s /d -30 /c "cmd /c del @path"',
    desc: '删除30天前的备份文件',
    cron_default: '0 4 * * *',
    icon: 'delete'
  },
  // 文件清理类
  {
    name: '清理临时文件',
    command: 'forfiles /p public\\uploads\\temp /s /d -7 /c "cmd /c del @path"',
    desc: '删除7天前的临时上传文件',
    cron_default: '0 2 * * *',
    icon: 'folder-delete'
  },
  {
    name: '清理日志文件',
    command: 'forfiles /p logs /s /d -14 /c "cmd /c del @path"',
    desc: '删除14天前的日志文件',
    cron_default: '0 5 * * 0',
    icon: 'document-delete'
  },
  // 数据库维护类
  {
    name: '清理回收站文章',
    command: 'sqlite3 data/catter.db "DELETE FROM posts WHERE deleted_at IS NOT NULL AND datetime(deleted_at) < datetime(\'now\', \'-30 days\')"',
    desc: '永久删除30天前已删除的文章',
    cron_default: '0 1 * * *',
    icon: 'delete'
  },
  {
    name: '清理过期Token',
    command: 'sqlite3 data/catter.db "DELETE FROM user_tokens WHERE expires_at < datetime(\'now\')"',
    desc: '删除过期的用户认证Token',
    cron_default: '0 0 * * *',
    icon: 'key'
  },
  // 监控统计类
  {
    name: '磁盘空间检查',
    command: 'wmic logicaldisk get size,freespace,caption',
    desc: '检查磁盘使用情况',
    cron_default: '0 */6 * * *',
    icon: 'monitor'
  },
  {
    name: 'CMS数据统计',
    command: 'sqlite3 -header -column data/catter.db "SELECT \'文章总数\' as 项目, COUNT(*) as 数量 FROM posts UNION ALL SELECT \'已发布\', COUNT(*) FROM posts WHERE status=\'published\' UNION ALL SELECT \'分类\', COUNT(*) FROM categories UNION ALL SELECT \'用户\', COUNT(*) FROM users;"',
    desc: '生成 CMS 数据统计报告',
    cron_default: '0 9 * * *',
    icon: 'data-analysis'
  },
  // 系统维护类
  {
    name: '优化数据库',
    command: 'sqlite3 data/catter.db "VACUUM;"',
    desc: '执行 SQLite VACUUM 优化存储',
    cron_default: '0 3 * * 0',
    icon: 'setting'
  },
  {
    name: '进程健康检查',
    command: 'tasklist | findstr "catter"',
    desc: '检查 CMS 服务进程是否运行',
    cron_default: '*/15 * * * *',
    icon: 'circle-check'
  },
]

// ========== Cron 人性化选择相关 ==========
const selectedFreq = ref('')
const dailyTime = ref(null)
const weeklyDay = ref('1')
const weeklyTime = ref(null)
const monthlyDay = ref(1)
const monthlyTime = ref(null)

// 是否显示时间设置
const showTimeSetting = computed(() => {
  return ['daily', 'weekly', 'monthly'].includes(selectedFreq.value)
})

// 时间设置标签
const timeSettingLabel = computed(() => {
  const labels = {
    'daily': '选择时间',
    'weekly': '选择星期和时间',
    'monthly': '选择日期和时间'
  }
  return labels[selectedFreq.value] || ''
})

// 频率改变
const onFreqChange = (freq) => {
  selectedFreq.value = freq
  switch (freq) {
    case 'every_minute':
      taskForm.cron_expr = '* * * * *'
      break
    case 'every_5_min':
      taskForm.cron_expr = '*/5 * * * *'
      break
    case 'every_30_min':
      taskForm.cron_expr = '*/30 * * * *'
      break
    case 'hourly':
      taskForm.cron_expr = '0 * * * *'
      break
    case 'daily':
      // 默认凌晨2点
      dailyTime.value = new Date(2000, 0, 1, 2, 0)
      onTimeChange()
      break
    case 'weekly':
      weeklyDay.value = '1'
      weeklyTime.value = new Date(2000, 0, 1, 2, 0)
      onTimeChange()
      break
    case 'monthly':
      monthlyDay.value = 1
      monthlyTime.value = new Date(2000, 0, 1, 2, 0)
      onTimeChange()
      break
  }
}

// 时间改变
const onTimeChange = () => {
  const freq = selectedFreq.value
  if (freq === 'daily' && dailyTime.value) {
    const h = dailyTime.value.getHours()
    const m = dailyTime.value.getMinutes()
    taskForm.cron_expr = `${m} ${h} * * *`
  } else if (freq === 'weekly' && weeklyTime.value) {
    const h = weeklyTime.value.getHours()
    const m = weeklyTime.value.getMinutes()
    taskForm.cron_expr = `${m} ${h} * * ${weeklyDay.value}`
  } else if (freq === 'monthly' && monthlyTime.value) {
    const h = monthlyTime.value.getHours()
    const m = monthlyTime.value.getMinutes()
    taskForm.cron_expr = `${m} ${h} ${monthlyDay.value} * *`
  }
}

// 选择模板时，解析 cron_default 来设置频率
const selectTemplate = (tpl) => {
  taskForm.command = tpl.command
  if (!taskForm.name) {
    taskForm.name = tpl.name
  }
  if (tpl.cron_default) {
    taskForm.cron_expr = tpl.cron_default
    // 尝试解析 cron 表达式，设置对应的频率
    parseCronToFreq(tpl.cron_default)
  }
}

// 解析 cron 表达式为人性化频率
const parseCronToFreq = (cron) => {
  const parts = cron.split(/\s+/)
  if (parts.length !== 5) return
  
  const [min, hour, day, month, weekday] = parts
  
  if (min === '*' && hour === '*') {
    selectedFreq.value = 'every_minute'
  } else if (min === '*/5' && hour === '*') {
    selectedFreq.value = 'every_5_min'
  } else if (min === '*/30' && hour === '*') {
    selectedFreq.value = 'every_30_min'
  } else if (min === '0' && hour === '*') {
    selectedFreq.value = 'hourly'
  } else if (min !== '*' && hour !== '*' && day === '*' && month === '*' && weekday === '*') {
    selectedFreq.value = 'daily'
    dailyTime.value = new Date(2000, 0, 1, parseInt(hour), parseInt(min))
  } else if (min !== '*' && hour !== '*' && weekday !== '*' && weekday !== '*') {
    selectedFreq.value = 'weekly'
    weeklyDay.value = weekday
    weeklyTime.value = new Date(2000, 0, 1, parseInt(hour), parseInt(min))
  } else if (min !== '*' && hour !== '*' && day !== '*' && day !== '*') {
    selectedFreq.value = 'monthly'
    monthlyDay.value = parseInt(day)
    monthlyTime.value = new Date(2000, 0, 1, parseInt(hour), parseInt(min))
  }
}

// 按类别分组模板
const groupedTemplates = computed(() => {
  const groups = {
    '💾 数据备份': [],
    '🗑️ 文件清理': [],
    '🗃️ 数据库维护': [],
    '📊 监控统计': [],
    '⚙️ 系统维护': []
  }
  builtinTemplates.forEach(tpl => {
    if (tpl.name.includes('备份')) {
      groups['💾 数据备份'].push(tpl)
    } else if (tpl.name.includes('清理') && !tpl.name.includes('数据库')) {
      groups['🗑️ 文件清理'].push(tpl)
    } else if (tpl.name.includes('数据库') || tpl.name.includes('Token')) {
      groups['🗃️ 数据库维护'].push(tpl)
    } else if (tpl.name.includes('检查') || tpl.name.includes('统计')) {
      groups['📊 监控统计'].push(tpl)
    } else {
      groups['⚙️ 系统维护'].push(tpl)
    }
  })
  // 过滤空分组
  return Object.fromEntries(Object.entries(groups).filter(([_, items]) => items.length > 0))
})

// 查看执行日志
const openExecutionLog = async (task) => {
  currentTask.value = task
  showExecutionLog.value = true
  try {
    const response = await request.get(`/api/admin/scheduled-tasks/${task.id}/executions?limit=20`)
    if (response.code === 0) {
      executions.value = response.data || []
    }
  } catch (error) {
    ElMessage.error('获取执行日志失败')
  }
}

// 表单验证规则
const taskRules = {
  name: [{ required: true, message: '请输入任务名称', trigger: 'blur' }],
  task_type: [{ required: true, message: '请选择任务类型', trigger: 'change' }],
  cron_expr: [
    { required: true, message: '请输入Cron表达式', trigger: 'blur' },
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
    const response = await request.get('/api/admin/scheduled-tasks')
    if (response.code === 0) {
      tasks.value = response.data || []
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
    const response = await request.get('/api/admin/scheduler/status')
    if (response.code === 0) {
      schedulerStatus.value = response.data
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
        await request.put(`/api/admin/scheduled-tasks/${taskForm.id}`, taskForm)
        ElMessage.success('任务更新成功')
      } else {
        await request.post('/api/admin/scheduled-tasks', taskForm)
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
// 操作下拉菜单
const handleAction = (command, row) => {
  switch (command) {
    case 'edit':
      openEditDialog(row)
      break
    case 'delete':
      handleDelete(row)
      break
  }
}

const handleDelete = async (task) => {
  try {
    await ElMessageBox.confirm(`确定要删除任务「${task.name}」吗？`, '确认删除', {
      type: 'warning'
    })

    await request.delete(`/api/admin/scheduled-tasks/${task.id}`)
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
    await ElMessageBox.confirm(`确定要立即执行「${task.name}」吗？`, '确认执行', {
      type: 'warning'
    })

    task.running = true
    const res = await request.post(`/api/admin/scheduled-tasks/${task.id}/run`)
    if (res.code === 0) {
      ElMessage.success('任务已开始执行，请查看执行日志确认结果')
      // 3秒后自动刷新任务列表
      setTimeout(async () => {
        await loadTasks()
        await loadSchedulerStatus()
      }, 3000)
    } else {
      ElMessage.error(res.message || '任务执行失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('任务执行失败：' + (error.message || error))
    }
  } finally {
    task.running = false
  }
}

// 启用/禁用任务
const handleEnableChange = async (task) => {
  task.switching = true
  try {
    if (task.enabled) {
      await request.post(`/api/admin/scheduled-tasks/${task.id}/enable`)
      ElMessage.success('任务已启用')
    } else {
      await request.post(`/api/admin/scheduled-tasks/${task.id}/disable`)
      ElMessage.success('任务已禁用')
    }
    await loadSchedulerStatus()
  } catch (error) {
    task.enabled = !task.enabled
    ElMessage.error(task.enabled ? '任务启用失败' : '任务禁用失败')
  } finally {
    task.switching = false
  }
}

// 重启调度器
const restartScheduler = async () => {
  try {
    await ElMessageBox.confirm('确定要重启调度器吗？', '确认重启', {
      type: 'warning'
    })

    await request.post('/api/admin/scheduler/restart')
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
  flex-wrap: wrap;
}

.scheduler-status .el-tag {
  display: inline-flex;
}

.status-tag-content {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  white-space: nowrap;
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

/* Cron 人性化选择样式 */
.cron-quick-select {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.freq-group {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.freq-label {
  font-size: 13px;
  color: #606266;
  white-space: nowrap;
}

.time-setting {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  background: #f5f7fa;
  border-radius: 4px;
  flex-wrap: wrap;
}

.setting-label {
  font-size: 13px;
  color: #606266;
  white-space: nowrap;
}

.cron-preview {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 8px;
}

.cron-tips {
  margin-top: 5px;
  color: #909399;
}

.tip-link {
  color: #409eff;
  cursor: pointer;
}

.tip-link:hover {
  text-decoration: underline;
}

/* 内置任务模板分组卡片样式 */
.template-cards {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.template-group-title {
  font-size: 13px;
  font-weight: 600;
  color: #606266;
  margin-bottom: 8px;
  padding-left: 4px;
}

.template-group-items {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
  gap: 8px;
}

.template-card {
  padding: 10px 12px;
  border: 1px solid #dcdfe6;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s;
  background: #fafafa;
}

.template-card:hover {
  border-color: #409eff;
  background: #ecf5ff;
}

.template-card.active {
  border-color: #67c23a;
  background: #f0f9eb;
}

.template-card-name {
  font-size: 13px;
  font-weight: 500;
  color: #303133;
  margin-bottom: 4px;
}

.template-card-desc {
  font-size: 11px;
  color: #909399;
  line-height: 1.4;
}

.command-help {
  margin-top: 5px;
  color: #909399;
}

.next-run {
  font-size: 13px;
  color: #606266;
}

/* Cron帮助样式 */
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

/* 执行日志样式 */
.log-header {
  display: flex;
  align-items: center;
  gap: 15px;
  padding-bottom: 15px;
  margin-bottom: 15px;
  border-bottom: 1px solid #ebeef5;
}

.log-meta {
  color: #909399;
  font-size: 13px;
}

.execution-list {
  max-height: 500px;
  overflow-y: auto;
}

.execution-item {
  margin-bottom: 15px;
  padding: 10px;
  background: #f5f7fa;
  border-radius: 4px;
}

.execution-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 8px;
}

.execution-time {
  color: #606266;
  font-size: 13px;
}

.execution-duration {
  color: #909399;
  font-size: 12px;
}

.execution-output {
  background: #fff;
  padding: 10px;
  border-radius: 4px;
  font-size: 12px;
  font-family: 'Consolas', monospace;
  white-space: pre-wrap;
  word-break: break-all;
  max-height: 200px;
  overflow-y: auto;
  margin: 0;
  color: #303133;
}

/* 响应式 */
@media (max-width: 768px) {
  .scheduled-tasks-container {
    padding: 0;
  }

  .header-content {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }

  .header-actions {
    width: 100%;
    flex-wrap: wrap;
  }

  .scheduler-status {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }

  .status-stats {
    display: none;
  }

  .cron-quick-select {
    overflow-x: auto;
  }

  .template-cards {
    grid-template-columns: 1fr;
  }
}
</style>
