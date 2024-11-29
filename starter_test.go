package cron_starter_test

import (
	cron_starter "github.com/kordar/cron-starter"
	"testing"
	"time"
)

var handle = cron_starter.NewCronModule("AA", nil, map[string]interface{}{})

var cfg = map[string]interface{}{
	"AAA": map[string]interface{}{
		"id":          "xxx",
		"remote":      "worker",
		"remote_host": "https://www.baidu.com",
	},
	"BBB": map[string]interface{}{
		"id":          "BBB",
		"remote":      "worker",
		"remote_host": "https://www.sina.com",
	},
}

func TestNewCronModule(t *testing.T) {
	handle.Load(cfg)

	time.Sleep(100 * time.Second)
}
