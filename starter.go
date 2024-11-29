package cron_starter

import (
	"github.com/go-resty/resty/v2"
	"github.com/kordar/gocron"
	goframeworkcron "github.com/kordar/goframework-cron"
	"github.com/kordar/goframework-resty"
	logger "github.com/kordar/gologger"
	"github.com/kordar/goresty"
	"github.com/spf13/cast"
	"time"
)

type InitializeFunction func(job gocron.Schedule) map[string]string
type RuntimeFunction func(job gocron.Schedule) bool

var (
	LBHandle           = map[string]LB{}
	RemoteWorkerHandle = map[string]*RemoteWorker{}
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

func (m CronModule) runFn() (InitializeFunction, RuntimeFunction) {
	var initializeFn InitializeFunction
	if m.args != nil && m.args["cfgFn"] != nil {
		initializeFn = m.args["cfgFn"].(InitializeFunction)
	} else {
		initializeFn = func(job gocron.Schedule) map[string]string {
			return map[string]string{
				"spec": "@every 10s",
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

func (m CronModule) loadLb(id string, cfg map[string]string) {
	nid := cfg["node_id"]
	if nid == "" {
		logger.Fatalf("[%s] to enable load balancing, you must configure the parameter \"node_id\"", m.Name())
	}

	if cfg["lb"] == "redis" {
		timeout := cast.ToInt(cfg["lb_redis_timeout"])
		heartbeat := cast.ToInt(cfg["lb_redis_heartbeat"])
		prefix, channel, weight := cfg["lb_redis_prefix"], cfg["lb_redis_channel"], cfg["lb_redis_weight"]
		LBHandle[id] = NewRedisNode(nid, prefix, channel)
		StartRedisRegistry(cfg["lb_redis"], prefix, nid, timeout, heartbeat, weight, channel, LBHandle[id])
	}
}

func (m CronModule) loadWorker(id string, cfg map[string]string) {

	nid := cfg["node_id"]
	if nid == "" {
		logger.Fatalf("[%s] you must configure the parameter \"node_id\"", m.Name())
	}

	host := cfg["remote_host"]
	if host == "" {
		logger.Fatalf("[%s] you must configure the parameter \"remote_host\"", m.Name())
	}

	if cfg["remote_feign"] != "" {
		client := goframework_resty.GetFeignClient(cfg["remote_feign"])
		RemoteWorkerHandle[id] = NewRemoteWorker(id, client)
		return
	}

	feign := goresty.NewFeign(nil).OnBeforeRequest(func(client *resty.Client, request *resty.Request) error {
		if cfg["remote_trace"] == "enable" {
			request.EnableTrace()
		}
		//
		// if cfg["remote_secret_key"] != "" {
		//
		// }
		//
		return nil
	}).Options(func(client *resty.Client) {
		client.SetBaseURL(cfg["remote_host"])
		if cfg["remote_debug"] == "enable" {
			client.SetDebug(true)
		}
		if cfg["remote_timeout"] != "" {
			remoteTimeout := cast.ToDuration(cfg["remote_timeout"])
			client.SetTimeout(time.Second * remoteTimeout)
		}
		if cfg["remote_retry_count"] != "" {
			remoteRetryCount := cast.ToInt(cfg["remote_retry_count"])
			client.SetRetryCount(remoteRetryCount)
		}
		if cfg["remote_retry_wait_time"] != "" {
			remoteRetryWaitTime := cast.ToDuration(cfg["remote_retry_wait_time"])
			client.SetRetryWaitTime(time.Second * remoteRetryWaitTime)
		}
	}).OnError(func(request *resty.Request, err error) {
		logger.Errorf("[%s] request err = %+v", m.Name(), err)
	})

	RemoteWorkerHandle[id] = NewRemoteWorker(id, feign)
}

func (m CronModule) _load(id string, cfg map[string]string) {
	if id == "" {
		logger.Fatalf("[%s] the attribute id cannot be empty.", m.Name())
		return
	}

	if cfg["lb"] != "" {
		m.loadLb(id, cfg)
	}

	fn1, fn2 := m.runFn()
	err := goframeworkcron.AddGocronInstance(id, func(job gocron.Schedule) map[string]string {
		return fn1(job)
	}, func(job gocron.Schedule) bool {
		if LBHandle[id+":"+job.GetId()] != nil {
			return LBHandle[id].Can()
		}
		return fn2(job)
	})

	// 配置节点类型为work时，生成rest请求对象，并插入提交定时任务进行心跳维护。
	if cfg["node_type"] == "worker" {
		m.loadWorker(id, cfg)
		goframeworkcron.AddJob(id, NewWorkerHeartSchedule(id, cfg["remote_worker_heart_spec"], cfg["remote_worker_heart_url"]))
		//goframeworkcron.AddJob(id, NewWorkerHeartSchedule(id, cfg["remote_worker_heart_spec"]))
	}

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
