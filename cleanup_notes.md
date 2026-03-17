# 代码清理笔记

## 需要清理的硬编码配置

### 1. actions/app.go
- **硬编码的ENV变量**：
  ```go
  var ENV = envy.Get("GO_ENV", "development")
  ```
  应该改为：
  ```go
  import "catter-cms/config"
  var ENV = config.GetEnvironment()
  ```

- **硬编码的SessionName**：
  ```go
  SessionName: "_catter_cms_session"
  ```
  应该改为：
  ```go
  SessionName: config.GetSessionSecret()
  ```

- **错误的中间件引用**：
  ```go
  authAdminGroup.Use(middleware.Auth)
  ```
  应该改为：
  ```go
  authAdminGroup.Use(middleware.RequireAuth)
  ```

### 2. actions/admin/login.go
- **未使用的render导入**：
  如果没有使用render包的功能，可以移除

### 3. actions/api/posts.go
- **硬编码的响应格式**：
  应该使用统一的DTO响应格式

### 4. models/models.go
- **硬编码的数据库连接**：
  应该通过配置管理

## 需要清理的未使用导入

### 1. 检查所有Go文件中的导入
- 使用`goimports`工具自动清理未使用的导入
- 使用`golangci-lint`工具检查代码质量

### 2. 移除注释掉的代码
- 搜索并移除所有注释掉的代码块

### 3. 移除调试代码
- 搜索并移除fmt.Println等调试语句

## 建议的清理步骤

1. 安装代码清理工具：
   ```bash
   go install golang.org/x/tools/cmd/goimports@latest
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
   ```

2. 运行导入清理：
   ```bash
   goimports -w ./...
   ```

3. 运行代码检查：
   ```bash
   golangci-lint run
   ```

4. 手动清理硬编码配置：
   - 将硬编码的值移到config包
   - 使用config.GetXXX()函数获取配置

5. 移除未使用的文件：
   - 检查是否有未使用的测试文件
   - 检查是否有重复的文件

## 清理后的预期效果

1. **代码质量提升**：
   - 减少硬编码，提高可维护性
   - 移除未使用的导入，减少包大小
   - 提高代码一致性

2. **可维护性提升**：
   - 配置集中管理
   - 代码更加清晰
   - 减少配置错误

3. **性能提升**：
   - 减少不必要的包导入
   - 减少编译时间
   - 减少二进制文件大小
