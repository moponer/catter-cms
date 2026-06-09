package main

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
)

func main() {
	themePath := "templates/themes/enterprise-classic"

	var templateFiles []string

	// modules
	modulesPath := filepath.Join(themePath, "modules")
	if info, err := os.Stat(modulesPath); err == nil && info.IsDir() {
		moduleFiles, _ := filepath.Glob(filepath.Join(modulesPath, "*.html"))
		templateFiles = append(templateFiles, moduleFiles...)
	}

	// pages
	pagesPath := filepath.Join(themePath, "pages")
	if info, err := os.Stat(pagesPath); err == nil && info.IsDir() {
		pageFiles, _ := filepath.Glob(filepath.Join(pagesPath, "*.html"))
		templateFiles = append(templateFiles, pageFiles...)
	}

	fmt.Printf("Found %d template files:\n", len(templateFiles))
	for _, f := range templateFiles {
		fmt.Printf("  %s\n", f)
	}

	tmpl := template.New("")
	parsedTmpl, err := tmpl.ParseFiles(templateFiles...)
	if err != nil {
		fmt.Printf("Parse error: %v\n", err)
		return
	}

	fmt.Printf("\nRegistered template names:\n")
	for _, t := range parsedTmpl.Templates() {
		fmt.Printf("  %q\n", t.Name())
	}

	// 测试 ExecuteTemplate("home.html")
	fmt.Printf("\nTesting ExecuteTemplate(\"home.html\"):\n")
	var buf1 [4096]byte
	_ = buf1
	// 用 Root 执行
	tmpl2 := parsedTmpl
	// 查找 home 相关
	for _, t := range tmpl2.Templates() {
		if filepath.Base(t.Name()) == "home.html" || t.Name() == "home.html" || t.Name() == "home" {
			fmt.Printf("  Found matching template: %q\n", t.Name())
		}
	}
}
