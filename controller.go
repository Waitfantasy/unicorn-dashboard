package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func index() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "index", nil)
	}
}

func machineIndex() gin.HandlerFunc {
	return func(c *gin.Context) {
		items, err := MachineService.GetMachineItemList()
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

		c.HTML(http.StatusOK, "machine/index", gin.H{
			"machinesItems": items,
		})
	}
}

func machineStore() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.PostForm("ip")
		if item, err := MachineService.PutMachineItem(ip); err != nil {
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
					"id":                strconv.Itoa(item.Id),
					"ip":                item.Ip,
					"created_timestamp": item.FormatCreatedTime(),
					"updated_timestamp": item.FormatUpdatedTime(),
				},
			})
		}
	}
}

func machineDelete() gin.HandlerFunc  {
	return func(c *gin.Context) {
		ip := c.PostForm("ip")
		if item, err := MachineService.DelMachineItem(ip); err != nil {
			c.JSON(http.StatusOK, map[string]interface{}{
				"success": false,
				"code":    0,
				"message": err.Error(),
			})
		} else {
			c.JSON(http.StatusOK, map[string]interface{}{
				"success": true,
				"code":    1,
				"message": "delete machine success.",
				"data": map[string]string{
					"id":                strconv.Itoa(item.Id),
					"ip":                item.Ip,
					"created_timestamp": item.FormatCreatedTime(),
					"updated_timestamp": item.FormatUpdatedTime(),
				},
			})
		}
	}
}
