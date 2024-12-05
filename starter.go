package cron_starter

import (
	"github.com/kordar/gocron"
	goframeworkcron "github.com/kordar/goframework-cron"
	logger "github.com/kordar/gologger"
	"github.com/spf13/cast"
)

type InitializeFunction func(job gocron.Schedule) map[string]string
type RuntimeFunction func(job gocron.Schedule) bool

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

func (m CronModule) runFn() (InitializeFunction, RuntimeFunction) {
	var initializeFn InitializeFunction
	if m.args != nil && m.args["cfgFn"] != nil {
		initializeFn = m.args["cfgFn"].(InitializeFunction)
	} else {
		initializeFn = func(job gocron.Schedule) map[string]string {
			return map[string]string{
				"spec": "@every 10m",
			}
		}
	}

	var runtimeFn RuntimeFunction
	if m.args != nil && m.args["validateFn"] != nil {
		runtimeFn = m.args["validateFn"].(RuntimeFunction)
	} else {
		runtimeFn = func(job gocron.Schedule) bool {
			return true
		}
	}

	return initializeFn, runtimeFn
}

func (m CronModule) _load(id string, cfg map[string]string) {

	if id == "" {
		logger.Fatalf("[%s] the attribute id cannot be empty.", m.Name())
		return
	}

	starterOptions := NewStarterOptions(id, cfg)
	lb := loadLb(m.Name(), starterOptions)

	fn1, fn2 := m.runFn()
	err := goframeworkcron.AddGocronInstance(id, func(job gocron.Schedule) map[string]string {
		return fn1(job)
	}, func(job gocron.Schedule) bool {
		if lb != nil {
			return lb.Can(job.GetId())
		}
		return fn2(job)
	})

	// 配置节点类型为work时，生成rest请求对象，并插入提交定时任务进行心跳维护。
	loadWorker(m.Name(), starterOptions)

	if err != nil {
		logger.Fatalf("[%s] failed to initialize gocron instance.", m.Name())
	}

	if m.load != nil {
		m.load(m.name, id, cfg)
		logger.Debugf("[%s] triggering custom loader completion", m.Name())
	}

	logger.Infof("[%s] loading module '%s' successfully", m.Name(), id)
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
