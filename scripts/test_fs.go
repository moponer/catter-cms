package main

import (
	"fmt"
	"io/fs"

	"catter-cms/templates"
)

func main() {
	fsys := templates.FS()
	
	// 列出根目录
	fmt.Println("=== Root directory ===")
	entries, err := fs.ReadDir(fsys, ".")
	if err != nil {
		fmt.Println("Error reading root:", err)
		return
	}
	for _, e := range entries {
		fmt.Printf("  %s (dir: %v)\n", e.Name(), e.IsDir())
	}

	// 列出 themes 目录
	fmt.Println("\n=== themes directory ===")
	entries2, err := fs.ReadDir(fsys, "themes")
	if err != nil {
		fmt.Println("Error reading themes:", err)
		return
	}
	for _, e := range entries2 {
		fmt.Printf("  %s (dir: %v)\n", e.Name(), e.IsDir())
	}

	// 列出 catter-blog/pages 目录
	fmt.Println("\n=== catter-blog/pages directory ===")
	entries3, err := fs.ReadDir(fsys, "themes/catter-blog/pages")
	if err != nil {
		fmt.Println("Error reading pages:", err)
		return
	}
	for _, e := range entries3 {
		fmt.Printf("  %s (dir: %v)\n", e.Name(), e.IsDir())
	}
}
