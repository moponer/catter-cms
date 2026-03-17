#!/bin/bash

echo "=== Catter CMS 全面测试验证 ==="
echo ""

# 1. 测试数据库连接
echo "1. 测试数据库连接..."
cd f:\Catter-CMS
go run -tags development -run "models" - <<'EOF'
package main
import (
    "fmt"
    "catter-cms/models"
)
func main() {
    if err := models.DB.Raw("SELECT COUNT(*) as count FROM users").Exec(); err != nil {
        fmt.Printf("数据库连接失败: %v\n", err)
    } else {
        fmt.Println("数据库连接成功")
    }
}
EOF

# 2. 测试API接口
echo ""
echo "2. 测试API接口..."
echo "2.1 测试文章列表API..."
curl -s -X GET http://localhost:3000/api/posts | head -c 100

echo "2.2 测试分类列表API..."
curl -s -X GET http://localhost:3000/api/categories | head -c 100

# 3. 测试前端页面
echo ""
echo "3. 测试前端页面..."
echo "3.1 测试首页..."
curl -s -X GET http://localhost:3000/ | head -c 50

echo "3.2 测试登录页面..."
curl -s -X GET http://localhost:3000/admin/login | head -c 50

# 4. 测试登录功能
echo ""
echo "4. 测试登录功能..."
echo "4.1 获取登录页面和CSRF令牌..."
CSRF_TOKEN=$(curl -s -X GET http://localhost:3000/admin/login | grep -oP 'csrf-token" content="\K[^"]+')
echo "CSRF令牌: $CSRF_TOKEN"

echo "4.2 提交登录表单..."
curl -s -X POST http://localhost:3000/admin/login/ \
  -d "email=admin@example.com" \
  -d "password=admin123" \
  -d "authenticity_token=$CSRF_TOKEN" \
  -L | head -c 100

echo ""
echo "=== 测试完成 ==="
