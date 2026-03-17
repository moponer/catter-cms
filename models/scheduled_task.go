package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
)

// ScheduledTask 定时任务模型
type ScheduledTask struct {
	ID          int       `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`                   // 任务名称
	TaskType    string    `db:"task_type" json:"task_type"`           // 任务类型: publish, backup, cleanup, cache_clear, custom
	CronExpr    string    `db:"cron_expr" json:"cron_expr"`           // Cron表达式
	Command     string    `db:"command" json:"command"`               // 自定义命令(仅custom类型)
	Enabled     bool      `db:"enabled" json:"enabled"`               // 是否启用
	LastRunAt   *time.Time `db:"last_run_at" json:"last_run_at"`       // 最后执行时间
	NextRunAt   *time.Time `db:"next_run_at" json:"next_run_at"`       // 下次执行时间
	RunCount    int       `db:"run_count" json:"run_count"`           // 执行次数
	Status      string    `db:"status" json:"status"`                 // 状态: idle, running, failed
	ErrorMsg    *string   `db:"error_msg" json:"error_msg,omitempty"` // 错误信息
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

// ScheduledTasks 定时任务集合
type ScheduledTasks []ScheduledTask

// TableName 指定表名
func (ScheduledTask) TableName() string {
	return "scheduled_tasks"
}

// Validate 验证定时任务数据
func (t *ScheduledTask) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Name: "Name", Field: t.Name},
		&validators.StringIsPresent{Name: "TaskType", Field: t.TaskType},
		&validators.StringInclusion{Name: "TaskType", Field: t.TaskType, List: []string{"publish", "backup", "cleanup", "cache_clear", "custom"}},
		&validators.StringIsPresent{Name: "CronExpr", Field: t.CronExpr},
	), nil
}
