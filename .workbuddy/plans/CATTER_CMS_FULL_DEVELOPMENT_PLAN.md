# Catter-CMS 完整对标Halo开发计划

## 项目目标
对标Halo CMS的所有核心功能,实现一个功能完备的现代化内容管理系统。

## 开发原则
- 原子级任务分解,每个任务可独立完成和测试
- 包含前端(Vue 3 + Element Plus)和后端(Go + Buffalo)实现
- 每个任务包含:数据库设计、后端API、前端界面、测试验证
- 按优先级分阶段实施

---

## 阶段一:P0核心功能(必须实现)

### 任务1:评论系统

#### 1.1 数据库设计
**文件**: `migrations/20260315000001_create_comments.sql`

```sql
CREATE TABLE comments (
    id UUID PRIMARY KEY,
    post_id UUID NOT NULL,
    user_id UUID,
    parent_id UUID DEFAULT NULL,
    author_name VARCHAR(255),
    author_email VARCHAR(255),
    author_url VARCHAR(500),
    content TEXT NOT NULL,
    status VARCHAR(20) DEFAULT 'pending', -- pending, approved, spam, trash
    ip_address VARCHAR(50),
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_comments_post_id ON comments(post_id);
CREATE INDEX idx_comments_user_id ON comments(user_id);
CREATE INDEX idx_comments_parent_id ON comments(parent_id);
CREATE INDEX idx_comments_status ON comments(status);
CREATE INDEX idx_comments_created_at ON comments(created_at);

ALTER TABLE comments ADD FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE;
ALTER TABLE comments ADD FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL;
ALTER TABLE comments ADD FOREIGN KEY (parent_id) REFERENCES comments(id) ON DELETE CASCADE;
```

**文件**: `models/comment.go`

```go
package models

import (
    "time"
    "github.com/gobuffalo/pop/v6"
    "github.com/gobuffalo/validate/v3"
    "github.com/gobuffalo/validate/v3/validators"
    "github.com/gofrs/uuid"
)

type Comment struct {
    ID          uuid.UUID  `db:"id" json:"id"`
    PostID      uuid.UUID  `db:"post_id" json:"post_id"`
    UserID      *uuid.UUID `db:"user_id" json:"user_id"`
    ParentID    *uuid.UUID `db:"parent_id" json:"parent_id"`
    AuthorName  string     `db:"author_name" json:"author_name"`
    AuthorEmail string     `db:"author_email" json:"author_email"`
    AuthorURL   string     `db:"author_url" json:"author_url"`
    Content     string     `db:"content" json:"content"`
    Status      string     `db:"status" json:"status"` // pending, approved, spam, trash
    IPAddress   string     `db:"ip_address" json:"ip_address"`
    UserAgent   string     `db:"user_agent" json:"user_agent"`
    CreatedAt   time.Time  `db:"created_at" json:"created_at"`
    UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`

    // Associations
    Post      *Post     `belongs_to:"post" json:"post,omitempty"`
    User      *User     `belongs_to:"user" json:"user,omitempty"`
    Parent    *Comment  `belongs_to:"comment" json:"parent,omitempty"`
    Replies   Comments  `has_many:"replies" json:"replies,omitempty"`
}

type Comments []*Comment

func (c *Comment) Validate(tx *pop.Connection) (*validate.Errors, error) {
    return validate.Validate(
        &validators.StringIsPresent{Field: c.Content, Name: "Content"},
        &validators.StringLengthInRange{Field: c.AuthorName, Name: "AuthorName", Min: 1, Max: 255},
        &validators.EmailLike{Field: c.AuthorEmail, Name: "AuthorEmail"},
    ), nil
}
```

#### 1.2 后端API实现

**文件**: `actions/admin/comments.go`

```go
package admin

import (
    "catter-cms/actions/api"
    "catter-cms/models"
    "strconv"

    "github.com/gobuffalo/buffalo"
    "github.com/gobuffalo/buffalo/render"
)

