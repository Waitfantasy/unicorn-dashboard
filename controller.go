package main

import (
	"errors"
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
		items, err := machineService.All()
		if err != nil {
			c.JSON(http.StatusOK, map[string]interface{}{
				"success": false,
				"message": "get all machine item fail.",
				"data": map[string]interface{}{
					"error": err.Error(),
				},
			})
		} else {
			c.JSON(http.StatusOK, map[string]interface{}{
				"success": true,
				"message": "get all machine item success.",
				"data": map[string]interface{}{
					"machines": items,
				},
			})
		}
	}
}

func machineStore() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip, err := validatorIp(c)
		if err != nil {
			c.JSON(http.StatusOK, map[string]interface{}{
				"success": false,
				"message": err.Error(),
			})
		}

		if item, err := machineService.Put(ip); err != nil {
			c.JSON(http.StatusOK, map[string]interface{}{
				"success": false,
				"message": fmt.Sprintf("put ip: %s fail.\n", ip),
				"data": map[string]interface{}{
					"error": err.Error(),
				},
			})
		} else {
			c.JSON(http.StatusOK, map[string]interface{}{
				"success": true,
				"message": fmt.Sprintf("put ip: %s success.\n", ip),
				"data": map[string]interface{}{
					"machine": item,
				},
			})
		}
	}
}

func machineDelete() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip, err := validatorIp(c)
		if err != nil {
			c.JSON(http.StatusOK, map[string]interface{}{
				"success": false,
				"message": err.Error(),
			})
		}

		if item, err := machineService.Del(ip); err != nil {
			c.JSON(http.StatusOK, map[string]interface{}{
				"success": false,
				"message": fmt.Sprintf("del ip: %s fail.\n", ip),
				"data": map[string]interface{}{
					"error": err.Error(),
				},
			})
		} else {
			c.JSON(http.StatusOK, map[string]interface{}{
				"success": true,
				"message": fmt.Sprintf("put ip: %s success.\n", ip),
				"data": map[string]interface{}{
					"machine": item,
				},
			})
		}
	}
}

func validatorIp(c *gin.Context) (string, error) {
	if ip := c.PostForm("ip"); ip == "" {
		return "", errors.New("missing parameters: ip")
	} else {
		return ip, nil
	}
}
