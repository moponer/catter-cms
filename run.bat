@echo off
chcp 65001 >nul
echo Starting Catter CMS...
echo.

:: 禁用 Buffalo 插件加载，避免 "not a valid Win32 application" 错误
set BUFFALO_PLUGIN_PATH=
set BUFFALO_PLUGINS_OFF=1

:: 设置 Go 环境
set GO_ENV=development

:: 运行应用
go run ./cmd/app/main.go

pause
