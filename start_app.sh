#!/bin/bash

# 停止占用端口3000的进程
echo "停止占用端口3000的进程..."
netstat -ano | findstr :3000 | awk '{print $2}' | xargs -I {} powershell -Command "Stop-Process -Id {} -Force" 2>/dev/null || true

# 等待端口释放
sleep 2

# 启动应用程序
echo "启动应用程序..."
cd f:\Catter-CMS
set CGO_ENABLED=1 && go run cmd/app/main.go
