package cron_starter

import (
	"github.com/spf13/cast"
)

type StarterOptions struct {
	Id            string `json:"id"`
	NodeId        string `json:"node_id"`
	NodeType      string `json:"node_type"`
	HashringSpots int    `json:"hashring_spots"`

	Lb               string `json:"lb"`
	LbRedis          string `json:"lb_redis"`
	LbRedisWeight    string `json:"lb_redis_weight"`
	LbRedisTimeout   int    `json:"lb_redis_timeout"`
	LbRedisHeartbeat int    `json:"lb_redis_heartbeat"`
	LbRedisPrefix    string `json:"lb_redis_prefix"`
	LbRedisChannel   string `json:"lb_redis_channel"`

	WorkerFeign              string `json:"worker_feign"`
	WorkerFeignTrace         string `json:"worker_feign_trace"`
	WorkerFeignDebug         string `json:"worker_feign_debug"`
	WorkerFeignHost          string `json:"worker_feign_host"`
	WorkerFeignTimeout       int    `json:"worker_feign_timeout"`
	WorkerFeignRetryCount    int    `json:"worker_feign_retry_count"`
	WorkerFeignRetryWaitTime int    `json:"worker_feign_retry_wait_time"`

	WorkerHeartbeatSpec string `json:"worker_heartbeat_spec"`
	WorkerHeartbeatUrl  string `json:"worker_heartbeat_url"`
}

func NewStarterOptions(id string, cfg map[string]string) *StarterOptions {
	options := StarterOptions{Id: id}
	options.LoadConfig(cfg)
	return &options
}

func (options *StarterOptions) LoadConfig(cfg map[string]string) {
	options.NodeId = cfg["node_id"]
	options.NodeType = cfg["node_type"]
	if cfg["hashring_spots"] == "" {
		options.HashringSpots = 5
	} else {
		options.HashringSpots = cast.ToInt(cfg["hashring_spots"])
	}

	options.Lb = cfg["lb"]
	options.LbRedis = cfg["lb_redis"]
	if cfg["lb_redis_timeout"] == "" {
		options.LbRedisTimeout = 30
	} else {
		options.LbRedisTimeout = cast.ToInt(cfg["lb_redis_timeout"])
	}

	if cfg["lb_redis_heartbeat"] == "" {
		options.LbRedisHeartbeat = 300
	} else {
		options.LbRedisHeartbeat = cast.ToInt(cfg["lb_redis_heartbeat"])
	}

	options.LbRedisPrefix = cast.ToString(cfg["lb_redis_prefix"])
	options.LbRedisChannel = cast.ToString(cfg["lb_redis_channel"])
	options.LbRedisWeight = cfg["lb_redis_weight"]

	options.WorkerHeartbeatUrl = cfg["worker_heartbeat_uri"]
	options.WorkerHeartbeatSpec = cfg["worker_heartbeat_spec"]
	options.WorkerFeign = cfg["worker_feign"]
	options.WorkerFeignTrace = cfg["worker_feign_trace"]
	options.WorkerFeignDebug = cfg["worker_feign_debug"]
	options.WorkerFeignHost = cfg["worker_feign_host"]
	options.WorkerFeignTimeout = cast.ToInt(cfg["worker_feign_timeout"])
	options.WorkerFeignRetryCount = cast.ToInt(cfg["worker_feign_retry_count"])

}
