package api

import (
	"catter-cms-gin/internal/models"
	"catter-cms-gin/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ListScheduledTasks 获取定时任务列表
func ListScheduledTasks(c *gin.Context) {
	tasks, err := models.GetAllScheduledTasks()
	if err != nil {
		InternalError(c, "获取任务列表失败: "+err.Error())
		return
	}

	// 计算每个任务的下次执行时间
	for i := range tasks {
		if tasks[i].Enabled && tasks[i].CronExpr != "" {
			if nextRun := services.CalculateNextRun(tasks[i].CronExpr); nextRun != nil {
				tasks[i].NextRunAt = nextRun
			}
		}
	}

	SuccessResponse(c, tasks)
}

// GetScheduledTask 获取单个定时任务
func GetScheduledTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		BadRequest(c, "无效的任务ID")
		return
	}

	task, err := models.GetScheduledTaskByID(uint(id))
	if err != nil {
		NotFound(c, "任务不存在")
		return
	}

	// 计算下次执行时间
	if task.Enabled && task.CronExpr != "" {
		if nextRun := services.CalculateNextRun(task.CronExpr); nextRun != nil {
			task.NextRunAt = nextRun
		}
	}

	SuccessResponse(c, task)
}

// CreateScheduledTask 创建定时任务
func CreateScheduledTask(c *gin.Context) {
	var task models.ScheduledTask
	if err := c.ShouldBindJSON(&task); err != nil {
		BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 验证 Cron 表达式
	if task.CronExpr == "" {
		BadRequest(c, "请输入Cron表达式")
		return
	}

	// 设置默认值
	task.Status = "idle"
	task.RunCount = 0

	// 创建任务
	if err := models.CreateScheduledTask(&task); err != nil {
		InternalError(c, "创建任务失败: "+err.Error())
		return
	}

	// 如果启用，添加到调度器
	if task.Enabled {
		scheduler := services.GetScheduler()
		scheduler.AddTask(&task)
	}

	SuccessWithMessage(c, "任务创建成功", task)
}

// UpdateScheduledTask 更新定时任务
func UpdateScheduledTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		BadRequest(c, "无效的任务ID")
		return
	}

	var updatedTask models.ScheduledTask
	if err := c.ShouldBindJSON(&updatedTask); err != nil {
		BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 获取现有任务
	task, err := models.GetScheduledTaskByID(uint(id))
	if err != nil {
		NotFound(c, "任务不存在")
		return
	}

	// 更新字段
	task.Name = updatedTask.Name
	task.TaskType = updatedTask.TaskType
	task.CronExpr = updatedTask.CronExpr
	task.Command = updatedTask.Command
	task.Enabled = updatedTask.Enabled

	if err := models.UpdateScheduledTask(task); err != nil {
		InternalError(c, "更新任务失败: "+err.Error())
		return
	}

	// 更新调度器中的任务
	scheduler := services.GetScheduler()
	scheduler.RemoveTask(uint(id))
	if task.Enabled {
		scheduler.AddTask(task)
	}

	SuccessWithMessage(c, "任务更新成功", task)
}

// DeleteScheduledTask 删除定时任务
func DeleteScheduledTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		BadRequest(c, "无效的任务ID")
		return
	}

	// 从调度器移除
	scheduler := services.GetScheduler()
	scheduler.RemoveTask(uint(id))

	// 删除任务
	if err := models.DeleteScheduledTask(uint(id)); err != nil {
		InternalError(c, "删除任务失败: "+err.Error())
		return
	}

	SuccessWithMessage(c, "任务删除成功", nil)
}

// ToggleScheduledTask 切换任务启用状态
func ToggleScheduledTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		BadRequest(c, "无效的任务ID")
		return
	}

	task, err := models.ToggleScheduledTask(uint(id))
	if err != nil {
		InternalError(c, "切换任务状态失败: "+err.Error())
		return
	}

	// 更新调度器
	scheduler := services.GetScheduler()
	scheduler.RemoveTask(uint(id))
	if task.Enabled {
		scheduler.AddTask(task)
	}

	SuccessResponse(c, task)
}

// EnableScheduledTask 启用定时任务
func EnableScheduledTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		BadRequest(c, "无效的任务ID")
		return
	}

	if err := models.EnableScheduledTask(uint(id)); err != nil {
		InternalError(c, "启用任务失败: "+err.Error())
		return
	}

	task, _ := models.GetScheduledTaskByID(uint(id))

	// 添加到调度器
	scheduler := services.GetScheduler()
	if task.Enabled {
		scheduler.AddTask(task)
	}

	SuccessWithMessage(c, "任务已启用", task)
}

// DisableScheduledTask 禁用定时任务
func DisableScheduledTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		BadRequest(c, "无效的任务ID")
		return
	}

	if err := models.DisableScheduledTask(uint(id)); err != nil {
		InternalError(c, "禁用任务失败: "+err.Error())
		return
	}

	// 从调度器移除
	scheduler := services.GetScheduler()
	scheduler.RemoveTask(uint(id))

	SuccessWithMessage(c, "任务已禁用", nil)
}

// RunScheduledTask 立即执行定时任务
func RunScheduledTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		BadRequest(c, "无效的任务ID")
		return
	}

	task, err := models.GetScheduledTaskByID(uint(id))
	if err != nil {
		NotFound(c, "任务不存在")
		return
	}

	// 异步执行任务
	scheduler := services.GetScheduler()
	if err := scheduler.RunTaskNow(uint(id)); err != nil {
		InternalError(c, "执行任务失败: "+err.Error())
		return
	}

	// 更新状态
	models.UpdateScheduledTaskStatus(uint(id), "running")

	SuccessWithMessage(c, "任务已开始执行", task)
}

// GetSchedulerStatus 获取调度器状态
func GetSchedulerStatus(c *gin.Context) {
	scheduler := services.GetScheduler()
	status := scheduler.GetStatus()
	SuccessResponse(c, status)
}

// RestartScheduler 重启调度器
func RestartScheduler(c *gin.Context) {
	scheduler := services.GetScheduler()
	scheduler.Stop()
	scheduler.Start()
	SuccessWithMessage(c, "调度器重启成功", nil)
}

// GetTaskExecutions 获取任务执行记录
func GetTaskExecutions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		BadRequest(c, "无效的任务ID")
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	executions, err := models.GetTaskExecutions(uint(id), limit)
	if err != nil {
		InternalError(c, "获取执行记录失败: "+err.Error())
		return
	}

	SuccessResponse(c, executions)
}

// GetAllExecutions 获取所有执行记录
func GetAllExecutions(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	executions, err := models.GetRecentExecutions(limit)
	if err != nil {
		InternalError(c, "获取执行记录失败: "+err.Error())
		return
	}

	SuccessResponse(c, executions)
}

// ValidateCronExpr 验证Cron表达式
func ValidateCronExpr(c *gin.Context) {
	expr := c.Query("expr")
	if expr == "" {
		BadRequest(c, "请输入Cron表达式")
		return
	}

	nextRun := services.CalculateNextRun(expr)
	if nextRun == nil {
		BadRequest(c, "无效的Cron表达式")
		return
	}

	SuccessResponse(c, gin.H{
		"valid":       true,
		"next_run_at": nextRun.Format("2006-01-02 15:04:05"),
		"description": describeCronExpr(expr),
	})
}

// describeCronExpr 描述Cron表达式的含义
func describeCronExpr(expr string) string {
	next := services.CalculateNextRun(expr)
	if next == nil {
		return "无法解析"
	}

	return "下次执行时间: " + next.Format("2006-01-02 15:04:05")
}