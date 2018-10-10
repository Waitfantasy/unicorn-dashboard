package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func index() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "index", nil)
	}
}

func machineIndex() gin.HandlerFunc {
	return func(c *gin.Context) {
		machines, err := etcdService.GetMachineList()
		if err != nil {
			c.HTML(http.StatusOK, "machine/index", gin.H{
				"message": map[string]string{
					"level": "danger",
					"title": "Error",
					"body":  fmt.Sprintf("Get machine list error: %v", err.Error()),
				},
			})
			return
		}

		fmt.Println(machines)

		c.HTML(http.StatusOK, "machine/index", gin.H{
			"machines": machines,
		})
	}
}

func machineStore() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.PostForm("id")
		ip := c.PostForm("ip")
		if data, err := etcdService.PutMachineId(ip, id); err != nil {
			c.JSON(http.StatusOK, map[string]interface{}{
				"success": false,
				"code":    0,
				"message": err.Error(),
			})
		} else {
			c.JSON(http.StatusOK, map[string]interface{}{
				"success": false,
				"code":    1,
				"message": "put etcd success.",
				"data": map[string]string{
					"id":                id,
					"ip":                ip,
					"created_timestamp": formatDate(data["created_timestamp"]),
					"updated_timestamp": formatDate(data["updated_timestamp"]),
				},
			})
		}
	}
}
