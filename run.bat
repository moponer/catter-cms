@echo off
chcp 65001 >nul
echo ================================================
echo   Catter CMS 启动脚本
echo ================================================
echo.

:: 设置 Go 环境
set GO_ENV=development
set APP_ENV=development

:: 启动应用
:: 从 gin-app 目录运行时
if exist "gin-app\cmd\main.go" (
    go run ./gin-app/cmd/main.go
) else if exist "cmd\main.go" (
    go run ./cmd/main.go
) else (
    echo [错误] 无法找到 main.go 文件
    echo 请确保在项目根目录运行此脚本
    pause
    exit /b 1
)

pause
