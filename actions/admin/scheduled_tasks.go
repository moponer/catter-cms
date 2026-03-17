package admin

import (
	"time"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"catter-cms/actions/api"
	"catter-cms/models"
	"catter-cms/services"
)

// ScheduledTasksListAPI 获取定时任务列表API
// GET /api/admin/scheduled-tasks
func ScheduledTasksListAPI(c buffalo.Context) error {
	tasks := &models.ScheduledTasks{}
	query := models.DB.Order("created_at DESC")

	// 筛选 - 从查询字符串中获取参数
	enabled := c.Request().URL.Query().Get("enabled")
	if enabled != "" {
		query = query.Where("enabled = ?", enabled == "true")
	}

	taskType := c.Request().URL.Query().Get("task_type")
	if taskType != "" {
		query = query.Where("task_type = ?", taskType)
	}

	status := c.Request().URL.Query().Get("status")
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.All(tasks); err != nil {
		return api.InternalError(c, "获取定时任务列表失败")
	}

	return api.Success(c, *tasks)
}

// ScheduledTaskDetailAPI 获取定时任务详情API
// GET /api/admin/scheduled-tasks/{id}
func ScheduledTaskDetailAPI(c buffalo.Context) error {
	taskID := c.Param("id")
	task := &models.ScheduledTask{}

	if err := models.DB.Find(task, taskID); err != nil {
		return api.NotFound(c, "定时任务不存在")
	}

	return api.Success(c, task)
}

// ScheduledTaskCreateAPI 创建定时任务API
// POST /api/admin/scheduled-tasks
func ScheduledTaskCreateAPI(c buffalo.Context) error {
	task := &models.ScheduledTask{}
	if err := c.Bind(task); err != nil {
		return api.BadRequest(c, "参数错误")
	}

	// 设置默认值
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()
	if task.Status == "" {
		task.Status = "idle"
	}
	task.RunCount = 0

	// 验证数据
	verrs, err := task.Validate(models.DB)
	if err != nil {
		return api.InternalError(c, "验证失败")
	}
	if verrs.HasAny() {
		return api.BadRequest(c, verrs.Error())
	}

	// 保存任务
	if err := models.DB.Create(task); err != nil {
		return api.InternalError(c, "创建定时任务失败")
	}

	// 如果任务启用,添加到调度器
	if task.Enabled {
		scheduler := services.GetScheduler()
		if err := scheduler.AddTask(task); err != nil {
			return api.InternalError(c, "添加任务到调度器失败: "+err.Error())
		}
	}

	return api.SuccessWithMessage(c, "定时任务创建成功", task)
}

// ScheduledTaskUpdateAPI 更新定时任务API
// PUT /api/admin/scheduled-tasks/{id}
func ScheduledTaskUpdateAPI(c buffalo.Context) error {
	taskID := c.Param("id")
	task := &models.ScheduledTask{}

	if err := models.DB.Find(task, taskID); err != nil {
		return api.NotFound(c, "定时任务不存在")
	}

	if err := c.Bind(task); err != nil {
		return api.BadRequest(c, "参数错误")
	}

	// 更新时间
	task.UpdatedAt = time.Now()

	// 验证数据
	verrs, err := task.Validate(models.DB)
	if err != nil {
		return api.InternalError(c, "验证失败")
	}
	if verrs.HasAny() {
		return api.BadRequest(c, verrs.Error())
	}

	// 更新任务
	if err := models.DB.Update(task); err != nil {
		return api.InternalError(c, "更新定时任务失败")
	}

	// 更新调度器
	scheduler := services.GetScheduler()
	if task.Enabled {
		// 如果任务启用,添加或更新任务
		if err := scheduler.AddTask(task); err != nil {
			return api.InternalError(c, "更新任务调度失败: "+err.Error())
		}
	} else {
		// 如果任务禁用,从调度器移除
		scheduler.RemoveTask(task.ID)
	}

	return api.SuccessWithMessage(c, "定时任务更新成功", task)
}

