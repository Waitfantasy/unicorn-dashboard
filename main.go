package main

import (
	"go.etcd.io/etcd/clientv3"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
)

func main() {
	route := gin.Default()
	route.SetHTMLTemplate(template.Must(template.ParseFiles(
		"./templates/base.html",
		"./templates/layout/nav.html",
		"./templates/layout/sidebar.html",
		"./templates/layout/content-header.html",
		"./templates/layout/content.html",
		"./templates/layout/footer.html")))
	route.Static("/static", "./static")
	route.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "base.html", nil)
	})
	route.Run(":8001")
}
