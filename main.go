package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

var Conf *Config

func main() {
	var (
		err error
	)
	initFlag()
	Conf = initConfig()
	etcdService := NewEtcdService(Conf)
	if err = etcdService.connection(); err != nil {
		log.Fatalf("connection etcd error: %v", err)
	}

	route := gin.Default()

	route.HTMLRender = createHTMLRender()

	route.Static("/static", "./static")
	route.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index", nil)
	})

	route.GET("/machine/create", func(c *gin.Context) {
	})

	route.GET("/machine/index", func(c *gin.Context) {
		machines, err := etcdService.GetMachineList()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(machines)
		c.HTML(http.StatusOK, "machine/index", gin.H{
			"machines": machines,
		})
	})
	route.Run(":8001")
}