// ScheduledTaskDeleteAPI 删除定时任务API
// DELETE /api/admin/scheduled-tasks/{id}
func ScheduledTaskDeleteAPI(c buffalo.Context) error {
	taskID := c.Param("id")
	task := &models.ScheduledTask{}

	if err := models.DB.Find(task, taskID); err != nil {
		return api.NotFound(c, "定时任务不存在")
	}

	// 从调度器移除
	scheduler := services.GetScheduler()
	scheduler.RemoveTask(task.ID)

	// 删除任务
	if err := models.DB.Destroy(task); err != nil {
		return api.InternalError(c, "删除定时任务失败")
	}

	return api.SuccessWithMessage(c, "定时任务删除成功", nil)
}

// ScheduledTaskRunNowAPI 立即执行定时任务API
// POST /api/admin/scheduled-tasks/{id}/run-now
func ScheduledTaskRunNowAPI(c buffalo.Context) error {
	taskID := c.Param("id")
	task := &models.ScheduledTask{}

	if err := models.DB.Find(task, taskID); err != nil {
		return api.NotFound(c, "定时任务不存在")
	}

	// 异步执行任务
	go func() {
		scheduler := services.GetScheduler()
		scheduler.ExecuteTask(task)
	}()

	return api.SuccessWithMessage(c, "定时任务已开始执行", nil)
}

// ScheduledTaskEnableAPI 启用定时任务API
// POST /api/admin/scheduled-tasks/{id}/enable
func ScheduledTaskEnableAPI(c buffalo.Context) error {
	taskID := c.Param("id")
	task := &models.ScheduledTask{}

	if err := models.DB.Find(task, taskID); err != nil {
		return api.NotFound(c, "定时任务不存在")
	}

	task.Enabled = true
	task.UpdatedAt = time.Now()

	if err := models.DB.Update(task); err != nil {
		return api.InternalError(c, "启用定时任务失败")
	}

	// 添加到调度器
	scheduler := services.GetScheduler()
	if err := scheduler.AddTask(task); err != nil {
		return api.InternalError(c, "添加任务到调度器失败: "+err.Error())
	}

	return api.SuccessWithMessage(c, "定时任务已启用", task)
}

// ScheduledTaskDisableAPI 禁用定时任务API
// POST /api/admin/scheduled-tasks/{id}/disable
func ScheduledTaskDisableAPI(c buffalo.Context) error {
	taskID := c.Param("id")
	task := &models.ScheduledTask{}

	if err := models.DB.Find(task, taskID); err != nil {
		return api.NotFound(c, "定时任务不存在")
	}

	task.Enabled = false
	task.UpdatedAt = time.Now()

	if err := models.DB.Update(task); err != nil {
		return api.InternalError(c, "禁用定时任务失败")
	}

	// 从调度器移除
	scheduler := services.GetScheduler()
	scheduler.RemoveTask(task.ID)

	return api.SuccessWithMessage(c, "定时任务已禁用", task)
}

// SchedulerStatusAPI 获取调度器状态API
// GET /api/admin/scheduler/status
func SchedulerStatusAPI(c buffalo.Context) error {
	scheduler := services.GetScheduler()

	// 统计任务状态
	tasks := []models.ScheduledTask{}
	models.DB.All(&tasks)

	stats := map[string]interface{}{
		"is_running":    scheduler.IsRunning(),
		"total_tasks":   len(tasks),
		"enabled_tasks": 0,
		"running_tasks": 0,
		"idle_tasks":    0,
		"failed_tasks":  0,
	}

	for _, task := range tasks {
		if task.Enabled {
			stats["enabled_tasks"] = stats["enabled_tasks"].(int) + 1
		}
		switch task.Status {
		case "running":
			stats["running_tasks"] = stats["running_tasks"].(int) + 1
		case "idle":
			stats["idle_tasks"] = stats["idle_tasks"].(int) + 1
		case "failed":
			stats["failed_tasks"] = stats["failed_tasks"].(int) + 1
		}
	}

	return c.Render(200, render.JSON(map[string]interface{}{
		"success": true,
		"message": "获取调度器状态成功",
		"data":    stats,
	}))
}

// SchedulerRestartAPI 重启调度器API
// POST /api/admin/scheduler/restart
func SchedulerRestartAPI(c buffalo.Context) error {
	scheduler := services.GetScheduler()

	if err := scheduler.Restart(); err != nil {
		return api.InternalError(c, "重启调度器失败: "+err.Error())
	}

	return api.SuccessWithMessage(c, "调度器重启成功", nil)
}
