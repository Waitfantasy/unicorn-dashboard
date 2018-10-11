package main

import (
	"strings"
	"time"
)

func formatMachineIP(ip string) string {
	return strings.TrimPrefix(ip, "/unicorn_machine/")
}

func formatDate(ts int64) string {
	return time.Unix(ts, 0).Format("2006-01-02 15:04:05")
}
