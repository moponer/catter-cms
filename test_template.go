// +build ignore

package main

import (
	"fmt"

	"catter-cms/templates"
)

func main() {
	fs := templates.FS()
	
	// 尝试打开模板文件
	file, err := fs.Open("themes/blog/pages/home.plush.html")
	if err != nil {
		fmt.Println("Error opening template:", err)
		return
	}
	defer file.Close()
	
	fmt.Println("Template found successfully!")
	
	// 列出 themes 目录
	entries, err := fs.(interface{ ReadDir(name string) ([]fs.DirEntry, error) }).ReadDir("themes")
	if err != nil {
		fmt.Println("Error reading themes dir:", err)
		return
	}
	
	fmt.Println("Themes available:")
	for _, e := range entries {
		fmt.Println(" -", e.Name())
	}
}
