# 修复 Category.PostCount 字段缺失问题

## 📅 修复时间
2026-03-16

## 🐛 问题描述

访问主题页面时出现错误：
```
themes/catter-blog/pages/home.plush.html: line 298:
'category' does not have a field or method named 'PostCount' (category.PostCount)
```

**原因**: 模板中使用了 `category.PostCount` 字段，但 Category 模型中没有定义该字段，控制器也没有计算并传递该数据。

---

## 🔧 修复方案

### 方案选择

有两种解决方案：

#### 方案 1: 添加 PostCount 字段到模型（不推荐）
- ❌ 需要数据库迁移
- ❌ 需要维护计数（增删文章时更新）
- ❌ 可能导致数据不一致

#### 方案 2: 在控制器中动态计算（推荐）✅
- ✅ 不需要修改数据库
- ✅ 数据始终准确
- ✅ 实现简单

**选择方案 2**

---

### 实现步骤

#### 1. 创建带 PostCount 的结构体

**文件**: `actions/front/home.go`

```go
// 为分类添加文章数量
type CategoryWithCount struct {
    ID          int
    Name        string
    Slug        string
    Description string
    PostCount   int
}

// 为标签添加文章数量
type TagWithCount struct {
    ID        int
    Name      string
    Slug      string
    PostCount int
}
```

#### 2. 在控制器中计算文章数量

**修改前**:
```go
// 获取所有分类
categories := models.Categories{}
if err := models.DB.All(&categories); err != nil {
    return c.Error(http.StatusInternalServerError, err)
}

// 获取所有标签
tags := models.Tags{}
if err := models.DB.All(&tags); err != nil {
    return c.Error(http.StatusInternalServerError, err)
}
```

**修改后**:
```go
// 获取所有分类（带文章数量）
categories := models.Categories{}
if err := models.DB.All(&categories); err != nil {
    return c.Error(http.StatusInternalServerError, err)
}

// 为每个分类计算文章数量
categoriesWithCount := []CategoryWithCount{}
for _, cat := range categories {
    var count int
    models.DB.Where("category_id = ? AND status = ?", cat.ID, "published").Count(&count)
    categoriesWithCount = append(categoriesWithCount, CategoryWithCount{
        ID:          cat.ID,
        Name:        cat.Name,
        Slug:        cat.Slug,
        Description: cat.Description,
        PostCount:   count,
    })
}

// 获取所有标签（带文章数量）
tags := models.Tags{}
if err := models.DB.All(&tags); err != nil {
    return c.Error(http.StatusInternalServerError, err)
}

// 为每个标签计算文章数量
tagsWithCount := []TagWithCount{}
for _, tag := range tags {
    var count int
    models.DB.RawQuery("SELECT COUNT(*) FROM post_tags WHERE tag_id = ?", tag.ID).Count(&count)
    tagsWithCount = append(tagsWithCount, TagWithCount{
        ID:        tag.ID,
        Name:      tag.Name,
        Slug:      tag.Slug,
        PostCount: count,
    })
}
```

#### 3. 传递数据给模板

**修改前**:
```go
c.Set("categories", categories)
c.Set("tags", tags)
```

**修改后**:
```go
c.Set("categories", categoriesWithCount)
c.Set("tags", tagsWithCount)
```

---

## 📝 模板使用示例

### home.plush.html

```html
<!-- Categories Widget -->
<div class="widget">
    <h4>📂 分类</h4>
    <ul class="category-list">
        <%= for (category) in categories { %>
        <li>
            <a href="/category/<%= category.Slug %>">
                <%= category.Name %> (<%= category.PostCount %>)  <!-- ✅ 正常工作 -->
            </a>
        </li>
        <% } %>
    </ul>
</div>
```

---

## 🎯 功能说明

### PostCount 计算逻辑

#### 分类文章数
```go
models.DB.Where("category_id = ? AND status = ?", cat.ID, "published").Count(&count)
```
- 只统计已发布的文章
- 根据分类 ID 过滤

#### 标签文章数
```go
models.DB.RawQuery("SELECT COUNT(*) FROM post_tags WHERE tag_id = ?", tag.ID).Count(&count)
```
- 通过 `post_tags` 关联表统计
- 包含所有状态的文章（可根据需求调整）