// CommentsListAPI 获取评论列表API
func CommentsListAPI(c buffalo.Context) error {
    page, _ := strconv.Atoi(c.Param("page"))
    pageSize, _ := strconv.Atoi(c.Param("page_size"))

    if page == 0 {
        page = 1
    }
    if pageSize == 0 {
        pageSize = 20
    }

    comments := &models.Comments{}
    query := models.DB.Paginate(page, pageSize)
    query = query.Order("created_at DESC")

    // 搜索
    keyword := c.Param("keyword")
    if keyword != "" {
        query = query.Where("content LIKE ? OR author_name LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
    }

    // 状态筛选
    status := c.Param("status")
    if status != "" {
        query = query.Where("status = ?", status)
    }

    // 文章筛选
    postID := c.Param("post_id")
    if postID != "" {
        query = query.Where("post_id = ?", postID)
    }

    // 预加载关联数据
    query = query.Eager("User", "Post")

    if err := query.All(comments); err != nil {
        return api.InternalError(c, "获取评论列表失败")
    }

    // 获取总数
    countQuery := models.DB.Q()
    if keyword != "" {
        countQuery = countQuery.Where("content LIKE ? OR author_name LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
    }
    if status != "" {
        countQuery = countQuery.Where("status = ?", status)
    }
    if postID != "" {
        countQuery = countQuery.Where("post_id = ?", postID)
    }
    total, _ := countQuery.Count(&models.Comment{})

    return api.PagedSuccess(c, comments, total, page, pageSize)
}

// CommentDetailAPI 获取评论详情API
func CommentDetailAPI(c buffalo.Context) error {
    commentID := c.Param("id")
    comment := &models.Comment{}

    // 预加载关联数据
    if err := models.DB.Eager("User", "Post", "Replies").Find(comment, commentID); err != nil {
        return api.NotFound(c, "评论不存在")
    }

    return api.Success(c, comment)
}

// CommentCreateAPI 创建评论API
func CommentCreateAPI(c buffalo.Context) error {
    comment := &models.Comment{}
    if err := c.Bind(comment); err != nil {
        return api.BadRequest(c, "参数错误")
    }

    // 设置IP和User-Agent
    comment.IPAddress = c.Request().RemoteAddr
    comment.UserAgent = c.Request().Header.Get("User-Agent")

    // 如果已登录,设置用户ID
    if currentUser, ok := c.Value("current_user").(*models.User); ok {
        comment.UserID = &currentUser.ID
        comment.AuthorName = currentUser.Name
        comment.AuthorEmail = currentUser.Email
    }

    // 默认状态为待审核
    if comment.Status == "" {
        comment.Status = "pending"
    }

    // 验证数据
    verrs, err := comment.Validate(models.DB)
    if err != nil {
        return api.InternalError(c, "验证失败")
    }
    if verrs.HasAny() {
        return api.BadRequest(c, verrs.Error())
    }

    // 保存评论
    if err := models.DB.Create(comment); err != nil {
        return api.InternalError(c, "创建评论失败")
    }

    return api.SuccessWithMessage(c, "评论创建成功", comment)
}

// CommentUpdateAPI 更新评论API
func CommentUpdateAPI(c buffalo.Context) error {
    commentID := c.Param("id")
    comment := &models.Comment{}

    if err := models.DB.Find(comment, commentID); err != nil {
        return api.NotFound(c, "评论不存在")
    }

    req := map[string]interface{}{}
    if err := c.Bind(&req); err != nil {
        return api.BadRequest(c, "参数错误")
    }

    // 更新字段
    if status, ok := req["status"].(string); ok {
        comment.Status = status
    }
    if content, ok := req["content"].(string); ok {
        comment.Content = content
    }
    if authorName, ok := req["author_name"].(string); ok {
        comment.AuthorName = authorName
    }
    if authorEmail, ok := req["author_email"].(string); ok {
        comment.AuthorEmail = authorEmail
    }

    // 验证数据
    verrs, err := comment.Validate(models.DB)
    if err != nil {
        return api.InternalError(c, "验证失败")
    }
    if verrs.HasAny() {
        return api.BadRequest(c, verrs.Error())
    }

    // 更新评论
    if err := models.DB.Update(comment); err != nil {
        return api.InternalError(c, "更新评论失败")
    }

    return api.SuccessWithMessage(c, "评论更新成功", comment)
}

// CommentDeleteAPI 删除评论API
func CommentDeleteAPI(c buffalo.Context) error {
    commentID := c.Param("id")
    comment := &models.Comment{}

    if err := models.DB.Find(comment, commentID); err != nil {
        return api.NotFound(c, "评论不存在")
    }

    if err := models.DB.Destroy(comment); err != nil {
        return api.InternalError(c, "删除评论失败")
    }

    return api.SuccessWithMessage(c, "评论删除成功", nil)
}

// CommentApproveAPI 批准评论API
func CommentApproveAPI(c buffalo.Context) error {
    commentID := c.Param("id")
    comment := &models.Comment{}

    if err := models.DB.Find(comment, commentID); err != nil {
        return api.NotFound(c, "评论不存在")
    }

    comment.Status = "approved"
    if err := models.DB.Update(comment); err != nil {
        return api.InternalError(c, "批准评论失败")
    }

    return api.SuccessWithMessage(c, "评论已批准", comment)
}

// CommentRejectAPI 拒绝评论API
func CommentRejectAPI(c buffalo.Context) error {
    commentID := c.Param("id")
    comment := &models.Comment{}

    if err := models.DB.Find(comment, commentID); err != nil {
        return api.NotFound(c, "评论不存在")
    }

    comment.Status = "spam"
    if err := models.DB.Update(comment); err != nil {
        return api.InternalError(c, "拒绝评论失败")
    }

    return api.SuccessWithMessage(c, "评论已拒绝", comment)
}

// CommentBatchDeleteAPI 批量删除评论API
func CommentBatchDeleteAPI(c buffalo.Context) error {
    ids := []string{}
    if err := c.Bind(&ids); err != nil {
        return api.BadRequest(c, "参数错误")
    }

    for _, id := range ids {
        comment := &models.Comment{}
        if err := models.DB.Find(comment, id); err != nil {
            continue
        }
        models.DB.Destroy(comment)
    }

    return api.SuccessWithMessage(c, "批量删除成功", nil)
}

// CommentBatchApproveAPI 批量批准评论API
func CommentBatchApproveAPI(c buffalo.Context) error {
    ids := []string{}
    if err := c.Bind(&ids); err != nil {
        return api.BadRequest(c, "参数错误")
    }

    comments := &models.Comments{}
    if err := models.DB.Where("id IN (?)", ids).All(comments); err != nil {
        return api.InternalError(c, "批量批准失败")
    }

    for _, comment := range *comments {
        comment.Status = "approved"
        models.DB.Update(comment)
    }

    return api.SuccessWithMessage(c, "批量批准成功", nil)
}

// CommentsIndexHandler 评论列表页面 (兼容旧API)
func CommentsIndexHandler(c buffalo.Context) error {
    return c.Render(200, render.HTML("<html><body><p>Comments SPA</p></body></html>"))
}

// unused - avoid compile error
var _ = strconv.Itoa
```

**文件**: `actions/front/comments.go`

```go
package front

import (
    "catter-cms/actions/api"
    "catter-cms/models"

    "github.com/gobuffalo/buffalo"
    "github.com/gobuffalo/buffalo/render"
)

// FrontCommentsListAPI 获取文章评论列表API(前台)
func FrontCommentsListAPI(c buffalo.Context) error {
    postID := c.Param("post_id")

    comments := &models.Comments{}
    query := models.DB.Where("post_id = ?", postID)
    query = query.Where("status = ?", "approved") // 只显示已批准的评论
    query = query.Order("created_at ASC") // 按时间正序

    // 预加载关联数据
    query = query.Eager("User", "Replies")

    if err := query.All(comments); err != nil {
        return api.InternalError(c, "获取评论列表失败")
    }

    return api.Success(c, comments)
}

// FrontCommentCreateAPI 创建评论API(前台)
func FrontCommentCreateAPI(c buffalo.Context) error {
    comment := &models.Comment{}
    if err := c.Bind(comment); err != nil {
        return api.BadRequest(c, "参数错误")
    }

    // 设置IP和User-Agent
    comment.IPAddress = c.Request().RemoteAddr
    comment.UserAgent = c.Request().Header.Get("User-Agent")

    // 如果已登录,设置用户ID
    if currentUser, ok := c.Value("current_user").(*models.User); ok {
        comment.UserID = &currentUser.ID
        comment.AuthorName = currentUser.Name
        comment.AuthorEmail = currentUser.Email
    }

    // 默认状态为待审核
    comment.Status = "pending"

    // 验证数据
    verrs, err := comment.Validate(models.DB)
    if err != nil {
        return api.InternalError(c, "验证失败")
    }
    if verrs.HasAny() {
        return api.BadRequest(c, verrs.Error())
    }

    // 保存评论
    if err := models.DB.Create(comment); err != nil {
        return api.InternalError(c, "创建评论失败")
    }

    return api.SuccessWithMessage(c, "评论提交成功,等待审核", comment)
}

// unused - avoid compile error
var _ = render.HTML
```

**文件**: `actions/app.go` (添加路由)

在 `app.go` 中添加评论相关路由:

```go
// Admin评论路由
admin.GET("/comments", admin.CommentsIndexHandler)
admin.GET("/comments/{id}", admin.CommentDetailAPI)
admin.POST("/comments", admin.CommentCreateAPI)
admin.PUT("/comments/{id}", admin.CommentUpdateAPI)
admin.DELETE("/comments/{id}", admin.CommentDeleteAPI)
admin.POST("/comments/{id}/approve", admin.CommentApproveAPI)
admin.POST("/comments/{id}/reject", admin.CommentRejectAPI)
admin.POST("/comments/batch/delete", admin.CommentBatchDeleteAPI)
admin.POST("/comments/batch/approve", admin.CommentBatchApproveAPI)

// 前台评论API
api.GET("/posts/{post_id}/comments", front.FrontCommentsListAPI)
api.POST("/comments", front.FrontCommentCreateAPI)
```

#### 1.3 前端管理界面实现

**文件**: `admin/src/views/comments/Index.vue`

```vue
<template>
  <div class="comments-page">
    <div class="page-header">
      <h1>评论管理</h1>
      <el-button type="primary" @click="refreshList">
        <el-icon><Refresh /></el-icon>
        刷新
      </el-button>
    </div>

    <!-- 搜索筛选 -->
    <el-card class="filter-card">
      <el-form :inline="true" :model="filters" class="filter-form">
        <el-form-item label="搜索">
          <el-input
            v-model="filters.keyword"
            placeholder="搜索评论内容或作者"
            clearable
            @clear="handleSearch"
            @keyup.enter="handleSearch"
          >
            <template #append>
              <el-button :icon="Search" @click="handleSearch" />
            </template>
          </el-input>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="filters.status" placeholder="全部状态" clearable @change="handleSearch">
            <el-option label="待审核" value="pending" />
            <el-option label="已批准" value="approved" />
            <el-option label="垃圾评论" value="spam" />
          </el-select>
        </el-form-item>
        <el-form-item label="文章">
          <el-select v-model="filters.post_id" placeholder="选择文章" clearable filterable @change="handleSearch">
            <el-option
              v-for="post in posts"
              :key="post.id"
              :label="post.title"
              :value="post.id"
            />
          </el-select>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 评论列表 -->
    <el-card class="comments-card">
      <template #header>
        <div class="card-header">
          <span>评论列表</span>
          <div class="header-actions">
            <el-checkbox v-model="selectAll" @change="handleSelectAll">全选</el-checkbox>
            <el-button
              type="success"
              size="small"
              :disabled="selectedComments.length === 0"
              @click="handleBatchApprove"
            >
              批量批准
            </el-button>
            <el-button
              type="danger"
              size="small"
              :disabled="selectedComments.length === 0"
              @click="handleBatchDelete"
            >
              批量删除
            </el-button>
          </div>
        </div>
      </template>

      <div v-loading="loading" class="comments-list">
        <el-empty v-if="!loading && comments.length === 0" description="暂无评论" />

        <div v-for="comment in comments" :key="comment.id" class="comment-item">
          <el-checkbox v-model="selectedComments" :value="comment.id" class="comment-checkbox" />
          <div class="comment-content">
            <div class="comment-header">
              <div class="comment-author">
                <el-avatar :size="32" :src="comment.user?.avatar" class="author-avatar">
                  {{ comment.author_name?.charAt(0).toUpperCase() }}
                </el-avatar>
                <div class="author-info">
                  <span class="author-name">{{ comment.author_name }}</span>
                  <el-tag :type="getStatusType(comment.status)" size="small" class="status-tag">
                    {{ getStatusText(comment.status) }}
                  </el-tag>
                </div>
              </div>
              <div class="comment-meta">
                <span class="comment-date">{{ formatDate(comment.created_at) }}</span>
                <span class="comment-post">{{ comment.post?.title }}</span>
              </div>
            </div>

            <div class="comment-text">{{ comment.content }}</div>

            <div class="comment-actions">
              <el-button
                v-if="comment.status === 'pending'"
                type="success"
                size="small"
                link
                @click="handleApprove(comment.id)"
              >
                <el-icon><Select /></el-icon>
                批准
              </el-button>
              <el-button
                v-if="comment.status === 'pending'"
                type="danger"
                size="small"
                link
                @click="handleReject(comment.id)"
              >
                <el-icon><Close /></el-icon>
                拒绝
              </el-button>
              <el-button type="primary" size="small" link @click="handleEdit(comment)">
                <el-icon><Edit /></el-icon>
                编辑
              </el-button>
              <el-button type="danger" size="small" link @click="handleDelete(comment.id)">
                <el-icon><Delete /></el-icon>
                删除
              </el-button>
            </div>

            <!-- 回复列表 -->
            <div v-if="comment.replies && comment.replies.length > 0" class="replies-list">
              <div v-for="reply in comment.replies" :key="reply.id" class="reply-item">
                <el-avatar :size="24" class="reply-avatar">
                  {{ reply.author_name?.charAt(0).toUpperCase() }}
                </el-avatar>
                <div class="reply-content">
                  <span class="reply-author">{{ reply.author_name }}</span>
                  <span class="reply-text">{{ reply.content }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 分页 -->
      <div class="pagination">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :total="pagination.total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @current-change="loadComments"
          @size-change="loadComments"
        />
      </div>
    </el-card>

    <!-- 编辑对话框 -->
    <el-dialog v-model="editDialogVisible" title="编辑评论" width="600px">
      <el-form :model="editForm" label-width="80px">
        <el-form-item label="作者姓名">
          <el-input v-model="editForm.author_name" />
        </el-form-item>
        <el-form-item label="作者邮箱">
          <el-input v-model="editForm.author_email" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="editForm.status">
            <el-option label="待审核" value="pending" />
            <el-option label="已批准" value="approved" />
            <el-option label="垃圾评论" value="spam" />
          </el-select>
        </el-form-item>
        <el-form-item label="评论内容">
          <el-input
            v-model="editForm.content"
            type="textarea"
            :rows="4"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSaveEdit">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Refresh,
  Search,
  Select,
  Close,
  Edit,
  Delete
} from '@element-plus/icons-vue'
import { get, post, put, del } from '@/utils/request'

const loading = ref(false)
const comments = ref([])
const posts = ref([])
const selectedComments = ref([])
const selectAll = ref(false)
const editDialogVisible = ref(false)
const editForm = reactive({
  id: '',
  author_name: '',
  author_email: '',
  status: '',
  content: ''
})

const filters = reactive({
  keyword: '',
  status: '',
  post_id: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

const loadPosts = async () => {
  try {
    const res = await get('/api/posts', { page_size: 100 })
    posts.value = res.data || res
  } catch (error) {
    console.error('加载文章列表失败:', error)
  }
}

const loadComments = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.pageSize,
      ...filters
    }
    const res = await get('/api/admin/comments', params)
    comments.value = res.data || res.items || []
    pagination.total = res.total || 0
  } catch (error) {
    ElMessage.error('加载评论列表失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.page = 1
  selectedComments.value = []
  selectAll.value = false
  loadComments()
}

const refreshList = () => {
  loadComments()
}

const handleSelectAll = () => {
  if (selectAll.value) {
    selectedComments.value = comments.value.map(c => c.id)
  } else {
    selectedComments.value = []
  }
}

const handleApprove = async (id) => {
  try {
    await post(`/api/admin/comments/${id}/approve`)
    ElMessage.success('评论已批准')
    loadComments()
  } catch (error) {
    ElMessage.error('批准评论失败')
  }
}

const handleReject = async (id) => {
  try {
    await post(`/api/admin/comments/${id}/reject`)
    ElMessage.success('评论已拒绝')
    loadComments()
  } catch (error) {
    ElMessage.error('拒绝评论失败')
  }
}

const handleEdit = (comment) => {
  Object.assign(editForm, {
    id: comment.id,
    author_name: comment.author_name,
    author_email: comment.author_email,
    status: comment.status,
    content: comment.content
  })
  editDialogVisible.value = true
}

const handleSaveEdit = async () => {
  try {
    await put(`/api/admin/comments/${editForm.id}`, editForm)
    ElMessage.success('评论更新成功')
    editDialogVisible.value = false
    loadComments()
  } catch (error) {
    ElMessage.error('更新评论失败')
  }
}

const handleDelete = async (id) => {
  try {
    await ElMessageBox.confirm('确定要删除这条评论吗?', '提示', {
      type: 'warning'
    })
    await del(`/api/admin/comments/${id}`)
    ElMessage.success('评论删除成功')
    loadComments()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除评论失败')
    }
  }
}

const handleBatchApprove = async () => {
  try {
    await post('/api/admin/comments/batch/approve', selectedComments.value)
    ElMessage.success('批量批准成功')
    selectedComments.value = []
    selectAll.value = false
    loadComments()
  } catch (error) {
    ElMessage.error('批量批准失败')
  }
}

const handleBatchDelete = async () => {
  try {
    await ElMessageBox.confirm(`确定要删除选中的 ${selectedComments.value.length} 条评论吗?`, '提示', {
      type: 'warning'
    })
    await post('/api/admin/comments/batch/delete', selectedComments.value)
    ElMessage.success('批量删除成功')
    selectedComments.value = []
    selectAll.value = false
    loadComments()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('批量删除失败')
    }
  }
}

const getStatusType = (status) => {
  const typeMap = {
    pending: 'warning',
    approved: 'success',
    spam: 'danger'
  }
  return typeMap[status] || 'info'
}

const getStatusText = (status) => {
  const textMap = {
    pending: '待审核',
    approved: '已批准',
    spam: '垃圾评论'
  }
  return textMap[status] || status
}

const formatDate = (date) => {
  return new Date(date).toLocaleString('zh-CN')
}

onMounted(() => {
  loadPosts()
  loadComments()
})
</script>

<style scoped>
.comments-page {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-lg);
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.page-header h1 {
  font-size: var(--font-size-2xl);
  font-weight: 600;
  margin: 0;
}

.filter-card {
  margin-bottom: var(--spacing-md);
}

.filter-form {
  display: flex;
  flex-wrap: wrap;
  gap: var(--spacing-sm);
}

.filter-form .el-form-item {
  margin-bottom: 0;
}

.comments-card {
  flex: 1;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-actions {
  display: flex;
  gap: var(--spacing-sm);
  align-items: center;
}

.comments-list {
  min-height: 400px;
}

.comment-item {
  display: flex;
  gap: var(--spacing-md);
  padding: var(--spacing-md);
  border-bottom: 1px solid var(--color-border);
  transition: background-color var(--transition-fast);
}

.comment-item:hover {
  background-color: var(--color-bg-hover);
}

.comment-checkbox {
  margin-top: var(--spacing-md);
}

.comment-content {
  flex: 1;
}

.comment-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: var(--spacing-sm);
}

.comment-author {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.author-avatar {
  flex-shrink: 0;
}

.author-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.author-name {
  font-weight: 500;
  color: var(--color-text-primary);
}

.status-tag {
  align-self: flex-start;
}

.comment-meta {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 4px;
  font-size: 12px;
  color: var(--color-text-secondary);
}

.comment-date {
  color: var(--color-text-secondary);
}

.comment-post {
  color: var(--color-primary);
}

.comment-text {
  margin: var(--spacing-md) 0;
  line-height: 1.6;
  color: var(--color-text-primary);
}

.comment-actions {
  display: flex;
  gap: var(--spacing-sm);
  margin-top: var(--spacing-sm);
}

.replies-list {
  margin-top: var(--spacing-md);
  padding-left: var(--spacing-lg);
  border-left: 2px solid var(--color-border);
}

.reply-item {
  display: flex;
  gap: var(--spacing-sm);
  margin-bottom: var(--spacing-sm);
}

.reply-avatar {
  flex-shrink: 0;
}

.reply-content {
  flex: 1;
  background-color: var(--color-bg-light);
  padding: var(--spacing-sm) var(--spacing-md);
  border-radius: var(--radius-sm);
}

.reply-author {
  font-weight: 500;
  color: var(--color-text-primary);
  margin-right: var(--spacing-sm);
}

.reply-text {
  color: var(--color-text-secondary);
}

.pagination {
  display: flex;
  justify-content: center;
  margin-top: var(--spacing-lg);
}
</style>
```

**文件**: `admin/src/router/index.js` (添加评论路由)

```javascript
{
  path: '/comments',
  name: 'Comments',
  component: () => import('@/views/comments/Index.vue'),
  meta: { title: '评论管理', requiresAuth: true }
}
```

#### 1.4 测试验证

```sql
-- 插入测试评论数据
INSERT INTO comments (id, post_id, author_name, author_email, content, status, created_at)
VALUES
    (gen_random_uuid(), (SELECT id FROM posts LIMIT 1), '张三', 'zhangsan@example.com', '这是一条测试评论', 'pending', NOW()),
    (gen_random_uuid(), (SELECT id FROM posts LIMIT 1), '李四', 'lisi@example.com', '这是一条已批准的评论', 'approved', NOW());
```

测试清单:
- ✅ 评论列表展示
- ✅ 评论搜索功能
- ✅ 评论状态筛选
- ✅ 批准/拒绝评论
- ✅ 批量操作
- ✅ 评论编辑
- ✅ 评论删除

---

## 继续下一个任务...
