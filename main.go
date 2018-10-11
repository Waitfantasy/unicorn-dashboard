package main

import (
	"github.com/Waitfantasy/unicorn-service/machine"
	"github.com/gin-gonic/gin"
	"log"
)

var Conf *Config
var MachineService *machine.Service

func main() {
	var (
		err error
	)
	initFlag()
	Conf = initConfig()

	MachineService = machine.NewService(createEtcdClientv3Config(Conf))

	if err = MachineService.EtcdConnection(); err != nil {
		log.Fatalf("connection etcd error: %v", err)
	}

	route := gin.Default()

	route.HTMLRender = createHTMLRender()

	route.Static("/static", "./static")
	route.GET("/", index())
	route.GET("/machine/index", machineIndex())
	route.POST("/machine/store", machineStore())
	route.POST("/machine/delete", machineDelete())
	route.Run(":8001")
}
