package main

import (
	"strconv"
	"strings"
	"time"
)

func formatMachineIP(ip string) string {
	return strings.TrimPrefix(ip, "/unicorn_machine/")
}

func formatDate(ts string) string {
	timestamp, err := strconv.Atoi(ts)
	if err != nil {
		return "timestamp format error"
	}

	return time.Unix(int64(timestamp), 0).Format("2006-01-02 15:04:05")
}
