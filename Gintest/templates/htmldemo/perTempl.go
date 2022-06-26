package main

import (
	"github.com/gin-gonic/gin"
	"html/template"
)

func main() {
	router := gin.Default()
	//自定义模板渲染器
	html := template.Must(template.ParseFiles("file1", "file2"))
	router.SetHTMLTemplate(html)
	//自定义分隔符
	router.Delims("{[{", "}]}")
	router.LoadHTMLGlob("/path/to/templates")
	router.Run(":8080")
}
