package main

import (
	"github.com/gin-gonic/gin"
	"log"
)

var Conf *Config
var etcdService = NewEtcdService(Conf)

func main() {
	var (
		err error
	)
	initFlag()
	Conf = initConfig()
	etcdService = NewEtcdService(Conf)

	if err = etcdService.connection(); err != nil {
		log.Fatalf("connection etcd error: %v", err)
	}

	route := gin.Default()

	route.HTMLRender = createHTMLRender()

	route.Static("/static", "./static")
	route.GET("/", index())
	route.GET("/machine/index", machineIndex())
	route.POST("/machine/store", machineStore())
	route.Run(":8001")
}
