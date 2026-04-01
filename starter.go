package cron_starter

import (
	"log/slog"

	"github.com/kordar/gocron"
	goframeworkcron "github.com/kordar/goframework-cron"
	"github.com/spf13/cast"
)

type CronModule struct {
	name string
	args map[string]interface{}
	load func(moduleName string, itemId string, item map[string]string)
}

func NewCronModule(name string, load func(moduleName string, itemId string, item map[string]string), args map[string]interface{}) *CronModule {
	return &CronModule{name, args, load}
}

func (m CronModule) Name() string {
	return m.name
}

func (m CronModule) runFn() (gocron.InitializeFunction, gocron.RuntimeFunction) {
	var initializeFn gocron.InitializeFunction
	if m.args != nil && m.args["cfgFn"] != nil {
		initializeFn = m.args["cfgFn"].(gocron.InitializeFunction)

	} else {
		initializeFn = func(job gocron.Schedule) map[string]string {
			return map[string]string{
				"spec": "@every 10m",
			}
		}
	}

	var runtimeFn gocron.RuntimeFunction
	if m.args != nil && m.args["validateFn"] != nil {
		runtimeFn = m.args["validateFn"].(gocron.RuntimeFunction)
	} else {
		runtimeFn = func(job gocron.Schedule) bool {
			return true
		}
	}

	return initializeFn, runtimeFn
}

func (m CronModule) _load(id string, cfg map[string]string) {

	if id == "" {
		slog.Error("the attribute id cannot be empty", "module", m.Name())
		return
	}

	fn1, fn2 := m.runFn()
	err := goframeworkcron.AddGocronInstance(id, func(job gocron.Schedule) map[string]string {
		return fn1(job)
	}, func(job gocron.Schedule) bool {
		return fn2(job)
	})

	if err != nil {
		slog.Error("failed to initialize gocron instance", "module", m.Name())
	}

	if m.load != nil {
		m.load(m.name, id, cfg)
		slog.Debug("triggering custom loader completion", "module", m.Name())
	}

	slog.Info("loading module successfully", "module", m.Name(), "id", id)
}

func (m CronModule) Load(value interface{}) {
	items := cast.ToStringMap(value)
	if items["id"] != nil {
		id := cast.ToString(items["id"])
		m._load(id, cast.ToStringMapString(value))
		return
	}

	for key, item := range items {
		m._load(key, cast.ToStringMapString(item))
	}
}

func (m CronModule) Close() {
}