---

## 🔍 性能优化

### 当前实现
- 为每个分类/标签执行一次 COUNT 查询
- N+1 查询问题（N = 分类/标签数量）

### 优化方案（可选）

#### 方案 1: 使用 JOIN 一次查询

```go
// 一次查询获取所有分类及文章数
categoriesWithCount := []CategoryWithCount{}
models.DB.RawQuery(`
    SELECT c.id, c.name, c.slug, c.description, COUNT(p.id) as post_count
    FROM categories c
    LEFT JOIN posts p ON c.id = p.category_id AND p.status = 'published'
    GROUP BY c.id
    ORDER BY c.name
`).All(&categoriesWithCount)
```

#### 方案 2: 使用缓存

```go
// 使用内存缓存
var cachedCategories []CategoryWithCount
var lastCacheTime time.Time

func getCategoriesWithCount() []CategoryWithCount {
    // 缓存 5 分钟
    if time.Since(lastCacheTime) < 5*time.Minute {
        return cachedCategories
    }
    
    // 重新计算
    // ... 计算逻辑
    
    cachedCategories = result
    lastCacheTime = time.Now()
    return result
}
```

---

## 📊 数据结构对比

### 修改前
```go
type Category struct {
    ID          int
    Name        string
    Slug        string
    Description string
    // ❌ 没有 PostCount
}
```

### 修改后
```go
type CategoryWithCount struct {
    ID          int
    Name        string
    Slug        string
    Description string
    PostCount   int  // ✅ 添加文章数量
}
```

---

## ✅ 验证清单

- [x] 创建 CategoryWithCount 和 TagWithCount 结构体
- [x] 在控制器中计算文章数量
- [x] 传递带 PostCount 的数据给模板
- [x] 编译成功
- [x] 模板可以正常使用 category.PostCount 和 tag.PostCount

---

## 📂 修改文件

```
F:\Catter-CMS\
└── actions/
    └── front/
        └── home.go  ← 主要修改
```

**修改内容**:
- 添加 CategoryWithCount 结构体（第 40-46 行）
- 添加 TagWithCount 结构体（第 61-66 行）
- 计算分类文章数量（第 48-58 行）
- 计算标签文章数量（第 70-78 行）
- 更新模板数据传递（第 81 行）

---

## 🎯 扩展建议

### 1. 添加到其他页面

**分类页面** (`CategoryShowHandler`):
```go
// 获取所有分类及文章数
categoriesWithCount := getCategoriesWithCount()
c.Set("categories", categoriesWithCount)
```

**文章详情页** (`BlogShowHandler`):
```go
// 获取侧边栏分类
categoriesWithCount := getCategoriesWithCount()
c.Set("categories", categoriesWithCount)
```

### 2. 封装为函数

```go
// 封装为可复用的函数
func GetCategoriesWithCount() []CategoryWithCount {
    categories := models.Categories{}
    models.DB.All(&categories)
    
    result := []CategoryWithCount{}
    for _, cat := range categories {
        var count int
        models.DB.Where("category_id = ? AND status = ?", cat.ID, "published").Count(&count)
        result = append(result, CategoryWithCount{
            ID:          cat.ID,
            Name:        cat.Name,
            Slug:        cat.Slug,
            Description: cat.Description,
            PostCount:   count,
        })
    }
    return result
}

func GetTagsWithCount() []TagWithCount {
    // 类似实现
}
```

### 3. 添加排序功能

```go
// 按文章数量排序
sort.Slice(categoriesWithCount, func(i, j int) bool {
    return categoriesWithCount[i].PostCount > categoriesWithCount[j].PostCount
})
```

---

## 🎉 总结

成功添加了 `PostCount` 字段到分类和标签数据中，通过在控制器中动态计算的方式实现，无需修改数据库结构。这种方式：

- ✅ 数据始终准确
- ✅ 实现简单
- ✅ 不需要维护
- ✅ 不占用数据库空间

模板现在可以正常显示每个分类和标签下的文章数量了！

---

**修复时间**: 2026-03-16  
**影响范围**: 首页控制器、分类和标签显示  
**状态**: ✅ 已完成
