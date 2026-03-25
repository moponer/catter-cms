package models

import (
	"time"

	"gorm.io/gorm"
)

// ScheduledTask 定时任务模型
type ScheduledTask struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"size:200;not null"`
	TaskType  string         `json:"task_type" gorm:"size:50;not null"` // publish, backup, cleanup, cache_clear, custom
	CronExpr  string         `json:"cron_expr" gorm:"size:100;not null"`
	Command   string         `json:"command" gorm:"size:500"`
	Enabled   bool           `json:"enabled" gorm:"default:false"`
	Status    string         `json:"status" gorm:"size:20;default:'idle'"` // idle, running, failed
	RunCount  int            `json:"run_count" gorm:"default:0"`
	LastRunAt *time.Time     `json:"last_run_at"`
	NextRunAt *time.Time     `json:"next_run_at"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName 表名
func (ScheduledTask) TableName() string {
	return "scheduled_tasks"
}

// ScheduledTaskExecution 任务执行记录
type ScheduledTaskExecution struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	TaskID       uint           `json:"task_id" gorm:"index;not null"`
	TaskName     string         `json:"task_name" gorm:"size:200"`
	Status       string         `json:"status" gorm:"size:20;not null"` // success, failed
	Message      string         `json:"message" gorm:"type:text"`
	StartedAt    time.Time      `json:"started_at"`
	FinishedAt   *time.Time     `json:"finished_at"`
	Duration     int            `json:"duration"` // 执行耗时（秒）
	CreatedAt    time.Time      `json:"created_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName 表名
func (ScheduledTaskExecution) TableName() string {
	return "scheduled_task_executions"
}

// CreateScheduledTask 创建定时任务
func CreateScheduledTask(task *ScheduledTask) error {
	return DB.Create(task).Error
}

// GetScheduledTaskByID 根据ID获取任务
func GetScheduledTaskByID(id uint) (*ScheduledTask, error) {
	var task ScheduledTask
	err := DB.First(&task, id).Error
	return &task, err
}

// GetAllScheduledTasks 获取所有定时任务
func GetAllScheduledTasks() ([]ScheduledTask, error) {
	var tasks []ScheduledTask
	err := DB.Order("created_at DESC").Find(&tasks).Error
	return tasks, err
}

// GetEnabledScheduledTasks 获取所有已启用的任务
func GetEnabledScheduledTasks() ([]ScheduledTask, error) {
	var tasks []ScheduledTask
	err := DB.Where("enabled = ?", true).Find(&tasks).Error
	return tasks, err
}

// UpdateScheduledTask 更新任务
func UpdateScheduledTask(task *ScheduledTask) error {
	return DB.Save(task).Error
}

// DeleteScheduledTask 删除任务
func DeleteScheduledTask(id uint) error {
	return DB.Delete(&ScheduledTask{}, id).Error
}

// ToggleScheduledTask 切换任务启用状态
func ToggleScheduledTask(id uint) (*ScheduledTask, error) {
	var task ScheduledTask
	if err := DB.First(&task, id).Error; err != nil {
		return nil, err
	}
	task.Enabled = !task.Enabled
	if err := DB.Save(&task).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

// EnableScheduledTask 启用任务
func EnableScheduledTask(id uint) error {
	return DB.Model(&ScheduledTask{}).Where("id = ?", id).Update("enabled", true).Error
}

// DisableScheduledTask 禁用任务
func DisableScheduledTask(id uint) error {
	return DB.Model(&ScheduledTask{}).Where("id = ?", id).Update("enabled", false).Error
}

// UpdateScheduledTaskStatus 更新任务状态
func UpdateScheduledTaskStatus(id uint, status string) error {
	return DB.Model(&ScheduledTask{}).Where("id = ?", id).Update("status", status).Error
}

// IncrementRunCount 增加执行次数并更新最后执行时间
func IncrementRunCount(id uint) error {
	now := time.Now()
	return DB.Model(&ScheduledTask{}).Where("id = ?", id).Updates(map[string]interface{}{
		"run_count":  gorm.Expr("run_count + 1"),
		"last_run_at": now,
		"status":     "idle",
	}).Error
}

// UpdateNextRunAt 更新下次执行时间
func UpdateNextRunAt(id uint, nextRunAt time.Time) error {
	return DB.Model(&ScheduledTask{}).Where("id = ?", id).Update("next_run_at", nextRunAt).Error
}

// CreateTaskExecution 创建执行记录
func CreateTaskExecution(execution *ScheduledTaskExecution) error {
	return DB.Create(execution).Error
}

// GetTaskExecutions 获取任务执行记录
func GetTaskExecutions(taskID uint, limit int) ([]ScheduledTaskExecution, error) {
	var executions []ScheduledTaskExecution
	err := DB.Where("task_id = ?", taskID).Order("created_at DESC").Limit(limit).Find(&executions).Error
	return executions, err
}

// GetRecentExecutions 获取最近的执行记录
func GetRecentExecutions(limit int) ([]ScheduledTaskExecution, error) {
	var executions []ScheduledTaskExecution
	err := DB.Order("created_at DESC").Limit(limit).Find(&executions).Error
	return executions, err
}