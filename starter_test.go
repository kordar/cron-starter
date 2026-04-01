package cron_starter_test

import (
	"testing"
	"time"
	"log/slog"

	cron_starter "github.com/kordar/cron-starter"
	"github.com/kordar/gocron"
	goframeworkcron "github.com/kordar/goframework-cron"
)

var initializeFn gocron.InitializeFunction = func(job gocron.Schedule) map[string]string {
	cfg := map[string]string{}
	cfg["spec"] = "@every 10s"
	slog.Info("xxxxxxxxxxx", "id", job.GetId())
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
	slog.Info("--------------test name--------------", "config", config)
}

func TestNewCronModule(t *testing.T) {

	handle.Load(cfg)

	s := &TestNameSchedule{}
	s.SetConfig(map[string]string{
		"spec": "@every 5s",
	})

	// goframeworkcron.AddJob("BBB", s)
	a := goframeworkcron.GetCronClient("BBB")
	a.Add(s)
	defer func() {
		// goframeworkcron.RemoveAllJob("BBB")
		// goframeworkcron.Stop("BBB")
	}()

	time.Sleep(10 * time.Second)
}
