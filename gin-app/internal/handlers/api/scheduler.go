package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ScheduledTask 定时任务结构
type ScheduledTask struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	TaskType    string `json:"task_type"`
	CronExpr    string `json:"cron_expr"`
	Command     string `json:"command"`
	Enabled     bool   `json:"enabled"`
	Status      string `json:"status"`
	RunCount    int    `json:"run_count"`
	LastRunAt   string `json:"last_run_at"`
	NextRunAt   string `json:"next_run_at"`
	CreatedAt   string `json:"created_at"`
}

// SchedulerStatus 调度器状态
type SchedulerStatus struct {
	IsRunning      bool `json:"is_running"`
	TotalTasks     int  `json:"total_tasks"`
	EnabledTasks   int  `json:"enabled_tasks"`
	RunningTasks   int  `json:"running_tasks"`
	IdleTasks      int  `json:"idle_tasks"`
	FailedTasks    int  `json:"failed_tasks"`
}

// 内存存储（实际应该使用数据库）
var scheduledTasks = []ScheduledTask{
	{
		ID:        1,
		Name:      "每日数据备份",
		TaskType:  "backup",
		CronExpr:  "0 2 * * *",
		Command:   "backup:database",
		Enabled:   true,
		Status:    "idle",
		RunCount:  0,
		NextRunAt: "2026-03-19 02:00:00",
	},
	{
		ID:        2,
		Name:      "缓存清理",
		TaskType:  "cache_clear",
		CronExpr:  "0 */6 * * *",
		Command:   "cache:clear",
		Enabled:   true,
		Status:    "idle",
		RunCount:  5,
		LastRunAt: "2026-03-18 12:00:00",
		NextRunAt: "2026-03-18 18:00:00",
	},
	{
		ID:        3,
		Name:      "系统清理",
		TaskType:  "cleanup",
		CronExpr:  "0 3 * * 0",
		Command:   "system:cleanup",
		Enabled:   false,
		Status:    "idle",
		RunCount:  0,
		NextRunAt: "2026-03-22 03:00:00",
	},
}

// ListScheduledTasks 获取定时任务列表
func ListScheduledTasks(c *gin.Context) {
	// 返回模拟数据
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    scheduledTasks,
	})
}

// GetScheduledTask 获取单个定时任务
func GetScheduledTask(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	for _, task := range scheduledTasks {
		if task.ID == uint(id) {
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"data":    task,
			})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{
		"success": false,
		"message": "任务不存在",
	})
}

// CreateScheduledTask 创建定时任务
func CreateScheduledTask(c *gin.Context) {
	var task ScheduledTask
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "参数错误",
		})
		return
	}

	// 分配新ID
	task.ID = uint(len(scheduledTasks) + 1)
	task.Status = "idle"
	task.RunCount = 0
	scheduledTasks = append(scheduledTasks, task)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "任务创建成功",
		"data":    task,
	})
}

// UpdateScheduledTask 更新定时任务
func UpdateScheduledTask(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var updatedTask ScheduledTask
	if err := c.ShouldBindJSON(&updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "参数错误",
		})
		return
	}

	for i, task := range scheduledTasks {
		if task.ID == uint(id) {
			updatedTask.ID = uint(id)
			updatedTask.Status = task.Status
			updatedTask.RunCount = task.RunCount
			scheduledTasks[i] = updatedTask
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"message": "任务更新成功",
				"data":    updatedTask,
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"success": false,
		"message": "任务不存在",
	})
}

// DeleteScheduledTask 删除定时任务
func DeleteScheduledTask(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	for i, task := range scheduledTasks {
		if task.ID == uint(id) {
			scheduledTasks = append(scheduledTasks[:i], scheduledTasks[i+1:]...)
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"message": "任务删除成功",
			})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{
		"success": false,
		"message": "任务不存在",
	})
}

// ToggleScheduledTask 启用/禁用定时任务
func ToggleScheduledTask(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	for i, task := range scheduledTasks {
		if task.ID == uint(id) {
			scheduledTasks[i].Enabled = !scheduledTasks[i].Enabled
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"message": "任务状态已切换",
				"data":    scheduledTasks[i],
			})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{
		"success": false,
		"message": "任务不存在",
	})
}

// RunScheduledTask 立即执行定时任务
func RunScheduledTask(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	for i, task := range scheduledTasks {
		if task.ID == uint(id) {
			scheduledTasks[i].Status = "running"
			scheduledTasks[i].RunCount++
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"message": "任务已开始执行",
				"data":    scheduledTasks[i],
			})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{
		"success": false,
		"message": "任务不存在",
	})
}

// EnableScheduledTask 启用定时任务
func EnableScheduledTask(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	for i, task := range scheduledTasks {
		if task.ID == uint(id) {
			scheduledTasks[i].Enabled = true
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"message": "任务已启用",
				"data":    scheduledTasks[i],
			})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{
		"success": false,
		"message": "任务不存在",
	})
}

// DisableScheduledTask 禁用定时任务
func DisableScheduledTask(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	for i, task := range scheduledTasks {
		if task.ID == uint(id) {
			scheduledTasks[i].Enabled = false
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"message": "任务已禁用",
				"data":    scheduledTasks[i],
			})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{
		"success": false,
		"message": "任务不存在",
	})
}

// GetSchedulerStatus 获取调度器状态
func GetSchedulerStatus(c *gin.Context) {
	enabledCount := 0
	runningCount := 0
	for _, task := range scheduledTasks {
		if task.Enabled {
			enabledCount++
		}
		if task.Status == "running" {
			runningCount++
		}
	}

	status := SchedulerStatus{
		IsRunning:    true,
		TotalTasks:   len(scheduledTasks),
		EnabledTasks: enabledCount,
		RunningTasks: runningCount,
		IdleTasks:    len(scheduledTasks) - runningCount,
		FailedTasks:  0,
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    status,
	})
}

// RestartScheduler 重启调度器
func RestartScheduler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "调度器重启成功",
	})
}
