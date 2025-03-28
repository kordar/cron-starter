package cron_starter_test

import (
	cron_starter "github.com/kordar/cron-starter"
	"github.com/kordar/gocron"
	goframeworkcron "github.com/kordar/goframework-cron"
	logger "github.com/kordar/gologger"
	"testing"
	"time"
)

var initializeFn gocron.InitializeFunction = func(job gocron.Schedule) map[string]string {
	cfg := map[string]string{}
	cfg["spec"] = "@every 10s"
	logger.Info("xxxxxxxxxxx", job.GetId())
	if job.GetId() == "AAA" {
		cfg["spec"] = "@every 5s"
	}
	return cfg
}

var handle = cron_starter.NewCronModule("AA", nil, map[string]interface{}{})

var cfg = map[string]interface{}{
	"AAA": map[string]interface{}{
		"node_id":           "xxx",
		"node_type":         "worker",
		"remote":            "worker",
		"worker_feign_host": "https://www.baidu.com",
	},
	"BBB": map[string]interface{}{
		"id":          "BBB",
		"remote":      "worker",
		"remote_host": "https://www.sina.com",
	},
}

type TestNameSchedule struct {
	gocron.BaseSchedule
}

func (s TestNameSchedule) GetId() string {
	return "test-name"
}

func (s TestNameSchedule) GetSpec() string {
	return "@every 5s"
}

func (s TestNameSchedule) Execute() {
	config := s.Config()
	logger.Infof("--------------test name--------------%v", config)
}

func TestNewCronModule(t *testing.T) {

	handle.Load(cfg)

	s := &TestNameSchedule{}
	s.SetConfig(map[string]string{
		"spec": "@every 5s",
	})

	goframeworkcron.AddJob("BBB", s)

	time.Sleep(100 * time.Second)
}
