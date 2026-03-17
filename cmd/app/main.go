package main

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"

	"catter-cms/actions"
)

func main() {
	// 捕获 panic
	defer func() {
		if r := recover(); r != nil {
			log.Printf("PANIC: %v\nStack trace:\n%s", r, debug.Stack())
			os.Exit(1)
		}
	}()

	// 禁用 Buffalo 插件加载，避免 Windows 上的插件加载错误
	os.Setenv("BUFFALO_PLUGIN_PATH", "")

	fmt.Println("Starting application...")
	fmt.Println("Getting App instance...")
	
	app := actions.App
	if app == nil {
		log.Fatal("App is nil!")
	}
	
	fmt.Printf("App initialized, Host: %s\n", app.Host)
	fmt.Println("Starting server...")
	
	if err := app.Serve(); err != nil {
		log.Fatal("Serve error:", err)
	}
}

/*
# Notes about `main.go`

## SSL Support

We recommend placing your application behind a proxy, such as
Apache or Nginx and letting them do the SSL heavy lifting
for you. https://gobuffalo.io/en/docs/proxy

## Buffalo Build

When `buffalo build` is run to compile your binary, this `main`
function will be at the heart of that binary. It is expected
that your `main` function will start your application using
the `app.Serve()` method.

*/
