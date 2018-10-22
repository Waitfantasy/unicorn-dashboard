package main

import (
	"github.com/Waitfantasy/unicorn/service/machine"
	"github.com/gin-gonic/gin"
	"log"
)

var Conf *Config
var machineService machine.Machiner
func main() {
	var (
		err error
	)
	initFlag()
	Conf = initConfig()

	factory := machine.MachineFactory{}

	machineService, err = factory.CreateEtcdMachine(CreateEtcdV3Client(Conf))

	if err != nil {
		log.Fatal(err)
	}



	route := gin.Default()

	api := route.Group("/api/v1")
	api.GET("/machine/list", machineIndex())
	api.POST("/machine/store", machineStore())
	api.POST("/machine/delete", machineDelete())
	route.Run(":8001")
}
