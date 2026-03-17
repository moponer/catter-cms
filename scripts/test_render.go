package main

import (
	"log"

	"catter-cms/templates"
	"catter-cms/public"

	"github.com/gobuffalo/buffalo/render"
)

func main() {
	R := render.New(render.Options{
		HTMLLayout:   "application.plush.html",
		TemplatesFS:  templates.FS(),
		AssetsFS:     public.FS(),
	})

	// 测试渲染
	f := R.HTML("themes/catter-blog/pages/home_simple.plush.html")

	// 尝试渲染
	data := map[string]interface{}{
		"posts":     []string{"Test1", "Test2"},
		"categories": []string{"Cat1", "Cat2"},
		"tags":       []string{"Tag1", "Tag2"},
		"site":       map[string]string{"title": "Test", "description": "Test desc"},
		"theme":      "catter-blog",
	}

	err := f.Render(log.Writer(), data)
	if err != nil {
		log.Fatal("Error rendering:", err)
	}
	
	log.Println("Success!")
}
